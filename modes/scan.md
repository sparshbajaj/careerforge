# Mode: scan — Portal Job Scanner

Scan pre-configured job boards and company career pages for new roles matching your profile.

## Configuration

Uses `portals.yml` for queries, filters, and company list.

## Workflow

1. **Execute Search**: Use `web-search` or Chrome DevTools MCP to hit portals (LinkedIn, Indeed, Otta, company pages).
2. **Scrape**: Extract titles, companies, and URLs.
3. **Filter**:
   - **Positive**: Keywords like "AI", "Agent", "Product Designer", "UX Engineer".
   - **Negative**: Keywords like "Intern", "Junior", "Java", "PHP".
   - **History**: Skip URLs in `data/scan-history.tsv` or `data/applications.md`.
4. **Save Results**: Append new valid URLs to `data/pipeline.md`.

## Output

Summary of new roles found:
```
🔍 Scan Complete:
- Found: 45 total roles
- Filtered: 38 (duplicates or low match)
- Added to pipeline: 7 new roles
```
Suggest running `/career-ops:pipeline` next.
