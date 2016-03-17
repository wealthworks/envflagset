package envflagset

/*

import (
	"flag"
	envcfg "tuluu.com/platform/envflagset"
)

type config struct {
	HttpListen   string
}

var (
	fs       *flag.FlagSet
	Settings *config = &config{}
)

func init() {
	fs = envcfg.New("app", "0.0.1")
	fs.StringVar(&Settings.HttpListen, "http-listen", "localhost:5000", "bind address and port")
}

func main() {
	envcfg.Parse()
	// start app
}

// env:
// export APP_HTTP_LISTEN="localhost:5002"
// ./app

*/
