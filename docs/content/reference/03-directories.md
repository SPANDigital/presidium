---
title: Directories
slug: directories
url: reference/directories
weight: 3
---

Articles are stored in the content directory. Associated resources, such as images, are stored in the media directory.
[TODO] // Check if the information about resources is still accurate
The `config.yaml` file is used to configure the project.

| **Directory** | **Description** |
|---------------|-----------------|
| `./config.yaml` | General options to configure the project |
| `./content` | Articles |
| `./static/media` | Various resources for the project (images, imported content, etc.) |

All content changes are monitored; any change triggers a regeneration of the content.
