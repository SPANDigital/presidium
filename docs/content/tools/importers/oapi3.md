---
title: OpenAPI3
weight: 2
---

Presidium includes a Golang tool ([presidium-oapi3](https://www.npmjs.com/package/presidium-oapi-3) for importing your OpenAPI 3 spec into Presidium documentation.

1. Add the [presidium-oapi3](https://www.npmjs.com/package/presidium-oapi-3) dependency to your site's `package.json` or run `npm install --save presidium-oapi3`.
1. Add a script that invokes the tool.
1. Run `npm run import-oapi` whenever you need to update your API documentation.

```json
{
    "scripts" : {
        "import-oapi" : "presidium-oapi3"
    },
}
```

Example:

```sh
$ npm run import-oapi convert -f <YOUR_API_SPEC> -o <THE_OUTPUT_DIRECTORY> -r <THE_PRESIDIUM_REFERENCE_URL>
```

The following options are available for `presidium-oapi3`:

| Option | Description
|:-------|:---
|  -n, \--apiName `string`       | The name under which the generated docs will be grouped |
|  -f, \--file `string`         |  OpenAPI 3 spec file |
|  -o, \--outputDir `string`     | The output directory |
|  -r, \--referenceURL `string`  | The reference URL (default "reference")|
|  -h, \--help                 | help for convert |