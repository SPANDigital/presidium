---
title: "DSS Introduction"
category: "DSS"
id: "DSS
author: "VirtualTraveler"
---
The visual display of diagrams produced by ADSS can be controlled and modified by diagram style sheets. Similar to CSS, A DSS rule-set consists of one or more selectors and a declaration block.

{% highlight CSS %}
Layer=Clients {	Shape : cylinder;
	fillColor : rgb(255,159,98);
	borderColor : rgb(95,95,95);
	borderWeight : 2;
	fontFamily : Arial;
	fontSize : 30;
	fontColor : rgb(255,255,255);
}
{% endhighlight %}

The selector identifies the diagram element(s) you want to style.
The declaration block contains one or more declarations separated by semicolons.
Each declaration includes a DSS property name and a value, separated by a colon.
A DSS declaration always ends with a semicolon, and declaration blocks are surrounded by curly braces.
