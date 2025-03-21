---
title: Docset Templates and Article Archetypes
weight: 30
---
Presidium provides five docset templates: Default, Blog, Design, Onboarding, and Requirements. Each one includes the following, tailored to the particular use case:
* Layouts and styling 
* A suggested set of top-level sections
* Archetypes for suggested articles

The archetypes always define the article's frontmatter, and in many cases provide suggestions about the content.

For example, the Onboarding docset template includes a top-level section "Tool Chain", and that contains eight article archetypes, one of which is the following:
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
Below are the docset templates and their top-level sections.

Default and Blog templates have the same set of sections:
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

Design template:
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

Onboarding Template:
* Organization Overview
* Solution Overview
* Technology Stack
* Tool Chain
* Dev Environment Setup
* Getting Started
* Reference
* Updates
  
Requirements Template:
* Overview
* Archetypes and Personas
* Entities and Relationships
* State Transitions
* Process Flows
* Capabilities
* Reference
* Updates