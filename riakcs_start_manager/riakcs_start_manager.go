package riakcs_start_manager

import (
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/riakcs_ctrl/os_helper"
)

type RiakCSStartManager struct {
	osHelper                 os_helper.OsHelper
	vmArgsFileLocation       string
	appConfigFileLocation    string
	riakCsExecutableLocation string
	riakCsPidFileLocation    string
}

func New(osHelper os_helper.OsHelper,
	vmArgsFileLocation string,
	appConfigFileLocation string,
	riakCsExecutableLocation string,
	riakCsPidFileLocation string,
) *RiakCSStartManager {
	return &RiakCSStartManager{
		osHelper:                 osHelper,
		vmArgsFileLocation:       vmArgsFileLocation,
		appConfigFileLocation:    appConfigFileLocation,
		riakCsExecutableLocation: riakCsExecutableLocation,
		riakCsPidFileLocation:    riakCsPidFileLocation,
	}
}

func (m *RiakCSStartManager) Execute() {
	ip, err := m.osHelper.GetIp()
	if err != nil {
		panic(err)
	}

	err = m.replaceAllIpInFiles(ip)
	if err != nil {
		panic(err)
	}

	_, err = m.osHelper.RunCommand(m.riakCsExecutableLocation, "start")
	if err != nil {
		panic(err)
	}

	// We sleep to ensure the PID exists when we query for it.
	time.Sleep(2 * time.Second)

	// Get the PID and write it to the pidfile
	var pid string
	pid, err = m.osHelper.RunCommand("pgrep", "-f", "'beam.*riak-cs'")
	if err != nil {
		panic(err)
	}
	err = m.osHelper.WriteStringToFile(m.riakCsPidFileLocation, pid)

	return
}

func (m *RiakCSStartManager) replaceAllIpInFiles(newIp string) error {
	vmArgsFileContents, err := m.osHelper.ReadFile(m.vmArgsFileLocation)
	if err != nil {
		return err
	}

	var appConfigFileContents string
	appConfigFileContents, err = m.osHelper.ReadFile(m.appConfigFileLocation)
	if err != nil {
		return err
	}

	newVmArgsFileContents := strings.Replace(vmArgsFileContents, "127.0.0.1", newIp, -1)
	newAppConfigFileContents := strings.Replace(appConfigFileContents, "127.0.0.1", newIp, -1)

	err = m.osHelper.WriteStringToFile(m.vmArgsFileLocation, newVmArgsFileContents)
	if err != nil {
		return err
	}

	err = m.osHelper.WriteStringToFile(m.appConfigFileLocation, newAppConfigFileContents)
	if err != nil {
		return err
	}

	return nil
}
