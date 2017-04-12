---
title: "Modeling Anonymous Containers"
author: "VirtualTraveler"
---
In many systems commonly used services are frequently given only a generic name often based on the implementation technology. For example many micro-services in a system may each have their own dedicated cache container. These caches may all be called "Redis".    

# Solution

1. Create a unique container with it's own name for each logical cacahe 
2. Assocate the each cache with the container(s) it serves

# Discussion

If you want to show the distinct cache containers on the diagram it is necessary to name each one you want to show uniquely.   

