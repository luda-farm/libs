package terminal

import (
	"os"
	"strings"

	"github.com/luda-farm/libs/std"
	"golang.org/x/term"
)

const (
	right    = "\033[C"
	left     = "\033[D"
	cleareos = "\033[J"
)

var (
	history []string
	index   int
	cursor  int
)

func init() {
	history = make([]string, 1)
}

func Listen(callback func(string)) {
	defer term.Restore(int(os.Stdin.Fd()), std.Must(term.MakeRaw(int(os.Stdin.Fd()))))
	buffer := make([]byte, 1)
	for {
		std.Must(os.Stdin.Read(buffer))
		switch buffer[0] {
		case 3, 4: // C-c, C-d
			return
		case 13: // Enter
			editHistoryCheck()
			os.Stdout.WriteString("\n")
			history = append(make([]string, 1), history...)
			resetCursor()
			callback(history[1])
		case 27: // \033
			std.Must(os.Stdin.Read(buffer))
			switch buffer[0] {
			case 91: // [
				std.Must(os.Stdin.Read(buffer))
				switch buffer[0] {
				case 65: // Up
					if index < len(history)-1 {
						index++
						resetCursor()
					}
				case 66: // Down
					if index > 0 {
						index--
						resetCursor()
					}
				case 67: // ->
					if cursor < len(history[index]) {
						os.Stdout.WriteString(right)
						cursor++
					}
				case 68: // <-
					if cursor > 0 {
						os.Stdout.WriteString(left)
						cursor--
					}
				}
			}
		case 127: // Backspace
			editHistoryCheck()
			if cursor == 0 {
				break
			}
			history[0] = history[0][:cursor-1] + history[0][cursor:]
			os.Stdout.WriteString(left)
			cursor--
		default:
			editHistoryCheck()
			history[0] = history[0][:cursor] + string(buffer) + history[0][cursor:]
			cursor++
		}

		os.Stdout.WriteString(strings.Repeat(left, cursor)) // Cursor to start of line
		os.Stdout.WriteString(cleareos)
		os.Stdout.WriteString(history[index])
		os.Stdout.WriteString(strings.Repeat(left, len(history[index])-cursor)) // Reset cursor
	}
}

func editHistoryCheck() {
	if index != 0 {
		history[0] = history[index]
		index = 0
	}
}

func resetCursor() {
	os.Stdout.WriteString(strings.Repeat(left, cursor))
	cursor = len(history[index])
	os.Stdout.WriteString(strings.Repeat(right, cursor))
}

func Write(s string) {
	for _, line := range strings.Split(s, "\n") {
		os.Stdout.WriteString(line + "\n")
		os.Stdout.WriteString(strings.Repeat(left, len(line))) // Cursor to start of line
	}
}
