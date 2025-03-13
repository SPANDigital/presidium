---
title: Layouts and Styling
weight: 50
---

Presidium stores layouts and styling separately. This allows multiple docsets to use a single source for layouts and styling.

For example, this section of a `config.yml` file points to the repositories that store layouts and styling:
```
module:
  imports:
  - path: github.com/spandigital/presidium-styling-base
  - path: github.com/spandigital/presidium-layouts-base
```
<span style="color:purple">**Reviewers**: Would it be good to list the items that these repos must or should contain? What would they be?</span>