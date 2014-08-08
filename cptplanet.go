package cptplanet

import (
	"flag"
	"fmt"
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

type ParseErr struct {
	MissingKeys []string
	ExtraKeys   []string
	ParseErrors []string
}

func (p *ParseErr) isError() bool {
	return len(p.MissingKeys) > 0 || len(p.ExtraKeys) > 0 || len(p.ParseErrors) > 0
}

func (p *ParseErr) addMissing(key string) {
	if p.MissingKeys == nil {
		p.MissingKeys = make([]string, 0, 1)
	}

	p.MissingKeys = append(p.MissingKeys, key)
}

func (p *ParseErr) addExtra(key string) {
	if p.ExtraKeys == nil {
		p.ExtraKeys = make([]string, 0, 1)
	}

	p.ExtraKeys = append(p.ExtraKeys, key)
}

func (p *ParseErr) addParseErr(key string) {
	if p.ParseErrors == nil {
		p.ParseErrors = make([]string, 0, 1)
	}

	p.ParseErrors = append(p.ParseErrors, key)
}

func (p *ParseErr) Error() string {
	return fmt.Sprintf("Missing keys: %v, Extra keys: %v, Parse Errors: %v", p.MissingKeys, p.ExtraKeys, p.ParseErrors)
}

type Value flag.Value

type Settings struct {
	Prefix             string
	ErrorOnExtraKeys   bool
	ErrorOnMissingKeys bool
	ErrorOnParseErrors bool
}

var DefaultSettings = Settings{Prefix: getAppName(os.Args[0]) + "_"}

type EnvSet struct {
	Settings

	*flag.FlagSet

	visited map[string]bool
}

func NewEnvironment(settings Settings) *EnvSet {
	return &EnvSet{
		Settings: settings,
		FlagSet:  flag.NewFlagSet("", flag.ExitOnError),
		visited:  make(map[string]bool, 0),
	}
}

func (e *EnvSet) PrintDefaults() {
	fmt.Printf("## Example Usage:\n")
	e.FlagSet.VisitAll(func(f *flag.Flag) {
		fmt.Printf("# %s\nexport %s%s=%q\n", f.Usage, e.Prefix, f.Name, f.DefValue)
	})
}

func (e *EnvSet) Parse() error {
	parseErr := &ParseErr{}

	for _, arg := range os.Environ() {
		if strings.HasPrefix(arg, e.Prefix) {
			// split env on "="
			s := strings.SplitN(arg, "=", 2)

			envKey := s[0]

			// remove the Prefix from the key
			key := strings.TrimPrefix(envKey, e.Prefix)
			value := s[1]

			l.Printf("setting %s => %s", key, value)

			err := e.Set(key, value)
			if err != nil {

				if strings.HasPrefix(err.Error(), "no such flag -") {
					if e.ErrorOnExtraKeys {
						parseErr.addExtra(envKey)
					}

					continue
				}

				if strings.Contains(err.Error(), "invalid syntax") {
					// replace the default value
					originalFlag := e.Lookup(key)
					e.Set(key, originalFlag.DefValue)

					if e.ErrorOnParseErrors {
						parseErr.addParseErr(envKey)
					}

					continue
				}

				return err
			}

			e.visited[envKey] = true
		}
	}

	if e.ErrorOnMissingKeys {
		e.VisitAll(func(f *flag.Flag) {
			if !e.visited[e.Prefix+f.Name] {
				parseErr.addMissing(e.Prefix + f.Name)
			}
		})
	}

	if parseErr.isError() {
		return parseErr
	}

	return nil
}

// default Prefix is the executable name
var Environment = NewEnvironment(DefaultSettings)

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
