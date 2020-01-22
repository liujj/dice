package handler

import (
	"dice/pkg/global"
	"dice/pkg/model"
	"dice/pkg/utils"
	"dice/pkg/version"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type versionInfo struct {
	GitCommit string
	BuildTime string
}
type basicInfo struct {
	Pid          int    // pid
	ServerTime   string // 当前服务器时间
	LaunchTime   string // 程序启动时间
	UpTime       string // 程序已启动时间
	NumGoroutine int    // go routine 数量
}
type memInfo struct {
	MemAllocated string // bytes allocated and still in use
	MemTotal     string // bytes allocated (even if freed)
	MemSys       string // bytes obtained from system (sum of XxxSys below)
	PointLookups uint64 // number of pointer lookups
	MemMallocs   uint64 // number of mallocs
	MemFrees     uint64 // number of frees
}
type gcInfo struct {
	NextGC       string // next run in HeapAlloc time (bytes)
	LastGC       string // last run in absolute time (ns)
	PauseTotalNs string
	PauseNs      string // circular buffer of recent GC pause times, most recent at [(NumGC+255)%256]
	NumGC        uint32
}
type heapInfo struct {
	HeapAlloc    string // bytes allocated and still in use
	HeapSys      string // bytes obtained from system
	HeapIdle     string // bytes in idle spans
	HeapInuse    string // bytes in non-idle span
	HeapReleased string // bytes released to the OS
	HeapObjects  uint64 // total number of allocated objects
}
type stackInfo struct {
	StackInuse  string // bootstrap stacks
	StackSys    string
	MSpanInuse  string // mspan structures
	MSpanSys    string
	MCacheInuse string // mcache structures
	MCacheSys   string
	BuckHashSys string // profiling bucket hash table
	GCSys       string // GC metadata
	OtherSys    string // other system allocations
}
type dbInfo struct {
	NumDbConn int //数据库连接数量
}

func AppVersion(c echo.Context) (err error) {
	return c.JSON(http.StatusOK, versionInfo{
		GitCommit: version.GitCommit,
		BuildTime: version.BuildTime,
	})
}

func ServerMonitor(c echo.Context) (err error) {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return c.JSON(http.StatusOK, struct {
		Version versionInfo
		Basic   basicInfo
		GC      gcInfo
		Mem     memInfo
		Heap    heapInfo
		Stack   stackInfo
		DB      dbInfo
	}{
		Version: versionInfo{
			GitCommit: version.GitCommit,
			BuildTime: version.BuildTime,
		},
		Basic: basicInfo{
			Pid:          os.Getpid(),
			ServerTime:   time.Now().Local().String(),
			LaunchTime:   model.Stat.Uptime.String(),
			UpTime:       utils.TimeSincePro(model.Stat.Uptime),
			NumGoroutine: runtime.NumGoroutine(),
		},
		GC: gcInfo{
			NextGC:       utils.FileSize(int64(m.NextGC)),
			LastGC:       fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000),
			PauseTotalNs: fmt.Sprintf("%.1fs", float64(m.PauseTotalNs)/1000/1000/1000),
			PauseNs:      fmt.Sprintf("%.3fs", float64(m.PauseNs[(m.NumGC+255)%256])/1000/1000/1000),
			NumGC:        m.NumGC,
		},
		Mem: memInfo{
			MemAllocated: utils.FileSize(int64(m.Alloc)),
			MemTotal:     utils.FileSize(int64(m.TotalAlloc)),
			MemSys:       utils.FileSize(int64(m.Sys)),
			PointLookups: m.Lookups,
			MemMallocs:   m.Mallocs,
			MemFrees:     m.Frees,
		},
		Heap: heapInfo{
			HeapAlloc:    utils.FileSize(int64(m.HeapAlloc)),
			HeapSys:      utils.FileSize(int64(m.HeapSys)),
			HeapIdle:     utils.FileSize(int64(m.HeapIdle)),
			HeapInuse:    utils.FileSize(int64(m.HeapInuse)),
			HeapReleased: utils.FileSize(int64(m.HeapReleased)),
			HeapObjects:  m.HeapObjects,
		},
		Stack: stackInfo{
			StackInuse:  utils.FileSize(int64(m.StackInuse)),
			StackSys:    utils.FileSize(int64(m.StackSys)),
			MSpanInuse:  utils.FileSize(int64(m.MSpanInuse)),
			MSpanSys:    utils.FileSize(int64(m.MSpanSys)),
			MCacheInuse: utils.FileSize(int64(m.MCacheInuse)),
			MCacheSys:   utils.FileSize(int64(m.MCacheSys)),
			BuckHashSys: utils.FileSize(int64(m.BuckHashSys)),
			GCSys:       utils.FileSize(int64(m.GCSys)),
			OtherSys:    utils.FileSize(int64(m.OtherSys)),
		},
		DB: dbInfo{
			NumDbConn: model.Gdb.DB().Stats().OpenConnections,
		},
	})
}

func AdminAuth(c echo.Context) error {
	if err := auth(c); err != nil {
		return err
	}
	sess, _ := session.Get("admin_session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}
	sess.Values["admin"] = "1"
	sess.Save(c.Request(), c.Response())
	return SimpleResponse(c, http.StatusOK, 0, "auth ok")
}

func Reload(c echo.Context) error {
	cmd := exec.Command("kill", "-USR2", strconv.Itoa(os.Getpid()))
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	if err == nil {
		return SimpleResponse(c, http.StatusOK, 0, "ok")
	}
	return DataResponse(c, http.StatusOK, -1, "error", err)
}

func ReadLog(c echo.Context) error {
	fileStr := c.QueryParam("file")
	file, ferr := os.OpenFile(global.LogDir+fileStr, os.O_RDONLY, 0644)
	defer file.Close()
	if ferr != nil {
		return DataResponse(c, http.StatusOK, -1, "file open error", ferr)
	}
	stat, statErr := file.Stat()
	if statErr != nil {
		return DataResponse(c, http.StatusOK, -1, "file stat error", statErr)
	}
	fr, frErr := strconv.Atoi(c.QueryParam("fr"))
	to, toErr := strconv.Atoi(c.QueryParam("to"))
	tail, tailErr := strconv.Atoi(c.QueryParam("tail"))
	if frErr == nil && toErr == nil {
		if int64(to) > stat.Size() || fr < 0 {
			return SimpleResponse(c, 200, -1, fmt.Sprintf("filesize is %d", stat.Size()))
		}
		if to-fr > 1048576 || to-fr <= 0 {
			return SimpleResponse(c, 200, -1, "tail should less than 1MB")
		}
		res := make([]byte, to-fr)
		file.Seek(int64(fr), 0)
		file.Read(res)
		c.String(http.StatusOK, string(res))
	} else if tailErr == nil {
		if int64(tail) > stat.Size() || tail <= 0 {
			return SimpleResponse(c, 200, -1, fmt.Sprintf("filesize is %d", stat.Size()))
		}
		if tail > 1048576 {
			return SimpleResponse(c, 200, -1, "tail should less than 1MB")
		}
		res := make([]byte, tail)
		file.Seek(-1*int64(tail), 2)
		file.Read(res)
		return c.String(http.StatusOK, string(res))
	} else {
		return DataResponse(c, http.StatusOK, 0, "stat", struct {
			Name    string
			Size    int64
			Mode    string
			ModTime time.Time
		}{
			Name:    stat.Name(),
			Size:    stat.Size(),
			Mode:    stat.Mode().String(),
			ModTime: stat.ModTime(),
		})
	}
	return nil
}

func SqlExec(c echo.Context) error {
	sql, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return SimpleResponse(c, http.StatusBadRequest, -1, err.Error())
	}
	if len(sql) == 0 {
		return SimpleResponse(c, http.StatusBadRequest, -1, "body length is 0")
	}
	dbResult := model.Gdb.Exec(string(sql))
	if dbResult.Error != nil {
		return SimpleResponse(c, http.StatusBadRequest, -1, dbResult.Error.Error())
	}
	return SimpleResponse(c, http.StatusBadRequest, 0, fmt.Sprintf("affected: %d", dbResult.RowsAffected))
}
