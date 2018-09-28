package g

import "path/filepath"

var Modules map[string]bool
var BinOf map[string]string
var cfgOf map[string]string
var ModuleApps map[string]string
var logpathOf map[string]string
var PidOf map[string]string
var AllModulesInOrder []string

func init() {
	Modules = map[string]bool{
		"app": true,
	}

	BinOf = map[string]string{
		"app": "./tm_app",
	}

	cfgOf = map[string]string{
		"app": "./config.json",
	}

	ModuleApps = map[string]string{
		"app": "tm_app",
	}

	logpathOf = map[string]string{
		"app": "./log/stdout.log",
	}

	PidOf = map[string]string{
		"app": "<NOT SET>",
	}

	// Modules are deployed in this order
	AllModulesInOrder = []string{
		"app",
	}
}

func Bin(name string) string {
	p, _ := filepath.Abs(BinOf[name])
	return p
}

func Cfg(name string) string {
	p, _ := filepath.Abs(cfgOf[name])
	return p
}

func LogPath(name string) string {
	p, _ := filepath.Abs(logpathOf[name])
	return p
}

func LogDir(name string) string {
	d, _ := filepath.Abs(filepath.Dir(logpathOf[name]))
	return d
}
