---
title: Focus on Containers first
id: "focus-on-containers-first"
author: "VirtualTraveler"
---

Containers are the core of the ADSS model. They are the bridge between the logical and physical views. Once you have completed the setup of your model, defined in the [Getting Started](/getting-started/) section, you should complete an inventory of Containers.

## Discussion 

Containers are physically deployable collections of componenets. They come in three variaties 

- **Functional** A logical service or application whose primary role is to implement application or business logic. This may be a client application or a backend services.    
- **Messaging** A logical messaging channel, topic, or queue. A messaging container carries messages, events or work items from one part of the system to another. 
- **Data** A logical collection of data implemented as a single container. Often the name of a database or schema. For example a an Postgres database containing a schema called Commerce would could be called Commerce_DB. It should not be called postgres, thats is the name of the implementation technology  

When creating containers add as many other details as you can. 

- **Layer** The layer the Container should be displayed in. **Note:** The container will not be displayed on the logical diagram if it doesn't have a layer.
- **Owner** The person who takes responsibility for the Container
- **Organization** The organization that takes responsibility for the Container
- **Functional Area** The area the container belongs to 
- **Lifecycle Stage** The Status of the Container

## See Also

* [REPLACE WITH a link description](http://www.google.com) 
