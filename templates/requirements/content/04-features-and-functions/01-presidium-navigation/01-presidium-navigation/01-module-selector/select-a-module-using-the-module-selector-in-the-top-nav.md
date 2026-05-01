---
title: Select a module using the module selector in the top navigation bar
---

scenario ID: 45

## Given 
- Given I am any type of user, on any page within Presidium

## When 
- When I click on the ""Go to Module"" dropdown in the top bar
And I select a module from the list of modules

## Then 
- Then I should be redirected to the selected module
- And the page should display content related to that module
- And the module dropdown should have the module I selected showing "