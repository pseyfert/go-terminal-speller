package terminalspeller

import (
	_ "bufio"
	_ "bytes"
	// "golang.org/x/text/unicode/norm"
	"golang.org/x/text/unicode/runenames"
	"io"
	"strings"
	"unicode"
)

type translator struct {
	w io.Writer
}

func NewTranslator(writer io.Writer) io.Writer {
	retval := translator{
		w: writer,
	}
	return retval
}

func (myself translator) Write(p []byte) (int, error) {
	p_string := string(p)
	var retval int
	for _, r := range p_string {
		if unicode.IsSymbol(r) {
			i, err := myself.w.Write([]byte{':'})
			if err != nil {
				// FIXME: check 0
				return 0, err
			}
			retval += i
			i, err = myself.w.Write([]byte(strings.ReplaceAll(strings.ToLower(runenames.Name(r)), " ", "_")))
			if err != nil {
				// FIXME: check 0
				return 0, err
			}
			retval += i
			i, err = myself.w.Write([]byte{':'})
			if err != nil {
				// FIXME: check 0
				return 0, err
			}
			retval += i
		} else {
			i, err := myself.w.Write([]byte(string(r)))
			if err != nil {
				// FIXME: check 0
				return 0, err
			}
			retval += i
		}
	}
	return retval, nil
}
