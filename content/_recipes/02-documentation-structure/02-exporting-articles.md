---
title:  Generating Article HTML
---

It can be useful to access an article's html without the navigation menu or headers and footers. This enables you to easily embed articles in other sites and systems.

When building a site, you can set whether the html for each article is generated separately.
To generate html for all articles in a section, set the `export-articles` property in the `site config`:

```
sections:
  - title: "Glossary"
    url: "/glossary/"
    collection: glossary
    export-articles: true
```

In this example, when `export-articles` is set to true, the generated site will include html files for all articles in the Glossary section:

```
  .
  └── glossary
      ├── article-a.html
      ├── article-b.html
      ├── article-c.html
      └── index.html
```
