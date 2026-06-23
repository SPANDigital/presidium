# Presidium Project Module Template

This is the SPAN Digital project module template for Presidium. Each SPAN client project should have its own instance of this module, which serves as the single source of truth for project context — for team members, stakeholders, and AI skills alike.

---

## Purpose

The project module provides a structured, navigable knowledge base for a SPAN engagement. It captures:

- Project overview, goals, scope, and key documents
- Plans, milestones, and risks
- Ongoing status and sprint reviews
- Solution design (domain model, architecture, UX, technical design)
- Team structure and resource allocation
- Technology decisions and dependencies
- Discovery material and client documents

---

## How to Use This Template

### 1. Instantiate the module

Copy this template to create a new Presidium module for your project. Rename the module to match the client project (e.g. `acme-project-module`).

Run the `project-setup-skill` (when available) to pre-populate static fields — client name, project type, key contacts, and integration identifiers. Until the skill is available, search for `REPLACE WITH` across all files and fill in manually.

### 2. Fill in content

Each file contains instructional placeholder text explaining what belongs there and where the data should come from. Follow the instructions in each article. The system of record for each section is noted at the top of every folder's index page.

### 3. Connect your tools

For AI skills and the corpus to work, the following SPAN connectors (or equivalent MCPs) must be connected to this module's Cowork session:

| Connector | Purpose |
|---|---|
| **GWS** (Google Workspace) | Overview, Plan, Risks & Issues, Discovery, Sprint Reviews |
| **ODS** (Operational Data Store) | Team & Roles, Allocation, Status |
| **Jira** | Risks & Issues, data/jira-epics.json |
| **Pipedrive** | Status |
| **Harvest** | Budget, Allocation (via ODS) |
| **Presidium** | Publishing and module management |
| **ADSS** | Technology and architecture documentation |

### 4. Run the required plugins

The following plugins populate specific sections of this module:

| Plugin | Section | Status |
|---|---|---|
| Domain Modeling Plugin | Solution > Domain Model | Available |
| Requirements & Testing Framework | (cross-cutting) | Available |
| Project Management & Status Reporting Plugin | Status, Plan | To be built |
| Architecture Plugin | Solution > Architecture | To be built |
| Software Design Plugin | Solution > Technical Design | To be built |
| UX Design Plugin | Solution > Experience Design | To be built |

Sections managed by a plugin will note this at the top of their index page. Do not edit plugin-managed content manually — it will be overwritten when the plugin next runs.

---

## The Corpus

When all connectors and plugins are active, this module is populated automatically from a **corpus** — a curated set of project data sources that AI skills read from and write to. The corpus includes:

### Google Workspace documents (via GWS MCP)

Documents associated with this project in Google Workspace, identified by label:

| Label | Maps To |
|---|---|
| **Project Brief** | Overview |
| **Estimation** | Overview, Plan |
| **Plan** | Plan (via ODS) |
| **Status** | Status (via GWS and Pipedrive) |
| **Contract** | Overview > Key Documents |
| **Deliverable** | Overview > Scope & Deliverables |
| **Design Artifact** | Solution > Design (via GWS or Miro) |
| **Meeting Notes** | Discovery (if relevant) |
| **General** | Discovery (if relevant) |

### Structured data (via connectors)

| Source | Data | Location in Module |
|---|---|---|
| Jira | Epics | `data/jira-epics.json`, Plan, Status |
| Harvest / ODS | Team members, time logs | `data/team-members.json`, Team, Budget |
| Pipedrive / ODS | Status reports | `data/status-reports.json`, Status |

---

## Folder and File Conventions

- **Numeric prefixes** on all folders and files (e.g. `01-overview/`, `02-plan/`) control display order in GitHub. They do not appear in rendered URLs — each file sets its own `url` in front matter.
- **`_index.md`** files define section headings. Each includes a "System of record" note explaining where the section's content originates.
- **`REPLACE WITH`** placeholders mark fields that require human input during project setup.
- **Plugin-managed sections** are noted in their `_index.md`. Do not edit these manually.

---

## Template Structure

```
templates/project/
├── README.md                        ← this file
├── config.yml                       ← Presidium/Hugo configuration
├── Makefile
├── go.mod.tmpl
├── data/                            ← Hugo data files (populated by AI skills)
│   └── README.md                    ← describes jira-epics.json, team-members.json, status-reports.json
├── discovery/                       ← client discovery documents (not a Presidium section)
│   └── README.md
└── content/
    ├── 01-overview/
    ├── 02-plan/
    │   └── 03-risks-and-issues/
    ├── 03-status/
    │   └── 06-sprint-reviews/
    ├── 04-solution/
    │   ├── 01-domain-model/
    │   │   ├── 02-logical-entity-model/
    │   │   └── 03-functional-decomposition/  ← Functional Areas, Features, Functions, Sub-functions
    │   ├── 02-architecture/
    │   │   ├── 01-system-context/
    │   │   └── 02-conceptual-architecture/
    │   └── 03-design/
    │       ├── 01-experience-design/
    │       └── 02-technical-design/
    ├── 05-team/
    └── 06-technology/
```
