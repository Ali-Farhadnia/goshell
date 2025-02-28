package commands_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Ali-Farhadnia/goshell/internal/service/shell/commands"
	"github.com/stretchr/testify/assert"
)

func TestEchoCommand_Execute(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  string
		setupEnv       map[string]string
		teardownEnv    []string
	}{
		{
			name:           "prints plain arguments",
			args:           []string{"Hello", "World"},
			expectedOutput: "Hello World\n",
			expectedError:  "",
		},
		{
			name:           "handles single-quoted literals",
			args:           []string{"'Hello'", "'World'"},
			expectedOutput: "Hello World\n",
			expectedError:  "",
		},
		{
			name:           "expands environment variables",
			args:           []string{"Hello", "$TEST_VAR"},
			expectedOutput: "Hello GoShell\n",
			expectedError:  "",
			setupEnv:       map[string]string{"TEST_VAR": "GoShell"},
			teardownEnv:    []string{"TEST_VAR"},
		},
		{
			name:           "handles multiple environment variables",
			args:           []string{"$VAR1", "and", "$VAR2"},
			expectedOutput: "Foo and Bar\n",
			expectedError:  "",
			setupEnv:       map[string]string{"VAR1": "Foo", "VAR2": "Bar"},
			teardownEnv:    []string{"VAR1", "VAR2"},
		},
		{
			name:           "ignores undefined environment variables",
			args:           []string{"Hello", "$UNDEFINED_VAR"},
			expectedOutput: "Hello\n",
			expectedError:  "",
		},
		{
			name:           "returns empty string when no arguments",
			args:           []string{},
			expectedOutput: "\n",
			expectedError:  "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuffer bytes.Buffer
			var errorBuffer bytes.Buffer
			ctx := context.Background()
			cmd := commands.NewEchoCommand()

			// Setup environment variables
			for key, value := range tc.setupEnv {
				os.Setenv(key, value)
			}

			err := cmd.Execute(ctx, tc.args, strings.NewReader(""), &outputBuffer, &errorBuffer)

			// Teardown environment variables
			for _, key := range tc.teardownEnv {
				os.Unsetenv(key)
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, outputBuffer.String())
			assert.Equal(t, tc.expectedError, errorBuffer.String())
		})
	}
}
