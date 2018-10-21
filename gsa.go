package gsa

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
)

type StoreUsage struct {
	Containers int64 `json:"total_bytes_containers"`
	Layers     int64 `json:"total_bytes_layers"`
	Active     int64 `json:"total_bytes_active_layers"`
}

type diskUsage struct {
	TotalBytesUsed     int64 `json:"total_bytes_used"`
	ExclusiveBytesUsed int64 `json:"exclusive_bytes_used"`
}

type imageStats struct {
	DiskUsage diskUsage `json:"disk_usage"`
}

type volumeMeta struct {
	Size int64 `json:"Size"`
}

func GrootStoreUsage(bin, store, config string) StoreUsage {
	var (
		err        error
		containers int64
		layers     int64
		active     int64
	)

	containers, err = getDiskTotalContainers(bin, store, config)
	if err != nil {
		log.Println("failed to get total container image disk usage: " + err.Error())
	}

	layers, err = getDiskTotalVolumes(store)
	if err != nil {
		log.Println("failed to get total volume/layer disk usage: " + err.Error())
	}

	active, err = getDiskTotalActiveVolumes(store)
	if err != nil {
		log.Println("failed to get active layers disk usage: " + err.Error())
	}

	return StoreUsage{
		Containers: containers,
		Layers:     layers,
		Active:     active,
	}
}

func getDiskTotalContainers(bin, store, config string) (int64, error) {
	contents, err := ioutil.ReadDir(filepath.Join(store, "images"))
	if err != nil {
		return 0, err
	}

	var total int64

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

func getDiskTotalVolumes(store string) (int64, error) {
	var total int64
	contents, err := ioutil.ReadDir(filepath.Join(store, "meta"))
	if err != nil {
		return 0, err
	}

	for _, dir := range contents {
		if !dir.IsDir() {
			size, err := readVolumeMeta(filepath.Join(store, "meta", dir.Name()))
			if err != nil {
				return 0, err
			}

			total += size
		}
	}

	return total, nil
}

func getDiskTotalActiveVolumes(store string) (int64, error) {
	var total int64
	contents, err := ioutil.ReadDir(filepath.Join(store, "meta", "dependencies"))
	if err != nil {
		return 0, err
	}

	var active []string

	for _, dir := range contents {
		if !dir.IsDir() {
			shas, err := ioutil.ReadFile(filepath.Join(store, "meta", "dependencies", dir.Name()))
			if err != nil {
				return 0, err
			}

			var s []string
			if err := json.Unmarshal(shas, &s); err != nil {
				return 0, err
			}
			active = append(active, s...)
		}
	}

	for _, a := range uniq(active) {
		size, err := readVolumeMeta(filepath.Join(store, "meta", fmt.Sprintf("volume-%s", a)))
		if err != nil {
			return 0, err
		}

		total += size
	}

	return total, nil
}

func readVolumeMeta(file string) (int64, error) {
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
