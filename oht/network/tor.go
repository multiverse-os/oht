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
)

type TorConfig struct {
	ListenPort  string
	WebUIPort   string
	SocksPort   string
	ControlPort string
}

type TorProcess struct {
	Online                bool
	process               *os.Process
	cmd                   *exec.Cmd
	Pid                   string
	OnionHost             string
	WebUIOnionHost        string
	ListenPort            string
	WebUIPort             string
	SocksPort             string
	ControlPort           string
	avoidDiskWrites       int
	hardwareAcceleration  int
	onionServiceDirectory string
	webUIOnionDirectory   string
	dataDirectory         string
	debugLog              string
	binaryFile            string
	configFile            string
	pidFile               string
	authCookie            string
}

func InitializeTor(config *TorConfig) (tor *TorProcess) {
	common.CreatePathUnlessExist("/tor", 0700)
	common.CreatePathUnlessExist("/tor/onion_service", 0700)
	common.CreatePathUnlessExist("/tor/onion_webui", 0700)
	tor = &TorProcess{
		Config:                config,
		Online:                false,
		ListenPort:            config.ListenPort,
		WebUIPort:             config.WebUIPort,
		SocksPort:             config.SocksPort,
		ControlPort:           config.ControlPort,
		avoidDiskWrites:       0,
		hardwareAcceleration:  1,
		dataDirectory:         common.AbsolutePath(common.DefaultDataDir(), "tor/data"),
		debugLog:              common.AbsolutePath(common.DefaultDataDir(), "tor/debug.log"),
		onionServiceDirectory: common.AbsolutePath(common.DefaultDataDir(), "tor/onion_service"),
		webUIOnionDirectory:   common.AbsolutePath(common.DefaultDataDir(), "tor/onion_webui"),
		binaryFile:            common.AbsolutePath("oht/network/tor/bin/linux/64/", "tor"),
		configFile:            common.AbsolutePath(common.DefaultDataDir(), "tor/torrc"),
		pidFile:               common.AbsolutePath(common.DefaultDataDir(), "tor/tor.pid"),
	}
	if runtime.GOOS == "darwin" {
		tor.binaryFile = common.AbsolutePath("oht/network/tor/bin/osx/", "tor")
	} else if runtime.GOOS == "windows" {
		log.Fatal("Tor: No windows binary in the source yet, sorry.")
	}
	if !common.FileExist(tor.configFile) {
		err := tor.initializeConfig()
		if err != nil {
			log.Fatal("Tor: Failed to write configuration: %s", err)
		}
	}
	return tor
}

func (tor *TorProcess) Start() bool {
	tor.cmd = exec.Command(tor.binaryFile, "-f", tor.configFile)
	stdout, _ := tor.cmd.StdoutPipe()
	err := tor.cmd.Start()
	if err != nil {
		log.Fatal("Tor: Failed to start: %s", err)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		//log.Println(line)
		if match, _ := regexp.Match("(100%|Is Tor already running?)", []byte(line)); match {
			tor.Online = true
			break
		}
	}
	tor.process = tor.cmd.Process
	tor.Pid = fmt.Sprintf("%d", tor.cmd.Process.Pid)
	tor.OnionHost = tor.readOnionHost("internal")
	tor.WebUIOnionHost = tor.readOnionHost("webui")
	tor.authCookie = tor.readAuthCookie()
	return tor.Online
}

func (tor *TorProcess) Stop() bool {
	if tor.Online {
		tor.cmd.Process.Kill()
	}
	tor.Online = false
	return tor.Online
}

func (tor *TorProcess) DeleteOnionFiles() bool {
	os.RemoveAll(tor.onionServiceDirectory)
	os.RemoveAll(tor.webUIOnionDirectory)
	return true
}

// CONFIGURATION
func (tor *TorProcess) initializeConfig() error {
	config := ""
	config += fmt.Sprintf("SOCKSPort 127.0.0.1:%s\n", tor.SocksPort)
	config += fmt.Sprintf("ControlPort 127.0.0.1:%s\n", tor.ControlPort)
	config += fmt.Sprintf("DataDirectory %s\n", tor.dataDirectory)
	config += fmt.Sprintf("HardwareAccel %d\n", tor.hardwareAcceleration)
	config += fmt.Sprintf("RunAsDaemon 0\n")
	config += fmt.Sprintf("HiddenServiceDir %s\n", tor.onionServiceDirectory)
	config += fmt.Sprintf("HiddenServicePort %s 127.0.0.1:%s\n", tor.ListenPort, tor.ListenPort)
	config += fmt.Sprintf("HiddenServiceDir %s\n", tor.webUIOnionDirectory)
	config += fmt.Sprintf("HiddenServicePort %s 127.0.0.1:%s\n", tor.WebUIPort, tor.WebUIPort)
	config += fmt.Sprintf("AvoidDiskWrites %d\n", tor.avoidDiskWrites)
	config += fmt.Sprintf("AutomapHostsOnResolve 1\n")
	config += fmt.Sprintf("CookieAuthentication 1\n")
	//config += fmt.Sprintf("SocksPolicy accept 192.168.0.0/16\n")
	//config += fmt.Sprintf("Log debug file %s\n", tor.debugLog)
	err := ioutil.WriteFile(tor.configFile, []byte(config), 0644)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (tor *TorProcess) readOnionHost(serviceType string) string {
	directory := tor.onionServiceDirectory + "/hostname"
	if serviceType == "webui" {
		directory = tor.webUIOnionDirectory + "/hostname"
	}
	onion, err := ioutil.ReadFile(directory)
	if err != nil {
		log.Fatal("Tor: Failed reading onion hostname file: %v", err)
	}
	return strings.Replace(string(onion), "\n", "", -1)
}

// TOR CONTROL
func (tor *TorProcess) readAuthCookie() string {
	cookie, err := ioutil.ReadFile(tor.dataDirectory + "/control_auth_cookie")
	if err != nil {
		log.Fatal("Tor: Failed to read authorization cookie: %s", err)
	}
	return strings.Replace(string(cookie), "\n", "", -1)
}

func (tor *TorProcess) controlCommand(command string) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+tor.ControlPort)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(conn, "AUTHENTICATE\r\n")
	fmt.Fprintf(conn, "%s\fmtr\n", command)
}

func (tor *TorProcess) Cycle() {
	tor.controlCommand("SIGNAL NEWNYM")
}

func (tor *TorProcess) Shutdown() {
	tor.controlCommand("SIGNAL HALT")
}
