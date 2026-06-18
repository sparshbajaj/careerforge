package screens

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sparshbajaj/careerforge/dashboard/internal/data"
	"github.com/sparshbajaj/careerforge/dashboard/internal/model"
	"github.com/sparshbajaj/careerforge/dashboard/internal/theme"
)

// PipelineClosedMsg is emitted when the pipeline screen is dismissed.
type PipelineClosedMsg struct{}

// PipelineOpenReportMsg is emitted when a report should be opened in FileViewer.
type PipelineOpenReportMsg struct {
	Path   string
	Title  string
	JobURL string
}

// PipelineOpenURLMsg is emitted when a job URL should be opened in browser.
type PipelineOpenURLMsg struct {
	URL string
}

// PipelineLoadReportMsg requests lazy loading of a report summary.
type PipelineLoadReportMsg struct {
	CareerOpsPath string
	ReportPath    string
}

// PipelineUpdateStatusMsg requests a status update for an application.
type PipelineUpdateStatusMsg struct {
	CareerOpsPath string
	App           model.CareerApplication
	NewStatus     string
}

type reportSummary struct {
	archetype string
	tldr      string
	remote    string
	comp      string
}

// Sort modes
const (
	sortScore   = "score"
	sortDate    = "date"
	sortCompany = "company"
	sortStatus  = "status"
)

// Filter modes
const (
	filterAll       = "all"
	filterEvaluated = "evaluated"
	filterApplied   = "applied"
	filterInterview = "interview"
	filterSkip      = "skip"
	filterTop       = "top"
)

type pipelineTab struct {
	filter string
	label  string
}

var pipelineTabs = []pipelineTab{
	{filterAll, "ALL"},
	{filterEvaluated, "EVALUATED"},
	{filterApplied, "APPLIED"},
	{filterInterview, "INTERVIEW"},
	{filterTop, "TOP \u22654"},
	{filterSkip, "SKIP"},
}

var sortCycle = []string{sortScore, sortDate, sortCompany, sortStatus}

var statusOptions = []string{"Evaluated", "Applied", "Responded", "Interview", "Offer", "Rejected", "Discarded", "SKIP"}

// statusGroupOrder defines display order for grouped view.
var statusGroupOrder = []string{"interview", "offer", "responded", "applied", "evaluated", "skip", "rejected", "discarded"}

// PipelineModel implements the career pipeline dashboard screen.
type PipelineModel struct {
	apps          []model.CareerApplication
	filtered      []model.CareerApplication
	metrics       model.PipelineMetrics
	cursor        int
	scrollOffset  int
	sortMode      string
	activeTab     int
	viewMode      string // "grouped" or "flat"
	width, height int
	theme         theme.Theme
	careerOpsPath string
	reportCache   map[string]reportSummary
	// Status picker sub-state
	statusPicker bool
	statusCursor int
	// Live update tracking
	lastAppUpdate time.Time
	
	onboarding bool
	setupMsg   string
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// NewPipelineModel creates a new pipeline screen.
func NewPipelineModel(t theme.Theme, apps []model.CareerApplication, metrics model.PipelineMetrics, careerOpsPath string, width, height int) PipelineModel {
	m := PipelineModel{
		apps:          apps,
		metrics:       metrics,
		sortMode:      sortScore,
		activeTab:     0,
		viewMode:      "grouped",
		width:         width,
		height:        height,
		theme:         t,
		careerOpsPath: careerOpsPath,
		reportCache:   make(map[string]reportSummary),
	}
	
	m.checkOnboarding()
	m.applyFilterAndSort()
	return m
}

func (m *PipelineModel) checkOnboarding() {
	missing := []string{}
	if _, err := os.Stat(filepath.Join(m.careerOpsPath, "cv.md")); os.IsNotExist(err) {
		missing = append(missing, "cv.md")
	}
	if _, err := os.Stat(filepath.Join(m.careerOpsPath, "config", "profile.yml")); os.IsNotExist(err) {
		missing = append(missing, "config/profile.yml")
	}
	
	appPath := filepath.Join(m.careerOpsPath, "applications.md")
	if _, err := os.Stat(appPath); os.IsNotExist(err) {
		appPath = filepath.Join(m.careerOpsPath, "data", "applications.md")
		if _, err := os.Stat(appPath); os.IsNotExist(err) {
			missing = append(missing, "applications.md")
		}
	}

	if len(missing) > 0 {
		m.onboarding = true
		m.setupMsg = "Missing required files: " + strings.Join(missing, ", ")
	} else {
		m.onboarding = false
	}
}

// Init implements tea.Model.
func (m PipelineModel) Init() tea.Cmd {
	return tickCmd()
}

// Resize updates dimensions.
func (m *PipelineModel) Resize(width, height int) {
	m.width = width
	m.height = height
}

// Width returns the current width.
func (m PipelineModel) Width() int { return m.width }

// Height returns the current height.
func (m PipelineModel) Height() int { return m.height }

// UpdateData updates the underlying data and recalculates without losing scroll or UI state.
func (m *PipelineModel) UpdateData(apps []model.CareerApplication, metrics model.PipelineMetrics) {
	m.apps = apps
	m.metrics = metrics
	m.applyFilterAndSort()
	m.adjustScroll()
}

// CopyReportCache copies the report cache from another pipeline model.
func (m *PipelineModel) CopyReportCache(other *PipelineModel) {
	for k, v := range other.reportCache {
		m.reportCache[k] = v
	}
}

// EnrichReport caches report summary data for preview.
func (m *PipelineModel) EnrichReport(reportPath, archetype, tldr, remote, comp string) {
	m.reportCache[reportPath] = reportSummary{
		archetype: archetype,
		tldr:      tldr,
		remote:    remote,
		comp:      comp,
	}
}

// CurrentApp returns the currently selected application, if any.
func (m PipelineModel) CurrentApp() (model.CareerApplication, bool) {
	if m.cursor < 0 || m.cursor >= len(m.filtered) {
		return model.CareerApplication{}, false
	}
	return m.filtered[m.cursor], true
}

// Update handles input for the pipeline screen.
func (m PipelineModel) Update(msg tea.Msg) (PipelineModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		// Reload apps if applications.md changed or we are onboarding
		m.checkOnboarding()
		appFile := filepath.Join(m.careerOpsPath, "data", "applications.md")
		if info, err := os.Stat(appFile); err == nil {
			if info.ModTime().After(m.lastAppUpdate) {
				m.lastAppUpdate = info.ModTime()
				apps := data.ParseApplications(m.careerOpsPath)
				if apps != nil {
					metrics := data.ComputeMetrics(apps)
					m.UpdateData(apps, metrics)
				}
			}
		} else if info, err := os.Stat(filepath.Join(m.careerOpsPath, "applications.md")); err == nil {
			if info.ModTime().After(m.lastAppUpdate) {
				m.lastAppUpdate = info.ModTime()
				apps := data.ParseApplications(m.careerOpsPath)
				if apps != nil {
					metrics := data.ComputeMetrics(apps)
					m.UpdateData(apps, metrics)
				}
			}
		}
		// Returning the tick cmd causes a View() re-render which naturally reads the latest ai.log!
		return m, tickCmd()
		
	case tea.KeyMsg:
		if m.statusPicker {
			return m.handleStatusPicker(msg)
		}
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	
	return m, nil
}

func (m PipelineModel) handleKey(msg tea.KeyMsg) (PipelineModel, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		return m, func() tea.Msg { return PipelineClosedMsg{} }

	case "down":
		if len(m.filtered) > 0 {
			m.cursor++
			if m.cursor >= len(m.filtered) {
				m.cursor = len(m.filtered) - 1
			}
			m.adjustScroll()
			return m, m.loadCurrentReport()
		}

	case "up":
		if len(m.filtered) > 0 {
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}
			m.adjustScroll()
			return m, m.loadCurrentReport()
		}

	case "s":
		// Cycle sort mode
		for i, s := range sortCycle {
			if s == m.sortMode {
				m.sortMode = sortCycle[(i+1)%len(sortCycle)]
				break
			}
		}
		m.applyFilterAndSort()
		m.cursor = 0
		m.scrollOffset = 0

	case "f", "right":
		m.activeTab++
		if m.activeTab >= len(pipelineTabs) {
			m.activeTab = 0
		}
		m.applyFilterAndSort()
		m.cursor = 0
		m.scrollOffset = 0

	case "left":
		m.activeTab--
		if m.activeTab < 0 {
			m.activeTab = len(pipelineTabs) - 1
		}
		m.applyFilterAndSort()
		m.cursor = 0
		m.scrollOffset = 0

	case "v":
		if m.viewMode == "grouped" {
			m.viewMode = "flat"
		} else {
			m.viewMode = "grouped"
		}

	case "enter":
		if app, ok := m.CurrentApp(); ok && app.ReportPath != "" {
			fullPath := filepath.Join(m.careerOpsPath, app.ReportPath)
			title := fmt.Sprintf("%s \u2014 %s", app.Company, app.Role)
			jobURL := app.JobURL
			return m, func() tea.Msg {
				return PipelineOpenReportMsg{Path: fullPath, Title: title, JobURL: jobURL}
			}
		}

	case "o":
		if app, ok := m.CurrentApp(); ok && app.JobURL != "" {
			return m, func() tea.Msg {
				return PipelineOpenURLMsg{URL: app.JobURL}
			}
		}

	case "c":
		if len(m.filtered) > 0 {
			m.statusPicker = true
			m.statusCursor = 0
		}

	case "pgdown", "ctrl+d":
		m.scrollOffset += m.height / 2
		return m, nil

	case "pgup", "ctrl+u":
		m.scrollOffset -= m.height / 2
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
		return m, nil
	}

	return m, nil
}

func (m PipelineModel) handleStatusPicker(msg tea.KeyMsg) (PipelineModel, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.statusPicker = false
		return m, nil

	case "down":
		m.statusCursor++
		if m.statusCursor >= len(statusOptions) {
			m.statusCursor = len(statusOptions) - 1
		}

	case "up":
		m.statusCursor--
		if m.statusCursor < 0 {
			m.statusCursor = 0
		}

	case "enter":
		m.statusPicker = false
		if app, ok := m.CurrentApp(); ok {
			newStatus := statusOptions[m.statusCursor]
			return m, func() tea.Msg {
				return PipelineUpdateStatusMsg{
					CareerOpsPath: m.careerOpsPath,
					App:           app,
					NewStatus:     newStatus,
				}
			}
		}
	}
	return m, nil
}

func (m PipelineModel) loadCurrentReport() tea.Cmd {
	app, ok := m.CurrentApp()
	if !ok || app.ReportPath == "" {
		return nil
	}
	if _, cached := m.reportCache[app.ReportPath]; cached {
		return nil
	}
	path := m.careerOpsPath
	report := app.ReportPath
	return func() tea.Msg {
		return PipelineLoadReportMsg{CareerOpsPath: path, ReportPath: report}
	}
}

// applyFilterAndSort rebuilds the filtered list from apps.
func (m *PipelineModel) applyFilterAndSort() {
	var filtered []model.CareerApplication

	currentFilter := pipelineTabs[m.activeTab].filter
	for _, app := range m.apps {
		norm := data.NormalizeStatus(app.Status)
		switch currentFilter {
		case filterAll:
			filtered = append(filtered, app)
		case filterTop:
			if app.Score >= 4.0 && norm != "no_aplicar" {
				filtered = append(filtered, app)
			}
		default:
			if norm == currentFilter {
				filtered = append(filtered, app)
			}
		}
	}

	// Sort
	switch m.sortMode {
	case sortScore:
		sort.SliceStable(filtered, func(i, j int) bool {
			return filtered[i].Score > filtered[j].Score
		})
	case sortDate:
		sort.SliceStable(filtered, func(i, j int) bool {
			return filtered[i].Date > filtered[j].Date
		})
	case sortCompany:
		sort.SliceStable(filtered, func(i, j int) bool {
			return strings.ToLower(filtered[i].Company) < strings.ToLower(filtered[j].Company)
		})
	case sortStatus:
		sort.SliceStable(filtered, func(i, j int) bool {
			return data.StatusPriority(filtered[i].Status) < data.StatusPriority(filtered[j].Status)
		})
	}

	// In grouped mode, always sort by status priority first, then by selected sort within groups
	if m.viewMode == "grouped" {
		sort.SliceStable(filtered, func(i, j int) bool {
			pi := data.StatusPriority(filtered[i].Status)
			pj := data.StatusPriority(filtered[j].Status)
			if pi != pj {
				return pi < pj
			}
			// Within same group, use selected sort
			switch m.sortMode {
			case sortScore:
				return filtered[i].Score > filtered[j].Score
			case sortDate:
				return filtered[i].Date > filtered[j].Date
			case sortCompany:
				return strings.ToLower(filtered[i].Company) < strings.ToLower(filtered[j].Company)
			default:
				return filtered[i].Score > filtered[j].Score
			}
		})
	}

	m.filtered = filtered
}

// adjustScroll updates scrollOffset so the cursor stays visible.
func (m *PipelineModel) adjustScroll() {
	availHeight := m.height - 12 // header + tabs(2) + metrics + sortbar + footer + preview
	if availHeight < 5 {
		availHeight = 5
	}
	line := m.cursorLineEstimate()
	margin := 3

	if line >= m.scrollOffset+availHeight-margin {
		m.scrollOffset = line - availHeight + margin + 1
	}
	if line < m.scrollOffset+margin {
		m.scrollOffset = line - margin
	}
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

func (m PipelineModel) cursorLineEstimate() int {
	if m.viewMode != "grouped" {
		return m.cursor
	}
	// Account for group headers
	line := 0
	prevStatus := ""
	for i, app := range m.filtered {
		norm := data.NormalizeStatus(app.Status)
		if norm != prevStatus {
			line++ // group header
			prevStatus = norm
		}
		if i == m.cursor {
			return line
		}
		line++
	}
	return line
}

// -- View --

// View renders the pipeline screen.
func (m PipelineModel) View() string {
	if m.width == 0 || m.height == 0 {
		return "Initializing..."
	}

	if m.onboarding {
		return m.renderOnboarding()
	}

	showRightPanel := m.width > 110
	rightWidth := 45
	if showRightPanel {
		m.width = m.width - rightWidth - 2
	}

	header := m.renderHeader()
	tabs := m.renderTabs()
	metricsBar := m.renderMetrics()
	sortBar := m.renderSortBar()
	body := m.renderBody()
	preview := m.renderPreview()
	commands := m.renderCommands()
	help := m.renderHelp()

	// Apply scroll to body
	bodyLines := strings.Split(body, "\n")
	if m.scrollOffset > 0 && m.scrollOffset < len(bodyLines) {
		bodyLines = bodyLines[m.scrollOffset:]
	}

	// Calculate available height for body
	previewLines := strings.Count(preview, "\n") + 1
	cmdLines := strings.Count(commands, "\n") + 1
	availHeight := m.height - 8 - previewLines - cmdLines
	if availHeight < 3 {
		availHeight = 3
	}
	if len(bodyLines) > availHeight {
		bodyLines = bodyLines[:availHeight]
	}
	body = strings.Join(bodyLines, "\n")

	// Status picker overlay
	if m.statusPicker {
		body = m.overlayStatusPicker(body)
	}

	leftSide := lipgloss.JoinVertical(lipgloss.Left,
		header,
		tabs,
		metricsBar,
		sortBar,
		body,
		preview,
		commands,
		help,
	)

	if !showRightPanel {
		return leftSide
	}

	rightSide := m.renderRightPanel(rightWidth)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftSide, "  ", rightSide)
}

func (m PipelineModel) renderOnboarding() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Peach).
		Padding(2, 4).
		Width(60)

	title := lipgloss.NewStyle().Foreground(m.theme.Peach).Bold(true).Render("🤖 SMART ONBOARDING MODE")
	
	msg := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		lipgloss.NewStyle().Foreground(m.theme.Text).Render("Welcome to CareerForge Studio!"),
		"",
		lipgloss.NewStyle().Foreground(m.theme.Subtext).Render(m.setupMsg),
		"",
		lipgloss.NewStyle().Foreground(m.theme.Green).Render("👉 To fix this automatically:"),
		lipgloss.NewStyle().Foreground(m.theme.Subtext).Render("1. Drop your raw resume or notes into the `context/` folder."),
		lipgloss.NewStyle().Foreground(m.theme.Subtext).Render("2. Switch to the `agy` pane on the right side."),
		lipgloss.NewStyle().Foreground(m.theme.Subtext).Render("3. Type: \"Generate my profile and CV\""),
		"",
		lipgloss.NewStyle().Foreground(m.theme.Overlay).Render("I am monitoring your files. This screen will disappear"),
		lipgloss.NewStyle().Foreground(m.theme.Overlay).Render("automatically as soon as the files are created!"),
	)

	centered := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, boxStyle.Render(msg))
	return centered
}

func (m PipelineModel) renderRightPanel(width int) string {
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Overlay).
		Width(width - 2).
		Height(m.height - 2)

	mascot := `     _____
    [◎_◎]
   /[___]\
    /   \ `

	mascotStyle := lipgloss.NewStyle().Foreground(m.theme.Peach).Bold(true).Padding(1, 0, 1, 4)
	titleStyle := lipgloss.NewStyle().Foreground(m.theme.Sky).Bold(true).Padding(0, 0, 1, 2)
	infoStyle := lipgloss.NewStyle().Foreground(m.theme.Subtext).Padding(0, 2)
	dimStyle := lipgloss.NewStyle().Foreground(m.theme.Overlay).Padding(0, 2)

	// Tail the AI log
	logLines := []string{"Waiting for AI output..."}
	if b, err := os.ReadFile(filepath.Join(m.careerOpsPath, "data", "ai.log")); err == nil {
		allLines := strings.Split(strings.TrimSpace(string(b)), "\n")
		if len(allLines) > 0 && allLines[0] != "" {
			tailCount := 12
			if len(allLines) < tailCount {
				tailCount = len(allLines)
			}
			logLines = allLines[len(allLines)-tailCount:]
		}
	}
	
	// Format tail logs
	var renderedLogs []string
	for _, l := range logLines {
		if len(l) > width-6 {
			l = l[:width-9] + "..."
		}
		renderedLogs = append(renderedLogs, dimStyle.Render("> "+l))
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		mascotStyle.Render(mascot),
		titleStyle.Render("CLAW-D // AI ACTIVE"),
		infoStyle.Render("I am monitoring the AI!"),
		infoStyle.Render("Pasting a URL instantly"),
		infoStyle.Render("triggers background processing."),
		"",
		titleStyle.Render("LIVE TERMINAL LOGS"),
		lipgloss.JoinVertical(lipgloss.Left, renderedLogs...),
	)

	return borderStyle.Render(content)
}

func (m PipelineModel) renderHeader() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Text).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(m.theme.Mauve).
		Width(m.width).
		Padding(0, 2)

	right := lipgloss.NewStyle().Foreground(m.theme.Subtext)
	avg := fmt.Sprintf("%.1f", m.metrics.AvgScore)
	info := right.Render(fmt.Sprintf("%d offers | Avg %s/5", m.metrics.Total, avg))

	title := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Blue).Render("CAREER PIPELINE")
	gap := m.width - lipgloss.Width(title) - lipgloss.Width(info) - 4
	if gap < 1 {
		gap = 1
	}

	return style.Render(title + strings.Repeat(" ", gap) + info)
}

func (m PipelineModel) renderTabs() string {
	var tabs []string
	var underParts []string

	for i, tab := range pipelineTabs {
		// Count items for this tab
		count := m.countForFilter(tab.filter)
		label := fmt.Sprintf(" %s (%d) ", tab.label, count)

		if i == m.activeTab {
			style := lipgloss.NewStyle().
				Bold(true).
				Foreground(m.theme.Blue).
				Padding(0, 0)
			tabs = append(tabs, style.Render(label))
			underParts = append(underParts, strings.Repeat("\u2501", lipgloss.Width(label)))
		} else {
			style := lipgloss.NewStyle().
				Foreground(m.theme.Subtext).
				Padding(0, 0)
			tabs = append(tabs, style.Render(label))
			underParts = append(underParts, strings.Repeat("\u2500", lipgloss.Width(label)))
		}
	}

	row := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	underline := lipgloss.NewStyle().Foreground(m.theme.Overlay).Render(strings.Join(underParts, ""))

	padStyle := lipgloss.NewStyle().Padding(0, 1)
	return padStyle.Render(row) + "\n" + padStyle.Render(underline)
}

func (m PipelineModel) countForFilter(filter string) int {
	count := 0
	for _, app := range m.apps {
		norm := data.NormalizeStatus(app.Status)
		switch filter {
		case filterAll:
			count++
		case filterTop:
			if app.Score >= 4.0 && norm != "no_aplicar" {
				count++
			}
		default:
			if norm == filter {
				count++
			}
		}
	}
	return count
}

func (m PipelineModel) renderMetrics() string {
	style := lipgloss.NewStyle().
		Background(m.theme.Surface).
		Width(m.width).
		Padding(0, 2)

	var parts []string
	statusColors := m.statusColorMap()

	for _, status := range statusGroupOrder {
		count, ok := m.metrics.ByStatus[status]
		if !ok || count == 0 {
			continue
		}
		color := statusColors[status]
		s := lipgloss.NewStyle().Foreground(color)
		parts = append(parts, s.Render(fmt.Sprintf("%s:%d", statusLabel(status), count)))
	}

	return style.Render(strings.Join(parts, "  "))
}

func (m PipelineModel) renderSortBar() string {
	style := lipgloss.NewStyle().
		Foreground(m.theme.Subtext).
		Width(m.width).
		Padding(0, 2)

	sortLabel := fmt.Sprintf("[Sort: %s]", m.sortMode)
	viewLabel := fmt.Sprintf("[View: %s]", m.viewMode)
	count := fmt.Sprintf("%d shown", len(m.filtered))

	return style.Render(fmt.Sprintf("%s  %s  %s", sortLabel, viewLabel, count))
}

func (m PipelineModel) renderBody() string {
	if len(m.filtered) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(m.theme.Subtext).
			Padding(1, 2)
		return emptyStyle.Render("No offers match this filter")
	}

	var lines []string
	prevStatus := ""
	padStyle := lipgloss.NewStyle().Padding(0, 2)

	for i, app := range m.filtered {
		norm := data.NormalizeStatus(app.Status)

		// Group header in grouped mode
		if m.viewMode == "grouped" && norm != prevStatus {
			count := m.countByNormStatus(norm)
			headerStyle := lipgloss.NewStyle().
				Bold(true).
				Foreground(m.theme.Subtext)
			lines = append(lines, padStyle.Render(
				headerStyle.Render(fmt.Sprintf("\u2500\u2500 %s (%d) %s",
					strings.ToUpper(statusLabel(norm)), count,
					strings.Repeat("\u2500", max(0, m.width-30-len(statusLabel(norm)))))),
			))
			prevStatus = norm
		}

		selected := i == m.cursor
		line := m.renderAppLine(app, selected)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (m PipelineModel) renderAppLine(app model.CareerApplication, selected bool) string {
	padStyle := lipgloss.NewStyle().Padding(0, 2)

	// Column widths
	scoreW := 5   // "4.5  "
	companyW := 20
	statusW := 12
	compW := 14
	// Role gets remaining space
	roleW := m.width - scoreW - companyW - statusW - compW - 10
	if roleW < 15 {
		roleW = 15
	}

	// Score with color
	scoreStyle := m.scoreStyle(app.Score)
	score := scoreStyle.Render(fmt.Sprintf("%.1f", app.Score))

	// Company (truncate)
	company := app.Company
	if len(company) > companyW {
		company = company[:companyW-3] + "..."
	}
	companyStyle := lipgloss.NewStyle().Foreground(m.theme.Text).Width(companyW)

	// Role (truncate)
	role := app.Role
	if len(role) > roleW {
		role = role[:roleW-3] + "..."
	}
	roleStyle := lipgloss.NewStyle().Foreground(m.theme.Subtext).Width(roleW)

	// Status with color -- fixed column
	norm := data.NormalizeStatus(app.Status)
	statusColor := m.statusColorMap()[norm]
	statusStyle := lipgloss.NewStyle().Foreground(statusColor).Width(statusW)
	statusText := statusStyle.Render(statusLabel(norm))

	// Comp from report cache -- fixed column
	compText := ""
	if summary, ok := m.reportCache[app.ReportPath]; ok && summary.comp != "" {
		comp := summary.comp
		if len(comp) > compW-1 {
			comp = comp[:compW-4] + "..."
		}
		compStyle := lipgloss.NewStyle().Foreground(m.theme.Yellow)
		compText = compStyle.Render(comp)
	}

	line := fmt.Sprintf(" %s %s %s %s %s",
		score,
		companyStyle.Render(company),
		roleStyle.Render(role),
		statusText,
		compText,
	)

	if selected {
		selStyle := lipgloss.NewStyle().
			Background(m.theme.Surface).
			Bold(true).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(m.theme.Pink).
			PaddingLeft(1).
			Width(m.width - 6)
		return padStyle.Render(selStyle.Render(strings.TrimPrefix(line, " ")))
	}
	return padStyle.Render("  " + strings.TrimPrefix(line, " "))
}

func (m PipelineModel) renderPreview() string {
	app, ok := m.CurrentApp()
	if !ok {
		return ""
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Mauve).
		Padding(0, 2).
		Width(m.width - 4)

	var lines []string
	labelStyle := lipgloss.NewStyle().Foreground(m.theme.Sky).Bold(true)
	valueStyle := lipgloss.NewStyle().Foreground(m.theme.Text)
	dimStyle := lipgloss.NewStyle().Foreground(m.theme.Subtext)

	// Check report cache
	if summary, ok := m.reportCache[app.ReportPath]; ok {
		if summary.archetype != "" {
			lines = append(lines, labelStyle.Render("Arquetipo: ")+valueStyle.Render(summary.archetype))
		}
		if summary.tldr != "" {
			lines = append(lines, labelStyle.Render("TL;DR: ")+valueStyle.Render(summary.tldr))
		}
		if summary.comp != "" {
			lines = append(lines, labelStyle.Render("Comp: ")+valueStyle.Render(summary.comp))
		}
		if summary.remote != "" {
			lines = append(lines, labelStyle.Render("Remote: ")+valueStyle.Render(summary.remote))
		}
	} else if app.Notes != "" {
		// Fallback: show notes
		notes := app.Notes
		if len(notes) > m.width-10 {
			notes = notes[:m.width-13] + "..."
		}
		lines = append(lines, dimStyle.Render(notes))
	} else {
		lines = append(lines, dimStyle.Render("Loading preview..."))
	}

	return lipgloss.NewStyle().Padding(0, 1).Render(boxStyle.Render(strings.Join(lines, "\n")))
}

func (m PipelineModel) renderCommands() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Overlay).
		Width(m.width - 2).
		Padding(0, 1)

	cmdStyle := lipgloss.NewStyle().Foreground(m.theme.Green).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(m.theme.Subtext)

	row1 := cmdStyle.Render("/careerforge:pipeline") + descStyle.Render(" (process inbox)  ") +
		cmdStyle.Render("/careerforge:evaluate") + descStyle.Render(" (score JD)  ") +
		cmdStyle.Render("/careerforge:pdf") + descStyle.Render(" (gen CV)")
	row2 := cmdStyle.Render("/careerforge:apply") + descStyle.Render(" (fill form)  ") +
		cmdStyle.Render("/careerforge:scan") + descStyle.Render(" (find jobs)  ") +
		cmdStyle.Render("/careerforge:compare") + descStyle.Render(" (compare offers)")
	row3 := cmdStyle.Render("/careerforge:outreach") + descStyle.Render(" (linkedin)  ") +
		cmdStyle.Render("/careerforge:deep") + descStyle.Render(" (research)  ") +
		cmdStyle.Render("/careerforge:batch") + descStyle.Render(" (batch)")

	return lipgloss.NewStyle().Padding(0, 1).Render(boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, row1, row2, row3)))
}

func (m PipelineModel) renderHelp() string {
	style := lipgloss.NewStyle().
		Foreground(m.theme.Subtext).
		Background(m.theme.Surface).
		Width(m.width).
		Padding(0, 1)

	keyStyle := lipgloss.NewStyle().Bold(true).Foreground(m.theme.Text)
	descStyle := lipgloss.NewStyle().Foreground(m.theme.Subtext)

	if m.statusPicker {
		return style.Render(
			keyStyle.Render("↑↓") + descStyle.Render(" navigate  ") +
				keyStyle.Render("Enter") + descStyle.Render(" confirm  ") +
				keyStyle.Render("Esc") + descStyle.Render(" cancel"))
	}

	brand := lipgloss.NewStyle().Foreground(m.theme.Overlay).Render("CareerForge Pipeline")

	keys := keyStyle.Render("↑↓") + descStyle.Render(" nav  ") +
		keyStyle.Render("←→") + descStyle.Render(" tabs  ") +
		keyStyle.Render("s") + descStyle.Render(" sort  ") +
		keyStyle.Render("Enter") + descStyle.Render(" report  ") +
		keyStyle.Render("o") + descStyle.Render(" open URL  ") +
		keyStyle.Render("c") + descStyle.Render(" change  ") +
		keyStyle.Render("v") + descStyle.Render(" view  ") +
		keyStyle.Render("Esc") + descStyle.Render(" quit")

	gap := m.width - lipgloss.Width(keys) - lipgloss.Width(brand) - 2
	if gap < 1 {
		gap = 1
	}

	return style.Render(keys + strings.Repeat(" ", gap) + brand)
}

func (m PipelineModel) overlayStatusPicker(body string) string {
	// Render status picker inline at bottom of body
	bodyLines := strings.Split(body, "\n")

	pickerWidth := 30
	padStyle := lipgloss.NewStyle().Padding(0, 2)
	borderStyle := lipgloss.NewStyle().
		Foreground(m.theme.Blue).
		Bold(true)

	var picker []string
	picker = append(picker, padStyle.Render(borderStyle.Render("Change status:")))

	for i, opt := range statusOptions {
		style := lipgloss.NewStyle().Foreground(m.theme.Text).Width(pickerWidth)
		if i == m.statusCursor {
			style = style.Background(m.theme.Overlay).Bold(true)
		}
		prefix := "  "
		if i == m.statusCursor {
			prefix = "> "
		}
		picker = append(picker, padStyle.Render(style.Render(prefix+opt)))
	}

	// Append picker to body
	bodyLines = append(bodyLines, picker...)
	return strings.Join(bodyLines, "\n")
}

// -- Helpers --

func (m PipelineModel) scoreStyle(score float64) lipgloss.Style {
	switch {
	case score >= 4.2:
		return lipgloss.NewStyle().Foreground(m.theme.Green).Bold(true)
	case score >= 3.8:
		return lipgloss.NewStyle().Foreground(m.theme.Yellow)
	case score >= 3.0:
		return lipgloss.NewStyle().Foreground(m.theme.Text)
	default:
		return lipgloss.NewStyle().Foreground(m.theme.Red)
	}
}

func (m PipelineModel) statusColorMap() map[string]lipgloss.Color {
	return map[string]lipgloss.Color{
		"interview": m.theme.Green,
		"offer":     m.theme.Green,
		"applied":   m.theme.Sky,
		"responded": m.theme.Blue,
		"evaluated": m.theme.Text,
		"skip":      m.theme.Red,
		"rejected":  m.theme.Subtext,
		"discarded": m.theme.Subtext,
	}
}

func (m PipelineModel) countByNormStatus(status string) int {
	count := 0
	for _, app := range m.filtered {
		if data.NormalizeStatus(app.Status) == status {
			count++
		}
	}
	return count
}

func statusLabel(norm string) string {
	switch norm {
	case "interview":
		return "Interview"
	case "offer":
		return "Offer"
	case "responded":
		return "Responded"
	case "applied":
		return "Applied"
	case "evaluated":
		return "Evaluated"
	case "skip":
		return "Skip"
	case "rejected":
		return "Rejected"
	case "discarded":
		return "Discarded"
	default:
		return norm
	}
}
