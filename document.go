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
    buffer := contentBuffer{}
    // insideForWardingHead := false
    for _, line := range inputSlice {
        buffer.push(line)
        
        buffer.parseLastLine()
        
        if buffer.atHeadStart {
            output += string(blackfriday.MarkdownBasic([]byte(buffer.content)))
            buffer.clear()
        }
        
        if buffer.atHeadEnd {
            output += string(blackfriday.MarkdownBasic([]byte(buffer.content)))
            
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
        
        // switch {
        // case isContentHeadStart(line):
        //     if contentBuffer != "" {
        //         output += string(blackfriday.MarkdownBasic([]byte(contentBuffer)))
        //         contentBuffer = ""
        //     }
        //     
        //     contentBuffer += line
        // case isCotentHeadEnd(line):
        
        // case isForwardingHeadStart(line):
        //     if contentBuffer != "" {
        //         output += string(blackfriday.MarkdownBasic([]byte(contentBuffer)))
        //         contentBuffer = ""
        //     }
        //     output += "<div class=\"forwarded_message_header\">\n"
        //     
        //     contentBuffer += line
        //     
        //     insideForWardingHead = true
        // case insideForWardingHead && isForwardingHeadEnd(line):
        //     output += string(blackfriday.MarkdownBasic([]byte(contentBuffer)))
        //     output += "</div>\n"
        // 
        //     contentBuffer = line
        // 
        //     insideForWardingHead = false
        // case isReplyLine(line):
        //     if contentBuffer != "" {
        //         output += string(blackfriday.MarkdownBasic([]byte(contentBuffer)))
        //         contentBuffer = ""
        //     }
        //     lineInHtml := string(blackfriday.MarkdownBasic([]byte(line)))
        //     output += fmt.Sprintf("<div class=\"reply\">\n<div class=\"reply_header\">\n%s</div>\n", lineInHtml)
        //     
        //     endTagCounter += 1
        // default:
        //     contentBuffer += line
        // }
    }
    
    if buffer.content != "" {
        output += string(blackfriday.MarkdownBasic([]byte(buffer.content)))
    }
    
    output += strings.Repeat("\n</div>\n", endTagCounter)
    
    return
}