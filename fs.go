package main

import (
	"os"
	"os/user"
	"strconv"
)

func chown(userName, groupName, path string) error {
	usr, err := user.Lookup(userName)
	if err != nil {
		return err
	}

	grp, err := user.LookupGroup(groupName)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(grp.Gid)
	if err != nil {
		return err
	}

	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(currentFilePath string, destFilePath string) error {
	destFile, err := os.OpenFile(destFilePath, os.O_RDWR|os.O_CREATE, 0750)
	if err != nil {
		return err
	}

	currentFileBytes, err := os.ReadFile(currentFilePath)
	if err != nil {
		return err
	}

	totalWrite, err := destFile.Write(currentFileBytes)
	if err != nil {
		return err
	}
	Logger.Debug("Wrote %d bytes", totalWrite)

	return nil
}

// listDir looks takes in a path and returns both the files
// in that path in one array and the directories in another
func listDir(path string) ([]string, []string, error) {
	entries, err := os.ReadDir(path)
	var files []string
	var dirs []string
	if err != nil {
		return files, dirs, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		} else {
			files = append(files, entry.Name())
		}
	}

	return files, dirs, nil
}
