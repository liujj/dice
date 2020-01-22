package global

import (
	"dice/pkg/version"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
)

func ParseFlag() {
	flag.BoolVarP(&Debug, "debug", "D", false, "print debug level log")
	flag.StringVar(&PidFile, "pid", PidFile, "rewrite pidfile location")
	flag.StringVarP(&ConfigPath, "config","c", ConfigPath, "config file location")
	flag.StringVar(&LogDir, "logdir", LogDir, "rewrite log dir location")
	showVersion := flag.BoolP("version", "v", false, "show version")
	neepHelp := flag.BoolP("help", "h", false, "show this help")
	flag.Parse()

	if *showVersion {
		version.ShowVersion()
		os.Exit(0)
	}
	if *neepHelp {
		flag.Usage()
		fmt.Println("\tstart\t\tstart daemon")
		fmt.Println("\tstop\t\tstop daemon")
		fmt.Println("\trestart\t\trestart daemon")
		fmt.Println("\tstatus\t\tshow daemon status")
		fmt.Println("\tremove\t\tremove all binary, scripts; reserve logs")
		os.Exit(0)
	}

	if Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}
