package os_helper

import (
	"io/ioutil"
	"net"
	"os"
	"os/exec"
)

type OsHelperImpl struct{}

func New() *OsHelperImpl {
	return &OsHelperImpl{}
}

type OsHelper interface {
	GetIp() (string, error)
	ReadFile(filename string) (string, error)
	WriteStringToFile(filename string, contents string) error
	RunCommand(executable string, args ...string) (string, error)
}

func (m OsHelperImpl) GetIp() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		return "", err
	}

	return addrs[0], nil
}

// Read the whole file, panic on err
func (m OsHelperImpl) ReadFile(filename string) (string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

// Overwrite the contents, creating if necessary. Panic on err
func (m OsHelperImpl) WriteStringToFile(filename string, contents string) error {
	err := ioutil.WriteFile(filename, []byte(contents), 0644)
	return err
}

// Runs command with stdout and stderr pipes connected to process
func (h OsHelperImpl) RunCommand(executable string, args ...string) (string, error) {
	cmd := exec.Command(executable, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), err
	}
	return string(out), nil
}
