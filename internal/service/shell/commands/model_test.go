package commands_test

import (
	"context"
	"io"
)

type MockCommand struct {
	NameVal string
	HelpVal string
}

func (m *MockCommand) Name() string {
	return m.NameVal
}

func (m *MockCommand) Help() string {
	return m.HelpVal
}

func (m *MockCommand) MaxArguments() int {
	return 0
}

func (m *MockCommand) Execute(ctx context.Context, args []string, inputReader io.Reader, outputWriter, errorOutputWriter io.Writer) error {
	return nil
}
