---
title: Hosting as a Sub-module
author: torabisu
---

If you want to store your documentation in a separate repository but still use it within your project, you can use Git submodules:

```bash

$ cd my-project
$ git submodule add https://github.com/my-company/my-project-docs docs
$ git submodule status
 6a1ed31b9cb215657a1bd4b4de6737c07b41c896 docs (heads/master)
```
