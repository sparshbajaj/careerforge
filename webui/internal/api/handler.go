package api

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Application represents a single job application from the tracker.
type Application struct {
	Number       int     `json:"number"`
	Date         string  `json:"date"`
	Company      string  `json:"company"`
	Role         string  `json:"role"`
	ScoreRaw     string  `json:"scoreRaw"`
	Score        float64 `json:"score"`
	Status       string  `json:"status"`
	StatusNorm   string  `json:"statusNorm"`
	HasPDF       bool    `json:"hasPdf"`
	ReportPath   string  `json:"reportPath"`
	ReportNumber string  `json:"reportNumber"`
	Notes        string  `json:"notes"`
	JobURL       string  `json:"jobUrl"`
}

// Metrics holds aggregate stats.
type Metrics struct {
	Total      int            `json:"total"`
	ByStatus   map[string]int `json:"byStatus"`
	AvgScore   float64        `json:"avgScore"`
	TopScore   float64        `json:"topScore"`
	WithPDF    int            `json:"withPdf"`
	Actionable int            `json:"actionable"`
}

// StatusUpdate is the request body for updating a status.
type StatusUpdate struct {
	Status string `json:"status"`
}

// NotesUpdate is the request body for updating notes.
type NotesUpdate struct {
	Notes string `json:"notes"`
}

var (
	reReportLink = regexp.MustCompile(`\[(\d+)\]\(([^)]+)\)`)
	reScoreValue = regexp.MustCompile(`(\d+\.?\d*)/5`)
	reReportURL  = regexp.MustCompile(`(?m)^\*\*URL:\*\*\s*(https?://\S+)`)
)

// Valid statuses
var validStatuses = map[string]bool{
	"Evaluated": true, "Applied": true, "Responded": true,
	"Interview": true, "Offer": true, "Rejected": true,
	"Discarded": true, "SKIP": true,
}

// Handler serves the Web UI API and static files.
type Handler struct {
	careerOpsPath string
	staticFS      embed.FS
	mu            sync.RWMutex
}

// NewHandler creates a new API handler.
func NewHandler(careerOpsPath string, staticFS embed.FS) http.Handler {
	h := &Handler{
		careerOpsPath: careerOpsPath,
		staticFS:      staticFS,
	}

	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("GET /api/applications", h.getApplications)
	mux.HandleFunc("PUT /api/applications/{number}/status", h.updateStatus)
	mux.HandleFunc("PUT /api/applications/{number}/notes", h.updateNotes)
	mux.HandleFunc("GET /api/metrics", h.getMetrics)
	mux.HandleFunc("GET /api/reports/{path...}", h.getReport)
	mux.HandleFunc("GET /api/profile", h.getProfile)

	// Static files
	staticSub, _ := fs.Sub(staticFS, "static")
	mux.Handle("GET /", http.FileServer(http.FS(staticSub)))

	return mux
}

func (h *Handler) getApplications(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	apps := h.parseApplications()
	writeJSON(w, apps)
}

func (h *Handler) getMetrics(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	apps := h.parseApplications()
	metrics := computeMetrics(apps)
	writeJSON(w, metrics)
}

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	profilePath := filepath.Join(h.careerOpsPath, "config", "profile.yml")
	content, err := os.ReadFile(profilePath)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(content)
}

func (h *Handler) updateStatus(w http.ResponseWriter, r *http.Request) {
	numStr := r.PathValue("number")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		http.Error(w, "Invalid application number", http.StatusBadRequest)
		return
	}

	var update StatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if !validStatuses[update.Status] {
		http.Error(w, "Invalid status. Must be one of: Evaluated, Applied, Responded, Interview, Offer, Rejected, Discarded, SKIP", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	apps := h.parseApplications()
	if num < 1 || num > len(apps) {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	app := apps[num-1]
	if err := h.updateApplicationField(app, "status", update.Status); err != nil {
		http.Error(w, "Failed to update: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"status": "ok", "newStatus": update.Status})
}

func (h *Handler) updateNotes(w http.ResponseWriter, r *http.Request) {
	numStr := r.PathValue("number")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		http.Error(w, "Invalid application number", http.StatusBadRequest)
		return
	}

	var update NotesUpdate
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Sanitize: strip pipe chars from notes to avoid breaking markdown table
	update.Notes = strings.ReplaceAll(update.Notes, "|", "-")

	h.mu.Lock()
	defer h.mu.Unlock()

	apps := h.parseApplications()
	if num < 1 || num > len(apps) {
		http.Error(w, "Application not found", http.StatusNotFound)
		return
	}

	app := apps[num-1]
	if err := h.updateApplicationField(app, "notes", update.Notes); err != nil {
		http.Error(w, "Failed to update: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{"status": "ok"})
}

func (h *Handler) getReport(w http.ResponseWriter, r *http.Request) {
	reportPath := r.PathValue("path")

	// Security: prevent path traversal
	cleaned := filepath.Clean(reportPath)
	if strings.Contains(cleaned, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join(h.careerOpsPath, cleaned)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(content)
}

// parseApplications reads applications.md and returns parsed applications.
func (h *Handler) parseApplications() []Application {
	filePath := filepath.Join(h.careerOpsPath, "data", "applications.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		filePath = filepath.Join(h.careerOpsPath, "applications.md")
		content, err = os.ReadFile(filePath)
		if err != nil {
			return nil
		}
	}

	lines := strings.Split(string(content), "\n")
	var apps []Application
	num := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "|---") || strings.HasPrefix(line, "| #") {
			continue
		}
		if !strings.HasPrefix(line, "|") {
			continue
		}

		var fields []string
		if strings.Contains(line, "\t") {
			line = strings.TrimPrefix(line, "|")
			line = strings.TrimSpace(line)
			parts := strings.Split(line, "\t")
			for _, p := range parts {
				fields = append(fields, strings.TrimSpace(strings.Trim(p, "|")))
			}
		} else {
			line = strings.Trim(line, "|")
			parts := strings.Split(line, "|")
			for _, p := range parts {
				fields = append(fields, strings.TrimSpace(p))
			}
		}

		if len(fields) < 8 {
			continue
		}

		num++
		app := Application{
			Number:  num,
			Date:    fields[1],
			Company: fields[2],
			Role:    fields[3],
			Status:  fields[5],
			HasPDF:  strings.Contains(fields[6], "\u2705"),
		}

		app.StatusNorm = normalizeStatus(app.Status)
		app.ScoreRaw = fields[4]
		if sm := reScoreValue.FindStringSubmatch(fields[4]); sm != nil {
			app.Score, _ = strconv.ParseFloat(sm[1], 64)
		}

		if rm := reReportLink.FindStringSubmatch(fields[7]); rm != nil {
			app.ReportNumber = rm[1]
			app.ReportPath = rm[2]
		}

		if len(fields) > 8 {
			app.Notes = fields[8]
		}

		// Try to get job URL from report
		if app.ReportPath != "" {
			fullReport := filepath.Join(h.careerOpsPath, app.ReportPath)
			reportContent, err := os.ReadFile(fullReport)
			if err == nil {
				header := string(reportContent)
				if len(header) > 1000 {
					header = header[:1000]
				}
				if m := reReportURL.FindStringSubmatch(header); m != nil {
					app.JobURL = m[1]
				}
			}
		}

		apps = append(apps, app)
	}

	return apps
}

// updateApplicationField updates a specific field in applications.md
func (h *Handler) updateApplicationField(app Application, field, value string) error {
	filePath := filepath.Join(h.careerOpsPath, "data", "applications.md")
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	found := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "|") || strings.HasPrefix(trimmed, "| #") || strings.HasPrefix(trimmed, "|---") {
			continue
		}

		// Match by report number
		if app.ReportNumber != "" && strings.Contains(line, fmt.Sprintf("[%s]", app.ReportNumber)) {
			switch field {
			case "status":
				lines[i] = replaceFieldInLine(line, 5, value)
			case "notes":
				lines[i] = replaceFieldInLine(line, 8, value)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("application #%d not found in tracker", app.Number)
	}

	return os.WriteFile(filePath, []byte(strings.Join(lines, "\n")), 0644)
}

// replaceFieldInLine replaces a specific field (0-indexed) in a markdown table line.
func replaceFieldInLine(line string, fieldIndex int, newValue string) string {
	// Handle pipe-separated format
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") {
		return line
	}

	parts := strings.Split(trimmed, "|")
	// parts[0] is empty (before first |), so field N is at parts[N+1]
	targetIdx := fieldIndex + 1
	if targetIdx >= len(parts)-1 {
		// Need to extend the line
		for len(parts) <= targetIdx+1 {
			parts = append(parts, " ")
		}
	}

	parts[targetIdx] = " " + newValue + " "
	return strings.Join(parts, "|")
}

// normalizeStatus normalizes raw status text to a canonical form.
func normalizeStatus(raw string) string {
	s := strings.ReplaceAll(raw, "**", "")
	s = strings.TrimSpace(strings.ToLower(s))
	if idx := strings.Index(s, " 202"); idx > 0 {
		s = strings.TrimSpace(s[:idx])
	}

	switch {
	case strings.Contains(s, "no aplicar") || strings.Contains(s, "no_aplicar") || s == "skip" || strings.Contains(s, "geo blocker"):
		return "skip"
	case strings.Contains(s, "interview") || strings.Contains(s, "entrevista"):
		return "interview"
	case s == "offer" || strings.Contains(s, "oferta"):
		return "offer"
	case strings.Contains(s, "responded") || strings.Contains(s, "respondido"):
		return "responded"
	case strings.Contains(s, "applied") || strings.Contains(s, "aplicado") || s == "enviada" || s == "aplicada" || s == "sent":
		return "applied"
	case strings.Contains(s, "rejected") || strings.Contains(s, "rechazado") || s == "rechazada":
		return "rejected"
	case strings.Contains(s, "discarded") || strings.Contains(s, "descartado") || s == "descartada" || s == "cerrada" || s == "cancelada":
		return "discarded"
	case strings.Contains(s, "evaluated") || strings.Contains(s, "evaluada") || s == "hold" || s == "monitor":
		return "evaluated"
	default:
		return s
	}
}

func computeMetrics(apps []Application) Metrics {
	m := Metrics{
		Total:    len(apps),
		ByStatus: make(map[string]int),
	}

	var totalScore float64
	var scored int

	for _, app := range apps {
		status := normalizeStatus(app.Status)
		m.ByStatus[status]++

		if app.Score > 0 {
			totalScore += app.Score
			scored++
			if app.Score > m.TopScore {
				m.TopScore = app.Score
			}
		}
		if app.HasPDF {
			m.WithPDF++
		}
		if status != "skip" && status != "rejected" && status != "discarded" {
			m.Actionable++
		}
	}

	if scored > 0 {
		m.AvgScore = totalScore / float64(scored)
	}

	// Sort status counts for consistent display
	sortedStatuses := []string{"interview", "offer", "responded", "applied", "evaluated", "skip", "rejected", "discarded"}
	sorted := make(map[string]int)
	for _, s := range sortedStatuses {
		if c, ok := m.ByStatus[s]; ok {
			sorted[s] = c
		}
	}
	m.ByStatus = sorted

	_ = sort.SliceIsSorted // import usage
	return m
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
