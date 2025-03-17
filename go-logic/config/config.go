package config

import (
    "os"
    "encoding/json"
    "log"
)

type Misc struct {
    AllowMisc        bool `json:"allow_misc"`
    MoveAllToMisc    bool `json:"move_all_to_misc"`
    ExcludeNamesFromMisc  []string `json:"exclude_names_from_misc"`
    ExcludeTypesFromMisc  []string `json:"exclude_types_from_misc"`
}

type Desktop struct {
	FolderTypes     []string `json:"folder_types"`
    ExcludeTypes    []string `json:"exclude_types"`
    ExcludeNames    []string  `json:"exclude_names"`
	SortScreenshots bool     `json:"sort_screenshots"`
    Misc            Misc     `json:"misc"`
}

type Downloads struct {
	FolderTypes []string `json:"folder_types"`
    ExcludeTypes    []string `json:"exclude_types"`
    ExcludeNames    []string  `json:"exclude_names"`
	AllowMisc       bool     `json:"allow_misc"`
}

type Crontab struct {
	Minute       string   `json:"minute"`
	Hour         string   `json:"hour"`
	MonthDay     string   `json:"month_day"`
	Month        string   `json:"month"`
	DayWeek      string   `json:"day_week"`
}

type Config struct {
	FolderDir    string   `json:"folder_dir"`
	FolderTypes  []string `json:"folder_types"`
	Desktop      Desktop  `json:"desktop"`
	Downloads    Downloads `json:"downloads"`
    Crontab      Crontab   `json:"crontab_settings"`
}

func HandleConfig() (Config, error) {
    var config Config

    file, err := os.Open("./config/config.json")
    if err != nil {
        return config, err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    err = decoder.Decode(&config)

    if err != nil {
        return config, err
    }

    log.Println(config)
    return config, nil

}
