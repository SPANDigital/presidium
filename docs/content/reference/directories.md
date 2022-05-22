---
title: Directories
weight: "4"
---

Articles are stored in the content directory.  Associated resources, such as images, are stored in the media directory. The `package.json` and `_config.yml` files are used to configure the project.

| Directory            | Description                                                      |
|----------------------|------------------------------------------------------------------|
| `./_config.yml`      | General options to configure the project                         |
| `./content`          | Articles                                                         |
| `./dist`             | Rendered site data and staging area for source files              |
| `./media`            | Various resources for the project (images, imported content, etc.) |
| `./package.json`     | Configuration settings                                          |

All content changes are monitored; any change triggers a regeneration of the content. 
