package mailthread

import (
    // "regexp"
    "fmt"
    "strings"
)

func Process(input string) (output string, err error) {
    inputSlice := strings.SplitAfter(input, "\n")
    endTagCounter := 1
    
    output += "<div class=\"message\">\n"
    insideForWardingBlock := false
    for _, line := range inputSlice {
        switch {
        case isForwardingBlockStart(line):
            output += "<div class=\"forwarded_message_header\">\n"
            output += line
            insideForWardingBlock = true
        case insideForWardingBlock && isForwardingBlockEnd(line):
            output += line
            output += "</div>\n"
            insideForWardingBlock = false
        case isReplyLine(line):
            endTagCounter += 1
            output += fmt.Sprintf("<div class=\"reply\">\n<div class=\"reply_header\">\n%s</div>\n", line)
        default:
            output += line
        }
    }
    
    output += strings.Repeat("\n</div>\n", endTagCounter)
    
    return
}