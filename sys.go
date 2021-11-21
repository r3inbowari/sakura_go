package Sakura

import (
	. "github.com/klauspost/cpuid/v2"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func ConfirmPermissions() {
	id := GetID()
	Log.Info("[SYS] Permissions key -> " + id)
	for _, v := range Defs.PDigests {
		if id == v {
			return
		}
	}
	ExitOops()
}

func ExitOops() {
	Log.Warn("[SYS] oops, you don't have permission. please contact the developer (⑉･̆-･̆⑉)")
	time.Sleep(time.Minute * 3)
	os.Exit(999)
}

func GetID() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("[SYS] panic" + err.Error())
	}
	add := ""
	for _, inter := range interfaces {
		add += inter.HardwareAddr.String()
	}
	add += MixCPUInfo()
	return GetMD5(add)
}

func MixCPUInfo() string {
	v := strings.Join(CPU.FeatureSet(), ",")
	r := CPU.PhysicalCores + CPU.ThreadsPerCore + CPU.LogicalCores + CPU.Family + CPU.Model + CPU.CacheLine + CPU.Cache.L1D*CPU.Cache.L1I
	t := strconv.Itoa(r<<3 + int(CPU.VendorID))
	return GetMD5(v + t)
}
