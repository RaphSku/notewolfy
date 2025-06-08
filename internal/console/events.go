package console

import (
	"fmt"
	"path/filepath"

	"github.com/RaphSku/cyclecmd"
	"github.com/RaphSku/notewolfy/internal/commands"
	"github.com/RaphSku/notewolfy/internal/structure"
	"github.com/RaphSku/notewolfy/internal/utility"
)

const NOTEWOLFY_GOODBYE_MESSAGE = "\r\nThank you for using notewolfy!\r\n"

type DefaultEvent struct{}

func (de *DefaultEvent) Handle(token string) (error, *cyclecmd.ControlEvent) {
	fmt.Print(token)
	return nil, nil
}

type BackspaceEvent struct{}

func (be *BackspaceEvent) Handle(token string) (error, *cyclecmd.ControlEvent) {
	eventHistory := getEventHistory()
	var filteredEventNames []string
	eventNamesFromHistory := eventHistory.GetLastEventsFromHistoryToEventReference("Enter")
	for _, eventName := range eventNamesFromHistory {
		if eventName == "NonDefault" {
			continue
		}
		if eventName != "Backspace" {
			filteredEventNames = append(filteredEventNames, eventName)
		}
	}

	if len(filteredEventNames) == 0 {
		return nil, nil
	}

	eventHistoryLength := eventHistory.Len()
	// Backspace will only work if we remove the backspace event itself and the event previous to the backspace
	eventHistory.RemoveNthEventFromHistory(eventHistoryLength - 1)
	eventHistory.RemoveNthEventFromHistory(eventHistoryLength - 2)

	fmt.Print("\b \b")
	return nil, nil
}

var mmf *structure.MetadataNoteWolfyFileHandle

func InitMetadataNoteWolfyFileHandle() error {
	if mmf == nil {
		homeDir, err := utility.GetHomeDir()
		if err != nil {
			return err
		}
		metadataFilePath := filepath.Join(homeDir, ".notewolfy")
		config := &structure.Config{
			MetadataFilePath: metadataFilePath,
		}
		newmmf, err := structure.NewMetadataNoteWolfyFileHandle(config)
		if err != nil {
			return err
		}
		mmf = newmmf
	}
	return nil
}

type EnterEvent struct{}

func (ee *EnterEvent) Handle(token string) (error, *cyclecmd.ControlEvent) {
	err := InitMetadataNoteWolfyFileHandle()
	if err != nil {
		return err, nil
	}
	statement := buildStatement()
	fmt.Print("\r")
	if checkExitCommand(statement) {
		fmt.Print(NOTEWOLFY_GOODBYE_MESSAGE)
		controlEvent := cyclecmd.NewControlEvent(cyclecmd.CYCLE_TERMINATE)
		return nil, controlEvent
	}
	handleEnter(mmf, statement)
	return nil, nil
}

type EscapeEvent struct{}

func (ee *EscapeEvent) Handle(token string) (error, *cyclecmd.ControlEvent) {
	if checkEscExitCondition(token) {
		fmt.Print(NOTEWOLFY_GOODBYE_MESSAGE)
		controlEvent := cyclecmd.NewControlEvent(cyclecmd.CYCLE_TERMINATE)
		return nil, controlEvent
	}
	return nil, nil
}

type CtrlCEvent struct{}

func (ce *CtrlCEvent) Handle(token string) (error, *cyclecmd.ControlEvent) {
	if checkCtrlCExitCondition(token) {
		fmt.Print(NOTEWOLFY_GOODBYE_MESSAGE)
		controlEvent := cyclecmd.NewControlEvent(cyclecmd.CYCLE_TERMINATE)
		return nil, controlEvent
	}
	return nil, nil
}

func checkEscExitCondition(token string) bool {
	return token == "\x03"
}

func checkCtrlCExitCondition(token string) bool {
	return token == "\x1b"
}

func buildStatement() string {
	eventHistory := getEventHistory()
	var statement string
	splicedEventEntries := eventHistory.MostRecentSpliceEventsOfHistory("Enter")
	for _, eventEntry := range splicedEventEntries {
		if _, ok := eventEntry.Event.(*DefaultEvent); ok {
			statement += eventEntry.Token
			continue
		}
		if _, ok := eventEntry.Event.(*BackspaceEvent); ok {
			if len(statement) == 0 {
				continue
			}
			statement = statement[:len(statement)-1]
			continue
		}
	}
	return statement
}

func handleEnter(mmf *structure.MetadataNoteWolfyFileHandle, statement string) {
	commands.MatchStatementToCommand(mmf, statement)
}

func checkExitCommand(statement string) bool {
	shouldBreak := false
	switch statement {
	case "exit":
		shouldBreak = true
	case "quit":
		shouldBreak = true
	}
	return shouldBreak
}
