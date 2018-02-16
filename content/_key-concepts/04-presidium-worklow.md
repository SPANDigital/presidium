---
title: Presidium Workflow
---

Presidium provides workflows templates and tools for building and managing a static technical documentation site.
As shown below, all user documentation for a system can be written in Markdown and built into a static site using
Presidium:

![Documentation Workflow]({{ site.baseurl }}/media/images/doc-workflow.png "Documentation Workflow")
*High-level Workflow and Dependencies*

# Content Owners

 - Write and structure their articles and user documentation in a `git` repo.
 - Include media assets (images, attachments) as required.
 - Include reference documentation sources as required.
 - Configure and publish their site.

# Presidium

  - Provides:
    - Base templates, styling, and themes
    - Common UI components
    - Build workflows via `npm`:
      - Import
      - Build
      - Serve
      - Publish
      - Version
      - Dependencies
    - [Jekyll](https://jekyllrb.com/) integration and workflows via `npm`
    - Reference documentation import:
      - Javadoc
      - Swagger
      - JSDoc
    - Linking and validation tools
