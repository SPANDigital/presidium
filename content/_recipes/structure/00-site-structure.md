---
title: Site Structure
---

Articles follow the natural alphabetic ordering of files and folders which may be grouped in sub-sections.

The following example demonstrates how you could order and organise articles and sub-sections using numeric prefixes:

```
    .
    ├── 01-article-2.1.md
    ├── 02-article-2.2.md
    ├── 02-Directory-2.3
    │   ├── 01-article-2.3.1.md
    │   ├── 01-article-2.3.2.md
    │   ├── index.md
    ├── 03-article-2.4.md
    └── index.md
```

The title of a sub-section is derived from the folder name, but may be set by providing 
an `index.md` file in the folder with the following `front matter`:

```
---
title: Sub Section Heading
---
```

>A maximum of four levels is supported by the main menu.
