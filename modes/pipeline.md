# Mode: pipeline — Pending URL Processor

Processes URLs stored in `data/pipeline.md`. This is the manual version of `batch` mode for focused triage.

## Workflow

1. **Read Pipeline**: Extract all URLs from `data/pipeline.md`.
2. **Triage**: For each URL:
   - Show Company and Role title.
   - Ask: "Evaluate, Skip, or Save for later?"
3. **Action**:
   - **Evaluate**: Run `auto-pipeline`.
   - **Skip**: Move to `data/scan-history.tsv` with `DISCARDED` status.
   - **Save**: Keep in pipeline.

## Triage UI

```
[1/5] Company: Google | Role: Staff AI Designer
> Action? (e)valuate / (s)kip / (n)ext
```

## Clean up
Remove processed URLs from `data/pipeline.md` at the end of the session.
