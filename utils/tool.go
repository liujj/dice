package utils

import (
	"crypto/sha1"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//当前项目根目录
var API_ROOT string

// 打印当前时间
func PrintTime(s string) {

	t := time.Now()
	fmt.Printf(s)
	fmt.Println(t)
}

// 打印出接口的时间类型
func PrintType(j interface{}) {

	fmt.Println(reflect.TypeOf(j))
}

//------------------------------------------------类型互转--------------------

func IntTurnString(i int) string {
	return strconv.Itoa(i)
}

//------------------------------------------------以下都是从接口类型进行类型断言转换---------------------

// 从接口类型转换到[]byte
func TurnByte(i interface{}) []byte {

	j, p := i.([]byte)
	if p {
		return j
	}

	return nil
}

// 从接口类型转换到map[string]interface{}
func TurnMapStringInterface(i interface{}) map[string]interface{} {

	j, p := i.(map[string]interface{})
	if p {
		return j
	}

	return nil
}

// 从接口类型转换到String
func TurnString(i interface{}) string {

	j, p := i.(string)
	if p {
		return j
	}

	return ""
}

// 从接口类型转换到Int
func TurnInt(i interface{}) int {

	j, p := i.(int)
	if p {
		return j
	}

	return 0
}

// 从接口类型转换到Int64
func TurnInt64(i interface{}) int64 {

	j, p := i.(int64)
	if p {
		return j
	}

	return 0
}

// 从接口类型转换到int64返回int类型
func Int64TurnInt(i interface{}) int {

	j, p := i.(int64)
	if p {
		return int(j)
	}

	return 0
}

// 从接口类型转换到Float64
func TurnFloat64(i interface{}) float64 {

	j, p := i.(float64)
	if p {
		return j
	}

	return 0
}

// 从接口类型转换到接口切片
func TurnSlice(i interface{}) []interface{} {

	j, p := i.([]interface{})
	if p {
		return j
	}

	return nil
}

//---------------------urlcode

// URL编码
func UrlEncode(urls string) (string, error) {

	//UrlEnCode编码
	urlStr, err := url.Parse(urls)
	if err != nil {
		return "", err
	}

	return urlStr.RequestURI(), nil
}

// URL解码
func UrlDecode(urls string) (string, error) {

	//UrlEnCode解码
	urlStr, err := url.Parse(urls)
	if err != nil {
		return "", err
	}

	return urlStr.Path, nil
}

// 获取项目路径
func GetPath() string {

	if API_ROOT != "" {
		return API_ROOT
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		print(err.Error())
	}

	API_ROOT = strings.Replace(dir, "\\", "/", -1)
	return API_ROOT
}

//判断文件目录否存在
func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

}

//创建文件
func MkdirFile(path string) error {

	err := os.Mkdir(path, os.ModePerm) //在当前目录下生成md目录
	if err != nil {
		return err
	}
	return nil
}

func GetRunPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	return path[:index+1]
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GenerateSecret() string {
	var b [16]int
	for i := 0; i < 16; i++ {
		b[i] = rand.Intn(256)
	}
	sk := ""
	for _, v := range b {
		sk += fmt.Sprintf("%02x", v) //默认添加0 显示2位 小写16进制
	}
	return sk
}

func GenerateRamdom(len int) string {
	t := time.Now().UnixNano()
	b := []byte(fmt.Sprintf("%d", t))
	s := sha1.New()
	s.Write(b)
	k := fmt.Sprintf("%x", s.Sum(nil))
	return k[0:len]
}
func Ip2Long(ipstr string) (ip uint32) {
	if ipstr == "::1" {
		ipstr = "127.0.0.1"
	}
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return
	}
	ips := reg.FindStringSubmatch(ipstr)
	if ips == nil {
		return
	}

	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	return
}

func IpRegion2Long(ipstr string) (ip uint32, sub uint8) {
	r := `^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})/(\d{1,2})`
	reg, err := regexp.Compile(r)
	if err != nil {
		return
	}
	ips := reg.FindStringSubmatch(ipstr)
	if ips == nil {
		return
	}

	ip1, _ := strconv.Atoi(ips[1])
	ip2, _ := strconv.Atoi(ips[2])
	ip3, _ := strconv.Atoi(ips[3])
	ip4, _ := strconv.Atoi(ips[4])
	subnet, _ := strconv.Atoi(ips[5])

	if ip1 > 255 || ip2 > 255 || ip3 > 255 || ip4 > 255 || sub > 32 {
		return
	}

	ip += uint32(ip1 * 0x1000000)
	ip += uint32(ip2 * 0x10000)
	ip += uint32(ip3 * 0x100)
	ip += uint32(ip4)

	sub = uint8(subnet)

	return
}

func Long2IpRegion(ip uint32, sub uint8) string {
	return fmt.Sprintf("%d.%d.%d.%d/%d", ip>>24, ip<<8>>24, ip<<16>>24, ip<<24>>24, sub)
}

func Long2Ip(ip uint32) string {
	return fmt.Sprintf("%d.%d.%d.%d", ip>>24, ip<<8>>24, ip<<16>>24, ip<<24>>24)
}

func TimeSincePro(then time.Time) string {
	now := time.Now()
	diff := now.Unix() - then.Unix()

	if then.After(now) {
		return "future"
	}

	var timeStr, diffStr string
	for {
		if diff == 0 {
			break
		}

		diff, diffStr = computeTimeDiff(diff)
		timeStr += ", " + diffStr
	}
	return strings.TrimPrefix(timeStr, ", ")
}

const (
	Minute = 60
	Hour   = 60 * Minute
	Day    = 24 * Hour
	Week   = 7 * Day
	Month  = 30 * Day
	Year   = 12 * Month
)

func computeTimeDiff(diff int64) (int64, string) {
	diffStr := ""
	switch {
	case diff <= 0:
		diff = 0
		diffStr = "now"
	case diff < 2:
		diff = 0
		diffStr = "1 second"
	case diff < 1*Minute:
		diffStr = fmt.Sprintf("%d seconds", diff)
		diff = 0

	case diff < 2*Minute:
		diff -= 1 * Minute
		diffStr = "1 minute"
	case diff < 1*Hour:
		diffStr = fmt.Sprintf("%d minutes", diff/Minute)
		diff -= diff / Minute * Minute

	case diff < 2*Hour:
		diff -= 1 * Hour
		diffStr = "1 hour"
	case diff < 1*Day:
		diffStr = fmt.Sprintf("%d hours", diff/Hour)
		diff -= diff / Hour * Hour

	case diff < 2*Day:
		diff -= 1 * Day
		diffStr = "1 day"
	case diff < 1*Week:
		diffStr = fmt.Sprintf("%d days", diff/Day)
		diff -= diff / Day * Day

	case diff < 2*Week:
		diff -= 1 * Week
		diffStr = "1 week"
	case diff < 1*Month:
		diffStr = fmt.Sprintf("%d weeks", diff/Week)
		diff -= diff / Week * Week

	case diff < 2*Month:
		diff -= 1 * Month
		diffStr = "1 month"
	case diff < 1*Year:
		diffStr = fmt.Sprintf("%d months", diff/Month)
		diff -= diff / Month * Month

	case diff < 2*Year:
		diff -= 1 * Year
		diffStr = "1 year"
	default:
		diffStr = fmt.Sprintf("%d years", diff/Year)
		diff = 0
	}
	return diff, diffStr
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(logn(float64(s), base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", val, suffix)
}

// FileSize calculates the file size and generate user-friendly string.
func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}

func TimeStepby10Min(t time.Time) int {
	return t.Hour()*6 + t.Minute()/10
}

func EthEip55(ethAddr string) string {
	//b := []byte(ethAddr)
	if strings.HasPrefix(ethAddr, "0x") {
		ethAddr = ethAddr[2:]
	}
	unchecksummed := strings.ToLower(ethAddr)
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)
	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}
