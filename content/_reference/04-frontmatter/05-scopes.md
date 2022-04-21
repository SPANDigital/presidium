---
title: "Scopes"
---

Articles and sections can be marked as `internal` and/or `external` via the use of the `scope` tag in article `front matter` or the `config.yaml`.


To enable it, use the following settings:

```
params:
  scopesEnabled: true
```

The `front matter` looks like this:
```
---
title: Welcome to SPAN!
weight: 2
slug: welcome
url: onboarding/welcome
id: _onboarding/01-welcome.md
scope: [internal]
---
```

Multiple scopes are used like this:

```
---
title: Welcome to SPAN!
weight: 2
slug: welcome
url: onboarding/welcome
id: _onboarding/01-welcome.md
scope: [internal, external]
---
```

Sections can be assigned scope in the `config.yaml` like so:

```
menu:
  Main:
  - identifier: Introduction
    name: Introduction
    url: /
    weight: 1
    params:
      scope: [internal]
```

Articles without scope will inherit from their section scope, while articles with scope will be unaffected.
The table below shows how these various interactions work, and what is and isn't included in the final site in each case.

<table>
<tr>
    <th>Build Scope</th>
    <th>Section</th>
    <th>Category</th>
    <th>Article</th>
</tr>
<tr>
    <td>unspecified</td>
    <td>
        <ul>
            <li>undefined : ✓</li>
            <li>internal : ✓</li>
            <li>external : ✓</li>
        </ul>
    </td>
    <td>✓</td>
    <td>✓</td>
</tr>
<tr>
    <td>internal</td>
    <td>
        <ul>
            <li>undefined : <code>if section.scope.includes("internal")</code></li>
            <li>internal : ✓</li>
            <li>external : ✗</li>
        </ul>
    </td>
    <td><code>if internal_articles > 0</code></td>
    <td>
        <ul>
            <li>undefined : <code>if section.scope.includes("internal")</code></li>
            <li>internal : ✓</li>
            <li>external : ✗</li>
        </ul>
    </td>
</tr>
<tr>
    <td>external</td>
    <td>
        <ul>
            <li>undefined : <code>if section.scope.includes("external")</code></li>
            <li>internal : ✗</li>
            <li>external : ✓</li>
        </ul>
    </td>
    <td><code>if external_articles > 0</code></td>
    <td>
        <ul>
            <li>undefined : <code>if section.scope.includes("external")</code></li>
            <li>internal : ✓</li>
            <li>external : ✗</li>
        </ul>
    </td>
</tr>
</table>