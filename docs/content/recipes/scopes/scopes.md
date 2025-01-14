---
title: Scopes
status: SHOULD BE DEPRECATING
---

Presidium allows content creators to scope the final site specific audiences. For example: _internal_, _public_, _customer_, or even something such as _head-office_. Presidium will take care of excluding/including only relevant articles and sections from the final site.

It is both possible to scope articles and the menu segments by adding the appropriate front matter.

To use scopes follow these 4 simple steps:

1. [Apply scope to articles](#scoping-articles).
2. [Apply scope to left side section menu](#scoping-section-menu)
3. [Declaring scopes](#declaring-applied-scopes)
4. [Build scoped Presidium](#building-presidium-site-with-scopes-or-without)

> It is important to realize that out of scope articles will not even appear in the final site.

## Scoping Rules

Presidium Hugo offers a very simple simple scoping rules, which are applied to both articles and menu configurations. The following table summarize them:

| Applied in `scopes.json` | Set on menu/article section       | Rendering | Rule                                                                  |
| ------------------------ | --------------------------------- | --------- | --------------------------------------------------------------------- |
| `[]`                     | `[internal, customer, developer]` | Y         | An empty rule set indicates that all rules have to applied.           |
| `[investor,corporate]`   | `[internal, customer, developer]` | N         | Neither _investor_, nor _corporate_ scopes applies to section/article |
| `[corporate,internal]`   | `[internal,developer]`            | Y         | _internal_ scoped, but also _developer_.                              |

## Scoping Articles

To scope articles start by adding front matter to your page:

```yaml
---
scope: [internal,customer]
---
```

Important to note the following

1. You may have multiple scopes.
2. Scopes must always be enclosed by squire (`[]`) brackets, even if there is only one scope.
3. You may only have front matter named `scope`.
4. Empty scope is allowed.

## Scoping Section Menu

It also possible to scope the left site menu separately by modifying the `menu.main` section in the `site/config.yaml` by injecting scope parameters directly on a section. For example:

```yaml
menu:
  main:
 - identifier: getting-started
      name: Getting Started
      url: /getting-started/
      weight: 4
      params: 
          scope: [private]  
```

## Declaring applied scopes

1. First make sure you have a `data` folder on your project. If not create one directly on under just beneath the `/content` folder. 
2. Within the data folder, create file called `scopes.json`.
3. Decided which scopes you want to apply, for example:

   ```json
    {
        "applied": [
            "internal",
            "develop",
            "customer",
            "investors"
        ],
        "enabled":false
    }
   ```

## Building Presidium site with scopes (or without)

1. Building the site with these scopes applied, just make sure the `"enabled":` flag in the `scopes.json` is set to `true`.
2. Run the presidium hugo to build the site:

   ```shell
   > presidium hugo
   ```

> To build the site with all scopes applied, just set `"enabled":false` in `data/scopes.json` configuration.

