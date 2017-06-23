---
title:  Generating Article HTML
---

It is sometimes useful to access just the article html, without the navigation menu or headers and footers. This allows articles to be easily embedded in other sites and systems. 

Each articles html may be generated as a seperate file when building a site.
To generated html for all articles in a section, the `export-articles` property of a section can be set in the `site config`:

```
sections:
  - title: "Glossary"
    url: "/glossary/"
    collection: glossary
    export-articles: true
```

The generated site will include html files for all articles in the glossary section: 

```
  .
  └── glossary
      ├── article-a.html
      ├── article-b.html
      ├── article-c.html
      └── index.html
```
