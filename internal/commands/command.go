package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/RaphSku/notewolfy/internal/structure"
)

type Strategy interface {
	Run() error
}

type Context struct {
	strategy Strategy
}

func NewContext(strategy Strategy) *Context {
	return &Context{strategy}
}

func (c *Context) RunStrategy() error {
	return c.strategy.Run()
}

func validateAndTrimStatement(statement string) string {
	if len(statement) == 0 {
		return ""
	}
	re := regexp.MustCompile(`\s+`)
	trimmedStatement := re.ReplaceAllString(statement, " ")
	trimmedStatement = strings.TrimSpace(trimmedStatement)

	return trimmedStatement
}

func MatchStatementToCommand(mmf *structure.MetadataNoteWolfyFileHandle, statement string) {
	validatedStatement := validateAndTrimStatement(statement)

	strategy := matchStatementToStrategy(mmf, validatedStatement)
	if strategy == nil {
		return
	}
	context := NewContext(strategy)
	if err := context.RunStrategy(); err != nil {
		fmt.Println(err)
	}
}
