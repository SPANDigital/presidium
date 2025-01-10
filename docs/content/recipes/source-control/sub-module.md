---
title: Sub Module
weight: "3"
status: GOOD
---

If you want to store your documentation in a separate repository or share documentation between projects you can use submodules. Use the following steps to set up a submodule.

1. Create a new repository for you module. 
2. Create a `config.yml` file and add the following
   ```yaml
   module:
     mounts:
       - source: content
         target: content
   ```
3. Create a `content` directory and add your markdown files. E.g.
   ```
   ├── config.yml
   └── content
       └── glossary
           ├── _index.md
           └── link.md

   ```
4. Commit and push your changes
5. To use your submodule, add it to the `imports` section of your project's `config.yml` file. E.g.
   ```yaml
   module:
     imports:
     - path: <REPO_URL_OF_SUBMODULE>
       mounts:
       - source: content
         target: content
   ```
