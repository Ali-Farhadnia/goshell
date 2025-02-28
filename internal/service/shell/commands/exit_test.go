package commands_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/stretchr/testify/assert"
)

func TestExitCommand_Execute(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedExit   int
		expectedOutput string
		expectedError  string
	}{
		{"Default exit code", []string{}, 0, "exit status 0\n", ""},
		{"Valid exit code", []string{"5"}, 5, "exit status 5\n", ""},
		{"Invalid exit code", []string{"abc"}, 0, "", "Invalid exit code: abc\n"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var capturedExitCode int
			exitCmd := commands.NewExitCommand(nil, func(code int) {
				capturedExitCode = code // Capture the exit code
			})

			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer

			err := exitCmd.Execute(context.Background(), tc.args, nil, &outputBuffer, &errorBuffer)

			assert.NoError(t, err)

			assert.Equal(t, tc.expectedExit, capturedExitCode)
			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())
		})
	}
}
