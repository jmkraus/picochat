package args

import "flag"

var (
	ConfigPath  = flag.String("config", "", "Loads a configuration file")
	HistoryFile = flag.String("history", "", "Loads a specific session")
	Quiet       = flag.Bool("quiet", false, "Suppresses all app messages")
	Model       = flag.String("model", "", "Overrides configured model")
	ShowVersion = flag.Bool("version", false, "Shows the version and exits")
	Image       = flag.String("image", "", "Sets a path to an image file")
	Output      = flag.String("output", "", "Sets the response output format (plain, json, json-pretty, yaml)")
	Format      = flag.String("format", "", "Sets the path to a JSON schema file")
)

func Parse() {
	flag.Parse()
}
