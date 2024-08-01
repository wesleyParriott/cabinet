package main

import (
	"os"
	"os/exec"
	"os/user"
	"strconv"
)

func run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	Logger.Debug("cmd: %v", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	Logger.Debug("output: %s", out)
	return nil
}

func addCabinetUser() error {
	usr, err := user.Lookup("cabinet")
	if err != nil {
		switch err.(type) {
		case user.UnknownUserError:
			break
		default:
			return err
		}
	}
	if usr == nil {
		err = run("useradd", "--no-create-home", "cabinet", "--shell=/usr/bin/nologin")
		if err != nil {
			return err
		}
	} else {
		Logger.Info("skipping creating user because it already exists")
	}

	return nil
}

func chownToCabinet(path string) error {
	usr, err := user.Lookup("cabinet")
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(usr.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(usr.Gid)
	if err != nil {
		return err
	}

	err = os.Chown(path, uid, gid)
	if err != nil {
		return err
	}

	return nil
}

func makeCabinetDirectory() error {
	err := os.Mkdir("/usr/local/share/Cabinet", 0750)
	if err != nil {
		if os.IsExist(err) {
			Logger.Info("skipping creating cabinet path because it already exists")
		} else {
			return err
		}
	}

	chownToCabinet("/usr/local/share/Cabinet")

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

func copyCabinetBinary() error {
	destFilePath := "/usr/local/bin/cabinet"

	err := copyFile("./cabinet", destFilePath)
	if err != nil {
		return err
	}

	err = chownToCabinet(destFilePath)
	if err != nil {
		return err
	}

	return nil
}

func Setup() {
	currentUser, err := user.Current()
	if err != nil {
		Logger.Fatal("error when getting current user: %s", err.Error())
	}
	Logger.Debug("%s %s", currentUser.Username, currentUser.Uid)
	if currentUser.Uid != "0" {
		Logger.Fatal("Not running as root. We're going to do some useradd and groupadd so please run as root")
	}

	Logger.Info("SETTING UP DAEMON")

	Logger.Info("adding cabinet user")
	err = addCabinetUser()
	if err != nil {
		Logger.Fatal("when adding cabinet user: %s", err.Error())
	}

	Logger.Info("creating cabinet path at /usr/local/share/Cabinet")
	err = makeCabinetDirectory()
	if err != nil {
		Logger.Fatal("when adding CabinetDirectory: %s", err.Error())
	}

	Logger.Info("copying ./cabinet to /usr/local/bin/cabinet")
	err = copyCabinetBinary()
	if err != nil {
		Logger.Fatal("when copying cabinet binary: %s", err.Error())
	}

	Logger.Info("copying ./setup/cabinet.service to /etc/systemd/system")
	err = copyFile("./setup/cabinet.service", "/etc/systemd/system/cabinet.service")
	if err != nil {
		Logger.Fatal("when copying service file: %s", err.Error())
	}

	Logger.Info("enabling cabinet service with systemd")
	err = run("systemctl", "enable", "cabinet")
	if err != nil {
		Logger.Fatal("when enabling systemd service: %s", err.Error())
	}

	Logger.Info("SETUP COMPLETED >:D")

	Logger.Info("please go ahead and run this: sudo systemctl start cabinet")

	os.Exit(0)
}
