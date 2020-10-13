package envy

import (
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
