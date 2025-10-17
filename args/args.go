package args

import "flag"

var (
	ConfigPath  = flag.String("config", "", "Path to configuration file")
	HistoryFile = flag.String("history", "", "Name of history file to preload (from config dir)")
	Quiet       = flag.Bool("quiet", false, "Suppress initial output")
	ShowVersion = flag.Bool("version", false, "Print picochat version")
)

func Parse() {
	flag.Parse()
}
