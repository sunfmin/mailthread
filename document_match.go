package mailthread

import (
    // "regexp"
    // "fmt"
    "strings"
)

func Process(input string) (output string, err error) {
    inputSlice := strings.Split(input, "\n")
    
    for _, line := range inputSlice {
        switch {
        case isForwardingBlockStart(line):
            
        case isForwardingBlockEnd(line):
            
        case isReplyLine(line):
            
        } 
    }
    
    return
}