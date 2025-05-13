// parameters/flags.go
package args

import "flag"

var (
	ConfigPath  = flag.String("config", "", "Path to configuration file")
	HistoryFile = flag.String("history", "", "Name of history file to preload (from config dir)")
)

func Parse() {
	flag.Parse()
}
