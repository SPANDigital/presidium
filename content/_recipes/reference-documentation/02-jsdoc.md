---
title: Jsdoc
---

Based off of [Jsdoc](http://usejsdoc.org/), Presidium offers a template-based tool, [presidium-jsdoc](https://www.npmjs.com/package/presidium-jsdoc), for importing Javascript comments into your Presidium documentation.

To use:

1. Add the [presidium-jsdoc](https://www.npmjs.com/package/presidium-jsdoc) dependency to your site's `package.json`. Or run `npm install --save presidium-jsdoc`. 
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
| -d,--directory `path`                      | The path to the output directory in `./content` e.g. `./content/_reference/mydocs`.  \n
| -h,--help                                  | Shows this help.
| -p,--path `path`                           | The path from which static files are served e.g. `./media/import/mydocs`. Defaults to `./media/jsdoc/<title>`.
| -s,--sourcepath `path`                     | The path to the project's source.
| -t,--title `string`                        | The title of the output folder. Defaults to the directory name supplied with -d if no package information can be found.

