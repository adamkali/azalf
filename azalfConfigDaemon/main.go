package main

/*

 ▄▀▀█▄   ▄▀▀▀▀▄    ▄▀▀█▄   ▄▀▀▀▀▄     ▄▀▀▀█▄        ▄▀▀▀▀▄  ▄▀▀█▄▄▄▄  ▄▀▀▄▀▀▀▄  ▄▀▀▄ ▄▀▀▄  ▄▀▀█▄▄▄▄  ▄▀▀▄▀▀▀▄
▐ ▄▀ ▀▄ █     ▄▀  ▐ ▄▀ ▀▄ █    █     █  ▄▀  ▀▄     █ █   ▐ ▐  ▄▀   ▐ █   █   █ █   █    █ ▐  ▄▀   ▐ █   █   █
  █▄▄▄█ ▐ ▄▄▀▀      █▄▄▄█ ▐    █     ▐ █▄▄▄▄          ▀▄     █▄▄▄▄▄  ▐  █▀▀█▀  ▐  █    █    █▄▄▄▄▄  ▐  █▀▀█▀
 ▄▀   █   █        ▄▀   █     █       █    ▐       ▀▄   █    █    ▌   ▄▀    █     █   ▄▀    █    ▌   ▄▀    █
█   ▄▀     ▀▄▄▄▄▀ █   ▄▀    ▄▀▄▄▄▄▄▄▀ █             █▀▀▀    ▄▀▄▄▄▄   █     █       ▀▄▀     ▄▀▄▄▄▄   █     █
▐   ▐          ▐  ▐   ▐     █        █              ▐       █    ▐   ▐     ▐               █    ▐   ▐     ▐
                            ▐        ▐                      ▐                              ▐


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

/*
	A few notes about the server:
	- This is a simple server that will be used to monitor the hardware
		and serve the config to many programs within the machine.
	- This is configurable through the ~/.config/azalf/.azalf.yaml file
		You can change the config file to your liking, primarly through
		html color codes. The config file also holds sizing information
		for any sizing elements.
	- In this way I want to be able to create a simple homogeneous design
		for the my operating system.
		for the time being I only have the config server in a state that
		I feel is good enough for production, but that is mosty because
		of the gettting gpu information. :/
	- Then the you can simple do a http://localhost:9999/config request to get
		the config file in your configurable application's request method.
	- You can also tell the server to refresh the config file by doing a
		azalf update command.

	Happy hacking!

*/

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	serverName = "AZALF"
	serverDesc = "Adam's Zillenial Arch Linux Flavor Server"
	serverPort = ":9999"
	serverVers = "0.2.4"
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
	// Make a log file for the server.
	logFile, err := os.Create("azalf.log")
	if err != nil {
		log.Fatal("Failed to create log file.")
	}

	// Check that the user is root and on linux
	if runtime.GOOS != "linux" {
		log.Fatalf("%s can only do his magic on linux", serverName)
	}

	// Loop until the user calls `azalf stop`
	for {
		// Initialize the config
		var conf Config

		// The server should be accept args to update the config file, even while running.
		// Allow the user to call the program with an argument to update the config
		// file while the program is running.
		if len(os.Args) > 1 {
			// Make a --help command to show the user how to use the program.
			if os.Args[1] == "--help" || os.Args[1] == "-h" {
				fmt.Printf("%s %s\n", serverName, serverVers)
				fmt.Printf("%s\n", serverDesc)
				fmt.Println("")
				// Print the ascii art to the user.
				fmt.Println(`
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
                            \\`)
				fmt.Println("A few notes about the server:")
				fmt.Println("- This is a simple server that will be used to monitor the hardware")
				fmt.Println("  and serve the config to many programs within the machine.")
				fmt.Println("- This is configurable through the ~/.config/azalf/.azalf.yaml file")
				fmt.Println("  You can change the config file to your liking, primarly through")
				fmt.Println("  html color codes. The config file also holds sizing information")
				fmt.Println("  for any sizing elements.")
				fmt.Println("- In this way I want to be able to create a simple homogeneous design")
				fmt.Println("  for the my operating system. for the time being I only have the config")
				fmt.Println("  server in a state that I feel is good enough for production, but that")
				fmt.Println("  is mostly because of the getting gpu information. :/")
				fmt.Println("- Then the you can simple do a http://localhost:9999/config request to get")
				fmt.Println("  the config file in your configurable application's request method.")
				fmt.Println("- You can also tell the server to refresh the config file by doing a")
				fmt.Println("  azalf update command.")
				fmt.Println("Happy hacking!")
				fmt.Println("")
				fmt.Println("Usage:")
				fmt.Println("  azalf update")
				fmt.Println("  azalf stop")
				fmt.Println("  azalf spells")
				fmt.Println("  azalf --help || -h")
				fmt.Println("  azalf --version || -v")
				fmt.Println("  azalf --liscense || -l")
				fmt.Println("")
				fmt.Println("Example Use of Server:")
				fmt.Println("    command-line:  curl -X GET \\ \n-H 'Content-type: application/json' \\ \n-H 'Accept: application/json' \\ \nhttp://localhost:9999/config")
				fmt.Println("")
				fmt.Println("    python: requests.get('http://localhost:9999/config')")
				fmt.Println("")
				fmt.Println("    go:  http.Get('http://localhost:9999/config')")
				fmt.Println("")
				fmt.Println("    java: HttpClient.get('http://localhost:9999/config')")
				fmt.Println("")
				fmt.Println("    lua: http.get('http://localhost:9999/config')")
				fmt.Println("")
				fmt.Println("    etc.")
				fmt.Println("")
				fmt.Println("    Once the you get the json you asked for you can use the json to")
				fmt.Println("    use in the application you want. For example Neovim, or Qtile.")
				fmt.Println("")
				fmt.Println("If you have any questions or comments, please email me at:")
				fmt.Println("    adamkali@outlook.com")
				fmt.Println("")
				fmt.Println("If you wish to contribute to the project, please visit:")
				fmt.Println("    github.com/adamkali/azalf")
				fmt.Println("Then make a pull request!")
				os.Exit(0)
			}
			if os.Args[1] == "update" {
				// Update the config file.
				conf, err = loadConfig(&conf)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
			}
			if os.Args[1] == "stop" {
				// Stop the server.
				os.Exit(0)
			}
			if os.Args[1] == "spells" {
				// print the endpoints.
				fmt.Println("")
				fmt.Println("Endpoints:")
				fmt.Println("  http://localhost:9999/config")
				fmt.Println("  http://localhost:9999/config/colors")
				fmt.Println("  http://localhost:9999/config/sizes")
				fmt.Println("  http://localhost:9999/config/fonts")
				fmt.Println("  http://localhost:9999/config/colors/{color}")
				fmt.Println("  http://localhost:9999/config/sizes/{size}")
				fmt.Println("  http://localhost:9999/config/fonts/{font}")
				fmt.Println("")
				fmt.Println("Please see your config for available {color}, {size}, and {font}")
				fmt.Println("   They will be the names of the colors, sizes, and fonts.")
				fmt.Println("")

			if os.Args[1] == "--version" || os.Args[1] == "-v" {
				// Print the version.
				fmt.Println(serverName + " " + serverVersion)
			}
			if os.Args[1] == "--license" || os.Args[1] == "-l" {
				// Print GNU GPLv3 license.
				fmt.Println(serverName + " " + serverVersion)
				fmt.Println("")
				fmt.Println("Copyright (C) 2022  Azalf")
				fmt.Println("This program is free software: you can redistribute it and/or modify")
				fmt.Println("it under the terms of the GNU General Public License as published by")
				fmt.Println("the Free Software Foundation, either version 3 of the License, or")
				fmt.Println("(at your option) any later version.")
				fmt.Println("")
				fmt.Println("This program is distributed in the hope that it will be useful,")
				fmt.Println("but WITHOUT ANY WARRANTY; without even the implied warranty of")
				fmt.Println("MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the")
				fmt.Println("GNU General Public License for more details.")
				fmt.Println("")
		} else if len (os.Args) == 0 {
			// If the program is called without an argument,
			// load the config, and then start the server.
			conf, err := loadConfig(&config)
			serveHTTP(&conf)
		} else {
			fmt.Printf("%s did not recognize that spell. Please use the --help command", serverName)
			os.Exit(1)
		}
	}
}

func () loadConfig(config *Config) (Config, err) {

	// Get the current User's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		config = Config{}
		return config, err
	}
	conf, err := = ioutil.ReadFile(fmt.Sprintf("%s/.config/azalf/.azalf.yml", homeDir))
	if err != nil {
		config = Config{}
		return config, err
	}
	err = yaml.Unmarshal(conf, &config)
	if err != nil {
		config = Config{}
		return config, err
	}
	return config, nil
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

// TODO: Finish this within arch linux
// get the GPU information by using lspci
func (h *HardwareInfo) getGPUInfo() {
	/*// Get the GPU info from lspci
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
	// now isolate the GPU info using the specified domian at the beginning of the string
	var gpuEntry string = strings.Split(strings.ToLower(gpuInfo), "vga compatible controller:")[0]
	// now run lspci -vs gpuEntry | grep -i -E "size|ram|memory|prefetchable" to get the VRAM info
	gpuInfo, err = exec.Command("lspci", "-vs", gpuEntry).Output()
	if err != nil {
		// write to the err file
		// then return
		log.Fatalf("%s could not scry lspci: %s", daemonName, err.Error())
		return
	}
	gpuInfoScan := bufio.NewScanner(bytes.NewReader(gpuInfo))
	//
	*/
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
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
	http.HandleFunc("/config", config.getConfig)
	http.HandleFunc("/config/colors", config.getConfigColor)
	http.HandleFunc("/config/fonts", config.getConfigFont)
	http.HandleFunc("/config/sizing", config.getConfigSizing)
	http.HandleFunc("/config/colors/", config.getConfigSpecificColor)
	http.HandleFunc("/config/fonts/", config.getConfigSpecificFont)
	http.HandleFunc("/config/sizing/", config.getConfigSpecificSizing)
	http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", daemonPort), nil)
}
