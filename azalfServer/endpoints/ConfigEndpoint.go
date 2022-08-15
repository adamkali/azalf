package endpoints

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

	"azalf/utils"
)

type ConfigHandler struct {
	config *utils.Config
}

// make a new Config Handler and load the config
func NewConfigHandler() *ConfigHandler {
	// Load the config using the LoadConfig function
	// and check if there is an error
	err := LoadConfig(utils.AzalfConfig)
	if err != nil {
		log.Fatalf("%s could not load the config: %s", utils.ServerName, err.Error())
		// make a new debug insstance if the server is in debug mode
		if utils.Debug {
			debug := utils.CreateDebug(5, "ERROR", fmt.Sprintf("%s could not load the config: %s", utils.ServerName, err.Error()))
			fmt.Println(debug.String())
		}
		os.Exit(1)
	}
	conf := *utils.AzalfConfig
	return &ConfigHandler{config: &conf}
}

func LoadConfig(config *utils.Config) error {

	// Get the current User's home directory
	// get the users os
	var homeDir string
	var err error
	var conf []byte

	// get the current users home directory like this:
	// homeDir should end up like this: /home/username
	// force to get the user's home directory even when
	// running root
	if runtime.GOOS == "linux" {
		// check if the user is root
		userCurrent := os.Getenv("SUDO_USER")
		homeDir = "/home/" + userCurrent
		conf, err = ioutil.ReadFile(
			fmt.Sprintf("%s/.config/azalf/azalf.yml", homeDir))
		if err != nil {
			if utils.Debug {
				d := utils.CreateDebug(59, utils.ERROR, err.Error())
				fmt.Println(d.String())
				fmt.Printf("%s is the current user home directory\n", homeDir)
			}
			return fmt.Errorf(`%s is having trouble reading his spellbook from
	%s\.config\azalf\azalf.yml`, utils.ServerName, homeDir)
		}
	}
	// DEVEL: THIS APPLICATION IS NOT SUPPORTED ON WINDOWS
	// DEVEL: THIS IS USED FOR DEVELOPMENT ONLY
	if runtime.GOOS == "windows" {
		homeDir = "C:\\Users\\" + os.Getenv("USERNAME")
		conf, err = os.ReadFile(utils.DevelopmentFile)
		if err != nil {
			if utils.Debug {
				d := utils.CreateDebug(75, utils.ERROR, err.Error())
				fmt.Println(d.String())
			}
			return fmt.Errorf(`%s is having trouble reading his spellbook from
%s\.config\azalf\.azalf.yml`, utils.ServerName, homeDir)
		}
	}
	log.Printf("%s is loading his spellbook from %s",
		utils.ServerName,
		fmt.Sprintf("%s\\.config\\azalf\\azalf.yml", homeDir))

	if utils.Debug {
		d := utils.CreateDebug(75, "", "Loading config \n"+string(conf))
		fmt.Println(d.String())
	}

	// Attempt to unmarshal the config file. and then store it into the
	// config struct initialized in the utils function.
	err = yaml.Unmarshal(conf, config)

	if err != nil {
		if utils.Debug {
			d := utils.CreateDebug(91, "error", err.Error())
			fmt.Println(d.String())
		}
		return err
	}

	if utils.AzalfConfig.Colors.Background == "" {
		if utils.Debug {
			d := utils.CreateDebug(91, utils.ERROR, "Config not loaded.\n This either means there is no config file, the config file is corrupted, or the file is not being loaded properly.")
			confString, err := json.Marshal(config)
			if err != nil {
				log.Fatalf("%s could not encode the utils.AzalfConfig object: %s", utils.ServerName, err.Error())
			}
			fmt.Println(d.String() + "\n" + string(confString))
		}
		return fmt.Errorf(`%s is missing vital spells from the spellbook.
		Ensure that the spellbook has valid colors, sizes, and fonts and try again`, utils.ServerName)
	} else {
		if runtime.GOOS == "linux" {
			if utils.Debug {
				configString, err := json.Marshal(&config)
				if err != nil {
					d := utils.CreateDebug(111, utils.ERROR, err.Error())
					fmt.Println(d.String())
				}
				d := utils.CreateDebug(91, utils.SUCCESS, "Config loaded:"+string(configString))
				fmt.Println(d.String())
			}
			log.Printf("%s has found spells in %s", utils.ServerName, fmt.Sprintf("%s/.config/azalf/.azalf.yml", homeDir))
			// print the config to the log formatted to json
		} else if runtime.GOOS == "windows" {
			log.Printf("%s has found spells in %s", utils.ServerName, fmt.Sprintf("%s\\.config\\azalf\\.azalf.yml", homeDir))
			if utils.Debug {
				configString, err := json.Marshal(&config)
				if err != nil {
					d := utils.CreateDebug(124, utils.ERROR, err.Error())
					fmt.Println(d.String())
				}
				d := utils.CreateDebug(91, utils.SUCCESS, "Config loaded:"+string(configString))
				fmt.Println(d.String())
			}
		}
	}
	return nil
}

//func (h *userHandler) post(w http.ResponseWriter, r *http.Request) {

func (h *ConfigHandler) getNormalColor(color string) string {
	var returnColor string

	config := h.config

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
		log.Fatalf("%s could not find the color: %s", utils.ServerName, color)
	}
	return returnColor
}

func (h *ConfigHandler) getBrightColor(color string) string {
	var returnColor string

	config := h.config

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
		log.Fatalf("%s could not find the color: %s", utils.ServerName, color)
	}
	return returnColor
}

func (h *ConfigHandler) getFont(font string) string {
	var returnFont string

	config := h.config

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
		log.Fatalf("%s could not find the font: %s", utils.ServerName, font)
	}
	return returnFont
}

func (h *ConfigHandler) getFontSizes(size string) string {
	var returnSize string

	config := h.config

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
		log.Fatalf("%s could not find the size: %s", utils.ServerName, size)
	}
	return returnSize
}

func (c *ConfigHandler) GetConfigObject(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the config from the handler
	config := c.config

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config)
	if err != nil {
		log.Fatalf("%s could not encode the utils.AzalfConfig object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Colors)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigFont(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a Get
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.FontFamilies)

	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigSizing(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Sizing)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigNormalColors(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Colors.Normal)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigBrightColors(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a Get
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Colors.Bright)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigSpecificNormalColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the path parameter for the color
	// the request path is /config/colors/normal/{:color}
	// the path parameter is {:color}
	color := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("http://localhost:%s/config/colors/bright/", utils.ServerPort))
	res := c.getNormalColor(color)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigBackgroundColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Colors.Background)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigForegroundColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Colors.Foreground)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigSpecificBrightColor(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// get the path parameter for the color
	// the request path is /config/colors/bright/{:color}
	// the path parameter is {:color}
	color := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("http://localhost:%s/config/colors/bright/", utils.ServerPort))
	res := c.getBrightColor(color)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigSpecificFont(w http.ResponseWriter, r *http.Request) {
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
	specificFont = c.getFont(specificFont)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(specificFont)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigSpecificFontSize(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sizing := strings.TrimPrefix(r.URL.Path, fmt.Sprintf("http://localhost:%s/config/sizing/fonts/", utils.ServerPort))
	res := c.getFontSizes(sizing)

	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigFontSizes(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Sizing.FontSizes)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigPadding(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Sizing.Padding)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigMargin(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Sizing.Margin)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}

func (c *ConfigHandler) GetConfigBorderRadius(w http.ResponseWriter, r *http.Request) {
	// ensure the request is a GET
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	config := c.config
	// set up the response header for a json object
	w.Header().Set("Content-Type", "application/json")
	// write the json object to the response
	err := json.NewEncoder(w).Encode(config.Sizing.BorderRadius)
	if err != nil {
		log.Fatalf("%s could not encode the config object: %s", utils.ServerName, err.Error())
	}
}
