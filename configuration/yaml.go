package configuration

type yamlDataCore struct {
	TimeZone string `yaml:"time_zone"`
	LogLevel int    `yaml:"log_level"`
	LogDir   string `yaml:"log_directory"`
}

type yamlDataGame struct {
	GalaxyName string `yaml:"galaxy_name"`
	InstallDir string `yaml:"install_dir"`
	DataDir    string `yaml:"data_dir"`
	GamePort   int    `yaml:"port"`
	PingPort   int    `yaml:"ping_port"`
}

type yamlDataDiscord struct {
	ChatChannel      string              `yaml:"channel"`
	BotsAllowed      bool                `yaml:"bots_allowed"`
	DiscordLink      string              `yaml:"invite"`
	Prefix           string              `yaml:"prefix"`
	Token            string              `yaml:"token"`
	AliasedCommands  map[string][]string `yaml:"aliased_commands"`
	DisabledCommands []string            `yaml:"disabled_commands,flow"`
}

type yamlDataRCON struct {
	Address string `yaml:"address"`
	Binary  string `yaml:"binary"`
}

type yamlData struct {
	Core    yamlDataCore    `yaml:"Core"`
	Game    yamlDataGame    `yaml:"Game"`
	RCON    yamlDataRCON    `yaml:"RCON"`
	Discord yamlDataDiscord `yaml:"Discord"`
}
