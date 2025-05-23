---
title: OpenAPI 3
weight: 10
---

Presidium includes a Golang tool ([presidium-oapi3](https://github.com/SPANDigital/presidium-oapi3)) for importing your OpenAPI 3 spec into Presidium documentation.

Install the [presidium-oapi3](https://github.com/SPANDigital/presidium-oapi3) tool using Homebrew:

```sh
   brew tap SPANDigital/homebrew-tap https://github.com/SPANDigital/homebrew-tap.git
   brew install presidium-oapi3
```

Example:

```sh
$ presidium-oapi3 convert -f <YOUR_API_SPEC> -o <THE_OUTPUT_DIRECTORY> -r <THE_PRESIDIUM_REFERENCE_URL>
```

The following options are available for `presidium-oapi3`:

| Option                       | Description                                             |
|:-----------------------------|:--------------------------------------------------------|
| -n, \--apiName `string`      | The name under which the generated docs will be grouped |
| -f, \--file `string`         | OpenAPI 3 spec file                                     |
| -o, \--outputDir `string`    | The output directory                                    |
| -r, \--referenceURL `string` | The reference URL (default is `reference`)                 |
| -h, \--help                  | Help for the conversion                                        |

| Option                       | Description                                             |
|:-----------------------------|:--------------------------------------------------------|
| -e, --allowExternalRefs      | Allow external references in the OpenAPI spec. |

| -n, --apiName string The name under which the generated docs will be grouped 
| -f, --file string OpenAPI 3 spec file 
| -h, --help help for convert --includeExamples Include a column on the schema for examples 
| --includeRestrictions Include a column on the schema for restrictions (default true) 
| -i, --inlineProperties Inline properties in the request and response schemas 
| -o, --outputDir string The output directory 
| -r, --referenceURL string The reference URL (default "reference") 
| -s, --sortFilePath Sort by filepath by prefixing a weight to the filename. Default is to use front matter weight 
| -t, --titleFormat string The template format used to create the title for each operation. Valid options are: - operationId: (Default) Uses the value of the operationId field. - MethodURL: Uses a combination of the Method property and the URL