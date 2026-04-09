# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0] - 2026-04-09

### Added
- **`evaluate` command** — dedicated A–F evaluation mode for single offers
- **`compare` command** — side-by-side weighted comparison matrix for 3–5 offers
- **`outreach` command** — LinkedIn outreach messages with proof-point templates
- **Web dashboard** (`webui/`) — browser-based pipeline viewer at `localhost:8080`
- `npm run dashboard` script to launch the web UI

### Changed
- Renamed Spanish commands to English: `contacto` → `outreach`, `oferta` → `evaluate`
- Streamlined all 14 mode instruction files for clarity and consistency
- Refactored `GEMINI.md` system prompt for cleaner agent bootstrapping
- Refactored `batch-prompt.md` for improved worker instructions
- Simplified state definitions in `templates/states.yml`
- Improved `dedup-tracker.mjs` and `merge-tracker.mjs` scripts
- Browser automation now uses Playwright directly (removed Chrome DevTools MCP dependency)

### Removed
- `.mcp.json` — Chrome DevTools MCP configuration
- Spanish-only commands: `contacto.toml`, `oferta.toml`, `ofertas.toml`
- Spanish-only modes: `contacto.md`, `oferta.md`, `ofertas.md`

## [1.0.0] - 2026-04-06

### Added
- Initial Gemini CLI migration from [santifer/career-ops](https://github.com/santifer/career-ops)
- 13 `.toml` command definitions for Gemini CLI
- Chrome DevTools MCP integration for browser automation
- Go TUI dashboard (Bubble Tea + Lipgloss)
- Batch processing with `gemini --yolo` workers
- ATS-optimized PDF generation via Playwright
- Portal scanner with 45+ pre-configured companies
- Pipeline integrity tools: merge, dedup, normalize, verify
