# data/

Hugo data directory for this project module. Data files placed here are available to Hugo templates and AI skills via structured data access.

## Intended contents

| File | Source | Description |
|---|---|---|
| `jira-epics.json` | Jira (via Jira connector) | Current Jira epics for this project |
| `team-members.json` | Harvest / ODS | Team members and roles |
| `status-reports.json` | Pipedrive / ODS | Status report summaries |

These files are populated by AI skills and connectors. Do not edit manually — changes will be overwritten on the next skill run.
