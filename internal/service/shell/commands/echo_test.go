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
	ctx := context.Background()

	cmd := commands.NewEchoCommand()

	t.Run("success - prints plain arguments", func(t *testing.T) {
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{"Hello", "World"}, nil, &outputBuffer, &errorBuffer)
		assert.NoError(t, err)
		assert.Equal(t, "Hello World\n", outputBuffer.String())
		assert.Empty(t, errorBuffer.String())
	})

	t.Run("success - handles single-quoted literals", func(t *testing.T) {
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{"'Hello'", "'World'"}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, errorBuffer.String())
		assert.Equal(t, "Hello World\n", outputBuffer.String())
	})

	t.Run("success - expands environment variables", func(t *testing.T) {
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		os.Setenv("TEST_VAR", "GoShell")
		defer os.Unsetenv("TEST_VAR")

		err := cmd.Execute(ctx, []string{"Hello", "$TEST_VAR"}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, errorBuffer.String())
		assert.Equal(t, "Hello GoShell\n", outputBuffer.String())
	})

	t.Run("success - handles multiple environment variables", func(t *testing.T) {
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		os.Setenv("VAR1", "Foo")
		os.Setenv("VAR2", "Bar")
		defer os.Unsetenv("VAR1")
		defer os.Unsetenv("VAR2")

		err := cmd.Execute(ctx, []string{"$VAR1", "and", "$VAR2"}, nil, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, errorBuffer.String())
		assert.Equal(t, "Foo and Bar\n", outputBuffer.String())
	})

	t.Run("success - ignores undefined environment variables", func(t *testing.T) {
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer

		err := cmd.Execute(ctx, []string{"Hello", "$UNDEFINED_VAR"}, nil, &outputBuffer, &errorBuffer)
		assert.NoError(t, err)
		assert.Empty(t, errorBuffer.String())
		assert.Equal(t, "Hello\n", outputBuffer.String())
	})

	t.Run("success - returns empty string when no arguments", func(t *testing.T) {
		var outputBuffer bytes.Buffer
		var errorBuffer bytes.Buffer
		input := strings.NewReader("")

		err := cmd.Execute(ctx, []string{}, input, &outputBuffer, &errorBuffer)

		assert.NoError(t, err)
		assert.Empty(t, errorBuffer.String())
		assert.Equal(t, "\n", outputBuffer.String())
	})

}
