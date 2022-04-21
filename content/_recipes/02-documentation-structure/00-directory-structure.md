---
title: Directory Structure
---

In Hugo, your content should be organized in a manner that reflects the rendered website.

While Hugo supports content nested at any level, the top levels (i.e. content/<DIRECTORIES>) are special in Hugo and are considered the content type used to determine layouts etc.

Without any additional configuration, the following will automatically work:

```
   └── content
    └── about
    |   └── index.md  // <- https://example.com/about/
    ├── posts
    |   ├── firstpost.md   // <- https://example.com/posts/firstpost/
    |   ├── happy
    |   |   └── ness.md  // <- https://example.com/posts/happy/ness/
    |   └── secondpost.md  // <- https://example.com/posts/secondpost/
    └── quote
        ├── first.md       // <- https://example.com/quote/first/
        └── second.md      // <- https://example.com/quote/second/
```

## Path Breakdown in Hugo 

The following demonstrates the relationships between your content organization and the output URL structure for your Hugo website when it renders. These examples assume you are using pretty URLs, which is the default behavior for Hugo. The examples also assume a key-value of baseURL = "https://example.com" in your site’s configuration file.
