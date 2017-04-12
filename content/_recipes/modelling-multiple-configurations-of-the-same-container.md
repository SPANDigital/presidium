---
title: "Modeling Multiple Configurations of the same Container"
author: "VirtualTraveler"
---
Sometimes the same container is deployed with multiple different configurations to handle different use cases. In these situations we often want to give each configuration a different name.  

# Solution

1. Create a container for each different configuration. Give each container a unique name.
2. Create containment relationships from the new containers to the same set of components. This results in multiple containers with different names that contain the same components. 

# Discussion

A container is best thought of as a executable version of a configured component. 

