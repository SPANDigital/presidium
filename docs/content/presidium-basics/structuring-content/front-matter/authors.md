---
title: Authors
weight: 20
author: helgard.meyer@spandigital.com
status: review
---

Every article can include an author in the frontmatter:

```
---
author: {author name}
---
```

To customize the label for authors, provide a new label in the `site config`:

```
params:
  author:
    label: Custom Label
```

To hide or show authors on your generated site, enable or disable the component in the `site config`:

```
params:
  show:
    author: true|false
```