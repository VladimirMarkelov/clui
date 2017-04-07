# Control Layout
CLUI control layout management is similar to management used in well-known libraries like Qt or FOX toolkit. The difference between them and CLUI that CLUI does not have dedicated controls for layout management - any control turns into layout(container) if it has children, and CLUI layout managent is much simpler - only one way to arrange controls is available.

### Layout basics
A control that has children can arrange them only in one direction: from left to right or from top to bottom in order of adding the children. So, any container is always either 1 control high or 1 control wide depending on layout type. The layout type for a control can be set by calling:
```
control.SetPack(newPackDirection)
```
where newPackDirection either Horizontal(default value) or Vertical.

Automatic control placement can be tuned up with extra container properties for nicer-looking result. It is padding and gap values. By default all of them are 0 (except paddings for Window and Frame with a border - in this case default paddings are 1). Padding is a number of character between the container edge and the first child, gap is a number of characters between children inside container. Please see a picture to get the idea:

<img src="/docs/img/layout.png" alt="Layout manager">

Please, keep in mind the following feature of CLUI layout manager:
* padding is always calculated from the very edge of a control. That is why Window and Frame with a border have default paddings equal 1 - to avoid overlapping their children with their border
* Fixed control placement and grid layout manager are not available - only horizontal and vertical automatic layout managers. It may make designing a complex Window layout difficult
* It is possible to set minimal width and height for a container control but if a container has at least one child then the real minimal size is calculated as a maximum of container's minimal values and the total space required to display all its children. In other words, you cannot set minimal size of a container less than the total minimal size of its children plus gaps and paddings
* To create a control aligned to bottom or right side, use the following trick: at first add a frameless Frame with scale equals 1 and after it add the control with scale equals Fixed. It makes the frame resizable when its parent is resized while the control will keep its size and will always stick to the container edge
