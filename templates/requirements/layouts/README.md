# Layouts (interim)

These layouts are inlined here as a starting point so a freshly generated requirements site renders end-to-end without external dependencies.

The plan is to extract them into a dedicated Hugo module — e.g. `github.com/spandigital/presidium-layouts-requirements` — alongside the existing `presidium-layouts-base` and `presidium-styling-base` modules. When that module ships, this directory is deleted and replaced with a single `module.imports` entry in `config.yaml`.

Until then, treat this folder as the source of truth for requirements-specific shortcodes (`entity-table`, `persona-details`, `archetype-details`, `discovery-stats`, `session-log`, etc.), partials, and section/taxonomy templates.

If you need to override a specific layout for your product, drop your version into your generated site's own `layouts/` directory — Hugo's lookup order means your override wins without touching this folder.
