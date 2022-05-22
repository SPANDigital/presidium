---
title: Sub Folder
slug: sub-folder
url: source-control/sub-folder
weight: 2
---

Presidium can exist within a sub-folder of an existing source code repository, for example, `/docs`.

## Getting Presidium

The easiest way to incorporate Presidium into your project is to run the wizard from your project root:

```shell
presidium-hugo init
```

And specify `Project Name` as `docs`, so that, Presidium creates and empty Presidium site under `docs/` folder.

The contents of your `docs/` folder should look something like this:

```
config.yaml
go.mod
content/
static
```

Add the following to your project’s .gitignore file:

```
docs/public
```

From this point on, you can follow instructions in the Getting Started section. The only difference is that your 
documentation root is `/docs`.
