package args

import "flag"

var (
	ConfigPath  = flag.String("config", "", "Path to configuration file")
	HistoryFile = flag.String("history", "", "Name of history file to preload (from config dir)")
	Quiet       = flag.Bool("quiet", false, "Suppress info messages")
	Model       = flag.String("model", "", "Set model for the request")
	ShowVersion = flag.Bool("version", false, "Print picochat version")
	Image       = flag.String("image", "", "Path to an image file")
	Output      = flag.String("output", "", "Output format of the LLM (plain, json, json-pretty, yaml)")
	Format      = flag.String("format", "", "Path to a JSON schema for the LLM content")
)

func Parse() {
	flag.Parse()
}
