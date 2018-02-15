---
title: Jsdoc
---

Presidium includes a template-based tool ( [presidium-jsdoc](https://www.npmjs.com/package/presidium-jsdoc), based off of [Jsdoc](http://usejsdoc.org/)) to import Javascript comments into your Presidium documentation.

1. Add the [presidium-jsdoc](https://www.npmjs.com/package/presidium-jsdoc) dependency to your site's `package.json` or run `npm install --save presidium-jsdoc`.
1. Add a script that invokes the tool.
1. Run `npm run import-jsdoc` whenever you need to update your API documentation.

```json
{
    "scripts" : {
        "import-jsdoc" : "presidium-jsdoc"
    },
    "devDependencies": {
        "presidium-core" : "#.#.#",
        "presidium-jsdoc" : "#.#.#"
    }
}
```

Example:

```sh
$ npm run import-jsdoc -- -s <path> -d <path> -t <string> -p <path>
```

The following options are available to `presidium-jsdoc`:

| Option | Description
|:---|:---
| -d,--directory `path`                      | The path to the output directory in `./content` (for example, `./content/_reference/mydocs`).
| -h,--help                                  | Displays this help.
| -p,--path `path`                           | The path from which static files are served (for example, `./media/import/mydocs`). Default is `./media/jsdoc/<title>`.
| -s,--sourcepath `path`                     | The path to the project's source.
| -t,--title `string`                        | The title of the output folder. Default is the directory name supplied with -d if no package information is found.
