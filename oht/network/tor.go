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
	"strconv"
	"strings"

	"./../../oht/common"
)

type TorProcess struct {
	process               *os.Process
	cmd                   *exec.Cmd
	Pid                   string
	OnionHost             string
	WebUIOnionHost        string
	ListenPort            int
	WebUIPort             int
	SocksPort             int
	ControlPort           int
	avoidDiskWrites       int
	hardwareAcceleration  int
	onionServiceDirectory string
	webUIOnionDirectory   string
	dataDirectory         string
	binaryFile            string
	configFile            string
	pidFile               string
	authCookie            string
}

func InitializeTor(listenPort int, socksPort int, controlPort int, webUIPort int) (tor *TorProcess) {
	common.CreatePathUnlessExist("/tor", 0700)
	common.CreatePathUnlessExist("/tor/onion_service", 0700)
	common.CreatePathUnlessExist("/tor/onion_webui", 0700)
	tor = &TorProcess{
		ListenPort:            listenPort,
		WebUIPort:             webUIPort,
		SocksPort:             socksPort,
		ControlPort:           controlPort,
		avoidDiskWrites:       0,
		hardwareAcceleration:  1,
		dataDirectory:         common.AbsolutePath(common.DefaultDataDir(), "tor/data"),
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
	torCmd := exec.Command(tor.binaryFile, "-f", tor.configFile)
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
	tor.process = torCmd.Process
	tor.cmd = torCmd
	tor.Pid = fmt.Sprintf("%d", torCmd.Process.Pid)
	tor.OnionHost = tor.readOnionHost("internal")
	tor.WebUIOnionHost = tor.readOnionHost("webui")
	tor.authCookie = tor.readAuthCookie()
	return tor
}

func (tor *TorProcess) readOnionHost(serviceType string) string {
	directory := tor.onionServiceDirectory + "/hostname"
	if serviceType == "webui" {
		directory = tor.webUIOnionDirectory + "/hostname"
	}
	onion, err := ioutil.ReadFile(directory)
	if err != nil {
		log.Fatal("Tor: Failed reading onion hostname file: %s", err)
	}
	return strings.Replace(string(onion), "\n", "", -1)
}

func (tor *TorProcess) initializeConfig() error {
	config := ""
	config += fmt.Sprintf("SOCKSPort 127.0.0.1:%s\n", strconv.Itoa(tor.SocksPort))
	config += fmt.Sprintf("ControlPort 127.0.0.1:%s\n", strconv.Itoa(tor.ControlPort))
	config += fmt.Sprintf("DataDirectory %s\n", tor.dataDirectory)
	config += fmt.Sprintf("HardwareAccel %d\n", tor.hardwareAcceleration)
	config += fmt.Sprintf("RunAsDaemon 0\n")
	config += fmt.Sprintf("HiddenServiceDir %s\n", tor.onionServiceDirectory)
	config += fmt.Sprintf("HiddenServicePort %s 127.0.0.1:%s\n", strconv.Itoa(tor.ListenPort), strconv.Itoa(tor.ListenPort))
	config += fmt.Sprintf("HiddenServiceDir %s\n", tor.webUIOnionDirectory)
	config += fmt.Sprintf("HiddenServicePort %s 127.0.0.1:%s\n", strconv.Itoa(tor.WebUIPort), strconv.Itoa(tor.WebUIPort))
	config += fmt.Sprintf("AvoidDiskWrites %d\n", tor.avoidDiskWrites)
	config += fmt.Sprintf("AutomapHostsOnResolve 1\n")
	config += fmt.Sprintf("CookieAuthentication 1\n")
	err := ioutil.WriteFile(tor.configFile, []byte(config), 0644)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (tor *TorProcess) controlCommand(command string) {
	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(tor.ControlPort))
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

func (tor *TorProcess) readAuthCookie() string {
	cookie, err := ioutil.ReadFile(tor.dataDirectory + "/control_auth_cookie")
	if err != nil {
		log.Fatal("Tor: Failed to read authorization cookie: %s", err)
	}
	return strings.Replace(string(cookie), "\n", "", -1)
}
