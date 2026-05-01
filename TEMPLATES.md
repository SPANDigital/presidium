# Templates Folder Overview

The `templates` folder contains Hugo-based documentation templates for different types of documentation sites. It includes **6 different template types**, each with pre-configured structure, content, and navigation.

## Template Types

### 1. Blog Template (`templates/blog`)
Blog/news site template with categorized content:
- General announcements
- Product news
- Project news
- Social events
- Archive functionality

### 2. Default Template (`templates/default`)
General documentation template with standard sections:
- Overview, Key Concepts, Prerequisites
- Getting Started, Best Practices
- Reference, Glossary, Recipes, Tools, Updates

### 3. Design Template (`templates/design`)
Design system documentation template featuring:
- Design principles & branding
- Components (buttons, modals, navigation, etc.)
- Visual elements (color, icons, data visualization)
- Design tokens, Typography, Motion
- Accessibility guidelines

### 4. Requirements Template (`templates/requirements`)
Software requirements documentation template, structured for AI agent workflows:

- **Content sections** (`content/`):
  - 01-overview
  - 02-archetypes-and-personas
  - 03-entities-and-relationships
  - 04-features-and-functions
  - 05-state-transitions
  - 06-process-flows
  - 07-capabilities
  - 08-reference
  - 09-updates
  - 10-glossary
  - 11-dashboard (live `{{< discovery-stats >}}` view of artifact counts and trends)
- **`discovery/`** — flat folder for input materials agents read (Confluence exports, domain models, schemas, screenshots).
- **`archetypes/requirements/`** — Hugo `hugo new` scaffolds for archetype, persona, entity, scenario, and story content types.
- **`layouts/`** — custom Hugo shortcodes/partials for the requirements site (entity-table, persona-details, archetype-details, discovery-stats, session-log, scenario-dependencies, etc.). *Interim — planned to move to a dedicated `presidium-layouts-requirements` Hugo module.*
- **`scripts/snapshot-stats.py`** — generates `data/discovery-stats.yaml` time-series from agent outputs.
- **`assets/_sass/client-custom/_client-custom.scss`** — overrides the empty stub in `presidium-styling-base` to style the requirements-specific shortcodes (scenario cards, persona panels, entity tables, etc.).
- **`static/`** — site-level static assets.

### 5. Onboarding Template (`templates/onboarding`)
Team/project onboarding documentation covering:
- Getting started
- Organization & Solution overview
- Dev environment setup
- Technology stack & Tool chain

### 6. Runbook Template (`templates/runbook`)
Operations/SRE documentation template with:
- Overview & Key Concepts
- Standard Operating Procedures
- Troubleshooting guides
- Response procedures
- SLAs & Glossary

## Common Structure

Each template typically contains:

- **`config.yaml/yml`** - Hugo configuration with menu structure and theme settings
- **`content/`** - Markdown content organized by sections
- **`archetypes/`** - Content templates for new pages
- **`static/`** - Static assets like images
- **`go.mod.tmpl`** - Templated Go module file (substitutes `{{ .ProjectName }}` at generation time)
- **`Makefile`** - Standard `serve`, `hugo`, `refresh`, `tidy` targets

## How templates are processed

When `presidium init` runs, every file under the chosen template is walked and either:

- **Templated** through Go's `text/template` (substituting `{{ .Title }}`, `{{ .ProjectName }}`, etc.) — the default for `content/`, `archetypes/`, `config.yaml`, `go.mod.tmpl`.
- **Copied verbatim** — applied to files under `layouts/`, `static/`, `data/`, `scripts/`, `assets/`, and `discovery/`. These paths typically contain Hugo template syntax (`{{ ... }}`) that must not be evaluated by Go's templater.

See `pkg/domain/service/template/template.go` for the full list of verbatim path prefixes.

## Technology

The templates use:
- **Hugo** static site generator
- **Presidium** styling and layout modules:
  - `presidium-styling-base`
  - `presidium-layouts-base`

## Usage

These templates serve as starting points for creating different types of documentation sites, each optimized for specific use cases and content organization patterns.
