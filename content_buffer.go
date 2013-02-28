package mailthread

import (
    "regexp"
    "fmt"
)

type contentBuffer struct {
    lastLine string
    content string
    atHeadStart bool
    atHeadEnd bool
    inHead bool
    bType string // fw | re
}

const (
    re_type = "re"
    fw_type = "fw"
)

func (buffer *contentBuffer) push(line string) {
    buffer.content += buffer.lastLine
    buffer.lastLine = line
}

func (buffer *contentBuffer) clear() {
    buffer.content = buffer.lastLine
    buffer.lastLine = ""
    buffer.bType = ""
    buffer.inHead = false
}

func (buffer *contentBuffer) parseLastLine() {
    buffer.atHeadStart = false
    buffer.atHeadEnd = false
    
    switch {
    case buffer.isFwHeadStart():
        buffer.lastLine = fw_type
        
        buffer.atHeadStart = true
        buffer.inHead = true
    case buffer.isReHeadStart():
        buffer.lastLine = re_type
        
        buffer.atHeadStart = true
        buffer.inHead = true
    case buffer.inHead && buffer.isHeadEnd():
        
        buffer.atHeadEnd = true
    }
}

var fwHeadStartExps = func() []*regexp.Regexp {
    expr1, err := regexp.Compile(`^-{10} Forwarded message -{10}\n`)
    if err != nil { 
        panic(err)
    }

    expr2, err := regexp.Compile(`^-{5} Forwarded message -{5}\n`)
    if err != nil { 
        panic(err)
    }
    
    return []*regexp.Regexp{expr1, expr2}
}()

func (buffer *contentBuffer) isFwHeadStart() (result bool) {
    // for _, exp := range fwHeadStartExps {
    //     if exp.MatchString(buffer.lastLine) {
    //         result = true
    //         break
    //     }
    // }
    return
}

func (buffer *contentBuffer) isReHeadStart() (result bool) {
    // for _, exp := range reHeadStartExps {
    //     if exp.MatchString(buffer.lastLine) {
    //         result = true
    //         break
    //     }
    // }
    return
}

var headEndExp = func() *regexp.Regexp {
    expr, err := regexp.Compile(`^\n`)
    if err != nil { 
        panic(err)
    }
    
    return expr
}()

func (buffer *contentBuffer) isHeadEnd() (bool) {
    return headEndExp.MatchString(buffer.lastLine)
}