package console

import (
	"context"

	"github.com/RaphSku/cyclecmd"
	"github.com/RaphSku/notewolfy/cmd/version"
)

var eventHistory *cyclecmd.EventHistory

func getEventHistory() *cyclecmd.EventHistory {
	if eventHistory == nil {
		eventHistory = cyclecmd.NewEventHistory()
	}
	return eventHistory
}

func StartConsoleApplication(ctx context.Context) {
	defaultEventInformation := cyclecmd.EventInformation{
		EventName: "Default",
		Event:     &DefaultEvent{},
	}
	eventRegistry := cyclecmd.NewEventRegistry(defaultEventInformation)

	backspaceEventInformation := cyclecmd.EventInformation{
		EventName: "Backspace",
		Event:     &BackspaceEvent{},
	}
	eventRegistry.RegisterEvent("\x7f", backspaceEventInformation)

	enterEventInformation := cyclecmd.EventInformation{
		EventName: "Enter",
		Event:     &EnterEvent{},
	}
	eventRegistry.RegisterEvent("\r", enterEventInformation)

	escEventInformation := cyclecmd.EventInformation{
		EventName: "Escape",
		Event:     &EscapeEvent{},
	}
	eventRegistry.RegisterEvent("\x03", escEventInformation)

	ctrlcEventInformation := cyclecmd.EventInformation{
		EventName: "Ctrl+C",
		Event:     &CtrlCEvent{},
	}
	eventRegistry.RegisterEvent("\x1b", ctrlcEventInformation)

	eventHistory := getEventHistory()

	consoleApp := cyclecmd.NewConsoleApp(
		context.Background(),
		"notewolfy",
		version.VERSION,
		"Creating organized notes is just easy with notewolfy",
		eventRegistry,
		eventHistory,
	)
	consoleApp.SetLineDelimiter("\n\r>>> ", "\r")
	consoleApp.Start()
}
