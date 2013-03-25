package mailthread

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

var SkipParseLeftError = errors.New("skip parse left")

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
	Skip(r io.Reader, w io.Writer) (err error)
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

func (dch *GoThroughHandler) Skip(r io.Reader, w io.Writer) (err error) {
	io.Copy(w, r)
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

func (dch *DefaultContentHandler) Skip(r io.Reader, w io.Writer) (err error) {
	io.Copy(w, r)
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
	var exit, eof bool
	for {
		if eof {
			return
		}

		l, err = r.ReadSlice('\n')

		if err != nil && err != io.EOF {
			return
		}

		if err == io.EOF {
			eof = true
		}

		if len(l) > 1 && l[len(l)-2] == '\r' {
			// line = strings.Replace(line, "\r\n", "\n", -1)
			l = append(l[:len(l)-2], '\n')
		}

		line := string(l)
		// fmt.Printf("%q\n", line)
		buffer.parseIn(line)
		if buffer.atHeadStart {
			// ch.Text(output, buffer.content)
			buffer.clean()
			continue
		}
		if buffer.atHeadEnd {
			switch buffer.bType {
			case fw_type:
				err = ch.StartForward(output)
				if err, exit = handleError(err, ch, r, output, endtaggers, ""); exit {
					return
				}
				endtaggers = append(endtaggers, ContentHandler.EndForward)

				err = ch.ForwardHeader(output, buffer.content)
				if err, exit = handleError(err, ch, r, output, endtaggers, buffer.headEndLineContent); exit {
					return
				}
			case re_type:
				err = ch.StartReply(output)
				if err, exit = handleError(err, ch, r, output, endtaggers, ""); exit {
					return
				}
				endtaggers = append(endtaggers, ContentHandler.EndReply)

				err = ch.ReplyHeader(output, buffer.content)
				if err, exit = handleError(err, ch, r, output, endtaggers, buffer.headEndLineContent); exit {
					return
				}
			}

			buffer.clean()
			err = ch.Text(output, buffer.headEndLineContent)
			if err, exit = handleError(err, ch, r, output, endtaggers, ""); exit {
				return
			}
			continue
		}

		if buffer.inHead {
			if !buffer.legalHeadContent {
				err = ch.Text(output, string(buffer.content))
				buffer.rewind()
			}

			continue
		}

		err = ch.Text(output, string(l))
		if err == nil && eof {
			err = io.EOF
		}
		if err, exit = handleError(err, ch, r, output, endtaggers, ""); exit {
			return
		}
	}

	return
}

func handleError(err error, ch ContentHandler, r io.Reader, output io.Writer, endtaggers []endTagger, headEndLineContent string) (newerr error, exit bool) {
	if err == nil {
		return nil, false
	}

	if len(headEndLineContent) > 0 {
		ch.Text(output, headEndLineContent)
	}

	if err == SkipParseLeftError {
		ch.Skip(r, output)
	}

	if err == SkipParseLeftError || err == io.EOF {
		for i := len(endtaggers) - 1; i >= 0; i-- {
			endtaggers[i](ch, output)
		}
		err = ch.EndThread(output)
		return nil, true
	}

	return err, true
}
