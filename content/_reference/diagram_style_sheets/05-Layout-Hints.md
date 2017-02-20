---
title: "Layout Hints"
category: "DSS"
author: "VirtualTraveler"
---
The layout of the diagram can be influenced by providing hints to the layout optimizer. These hints may be ignored if the optimizer decides they are too difficult to follow.

| HintName              | Description                                    |
|-----------------------|------------------------------------------------|
| distribute-horizontally| Instructs the layout engine to favor expanding the width of the element to accommodate its contents. May expand to the full width of the containing element.
| distribute-vertically | Instructs the layout engine to favor expanding the height of the element to accommodate its contents. May expand to the full height of the containing element. 
| corral-elements| Instructs the layout engine to gather together contents as compactly as reasonable.
| concentrate-lines| Instructs the layout engine to merge lines to reduce the complexity. This may reduce the fidelity of the diagram 
