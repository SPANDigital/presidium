---
title: Authors
weight: 2
---

Every article can include an author in the `frontmatter`:

```
---
author: {author name}
---
```

Optional links to authors can be enabled by setting a base URL in the `site config`:

```
params:
  author:
    external:
      url: https://github.com/
      newTab: false
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
