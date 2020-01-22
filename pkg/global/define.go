package global

import (
	"os"
	"path/filepath"
)

// 全局常量
const ()

// 全局变量
var (
	Debug         bool
	ConfigPath, _ = filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "config.json"))
	PidFile, _    = filepath.Abs(filepath.Join(filepath.Dir(os.Args[0]), "pid"))
	LogDir        = filepath.Join(os.TempDir(), "dice")
)
