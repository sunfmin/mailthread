package mailthread

import (
    "regexp"
)

// ===== forwarding block regexpes ===== //
var forwardingBlockStart = func() *regexp.Regexp {
    expr, err := regexp.Compile(`-{10} Forwarded message -{10}\n`)
    if err != nil { 
        panic(err)
    }
    
    return expr
}()

func isForwardingBlockStart(str string) bool {
    return forwardingBlockStart.MatchString(str)
}


var forwardingBlockEnd = func() *regexp.Regexp {
    expr, err := regexp.Compile(`^\n`)
    if err != nil { 
        panic(err)
    }
    
    return expr
}()

func isForwardingBlockEnd(str string) bool {
    return forwardingBlockEnd.MatchString(str)
}
