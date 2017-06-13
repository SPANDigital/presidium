---
title: Embed
---

A fallback approach to import generated documentation is to embed documentation within an iframe.
This approach is not advised as items are not indexed or available on the main menu but 
will work for certain cases where an importer is not yet available.

> Where possible use a simple template when embedding documentation in an iframe

To included documentation in an iframe:
1. Generate your static site documentation for your component
1. Place the documentation within the `/media` folder so that it is statically served. 
The Presidium convention is to place it under `/media/import/{my-reference}`
1. Add a reference article to the reference section:

```markdown

---
title: My Reference
---

# foo.bar

<div>
    <iframe
            src='{{site.baseurl}}/media/import/{my-reference}/foo/bar/package-summary.html'
    </iframe>
</div>
```

>You can create multiple markdown files for different components as required.
