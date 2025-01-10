---
title: User Roles
slug: user-roles
url: front-matter/user-roles
weight: 3
status: SHOULD BE DEPRECATING
---

Various roles can be added to a site to allow readers to filter articles and menu items by a target audience.

This optional feature can be enabled by defining user roles in the `site config`.

```
params:
  roles:
    label: "Show documentation for"
    all: "All Roles"
    options:
      - "Business Analyst"
      - "Developers"
      - "Testers"
```

If a role is not specified, articles default to `roles.all`

Articles can have one or more roles defined in `front-matter`:

```
---
roles: [Developer, Business Analyst]
---
```

```
---
roles: [Developer]
---
```

To show or hide roles on articles in your generated site, use the following setting in the `site config`:

```
params:
  show:
    roles: true|false
```
