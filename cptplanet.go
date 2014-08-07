package cptplanet

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

type logger bool

func (l *logger) Printf(s string, v ...interface{}) {
	if *l {
		log.Printf(s, v...)
	}
}

var l = logger(false)

type Value flag.Value

type EnvSet struct {
	prefix string
	*flag.FlagSet
}

func NewEnvironment(prefix string) *EnvSet {
	return &EnvSet{
		prefix:  prefix,
		FlagSet: flag.NewFlagSet(prefix, flag.ExitOnError),
	}
}

func (e *EnvSet) Parse() error {
	for _, arg := range os.Environ() {
		if strings.HasPrefix(arg, e.prefix) {
			// split env on "="
			s := strings.Split(arg, "=")

			// remove the prefix from the key
			key := strings.TrimPrefix(s[0], e.prefix)
			value := s[1]

			l.Printf("setting %s => %s", key, value)

			err := e.Set(key, value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// default prefix is the executable name
var Environment = NewEnvironment(getAppName(os.Args[0]) + "_")

func Int(name string, value int, usage string) *int {
	return Environment.Int(name, value, usage)
}

func IntVar(p *int, name string, value int, usage string) {
	Environment.IntVar(p, name, value, usage)
}

func String(name string, value string, usage string) *string {
	return Environment.String(name, value, usage)
}

func StringVar(p *string, name string, value string, usage string) {
	Environment.StringVar(p, name, value, usage)
}

func Bool(name string, value bool, usage string) *bool {
	return Environment.Bool(name, value, usage)
}

func BoolVar(p *bool, name string, value bool, usage string) {
	Environment.BoolVar(p, name, value, usage)
}

func Duration(name string, value time.Duration, usage string) *time.Duration {
	return Environment.Duration(name, value, usage)
}

func DurationVar(p *time.Duration, name string, value time.Duration, usage string) {
	Environment.DurationVar(p, name, value, usage)
}

func Var(value Value, name string, usage string) {
	Environment.Var(value, name, usage)
}

func Parse() error {
	return Environment.Parse()
}

func getAppName(appName string) string {
	return strings.ToUpper(path.Base(appName))
}