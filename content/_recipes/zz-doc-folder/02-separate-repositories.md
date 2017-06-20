---
title: Separate Repositories
author: torabisu
---

# Separate Repositories

If you wish to keep your documentation stored in a separate repository, but still use it within your project, you can use git submodules:

```bash

$ cd my-project
$ git submodule add https://github.com/my-company/my-project-docs docs
$ git submodule status
 6a1ed31b9cb215657a1bd4b4de6737c07b41c896 docs (heads/master)
```

