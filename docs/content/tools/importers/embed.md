---
title: Embed
weight: 5
---

A fallback approach to importing generated documentation is to embed documentation in an iframe.
This approach is not recommended because items are not indexed or available on the main menu. However, it
will work for certain cases when an importer is not yet available.

> When possible, use a simple template when embedding documentation in an iframe.

To include documentation in an iframe:
1. Generate the static site documentation for your component.
1. Put the documentation in the `/static` folder so that it's statically served.
The Presidium convention is to place it under `/static/import/{my-reference}`.
1. Add a reference article to the Reference section:

```markdown

---
title: My Reference
---

# foo.bar

<div>
    <iframe>
            src='{{% baseurl %}}/static/import/{my-reference}/foo/bar/package-summary.html'
    </iframe>
</div>
```

>You can create multiple Markdown files for different components as required.
