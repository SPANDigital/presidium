---
title: User Roles
weight: "3"
---

Various roles can be added to a site to allow readers to filter articles and menu items by a target audience.  

This optional feature can be enabled by defining user roles in the `site config`.

```
roles:
  label: "Show documentation for role"
  all: "All Roles"
  options:
    - "Business Analyst"
    - "Developer"
    - "Tester"
    - "..."
```


If a role is not specified, articles default to `roles.all`.

Articles can have one or more roles defined in `front matter`:

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

```
---
roles: Business Analyst
---
```

To show or hide roles on articles in your generated site, use the following setting in the `site config`:

```
show:
    roles: true|false
```
