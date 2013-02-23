package mailthread

import (
    // "regexp"
    "fmt"
    "strings"
    "github.com/russross/blackfriday"
)

func Process(input string) (output string, err error) {
    inputSlice := strings.SplitAfter(input, "\n")
    endTagCounter := 1
    
    output += "<div class=\"message\">\n"
    contentBuffer := ""
    insideForWardingBlock := false
    for _, line := range inputSlice {
        switch {
        case isForwardingBlockStart(line):
            if contentBuffer != "" {
                output += string(blackfriday.MarkdownBasic([]byte(contentBuffer)))
                contentBuffer = ""
            }
            output += "<div class=\"forwarded_message_header\">\n"
            
            contentBuffer += line
            
            insideForWardingBlock = true
        case insideForWardingBlock && isForwardingBlockEnd(line):
            output += string(blackfriday.MarkdownBasic([]byte(contentBuffer)))
            output += "</div>\n"

            contentBuffer = line

            insideForWardingBlock = false
        case isReplyLine(line):
            if contentBuffer != "" {
                output += string(blackfriday.MarkdownBasic([]byte(contentBuffer)))
                contentBuffer = ""
            }
            lineInHtml := string(blackfriday.MarkdownBasic([]byte(line)))
            output += fmt.Sprintf("<div class=\"reply\">\n<div class=\"reply_header\">\n%s</div>\n", lineInHtml)
            
            endTagCounter += 1
        default:
            contentBuffer += line
        }
    }
    
    if contentBuffer != "" {
        output += string(blackfriday.MarkdownBasic([]byte(contentBuffer)))
    }
    
    output += strings.Repeat("\n</div>\n", endTagCounter)
    
    return
}