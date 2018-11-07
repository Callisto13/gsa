package gsa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type StoreUsage struct {
	Containers uint64 `json:"total_bytes_containers"`
	Layers     uint64 `json:"total_bytes_layers"`
	Active     uint64 `json:"total_bytes_active_layers"`
	Total      uint64 `json:"total_bytes_store"`
}

type diskUsage struct {
	TotalBytesUsed     uint64 `json:"total_bytes_used"`
	ExclusiveBytesUsed uint64 `json:"exclusive_bytes_used"`
}

type imageStats struct {
	DiskUsage diskUsage `json:"disk_usage"`
}

type volumeMeta struct {
	Size uint64 `json:"Size"`
}

type config struct {
	StorePath string `yaml:"store"`
}

func GrootStoreUsage(bin, config string) StoreUsage {
	var (
		err        error
		store      string
		containers uint64
		layers     uint64
		active     uint64
	)

	store, err = grootfsStorePath(config)
	if err != nil {
		fmt.Println("failed to read grootfs store path from config: " + err.Error())
		os.Exit(1)
	}

	containers, err = getDiskTotalContainers(bin, store, config)
	if err != nil {
		log.Println("failed to get total container image disk usage: " + err.Error())
	}

	layers, active, err = getDiskTotalVolumes(store)
	if err != nil {
		log.Println("failed to get active layers disk usage: " + err.Error())
	}

	return StoreUsage{
		Containers: containers,
		Layers:     layers,
		Active:     active,
		Total:      containers + layers,
	}
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

func getDiskTotalContainers(bin, store, config string) (uint64, error) {
	contents, err := ioutil.ReadDir(filepath.Join(store, "images"))
	if err != nil {
		return 0, err
	}

	var total uint64

	for _, dir := range contents {
		args := []string{"--config", config, "stats", dir.Name()}
		grootCmd := exec.Command(bin, args...)
		output, err := grootCmd.CombinedOutput()
		if err != nil {
			return 0, fmt.Errorf("could not execute grootfs: %s : %#v", err.Error(), string(output))
		}

		var is *imageStats
		if err := json.Unmarshal(output, &is); err != nil {
			return 0, fmt.Errorf("mangled response from grootfs: %s : %#v", err.Error(), string(output))
		}

		total += is.DiskUsage.ExclusiveBytesUsed
	}

	return total, nil
}

func getDiskTotalVolumes(store string) (uint64, uint64, error) {
	vols, err := getListVolumes(store)
	if err != nil {
		return 0, 0, err
	}

	aVols, err := getListActiveVolumes(store)
	if err != nil {
		return 0, 0, err
	}

	var total, active uint64

	for _, dir := range difference(vols, aVols) {
		size, err := readVolumeMeta(filepath.Join(store, "meta", dir))
		if err != nil {
			return 0, 0, err
		}

		total += size
	}

	for _, dir := range aVols {
		size, err := readVolumeMeta(filepath.Join(store, "meta", dir))
		if err != nil {
			return 0, 0, err
		}

		total += size
		active += size
	}

	return total, active, nil
}

func getListVolumes(store string) ([]string, error) {
	files, err := ioutil.ReadDir(filepath.Join(store, "meta"))
	if err != nil {
		return []string{}, err
	}

	var vols []string
	for _, dir := range files {
		if !dir.IsDir() {
			vols = append(vols, dir.Name())
		}
	}

	return vols, nil
}

func getListActiveVolumes(store string) ([]string, error) {
	contents, err := ioutil.ReadDir(filepath.Join(store, "meta", "dependencies"))
	if err != nil {
		return []string{}, err
	}

	var shas []string
	for _, dir := range contents {
		if !dir.IsDir() {
			contents, err := ioutil.ReadFile(filepath.Join(store, "meta", "dependencies", dir.Name()))
			if err != nil {
				return []string{}, err
			}

			if err := json.Unmarshal(contents, &shas); err != nil {
				return []string{}, err
			}
		}
	}

	var active []string
	for _, sha := range uniq(shas) {
		active = append(active, fmt.Sprintf("volume-%s", sha))
	}

	return active, nil
}

func readVolumeMeta(file string) (uint64, error) {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}

	var vm *volumeMeta
	if err := json.Unmarshal(contents, &vm); err != nil {
		return 0, err
	}

	return vm.Size, nil
}

func difference(a, b []string) []string {
	inB := map[string]bool{}
	for _, x := range b {
		inB[x] = true
	}

	diff := []string{}
	for _, x := range a {
		if _, ok := inB[x]; !ok {
			diff = append(diff, x)
		}
	}

	return diff
}

func uniq(elements []string) []string {
	found := map[string]bool{}
	result := []string{}

	for v := range elements {
		if found[elements[v]] == true {
		} else {
			found[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}
