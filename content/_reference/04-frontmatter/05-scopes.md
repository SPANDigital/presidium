---
title: "Scopes"
---

Articles and sections can be marked as `internal` and/or `external` via the use of the `scope` tag in article `front matter` or the `site config`.


Different versions of the site may then be built using:
```
presidium build -s [internal|external]
```

The `front matter` looks like this:
```
---
scope: internal
---
```

Multiple scopes are used like this:

```
---
scope: [internal, external]
---
```

Sections can be assigned scope in the `site config` like so:

```
sections:
  - title: Internal API
    url: /internal-api/
    collection: internal-api
    scope: internal
    
  - title: Public Section
    url: /contact-info/
    collection: contact-info
    scope: [internal, external]
```

To show or hide scopes on articles in your generated site, use the following setting in the `site config`:

```
show:
    scope: true|false
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