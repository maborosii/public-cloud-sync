package setting

type AccountName string

type Config struct {
	Accounts   map[AccountName]AccountInfo `toml:"accounts" mapstructure:"accounts"`
	MailConfig *MailConf                   `toml:"mailconfig"`
}

type AccountInfo struct {
	AK       string `toml:"ak"`
	SK       string `toml:"sk"`
	Region   string `toml:"region"`
	Project  string `toml:"project"`
	Provider string `toml:"provider"`
}

type MailConf struct {
	Host     string   `toml:"host"`
	UserName string   `toml:"username"`
	Password string   `toml:"password"`
	Port     int      `toml:"port"`
	Subject  string   `toml:"subject"`
	To       []string `toml:"to"`
}
