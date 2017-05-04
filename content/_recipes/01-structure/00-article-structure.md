---
title: Article Structure
---

Articles follow the natural alphabetic ordering of files and folders. 
Articles in sub directories use the folder name to group articles in the main menu.
The name of a category may be overridden using an `index.md` file 
The following example demonstrates how to order and organise articles and subsections:

```
    .
    ├── 00-article-structure.md
    ├── 01-article-2.1.md
    ├── 02-article-2.2.md
    ├── 02-article-2.3.md
    ├── Directory\ With\ Index
    │   ├── index.md
    │   └── with-index.md
    ├── Directory\ Without\ Index
    │   └── without-index.md
    └── index.md
```