#

Presidium allows content creators to scope articles and sections to certain audience. For example: _internal_, _public_, _customer_. Presidium will take care of excluding/including only relevant articles and sections.

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


