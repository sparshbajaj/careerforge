# CareerForge -- AI Job Search Pipeline

## What is CareerForge

AI-powered job search automation built on Antigravity CLI: pipeline tracking, offer evaluation, CV generation, portal scanning, batch processing.

### Main Files

| File | Function |
|------|----------|
| `data/applications.md` | Application tracker |
| `data/pipeline.md` | Inbox of pending URLs |
| `data/scan-history.tsv` | Scanner dedup history |
| `portals.yml` | Query and company config |
| `templates/cv-template.html` | HTML template for CVs |
| `generate-pdf.mjs` | Puppeteer: HTML to PDF |
| `article-digest.md` | Compact proof points from portfolio (optional) |
| `interview-prep/story-bank.md` | Accumulated STAR+R stories across evaluations |
| `reports/` | Evaluation reports (format: `{###}-{company-slug}-{YYYY-MM-DD}.md`) |

### Core Workflow — Environment Check (CRITICAL)

**Before performing any action (evaluations, scans, or mode execution), you MUST verify the environment is ready.** This check MUST be performed **silently and automatically** at the start of every session, new conversation, or when a new job URL is provided. **Do NOT ask the user if they have these files — read them directly.**

1.  **Read `cv.md` using shell (`Get-Content` or `cat`)**: Load it silently. It is in `.gitignore` so `read_file` will fail. This is the canonical source of truth.
2.  **Read `config/profile.yml` using shell**: Load candidate identity and targets.
3.  **Read `portals.yml` using shell**: Load search configuration.
4.  **Read `data/applications.md` using shell**: Load application tracker.

**If any file is missing, enter Onboarding Mode immediately.** If they exist, proceed with the requested task. **Do NOT ask the user if they already have these files if the check passes. Do NOT prompt the user to share their CV if cv.md exists.**

### Onboarding Steps (If missing files)

If `cv.md` or `config/profile.yml` are missing, trigger **Smart Onboarding Mode**:

1. **Read Context**: Scan the `context/` folder for any user-provided files (raw CV text, files with URLs, LinkedIn links, or notes).
2. **Auto-Generate**: Use the provided context to automatically generate the `cv.md` and `config/profile.yml`. **DO NOT** make the user fill out the YAML manually if you have enough context to infer their name, email, and target roles.
3. **Ask for Missing Context**: If the `context/` folder is empty or the provided context is insufficient to build a complete profile/CV, ask the user specific questions to fill in the gaps. Let them know they can drop files or links directly into the `context/` folder.
4. **Save & Setup Portals**: Once you have a complete picture, save the `profile.yml` and `cv.md`. Also, copy `templates/portals.example.yml` to `portals.yml` and create `data/applications.md` with the standard header if they don't exist. Confirm onboarding is complete.

---

## Candidate Specific Strategy (CRITICAL)

**Read the candidate's strategy from `config/profile.yml` and apply it to evaluations and CV generation.**
Do not assume the candidate is a UX/Product Designer. Dynamically read their priority dimensions, positioning, and keyword rules from the strategy block.

---

## Command Modes


Commands use Antigravity CLI's `/group:command` syntax. Each maps to a `.toml` file in `.agents/commands/careerforge/`.

| If the user... | Command |
|----------------|---------|
| Pastes JD or URL | `/careerforge:auto-pipeline` (evaluate + report + PDF + tracker) |
| Asks to evaluate offer | `/careerforge:evaluate` |
| Asks to compare offers | `/careerforge:compare` |
| Wants LinkedIn outreach | `/careerforge:outreach` |
| Asks for company research | `/careerforge:deep` |
| Wants to generate CV/PDF | `/careerforge:pdf` |
| Evaluates a course/cert | `/careerforge:training` |
| Evaluates portfolio project | `/careerforge:project` |
| Asks about application status | `/careerforge:tracker` |
| Fills out application form | `/careerforge:apply` |
| Searches for new offers | `/careerforge:scan` |
| Processes pending URLs | `/careerforge:pipeline` |
| Batch processes offers | `/careerforge:batch` |

### CV Source of Truth

- `cv.md` in project root is the canonical CV
- `article-digest.md` has detailed proof points (optional)
- **NEVER hardcode metrics** -- read them from these files at evaluation time

---

## Ethical Use -- CRITICAL

**This system is designed for quality, not quantity.** The goal is to help the user find and apply to roles where there is a genuine match -- not to spam companies with mass applications.

- **NEVER submit an application without the user reviewing it first.** Fill forms, draft answers, generate PDFs -- but always STOP before clicking Submit/Send/Apply. The user makes the final call.
- **Discourage low-fit applications.** If a score is below 3.0/5, explicitly tell the user this is a weak match and recommend skipping unless they have a specific reason.
- **Quality over speed.** A well-targeted application to 5 companies beats a generic blast to 50. Guide the user toward fewer, better applications.
- **Respect recruiters' time.** Every application a human reads costs someone's attention. Only send what's worth reading.

---

## Offer Verification -- MANDATORY

**Use Chrome DevTools MCP to verify offers in a real browser.** This is the most reliable method:
1. `navigate_page` to the URL
2. `take_snapshot` to read the page content
3. Only footer/navbar without JD content = closed. Title + description + Apply = active.
4. 404 or redirect to generic careers page = closed.

**Fallback:** If Chrome DevTools MCP is unavailable, use `web-fetch` to check the URL.

**For batch workers (`agy --yolo`):** Use `web-fetch` and mark the report header with `**Verification:** unconfirmed (batch mode)` if uncertain. The user can verify manually later.

---

## Native Tool Optimization

Antigravity CLI provides built-in tools, and Chrome DevTools MCP adds real browser automation:

| Tool | Use |
|------|-----|
| `web-search` | Comp research, trends, company culture, LinkedIn contacts, broad job discovery |
| `web-fetch` | Extract JDs from URLs (static pages), verify offer status, read company pages |
| Chrome DevTools MCP | Navigate portals, scan job listings, fill application forms, verify offers in a real Chrome browser (SPAs, JS-heavy pages) |
| `shell` | Run Node.js scripts (`node scripts/generate-pdf.mjs`, `node scripts/merge-tracker.mjs`) |
| File I/O | Read cv.md, article-digest.md, cv-template.html; Write reports, tracker TSVs |

**When to use Chrome DevTools MCP vs web-fetch:**
- **Chrome DevTools MCP**: SPA career pages (Ashby, Lever, Workday), form filling, pages requiring JS rendering, offer verification
- **web-fetch**: Static pages, APIs, quick URL checks, batch worker mode (no browser available)

The Chrome DevTools MCP server is configured in `.mcp.json` and starts automatically when Antigravity CLI calls a browser tool.

---

## Stack and Conventions

- Node.js (mjs modules), Playwright (PDF generation), YAML (config), HTML/CSS (template), Markdown (data)
- Scripts in `.mjs`, configuration in YAML
- Output in `output/` (gitignored), Reports in `reports/`
- JDs in `jds/` (referenced as `local:jds/{file}` in pipeline.md)
- Batch in `batch/` (gitignored except scripts and prompt)
- Report numbering: sequential 3-digit zero-padded, max existing + 1
- **RULE: After each batch of evaluations, run `node scripts/merge-tracker.mjs`** to merge tracker additions and avoid duplications.
- **RULE: NEVER create new entries in applications.md if company+role already exists.** Update the existing entry.

### TSV Format for Tracker Additions

Write one TSV file per evaluation to `batch/tracker-additions/{num}-{company-slug}.tsv`. Single line, 9 tab-separated columns:

```
{num}\t{date}\t{company}\t{role}\t{status}\t{score}/5\t{pdf_emoji}\t[{num}](reports/{num}-{slug}-{date}.md)\t{note}
```

**Column order (IMPORTANT -- status BEFORE score):**
1. `num` -- sequential number (integer)
2. `date` -- YYYY-MM-DD
3. `company` -- short company name
4. `role` -- job title
5. `status` -- canonical status (e.g., `Evaluated`)
6. `score` -- format `X.X/5` (e.g., `4.2/5`)
7. `pdf` -- `✅` or `❌`
8. `report` -- markdown link `[num](reports/...)`
9. `notes` -- one-line summary

**Note:** In applications.md, score comes BEFORE status. The merge script handles this column swap automatically.

### Pipeline Integrity

1. **NEVER edit applications.md to ADD new entries** -- Write TSV in `batch/tracker-additions/` and `merge-tracker.mjs` handles the merge.
2. **YES you can edit applications.md to UPDATE status/notes of existing entries.**
3. All reports MUST include `**URL:**` in the header (between Score and PDF).
4. All statuses MUST be canonical (see `templates/states.yml`).
5. Health check: `node scripts/verify-pipeline.mjs`
6. Normalize statuses: `node scripts/normalize-statuses.mjs`
7. Dedup: `node scripts/dedup-tracker.mjs`

### Canonical States (applications.md)

**Source of truth:** `templates/states.yml`

| State | When to use |
|-------|-------------|
| `Evaluated` | Report completed, pending decision |
| `Applied` | Application sent |
| `Responded` | Company responded |
| `Interview` | In interview process |
| `Offer` | Offer received |
| `Rejected` | Rejected by company |
| `Discarded` | Discarded by candidate or offer closed |
| `SKIP` | Doesn't fit, don't apply |

**RULES:**
- No markdown bold (`**`) in status field
- No dates in status field (use the date column)
- No extra text (use the notes column)
