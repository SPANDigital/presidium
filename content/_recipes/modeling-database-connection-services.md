---
title: "Modeling Database Connection Services"
author: "VirtualTraveler"
---
Database connection services are common patterns in large systems. These services broker connections from multiple clients to multiple databases. Usually these types of service handle connection pooling, caching, load balancing, and failover.   

# Solution

1. Create a messaging container for each distinct logical database connection service
2. If the connection service exists witin a specific functional area then associate it with that area
3. Create calls relationships from client services to the connection service and from the connection service to the databases that it handles. 

# Discussion

A messaging container will spread out to cover the width of it's enclosong functional area. This is a convenient way of modeling containers that connect many clients services to many provider services.  

