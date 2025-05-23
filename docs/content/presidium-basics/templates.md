---
title: Module Templates and Demo Articles
weight: 30
---
Presidium provides five module templates. Each one includes the following, tailored to the particular use case:
* Layouts and styling 
* A suggested set of top-level sections
* Demo articles, many with suggested content
* Article archetypes for use with the [hugo new content](https://gohugo.io/commands/hugo_new_content/) command.

### Module Templates
Below are the module templates and their top-level sections.

**Default templates**
* Overview
* Key Concepts
* Prerequisites
* Getting Started
* Best Practices
* Reference
* Glossary
* Recipes
* Tools
* Updates

**Design template:**
* Introduction
* Design Principles
* Visual Elements
* Typography
* Components
* Motion
* Design Tokens
* Resources
* Accessibility
* Contributing
* Reference
* FAQs
* Updates

**Onboarding Template:**
* Organization Overview
* Solution Overview
* Technology Stack
* Tool Chain
* Dev Environment Setup
* Getting Started
* Reference
* Updates
  
**Requirements Template:**
* Overview
* Archetypes and Personas
* Entities and Relationships
* State Transitions
* Process Flows
* Capabilities
* Reference
* Updates

**Blog Template**
* Archive
* General Announcements
* Product News
* Project News
* Social Events

### Demo Articles
The demo articles always define the article's frontmatter, and in many cases also provide suggestions about the content. For example, the Onboarding module template includes a top-level section "Tool Chain" which contains eight demo articles, one of which is the following:
```
---
title: Version Control
author: author
weight: 2
---

This article should contain information about the version control system that tracks all changes to files and the user who made the change.

Examples include:

* Github
* Atlassian Bitbucket
* Atlassian Fisheye/Crucible
* Perforce
* Mercurial

Optionally, add links to the official websites.
```