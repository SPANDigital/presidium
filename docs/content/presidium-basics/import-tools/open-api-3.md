---
title: OpenAPI 3
weight: 10
---
OpenAPI3

Presidium includes a Golang tool (presidium-oapi3) for importing your OpenAPI 3 spec into Presidium documentation.
1. Add the presidium-oapi3 dependency to your site’s `package.json` or run `npm install --save presidium-oapi3`.
1. Add a script that invokes the tool.
1. Run npm `run import-oapi` whenever you need to update your API documentation.

```
{
   "scripts" : {
       "import-oapi" : "presidium-oapi3"
   },
}
```
Example:
```
$ npm run import-oapi convert -f <YOUR_API_SPEC> -o <THE_OUTPUT_DIRECTORY> -r <THE_PRESIDIUM_REFERENCE_URL>
```
The following options are available for presidium-oapi3:

| Option | Description |
|--------|-------------|
| -n, --apiName  `string` | The name under which the generated docs will be grouped |
| -f, --file  `string` | The OpenAPI 3 spec file |
| -o, --outputDir  `string` | The output directory |
| -r, --referenceURL  `string` | The reference URL (default is `reference`) |
| -h, --help | Help for using the tool |