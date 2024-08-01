package main

import "os"

func listDir(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	var ret []string
	if err != nil {
		return ret, err
	}
	for _, entry := range entries {
		ret = append(ret, entry.Name())
	}

	return ret, nil
}
