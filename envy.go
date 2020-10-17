package envy

import (
	"flag"
	"os"
	"strings"
	"sync"
)

type Env interface {
	Getenv(key string) string
	LookupEnv(key string) (string, bool)
	Setenv(key string, value string) error
	Unsetenv(key string) error
	Clearenv()
	Environ() []string
	ExpandEnv(key string) string
	Expand(key string, mapping func(string) string) string
}

// DefaultMappingconvert string to uppercase and replace any occurence of - with _
func DefaultMapping(prefix string) func(string) string {
	return func(name string) string {
		return strings.ReplaceAll(strings.ToUpper(prefix+name), "-", "_")
	}
}

// FillFlags is equivalent to
// FillFlagsLookup(f , NewOsEnv(), DefaultMapping(""))
func FillFlags(f *flag.FlagSet) error {
	return FillFlagsLookup(f, NewOsEnv(), DefaultMapping(""))
}

// FillFlagsLookup set values in a flag.FlagSet with environment variables
// variable names are get calling mapping(flag.Name). An error is returned if
// flag.Value.Set return an error. This method should be called before flag.Parse
func FillFlagsLookup(f *flag.FlagSet, e Env, mapping func(string) string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	setted := make(map[string]struct{})
	f.Visit(func(fn *flag.Flag) {
		setted[fn.Name] = struct{}{}
	})
	f.VisitAll(func(fn *flag.Flag) {
		if _, isSet := setted[fn.Name]; !isSet {
			name := mapping(fn.Name)
			if v, ok := e.LookupEnv(name); ok {
				if err := fn.Value.Set(v); err != nil {
					panic(err)
				}
			}
		}
	})
	return nil
}

var _ Env = (*MapEnv)(nil)

type MapEnv struct {
	l   sync.RWMutex
	env map[string]string
}

func (m *MapEnv) Getenv(key string) string {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.env[key]
}

func (m *MapEnv) LookupEnv(key string) (string, bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	v, ok := m.env[key]
	return v, ok
}
func (m *MapEnv) Setenv(key, value string) error {
	m.l.Lock()
	defer m.l.Unlock()
	m.env[key] = value
	return nil
}
func (m *MapEnv) Unsetenv(key string) error {
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.env, key)
	return nil
}
func (m *MapEnv) Clearenv() {
	m.l.Lock()
	defer m.l.Unlock()
	m.env = map[string]string{}
}
func (m *MapEnv) Environ() []string {
	list := make([]string, 0, len(m.env))
	for k, v := range m.env {
		list = append(list, strings.Join([]string{k, v}, "="))
	}
	return list
}
func (m *MapEnv) Expand(key string, mapping func(string) string) string {
	return os.Expand(key, mapping)
}
func (m *MapEnv) ExpandEnv(key string) string {
	return os.Expand(key, m.Getenv)
}

func NewMapEnv(e map[string]string) *MapEnv {
	m := make(map[string]string, len(e))
	for k, v := range e {
		m[k] = v
	}
	return &MapEnv{env: m}
}

var _ Env = (*OsEnv)(nil)

type OsEnv struct{}

func (o *OsEnv) Getenv(key string) string {
	return os.Getenv(key)
}
func (o *OsEnv) LookupEnv(key string) (string, bool) {
	return os.LookupEnv(key)
}
func (o *OsEnv) Setenv(key, value string) error {
	return os.Setenv(key, value)
}
func (o *OsEnv) Unsetenv(key string) error {
	return os.Unsetenv(key)
}
func (o *OsEnv) Clearenv() {
	os.Clearenv()
}
func (o *OsEnv) Environ() []string {
	return os.Environ()
}
func (o *OsEnv) Expand(key string, mapping func(string) string) string {
	return os.Expand(key, mapping)
}
func (o *OsEnv) ExpandEnv(key string) string {
	return os.ExpandEnv(key)
}
func NewOsEnv() *OsEnv {
	return &OsEnv{}
}
