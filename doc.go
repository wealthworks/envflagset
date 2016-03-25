package envflagset

/*

import (
	"flag"
	envcfg "lcgc/platform/envflagset"
)

var (
	fs           *flag.FlagSet
	HttpListen   string
	CoreDSN      string
)

func init() {
	fs = envcfg.New("app", "0.0.1")
	fs.StringVar(&HttpListen, "http-listen", "localhost:5000", "bind address and port")
	fs.StringVar(&CoreDSN, "core-dsn", "mysql://user:pass@localhost:3306/appdb", "core database connecting string")
}

func main() {
	envcfg.Parse()
	// start app
}

// env:
// export APP_HTTP_LISTEN="localhost:5002"
// ./app

*/
