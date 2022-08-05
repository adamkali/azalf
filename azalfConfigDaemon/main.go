package main

/*

 ▄▀▀█▄   ▄▀▀▀▀▄    ▄▀▀█▄   ▄▀▀▀▀▄     ▄▀▀▀█▄        ▄▀▀█▄▄   ▄▀▀█▄   ▄▀▀█▄▄▄▄  ▄▀▀▄ ▄▀▄  ▄▀▀▀▀▄   ▄▀▀▄ ▀▄
▐ ▄▀ ▀▄ █     ▄▀  ▐ ▄▀ ▀▄ █    █     █  ▄▀  ▀▄     █ ▄▀   █ ▐ ▄▀ ▀▄ ▐  ▄▀   ▐ █  █ ▀  █ █      █ █  █ █ █
  █▄▄▄█ ▐ ▄▄▀▀      █▄▄▄█ ▐    █     ▐ █▄▄▄▄       ▐ █    █   █▄▄▄█   █▄▄▄▄▄  ▐  █    █ █      █ ▐  █  ▀█
 ▄▀   █   █        ▄▀   █     █       █    ▐         █    █  ▄▀   █   █    ▌    █    █  ▀▄    ▄▀   █   █
█   ▄▀     ▀▄▄▄▄▀ █   ▄▀    ▄▀▄▄▄▄▄▄▀ █             ▄▀▄▄▄▄▀ █   ▄▀   ▄▀▄▄▄▄   ▄▀   ▄▀     ▀▀▀▀   ▄▀   █
▐   ▐          ▐  ▐   ▐     █        █             █     ▐  ▐   ▐    █    ▐   █    █             █    ▐
                            ▐        ▐             ▐                 ▐        ▐    ▐             ▐


                                  ....
                                .'' .'''
.                             .'   :
\\                          .:    :
 \\                        _:    :       ..----.._
  \\                    .:::.....:::.. .'         ''.
   \\                 .'  #-. .-######'     #        '.
    \\                 '.##'/ ' ################       :
     \\                  #####################         :
      \\               ..##.-.#### .''''###'.._        :
       \\             :--:########:            '.    .' :
        \\..__...--.. :--:#######.'   '.         '.     :
        :     :  : : '':'-:'':'::        .         '.  .'
        '---'''..: :    ':    '..'''.      '.        :'
           \\  :: : :     '      ''''''.     '.      .:
            \\ ::  : :     '            '.      '      :
             \\::   : :           ....' ..:       '     '.
              \\::  : :    .....####\\ .~~.:.             :
               \\':.:.:.:'#########.===. ~ |.'-.   . '''.. :
                \\    .'  ########## \ \ _.' '. '-.       '''.
                :\\  :     ########   \ \      '.  '-.        :
               :  \\'    '   #### :    \ \      :.    '-.      :
              :  .'\\   :'  :     :     \ \       :      '-.    :
             : .'  .\\  '  :      :     :\ \       :        '.   :
             ::   :  \\'  :.      :     : \ \      :          '. :
             ::. :    \\  : :      :    ;  \ \     :           '.:
              : ':    '\\ :  :     :     :  \:\     :        ..'
                 :    ' \\ :        :     ;  \|      :   .'''
                 '.   '  \\:                         :.''
                  .:..... \\:       :            ..''
                 '._____|'.\\......'''''''.:..'''
                            \\
*/

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/sevlyar/go-daemon"
)

var (
	daemonName = "AZALF"
	daemonDesc = "Adam's Zillenial Arch Linux Flavor Daemon"
	daemonUser = "root"
	daemonPid  = "/var/run/azalfConfigDaemon.pid"
	daemonLog  = "/var/log/azalfConfigDaemon.log"
	daemonPwd  = "/var/run/azalfConfigDaemon.pwd"
	daemonHar  = "/var/run/azalfConfigDaemon/hardware.json"

	signal = flag.String("s", "", `Send signal to the daemon:
	quit — graceful shutdown
	stop — fast shutdown
	reload — reloading the configuration file`)
)

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

type Config struct {
	Colors struct {
		Background         string `yaml:"background"`
		Foreground         string `yaml:"foreground"`
		Important          string `yaml:"important"`
		Active             string `yaml:"active"`
		Primary            string `yaml:"primary"`
		Secondary          string `yaml:"secondary"`
		Highlight          string `yaml:"highlight"`
		HighlightSecondary string `yaml:"highlight-secondary"`
		WidgetBackground   string `yaml:"widget-background"`
		WidgetForeground   string `yaml:"widget-foreground"`
		WidgetWarn         string `yaml:"widget-warn"`
		WidgetDanger       string `yaml:"widget-danger"`
		Info               string `yaml:"info"`
		Success            string `yaml:"success"`
		ShellBackground    string `yaml:"shell-background"`
		ShellBackgroundAlt string `yaml:"shell-background-alt"`
		ShellInfo          string `yaml:"shell-info"`
		Other              string `yaml:"other"`
		ShellDark          string `yaml:"shell-dark"`
	} `yaml:"colors"`
	Fonts struct {
		Monospace string `yaml:"monospace"`
		SansSerif string `yaml:"sans-serif"`
		Serif     string `yaml:"serif"`
		Emoji     string `yaml:"emoji"`
	} `yaml:"font-families"`
	Sizing struct {
		FontSize struct {
			Shell          int `yaml:"shell"`
			ShellLarge     int `yaml:"shell-large"`
			Normal         int `yaml:"normal"`
			Bar            int `yaml:"bar"`
			Dashboard      int `yaml:"dashboard"`
			DashboardLarge int `yaml:"dashboard-large"`
		} `yaml:"font-size"`
		Padding struct {
			X int `yaml:"x"`
			Y int `yaml:"y"`
		} `yaml:"padding"`
		Spacing      int `yaml:"spacing"`
		BorderRadius int `yaml:"border-radius"`
		BorderWidth  int `yaml:"border-width"`
	} `yaml:"sizing"`
}

type HardwareInfo struct {
	CPU struct {
		Model    string `yaml:"model"; json:"model"`
		Cores    int    `yaml:"cores"; json: "cores"`
		Speed    int    `yaml:"speed"; json: speed"`
		CPUUsage string `yaml:"cpu_usage; json:"cpu_usage"`
	} `yaml:"cpu"; json:"cpu"`
	GPU struct {
		Model       string `yaml:"model"; json:"model"`
		Driver      string `yaml:"driver"; json:"driver"`
		GPUUsage    string `yaml:"gpu_usage"; json:"gpu_usage"`
		VRamUsed    string `yaml:"vram_usage"; json:"vram_usage"`
		VRamTotal   string `yaml:"vram_total"; json:"vram_total"`
		VRamPercent string `yaml:"vram_usage"; json:"vram_usage_percent"`
	} `yaml:"gpu"; json:"gpu"`
	RAM struct {
		Total    int `yaml:"total"; json:"total"`
		Used     int `yaml:"used"; json:"used"`
		Frequncy int `yaml:"frequency"; json:"frequency"`
	} `yaml:"ram"; json:"ram"`
	Drives []struct {
		Name      string `json:"name"`
		Size      string `json:"size"`
		Used      string `json:"used"`
		Available string `json:"available"`
		Usage     string `json:"usage"`
	} `json:"drives"`
	Network []struct {
		Name     string `json:"name"`
		Speed    string `json:"speed"`
		Download string `json:"download"`
		Upload   string `json:"upload"`
	} `json:"network"`
	Power struct {
		Battery        string `json:"battery"`
		BatteryPercent string `json:"battery_percent"`
		CPUPower       string `json:"cpu_power"`
		GPUPower       string `json:"gpu_power"`
		TotalPower     string `json:"total_power"`
	} `json:"power"`
	Temperature struct {
		CPU string `json:"cpu"`
		GPU string `json:"gpu"`
	} `json:"temperature"`
}

type CPUUsageTunnel struct {
	Usage int `json:"usage"`
	Total int `json:"total"`
}

func main() {

	flag.Parse()
	daemon.AddCommand(daemon.StringFlag(signal, "quit"), syscall.SIGQUIT, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "stop"), syscall.SIGTERM, termHandler)
	daemon.AddCommand(daemon.StringFlag(signal, "reload"), syscall.SIGHUP, reloadHandler)

	context :=
		&daemon.Context{
			PidFileName: daemonPid,
			PidFilePerm: 0644,
			LogFileName: daemonLog,
			LogFilePerm: 0640,
			WorkDir:     "/var/run/azalfConfigDaemon",
			Umask:       027,
			Args:        []string{"azalfConfigDaemon"},
		}

	if len(daemon.ActiveFlags()) > 0 {
		d, err := context.Search()
		if err != nil {
			log.Fatalf("Unable send signal to the daemon: ", err)
		}
		daemon.SendCommands(d)
		return
	}

	d, err := context.Reborn()
	if err != nil {
		log.Fatalf("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer context.Release()

	log.Printf("%s has started brewing spells", daemonName)
	log.Printf("%s is now casting spells!", daemonDesc)

	//
	// get the .azalf.yml file
	// the .azalf.yml file is going to hold
	// the look and feel of the OS
	//
	// neovim could call this daemon by an http request
	// and then depending on the request, it could
	//
	var config Config
	conf, err := ioutil.ReadFile("~/.azalf.yml")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = yaml.Unmarshal(conf, &config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// create hardware info to be written to
	// /var/run/azalfConfigDaemon/hardware.json
	//
	// this is going to be used by dashboard and the bar
	// to display hardware info
	hardwareInfo := HardwareInfo{}
	go hardwareListener(&hardwareInfo)

	go worker(&config)

	// make an anonymous go routine to listen for server requests

	// wg := sync.WaitGroup{}
	// wg.Add(1) // one for now
	// go func() {
	// 	defer wg.Done()
	// 	worker()
	// }()
	// wg.Wait()

	err = daemon.ServeSignals()
	if err != nil {
		log.Printf("AZALF's spells blew up in his face: ", err)
	}

	log.Printf("%s has gone to sleep.", daemonName)

}

func worker(c *Config) {
	for {
		// This will be where the server is started
		// and the server is listening for requests
		// and then the server is closed if the daemon
		// is stopped or if the daemon is killed
		//
		// The server is started here
		// the server will be listening on localhost:9999
		serveHTTP(c)

		select {
		case <-stop:
			done <- struct{}{}
			return
		case <-time.After(time.Second):
			log.Printf("%s is still alive and brewing your config.", daemonName)
		}
	}
}

func hardwareListener(h *HardwareInfo) {
	for {
		// This will be where the Daemon listens to the hardware
		// then it will write the hardware info to the HardwareInfo struct
		// which can be accessed by the server
		// then the user can make calls to th server localy and get the
		// hardware info that they want.
		//
		time.Sleep(time.Second / 2)

		select {
		case <-stop:
			done <- struct{}{}
			return
		case <-time.After(time.Second):
			log.Printf("%s is monotoring the magic.", daemonName)
		}
	}
}

// TODO: Finish this
func (h *HardwareInfo) getHardwareInfo() {
	// Check that the os is linux
	// if not, return and write Fatalf

	if runtime.GOOS != "linux" {
		log.Fatalf("%s can only do his magic on linux", daemonName)
		return
	}

	// Get the hardware info
	h.getCPUInfo()
	h.getGPUInfo()

	// Write the hardware info to the hardware.json file
	//
	// return the hardware info
}

func (h *HardwareInfo) getCPUInfo() {
	// Get the CPU info from /proc/cpuinfo
	//
	// first get the name of the CPU
	// then get the number of cores
	// then get the frequency of the CPU
	//
	// then write to the hardware.json file

	// first read the cpuinfo file
	// then parse the cpuinfo file

	cpuInfoFile, err := os.Open("/proc/cpuinfo")
	if err != nil {
		log.Fatalf("%s could not open /proc/cpuinfo: ", daemonName, err)
	}

	cpuInfo, err := ioutil.ReadAll(cpuInfoFile)

	cpuscan := bufio.NewScanner(bytes.NewBuffer(cpuInfo))
	for cpuscan.Scan() {

		/* CPU {
			Model string 'yaml:"model"; json:"model"'
			Cores int 'yaml:"cores"; json: "cores"'
			Speed int 'yaml:"speed"; json: speed"'
			CPUUsage string `yaml:"cpu_usage; json:"cpu_usage"`
			CPUUsagePercent string `yaml:"cpu_usage_prcent"; json:"cpu_usage_percent"`
		} */

		if strings.Contains(strings.ToLower(cpuscan.Text()), "model name") {
			h.CPU.Model = strings.Split(cpuscan.Text(), ":")[1]
		} else if strings.Contains(strings.ToLower(cpuscan.Text()), "cpu cores") {
			h.CPU.Cores, err = strconv.Atoi(strings.Split(cpuscan.Text(), ":")[1])
		} else if strings.Contains(strings.ToLower(cpuscan.Text()), "cpu mhz") {
			h.CPU.Speed, err = strconv.Atoi(strings.Split(cpuscan.Text(), ":")[1])
		}
		if err != nil {
			log.Fatalf("%s could not convert the cpu info: ", daemonName, err)
		}
	}
	var d1 CPUUsageTunnel = calculateInstantaneousUsage()
	var d2 CPUUsageTunnel = calculateInstantaneousUsage()

	h.CPU.CPUUsage = string((d2.Usage - d1.Usage) / (d2.Total - d1.Total))

	defer cpuInfoFile.Close()
}

func calculateInstantaneousUsage() CPUUsageTunnel {

	// load /proc/stat
	// then parse the /proc/stat file
	// then put the information in the HardwareInfo struct
	// then write to the hardware.json file
	// then return
	cpustat, err := os.Open("/proc/stat")
	if err != nil {
		// write to the err file
		// then return
		log.Fatalf("%s could not scry /proc/stat: %s", daemonName, err.Error())
		return CPUUsageTunnel{}
	}
	// then parse the cpuinfo file
	cpustatscan := bufio.NewScanner(cpustat)

	// calculate the CPU usage
	// total: user + nice + system + idle + iowait + irq + softirq + (steal ? steal : 0)
	// idle : idle + iowait
	// used : total - idle
	// then we can calculate the CPU usage
	// CPUUsage = used / total * 100

	var cpuUsageInfo []string
	for cpustatscan.Scan() {
		if strings.Contains(strings.ToLower(cpustatscan.Text()), "cpu") {
			cpuUsageInfo = strings.Split(cpustatscan.Text(), " ")
		}
	}

	var totalTime int
	var idleTime int
	var usedTime int

	for i := 1; i < 9; i++ {
		c, err := strconv.Atoi(cpuUsageInfo[i])
		if err != nil {
			log.Fatalf("%s could not convert the cpu info: ", daemonName, err)
			break
		}
		totalTime += c
	}

	idle, err := strconv.Atoi(cpuUsageInfo[4])
	if err != nil {
		log.Fatalf("%s could not convert the cpu info: ", daemonName, err)
	}
	iowait, err := strconv.Atoi(cpuUsageInfo[5])
	if err != nil {
		log.Fatalf("%s could not convert the cpu info: ", daemonName, err)
	}

	idleTime = idle + iowait // idle + iowait
	usedTime = totalTime - idleTime

	return CPUUsageTunnel{
		Usage: usedTime,
		Total: totalTime,
	}
}

/* GPU {
	Model string 'yaml:"model"; json:"model"'
	Driver string 'yaml:"driver"; json:"driver"'
	GPUUsage string `yaml:"gpu_usage"; json:"gpu_usage"`
	VRamUsed string `yaml:"vram_usage"; json:"vram_usage"`
	VRamTotal string `yaml:"vram_total"; json:"vram_total"`
	VRamPercent string `yaml:"vram_usage"; json:"vram_usage_percent"`
} `yaml:"gpu"; json:"gpu"` */

// TODO: Finish this
// get the GPU information by using lspci
func (h *HardwareInfo) getGPUInfo() {
	// Get the GPU info from lspci
	//
	// first get the name of the GPU
	// then get the driver of the GPU
	// then get the total amount of VRAM
	pciDevices, err := exec.Command("lspci", "-v").Output()
	if err != nil {
		// write to the err file
		// then return
		log.Fatalf("%s could not scry lspci: %s", daemonName, err.Error())
		return
	}

	// then parse the lspci output
	// then put the information in the HardwareInfo struct
	// then write to the hardware.json file
	// then return
	pciDevicesScan := bufio.NewScanner(bytes.NewReader(pciDevices))
	var gpuInfo string
	var breaker bool
	for pciDevicesScan.Scan() {

		if strings.Contains(strings.ToLower(pciDevicesScan.Text()), "vga compatible controller") {
			breaker = true
			gpuInfo = pciDevicesScan.Text()
		}
		if breaker {
			break
		}
	}

	// make a list of strings to hold the information
	var gpuInfoList []string
	gpuInfoList = strings.Split(gpuInfo, "\n")

	// loop through the list of strings
	// and get the information
	for i := 0; i < len(gpuInfoList); i++ {

	}
}

func (config *Config) GetColor(color string) string {
	var returnColor string
	switch color {
	// match the color string to the correct config.Color value
	// Then get the color value from the config.Color struct
	// there should be a case for each color in config.Color struct
	// then return what the color value is from the struct
	// given that these are the colors
	// 		Background         string `yaml:"background"`
	//		Foreground         string `yaml:"foreground"`
	//		Important          string `yaml:"important"`
	//		Active             string `yaml:"active"`
	//		Primary            string `yaml:"primary"`
	//		Secondary          string `yaml:"secondary"`
	//		Highlight          string `yaml:"highlight"`
	//		HighlightSecondary string `yaml:"highlight-secondary"`
	//		WidgetBackground   string `yaml:"widget-background"`
	//		WidgetForeground   string `yaml:"widget-foreground"`
	//		WidgetWarn         string `yaml:"widget-warn"`
	//		WidgetDanger       string `yaml:"widget-danger"`
	//		Info               string `yaml:"info"`
	//		Success            string `yaml:"success"`
	//		ShellBackground    string `yaml:"shell-background"`
	//		ShellBackgroundAlt string `yaml:"shell-background-alt"`
	//		ShellInfo          string `yaml:"shell-info"`
	//		other              string `yaml:"other"`
	//		ShellDark          string `yaml:"shell-dark"`
	case "background":
		returnColor = config.Colors.Background
	case "forground":
		returnColor = config.Colors.Foreground
	case "important":
		returnColor = config.Colors.Important
	case "active":
		returnColor = config.Colors.Active
	case "primary":
		returnColor = config.Colors.Primary
	case "secondary":
		returnColor = config.Colors.Secondary
	case "highlight":
		returnColor = config.Colors.Highlight
	case "highlight-secondary":
		returnColor = config.Colors.HighlightSecondary
	case "widget-background":
		returnColor = config.Colors.WidgetBackground
	case "widget-foreground":
		returnColor = config.Colors.WidgetForeground
	case "widget-warn":
		returnColor = config.Colors.WidgetWarn
	case "widget-danger":
		returnColor = config.Colors.WidgetDanger
	case "info":
		returnColor = config.Colors.Info
	case "success":
		returnColor = config.Colors.Success
	case "shell-background":
		returnColor = config.Colors.ShellBackground
	case "shell-background-alt":
		returnColor = config.Colors.ShellBackgroundAlt
	case "shell-info":
		returnColor = config.Colors.ShellInfo
	case "other":
		returnColor = config.Colors.Other
	case "shell-dark":
		returnColor = config.Colors.ShellDark
	default:
		log.Fatalf("%s could not find the color: %s", daemonName, color)
	}
	return returnColor
}

func (config *Config) GetFont(font string) string {
	var returnFont string
	switch font {
	// match the font string to the correct config.Font value
	// Then get the font value from the config.Font struct
	// there should be a case for each font in config.Font struct
	// then return what the font value is from the struct
	// given that these are the fonts are
	// 		Monospace string `yaml:"monospace"`
	// 		SansSerif string `yaml:"sans-serif"`
	// 		Serif     string `yaml:"serif"`
	// 		Emoji     string `yaml:"emoji"`

	case "monospace":
		returnFont = config.Fonts.Monospace
	case "sans-serif":
		returnFont = config.Fonts.SansSerif
	case "serif":
		returnFont = config.Fonts.Serif
	case "emoji":
		returnFont = config.Fonts.Emoji
	default:
		log.Fatalf("%s could not find the font: %s", daemonName, font)
	}
	return returnFont
}

func (config *Config) GetSizes(size string) string {
	// match the size string to the correct config.Sizing value
	// Then get the size value from the config.Sizing struct
	// given a valid size call
	// given that these are the sizes
	//		Sizing struct {
	//			FontSize struct {
	//				Shell          int `yaml:"shell"`
	//				ShellLarge     int `yaml:"shell-large"`
	//				Normal         int `yaml:"normal"`
	//				bar            int `yaml:"bar"`
	//				dashboard      int `yaml:"dashboard"`
	//				dashboardLarge int `yaml:"dashboard-large"`
	//			} `yaml:"font-size"`
	//			Padding struct {
	//				X int `yaml:"x"`
	//				Y int `yaml:"y"`
	//			} `yaml:"padding"`
	//			Spacing      int `yaml:"spacing"`
	//			BorderRadius int `yaml:"border-radius"`
	//			BorderWidth  int `yaml:"border-width"`
	//		} `yaml:"sizing"`
	//
	// so to get the Config.Sizing.FontSize.Shell value in the size string
	// you would do sizing-fontsize-shell
	// then you would get the value from the config.Sizing.FontSize struct
	var returnSize string
	switch size {
	case "sizing-fontsize-shell":
		returnSize = strconv.Itoa(config.Sizing.FontSize.Shell)
	case "sizing-fontsize-shell-large":
		returnSize = strconv.Itoa(config.Sizing.FontSize.ShellLarge)
	case "sizing-fontsize-normal":
		returnSize = strconv.Itoa(config.Sizing.FontSize.Normal)
	case "sizing-fontsize-bar":
		returnSize = strconv.Itoa(config.Sizing.FontSize.Bar)
	case "sizing-fontsize-dashboard":
		returnSize = strconv.Itoa(config.Sizing.FontSize.Dashboard)
	case "sizing-fontsize-dashboard-large":
		returnSize = strconv.Itoa(config.Sizing.FontSize.DashboardLarge)
	case "sizing-padding-x":
		returnSize = strconv.Itoa(config.Sizing.Padding.X)
	case "sizing-padding-y":
		returnSize = strconv.Itoa(config.Sizing.Padding.Y)
	case "sizing-spacing":
		returnSize = strconv.Itoa(config.Sizing.Spacing)
	case "sizing-border-radius":
		returnSize = strconv.Itoa(config.Sizing.BorderRadius)
	case "sizing-border-width":
		returnSize = strconv.Itoa(config.Sizing.BorderWidth)
	default:
		log.Fatalf("%s could not find the size: %s", daemonName, size)
	}
	return returnSize
}

// make a handler for the server
func (config *Config) AzalfHandler(w http.ResponseWriter, r *http.Request) {
	// Check the method of the request
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	} else {
		// Check the path of the request
		// There should be paths of
		// - /config 					
		// - /config/colors 			
		// - /config/colors/{:color} 	
		// - /config/fonts				
		// - /config/fonts/{:font} 		
		// - /config/sizing 			
		// - /config/sizing/{:size} 	

		switch {
		case r.URL.Path == "/config":
			config.getConfigObject(w, r)
		case r.URL.Path == "/config/colors":
			config.getConfigColor(w, r)
		case r.URL.Path == "/config/fonts":
			config.getConfigFont(w, r)
		case r.URL.Path == "/config/sizing":
			config.getConfigSizing(w, r)
		case r.URL.Path == "/config/colors/":
			config.getConfigColor(w, r)
		case r.URL.Path == "/config/fonts/":
			config.getConfigFont(w, r)
		case r.URL.Path == "/config/sizing/":
			config.getConfigSizing(w, r)
		default:
			http.Error(w, "Not found", http.StatusNotFound)
			fmt.Errorf("AZALF could not find that configuration.", r.URL.Path)
		}
	}
	return // return the config object
}

func (config *Config) getConfigObject(w http.ResponseWriter, r *http.Request) {
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", daemonName, err.Error())
	}

	// return the json object as bytes for the response
	return
}

func (config *Config) getConfigColor(w http.ResponseWriter, r *http.Request) {
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Colors)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", daemonName, err.Error())
	}

	// return the json object as bytes for the response
	return
}

func (config *Config) getConfigFont(w http.ResponseWriter, r *http.Request) {
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Fonts)

	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", daemonName, err.Error())
	}
	// return the json object as bytes for the response

	return
}

func (config *Config) getConfigSizing(w http.ResponseWriter, r *http.Request) {
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Sizing)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", daemonName, err.Error())
	}

	return
}

func (config *Config) getConfigSpcecificColor(w http.ResponseWriter, r *http.Request) {
	// GET THE COLOR NAME FROM THE URL
	// seperate the path from the url
	// get the color name from the path
	specificColorPath := strings.Split(r.URL.Path, "/")
	specificColor := specificColorPath[3]

	// get the color from the config object
	specificColor = config.GetColor(specificColor)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&specificColor)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", daemonName, err.Error())
	}

	return
}

func (config *Config) getConfigSpecificFont(w http.ResponseWriter, r *http.Request) {
	// GET THE FONT NAME FROM THE URL
	// seperate the path from the url
	// get the font name from the path
	specificFontPath := strings.Split(r.URL.Path, "/")
	specificFont := specificFontPath[3]

	// get the font from the config object
	specificFont = config.GetFont(specificFont)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&specificFont)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", daemonName, err.Error())
	}

	return
}

func (config *Config) getConfigSpecificSizing(w http.ResponseWriter, r *http.Request) {
	// GET THE SIZE NAME FROM THE URL
	// seperate the path from the url
	// get the size name from the path
	specificSizePath := strings.Split(r.URL.Path, "/")
	specificSize := specificSizePath[3]

	// get the size from the config object
	specificSize = config.GetSizes(specificSize)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&specificSize)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", daemonName, err.Error())
	}

	return
}

func serveHTTP(config *Config) {
	http.HandleFunc("/config", config.AzalfHandler)
	http.ListenAndServe("127.0.0.1:9999", nil)
}

func termHandler(sig os.Signal) error {
	log.Println(
		"%s shall be terminated. Hasta La Vista... %s", daemonName, daemonName,
	)
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
	return daemon.ErrStop
}

func reloadHandler(sig os.Signal) error {
	log.Printf("%s fell asleep", daemonName)
	return nil
}
