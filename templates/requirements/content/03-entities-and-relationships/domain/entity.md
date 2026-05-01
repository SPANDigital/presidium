---
title: An Entity
author: author
name: "ExampleEntity"
description: "A short description of what this entity represents in the domain."
attributes:
  - name: "id"
    key_role: "primary"
    optionality: "required"
    description: "Unique identifier for this entity."
  - name: "name"
    key_role: ""
    optionality: "required"
    description: "Display name shown in the UI and reports."
  - name: "status"
    key_role: ""
    optionality: "optional"
    description: "Current lifecycle status of this entity."
    enum: "active | archived | deleted"
---

{{ `{{< entity-table >}}` }}
