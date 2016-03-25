package envflagset

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	cn          string
	ver         string
	fs          *flag.FlagSet
	showVersion bool
	parsed      bool
)

func New(name, version string) *flag.FlagSet {
	cn = name
	ver = version
	fs = flag.NewFlagSet(name, flag.ExitOnError)
	fs.BoolVar(&showVersion, "version", false, "Print the version and exit")

	return fs
}

func Parse() {
	if parsed {
		return
	}
	parsed = true
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
	if len(fs.Args()) != 0 {
		log.Fatalf("'%s' is not a valid flag", fs.Arg(0))
	}

	if showVersion {
		fmt.Printf("%s version %s\n", cn, ver)
		os.Exit(0)
	}

	err := ParseEnv(fs, cn+"_")
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

	if prefix == "" {
		prefix = "_"
	} else {
		prefix = strings.ToUpper(strings.Replace(prefix, "-", "_", -1))
	}

	fs.VisitAll(func(f *flag.Flag) {
		if !alreadySet[f.Name] {
			key := prefix + strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
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
