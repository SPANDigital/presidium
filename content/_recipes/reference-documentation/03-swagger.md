---
title: Swagger
---

Based off of [Swagger2Markup](https://github.com/Swagger2Markup/swagger2markup), Presidium offers a Java-based tool, [presidium-swagger](https://www.npmjs.com/package/presidium-swagger), for importing your API's Swagger into your Presidium documentation.

To use:

1. Add the [presidium-swagger](https://www.npmjs.com/package/presidium-swagger) dependency to your site's `package.json`. Or run `npm install --save presidium-swagger`.
1. Add a script that invokes the tool.
1. Run `npm run import-swagger` whenever you need to update your API documentation.

```json
{
    "scripts" : {
        "import-swagger" : "presidium-swagger"
    },
    "devDependencies": {
        "presidium-core" : "#.#.#",
        "presidium-jsdoc" : "#.#.#"
    }
}
```

Example:

```sh
$ npm run import-swagger -- -u <url> -d <path> -t <string>
```

The following options are available to `presidium-swagger`:

| Option | Description
|:-------|:---
| -d,--directory `path`                      | The destination directory to save the generated documentation to, defaults to: './docs'.
| -h,--help                                  | Shows this help.
| -s,--sourcepath `path`                     | Swagger source path.
| -t,--title `string`                        | Title of your docs folder, defaults to directory name supplied with '-d'.
| -u,--sourceurl  `url`                      | URL to your Swagger Json file.
