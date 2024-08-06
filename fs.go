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
