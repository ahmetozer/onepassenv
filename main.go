package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type profile struct {
	ProfileName string   `json:"profileName"`
	Description string   `json:"description"`
	Variables   []string `json:"variables"`
}

type config struct {
	Profiles    []profile `json:"profiles"`
	OpPath      string    `json:"opPath"`
	AllowedBins []string  `json:"allowedBins"`
}

type onePassSecret struct {
	Id    string `json:"id"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type onePassOutput struct {
	Fields []onePassSecret `json:"fields"`
}

//! Don't trust env or os calls to find home path
//* Change this path for yours
const configPath = "/Users/ahmet/.onepassenv.json"

func main() {

	// config file control
	info, err := os.Stat(configPath)
	if err != nil {
		log.Printf("stat: %v", err)
		os.Exit(1)
	}

	if !(info.Mode().String() == "-rw----r--" || info.Mode().String() == "-rw-r--r--") {
		log.Printf(`error: file is too open
sudo chown root '%v'
sudo chmod u+w '%v'
sudo chmod a+r '%v'`, configPath, configPath, configPath)
		os.Exit(1)
	}

	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Printf("file open: %v", err)
		os.Exit(1)
	}

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Printf("io read: %v", err)
		os.Exit(1)
	}

	var config config
	json.Unmarshal(byteValue, &config)
	onepenv := os.Getenv("onepenv")
	if onepenv == "" {
		log.Fatalf("environment variable 'onepenv' is empty")
	}

	onepenvBin := os.Getenv("onepenvbin")
	if onepenv == "" {
		log.Fatalf("environment variable 'onepenvbin' is empty")
	}

	onepenvBin, err = exec.LookPath(onepenvBin)
	if err != nil {
		log.Fatalf("look path: %s", err)
	}
	if !(contains(config.AllowedBins, onepenvBin)) {
		log.Fatalf("Binnary '%v' is not in allowedBins", onepenvBin)
	}

	currentProfile, err := getProfile(config.Profiles, onepenv)
	if err != nil {
		log.Fatalf("profile Error %s", err)
	}

	// Execute onepass to get password
	onePassword := exec.Command(config.OpPath, "item", "get", onepenv, "--format", "json")
	onePasswordOut, err := onePassword.Output()
	if err != nil {
		log.Fatalf("onePassword Error %s", err)
	}
	var secrets onePassOutput
	json.Unmarshal(onePasswordOut, &secrets)

	for _, element := range secrets.Fields {
		if contains(currentProfile.Variables, element.Label) {
			os.Setenv(element.Label, element.Value)
		}
	}

	executeArgs := os.Args
	executeArgs[0] = onepenvBin
	err = syscall.Exec(onepenvBin, executeArgs, os.Environ())
	if err != nil {
		log.Fatalf("proccess execute failed for '%v': %v", onepenvBin, err)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func getProfile(p []profile, name string) (profile, error) {
	for _, v := range p {
		if v.ProfileName == name {
			return v, nil
		}
	}
	return profile{}, errors.New("profile not found")
}
