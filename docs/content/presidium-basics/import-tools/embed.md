---
title: Embedding
weight: 60
---

A fallback approach to importing generated documentation is to embed documentation in an iframe, using the `iframe` shortcode supplied with Presidium.

When possible, use a simple template when embedding documentation in an iframe.

To include documentation in an iframe:
1. Generate the static site documentation for your component.
2. Put the documentation in the `/static` folder so that it’s statically served. The Presidium convention is to place it under `/static/import/{my-reference}`.
3. Add a reference article to the Reference section

```
    ---
    title: My Reference
    ---
    # foo.bar
    {{/< iframe src="/static/import/{my-reference}" >}}
```

>  **Note:** Ignore the slash character after the second curly bracket in the above example. It is only there to keep the shortcode from being interpreted.

You can create multiple Markdown files for different components as required.
