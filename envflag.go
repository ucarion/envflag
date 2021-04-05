// Package envflag enhances the flag package with the ability to read from
// environment variables.
package envflag

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Parse loads any environemnt variables it can find into flag.CommandLine.
// Environment variables are expected to be prefixed with the filename of
// os.Args[0].
//
// flag.CommandLine is the "default" / "global" FlagSet. The top-level
// flag.String, flag.Int, etc. functions ultimately get added to
// flag.CommandLine.
//
// Parse calls Load under the hood. If you want to parse into another FlagSet
// than flag.CommandLine, or if you would like to customize or remove the
// os.Args[0] prefix, then consider using Load instead.
//
// Parse will call flag.Parse(). Though there are no negative consequences to
// calling flag.Parse() after calling flagenv.Parse(), there are no benefits
// either. If you don't want this package to call flag.Parse(), then use Load
// instead.
func Parse() {
	Load(filepath.Base(os.Args[0]), flag.CommandLine)
	flag.Parse()
}

// Load loads environment variables into a flag.FlagSet.
//
// Environment variables ("env vars") are expected to be named after their
// corresponding flag's name in upper-case letters, with dashes converted to
// underscores. If prefix is non-empty, then the env var must be prefixed by
// prefix (in all caps, with dashes converted to underscores) and an underscore.
//
// For example, if prefix is empty, then for a flag named "user-id", Load will
// look for an env var named "USER_ID". If prefix were instead "count-users",
// then Load would instead look for an env var named "COUNT_USERS_USER_ID".
//
// If an env var for a flag is not found, then that flag is untouched. Whatever
// value it had before calling Load is preserved.
//
// If an env var for a flag is found, but its value is incompatible with the
// flag (for example, if an Int flag has a corresponding env var whose value
// isn't parsable as an int), then Load will trigger an error in correspondence
// with the ErrorHandling of the given FlagSet.
func Load(prefix string, fs *flag.FlagSet) error {
	var err error

	fs.VisitAll(func(f *flag.Flag) {
		if err == nil {
			if env, ok := os.LookupEnv(flagNameToEnvKey(prefix, f.Name)); ok {
				err = f.Value.Set(env)
			}
		}
	})

	if err != nil {
		switch fs.ErrorHandling() {
		case flag.ContinueOnError:
			return err
		case flag.ExitOnError:
			os.Exit(2)
		case flag.PanicOnError:
			panic(err)
		}
	}

	return nil
}

func flagNameToEnvKey(prefix, name string) string {
	base := name
	if prefix != "" {
		base = fmt.Sprintf("%s-%s", prefix, name)
	}

	return strings.ToUpper(strings.ReplaceAll(base, "-", "_"))
}
