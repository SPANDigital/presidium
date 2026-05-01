---
title: "A Scenario"
type: scenario
author: "author"
scenarioStatus: "proposed"
scenario_id: ""
weight: 1
testSuites:
  - "smoke"
  - "regression"
depends_on:
  - scenario_id: ""
    confidence_level: 0.85
---

### Logical

```gherkin
Scenario: {{ `{{< article-title >}}` }}
  Given [precondition]
  When [action]
  Then [expected outcome]
```

### Test

```gherkin
Scenario: {{ `{{< article-title >}}` }}
  Given [precondition with UI context]
  When [UI action]
  Then [observable outcome]
```
