package clui

import (
	term "github.com/nsf/termbox-go"
	мИнт "./пакИнтерфейсы"
)

// Composer is a service object that manages Views and console, processes
// events, and provides service methods. One application must have only
// one object of this type
type mainLoop struct {
	// a channel to communicate with View(e.g, Views send redraw event to this channel)
	channel chan мИнт.ИСобытие
}

var (
	loop *mainLoop
)

func initMainLoop() {
	loop = new(mainLoop)
	loop.channel = make(chan мИнт.ИСобытие)
}

// MainLoop starts the main application event loop
func MainLoop() {
	RefreshScreen()

	eventQueue := make(chan term.Event)
	go func() {
		for {
			eventQueue <- term.PollEvent()
		}
	}()

	for {
		RefreshScreen()

		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case term.EventError:
				panic(ev.Err)
			default:
				ProcessEvent(termboxEventToLocal(&ev))
			}
		case cmd := <-loop.channel:
			if cmd.Type() == мИнт.EventQuit {
				return
			}
			ProcessEvent(cmd)
		}
	}
}

func _putEvent(ev мИнт.ИСобытие) {
	loop.channel <- ev
}

// PutEvent send event to a Composer directly.
// Used by Views to ask for repainting or for quitting the application
func PutEvent(ev мИнт.ИСобытие) {
	go _putEvent(ev)
}
