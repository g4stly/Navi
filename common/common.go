package common

import (
	"flag"
	"fmt"
	"io"
	"os"
	"io/ioutil"
	"errors"
	"encoding/json"
)

func init() {
	flag.Parse()
	Args = flag.Args()

	var err error
	Config, err = loadFile(*configLocation)
	if err != nil {
		Fatal("common.init(): %v", err)
	}
}

// important stuffs
var Args []string
var Config map[string]interface{}

var silent		= flag.Bool("s", false, "silent: surpress all output")
var verbose		= flag.Bool("v", true, "verbose: print debug output")
var configLocation	= flag.String("f", "config.json", "config-file: path to configuration")

// config
func loadFile(fileName string) (map[string]interface{}, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		return nil, err
	}
	return jsonData.(map[string]interface{}), nil
}

// makes error handling a little sexier
func NewError(fmtstring string, err error) error {
	return errors.New(fmt.Sprintf(fmtstring, err))
}

// logging
func msg(w io.Writer, badge string, fmtstring string, args ...interface{}) {
	if *silent { return }

	if len(args) < 1 {
		w.Write([]byte(fmt.Sprintf("%v: %v\n", badge, fmtstring)))
		return
	}

	w.Write([]byte(fmt.Sprintf("%v: %v\n", badge, fmt.Sprintf(fmtstring, args...))))
}

// for debug only output
func Log(fmtstring string, args ...interface{}) {
	if !*verbose { return }
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

// fatal out
func Fatal(fmtstring string, args ...interface{}) {
	msg(os.Stderr, "FATAL", fmtstring, args...)
	os.Exit(-1)
}










