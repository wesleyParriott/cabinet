package main

import (
	"os"
	"os/exec"
	"os/user"
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
	return chown("cabinet", "cabinet", path)
}

func makeCabinetDirectory() error {
	err := os.Mkdir("/usr/local/share/Cabinet", 0755)
	if err != nil {
		if os.IsExist(err) {
			Logger.Info("skipping creating cabinet path because it already exists")
		} else {
			return err
		}
	}

	err = chownToCabinet("/usr/local/share/Cabinet")
	if err != nil {
		return err
	}

	return nil
}

func makeCabinetDataDirectory() error {
	err := os.Mkdir("/usr/local/share/CabinetData/", 0755)
	if err != nil {
		if os.IsExist(err) {
			Logger.Info("skipping creating cabinet data path because it already exists")
		} else {
			return err
		}
	}

	err = chownToCabinet("/usr/local/share/CabinetData/")
	if err != nil {
		return err
	}

	err = os.Mkdir("/usr/local/share/CabinetData/tmpls", 0755)
	if err != nil {
		if os.IsExist(err) {
			Logger.Info("skipping creating cabinet data path because it already exists")
		} else {
			return err
		}
	}

	err = chownToCabinet("/usr/local/share/CabinetData/tmpls")
	if err != nil {
		return err
	}

	// list every file in ./tmpls
	fileNames, _, err := listDir("./tmpls")
	for _, fileName := range fileNames {
		Logger.Debug("copying %s", fileName)
		filePath := "./tmpls/" + fileName
		destFilePath := "/usr/local/share/CabinetData/tmpls/" + fileName
		err = copyFile(filePath, destFilePath)
		if err != nil {
			return err
		}
	}

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

	Logger.Info("making tmpls directory at /usr/local/share/CabinetData/tmpls based on ./tmpls")
	err = makeCabinetDataDirectory()
	if err != nil {
		Logger.Fatal("when making the tmpls directory: %s", err.Error())
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

	Logger.Info("copying ./setup/favicon.ico to /usr/local/share/CabinetData/favicon.ico")
	err = copyFile("./setup/favicon.ico", "/usr/local/share/CabinetData/favicon.ico")
	if err != nil {
		Logger.Fatal("when copying favicon: %s", err.Error())
	}
	err = chown("cabinet", "cabinet", "/usr/local/share/CabinetData/favicon.ico")
	if err != nil {
		Logger.Fatal("when chowning favicon: %s", err.Error())
	}

	Logger.Info("copying ./setup/passcode.jpg to /usr/local/share/CabinetData/passcode.jpg")
	err = copyFile("./setup/passcode.jpg", "/usr/local/share/CabinetData/passcode.jpg")
	if err != nil {
		Logger.Fatal("when copying passcode jpg: %s", err.Error())
	}
	err = chown("cabinet", "cabinet", "/usr/local/share/CabinetData/passcode.jpg")
	if err != nil {
		Logger.Fatal("when chowning passcode jpg: %s", err.Error())
	}

	Logger.Info("copying .passcode to /usr/local/share/CabinetData/.passcode")
	err = copyFile("./.passcode", "/usr/local/share/CabinetData/.passcode")
	if err != nil {
		Logger.Fatal("when copying passcode file: %s", err.Error())
	}
	err = chown("cabinet", "cabinet", "/usr/local/share/CabinetData/.passcode")
	if err != nil {
		Logger.Fatal("when chowning passcode file: %s", err.Error())
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
