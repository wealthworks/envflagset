package envflagset

/*

import (
	"flag"
	envcfg "lcgc/platform/envflagset"
)

var (
	fs           *flag.FlagSet
	HttpListen   string
)

func init() {
	fs = envcfg.New("app", "0.0.1")
	fs.StringVar(&HttpListen, "http-listen", "localhost:5000", "bind address and port")
}

func main() {
	envcfg.Parse()
	// start app
}

// env:
// export APP_HTTP_LISTEN="localhost:5002"
// ./app

*/
