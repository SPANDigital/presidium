---
title: Repository Structure
weight: 20
---
A Presidium repository typically contains the following:

**Items You'll Frequently Access**

Folders:
* `content`: Holds the Markdown files that are the content of the docsite, and the  subfolders that contribute to determining the docsite structure
* `static`: Contains assets, like images, that are integral to the site but don't change frequently

File: [config.yml]({{< ref "presidium-basics/structuring-content/#config.yml" >}}): Defines the module's navigational structure, along with global site settings like the module title, base URL, layouts, and styling.

**Items You'll Access Less**

This list is not exhaustive.

Folders:
* `build`: Contains additional scripts that must run when the module is built.
* `public`: Created when you build the module. It is the published website.
* `resources`: Created when you build the module; contains cached output from Hugo’s asset pipelines.
* `data`: Holds data drawn on by articles.

Files:
* `go.mod`: Contains layouts and styling configuration
* `go.sum`: Generated from the go.mod file
* `README.md`: Information for contributors
  
… and any other files you want to include, such as `.gitignore`, license, and so on.
<!--<span style="color:purple">**Reviewers:** -->