# Mode: apply — Live Application Assistant

Interactive mode for when the candidate is filling out an application form in Chrome. Reads what's on screen, loads the offer's previous context, and generates personalized answers for each form question.

## Requirements

- **Best with Chrome DevTools MCP**: Gemini controls a real Chrome browser via the Chrome DevTools MCP server. The candidate sees the browser and Gemini can read page content, fill forms, and interact with dropdowns.
- **Without Chrome**: The candidate shares a screenshot or pastes the questions manually.

## Workflow

```
1. DETECT    → Read active Chrome tab (snapshot/URL/title)
2. IDENTIFY  → Extract company + role from the page
3. SEARCH    → Match against existing reports in reports/
4. LOAD      → Read full report + Section G (if it exists)
5. COMPARE   → Does the role on screen match the evaluated one? If changed → warn
6. ANALYZE   → Identify ALL visible form questions
7. GENERATE  → For each question, generate personalized answer
8. PRESENT   → Show answers formatted for copy-paste
```

## Step 1 — Detect the offer

**With Chrome DevTools MCP:** Use `take_snapshot` to read the active page. Extract title, URL, and visible content. Use `list_pages` to find the right tab if multiple are open.

**Without Chrome:** Ask the candidate to:
- Share a screenshot of the form (Read tool reads images)
- Or paste the form questions as text
- Or say company + role for us to search

## Step 2 — Identify and search context

1. Extract company name and role title from the page
2. Search in `reports/` by company name (Case-insensitive Grep)
3. If match → load full report
4. If Section G → load previous draft answers as base
5. If NO match → warn and offer to run a quick auto-pipeline

## Step 3 — Detect changes in the role

If the role on screen differs from the evaluated one:
- **Warn the candidate**: "The role has changed from [X] to [Y]. Do you want me to re-evaluate or adapt the answers to the new title?"
- **If adapt**: Adjust answers to the new role without re-evaluating
- **If re-evaluate**: Execute full A-F evaluation, update report, regenerate Section G
- **Update tracker**: Change role title in applications.md if appropriate

## Step 4 — Analyze form questions

Identify ALL visible questions:
- Free text fields (cover letter, why this role, etc.)
- Dropdowns (how did you hear, work authorization, etc.)
- Yes/No (relocation, visa, etc.)
- Salary fields (range, expectation)
- Upload fields (resume, cover letter PDF)

Classify each question:
- **Already answered in Section G** → adapt existing answer
- **New question** → generate answer from report + cv.md

## Step 5 — Generate answers

For each question, generate the answer following:

1. **Report context**: Use block B proof points, block F STAR stories
2. **Previous Section G**: If a draft answer exists, use it as a base and refine
3. **"I'm choosing you" tone**: Same framework as auto-pipeline
4. **Specificity**: Reference something concrete from the visible JD on screen
5. **career-ops proof point**: Include in "Additional info" if there is a field for it

**Output format:**

```
## Answers for [Company] — [Role]

Based on: Report #NNN | Score: X.X/5 | Archetype: [type]

---

### 1. [Exact form question]
> [Answer ready for copy-paste]

### 2. [Next question]
> [Answer]

...

---

Notes:
- [Any observations about the role, changes, etc.]
- [Personalization suggestions the candidate should review]
```

## Step 6 — Post-apply (optional)

If the candidate confirms they sent the application:
1. Update status in `applications.md` from "Evaluated" to "Applied"
2. Update report's Section G with final answers
3. Suggest next step: `/career-ops outreach` for LinkedIn outreach

## Scroll handling

If the form has more questions than those visible:
- Ask the candidate to scroll and share another screenshot
- Or paste the remaining questions
- Process in iterations until the entire form is covered
