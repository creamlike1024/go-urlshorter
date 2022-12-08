package jsonedit

import (
	"encoding/json"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
)

type PathUrl struct {
	Db []struct {
		Path string `json:"path"`
		URL  string `json:"url"`
	} `json:"db"`
}

func AddPath(path string, url string, jsonFileName string) error {
	jsonFileContent, err := os.ReadFile(jsonFileName)
	if err != nil {
		log.Fatal("Error reading json file")
	}
	if CheckPath(path, &jsonFileContent) {
		return errors.New("Path:" + path + "already exists")
	}
	pathUrls := ParseJsonFile(jsonFileContent)
	pathUrls.Db = append(pathUrls.Db, struct {
		Path string `json:"path"`
		URL  string `json:"url"`
	}{path, url})
	jsonFileContent, _ = json.Marshal(pathUrls)
	err = os.WriteFile(jsonFileName, jsonFileContent, 0644)
	if err != nil {
		log.Fatal("Error writing to json file")
	}
	return nil
}

func DelPath(path string, jsonFileName string) error {
	jsonFileContent, err := os.ReadFile(jsonFileName)
	if err != nil {
		log.Fatal("Error reading json file")
	}
	if !CheckPath(path, &jsonFileContent) {
		return errors.New("Path:" + path + "does not exist")
	}
	pathUrls := ParseJsonFile(jsonFileContent)
	for i, pathUrl := range pathUrls.Db {
		if path == pathUrl.Path {
			pathUrls.Db = append(pathUrls.Db[:i], pathUrls.Db[i+1:]...)
			break
		}
	}
	jsonFileContent, _ = json.Marshal(pathUrls)
	err = os.WriteFile(jsonFileName, jsonFileContent, 0644)
	if err != nil {
		log.Fatal("Error writing to json file")
	}
	return nil
}

func ParseJsonFile(jsonFileContent []byte) PathUrl {
	var pathUrls PathUrl
	json.Unmarshal(jsonFileContent, &pathUrls)
	return pathUrls
}

func CheckPath(path string, jsonFileContent *[]byte) bool {
	pathUrls := ParseJsonFile(*jsonFileContent)
	for _, pathUrl := range pathUrls.Db {
		if path == pathUrl.Path {
			return true
		}
	}
	return false
}

func InitJson(jsonFileName string) {
	var pathUrls PathUrl
	jsonFileContent, _ := json.Marshal(pathUrls)
	err := os.WriteFile(jsonFileName, jsonFileContent, 0644)
	if err != nil {
		log.Fatal("Error init json file")
	}
}
