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

func New(name, version string) *flag.FlagSet {
	cn = name
	SetPrefix(name)
	ver = version
	fs = flag.NewFlagSet(name, flag.ExitOnError)
	fs.BoolVar(&showVersion, "version", false, "Print the version and exit")
	fs.BoolVar(&dumpDefault, "dump-env", false, "output all default env values")

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

	if len(os.Args) > 1 && os.Args[1] == "-dump-env" {
		Dump(fs, prefix)
		os.Exit(0)
	}

	// patch for flag.CommandLine (ex: golang/glog)
	flag.VisitAll(func(ff *flag.Flag) {
		fs.Var(ff.Value, ff.Name, ff.Usage)
	})

	perr := fs.Parse(os.Args[1:])
	switch perr {
	case nil:
	case flag.ErrHelp:
		os.Exit(0)
	default:
		os.Exit(2)
	}
	// if len(fs.Args()) != 0 {
	// 	log.Fatalf("'%s' is not a valid flag", fs.Arg(0))
	// }

	if showVersion {
		fmt.Printf("%s version %s %s\n", cn, ver, runtime.Version())
		os.Exit(0)
	}
	flag.Parse()

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
		if strings.HasPrefix(f.Name, "dump") || f.Name == "version" {
			return
		}
		key := prefix + toEnvName(f.Name)
		fmt.Printf("%s=%q\n", key, f.DefValue)
	})
}
