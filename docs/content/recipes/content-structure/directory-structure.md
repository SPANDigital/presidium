---
title: Directory Structure
weight: "1"
---
When you create a Presidium site using the CLI `init` command, Presidium creates the directory structure for the 
selected template.

Sections and articles are arranged using the `weight` key in the `front matter` of each file, and to specify section
level title and ordering you may use the `_index.md` file inside the directory containing a section.

The following is an example of how you can order and organize files and directories:

```
    .
    ├── article-1.md // Specify weight 1 here in front matter
    ├── Directory-2
    │   ├── article-2.1.md // Specify weight 1 here in front matter
    │   ├── article-2.2.md // Specify weight 2 here in front matter
    │   ├── _index.md // Specify weight 2 here in front matter, this will set `Directory-2` as the second item in the parent section
    ├── article-3.md // Specify weight 3 here in front matter
    ├── article-4.md // Specify weight 4 here in front matter
    └── _index.md // Specify weight 1 here in front matter
```

Please note that for `Directory-2/_index.md` file we specify `weight` as `2`, as in this case, this weight 
determines the ordering for the entire section with respect of its siblings (`article-1.md`, `article-2.md`, `article-3.md`, `article-4.md`)
