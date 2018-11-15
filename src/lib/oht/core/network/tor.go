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

	"github.com/multiverse-os/libs/oht/core/common"
)

// TODO: We are definitely going to move to a model where we package within
// our signle binaries other binaries like Tor to simplify release. but
// if we do that we need to build in the checkls for updates against tor
// so this changeble binary code can be updated. but also need to be caituosu
// about the execution of thiks since its a massive security fuck up waiting
// to happen if you are not strict about signatures and checksums and so on

type InitializeConfig struct {
	BinaryPath          string
	Directory           string
	SocksPort           string
	ControlPort         string
	OnionServiceConfigs []*OnionServiceConfig
}

type OnionServiceConfig struct {
	DirectoryName    string
	OnionHost        string
	RemoteListenPort string
	LocalListenPort  string
}

type TorRuntimeConfig struct {
	OnionServiceConfigs    []*OnionServiceConfig
	SocksPort              string
	ControlPort            string
	DNSPort                string
	oRPort                 []string
	extORPort              []string
	dataDirectory          string
	debugLogFile           string
	serverTransportPlugin  []string
	virtualAddrNetworkIPv4 []string
	transPort              []string
	exitPolicy             []string
	dnsPort                []string
	socksPolicy            []string
	runAsDaemon            int
	avoidDiskWrites        int
	hardwareAcceleration   int
	cookieAuthentication   int
	automapHostsOnResolve  int
	bridgeRelay            int
}

type TorProcess struct {
	Online              bool
	authCookie          string
	binaryFile          string
	configFile          string
	pidFile             string
	OnionServiceConfigs []*OnionServiceConfig
	process             *common.Process
	torRC               *TorRuntimeConfig
}

func InitializeTor(config *InitializeConfig) (tor *TorProcess) {
	common.CreatePathUnlessExist("/tor", 0700)
	// Iterate for each onion service and initialize folders
	for i := 0; i < len(tor.OnionServiceConfigs); i++ {
		common.CreatePathUnlessExist(("/tor/" + tor.OnionServiceConfigs[i].DirectoryName), 0700)
	}
	tor = &TorProcess{
		Online:     false,
		configFile: common.AbsolutePath(config.Directory, "tor/torrc"),
		pidFile:    common.AbsolutePath(config.Directory, "tor/tor.pid"),
		torRC: &TorRuntimeConfig{
			SocksPort:            config.SocksPort,
			ControlPort:          config.ControlPort,
			dataDirectory:        common.AbsolutePath(config.Directory, "tor/data"),
			avoidDiskWrites:      0,
			hardwareAcceleration: 1,
			debugLogFile:         common.AbsolutePath(config.Directory, "tor/debug.log"),
			cookieAuthentication: 1,
		},
	}
	if runtime.GOOS == "darwin" {
		tor.binaryFile = common.AbsolutePath((config.BinaryPath + "/osx/"), "tor")
	} else if runtime.GOOS == "windows" {
		log.Fatal("Tor: No windows binary is included the source yet, sorry. It is more secure to build it locally.")
	} else {
		tor.binaryFile = common.AbsolutePath((config.BinaryPath + "/linux/64/"), "tor")
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
	tor.process.cmd = exec.Command(tor.binaryFile, "-f", tor.configFile)
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
	for i := 0; i < len(tor.OnionServerConfigs); i++ {
		tor.OnionServerConfigs[i].OnionHost = tor.readOnionHost(tor.OnionServerConfigs[i].DirectoryName)
	}
	tor.authCookie = tor.readAuthCookie()
	return tor.Online
}

func (tor *TorProcess) Stop(kill bool) bool {
	if tor.Online {
		if kill {
			tor.cmd.Process.Kill()
		} else {
			tor.cmd.Process.Signal(os.Interrupt)
		}
	}
	tor.Online = false
	return tor.Online
}

func (tor *TorProcess) DeleteOnionFiles() bool {
	for i := 0; i < len(tor.OnionServerConfigs); i++ {
		os.RemoveAll(tor.OnionServerConfigs[i].DirectoryName)
	}
	return true
}

// CONFIGURATION
func (tor *TorProcess) initializeConfig() error {
	config := ""
	config += fmt.Sprintf("SOCKSPort 127.0.0.1:%s\n", tor.torRC.SocksPort)
	config += fmt.Sprintf("ControlPort 127.0.0.1:%s\n", tor.torRC.ControlPort)
	config += fmt.Sprintf("DataDirectory %s\n", tor.torRC.dataDirectory)
	config += fmt.Sprintf("HardwareAccel %d\n", tor.torRC.hardwareAcceleration)
	config += fmt.Sprintf("RunAsDaemon %d\n", tor.torRC.runAsDaemon)
	for i := 0; i < len(tor.torRC.OnionServiceConfigs); i++ {
		config += fmt.Sprintf("HiddenServiceDir %s\n", tor.torRC.OnionServiceConfigs.DirectoryName)
		config += fmt.Sprintf("HiddenServicePort %s 127.0.0.1:%s\n", tor.torRC.OnionServerConfigs[i].RemoteListenPort, tor.torRC.OnionServerConfigs[i].LocalListenPort)
	}
	config += fmt.Sprintf("AvoidDiskWrites %d\n", tor.torRC.avoidDiskWrites)
	config += fmt.Sprintf("AutomapHostsOnResolve %d\n", tor.torRC.automapHostsOnResolve)
	config += fmt.Sprintf("CookieAuthentication %d\n", tor.torRC.cookieAuthentication)
	// Should restrict socks access
	//config += fmt.Sprintf("SocksPolicy accept 192.168.0.0/16\n")
	// This does not work because we are using the hacky method of watching stdout
	//config += fmt.Sprintf("Log debug file %s\n", tor.debugLog)
	err := ioutil.WriteFile(tor.configFile, []byte(config), 0644)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (tor *TorProcess) readOnionHost(directoryName string) string {
	directory := tor.dataDirectory + "/" + directoryName + "/hostname"
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
