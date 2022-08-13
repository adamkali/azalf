package endpoints

type HardwareInfo struct {
	CPU struct {
		Model    string `yaml:"model"; json:"model"`
		Cores    int    `yaml:'cores'; json:'cores'`
		Speed    int    `yaml:"speed"; json:speed"`
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
