package envflagset

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

var (
	cn          string
	prefix      string
	ver         string
	fs          *flag.FlagSet
	showVersion bool
	dumpDefault bool
	parsed      bool
)

const (
	fVersion = "Ver"
	fDumpEnv = "DumpEnv"
)

func New(name, version string) *flag.FlagSet {
	cn = name
	SetPrefix(name)
	ver = version
	// fs = flag.NewFlagSet(name, flag.ExitOnError)
	fs = flag.CommandLine
	fs.BoolVar(&showVersion, fVersion, false, "Print the version and exit")
	fs.BoolVar(&dumpDefault, fDumpEnv, false, "output all default env values")

	return fs
}

func SetPrefix(name string) {
	if name != "" {
		prefix = toEnvName(name) + "_"
	}
}

func toEnvName(name string) string {
	return strings.ToUpper(strings.Replace(name, "-", "_", -1))
}

func Parse() {
	if parsed {
		return
	}
	parsed = true

	if len(os.Args) > 1 && strings.HasSuffix(os.Args[1], fDumpEnv) {
		Dump(fs, prefix)
		os.Exit(0)
	}

	perr := fs.Parse(os.Args[1:])
	switch perr {
	case nil:
	case flag.ErrHelp:
		os.Exit(0)
	default:
		os.Exit(2)
	}

	if showVersion {
		fmt.Printf("%s version %s %s\n", cn, ver, runtime.Version())
		os.Exit(0)
	}

	err := ParseEnv(fs, prefix)
	if err != nil {
		log.Fatalf("ParseEnv error: %v", err)
	}

}

// ParseEnv parses all registered flags in the given flagset,
// and if they are not already set it attempts to set their values from
// environment variables. Environment variables take the name of the flag but
// are UPPERCASE, have the prefix "PREFIX_", and any dashes are replaced by
// underscores - for example: some-flag => PREFIX_SOME_FLAG
func ParseEnv(fs *flag.FlagSet, prefix string) error {
	var err error
	alreadySet := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		alreadySet[f.Name] = true
	})

	fs.VisitAll(func(f *flag.Flag) {
		if !alreadySet[f.Name] {
			key := prefix + toEnvName(f.Name)
			val := os.Getenv(key)
			if val != "" {
				if serr := fs.Set(f.Name, val); serr != nil {
					err = fmt.Errorf("invalid value %q for %s: %v", val, key, serr)
				}
			}
		}
	})
	return err
}

func Dump(fs *flag.FlagSet, prefix string) {
	fs.VisitAll(func(f *flag.Flag) {
		if len(f.Name) < 2 || strings.HasPrefix(f.Name, "-") {
			return
		}
		key := prefix + toEnvName(f.Name)
		fmt.Printf("%s=%q\n", key, f.DefValue)
	})
}
