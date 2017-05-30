---
title: Exporting Article HTML
---

The raw html of an article may be exported when building a site.
To export the generated html for all articles in a section, the `export-articles` property of a section can be set in the `site config`:

```
sections:
  - title: "Glossary"
    url: "/glossary/"
    collection: glossary
    export-articles: true
```

This generated static site will include html files for all articles: 

```
├── glossary
│   ├── article-a.html
│   ├── article-b.html
│   ├── article-c.html
│   └── index.html
```
