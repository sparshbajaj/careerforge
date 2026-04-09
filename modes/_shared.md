# Shared Context -- career-ops

<!-- ============================================================
     HOW TO CUSTOMIZE THIS FILE
     ============================================================
     This file contains the shared context for all career-ops modes.
     Before using career-ops, you MUST:
     1. Fill in config/profile.yml with your personal data
     2. Create your cv.md in the project root
     3. (Optional) Create article-digest.md with your proof points
     4. Customize the sections below marked with [CUSTOMIZE]
     ============================================================ -->

## Sources of Truth (ALWAYS read before evaluating)

| File | Path | When |
|------|------|------|
| cv.md | `cv.md` (project root) | ALWAYS |
| article-digest.md | `article-digest.md` (if exists) | ALWAYS (detailed proof points) |
| profile.yml | `config/profile.yml` | ALWAYS (candidate identity and targets) |

**RULE: NEVER hardcode metrics from proof points.** Read them from cv.md + article-digest.md at evaluation time.
**RULE: For article/project metrics, article-digest.md takes precedence over cv.md** (cv.md may have older numbers).

---

## Automatic Environment Check (CRITICAL — runs silently every session)

**At the start of EVERY session, new conversation, or when a link/URL is provided, you MUST silently verify the environment is ready.** Do NOT ask the user if files exist — check them yourself. This check must complete before any other action.

### Check Sequence (run automatically, do NOT prompt user)

1. **Read `cv.md`** — If it exists, load it silently. If missing, enter Onboarding Mode.
2. **Read `config/profile.yml`** — If it exists, load candidate identity. If missing, copy from `config/profile.example.yml` and ask user to fill in details.
3. **Read `portals.yml`** — If it exists, load it. If missing, copy from `templates/portals.example.yml`.
4. **Read `data/applications.md`** — If it exists, load tracker. If missing, create with standard header.
5. **Run `node cv-sync-check.mjs`** — If it reports warnings, notify the candidate before continuing.

### Behavior Rules

- **NEVER ask the user to provide cv.md or profile.yml if they already exist.** Just read them.
- **NEVER say "I need your CV"** if cv.md is present in the project root.
- **Load these files at the start of every session** — do not wait until an evaluation is requested.
- **Cache in context**: Once loaded, keep cv.md and profile.yml contents in your working context for the entire session.
- **If files change mid-session**: Re-read when the user says "I updated my CV" or similar.

---

## North Star -- Target Roles

The skill applies with EQUAL rigor to ALL target roles. None is primary or secondary -- any is a success if compensation and growth are right:

| Archetype | Thematic axes | What they buy |
|-----------|---------------|---------------|
| **AI Platform / LLMOps Engineer** | Evaluation, observability, reliability, pipelines | Someone who puts AI in production with metrics |
| **Agentic Workflows / Automation** | HITL, tooling, orchestration, multi-agent | Someone who builds reliable agent systems |
| **Technical AI Product Manager** | GenAI/Agents, PRDs, discovery, delivery | Someone who translates business to AI product |
| **AI Solutions Architect** | Hyperautomation, enterprise, integrations | Someone who designs end-to-end AI architectures |
| **AI Forward Deployed Engineer** | Client-facing, fast delivery, prototyping | Someone who delivers AI solutions to clients fast |
| **AI Transformation Lead** | Change management, adoption, org enablement | Someone who leads AI transformation in an org |
| **AI/UX Designer / Product Builder** | Intelligent Agents, Conversational Design, Prototyping in code | Someone who bridges design and engineering for AI |

<!-- [CUSTOMIZE] Edit the archetypes above to match YOUR target roles. -->

### Adaptive Framing by Archetype

> **Concrete metrics: read from `cv.md` + `article-digest.md` at evaluation time. NEVER hardcode numbers here.**

| If the role is... | Emphasize about the candidate... | Proof point sources |
|-------------------|----------------------------------|---------------------|
| Platform / LLMOps | Production systems builder, observability, evals, closed-loop | article-digest.md + cv.md |
| Agentic / Automation | Multi-agent orchestration, HITL, reliability, cost | article-digest.md + cv.md |
| Technical AI PM | Product discovery, PRDs, metrics, stakeholder mgmt | cv.md + article-digest.md |
| Solutions Architect | System design, integrations, enterprise-ready | article-digest.md + cv.md |
| Forward Deployed Engineer | Fast delivery, client-facing, prototype to prod | cv.md + article-digest.md |
| AI Transformation Lead | Change management, team enablement, adoption | cv.md + article-digest.md |
| AI/UX Designer | Intelligent Agents, multi-turn interaction patterns, technical feasibility | cv.md + article-digest.md |

### Exit Narrative (use in ALL framings)

Use the candidate's exit story from `config/profile.yml` to frame ALL content:
- **In PDF Summaries:** Bridge from past to future -- "Now applying the same [skill] to [JD domain]."
- **In STAR stories:** Reference proof points from article-digest.md
- **In Draft Answers (Section G):** The transition narrative should appear in the first response.
- **When the JD asks for "entrepreneurial", "ownership", "builder", "end-to-end":** This is the #1 differentiator. Increase match weight.

### Cross-cutting Advantage

Frame profile as **"Technical builder with real-world proof"** that adapts framing to the role:
- For PM: "builder who reduces uncertainty with prototypes then productionizes with discipline"
- For Designer: "builder who ensures technical feasibility by prototyping with code and Figma API automation"
- For Engineer: "builder who delivers fast with observability and metrics from day 1"

Convert "builder" into a professional signal, not a "hobby maker". Real proof points make this credible.

### Portfolio as Proof Point (use in high-value applications)

If the candidate has a live demo/dashboard (check profile.yml), offer access in applications for relevant roles.

### Compensation Intelligence

**General guidance:**
- Use WebSearch for current market data (Glassdoor, Levels.fyi, Blind)
- Frame by role title, not by skills -- titles determine compensation bands
- Contractor rates are typically 30-50% higher than employee base to account for benefits
- Geographic arbitrage works for remote roles: lower cost of living = better net

### Negotiation Scripts

**Salary expectations (general framework):**
> "Based on market data for this role, I'm targeting [RANGE from profile.yml]. I'm flexible on structure -- what matters is the total package and the opportunity."

**Geographic discount pushback:**
> "The roles I'm competitive for are output-based, not location-based. My track record doesn't change based on postal code."

**When offered below target:**
> "I'm comparing with opportunities in the [higher range]. I'm drawn to [company] because of [reason]. Can we explore [target]?"

### Location Policy

**In forms:**
- Binary "can you be on-site?" questions: follow your actual availability from profile.yml
- In free-text fields: specify your timezone overlap and availability

**In evaluations (scoring):**
- Remote dimension for hybrid outside your country: score **3.0** (not 1.0)
- Only score 1.0 if JD explicitly says "must be on-site 4-5 days/week, no exceptions"

### Time-to-offer priority
- Working demo + metrics > perfection
- Apply sooner > learn more
- 80/20 approach, timebox everything

---

## Global Rules

### NEVER

1. Invent experience or metrics
2. Modify cv.md or portfolio files
3. Submit applications on behalf of the candidate
4. Share phone number in generated messages
5. Recommend compensation below market rate
6. Generate a PDF without reading the JD first
7. Use corporate-speak
8. Ignore the tracker (every evaluated offer gets registered)

### ALWAYS

0. **Mandatory Environment Check:** Silently verify `cv.md`, `config/profile.yml`, `portals.yml`, and `data/applications.md` exist at session start. Do NOT ask the user — just read them.
1. Read cv.md and article-digest.md (if exists) before evaluating any offer
1b. **First evaluation of each session:** Run `node cv-sync-check.mjs` with shell. If it reports warnings, notify the candidate before continuing.
2. Detect the role archetype and adapt framing (favor "Product Builder" / "UI/UX Engineer" for Sparsh).
3. Cite exact lines from CV when matching.
4. Use WebSearch for compensation and company data.
5. Register in tracker after evaluating.
6. Generate ALL content in English. All modes, reports, evaluations, messages, and tracker entries MUST be in English.
7. Be direct and actionable -- no fluff.
8. When generating English text (PDF summaries, bullets, LinkedIn messages, STAR stories): native tech English, not translated. Short sentences, action verbs, no unnecessary passive voice.
8b. **Case study URLs in PDF Professional Summary:** If the PDF mentions case studies or demos, URLs MUST appear in the first paragraph (Professional Summary). The recruiter may only read the summary. All URLs with `white-space: nowrap` in HTML.
9. **Tracker additions as TSV** -- NEVER edit applications.md to add new entries. Write TSV in `batch/tracker-additions/` and `merge-tracker.mjs` handles the merge.
10. **Include `**URL:**` in every report header** -- between Score and PDF.

### Tools

| Tool | Use |
|------|-----|
| web-search | Comp research, trends, company culture, LinkedIn contacts, broad job discovery |
| web-fetch | Extract JDs from static URLs, verify offers, read company pages. Fallback when Chrome DevTools MCP is unavailable |
| Chrome DevTools MCP | Navigate portals (`navigate_page`), read page content (`take_snapshot`), fill forms (`fill`, `fill_form`), click buttons (`click`), interact with SPAs. Preferred for scanning and applying |
| shell | Run Node.js scripts: `node generate-pdf.mjs`, `node merge-tracker.mjs`, etc. |
| File I/O | Read cv.md, article-digest.md, cv-template.html; Write reports, tracker TSVs, temp HTML |
