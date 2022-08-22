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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"

	"azalf/endpoints"
	"azalf/utils"
)

// Use ../.config/azalf/.azalf.yaml as a model for what this struct
// should look like.

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
		// get the current user home directory
		userName, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		user := userName.HomeDir
		// make an entry point for the server in the home directory
		serverFile, err = os.OpenFile(user+"/.local/azalf/server.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal("Failed to create server log file")
		}
		/* TODO:
		// also create a file for authentication if it doesn't exist
		authDir := user + "/.local/azalf/auth"
		authFile := authDir + "/auth.lck"
		if _, err := os.Stat(authDir); os.IsNotExist(err) {
			os.Mkdir(authDir, 0755)
		}
		if _, err := os.Stat(authFile); os.IsNotExist(err) {
			os.Create(authFile)
			// ask the user for a password to use for authentication
			fmt.Print("Enter a password to use for authentication: ")
			var password string
			fmt.Scanln(&password)
			// write the password to the file
			ioutil.WriteFile(authFile, []byte(password), 0644)
		}

		// print a message to the user that the server can change the password
		// by making a post request to the server with the password as the body
		// of the request. { password: "<password>" }

		*/
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

		/* TODO:
		// also create a file for authentication if it doesn't exist
		authDir := "C:\\azalf\\auth"
		authFile := authDir + "\\auth.lck"
		if _, err := os.Stat(authDir); os.IsNotExist(err) {
			os.Mkdir(authDir, 0755)
		}
		if _, err := os.Stat(authFile); os.IsNotExist(err) {
			os.Create(authFile)
			// ask the user for a password to use for authentication
			fmt.Print("Enter a password to use for authentication: ")
			var password string
			fmt.Scanln(&password)
			// write the password to the file
			ioutil.WriteFile(authFile, []byte(password), 0644)
		}
		*/
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
			fmt.Printf("%s %s\n", utils.ServerName, utils.ServerVers)
			fmt.Printf("%s\n", utils.ServerDesc)
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
			log.Printf("%s is getting the logs from %s", utils.ServerName, serverFile.Name())
			logs, err := ioutil.ReadFile(serverFile.Name())
			if err != nil {
				log.Printf("%s is having trouble getting the logs from %s >> %s", utils.ServerName, serverFile.Name(), err)
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
				log.Printf("%s is having trouble getting the response from %s >> %s", utils.ServerName, "http://localhost:9999/", err)
				os.Exit(1)
			}
			defer response.Body.Close()

			// read the response body
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Printf("%s is having trouble reading the response body from %s >> %s", utils.ServerName, "http://localhost:9999/", err)
				os.Exit(1)
			}
			fmt.Println(string(body))
		}
		if os.Args[1] == "--version" || os.Args[1] == "-v" {
			// Print the version.
			fmt.Println(utils.ServerName + " " + utils.ServerVers)
			os.Exit(0)
		}
		if os.Args[1] == "--license" || os.Args[1] == "-l" {
			// Print GNU GPLv3 license.
			fmt.Println(utils.ServerName + " " + utils.ServerVers)
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
			utils.Debug = true

			if utils.Debug {
				log.Printf("%s is starting in Debug mode.", utils.ServerName)
				d := utils.CreateDebug(467, "", "Started in Debug mode.")
				fmt.Println(d.String())
			}

			flag = true
		}
	} else if len(os.Args) == 1 {
		flag = true
	} else {
		fmt.Printf("%s did not recognize that spell. Please use the --help command", utils.ServerName)
		os.Exit(1)
	}

	// If the config is loaded, start the server.
	if flag {
		// Start the server.
		// Get time now
		serveHTTP(utils.AzalfConfig)
	} else {
		log.Printf("%s was consulted but not asked to start slinging spells.", utils.ServerName)
	}

	log.Printf("%s is exiting the building...", utils.ServerName)

	// Exit the program.
	os.Exit(0)
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

func serveHTTP(config *utils.Config) {
	// Get the hanfler for the Config
	h := endpoints.NewConfigHandler()

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			// Just return a 200 OK response
			w.WriteHeader(http.StatusOK)
			confirmation := fmt.Sprintf("%s is still slingin spells!", utils.ServerName)
			w.Write([]byte(confirmation))
		})
	http.HandleFunc("/config", h.GetConfigObject)
	http.HandleFunc("/config/colors", h.GetConfigColor)
	http.HandleFunc("/config/fonts", h.GetConfigFont)
	http.HandleFunc("/config/sizing", h.GetConfigSizing)
	http.HandleFunc("/config/colors/normal", h.GetConfigNormalColors)
	http.HandleFunc("/config/colors/bright", h.GetConfigBrightColors)
	http.HandleFunc("/config/colors/background", h.GetConfigBackgroundColor)
	http.HandleFunc("/config/colors/foreground", h.GetConfigForegroundColor)
	http.HandleFunc("/config/colors/normal/", h.GetConfigSpecificNormalColor)
	http.HandleFunc("/config/colors/bright/", h.GetConfigSpecificBrightColor)
	http.HandleFunc("/config/fonts/", h.GetConfigSpecificFont)
	http.HandleFunc("/config/sizing/fonts", h.GetConfigFontSizes)
	http.HandleFunc("/config/sizing/fonts/", h.GetConfigSpecificFontSize)
	http.HandleFunc("/config/sizing/padding", h.GetConfigPadding)
	http.HandleFunc("/config/sizing/margin", h.GetConfigMargin)
	http.HandleFunc("/config/sizing/borderRadius", h.GetConfigBorderRadius)

	// start a TCP listener on the specified port

	// start the server
	err := http.ListenAndServe(utils.ServerPort, nil)
	if err != nil {
		log.Fatalf("%s could not start the server: %s", utils.ServerName, err.Error())
	} else {
		log.Printf("%s started the server on port %s", utils.ServerName, utils.ServerPort)
	}
}
