---
title: "title"
type: scenario
author: "email@address.com"
scenarioStatus: "proposed"
scenarioId: ""
testSuites:
  - ""
depends_on:
  - scenario_id: ""
    confidence_level: 0.00
---

### Logical

```gherkin
Scenario | Scenario Outline: {{ `{{` }}< article-title >{{ `}}` }}
  Given [condition]
  And [additional condition]
  When [action]
  And [additional action]
  Then [expected result]
  And [additional expected result]

  Examples:
  | parameter | value |
  | example   | data  |

```

### Test

```gherkin
Scenario | Scenario Outline: {{ `{{` }}< article-title >{{ `}}` }}
  Given [precondition with UI context]
  When [UI action]
  And [field interaction]
  Then [observable outcome]

  Examples:
  | parameter | value |
  | example   | data  |

```
