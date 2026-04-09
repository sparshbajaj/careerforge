# Mode: tracker — Application Status Triage

View and update the status of current applications in `data/applications.md`.

## Commands

- **List**: Show recent applications and their current status.
- **Update**: Change status of an entry (e.g., `Evaluated` → `Applied`).
- **Health Check**: Run `node verify-pipeline.mjs` to check for broken links or inconsistent states.
- **Reminders**: Identify applications that haven't responded in > 7 days.

## Update Workflow

1. User provides ID or Company.
2. Ask for new status (use canonical states from `templates/states.yml`).
3. Update `data/applications.md` and the corresponding report header.
4. Suggest next mode (e.g., if status is `Interview`, suggest `deep` or `story-bank`).
