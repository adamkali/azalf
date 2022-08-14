package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"azalf/utils"
)

func LoadConfig() error {

	// Get the current User's home directory
	// get the users os
	var homeDir string
	var err error
	var conf []byte

	utils.AzalfConfig = new(utils.Config)

	// get the current users home directory like this:
	// homeDir should end up like this: /home/username
	// force to get the user's home directory even when
	// running root
	if runtime.GOOS == "linux" {
		// check if the user is root
		if os.Getuid() == 0 {
			//
			*
		userCurrent, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		homeDir = userCurrent.HomeDir
		conf, err = ioutil.ReadFile(
			fmt.Sprintf("%s/.config/azalf/.azalf.yml", homeDir))
		if err != nil {
			if utils.Debug {
				d := utils.CreateDebug(35, utils.ERROR, err.Error())
				fmt.Println(d.String())
				fmt.Printf("%s is the current user home directory\n", homeDir)
			}
			return fmt.Errorf(`%s is having trouble reading his spellbook from
	%s\.config\azalf\.azalf.yml`, utils.ServerName, homeDir)
		}
	}
	// DEVEL: THIS APPLICATION IS NOT SUPPORTED ON WINDOWS
	// DEVEL: THIS IS USED FOR DEVELOPMENT ONLY
	if runtime.GOOS == "windows" {
		homeDir = "C:\\Users\\" + os.Getenv("USERNAME")
		conf, err = os.ReadFile(utils.DevelopmentFile)
		if err != nil {
			if utils.Debug {
				d := utils.CreateDebug(50, utils.ERROR, err.Error())
				fmt.Println(d.String())
			}
			return fmt.Errorf(`%s is having trouble reading his spellbook from
%s\.config\azalf\.azalf.yml`, utils.ServerName, homeDir)
		}
	}
	log.Printf("%s is loading his spellbook from %s",
		utils.ServerName,
		fmt.Sprintf("%s\\.config\\azalf\\.azalf.yml", homeDir))

	// Attempt to unmarshal the config file. and then store it into the
	// config struct initialized in the utils function.
	err = yaml.Unmarshal(conf, &utils.AzalfConfig)

	if err != nil {
		if utils.Debug {
			d := utils.CreateDebug(548, "error", err.Error())
			fmt.Println(d.String())
		}
		return err
	}

	if utils.AzalfConfig.Colors.Background == "" {
		if utils.Debug {
			d := utils.CreateDebug(81, utils.ERROR, "Config not loaded.\n This either means there is no config file, the config file is corrupted, or the file is not being loaded properly.")
			fmt.Printf("%s; \nThe following string was what was loaded as an input.\n%s\nPlease review the above and continue.", d.String(), string(utils.AzalfConfig.String()))
		}
		return fmt.Errorf(`%s is missing vital spells from the spellbook.
		Ensure that the spellbook has valid colors, sizes, and fonts and try again`, utils.ServerName)
	} else {
		if runtime.GOOS == "linux" {
			if utils.Debug {
				configString, err := json.Marshal(&utils.AzalfConfig)
				if err != nil {
					d := utils.CreateDebug(86, utils.ERROR, err.Error())
					fmt.Println(d.String())
				}
				d := utils.CreateDebug(66, utils.SUCCESS, "Config loaded:"+string(configString))
				fmt.Println(d.String())
			}
			log.Printf("%s has found spells in %s", utils.ServerName, fmt.Sprintf("%s/.config/azalf/.azalf.yml", homeDir))
			// print the config to the log formatted to json
		} else if runtime.GOOS == "windows" {
			log.Printf("%s has found spells in %s", utils.ServerName, fmt.Sprintf("%s\\.config\\azalf\\.azalf.yml", homeDir))
			if utils.Debug {
				configString, err := json.Marshal(utils.AzalfConfig)
				if err != nil {
					d := utils.CreateDebug(99, utils.ERROR, err.Error())
					fmt.Println(d.String())
				}
				d := utils.CreateDebug(86, utils.SUCCESS, "Config loaded:"+string(configString))
				fmt.Println(d.String())
			}
		}
	}
	return nil
}

func GetNormalColor(color string) string {
	var returnColor string
	switch color {
	case "black":
		returnColor = utils.AzalfConfig.Colors.Normal.Black
	case "red":
		returnColor = utils.AzalfConfig.Colors.Normal.Red
	case "green":
		returnColor = utils.AzalfConfig.Colors.Normal.Green
	case "yellow":
		returnColor = utils.AzalfConfig.Colors.Normal.Yellow
	case "blue":
		returnColor = utils.AzalfConfig.Colors.Normal.Blue
	case "magenta":
		returnColor = utils.AzalfConfig.Colors.Normal.Magenta
	case "cyan":
		returnColor = utils.AzalfConfig.Colors.Normal.Cyan
	case "white":
		returnColor = utils.AzalfConfig.Colors.Normal.White
	default:
		log.Fatalf("%s could not find the color: %s", utils.ServerName, color)
	}
	return returnColor
}

func GetBrightColor(color string) string {
	var returnColor string
	switch color {
	case "black":
		returnColor = utils.AzalfConfig.Colors.Bright.Black
	case "red":
		returnColor = utils.AzalfConfig.Colors.Bright.Red
	case "green":
		returnColor = utils.AzalfConfig.Colors.Bright.Green
	case "yellow":
		returnColor = utils.AzalfConfig.Colors.Bright.Yellow
	case "blue":
		returnColor = utils.AzalfConfig.Colors.Bright.Blue
	case "magenta":
		returnColor = utils.AzalfConfig.Colors.Bright.Magenta
	case "cyan":
		returnColor = utils.AzalfConfig.Colors.Bright.Cyan
	case "white":
		returnColor = utils.AzalfConfig.Colors.Bright.White
	default:
		log.Fatalf("%s could not find the color: %s", utils.ServerName, color)
	}
	return returnColor
}

func GetFont(font string) string {
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
		returnFont = utils.AzalfConfig.FontFamilies.Monospace
	case "sans-serif":
		returnFont = utils.AzalfConfig.FontFamilies.SansSerif
	case "serif":
		returnFont = utils.AzalfConfig.FontFamilies.Serif
	case "emoji":
		returnFont = utils.AzalfConfig.FontFamilies.Emoji
	default:
		log.Fatalf("%s could not find the font: %s", utils.ServerName, font)
	}
	return returnFont
}

func GetFontSizes(size string) string {
	var returnSize string
	switch size {
	case "small":
		// convert the size to a string
		// then return the string
		returnSize = strconv.Itoa(utils.AzalfConfig.Sizing.FontSizes.Small)
	case "medium":
		returnSize = strconv.Itoa(utils.AzalfConfig.Sizing.FontSizes.Medium)
	case "large":
		returnSize = strconv.Itoa(utils.AzalfConfig.Sizing.FontSizes.Large)
	case "x-large":
		returnSize = strconv.Itoa(utils.AzalfConfig.Sizing.FontSizes.XLarge)
	case "xx-large":
		returnSize = strconv.Itoa(utils.AzalfConfig.Sizing.FontSizes.XXLarge)
	case "Huge":
		returnSize = strconv.Itoa(utils.AzalfConfig.Sizing.FontSizes.Huge)
	default:
		log.Fatalf("%s could not find the size: %s", utils.ServerName, size)
	}
	return returnSize
}

func GetConfigObject(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig)
	if err != nil {
		log.Fatalf("%s could not encode the utils.AzalfConfig object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Colors)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigFont(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a Get
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.FontFamilies)

	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigSizing(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Sizing)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigNormalColors(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Colors.Normal)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigBrightColors(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a Get
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Colors.Bright)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigSpecificNormalColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the path parameter for the color
	// the request path is /config/colors/normal/{:color}
	// the path parameter is {:color}
	color := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("http://localhost:%s/config/colors/bright/", utils.ServerPort))
	res := GetNormalColor(color)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigBackgroundColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Colors.Background)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigForegroundColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Colors.Foreground)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigSpecificBrightColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the path parameter for the color
	// the request path is /config/colors/bright/{:color}
	// the path parameter is {:color}
	color := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("http://localhost:%s/config/colors/bright/", utils.ServerPort))
	res := GetBrightColor(color)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigSpecificFont(w http.ResponseWriter, r *http.Request) {
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
	specificFont = GetFont(specificFont)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&specificFont)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigSpecificFontSize(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sizing := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("http://localhost:%s/config/sizing/fonts/", utils.ServerPort))
	res := GetFontSizes(sizing)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigFontSizes(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Sizing.FontSizes)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigPadding(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Sizing.Padding)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigMargin(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Sizing.Margin)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func GetConfigBorderRadius(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(&utils.AzalfConfig.Sizing.BorderRadius)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}
