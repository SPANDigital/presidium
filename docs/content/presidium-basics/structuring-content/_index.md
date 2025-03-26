---
title: Structuring Content
weight: 35
---
The elements that structure Presidium content are:
* The `config.yml` file: The top level of structure. In it you 
  * [Set the order]({{< ref "presidium-basics/structuring-content/#ordering-items" >}}) to top-level sections
  * Assign URL paths
* Folders: Folders in the `content/` folder organize documentation into sections. Folders can themselves contain folders; however it's recommended not to have more than three or four levels of subfolders.
* `_index.md` Files: One of these is required in each folder, to specify the section's title and other metadata. If the folder is contained in another folder, _index.md [sets its order]({{< ref "presidium-basics/structuring-content/#ordering-items" >}}) relative to other items in the containing folder. It can also contain introductory text for its (sub)section.
* [Frontmatter]({{< ref "presidium-basics/structuring-content/#front-matter" >}}): This includes the `weight` attribute, which is the recommended way of [setting the order of articles]({{< ref "presidium-basics/#ordering-items" >}}).