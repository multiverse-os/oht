package network

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"../common"
)

type TorProcess struct {
	Process               *os.Process
	Cmd                   *exec.Cmd
	Pid                   string
	OnionHost             string
	OnionWebUIHost        string
	ListenPort            string
	WebUIPort             string
	SocksPort             string
	ControlPort           string
	AvoidDiskWrites       int
	HardwareAcceleration  int
	OnionServiceDirectory string
	OnionWebUIDirectory   string
	DataDirectory         string
	BinaryFile            string
	ConfigFile            string
	PidFile               string
	AuthCookie            string
}

func InitializeTor(listenPort string, socksPort string, controlPort string, webUIPort string) (tor *TorProcess) {

	if !common.FileExist(common.DefaultDataDir() + "/tor") {
		os.MkdirAll(common.DefaultDataDir()+"/tor", os.ModePerm)
	}
	tor = &TorProcess{
		ListenPort:            listenPort,
		WebUIPort:             webUIPort,
		SocksPort:             socksPort,
		ControlPort:           controlPort,
		AvoidDiskWrites:       0,
		HardwareAcceleration:  1,
		DataDirectory:         common.AbsolutePath(common.DefaultDataDir(), "tor/data"),
		OnionServiceDirectory: common.AbsolutePath(common.DefaultDataDir(), "tor/onion_service"),
		OnionWebUIDirectory:   common.AbsolutePath(common.DefaultDataDir(), "tor/onion_webui"),
		BinaryFile:            common.AbsolutePath("network/tor/bin/linux/64/", "tor"),
		ConfigFile:            common.AbsolutePath(common.DefaultDataDir(), "tor/torrc"),
		PidFile:               common.AbsolutePath(common.DefaultDataDir(), "tor/tor.pid"),
	}
	if runtime.GOOS == "darwin" {
		tor.BinaryFile = common.AbsolutePath("network/tor/bin/osx/", "tor")
	} else if runtime.GOOS == "windows" {
		log.Fatal("Tor: No windows binary in the source yet, sorry.")
	}
	if !common.FileExist(tor.ConfigFile) {
		err := tor.InitializeConfig()
		if err != nil {
			log.Fatal("Tor: Failed to write configuration: %s", err)
		}
	}
	torCmd := exec.Command(tor.BinaryFile, "-f", tor.ConfigFile)
	stdout, _ := torCmd.StdoutPipe()
	err := torCmd.Start()
	if err != nil {
		log.Fatal("Tor: Failed to start: %s", err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if match, _ := regexp.Match("(100%|Is Tor already running?)", []byte(line)); match {
			break
		}
	}
	tor.Process = torCmd.Process
	tor.Cmd = torCmd
	tor.Pid = fmt.Sprintf("%d", torCmd.Process.Pid)
	tor.OnionHost = tor.ReadOnionHost("internal")
	tor.OnionWebUIHost = tor.ReadOnionHost("webui")
	tor.AuthCookie = tor.ReadAuthCookie()
	return tor
}

func (tor *TorProcess) ReadOnionHost(serviceType string) string {
	directory := tor.OnionServiceDirectory + "/hostname"
	if serviceType == "webui" {
		directory = tor.OnionWebUIDirectory + "/hostname"
	}
	onion, err := ioutil.ReadFile(directory)
	if err != nil {
		log.Fatal("Tor: Failed reading onion hostname file: %s", err)
	}
	return strings.Replace(string(onion), "\n", "", -1)
}

func (tor *TorProcess) InitializeConfig() error {
	config := ""
	config += fmt.Sprintf("SOCKSPort 127.0.0.1:%s\n", tor.SocksPort)
	config += fmt.Sprintf("ControlPort 127.0.0.1:%s\n", tor.ControlPort)
	config += fmt.Sprintf("DataDirectory %s\n", tor.DataDirectory)
	config += fmt.Sprintf("HardwareAccel %d\n", tor.HardwareAcceleration)
	config += fmt.Sprintf("RunAsDaemon 0\n")
	config += fmt.Sprintf("HiddenServiceDir %s\n", tor.OnionServiceDirectory)
	config += fmt.Sprintf("HiddenServicePort %s 127.0.0.1:%s\n", tor.ListenPort, tor.ListenPort)
	config += fmt.Sprintf("HiddenServiceDir %s\n", tor.OnionWebUIDirectory)
	config += fmt.Sprintf("HiddenServicePort %s 127.0.0.1:%s\n", tor.WebUIPort, tor.WebUIPort)
	config += fmt.Sprintf("AvoidDiskWrites %d\n", tor.AvoidDiskWrites)
	config += fmt.Sprintf("AutomapHostsOnResolve 1\n")
	config += fmt.Sprintf("CookieAuthentication 1\n")
	err := ioutil.WriteFile(tor.ConfigFile, []byte(config), 0644)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (tor *TorProcess) ControlCommand(command string) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+tor.ControlPort)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(conn, "AUTHENTICATE\r\n")
	fmt.Fprintf(conn, "%s\fmtr\n", command)
}

func (tor *TorProcess) Cycle() {
	tor.ControlCommand("SIGNAL NEWNYM")
}

func (tor *TorProcess) Shutdown() {
	tor.ControlCommand("SIGNAL HALT")
}

func (tor *TorProcess) ReadAuthCookie() string {
	cookie, err := ioutil.ReadFile(tor.DataDirectory + "/control_auth_cookie")
	if err != nil {
		log.Fatal("Tor: Failed to read authorization cookie: %s", err)
	}
	return strings.Replace(string(cookie), "\n", "", -1)
}

func FreePort() int {
	address, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
