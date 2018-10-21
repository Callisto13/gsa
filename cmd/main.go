package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Callisto13/gsa"
	humanize "github.com/dustin/go-humanize"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	var (
		grootfsConfig string
		grootfsBin    string
		humanReadble  bool
	)

	flag.StringVar(&grootfsConfig, "grootfs-config", "/var/vcap/jobs/garden/config/grootfs_config.yml", "path to grootfs' config")
	flag.StringVar(&grootfsBin, "grootfs-bin", "/var/vcap/packages/grootfs/bin/grootfs", "path to the grootfs bin")
	flag.BoolVar(&humanReadble, "r", false, "human readable result")
	flag.Parse()

	if _, err := os.Stat(grootfsBin); os.IsNotExist(err) {
		fmt.Println("grootfs not found")
		os.Exit(1)
	}

	if _, err := os.Stat(grootfsConfig); os.IsNotExist(err) {
		fmt.Println("grootfs config not found")
		os.Exit(1)
	}

	store, err := grootfsStorePath(grootfsConfig)
	if err != nil {
		fmt.Println("failed to read grootfs store path from config: " + err.Error())
		os.Exit(1)
	}

	usage := gsa.GrootStoreUsage(grootfsBin, store, grootfsConfig)

	if humanReadble {
		fmt.Printf("Containers: %s\n", humanize.Bytes(usage.Containers))
		fmt.Printf("Layers: %s (of which Active: %s)\n", humanize.Bytes(usage.Layers), humanize.Bytes(usage.Active))
		fmt.Println(humanize.Bytes(usage.Total))
		os.Exit(0)
	}

	result, err := json.Marshal(usage)
	if err != nil {
		fmt.Println("failed to marshal result: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(string(result))
	os.Exit(0)
}

type config struct {
	StorePath string `yaml:"store"`
}

func grootfsStorePath(path string) (string, error) {
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	var c *config
	if err = yaml.Unmarshal(yml, &c); err != nil {
		return "", err
	}

	return c.StorePath, nil
}
