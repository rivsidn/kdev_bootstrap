package config

var UbuntuSuiteMap = map[string]string{
	"5.10":  "breezy",
	"16.04": "xenial",
	"18.04": "bionic",
	"20.04": "focal",
	"22.04": "jammy",
	"24.04": "noble",
}

func (c *Config) GetSuite() string,err {
	if suite, ok := UbuntuSuiteMap[c.Version]; ok {
		return suite
	} else {
		return ""
	}
}

