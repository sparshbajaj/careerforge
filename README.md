# Career-Ops Gemini

AI-powered job search pipeline built on [Google Gemini CLI](https://github.com/google-gemini/gemini-cli). Evaluate offers, generate tailored CVs, scan portals, fill applications, and track everything — powered by AI agents with real browser automation.

[![Gemini CLI](https://img.shields.io/badge/Gemini_CLI-4285F4?style=flat&logo=google&logoColor=white)](https://github.com/google-gemini/gemini-cli)
[![Chrome DevTools MCP](https://img.shields.io/badge/Chrome_DevTools_MCP-4285F4?style=flat&logo=googlechrome&logoColor=white)](https://github.com/AjayKannan97/chrome-devtools-mcp)
[![Node.js](https://img.shields.io/badge/Node.js-339933?style=flat&logo=node.js&logoColor=white)](https://nodejs.org)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> **Fork of [santifer/career-ops](https://github.com/santifer/career-ops)** — originally built on Claude Code by [Santiago Fernández de Valderrama](https://santifer.io). This fork migrates the agent layer to Gemini CLI's `.toml` command architecture, 1M token context, native `web-search`/`web-fetch` tools, and adds [Chrome DevTools MCP](https://github.com/AjayKannan97/chrome-devtools-mcp) for real browser automation.

---

## What It Does

Career-Ops turns Gemini CLI into a job search command center. Instead of manually tracking applications in a spreadsheet, you get an AI pipeline that:

- **Evaluates offers** with a structured A–F scoring system across 10 weighted dimensions
- **Generates tailored PDFs** — ATS-optimized CVs customized per job description
- **Scans portals** automatically (Greenhouse, Ashby, Lever, company career pages)
- **Fills application forms** using Chrome DevTools MCP for real browser interaction
- **Processes in batch** — evaluate 10+ offers in parallel with sub-agents
- **Tracks everything** in a single source of truth with automated integrity checks

> **This is NOT a spray-and-pray tool.** The scoring system helps you focus on high-fit opportunities. Always review before submitting.

## Features

| Feature | Description |
|---------|-------------|
| **Auto-Pipeline** | Paste a URL → full evaluation + tailored PDF + tracker entry |
| **6-Block Evaluation** | Role summary, CV match, level strategy, comp research, personalization, interview prep |
| **Story Bank** | Accumulates STAR+R stories across evaluations for interview prep |
| **ATS PDF Generation** | Keyword-injected CVs built from an HTML template |
| **Browser Automation** | Chrome DevTools MCP for portal scanning, form filling, and verification |
| **Portal Scanner** | 45+ companies pre-configured + custom queries across major job boards |
| **Batch Processing** | Parallel evaluation with `gemini --yolo` workers |
| **Dashboard TUI** | Go-based terminal UI to browse, filter, and sort your pipeline |
| **Pipeline Integrity** | Automated merge, dedup, status normalization, and health checks |

## Quick Start

```bash
# 1. Clone and install
git clone https://github.com/sparshbajaj/career-ops-gemini.git
cd career-ops-gemini
npm install

# 2. Install Gemini CLI
npm install -g @google/gemini-cli

# 3. Configure your profile
cp config/profile.example.yml config/profile.yml   # Edit with your details
cp templates/portals.example.yml portals.yml        # Customize target companies

# 4. Add your CV
# Create cv.md in the project root with your CV in markdown format

# 5. Launch
gemini
```

Then paste a job URL or use any command:

```
/career-ops:auto-pipeline {JD or URL}   → Full pipeline (evaluate + PDF + tracker)
/career-ops:scan                        → Scan portals for new offers
/career-ops:pdf                         → Generate ATS-optimized CV
/career-ops:apply                       → Fill application forms with AI
/career-ops:batch                       → Batch evaluate multiple offers
/career-ops:tracker                     → View application status
/career-ops:pipeline                    → Process pending URLs
/career-ops:deep                        → Deep company research (6 axes)
/career-ops:contacto                    → LinkedIn outreach messages
/career-ops:training                    → Evaluate a course/certification
/career-ops:project                     → Evaluate a portfolio project idea
```

Or just paste a job URL directly — career-ops auto-detects it and runs the full pipeline.

## How It Works

```
Paste a job URL or description
        │
        ▼
┌──────────────────┐
│  Archetype       │  Classifies role type and seniority
│  Detection       │
└────────┬─────────┘
         │
┌────────▼─────────┐
│  A–F Evaluation   │  Match analysis, gaps, comp research, STAR stories
│  (reads cv.md)    │
└────────┬─────────┘
         │
    ┌────┼────┐
    ▼    ▼    ▼
 Report  PDF  Tracker
  .md   .pdf   .tsv
```

## Browser Automation

Career-Ops uses [Chrome DevTools MCP](https://github.com/AjayKannan97/chrome-devtools-mcp) for real browser interaction:

- **Scanning** — Navigate SPAs, click through pagination, extract job listings
- **Applications** — Read form fields, generate tailored answers, fill forms (never auto-submits)
- **Verification** — Load offer pages to confirm details

Start Chrome with remote debugging before using browser features:

```bash
chrome --remote-debugging-port=9222
```

## Pre-configured Portals

The scanner includes **45+ companies** and **19 search queries** across major job boards. Copy and customize:

| Category | Companies |
|----------|-----------|
| **AI Labs** | Anthropic, OpenAI, Mistral, Cohere, LangChain, Pinecone |
| **Voice AI** | ElevenLabs, PolyAI, Parloa, Hume AI, Deepgram, Vapi |
| **AI Platforms** | Retool, Airtable, Vercel, Temporal, Glean, Arize AI |
| **Contact Center** | Ada, LivePerson, Sierra, Decagon, Talkdesk, Genesys |
| **Enterprise** | Salesforce, Twilio, Gong, Dialpad |
| **Automation** | n8n, Zapier, Make.com |
| **Job Boards** | Ashby, Greenhouse, Lever, Wellfound, Workable, RemoteFront |

## Project Structure

```
career-ops-gemini/
├── GEMINI.md                    # Agent system instructions
├── .gemini/commands/career-ops/ # 13 .toml command definitions
├── .mcp.json                    # Chrome DevTools MCP config
├── modes/                       # 14 skill modes
│   ├── _shared.md               # Shared context & archetypes
│   ├── scan.md                  # Portal scanner
│   ├── apply.md                 # Application assistant
│   ├── batch.md                 # Batch processing
│   └── ...                      # oferta, pdf, deep, etc.
├── config/
│   └── profile.example.yml      # Profile template
├── templates/
│   ├── cv-template.html         # ATS CV template
│   ├── portals.example.yml      # Scanner config template
│   └── states.yml               # Canonical application statuses
├── batch/
│   ├── batch-prompt.md          # Worker prompt
│   └── batch-runner.sh          # Batch orchestrator
├── dashboard/                   # Go TUI pipeline viewer
├── examples/                    # Sample CV, report, article digest
├── docs/                        # Setup, architecture, customization
└── fonts/                       # PDF typography
```

## Tech Stack

| Layer | Technology |
|-------|------------|
| **Agent** | [Gemini CLI](https://github.com/google-gemini/gemini-cli) with `.toml` custom commands |
| **Browser** | [Chrome DevTools MCP](https://github.com/AjayKannan97/chrome-devtools-mcp) |
| **PDF** | Playwright + HTML template |
| **Search** | Gemini native `web-search` / `web-fetch` + Greenhouse API |
| **Dashboard** | Go + Bubble Tea + Lipgloss (Catppuccin Mocha) |
| **Data** | Markdown tables, YAML config, TSV batch files |

## Documentation

- [Setup Guide](docs/SETUP.md) — Installation and configuration
- [Architecture](docs/ARCHITECTURE.md) — How the system works
- [Customization](docs/CUSTOMIZATION.md) — Adapting modes, archetypes, and scoring

## Acknowledgments

This project is a fork of [santifer/career-ops](https://github.com/santifer/career-ops). All credit for the original system design, scoring logic, evaluation framework, and pipeline architecture goes to [Santiago Fernández de Valderrama](https://santifer.io). Read his [case study](https://santifer.io/career-ops-system) on how he used the original system to evaluate 740+ offers and land a Head of Applied AI role.

## License

MIT
