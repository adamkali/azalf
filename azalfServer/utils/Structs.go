package utils

import "fmt"

type Config struct {
	Colors struct {
		Background string `yaml:"background" `
		Foreground string `yaml:"foreground" `
		Normal     struct {
			Black   string `yaml:"black" `
			Red     string `yaml:"red" `
			Green   string `yaml:"green" `
			Yellow  string `yaml:"yellow"`
			Blue    string `yaml:"blue" `
			Magenta string `yaml:"magenta"`
			Cyan    string `yaml:"cyan" `
			White   string `yaml:"white"`
		} `yaml:"normal"`
		Bright struct {
			Black   string `yaml:"black"`
			Red     string `yaml:"red"`
			Green   string `yaml:"green"`
			Yellow  string `yaml:"yellow"`
			Blue    string `yaml:"blue"`
			Magenta string `yaml:"magenta"`
			Cyan    string `yaml:"cyan"`
			White   string `yaml:"white"`
		} `yaml:"bright"`
	} `yaml:"colors"`
	FontFamilies struct {
		Monospace string `yaml:"monospace"`
		SansSerif string `yaml:"sans-serif"`
		Serif     string `yaml:"serif"`
		Emoji     string `yaml:"emoji"`
	} `yaml:"font-families"`
	Sizing struct {
		FontSizes struct {
			Small   int `yaml:"small"`
			Medium  int `yaml:"medium"`
			Large   int `yaml:"large"`
			XLarge  int `yaml:"x-large"`
			XXLarge int `yaml:"xx-large"`
			Huge    int `yaml:"huge"`
		} `yaml:"font-sizes"`
		Padding      int `yaml:"padding"`
		BorderRadius int `yaml:"border-radius"`
		Margin       int `yaml:"margin"`
	} `yaml:"sizing"`
}

func (c *Config) String() string {
	return fmt.Sprintf(`
Colors:
	Background: %s
	Foreground: %s
	Normal:
		Black: %s
		Red: %s
		Green: %s
		Yellow: %s
		Blue: %s
		Magenta: %s
		Cyan: %s
		White: %s
	Bright:
		Black: %s
		Red: %s
		Green: %s
		Yellow: %s
		Blue: %s
		Magenta: %s
		Cyan: %s
		White: %s
FontFamilies:
	Monospace: %s
	SansSerif: %s
	Serif: %s
	Emoji: %s
Sizing:
	FontSizes:
		Small: %d
		Medium: %d
		Large: %d
		XLarge: %d
		XXLarge: %d
		Huge: %d
	Padding: %d
	BorderRadius: %d
	Margin: %d
	`,
		c.Colors.Background,
		c.Colors.Foreground,
		c.Colors.Normal.Black,
		c.Colors.Normal.Red,
		c.Colors.Normal.Green,
		c.Colors.Normal.Yellow,
		c.Colors.Normal.Blue,
		c.Colors.Normal.Magenta,
		c.Colors.Normal.Cyan,
		c.Colors.Normal.White,
		c.Colors.Bright.Black,
		c.Colors.Bright.Red,
		c.Colors.Bright.Green,
		c.Colors.Bright.Yellow,
		c.Colors.Bright.Blue,
		c.Colors.Bright.Magenta,
		c.Colors.Bright.Cyan,
		c.Colors.Bright.White,
		c.FontFamilies.Monospace,
		c.FontFamilies.SansSerif,
		c.FontFamilies.Serif,
		c.FontFamilies.Emoji,
		c.Sizing.FontSizes.Small,
		c.Sizing.FontSizes.Medium,
		c.Sizing.FontSizes.Large,
		c.Sizing.FontSizes.XLarge,
		c.Sizing.FontSizes.XXLarge,
		c.Sizing.FontSizes.Huge,
		c.Sizing.Padding,
		c.Sizing.BorderRadius,
		c.Sizing.Margin,
	)
}
