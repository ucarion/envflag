package envflag_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucarion/envflag"
)

// func withArgv(fn func()) {
// 	defer func(args []string) {
// 		os.Args = args
// 	}(os.Args)

// 	fn()
// }

// func TestParse(t *testing.T) {
// 	initialArgv := os.Args
// 	initialEnv := os.Environ()

// 	defer func() {
// 		os.Args = initialArgv

// 		os.Clearenv()
// 		for _, env := range initialEnv {
// 			parts := strings.SplitN(env, "=", 2)
// 			os.Setenv(parts[0], parts[1])
// 		}
// 	}()

// 	a := flag.String("a", "default-a", "")
// 	b := flag.Int("b", 123, "")

// 	os.Args = []string{"testparse", "--a=foo"}
// 	os.Clearenv()
// 	os.Setenv("TESTPARSE_B", "456")

// 	envflag.Parse()
// 	flag.Parse()

// 	if *a != "foo" {
// 		t.Errorf("expected a to be foo, was: %s", *a)
// 	}

// 	if *b != 456 {
// 		t.Errorf("expected b to be 456, was: %d", *b)
// 	}
// }

func TestParseFlagSet(t *testing.T) {
	// You can't un-set flags. So we set up the flags.CommandLine flags here,
	// outside of the individual tests.
	a := flag.String("a", "default-a", "")
	b := flag.Int("b", 123, "")
	c := flag.Bool("has-dashes", false, "")

	type testCase struct {
		name string
		args []string
		env  map[string]string
		fn   func(t *testing.T)
	}

	testCases := []testCase{
		testCase{
			name: "flags.CommandLine, all defaults",
			args: []string{"cmd"},
			env:  map[string]string{},
			fn: func(t *testing.T) {
				envflag.Parse()
				flag.Parse()

				assert.Equal(t, "default-a", *a)
				assert.Equal(t, 123, *b)
				assert.Equal(t, false, *c)
			},
		},

		testCase{
			name: "flags.CommandLine, override partially via env",
			args: []string{"cmd"},
			env: map[string]string{
				"CMD_A":          "from-env",
				"CMD_HAS_DASHES": "true",
			},
			fn: func(t *testing.T) {
				envflag.Parse()
				flag.Parse()

				assert.Equal(t, "from-env", *a)
				assert.Equal(t, 123, *b)
				assert.Equal(t, true, *c)
			},
		},

		testCase{
			name: "flags.CommandLine, override partially via env, override that with argv",
			args: []string{"cmd", "--a=from-argv"},
			env: map[string]string{
				"CMD_A":          "from-env",
				"CMD_HAS_DASHES": "true",
			},
			fn: func(t *testing.T) {
				envflag.Parse()
				flag.Parse()

				assert.Equal(t, "from-argv", *a)
				assert.Equal(t, 123, *b)
				assert.Equal(t, true, *c)
			},
		},

		testCase{
			name: "custom flagset, all defaults",
			args: []string{},
			env:  map[string]string{},
			fn: func(t *testing.T) {
				fs := flag.NewFlagSet("", flag.ExitOnError)
				a := fs.String("a", "default", "")
				b := fs.Int("b", 123, "")
				c := fs.Bool("has-dashes", false, "")

				assert.NoError(t, envflag.ParseFlagSet("", fs))
				assert.NoError(t, fs.Parse([]string{}))

				assert.Equal(t, "default", *a)
				assert.Equal(t, 123, *b)
				assert.Equal(t, false, *c)
			},
		},

		testCase{
			name: "custom flagset, override partially via env",
			args: []string{},
			env: map[string]string{
				"A":          "from-env",
				"HAS_DASHES": "true",
			},
			fn: func(t *testing.T) {
				fs := flag.NewFlagSet("", flag.ExitOnError)
				a := fs.String("a", "default", "")
				b := fs.Int("b", 123, "")
				c := fs.Bool("has-dashes", false, "")

				assert.NoError(t, envflag.ParseFlagSet("", fs))
				assert.NoError(t, fs.Parse([]string{}))

				assert.Equal(t, "from-env", *a)
				assert.Equal(t, 123, *b)
				assert.Equal(t, true, *c)
			},
		},

		testCase{
			name: "custom flagset, override partially via env, override that with argv",
			args: []string{},
			env: map[string]string{
				"A":          "from-env",
				"HAS_DASHES": "true",
			},
			fn: func(t *testing.T) {
				fs := flag.NewFlagSet("", flag.ExitOnError)
				a := fs.String("a", "default", "")
				b := fs.Int("b", 123, "")
				c := fs.Bool("has-dashes", false, "")

				assert.NoError(t, envflag.ParseFlagSet("", fs))
				assert.NoError(t, fs.Parse([]string{"--a=from-argv"}))

				assert.Equal(t, "from-argv", *a)
				assert.Equal(t, 123, *b)
				assert.Equal(t, true, *c)
			},
		},

		testCase{
			name: "custom flagset with prefix, override partially via env, override that with argv",
			args: []string{},
			env: map[string]string{
				"SOME_PREFIX_A":          "from-env",
				"SOME_PREFIX_HAS_DASHES": "true",
			},
			fn: func(t *testing.T) {
				fs := flag.NewFlagSet("", flag.ExitOnError)
				a := fs.String("a", "default", "")
				b := fs.Int("b", 123, "")
				c := fs.Bool("has-dashes", false, "")

				assert.NoError(t, envflag.ParseFlagSet("some-prefix", fs))
				assert.NoError(t, fs.Parse([]string{"--a=from-xargv"}))

				assert.Equal(t, "from-argv", *a)
				assert.Equal(t, 123, *b)
				assert.Equal(t, true, *c)
			},
		},

		testCase{
			name: "custom flagset, panic on error",
			args: []string{},
			env: map[string]string{
				"B": "not-an-int",
			},
			fn: func(t *testing.T) {
				fs := flag.NewFlagSet("", flag.PanicOnError)
				fs.Int("b", 123, "")

				assert.Panics(t, func() {
					envflag.ParseFlagSet("", fs)
				})
			},
		},

		testCase{
			name: "custom flagset, continue on error",
			args: []string{},
			env: map[string]string{
				"B": "not-an-int",
			},
			fn: func(t *testing.T) {
				fs := flag.NewFlagSet("", flag.ContinueOnError)
				fs.Int("b", 123, "")

				assert.Equal(t, "parse error", envflag.ParseFlagSet("", fs).Error())
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			defer func(args []string) {
				os.Args = args
			}(os.Args)

			defer func(env []string) {
				os.Clearenv()
				for _, v := range env {
					parts := strings.SplitN(v, "=", 2)
					assert.NoError(t, os.Setenv(parts[0], parts[1]))
				}
			}(os.Environ())

			os.Args = tt.args

			os.Clearenv()
			for k, v := range tt.env {
				assert.NoError(t, os.Setenv(k, v))
			}

			tt.fn(t)
		})
	}
}
