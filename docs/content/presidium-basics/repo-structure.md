---
title: Repository Structure
weight: 20
---
At minimum (default installation?) a Presidium repo contains the following

### Items You'll Frequently Access
**Folders**
* content: Holds the Markdown files that are the content of the docsite, and the  subfolders that contribute to determining the docsite structure.
* static: Contains assets, like images, that are integral to the site but don't change frequently.

**File:** [config.yml]({{< ref "presidium-basics/structuring-content/#config.yml" >}}): Defines the docset's navigational structure, along with global site settings like the docset title, base URL, layouts, and styling.

### Items You'll Access Less
This list is not exhaustive.

**Folders**
* build: Contains additional scripts that must run when the docset is built
* public: Created when you build the docset ("the published website, generated when you run the hugo or hugo server command")
* resources: Created when you build the docset. "contains cached output from Hugo’s asset pipelines, generated when you run the hugo or hugo server commands."

**Files**
* .gitignore
* go.mod
* go.sum
* LICENSE
* README.md
* rio.yml

<!--<span style="color:purple">**Reviewers:** -->