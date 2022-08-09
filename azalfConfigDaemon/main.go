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
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	serverName = "AZALF"
	serverDesc = "Adam's Zillenial Arch Linux Flavor Server"
	serverPort = ":9999"
	serverVers = "0.2.4"
)

// Use ../.config/azalf/.azalf.yaml as a model for what this struct
// should look like.
type Config struct {
	Colors struct {
		Background string `yaml:"background" json:"background"`
		Foreground string `yaml:"foreground" json:"foreground"`
		Normal     struct {
			Black   string `yaml:"black" json:"black"`
			Red     string `yaml:"red" json:"red"`
			Green   string `yaml:"green" json:"green"`
			Yellow  string `yaml:"yellow" json:"yellow"`
			Blue    string `yaml:"blue" json:"blue"`
			Magenta string `yaml:"magenta" json:"magenta"`
			Cyan    string `yaml:"cyan" json:"cyan"`
			White   string `yaml:"white" json:"white"`
		} `yaml:"normal" json:"normal"`
		Bright struct {
			Black   string `yaml:"black" json:"black"`
			Red     string `yaml:"red" json:"red"`
			Green   string `yaml:"green" json:"green"`
			Yellow  string `yaml:"yellow" json:"yellow"`
			Blue    string `yaml:"blue" json:"blue"`
			Magenta string `yaml:"magenta" json:"magenta"`
			Cyan    string `yaml:"cyan" json:"cyan"`
			White   string `yaml:"white" json:"white"`
		} `yaml:"bright" json:"bright"`
	} `yaml:"colors" json:"colors"`
	FontFamilies struct {
		Monospace string `yaml:"monospace" json:"monospace"`
		SansSerif string `yaml:"sans-serif" json:"sans-serif"`
		Serif     string `yaml:"serif" json:"serif"`
		Emoji     string `yaml:"emoji" json:"emoji"`
	} `yaml:"font-families" json:"font-families"`
	Sizing struct {
		FontSizes struct {
			Small   int `yaml:"small" json:"small"`
			Medium  int `yaml:"medium" json:"medium"`
			Large   int `yaml:"large" json:"large"`
			XLarge  int `yaml:"x-large" json:"x-large"`
			XXLarge int `yaml:"xx-large" json:"xx-large"`
			Huge    int `yaml:"huge" json:"huge"`
		} `yaml:"font-sizes" json:"font-sizes"`
		Padding      int `yaml:"padding" json:"padding"`
		BorderRadius int `yaml:"border-radius" json:"border-radius"`
		Margin       int `yaml:"margin" json:"margin"`
	} `yaml:"sizing" json:"sizing"`
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

	var serverFile *os.File
	var err error

	// Check that the user is root and on linux
	if runtime.GOOS == "linux" {
		// create a directory for the server if it doesn't exist
		if _, err := os.Stat("/var/run/azalf"); os.IsNotExist(err) {
			os.Mkdir("/var/run/azalf", 0755)
		}
		serverFile, err = os.Create("/var/run/azalf/server.log")
		if err != nil {
			log.Fatal("Failed to create server log file")
		}

	} else if runtime.GOOS == "windows" {
		// create a directory for the server log file
		err = os.MkdirAll("C:\\azalf", 0755)
		if err != nil {
			log.Fatal("Failed to create server log directory")
		}
		serverFile, err = os.Create("C:\\azalf\\server.log")
		if err != nil {
			log.Fatal("Failed to create server log file")
		}

	} else {
		log.Fatal("This program is only supported on Linux and Windows")
		os.Exit(1)
	}

	// Specify the log file as the output for the server.
	log.SetOutput(serverFile)

	// Loop until the user calls `azalf stop`

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
			conf, err := loadConfig()
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}
			var wg sync.WaitGroup
			wg.Add(1)
			run(&conf, &wg)
			wg.Wait()
		}
		if os.Args[1] == "stop" {
			// Stop the server.
			os.Exit(0)
		}
		if os.Args[1] == "spells" {

			fmt.Println("")
			fmt.Println("Endpoints:")
			fmt.Println("  http://localhost:9999/config")
			fmt.Println("  http://localhost:9999/config/colors")
			fmt.Println("  http://localhost:9999/config/fonts")
			fmt.Println("  http://localhost:9999/config/sizing")
			fmt.Println("  http://localhost:9999/config/colors/normal")
			fmt.Println("  http://localhost:9999/config/colors/bright")
			fmt.Println("  http://localhost:9999/config/colors/background")
			fmt.Println("  http://localhost:9999/config/colors/foreground")
			fmt.Println("  http://localhost:9999/config/colors/normal/{color}")
			fmt.Println("  http://localhost:9999/config/colors/bright/{color}")
			fmt.Println("  http://localhost:9999/config/fonts/{font}")
			fmt.Println("  http://localhost:9999/config/sizing/fonts/{size}")
			fmt.Println("")
			fmt.Println("Please see your config for available {color}, {size}, and {font}")
			fmt.Println("   They will be the names of the specific colors, sizes, and fonts.")
			fmt.Println("")
			os.Exit(0)
		}
		if os.Args[1] == "--version" || os.Args[1] == "-v" {
			// Print the version.
			fmt.Println(serverName + " " + serverVers)
			os.Exit(0)
		}
		if os.Args[1] == "--license" || os.Args[1] == "-l" {
			// Print GNU GPLv3 license.
			fmt.Println(serverName + " " + serverVers)
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
			os.Exit(0)
		}
	} else if len(os.Args) == 1 {
		// If the program is called without an argument,
		// load the config, and then start the server.
		conf, err := loadConfig()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// make a waitgroup to wait for the server to finish.
		var wg sync.WaitGroup
		wg.Add(1)
		run(&conf, &wg)
		wg.Wait()
	} else {
		fmt.Printf("%s did not recognize that spell. Please use the --help command", serverName)
		os.Exit(1)
	}

}

func run(conf *Config, wg *sync.WaitGroup) {
	// Run the server.
	// Check if the server is running on port 9999.
	// If it is, then exit. Then start the server.
	// If it is not, then start the server.
	// Wait for the server to finish.
	// then call wg.Done()
	var err error

	go func() {
		for {
			// check if the server is running on port 9999.
			// if it is, then exit.
			err = serveHTTP(conf)
			if err != nil {
				wg.Done()
				break
			}
		}
	}()

}

func loadConfig() (Config, error) {
	var newConfig Config

	// Get the current User's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		newConfig = Config{}
		return newConfig, err
	}
	conf, err := ioutil.ReadFile(fmt.Sprintf("%s/.config/azalf/.azalf.yml", homeDir))
	if err != nil {
		newConfig = Config{}
		return newConfig, err
	}
	err = yaml.Unmarshal(conf, &newConfig)
	if err != nil {
		newConfig = Config{}
		return newConfig, err
	}
	return newConfig, nil
}

// TODO: Finish this
func (h *HardwareInfo) getHardwareInfo() {
	// Check that the os is linux
	// if not, return and write Fatalf

	if runtime.GOOS != "linux" {
		log.Fatalf("%s can only do his magic on linux", serverName)
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
		log.Fatalf("%s could not open /proc/cpuinfo: ", serverName, err)
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
			log.Fatalf("%s could not convert the cpu info: ", serverName, err)
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
		log.Fatalf("%s could not scry /proc/stat: %s", serverName, err.Error())
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
			log.Fatalf("%s could not convert the cpu info: ", serverName, err)
			break
		}
		totalTime += c
	}

	idle, err := strconv.Atoi(cpuUsageInfo[4])
	if err != nil {
		log.Fatalf("%s could not convert the cpu info: ", serverName, err)
	}
	iowait, err := strconv.Atoi(cpuUsageInfo[5])
	if err != nil {
		log.Fatalf("%s could not convert the cpu info: ", serverName, err)
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
		log.Fatalf("%s could not scry lspci: %s", serverName, err.Error())
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
		log.Fatalf("%s could not scry lspci: %s", serverName, err.Error())
		return
	}
	gpuInfoScan := bufio.NewScanner(bytes.NewReader(gpuInfo))
	//
	*/
}

func (config *Config) GetNormalColor(color string) string {
	var returnColor string
	switch color {
	case "black":
		returnColor = config.Colors.Normal.Black
	case "red":
		returnColor = config.Colors.Normal.Red
	case "green":
		returnColor = config.Colors.Normal.Green
	case "yellow":
		returnColor = config.Colors.Normal.Yellow
	case "blue":
		returnColor = config.Colors.Normal.Blue
	case "magenta":
		returnColor = config.Colors.Normal.Magenta
	case "cyan":
		returnColor = config.Colors.Normal.Cyan
	case "white":
		returnColor = config.Colors.Normal.White
	default:
		log.Fatalf("%s could not find the color: %s", serverName, color)
	}
	return returnColor
}

func (config *Config) GetBrightColor(color string) string {
	var returnColor string
	switch color {
	case "black":
		returnColor = config.Colors.Bright.Black
	case "red":
		returnColor = config.Colors.Bright.Red
	case "green":
		returnColor = config.Colors.Bright.Green
	case "yellow":
		returnColor = config.Colors.Bright.Yellow
	case "blue":
		returnColor = config.Colors.Bright.Blue
	case "magenta":
		returnColor = config.Colors.Bright.Magenta
	case "cyan":
		returnColor = config.Colors.Bright.Cyan
	case "white":
		returnColor = config.Colors.Bright.White
	default:
		log.Fatalf("%s could not find the color: %s", serverName, color)
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
		returnFont = config.FontFamilies.Monospace
	case "sans-serif":
		returnFont = config.FontFamilies.SansSerif
	case "serif":
		returnFont = config.FontFamilies.Serif
	case "emoji":
		returnFont = config.FontFamilies.Emoji
	default:
		log.Fatalf("%s could not find the font: %s", serverName, font)
	}
	return returnFont
}

func (config *Config) GetFontSizes(size string) string {
	var returnSize string
	switch size {
	case "small":
		// convert the size to a string
		// then return the string
		returnSize = strconv.Itoa(config.Sizing.FontSizes.Small)
	case "medium":
		returnSize = strconv.Itoa(config.Sizing.FontSizes.Medium)
	case "large":
		returnSize = strconv.Itoa(config.Sizing.FontSizes.Large)
	case "x-large":
		returnSize = strconv.Itoa(config.Sizing.FontSizes.XLarge)
	case "xx-large":
		returnSize = strconv.Itoa(config.Sizing.FontSizes.XXLarge)
	case "Huge":
		returnSize = strconv.Itoa(config.Sizing.FontSizes.Huge)
	default:
		log.Fatalf("%s could not find the size: %s", serverName, size)
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
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
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
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
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
	err := json.NewEncoder(w).Encode(&config.FontFamilies)

	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
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
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigNormalColors(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Colors.Normal)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigBrightColors(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Colors.Bright)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigSpecificNormalColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the path parameter for the color
	// the request path is /config/colors/normal/{:color}
	// the path parameter is {:color}
	color := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("http://localhost:%s/config/colors/bright/", serverPort))
	res := config.GetNormalColor(color)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigBackgroundColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Colors.Background)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigForegroundColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Colors.Foreground)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigSpecificBrightColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the path parameter for the color
	// the request path is /config/colors/bright/{:color}
	// the path parameter is {:color}
	color := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("http://localhost:%s/config/colors/bright/", serverPort))
	res := config.GetBrightColor(color)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
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
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) GetConfigSpecificFontSize(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sizing := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("http://localhost:%s/config/sizing/fonts/", serverPort))
	res := config.GetFontSizes(sizing)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigFontSizes(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Sizing.FontSizes)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigPadding(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Sizing.Padding)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigMargin(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Sizing.Margin)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func (config *Config) getConfigBorderRadius(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&config.Sizing.BorderRadius)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", serverName, err.Error())
	}

	return
}

func serveHTTP(config *Config) error {
	http.HandleFunc("/config", config.getConfigObject)
	http.HandleFunc("/config/colors", config.getConfigColor)
	http.HandleFunc("/config/fonts", config.getConfigFont)
	http.HandleFunc("/config/sizing", config.getConfigSizing)
	http.HandleFunc("/config/colors/normal", config.getConfigNormalColors)
	http.HandleFunc("/config/colors/bright", config.getConfigBrightColors)
	http.HandleFunc("/config/colors/background", config.getConfigBackgroundColor)
	http.HandleFunc("/config/colors/foreground", config.getConfigForegroundColor)
	http.HandleFunc("/config/colors/normal/", config.getConfigSpecificNormalColor)
	http.HandleFunc("/config/colors/bright/", config.getConfigSpecificBrightColor)
	http.HandleFunc("/config/fonts/", config.getConfigSpecificFont)
	http.HandleFunc("/config/sizing/fonts", config.getConfigFontSizes)
	http.HandleFunc("/config/sizing/fonts/", config.GetConfigSpecificFontSize)
	http.HandleFunc("/config/sizing/padding", config.getConfigPadding)
	http.HandleFunc("/config/sizing/margin", config.getConfigMargin)
	http.HandleFunc("/config/sizing/borderRadius", config.getConfigBorderRadius)

	err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", serverPort), nil)
	return err
}
