---
title: Swagger
weight: "2"
---

Presidium includes a Java-based tool ([presidium-swagger](https://www.npmjs.com/package/presidium-swagger), based on [Swagger2Markup](https://github.com/Swagger2Markup/swagger2markup)) to import your API's Swagger into your Presidium documentation.

1. Add the [presidium-swagger](https://www.npmjs.com/package/presidium-swagger) dependency to your site's `package.json` or run `npm install --save presidium-swagger`.
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

The following options are available for `presidium-swagger`:

| Option | Description
|:-------|:---
| -d,--directory `path`                      | The destination directory in which to save the generated documentation. (Default is './docs'.)
| -h,--help                                  | Displays this help.
| -s,--sourcepath `path`                     | Swagger source path.
| -t,--title `string`                        | The title of your docs folder. (Default is directory name supplied with '-d'.)
| -u,--sourceurl  `url`                      | URL to your Swagger Json file.
