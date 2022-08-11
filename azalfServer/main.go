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
	serverVers = "0.2.7"

	developmentFile = "C:\\Users\\adam\\.config\\azalf\\.azalf.yml"
	debug           = false

	config  = new(Config)
	ERROR   = "error"
	WARN    = "warning"
	SUCCESS = "success"
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

// My own type for debugging the application
//    - Extensible for any type of data
//    - shows lines of code that are executed
//    - shows what type of debug message it is
type DebugStruct struct {
	FileLine  int32
	DebugType string
	Message   string
}

func CreateDebug(fileLine int32, debugType string, message string) DebugStruct {
	return DebugStruct{
		FileLine:  fileLine,
		DebugType: debugType,
		Message:   message,
	}
}

func (d *DebugStruct) String() string {
	if d.DebugType == "error" {
		// print in red
		return fmt.Sprintf("\x1b[31m%d: ERROR\t\x1b[0m%s\n", d.FileLine, d.Message)
	} else if d.DebugType == "warning" {
		return fmt.Sprintf("\x1b[33m%d: WARN\t\x1b[0m%s\n", d.FileLine, d.Message)
	} else if d.DebugType == "success" {
		// print in blue
		return fmt.Sprintf("\x1b[34m%d: SUCCESS\t\x1b[0m%s\n", d.FileLine, d.Message)
	} else {
		// print in gray
		return fmt.Sprintf("\x1b[37m%d: DEBUG\t\x1b[0m%s\n", d.FileLine, d.Message)
	}
}

func main() {

	var serverFile *os.File
	var err error

	// These are values that will be used in the server
	// if var flag bool = true; then the server will
	// after parsing the command line arguments. Then
	// the server will use the values in the config
	// file, by using the var config Config.
	var flag bool

	// Parse the command line arguments.
	// Check that the user is root and on linux
	if runtime.GOOS == "linux" {
		// create a directory for the server if it doesn't exist
		if _, err := os.Stat("/var/run/azalf"); os.IsNotExist(err) {
			os.Mkdir("/var/run/azalf", 0755)
		}
		if err != nil {
			fmt.Print(err.Error())
		}
		// If the log exists, delete it.
		if _, err := os.Stat("/var/run/azalf/azalf.log"); os.IsNotExist(err) {
			os.Remove("/var/run/azalf/azalf.log")
		}
		if err != nil {
			fmt.Print(err.Error())
		}
		serverFile, err = os.Create("/var/run/azalf/server.log")
		if err != nil {
			log.Fatal("Failed to create server log file")
		}

	} else if runtime.GOOS == "windows" {
		// DEVEL: THIS APPLICATION IS ONLY FOR LINUX
		// DEVEL: WINDOWS IS USED ONLY FOR TESTING
		// create a directory for the server log file
		err = os.MkdirAll("C:\\azalf", 0755)
		if err != nil {
			fmt.Print("Failed to create server log directory")
		}
		// If the log exists, delete it.
		if _, err := os.Stat("C:\\azalf\\azalf.log"); os.IsNotExist(err) {
			os.Remove("C:\\azalf\\azalf.log")
		}
		if err != nil {
			fmt.Print(err.Error())
		}
		serverFile, err = os.Create("C:\\azalf\\server.log")
		if err != nil {
			fmt.Print("Failed to create server log file")
		}

	} else {
		log.Fatal("This program is only supported on Linux and Windows")
		os.Exit(1)
	}

	// Specify the log file as the output for the server.
	log.SetOutput(serverFile)

	// The server should be accept args to update the config file, even while running.
	// Allow the user to call the program with an argument to update the config
	// file while the program is running.
	if len(os.Args) > 1 {
		// Make a --help command to show the user how to use the program.
		if os.Args[1] == "--help" || os.Args[1] == "-h" {
			fmt.Printf("%s %s\n", serverName, serverVers)
			fmt.Printf("%s\n", serverDesc)
			fmt.Println("")
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
			fmt.Println("  azalf logs")
			fmt.Println("  azalf --help || -h")
			fmt.Println("  azalf --version || -v")
			fmt.Println("  azalf --liscense || -l")
			fmt.Println("")
			fmt.Println("Example Use of Server:")
			fmt.Println(`    command-line:  curl -X GET \\ 
											-H 'Content-type: application/json' \\ 
											-H 'Accept: application/json' \\ 
												http://localhost:9999/config`)
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

			// check if the server is running
			// to do this we will do a get request to the server
			// on / endpoint
			resp, err := http.Get("http://localhost:9999/")
			if err != nil {
				fmt.Println("The server is not running, please start the server")
				os.Exit(1)
			}
			defer resp.Body.Close()

			if strings.Contains(strconv.Itoa(resp.StatusCode), "20") {
				flag = false
			}

			log.Printf("%s is updating his spells", serverName)

			// set var config Config using loadConfig()
			err = loadConfig(config)
			if err != nil {
				// get the home directory of the user
				homeDir, err := os.UserHomeDir()
				if err != nil {
					log.Fatal(err)
				}

				log.Printf("%s is having trouble loading his spellbook in %s >> %s", serverName, homeDir, err)
				os.Exit(1)
			}
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
		if os.Args[1] == "logs" {
			// Get the logs from serverFile.
			// Print the logs to the screen.
			log.Printf("%s is getting the logs from %s", serverName, serverFile.Name())
			logs, err := ioutil.ReadFile(serverFile.Name())
			if err != nil {
				log.Printf("%s is having trouble getting the logs from %s >> %s", serverName, serverFile.Name(), err)
				os.Exit(1)
			}
			fmt.Println(string(logs))
		}
		if os.Args[1] == "are-you-alive" {
			// check if the server is running
			// call the server / endpoint
			// and print the response to the screen.
			response, err := http.Get("http://localhost:9999/")
			if err != nil {
				log.Printf("%s is having trouble getting the response from %s >> %s", serverName, "http://localhost:9999/", err)
				os.Exit(1)
			}
			defer response.Body.Close()

			// read the response body
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Printf("%s is having trouble reading the response body from %s >> %s", serverName, "http://localhost:9999/", err)
				os.Exit(1)
			}
			fmt.Println(string(body))
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
		if os.Args[1] == "--debug" || os.Args[1] == "-d" {
			debug = true

			if debug {
				log.Printf("%s is starting in debug mode.", serverName)
				d := CreateDebug(467, "", "Started in debug mode.")
				fmt.Println(d.String())
			}

			err = loadConfig(config)
			if err != nil {
				log.Fatal(err)
				os.Exit(1)
			}

			flag = true
		}
	} else if len(os.Args) == 1 {
		// If the program is called without an argument,
		// load the config, and then start the server.

		err = loadConfig(config)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		flag = true
	} else {
		fmt.Printf("%s did not recognize that spell. Please use the --help command", serverName)
		os.Exit(1)
	}

	// If the config is loaded, start the server.
	if flag {
		// Start the server.
		// Get time now
		log.Printf("%s has been summoned.", serverName)

		serveHTTP(config)
	} else {
		log.Printf("%s was consulted but not asked to start slinging spells.", serverName)
	}

	log.Printf("%s is exiting the building...", serverName)

	// Exit the program.
	os.Exit(0)
}

func loadConfig(config *Config) error {

	// Get the current User's home directory
	// get the users os
	var homeDir string
	var err error
	var conf []byte

	if runtime.GOOS == "linux" {
		homeDir, err = os.UserHomeDir()
		conf, err = ioutil.ReadFile(
			fmt.Sprintf("%s/.config/azalf/.azalf.yml", homeDir))
		if err != nil {
			if debug {
				d := CreateDebug(524, ERROR, err.Error())
				fmt.Println(d.String())
			}
			return fmt.Errorf(`%s is having trouble reading his spellbook from
	%s\.config\azalf\.azalf.yml`, serverName, homeDir)
		}
	}
	// DEVEL: THIS APPLICATION IS NOT SUPPORTED ON WINDOWS
	// DEVEL: THIS IS USED FOR DEVELOPMENT ONLY
	if runtime.GOOS == "windows" {
		homeDir = "C:\\Users\\" + os.Getenv("USERNAME")
		conf, err = os.ReadFile(developmentFile)
		if err != nil {
			if debug {
				d := CreateDebug(539, ERROR, err.Error())
				fmt.Println(d.String())
			}
			return fmt.Errorf(`%s is having trouble reading his spellbook from
%s\.config\azalf\.azalf.yml`, serverName, homeDir)
		}
	}
	log.Printf("%s is loading his spellbook from %s",
		serverName,
		fmt.Sprintf("%s\\.config\\azalf\\.azalf.yml", homeDir))
	err = yaml.Unmarshal(conf, &config)
	if err != nil {
		if debug {
			d := CreateDebug(548, "error", err.Error())
			fmt.Println(d.String())
		}
		return err
	}

	if config.Colors.Background == "" {
		if debug {
			d := CreateDebug(552, ERROR, "Config not loaded.\n This either means there is no config file, the config file is corrupted, or the file is not being loaded properly.")
			fmt.Println(d.String())
		}
		return fmt.Errorf(`%s is missing vital spells from the spellbook.
		Ensure that the spellbook has valid colors, sizes, and fonts and try again`, serverName)
	} else {
		if runtime.GOOS == "linux" {
			if debug {
				configString, err := json.Marshal(config)
				if err != nil {
					d := CreateDebug(571, ERROR, err.Error())
					fmt.Println(d.String())
				}
				d := CreateDebug(552, SUCCESS, "Config loaded:"+string(configString))
				fmt.Println(d.String())
			}
			log.Printf("%s has found spells in %s", serverName, fmt.Sprintf("%s/.config/azalf/.azalf.yml", homeDir))
			// print the config to the log formatted to json
		} else if runtime.GOOS == "windows" {
			log.Printf("%s has found spells in %s", serverName, fmt.Sprintf("%s\\.config\\azalf\\.azalf.yml", homeDir))
			if debug {
				configString, err := json.Marshal(config)
				if err != nil {
					d := CreateDebug(584, ERROR, err.Error())
					fmt.Println(d.String())
				}
				d := CreateDebug(552, SUCCESS, "Config loaded:"+string(configString))
				fmt.Println(d.String())
			}
		}
	}
	return nil
}

// TODO: Implement getting Hardware info.
/*func (h *HardwareInfo) getHardwareInfo() {
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
}*
/

// TODO: Finish this within arch linux
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

}
*/

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
			log.Printf("AZALF could not find that %s in .azalf.yml.", r.URL.Path)
		}
	}
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
}

func serveHTTP(config *Config) {
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			// Just return a 200 OK response
			w.WriteHeader(http.StatusOK)
			confirmation := fmt.Sprintf("%s is still slingin spells!", serverName)
			w.Write([]byte(confirmation))
		})
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

	// start a TCP listener on the specified port

	// start the server
	err := http.ListenAndServe(serverPort, nil)
	if err != nil {
		log.Fatalf("%s could not start the server: %s", serverName, err.Error())
	} else {
		log.Printf("%s started the server on port %s", serverName, serverPort)
	}
}
