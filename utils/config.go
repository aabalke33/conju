package utils

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DatabaseDirectory string
	AudioDirectory    string
	ExportLocation    string
	FFPlayLocation    string
	Languages         map[string]ConfigLanguage
}

type ConfigLanguage struct {
	DefaultConjugations map[string]bool
}

type ConfigFile struct {
	DatabaseDirectory string              `json:"database_directory"`
	AudioDirectory    string              `json:"audio_directory"`
	ExportLocation    string              `json:"export_location"`
	FFPlayLocation    string              `json:"ffplay_location"`
	Languages         map[string]Language `json:"languages"`
}

type Language struct {
	DefaultConjugations map[string]interface{} `json:"default_conjugations"`
}

func ReadConfig() Config {

	configFilePath := "./data/config.json"
	var configFile ConfigFile

	jsonFile, err := os.Open(configFilePath)
	if err != nil {
		panic("Could Not Read Config.json")
	}

	defer jsonFile.Close()

	bytes, _ := io.ReadAll(jsonFile)

	err = json.Unmarshal(bytes, &configFile)
	if err != nil {
		panic("Error Unmarshalling JSON in Config.json")
	}

	config := Config{
		DatabaseDirectory: configFile.DatabaseDirectory,
		AudioDirectory:    configFile.AudioDirectory,
		ExportLocation:    configFile.ExportLocation,
		FFPlayLocation:    configFile.FFPlayLocation,
		Languages:         make(map[string]ConfigLanguage),
	}

	for language, languageData := range configFile.Languages {
		var conjugations = languageData.DefaultConjugations

		var tempLang ConfigLanguage
		tempLang.DefaultConjugations = make(map[string]bool)

		for conjugation, value := range conjugations {
			switch value := value.(type) {
			case bool:
				tempLang.DefaultConjugations[conjugation] = value
			default:
				panic("Values language conjugation config.json are not bools")
			}
		}
		config.Languages[language] = tempLang
	}
	return config
}
