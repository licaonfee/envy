package envy_test

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/licaonfee/envy/v2"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name string
		env  envy.Env
		key  string
		want string
	}{
		{
			name: "Get existent env",
			env:  envy.NewMapEnv(map[string]string{"existent_key": "value"}),
			key:  "existent_key",
			want: "value",
		},
		{
			name: "Get not existent env",
			env:  envy.NewMapEnv(map[string]string{"existent_key": "value"}),
			key:  "missing_key",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.env.Getenv(tt.key)
			if got != tt.want {
				t.Errorf("GetEnv() got = %s , want = %s", got, tt.want)
			}
		})
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
	tests := []struct {
		name  string
		env   envy.Env
		key   string
		value string
	}{
		{
			name:  "Set key",
			env:   envy.NewMapEnv(map[string]string{}),
			key:   "setted_key",
			value: "xxxxx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.env.Setenv(tt.key, tt.value)
			got := tt.env.Getenv(tt.key)
			if got != tt.value {
				t.Errorf("GetEnv() got = %s , want = %s", got, tt.value)
			}
		})
	}
}

func TestUnsetEnv(t *testing.T) {
	tests := []struct {
		name  string
		env   envy.Env
		key   string
		value string
	}{
		{
			name:  "unset key",
			env:   envy.NewMapEnv(map[string]string{"deleted_key": "asdf"}),
			key:   "deleted_key",
			value: "asdf",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = tt.env.Unsetenv(tt.key)
			got := tt.env.Getenv(tt.key)
			if got != "" {
				t.Errorf("Unsetenv() got = %s , want = %s", got, tt.value)
			}
		})
	}
}

func TestCleanenv(t *testing.T) {
	tests := []struct {
		name string
		env  envy.Env
	}{
		{
			name: "clean",
			env:  envy.NewMapEnv(map[string]string{"cleaned_key01": "asdf", "cleaned_key02": "asdf"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.env.Clearenv()
			got := tt.env.Getenv("cleaned_key01")
			got += tt.env.Getenv("cleaned_key02")
			if got != "" {
				t.Errorf("Cleanenv() got = %s , want = \"\"", got)
			}
		})
	}

}

func TestMapEnvEnviron(t *testing.T) {
	tests := []struct {
		name string
		env  envy.Env
		want []string
	}{
		{
			name: "Environ",
			env:  envy.NewMapEnv(map[string]string{"env_key_01": "01", "env_key_02": "02"}),
			want: []string{"env_key_01=01", "env_key_02=02"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.env.Environ()
			sort.Strings(got)
			sort.Strings(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Environ() got = %v , want = %v", got, tt.want)
			}
		})
	}

}

func TestFillFlagsLookup(t *testing.T) {
	tests := []struct {
		name    string
		env     envy.Env
		lookup  func(string) string
		want    map[string]interface{}
		wantErr bool
	}{
		{name: "fill all",
			env: envy.NewMapEnv(map[string]string{
				"FLAG01": "value1",
				"FLAG02": "55",
			}),
			lookup: strings.ToUpper,
			want: map[string]interface{}{
				"flag01": "value1",
				"flag02": 55,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flag.NewFlagSet(tt.name, flag.ContinueOnError)
			got := make(map[string]interface{})
			for k, v := range tt.want {
				var value interface{}
				switch v.(type) {
				case string:
					value = new(string)
					fs.StringVar(value.(*string), k, "", "")
				case int:
					value = new(int)
					fs.IntVar(value.(*int), k, 0, "")
				}
				got[k] = value
			}
			if err := fs.Parse([]string{tt.name}); err != nil {
				t.Fatal(err)
			}
			err := envy.FillFlagsLookup(fs, tt.env, tt.lookup)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FillFlagsLookup() err = %v , wantErr = %v", err, tt.wantErr)
			}
			for k, v := range tt.want {
				g, ok := got[k]
				if !ok {
					t.Fatalf("FillFlagLookup() missing flag %s", k)
				}
				var val interface{}
				switch t := g.(type) {
				case *string:
					val = *t
				case *int:
					val = *t
				}
				if !reflect.DeepEqual(val, v) {
					t.Fatalf("FillFlagLookup() got = %v , want = %v", val, v)
				}
			}
		})
	}
}

func TestFillMap(t *testing.T) {
	tests := map[string]struct {
		env    envy.Env
		filter func(string) (envy.VarName, envy.VarType)
		want   map[string]any
	}{
		"only strings": {
			env: envy.NewMapEnv(map[string]string{
				"TEST01_FOO": "var 01",
				"TEST01_BAR": "var 02",
				"TEST01_BAZ": "var 03"}),
			filter: envy.FilterPrefix("TEST01_"),
			want: map[string]any{
				"foo": "var 01",
				"bar": "var 02",
				"baz": "var 03",
			},
		},
		"array": {
			env: envy.NewMapEnv(map[string]string{
				"TEST_CONFIG":  "path",
				"TEST_USERS_0": "user00",
				"TEST_USERS_1": "user01",
				"TEST_USERS_2": "user02",
				"TEST_USERS_3": "user03",
			}),
			filter: envy.FilterPrefix("TEST_"),
			want: map[string]any{
				"config": "path",
				"users":  []string{"user00", "user01", "user02", "user03"},
			},
		},
		"map": {
			env: envy.NewMapEnv(map[string]string{
				"TEST_CONFIG":        "path",
				"TEST_USERS_0":       "user00",
				"TEST_USERS_1":       "user01",
				"TEST_ADDR_DEFAULT":  "localhost",
				"TEST_ADDR_FALLBACK": "127.0.0.1",
			}),
			filter: envy.FilterPrefix("TEST_", "addr"),
			want: map[string]any{
				"config": "path",
				"users":  []string{"user00", "user01"},
				"addr":   map[string]string{"default": "localhost", "fallback": "127.0.0.1"},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := envy.FillMap(tt.env, tt.filter)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("FillMap got %#v, want %#v", got, tt.want)
			}
		})
	}
}
