# clui
Command Line User Interface (Console UI inspired by TurboVision).

WARNING: the library is experimental. Use it at your own risk.

## Introduction
It includes a few number of controls that is enough to create an application for every day task. More controls are to come later.
The current list of controls:
* Label
* Button
* EditFiled
* ComboBox
* ListBox
* TextScroll: simple control to display scrolling text (e.g, for tail output)
* ProgressBar
* Frame
* CheckBox
* Radio: RadioGroup

Built-in theme support. Now it is very basic: no way to load any theme from file, only one predefined theme

A set of global hotkeys. Windows version has more features and hotkeys - it is current limitation of 'termbox' library: https://github.com/VladimirMarkelov/termbox-go

## Screenshots
The screencast of demo included in the library:

<img src="./demos/demo.gif" alt="Library Demo">

The library is in the very beginning but it can be used to create working utilities: below is the example of my Dilbert comix downloader:

<img src="./demos/dilbert_demo.gif" alt="Dilbert Downloader">
