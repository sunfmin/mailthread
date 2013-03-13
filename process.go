package mailthread

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ContentHandler interface {
	Text(w io.Writer, text string) (err error)
	StartThread(w io.Writer) (err error)
	EndThread(w io.Writer) (err error)
	StartForward(w io.Writer) (err error)
	ForwardHeader(w io.Writer, header string) (err error)
	EndForward(w io.Writer) (err error)
	StartReply(w io.Writer) (err error)
	ReplyHeader(w io.Writer, header string) (err error)
	EndReply(w io.Writer) (err error)
}

type GoThroughHandler struct {
}

func (dch *GoThroughHandler) Text(w io.Writer, text string) (err error) {
	_, err = io.WriteString(w, text)
	return
}

func (dch *GoThroughHandler) StartThread(w io.Writer) (err error) {
	return
}

func (dch *GoThroughHandler) EndThread(w io.Writer) (err error) {
	return
}

func (dch *GoThroughHandler) StartForward(w io.Writer) (err error) {
	return
}

func (dch *GoThroughHandler) EndForward(w io.Writer) (err error) {
	return
}

func (dch *GoThroughHandler) ForwardHeader(w io.Writer, header string) (err error) {
	_, err = io.WriteString(w, header)
	return
}

func (dch *GoThroughHandler) StartReply(w io.Writer) (err error) {
	return
}

func (dch *GoThroughHandler) EndReply(w io.Writer) (err error) {
	return
}

func (dch *GoThroughHandler) ReplyHeader(w io.Writer, header string) (err error) {
	_, err = io.WriteString(w, header)
	return
}

type DefaultContentHandler struct {
}

func (dch *DefaultContentHandler) Text(w io.Writer, text string) (err error) {
	_, err = io.WriteString(w, text)
	return
}

func (dch *DefaultContentHandler) StartThread(w io.Writer) (err error) {
	_, err = io.WriteString(w, "<div class=\"message\">\n")
	return
}

func (dch *DefaultContentHandler) EndThread(w io.Writer) (err error) {
	_, err = io.WriteString(w, "\n</div>")
	return
}

func (dch *DefaultContentHandler) StartForward(w io.Writer) (err error) {
	_, err = io.WriteString(w, "<div class=\"forward\">\n")
	return
}

func (dch *DefaultContentHandler) EndForward(w io.Writer) (err error) {
	_, err = io.WriteString(w, "\n</div>")
	return
}

func (dch *DefaultContentHandler) ForwardHeader(w io.Writer, header string) (err error) {
	_, err = fmt.Fprintf(w, "<div class=\"forwarded_message_header\">\n%s</div>\n", header)
	return
}

func (dch *DefaultContentHandler) StartReply(w io.Writer) (err error) {
	_, err = io.WriteString(w, "<div class=\"reply\">\n")
	return
}

func (dch *DefaultContentHandler) EndReply(w io.Writer) (err error) {
	_, err = io.WriteString(w, "\n</div>")
	return
}

func (dch *DefaultContentHandler) ReplyHeader(w io.Writer, header string) (err error) {
	_, err = fmt.Fprintf(w, "<div class=\"reply_header\">\n%s</div>\n", header)
	return
}

func ProcessString(input string) (output string) {
	ch := &DefaultContentHandler{}
	return ProcessStringWithHandler(input, ch)
}

func ProcessStringWithHandler(input string, ch ContentHandler) (output string) {
	var err error
	out := bytes.NewBuffer(nil)
	err = Process(strings.NewReader(input), out, ch)
	if err != nil {
		return input
	}
	output = out.String()
	return
}

type endTagger func(ch ContentHandler, w io.Writer) error

func Process(input io.Reader, output io.Writer, ch ContentHandler) (err error) {

	r := bufio.NewReader(input)

	err = ch.StartThread(output)
	if err != nil {
		return
	}

	endtaggers := make([]endTagger, 0)
	buffer := contentBuffer{}
	var l []byte
	for {
		l, err = r.ReadSlice('\n')

		if err != nil && err != io.EOF {
			return
		}

		buffer.parseIn(string(l))
		if buffer.atHeadStart {
			// ch.Text(output, buffer.content)
			buffer.clear()
			continue
		}
		if buffer.atHeadEnd {
			switch buffer.bType {
			case fw_type:
				err = ch.StartForward(output)
				if err != nil {
					return
				}
				err = ch.ForwardHeader(output, buffer.content)
				if err != nil {
					return
				}
				endtaggers = append(endtaggers, ContentHandler.EndForward)
			case re_type:
				err = ch.StartReply(output)
				if err != nil {
					return
				}
				err = ch.ReplyHeader(output, buffer.content)
				if err != nil {
					return
				}
				endtaggers = append(endtaggers, ContentHandler.EndReply)
			}

			buffer.clear()
			io.WriteString(output, "\n")
			continue
		}

		if !buffer.inHead {
			ch.Text(output, string(l))
		}

		if err == io.EOF {
			err = nil
			break
		}
	}

	// if buffer.content != "" {
	// 	err = ch.Text(output, buffer.content+buffer.lastLine)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	for i := len(endtaggers) - 1; i >= 0; i-- {
		endtaggers[i](ch, output)
	}
	err = ch.EndThread(output)

	return
}
