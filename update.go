package Sakura

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

func CreateInstallBatch(name string) {
	file, e := os.OpenFile("install.bat", os.O_CREATE|os.O_WRONLY, 0666)
	if e != nil {
		fmt.Println("failed")
		os.Exit(1004)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	writer.WriteString("taskkill /f /pid " + strconv.Itoa(os.Getpid()) + "\n")
	writer.WriteString("start \"" + name + "\" " + name + ".exe -a\n")
	writer.WriteString("exit\n")
	writer.Flush()
}

func ExecBatchFromWindows(path string) error {
	return exec.Command("cmd.exe", "/c", "start "+path+".bat").Start()
}

func Reload(path string) error {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		exec.Command("chmod", "777", path)
		path = "./" + path
	}
	// init接管
	cmd := exec.Command(path, "-a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

func DigestVerify(path string, ver string, digestStr string) bool {
	if runtime.GOOS == "windows" {
		path += "_" + ver + ".exe"
	}
	file, err := os.Open(path)
	if err != nil {
		Log.Error("[UP] file not exist")
		return false
	}
	md5f := md5.New()
	_, err = io.Copy(md5f, file)
	if err != nil {
		Log.Error("[UP] file open error")
		return false
	}

	ok := digestStr == hex.EncodeToString(md5f.Sum([]byte("")))
	if ok {
		Log.Info("[UP] file digest match", logrus.Fields{"digest": digestStr, "file": hex.EncodeToString(md5f.Sum([]byte("")))})
	} else {
		Log.WithFields(logrus.Fields{"digest": digestStr, "file": hex.EncodeToString(md5f.Sum([]byte("")))}).Warn("[UP] file digest mismatch")
	}
	return ok
}

func CheckAndUpdateAndReload() {
	if GetConfig(true).AutoUpdate {
		return
	}
	Log.Warn("[UP] Checking for updates")
	defer func() {
		Log.Info("[UP] update check completed")
	}()

	systemType := runtime.GOOS

	// systemType = "linux"

	// get md5 digest
	ok, digest, verStr := CheckUpdate()
	if !ok {
		return
	}
	// var digest = "fba41dcef7634ed0b6c92a22c32ea2f8"

	// download
	err := DownloadExec(Up.AppName, verStr)
	if err != nil {
		return
	}

	// verify
	verify := DigestVerify(Up.AppName, verStr, digest)
	if !verify {
		return
	}

	// reload
	Log.WithFields(logrus.Fields{"os": runtime.GOOS, "arch": runtime.GOARCH}).Info("[UP] reloading")
	if systemType == "linux" || systemType == "darwin" {
		// 执行成功回调
		if Up.SucceedCallback != nil {
			Up.SucceedCallback()
		}
		time.Sleep(time.Second * 3)
		Reload(Up.AppName)
		os.Exit(1010)
	} else if systemType == "windows" {
		CreateInstallBatch(Up.AppName + "_" + verStr)
		ExecBatchFromWindows("install")
	}
}

var host = "http://r3in.top:3000/"

// "https://cdn.jsdelivr.net/gh/r3inbowari/hbuilderx_cli@v1.0.16/meiwobuxing_darwin_amd64_v1.0.15"
var speedup = "https://cdn.jsdelivr.net/gh/r3inbowari/hbuilderx_cli@"

func DownloadExec(name, version string) error {

	goarch := runtime.GOARCH
	goos := runtime.GOOS
	// dUrl := host + name + "/bin/" + name + "_" + goos + "_" + goarch + "_" + version
	dUrl := speedup + version + "/" + name + "_" + goos + "_" + goarch + "_" + version

	if runtime.GOOS == "windows" {
		// dUrl += ".exe"
		name += "_" + version + ".exe"
	}
	// Info(dUrl)

	var bar *ProgressBar

	err := Download(dUrl, name, func(fileLength int64) {
		Log.WithFields(logrus.Fields{"size": fileLength}).Info("[UP] downloading... collected file size")
		bar = NewProgressBar(fileLength)
	}, func(length, downLen int64) {
		bar.Play(downLen)
	})
	if err != nil {
		Log.WithFields(logrus.Fields{"err": err.Error()}).Warn("[UP] download failed...")
		return err
	}
	bar.Finish()
	return nil
}

func Download(url, name string, lenCall func(fileLength int64), fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	client := new(http.Client)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		return err
	}

	if fsize < 1000000 {
		return errors.New("error update file")
	}

	lenCall(fsize)

	file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("null")
	}
	defer resp.Body.Close()
	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := file.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		fb(fsize, written)
	}
	return err
}

type Default struct {
	Name     string   `json:"name"`
	Major    int      `json:"major"`
	Minor    int      `json:"minor"`
	Patch    int      `json:"patch"`
	Types    []string `json:"types"`
	Digests  []string `json:"digests"`
	PDigests []string `json:"pDigests"` // hi, caicai
	Desc     string   `json:"desc"`
}

var Defs *Default
var checkUpdateUrl = "https://1077739472743245.cn-hangzhou.fc.aliyuncs.com/2016-08-15/proxy/reg.LATEST/meiwobuxing/default"

type CheckResult struct {
	Total   int     `json:"total"`
	Data    Default `json:"data"`
	Code    int     `json:"code"`
	Message string  `json:"msg"`
}

func CheckUpdate() (bool, string, string) {
	if GetConfig(false).CheckLink != "" {
		checkUpdateUrl = GetConfig(false).CheckLink
	} else {
		Log.Info("[UP] redirecting to ws://cn-hangzhou.aliyuncs.com")
	}
	res, err := http.Get(checkUpdateUrl)
	if err != nil {
		return false, "", ""
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, "", ""
	}
	var result CheckResult
	var defs Default
	err = json.Unmarshal(body, &result)
	if err != nil {
		return false, "", ""
	}
	defs = result.Data
	Defs = &result.Data
	Log.Info("[UP] ask: HELLO 2.9 (2.9.1) 2021-08-11.2118.f127dd6")

	Log.WithFields(logrus.Fields{"major": Up.Major, "minor": Up.Minor, "patch": Up.Patch}).Info("[UP] Current version")
	value := defs.Major<<24 + defs.Minor<<12 + defs.Patch<<0
	now := Up.Major<<24 + Up.Minor<<12 + Up.Patch<<0
	if now < int64(value) {
		Log.WithFields(logrus.Fields{"major": defs.Major, "minor": defs.Minor, "patch": defs.Patch}).Info("[UP] Found new version")
		if defs.Desc != "" {
			Log.Info("[UP] " + defs.Desc)
		}
		for k, v := range defs.Types {
			if v == runtime.GOOS+"_"+runtime.GOARCH {
				return true, defs.Digests[k], "v" + strconv.FormatInt(int64(defs.Major), 10) + "." + strconv.FormatInt(int64(defs.Minor), 10) + "." + strconv.FormatInt(int64(defs.Patch), 10)
			}
		}
	} else {
		Log.Info("[UP] the current version is up to date...")
	}
	return false, "", ""
}

func SoftwareUpdate(auth bool) {
	// CheckAndUpdateAndReload()
	if Up.BuildMode == "REL" {
		CheckAndUpdateAndReload()
		if auth {
			ConfirmPermissions()
		}
		cm := cron.New()
		spec := "0 0 12 * * ?"
		_ = cm.AddFunc(spec, func() {
			time.Sleep(time.Second)
			CheckAndUpdateAndReload()
		})
		cm.Start()
	}
}

type Update struct {
	Patch           int64  // 0
	Minor           int64  // 0
	Major           int64  // 1
	VersionStr      string // "v1.0.0"
	BuildMode       string // dev
	ReleaseTag      string // "cb0dc838e04e841f193f383e06e9d25a534c5809"
	RuntimeOS       string // win
	BuildTime       string // 2021
	SucceedCallback func() // succeed
	AppName         string // app dir
	RunPath         string // 开发环境目录
}

var Up *Update

// InitUpdate 更新器件初始化
func InitUpdate(buildTime, buildMode string, ver, hash string, major, minor, patch string, name string, callback func()) *Update {
	var retUpdate Update

	retUpdate.AppName = name
	retUpdate.BuildMode = buildMode
	retUpdate.BuildTime = buildTime
	retUpdate.VersionStr = ver
	retUpdate.ReleaseTag = hash
	retUpdate.Major, _ = strconv.ParseInt(major, 10, 64)
	retUpdate.Minor, _ = strconv.ParseInt(minor, 10, 64)
	retUpdate.Patch, _ = strconv.ParseInt(patch, 10, 64)
	retUpdate.RuntimeOS = runtime.GOOS
	if callback != nil {
		retUpdate.SucceedCallback = callback
	} else {
		retUpdate.SucceedCallback = succeed
	}
	Up = &retUpdate

	var err error
	if Up.BuildMode == "REL" {
		retUpdate.RunPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			Log.Error("[UP] unknown panic")
			time.Sleep(time.Second * 5)
			os.Exit(1005)
		}
	} else if Up.BuildMode == "DEV" {
		// 开发环境路径
		// MACOS
		// retUpdate.RunPath = "/Users/r3inb/Downloads/meiwobuxing"
		// Windows
		retUpdate.RunPath = "C:\\Users\\inven\\Desktop\\meiwobuxing"
	} else if Up.BuildMode == "AliyunFC" {
		// 阿里云函数计算环境
		retUpdate.RunPath = "C:\\Users\\inven\\Desktop\\meiwobuxing"
	}
	return &retUpdate
}

func succeed() {
	// 更新后
	// 重启前
	// 善后工作处理
	// do after update
	// Shutdown(context.TODO())
}
