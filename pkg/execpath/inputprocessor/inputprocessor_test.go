package inputprocessor

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArguments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
		hasError bool
	}{
		{
			name:     "simple command",
			input:    "ls -la",
			expected: []string{"ls", "-la"},
			hasError: false,
		},
		{
			name:     "with quotes",
			input:    `echo "hello world"`,
			expected: []string{"echo", "hello world"},
			hasError: false,
		},
		{
			name:     "with escaped characters",
			input:    `echo hello\ world`,
			expected: []string{"echo", "hello world"},
			hasError: false,
		},
		{
			name:     "with escaped quotes",
			input:    `echo "hello \"world\""`,
			expected: []string{"echo", `hello "world"`},
			hasError: false,
		},
		{
			name:     "with escaped special characters",
			input:    `echo "hello \$USER"`,
			expected: []string{"echo", "hello $USER"},
			hasError: false,
		},
		{
			name:     "with escaped backslash",
			input:    `echo "hello \\ world"`,
			expected: []string{"echo", "hello \\ world"},
			hasError: false,
		},
		{
			name:     "with escaped backtick",
			input:    `echo "hello \` + "`" + `world"`,
			expected: []string{"echo", "hello `world"},
			hasError: false,
		},
		{
			name:     "with tabs",
			input:    "ls\t-la",
			expected: []string{"ls", "-la"},
			hasError: false,
		},
		{
			name:     "with multiple spaces",
			input:    "ls   -la",
			expected: []string{"ls", "-la"},
			hasError: false,
		},
		{
			name:     "unterminated quote",
			input:    `echo "hello world`,
			expected: nil,
			hasError: true,
		},
		{
			name:     "non-special escapes outside quotes",
			input:    `echo hello\world`,
			expected: []string{"echo", "helloworld"},
			hasError: false,
		},
		{
			name:     "non-special escapes inside quotes",
			input:    `echo "hello\world"`,
			expected: []string{"echo", "hello\\world"},
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseArguments(tt.input)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessRedirections(t *testing.T) {
	// Create temporary files for testing
	tmpIn, err := os.CreateTemp("", "input")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpIn.Name())

	tmpOut, err := os.CreateTemp("", "output")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpOut.Name())

	tmpErr, err := os.CreateTemp("", "error")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpErr.Name())

	// Write test data to input file
	inputData := "test input data"
	tmpIn.WriteString(inputData)
	tmpIn.Close()

	tests := []struct {
		name           string
		args           []string
		expectedArgs   []string
		checkInput     bool
		checkOutput    bool
		checkErrOutput bool
	}{
		{
			name:           "no redirections",
			args:           []string{"echo", "hello"},
			expectedArgs:   []string{"echo", "hello"},
			checkInput:     false,
			checkOutput:    false,
			checkErrOutput: false,
		},
		{
			name:           "output redirection",
			args:           []string{"echo", "hello", ">", tmpOut.Name()},
			expectedArgs:   []string{"echo", "hello"},
			checkInput:     false,
			checkOutput:    true,
			checkErrOutput: false,
		},
		{
			name:           "input redirection",
			args:           []string{"cat", "<", tmpIn.Name()},
			expectedArgs:   []string{"cat"},
			checkInput:     true,
			checkOutput:    false,
			checkErrOutput: false,
		},
		{
			name:           "error redirection",
			args:           []string{"ls", "nonexistent", "2>", tmpErr.Name()},
			expectedArgs:   []string{"ls", "nonexistent"},
			checkInput:     false,
			checkOutput:    false,
			checkErrOutput: true,
		},
		{
			name:           "append output",
			args:           []string{"echo", "hello", ">>", tmpOut.Name()},
			expectedArgs:   []string{"echo", "hello"},
			checkInput:     false,
			checkOutput:    true,
			checkErrOutput: false,
		},
		{
			name:           "append error",
			args:           []string{"ls", "nonexistent", "2>>", tmpErr.Name()},
			expectedArgs:   []string{"ls", "nonexistent"},
			checkInput:     false,
			checkOutput:    false,
			checkErrOutput: true,
		},
		{
			name:           "multiple redirections",
			args:           []string{"cat", "<", tmpIn.Name(), ">", tmpOut.Name(), "2>", tmpErr.Name()},
			expectedArgs:   []string{"cat"},
			checkInput:     true,
			checkOutput:    true,
			checkErrOutput: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentDir, _ := os.Getwd()
			reader, writer, errWriter, cleanArgs, cleanup := ProcessRedirections(tt.args, currentDir)
			defer cleanup()

			// Check that arguments were cleaned correctly
			assert.Equal(t, tt.expectedArgs, cleanArgs)

			// Check that stdin was changed if we expect it
			if tt.checkInput {
				assert.NotEqual(t, os.Stdin, reader)

				// Verify we can read from the file
				buf := make([]byte, len(inputData))
				n, err := reader.Read(buf)
				if err != nil && err != io.EOF {
					t.Errorf("unexpected error reading from input: %v", err)
				}
				assert.Equal(t, len(inputData), n)
				assert.Equal(t, inputData, string(buf[:n]))
			} else {
				assert.Equal(t, os.Stdin, reader)
			}

			// Check that stdout was changed if we expect it
			if tt.checkOutput {
				assert.NotEqual(t, os.Stdout, writer)
			} else {
				assert.Equal(t, os.Stdout, writer)
			}

			// Check that stderr was changed if we expect it
			if tt.checkErrOutput {
				assert.NotEqual(t, os.Stderr, errWriter)
			} else {
				assert.Equal(t, os.Stderr, errWriter)
			}

			// Test writing to output if redirected
			if tt.checkOutput {
				testOut := []byte("test output data")
				_, err := writer.Write(testOut)
				assert.NoError(t, err)
			}

			// Test writing to error output if redirected
			if tt.checkErrOutput {
				testErr := []byte("test error data")
				_, err := errWriter.Write(testErr)
				assert.NoError(t, err)
			}
		})
	}
}

// TestProcessRedirectionsInvalidFiles tests error handling for invalid files
func TestProcessRedirectionsInvalidFiles(t *testing.T) {
	// Redirect stdout and stderr temporarily to capture error messages
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with invalid file paths
	currentDir, _ := os.Getwd()
	args := []string{"cat", "<", "/invalid/path/that/should/not/exist"}
	_, _, _, _, cleanup := ProcessRedirections(args, currentDir)
	cleanup()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Verify we got an error message
	assert.Contains(t, buf.String(), "error:")
}
