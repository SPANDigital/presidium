# Requirements site

A Presidium-based requirements site, generated from the `requirements` template. The site captures the bones of a product spec — archetypes, entities, features, state transitions, process flows, capabilities, reference material, updates, and a glossary — and is structured so that AI agents can populate it from your discovery materials.

## Layout

```
content/        # The site you publish (organised 01–10)
discovery/      # Drop input materials here for agents to read
archetypes/     # `hugo new` scaffolds for archetype, persona, entity, scenario, story
layouts/        # Custom Hugo shortcodes/partials shipped with this template (interim — see below)
scripts/        # Helper Python utilities (e.g. snapshot-stats.py)
assets/         # SCSS overrides compiled by presidium-styling-base
static/         # Static assets that ship with the site
config.yaml     # Hugo config (menu, taxonomies, output formats, modules)
Makefile        # serve / build / refresh / stats targets
```

## Getting started

```bash
make serve         # build and serve locally with live reload
make stats         # regenerate data/discovery-stats.yaml from agent outputs (if present)
```

## Working with agents

1. Drop your discovery materials (Confluence/Notion exports, domain models, schema files, screenshots) into `discovery/`. The folder is intentionally flat.
2. Install the Presidium agent suite of your choice (e.g. via the Claude Code marketplace, or by copying agent definitions from a reference repo). This template does not ship `.claude/` or `agent-artifacts/` directly — bring your own.
3. Run agents to populate `content/`. Agents typically write structured outputs (YAML models, generated articles, session logs) into an `agent-artifacts/` directory you create at the site root, and the `make stats` target inventories that directory into `data/discovery-stats.yaml` for the dashboard shortcodes to render.

## Future work

The custom layouts under `layouts/` are inlined as a starting point. They will move to a dedicated Hugo module (`presidium-layouts-requirements`) so they can be versioned and shared across products. When that module ships, the `layouts/` folder is removed and a single `module.imports` entry in `config.yaml` takes its place.
