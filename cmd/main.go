package main

import (
	Sakura "github.com/r3inbowari/sakura_go"
	"os"
	"runtime"
	"strconv"
)

var (
	gitHash        string
	buildTime      string
	goVersion      string
	releaseVersion string
	major          string
	minor          string
	patch          string
)

var mode = "DEV"

func main() {
	dev()
	Sakura.InitConfig()
	Sakura.InitUpdate(buildTime, mode, releaseVersion, gitHash, major, minor, patch, "sakura", nil)
	Sakura.InitLogger(Sakura.Up.BuildMode, AppInfo)
	Sakura.Init()
}

func dev() {
	if mode == "DEV" {
		buildTime = "Thu Oct 01 00:00:00 1970 +0800"
		gitHash = "cb0dc838e04e841f193f383e06e9d25a534c5809"
		goVersion = runtime.Version()
		releaseVersion = "ver[DEV]"
	}
}

func AppInfo() {
	au := "cyt(r3inbowari)"
	Sakura.Blue("                                  ______             _                 ")
	Sakura.Blue("                                 |  ____|           (_)                PACKAGER #UNOFFICIAL " + Sakura.Up.ReleaseTag[:7] + "..." + Sakura.Up.ReleaseTag[33:])
	Sakura.Blue("  _ __ ___   ___  _ __ ___   ___ | |__   _ __   __ _ _ _ __   ___      -... .. .-.. .. -.-. --- .. -. " + Sakura.Up.VersionStr)
	Sakura.Blue(" | '_ ` _ \\ / _ \\| '_ ` _ \\ / _ \\|  __| | '_ \\ / _` | | '_ \\ / _ \\     Running: CLI Server" + " by " + au)
	Sakura.Blue(" | | | | | | (_) | | | | | | (_) | |____| | | | (_| | | | | |  __/     Port: " + Sakura.GetConfig(false).APIAddr[1:])
	Sakura.Blue(" |_| |_| |_|\\___/|_| |_| |_|\\___/|______|_| |_|\\__, |_|_| |_|\\___|     PID: " + strconv.Itoa(os.Getpid()))
	Sakura.Blue("                                                __/ |                  built on " + Sakura.Up.BuildTime)
	Sakura.Blue("                                               |___/              ")
	Sakura.Blue("")
}
