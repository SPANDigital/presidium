---
title: Dashboard Statistics
author: author
showHistorySection: false
---

## Current Statistics

{{ `{{< discovery-stats >}}` }}

{{ `{{% if-page-param "showHistorySection" %}}` }}
## Historical Data

{{ `{{< discovery-stats show-history="true" show-status="false" >}}` }}
{{ `{{% /if-page-param %}}` }}
