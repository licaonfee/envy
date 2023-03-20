package envy

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

type VarType int

const (
	VarTypeUnknown = VarType(iota)
	VarTypeString
	VarTypeMap
	VarTypeArray
)

// FilterPrefix returns the name and type only if the environment variable have the given preffix
// it returns as name , the lower case version of the variable with preffix stripped
func FilterPrefix(preffix string) func(string) (string, VarType) {
	isArray := regexp.MustCompile(`_[0-9]*$`)
	return func(s string) (string, VarType) {
		ok := strings.HasPrefix(s, preffix)
		if !ok {
			return "", VarTypeUnknown
		}
		s = strings.ToLower(strings.TrimPrefix(s, preffix))
		switch {
		case isArray.MatchString(s):
			return s, VarTypeArray
		default:
			return s, VarTypeString
		}
	}
}

// FillMap return a map with all environment variables that match filter function
func FillMap(e Env, filter func(string) (string, VarType)) map[string]any {
	values := make(map[string]any)
	for _, v := range e.Environ() {
		key, value, _ := strings.Cut(v, "=")
		name, t := filter(key)
		switch t {
		case VarTypeUnknown:
			continue
		case VarTypeString:
			values[name] = value
		case VarTypeArray:
			l := strings.LastIndex(name, "_")
			number, _ := strconv.Atoi(name[l+1:])
			name := name[:l]
			x, ok := values[name]
			if !ok {
				y := make([]string, 0)
				x = y
			}
			z, ok := x.([]string)
			if !ok {
				z = []string{fmt.Sprint(z)} // coerce string
			}
			if len(z)-1 < number {
				w := make([]string, number+1)
				copy(w, z)
				z = w
			}
			z[number] = value
			values[name] = z
		}
	}
	return values
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
