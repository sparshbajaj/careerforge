package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/santifer/career-ops/webui/internal/api"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	pathFlag := flag.String("path", "..", "Path to career-ops directory")
	portFlag := flag.Int("port", 8080, "Port to listen on")
	openFlag := flag.Bool("open", true, "Open browser on start")
	flag.Parse()

	careerOpsPath := *pathFlag

	// Verify the path has applications.md
	if _, err := os.Stat(careerOpsPath + "/data/applications.md"); err != nil {
		if _, err2 := os.Stat(careerOpsPath + "/applications.md"); err2 != nil {
			fmt.Fprintf(os.Stderr, "Error: could not find applications.md in %s or %s/data/\n", careerOpsPath, careerOpsPath)
			os.Exit(1)
		}
	}

	handler := api.NewHandler(careerOpsPath, staticFiles)

	addr := fmt.Sprintf(":%d", *portFlag)
	url := fmt.Sprintf("http://localhost:%d", *portFlag)
	fmt.Printf("Career-Ops Web Dashboard running at %s\n", url)
	fmt.Printf("Data directory: %s\n", careerOpsPath)

	if *openFlag {
		go openBrowser(url)
	}

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}
