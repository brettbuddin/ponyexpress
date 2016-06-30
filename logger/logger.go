package logger

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))

func Infof(format string, args ...interface{}) {
	log.Printf("INFO -- "+format, args...)
}

func Debugf(format string, args ...interface{}) {
	if debug {
		log.Printf("DEBUG -- "+format, args...)
	}
}

func Errorf(format string, args ...interface{}) {
	log.Printf("ERROR -- "+format, args...)
}

type Fields map[string]interface{}

func (f Fields) String() string {
	output := []string{}
	for k, v := range f {
		output = append(output, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(output, " ")
}
