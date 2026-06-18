# Mode: batch — Automated Batch Processing

Mode for evaluating multiple job offers from `data/pipeline.md` or a local folder of JDs. Used to "clean the inbox" and identify the highest match scores.

## Input Sources

1. **`data/pipeline.md`**: List of URLs (one per line).
2. **`jds/` folder**: Local markdown/text files containing JDs.

## Workflow

```
1. READ    → Load the list of URLs or files.
2. FILTER  → Skip URLs already in data/scan-history.tsv or applications.md.
3. LOOP    → For each item:
   a. FETCH    → Extract JD (WebFetch or Chrome DevTools MCP).
   b. EVAL     → Execute Blocks A-D (Summary, Match, Strategy, Comp).
   c. SCORE    → Calculate 1-5 match score.
   d. REPORT   → Save summary report in reports/.
   e. TRACK    → Write temporary TSV in batch/tracker-additions/.
4. MERGE   → Run `node merge-tracker.mjs` to consolidate into applications.md.
```

## Batch Evaluation Prompt

When in batch mode, use a condensed version of the evaluation to save tokens:

- **Archetype**: Single word.
- **Score**: X.X/5.
- **Match**: Top 3 matching proof points.
- **Gaps**: Top 2 critical gaps.
- **Strategy**: 1-sentence "Go/No-Go" recommendation.

## Verification

For batch workers, use `web-fetch` and mark the report header with `**Verification:** unconfirmed (batch mode)` if uncertain. The user can verify manually later.

## Execution Command

Usually triggered via:
`gemini /careerforge:batch`
Or for high-volume automated processing:
`gemini --yolo /careerforge:batch`
