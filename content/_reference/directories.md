---
title: Directories
---

Articles are stored in the content directory, with associated resources like images stored in media. The package.json and _config.yml files are used to configure the project.
./_config.yml General options to configure the project
./content Articles
./dist Rendered site data.
./media Various resources for the project, images, imported content etc.
./package.json Configuration settings.
All content changes are watched, and any changes trigger a regeneration of the content. If there is a structural change, i.e. adding a sub-directory, or a new article, then a restart is required:


