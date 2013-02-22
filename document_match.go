package mailthread

import (
    // "regexp"
    "fmt"
    "strings"
)

func Process(input string) (output string, err error) {
    inputSlice := strings.SplitAfter(input, "\n")
    endTagCounter := 1
    // * remove `\n` match in simple_forward.go
    // * endTagStack
    // * output
        
    // append top tag
    // push outter end tag into endTagStack
    // ... process
    // pop end tag
    
    output += "<div class=\"message\">\n"
    insideForWardBlock := false
    for _, line := range inputSlice {
        switch {
        case isForwardingBlockStart(line):
            output += "<div class=\"forwarded_message_header\">\n"
            output += line
            insideForWardBlock = true
        case insideForWardBlock && isForwardingBlockEnd(line):
            output += line
            output += "</div>\n"
            insideForWardBlock = false
        case isReplyLine(line):
            endTagCounter += 1
            output += fmt.Sprintf("<div class=\"reply\"><div class=\"reply_header\">\n%s\n</div>\n", line)
        default:
            output += line
        }
    }
    
    output += strings.Repeat("\n</div>\n", endTagCounter)
    
    return
}