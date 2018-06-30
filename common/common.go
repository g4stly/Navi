package common

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func init() {
	flag.Parse()
	Args = flag.Args()

	err := LoadConfig()
	if err != nil {
		Fatal("common.init(): %v", err)
	}
}

// important stuffs
var Args []string
var Config map[string]interface{}

var silent = flag.Bool("s", false, "silent: surpress all output")
var verbose = flag.Bool("v", true, "verbose: print debug output")
var configLocation = flag.String("f", "/etc/navi.conf", "config-file: path to configuration")

// config
func LoadConfig() error {
	file, err := ioutil.ReadFile(*configLocation)
	if err != nil {
		return err
	}

	var jsonData interface{}
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		return err
	}
	Config = jsonData.(map[string]interface{})
	return nil
}

// makes error handling a little sexier
func NewError(fmtstring string, err error) error {
	return errors.New(fmt.Sprintf(fmtstring, err))
}

// logging
func msg(w io.Writer, badge string, fmtstring string, args ...interface{}) {
	if *silent {
		return
	}

	if len(args) < 1 {
		w.Write([]byte(fmt.Sprintf("%v: %v\n", badge, fmtstring)))
		return
	}

	w.Write([]byte(fmt.Sprintf("%v: %v\n", badge, fmt.Sprintf(fmtstring, args...))))
}

// for debug only output
func Log(fmtstring string, args ...interface{}) {
	if !*verbose {
		return
	}
	msg(os.Stdout, "DEBUG", fmtstring, args...)
}

// not necesarily debug output
func Out(fmtstring string, args ...interface{}) {
	msg(os.Stdout, "info", fmtstring, args...)
}

// print usage
func Usage(fmtstring string, args ...interface{}) {
	msg(os.Stdout, "usage", fmtstring, args...)
}

// error msg
func Error(fmtstring string, args ...interface{}) {
	if !*verbose {
		return
	}
	msg(os.Stdout, "ERROR", fmtstring, args...)
}

// fatal out
func Fatal(fmtstring string, args ...interface{}) {
	msg(os.Stderr, "FATAL", fmtstring, args...)
	os.Exit(-1)
}
