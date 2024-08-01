package main

import (
	"os"
	"os/user"
)

func Breakdown() {
	currentUser, err := user.Current()
	if err != nil {
		Logger.Fatal("error when getting current user: %s", err.Error())
	}
	Logger.Debug("%s %s", currentUser.Username, currentUser.Uid)
	if currentUser.Uid != "0" {
		Logger.Fatal("Not running as root. We're going to be removing users and deleting files.")
	}

	Logger.Info("BREAKING DOWN DAEMON")

	Logger.Info("disabling cabinet service with systemd")
	err = run("systemctl", "stop", "cabinet")
	if err != nil {
		Logger.Fatal("when trying to stop cabinet: %s", err.Error())
	}
	err = run("systemctl", "disable", "cabinet")
	if err != nil {
		Logger.Fatal("when trying to disable cabinet: %s", err.Error())
	}

	Logger.Info("deleting /etc/systemd/system/service.cabinet")
	err = os.Remove("/etc/systemd/system/cabinet.service")
	if err != nil {
		Logger.Fatal("when trying to remove service file: %s", err.Error())
	}

	Logger.Info("deleting /usr/local/bin/cabinet")
	err = os.Remove("/usr/local/bin/cabinet")
	if err != nil {
		Logger.Fatal("when trying to remove binary file: %s", err.Error())
	}

	Logger.Info("deleting /usr/local/share/Cabinet")
	err = os.Remove("/usr/local/share/Cabinet")
	if err != nil {
		Logger.Fatal("when trying to remove Cabinet share: %s", err.Error())
	}

	Logger.Info("deleting cabinet user")
	err = run("userdel", "cabinet")
	if err != nil {
		Logger.Fatal("when trying to run userdel: %s", err.Error())
	}

	Logger.Info("BREAK DOWN COMPLETE")

	os.Exit(0)
}
