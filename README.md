# clui
Command Line User Interface (Console UI inspired by TurboVision) with built-in theme support.

## Introduction
The list of available controls:
* Label (Horizontal and Vertical with basic color control tags)
* Button (Simple push button control)
* EditFiled (One line text edit control with basic clipboard control)
* ListBox (string list control with vertical scroll)
* TextView (ListBox-alike control with vertical and horizontal scroll, and wordwrap mode)
* ProgressBar (Vertical and horizontal. The latter one supports custom text over control)
* Frame (A decorative control that can be a container for other controls as well)
* CheckBox (Simple check box)
* Radio (Simple radio button. Useless alone - should be used along with RadioGroup)
* RadioGroup (Non-visual control to manage a group of a few RadioButtons)

### TODO
* BarChart (Horizontal bar chart without scroll)
* Diagram (Show tabular data as a line graph or sparkle one)
* GridView (Table to show structured data - only virtual and readonly mode with scroll support)

## Screenshots
The screencast of demo (based on custom termbox-go Window build) included in the library:

<img src="./demos/demo.gif" alt="Library Demo">

The library is in the very beginning but it can be used to create working utilities: below is the example of my Dilbert comix downloader:

<img src="./demos/dilbert_demo.gif" alt="Dilbert Downloader">
