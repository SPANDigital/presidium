title: "Combining Selectors"
category: "DSS"
author: "VirtualTraveler"
---
Selectors can be combined in two ways

```Layer=Services + Category=FunctionalContainer```

Selects all objects that are in the service layer with their relationships AND all objects that are functional containers and their relationships. This is a boolean AND operation

```Layer=Services , Category=FunctionalContainer```

Selects all objects that are in the services layer and their relationships OR all objects that are functional containers and their relationships. This is a boolean OR operation.
