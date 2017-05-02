---
title: "Head"
---

Articles follow the natural alphabetic ordering of files and folders. 
Articles in sub directories use the folder name to group articles in the main menu.
The name of a category may be overridden using an `index.md` file 
The following example demonstrates how to order articles and categories:

```
    ./
    ├── 00-head.md
    ├── 01-sub.md
    ├── 02-sub.md
    ├── 03-category
    │   ├── 01-subcat.md
    │   └── index.md
    ├── 04-sub.md
    ├── 05-tail.md
    └── index.md
```