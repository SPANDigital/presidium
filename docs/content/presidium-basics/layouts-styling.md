---
title: Layouts and Styling
weight: 50
---

Presidium stores layouts and styling separately. This allows multiple modules to use a single source for layouts and styling. Typically you won't need to modify them.

For example, this section of a `config.yml` file points to the repositories that store layouts and styling:
```
module:
  imports:
  - path: github.com/spandigital/presidium-styling-base
  - path: github.com/spandigital/presidium-layouts-base
```