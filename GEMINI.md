# Career-Ops -- AI Job Search Pipeline

## Origin

This system was built and used by [santifer](https://santifer.io) to evaluate 740+ job offers, generate 100+ tailored CVs, and land a Head of Applied AI role. The archetypes, scoring logic, negotiation scripts, and proof point structure all reflect his specific career search in AI/automation roles.

The portfolio that goes with this system is also open source: [cv-santiago](https://github.com/santifer/cv-santiago).

**It will work out of the box, but it's designed to be made yours.** If the archetypes don't match your career, the modes are in the wrong language, or the scoring doesn't fit your priorities -- just ask. You (Gemini) can edit any file in this system. The user says "change the archetypes to data engineering roles" and you do it. That's the whole point.

## What is career-ops

AI-powered job search automation built on Gemini CLI: pipeline tracking, offer evaluation, CV generation, portal scanning, batch processing.

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

1.  **Read `cv.md`**: Load it silently. This is the canonical source of truth.
2.  **Read `config/profile.yml`**: Load candidate identity and targets.
3.  **Read `portals.yml`**: Load search configuration.
4.  **Read `data/applications.md`**: Load application tracker.

**If any file is missing, enter Onboarding Mode immediately.** If they exist, proceed with the requested task. **Do NOT ask the user if they already have these files if the check passes. Do NOT prompt the user to share their CV if cv.md exists.**

### Onboarding Steps (If missing)

#### Step 1: CV (Required)
If `cv.md` is missing, ask:
> "I don't have your CV yet. You can either:
> 1. Paste your CV here and I'll convert it to markdown
> 2. Paste your LinkedIn URL and I'll extract the key info
> 3. Tell me about your experience and I'll draft a CV for you
>
> Which do you prefer?"

#### Step 2: Profile (Required)
If `config/profile.yml` is missing, copy from `config/profile.example.yml` and then ask for missing details (Full name, email, target roles, salary range).

#### Step 3: Portals (Recommended)
If `portals.yml` is missing, copy from `templates/portals.example.yml`.

#### Step 4: Tracker
If `data/applications.md` is missing, create it with the standard header.

---

### Command Modes

Commands use Gemini CLI's `/group:command` syntax. Each maps to a `.toml` file in `.gemini/commands/career-ops/`.

| If the user... | Command |
|----------------|---------|
| Pastes JD or URL | `/career-ops:auto-pipeline` (evaluate + report + PDF + tracker) |
| Asks to evaluate offer | `/career-ops:evaluate` |
| Asks to compare offers | `/career-ops:compare` |
| Wants LinkedIn outreach | `/career-ops:outreach` |
| Asks for company research | `/career-ops:deep` |
| Wants to generate CV/PDF | `/career-ops:pdf` |
| Evaluates a course/cert | `/career-ops:training` |
| Evaluates portfolio project | `/career-ops:project` |
| Asks about application status | `/career-ops:tracker` |
| Fills out application form | `/career-ops:apply` |
| Searches for new offers | `/career-ops:scan` |
| Processes pending URLs | `/career-ops:pipeline` |
| Batch processes offers | `/career-ops:batch` |

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

**For batch workers (`gemini --yolo`):** Use `web-fetch` and mark the report header with `**Verification:** unconfirmed (batch mode)` if uncertain. The user can verify manually later.

---

## Native Tool Optimization

Gemini CLI provides built-in tools, and Chrome DevTools MCP adds real browser automation:

| Tool | Use |
|------|-----|
| `web-search` | Comp research, trends, company culture, LinkedIn contacts, broad job discovery |
| `web-fetch` | Extract JDs from URLs (static pages), verify offer status, read company pages |
| Chrome DevTools MCP | Navigate portals, scan job listings, fill application forms, verify offers in a real Chrome browser (SPAs, JS-heavy pages) |
| `shell` | Run Node.js scripts (`node generate-pdf.mjs`, `node merge-tracker.mjs`) |
| File I/O | Read cv.md, article-digest.md, cv-template.html; Write reports, tracker TSVs |

**When to use Chrome DevTools MCP vs web-fetch:**
- **Chrome DevTools MCP**: SPA career pages (Ashby, Lever, Workday), form filling, pages requiring JS rendering, offer verification
- **web-fetch**: Static pages, APIs, quick URL checks, batch worker mode (no browser available)

The Chrome DevTools MCP server is configured in `.mcp.json` and starts automatically when Gemini CLI calls a browser tool.

---

## Stack and Conventions

- Node.js (mjs modules), Playwright (PDF generation), YAML (config), HTML/CSS (template), Markdown (data)
- Scripts in `.mjs`, configuration in YAML
- Output in `output/` (gitignored), Reports in `reports/`
- JDs in `jds/` (referenced as `local:jds/{file}` in pipeline.md)
- Batch in `batch/` (gitignored except scripts and prompt)
- Report numbering: sequential 3-digit zero-padded, max existing + 1
- **RULE: After each batch of evaluations, run `node merge-tracker.mjs`** to merge tracker additions and avoid duplications.
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
5. Health check: `node verify-pipeline.mjs`
6. Normalize statuses: `node normalize-statuses.mjs`
7. Dedup: `node dedup-tracker.mjs`

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
