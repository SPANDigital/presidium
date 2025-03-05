---
title: Directory Structure
weight: "1"
---

When you create a Presidium site using the CLI `init` command, Presidium creates the directory structure for the 
selected template.

Sections and articles by default are ordered by file path. To make it easier to track the order of articles you can prefix your filenames and directories (not `_index.md` files) with a number, for example, `01-article.md`.

> **Note**: The number is added to the URL, for example, `reference/01-section`. To remove this add `url: reference/section` to the frontmatter of the section's `_index.md` file.

The main sections (for example, Reference and Overview) are ordered by their `weight` value in the `config.yml`. For more information on `weight`, see below:


### Sort Using Weight

In the project's `conifg.yml`, under `params:`,  change `sortByFilePath: true` to `sortByFilePath: false` to disable sorting by file path. 

```yml
params:
    sortByFilePath: false
```

Sections and articles can be arranged using the `weight` key in the `frontmatter` of each file. For specifying section level titles and ordering use the `_index.md` file inside the directory for that section. Should a `weight` not be specified and `sortByFilePath: false` Hugo will fall back to the following to order content: Data > Link Title > FilePath. For more information, see [Order Content](https://gohugo.io/templates/lists/#order-content) in the Hugo documentation.

The following is an example of how you can order and organize files and directories:

```
    .
    ├── article-1.md // Specify weight 1 here in frontmatter
    ├── Directory-2
    │   ├── article-2.1.md // Specify weight 1 here in frontmatter
    │   ├── article-2.2.md // Specify weight 2 here in frontmatter
    │   ├── _index.md // Specify weight 2 here in frontmatter, this will set `Directory-2` as the second item in the parent section
    ├── article-3.md // Specify weight 3 here in frontmatter
    ├── article-4.md // Specify weight 4 here in frontmatter
    └── _index.md // Specify weight 1 here in frontmatter
```

> **Note**: For `Directory-2/_index.md` file we specify `weight` as `2`, as in this case, this weight determines the ordering for the entire section with respect of its siblings (`article-1.md`, `article-2.md`, `article-3.md`, `article-4.md`)
