package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

// StringMap is used to handle multiple flag.
type StringMap map[string]string

func (s *StringMap) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *StringMap) Get() interface{} { return StringMap(*s) }

func (s *StringMap) Set(value string) error {
	if (*s) == nil {
		*s = make(StringMap)
	}

	// Backwards compatibility
	if strings.Contains(value, ",") {
		for _, keyValue := range strings.Split(value, ",") {
			s.Set(keyValue)
		}
		return nil
	}

	i := strings.Index(value, ":")
	if i != -1 && value[i+1:] != "" {
		(*s)[value[:i]] = value[i+1:]
	}

	return nil
}

func newFlagSet(name string, errorHandling flag.ErrorHandling) *flag.FlagSet {
	flag := flag.NewFlagSet(name, errorHandling)
	flag.SetOutput(ioutil.Discard)
	return flag
}
