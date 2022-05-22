---
title: Authors
weight: "2"
---

Every article can include an author in the `front matter`:
```
---
author: {author name}
---
```

Optional links to authors can be enabled by setting a base `authors-url` in the `site config`:

```
external:
  authors-url: https://github.com/orgs/SPANDigital/people/
```

```
external:
  authors-url: https://bitbucket.org/
```

To hide or show authors on your generated site,
enable or disable the component in the `site config`:

```
show:
  author: true|false
```
