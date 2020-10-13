package envy_test

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"

	"github.com/licaonfee/envy"
)

func TestGetEnv(t *testing.T) {
	setOsEnv := func() *envy.OsEnv {
		_ = os.Setenv("existent_key", "value")
		_ = os.Unsetenv("missing_key")
		return envy.NewOsEnv()
	}
	objects := []envy.Env{
		envy.NewMapEnv(map[string]string{"existent_key": "value"}),
		setOsEnv(),
	}
	tests := []struct {
		name string
		key  string
		want string
	}{
		{
			name: "Get existent env",
			key:  "existent_key",
			want: "value",
		},
		{
			name: "Get not existent env",
			key:  "missing_key",
			want: "",
		},
	}
	for _, m := range objects {
		for _, tt := range tests {
			name := fmt.Sprintf("%s_%T", tt.name, m)
			t.Run(name, func(t *testing.T) {
				got := m.Getenv(tt.key)
				if got != tt.want {
					t.Errorf("GetEnv() got = %s , want = %s", got, tt.want)
				}
			})
		}
	}
}
func TestLookupEnv(t *testing.T) {
	setOsEnv := func() *envy.OsEnv {
		_ = os.Setenv("existent_key", "value")
		_ = os.Unsetenv("missing_key")
		return envy.NewOsEnv()
	}
	objects := []envy.Env{
		envy.NewMapEnv(map[string]string{"existent_key": "value"}),
		setOsEnv(),
	}
	tests := []struct {
		name   string
		env    envy.Env
		key    string
		want   string
		wantOk bool
	}{
		{
			name:   "Get existent env",
			key:    "existent_key",
			want:   "value",
			wantOk: true,
		},
		{
			name:   "Get not existent env",
			key:    "missing_key",
			want:   "",
			wantOk: false,
		},
	}
	for _, m := range objects {
		for _, tt := range tests {
			name := fmt.Sprintf("%s_%T", tt.name, m)
			t.Run(name, func(t *testing.T) {
				got, ok := m.LookupEnv(tt.key)
				if got != tt.want || ok != tt.wantOk {
					t.Errorf("GetEnv() got = (%s,%v) , want = (%s,%v)", got, ok, tt.want, tt.wantOk)
				}
			})
		}
	}
}

func TestSetEnv(t *testing.T) {
	setOsEnv := func() *envy.OsEnv {
		_ = os.Unsetenv("setted_key")
		return envy.NewOsEnv()
	}
	objects := []envy.Env{
		envy.NewMapEnv(map[string]string{}),
		setOsEnv(),
	}
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "Set key",
			key:   "setted_key",
			value: "xxxxx",
		},
	}
	for _, m := range objects {
		for _, tt := range tests {
			name := fmt.Sprintf("%s_%T", tt.name, m)
			t.Run(name, func(t *testing.T) {
				_ = m.Setenv(tt.key, tt.value)
				got := m.Getenv(tt.key)
				if got != tt.value {
					t.Errorf("GetEnv() got = %s , want = %s", got, tt.value)
				}
			})
		}
	}
}

func TestUnsetEnv(t *testing.T) {
	setOsenv := func() *envy.OsEnv {
		_ = os.Setenv("deleted_key", "asdf")
		return envy.NewOsEnv()
	}
	objects := []envy.Env{
		envy.NewMapEnv(map[string]string{"deleted_key": "asdf"}),
		setOsenv(),
	}
	tests := []struct {
		name  string
		key   string
		value string
	}{
		{
			name:  "unset key",
			key:   "deleted_key",
			value: "asdf",
		},
	}
	for _, m := range objects {
		for _, tt := range tests {
			name := fmt.Sprintf("%s_%T", tt.name, m)
			t.Run(name, func(t *testing.T) {
				_ = m.Unsetenv(tt.key)
				got := m.Getenv(tt.key)
				if got != "" {
					t.Errorf("Unsetenv() got = %s , want = %s", got, tt.value)
				}
			})
		}
	}
}

func TestCleanenv(t *testing.T) {
	setOsenv := func() *envy.OsEnv {
		_ = os.Setenv("cleaned_key01", "asdf")
		_ = os.Setenv("cleaned_key02", "asdf")

		return envy.NewOsEnv()
	}
	objects := []envy.Env{
		envy.NewMapEnv(map[string]string{"cleaned_key01": "asdf", "cleaned_key02": "asdf"}),
		setOsenv(),
	}
	tests := []struct {
		name string
	}{
		{
			name: "clean",
		},
	}
	for _, m := range objects {
		for _, tt := range tests {
			name := fmt.Sprintf("%s_%T", tt.name, m)
			t.Run(name, func(t *testing.T) {
				m.Clearenv()
				got := m.Getenv("cleaned_key01")
				got = got + m.Getenv("cleaned_key02")
				if got != "" {
					t.Errorf("Cleanenv() got = %s , want = \"\"", got)
				}
			})
		}
	}
}

func TestEnviron(t *testing.T) {
	setOsEnv := func() *envy.OsEnv {
		os.Clearenv()
		os.Setenv("env_key_01", "01")
		os.Setenv("env_key_02", "02")
		return envy.NewOsEnv()
	}
	objects := []envy.Env{
		envy.NewMapEnv(map[string]string{"env_key_01": "01", "env_key_02": "02"}),
		setOsEnv(),
	}
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "Environ",
			want: []string{"env_key_01=01", "env_key_02=02"},
		},
	}
	for _, m := range objects {
		for _, tt := range tests {

			name := fmt.Sprintf("%s_%T", tt.name, m)
			t.Run(name, func(t *testing.T) {
				got := m.Environ()
				sort.Strings(got)
				sort.Strings(tt.want)
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("Environ() got = %v , want = %v", got, tt.want)
				}
			})
		}
	}

}
