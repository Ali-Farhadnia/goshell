package inputprocessor

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ParseArguments parses the input string into a slice of arguments, handling quotes and escaping.
func ParseArguments(input string) ([]string, error) {
	var args []string
	var current strings.Builder
	inQuotes := false
	escaped := false
	quoteChar := byte(0)

	for i := 0; i < len(input); i++ {
		c := input[i]

		if escaped {
			// Allow escaping only specific characters inside quotes
			if inQuotes && (c == '$' || c == '`' || c == '"' || c == '\\') {
				current.WriteByte(c)
			} else if !inQuotes {
				current.WriteByte(c) // keep escaped character
			} else {
				current.WriteByte('\\') // Keep backslash
				current.WriteByte(c)
			}
			escaped = false
			continue
		}

		if c == '\\' {
			escaped = true
			continue
		}

		if c == '"' {
			if inQuotes && quoteChar == '"' {
				inQuotes = false // End quote
			} else if !inQuotes {
				inQuotes = true
				quoteChar = '"'
			} else {
				current.WriteByte(c) // Inside different quote type, treat as normal character
			}
			continue
		}

		if (c == ' ' || c == '\t') && !inQuotes {
			// End
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
			continue
		}

		current.WriteByte(c)
	}

	// Append last argument if exists
	if current.Len() > 0 {
		args = append(args, current.String())
	}

	if inQuotes {
		return nil, fmt.Errorf("unterminated quote detected")
	}

	return args, nil
}

// ProcessRedirections processes the input arguments for redirections and returns the appropriate readers/writers.
func ProcessRedirections(args []string) (io.Reader, io.Writer, io.Writer, []string, func()) {
	inputReader := os.Stdin
	outputWriter := os.Stdout
	errorOutputWriter := os.Stderr
	var inputFile, outputFile, errorFile *os.File
	cleanArgs := []string{}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case ">":
			if i+1 < len(args) {
				f, err := os.Create(args[i+1])
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				outputWriter = f
				outputFile = f
				i++
			}
		case ">>":
			if i+1 < len(args) {
				f, err := os.OpenFile(args[i+1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				outputWriter = f
				outputFile = f
				i++
			}
		case "2>":
			if i+1 < len(args) {
				f, err := os.Create(args[i+1])
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				errorOutputWriter = f
				errorFile = f
				i++
			}
		case "2>>":
			if i+1 < len(args) {
				f, err := os.OpenFile(args[i+1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				errorOutputWriter = f
				errorFile = f
				i++
			}
		case "<":
			if i+1 < len(args) {
				f, err := os.Open(args[i+1])
				if err != nil {
					fmt.Println("error:", err)
					continue
				}
				inputReader = f
				inputFile = f
				i++
			}
		default:
			cleanArgs = append(cleanArgs, args[i])
		}
	}

	// Cleanup function to close files
	cleanup := func() {
		if inputFile != nil {
			inputFile.Close()
		}
		if outputFile != nil {
			outputFile.Close()
		}
		if errorFile != nil {
			errorFile.Close()
		}
	}

	return inputReader, outputWriter, errorOutputWriter, cleanArgs, cleanup
}
