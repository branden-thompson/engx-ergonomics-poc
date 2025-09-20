package components

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bthompso/engx-ergonomics-poc/internal/tui/design"
)

// CommandFormatter provides smart syntax highlighting for engx commands
// Follows the pattern: `engx [cmd] [subcmd] <cmd-or-subcmd-arg> [-flags] --option <option-arg>`
type CommandFormatter struct{}

// NewCommandFormatter creates a new command formatter
func NewCommandFormatter() *CommandFormatter {
	return &CommandFormatter{}
}

// FormatCommand applies smart syntax highlighting to an engx command
func (cf *CommandFormatter) FormatCommand(command string) string {
	// Remove backticks if present
	command = strings.Trim(command, "`")


	result := command

	// Apply colors in specific order to avoid conflicts
	// Start with the most specific patterns first

	// 1. Handle 'engx' first
	engxRegex := regexp.MustCompile(`\bengx\b`)
	result = engxRegex.ReplaceAllStringFunc(result, func(match string) string {
		return fmt.Sprintf("\033[38;5;%sm%s\033[0m", design.ColorEngxPink, match)
	})

	// 2. Handle angle bracket arguments
	argRegex := regexp.MustCompile(`<[^>]+>`)
	result = argRegex.ReplaceAllStringFunc(result, func(match string) string {
		return fmt.Sprintf("\033[%sm%s\033[0m", design.ColorBrightBlue, match)
	})

	// 3. Handle curly brace placeholders
	placeholderRegex := regexp.MustCompile(`\{[^}]+\}`)
	result = placeholderRegex.ReplaceAllStringFunc(result, func(match string) string {
		return fmt.Sprintf("\033[%sm%s\033[0m", design.ColorBrightBlue, match)
	})

	// 4. Handle quoted option values
	quoteRegex := regexp.MustCompile(`'[^']*'`)
	result = quoteRegex.ReplaceAllStringFunc(result, func(match string) string {
		return fmt.Sprintf("\033[%sm%s\033[0m", design.ColorBrightBlue, match)
	})

	// 5. Handle long options --option
	longOptionRegex := regexp.MustCompile(`--[\w-]+`)
	result = longOptionRegex.ReplaceAllStringFunc(result, func(match string) string {
		return fmt.Sprintf("\033[38;5;%sm%s\033[0m", design.ColorFlagGreen, match)
	})

	// 6. Handle short flags -f
	shortFlagRegex := regexp.MustCompile(`-[a-zA-Z]\b`)
	result = shortFlagRegex.ReplaceAllStringFunc(result, func(match string) string {
		return fmt.Sprintf("\033[38;5;%sm%s\033[0m", design.ColorFlagGreen, match)
	})

	// 7. Handle commands and subcommands (more complex logic)
	result = cf.highlightCommandsAndSubcommands(result)

	// 8. Handle option arguments that come after options
	optionArgRegex := regexp.MustCompile(`(--[\w-]+\s+)([^\s]+)`)
	result = optionArgRegex.ReplaceAllStringFunc(result, func(match string) string {
		parts := optionArgRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			// parts[1] is the option (already colored), parts[2] is the argument
			return parts[1] + fmt.Sprintf("\033[%sm%s\033[0m", design.ColorBrightBlue, parts[2])
		}
		return match
	})

	return result
}

// highlightCommandsAndSubcommands applies color to cmd and subcmd parts
func (cf *CommandFormatter) highlightCommandsAndSubcommands(text string) string {
	// Split on spaces to analyze word by word
	words := strings.Fields(text)
	result := make([]string, len(words))

	engxFound := false
	cmdFound := false

	for i, word := range words {
		// Skip if word already has ANSI codes (already processed)
		if strings.Contains(word, "\033[") {
			result[i] = word
			if strings.Contains(word, "engx") {
				engxFound = true
			}
			continue
		}

		// Look for the sequence: engx -> cmd -> subcmd
		if word == "engx" {
			engxFound = true
			result[i] = word // Already colored above
		} else if engxFound && !cmdFound && isCommand(word) {
			// This is the command (first word after engx)
			result[i] = fmt.Sprintf("\033[38;5;%sm%s\033[0m", design.ColorOrange, word)
			cmdFound = true
		} else if engxFound && cmdFound && isSubcommand(word) {
			// This is the subcommand (second word after engx)
			result[i] = fmt.Sprintf("\033[%sm%s\033[0m", design.ColorBrightYellow, word)
		} else {
			result[i] = word
		}
	}

	return strings.Join(result, " ")
}

// isCommand checks if a word looks like a command
func isCommand(word string) bool {
	// Commands are typically alphanumeric words without special characters
	return regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`).MatchString(word)
}

// isSubcommand checks if a word looks like a subcommand
func isSubcommand(word string) bool {
	// Subcommands are similar to commands
	return regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]*$`).MatchString(word)
}

// FormatCommandInBackticks wraps the formatted command in backticks with magenta color
func (cf *CommandFormatter) FormatCommandInBackticks(command string) string {
	formattedCommand := cf.FormatCommand(command)
	return fmt.Sprintf("\033[%sm`\033[0m%s\033[%sm`\033[0m",
		design.ColorBrightMagenta, formattedCommand, design.ColorBrightMagenta)
}

// Convenience function for quick command formatting
func FormatEngxCommand(command string) string {
	formatter := NewCommandFormatter()
	return formatter.FormatCommand(command)
}

// Convenience function for backtick-wrapped command formatting
func FormatEngxCommandInBackticks(command string) string {
	formatter := NewCommandFormatter()
	return formatter.FormatCommandInBackticks(command)
}