package mailthread

import (
	// "regexp"
	"fmt"
	"github.com/russross/blackfriday"
	"strings"
)

func Process(input string) (output string, err error) {
	inputSlice := strings.SplitAfter(input, "\n")
	endTagCounter := 1

	output += "<div class=\"message\">\n"
	buffer := contentBuffer{}
	for _, line := range inputSlice {
		buffer.push(line)

		buffer.parseLastLine()

		if buffer.atHeadStart {
			output += string(blackfriday.MarkdownBasic([]byte(buffer.content)))
			buffer.clear()
		}

		if buffer.atHeadEnd {
			switch buffer.bType {
			case fw_type:
				output += fmt.Sprintf(
					"<div class=\"forwarded_message_header\">\n%s</div>\n",
					string(blackfriday.MarkdownBasic([]byte(buffer.content))),
				)
			case re_type:
				output += fmt.Sprintf(
					"<div class=\"reply\">\n<div class=\"reply_header\">\n%s</div>\n",
					string(blackfriday.MarkdownBasic([]byte(buffer.content))),
				)

				endTagCounter += 1
			}

			buffer.clear()
		}
	}

	if buffer.content != "" {
		output += string(blackfriday.MarkdownBasic([]byte(buffer.content + buffer.lastLine)))
	}

	output += strings.Repeat("\n</div>\n", endTagCounter)

	return
}
