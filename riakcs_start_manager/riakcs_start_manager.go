package riakcs_start_manager

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cloudfoundry-incubator/riakcs_ctrl/os_helper"
)

func (m *RiakCSStartManager) printWithTimestamp(output string) {
	fmt.Printf("%v ----- %v\n", time.Now().Local(), output)
}

type RiakCSStartManager struct {
	osHelper                 os_helper.OsHelper
	vmArgsFileLocation       string
	appConfigFileLocation    string
	riakCsExecutableLocation string
	riakCsPidFileLocation    string
	ip                       string
}

func New(osHelper os_helper.OsHelper,
	vmArgsFileLocation string,
	appConfigFileLocation string,
	riakCsExecutableLocation string,
	riakCsPidFileLocation string,
	ip string,
) *RiakCSStartManager {
	return &RiakCSStartManager{
		osHelper:                 osHelper,
		vmArgsFileLocation:       vmArgsFileLocation,
		appConfigFileLocation:    appConfigFileLocation,
		riakCsExecutableLocation: riakCsExecutableLocation,
		riakCsPidFileLocation:    riakCsPidFileLocation,
		ip: ip,
	}
}

func (m *RiakCSStartManager) Execute() {
	err := m.replaceAllIpInFiles(m.ip)
	if err != nil {
		panic(err)
	}

	out, err := m.osHelper.RunCommand(m.riakCsExecutableLocation, "start")
	if err != nil && !strings.HasPrefix(out, "Node is already running!") {
		m.printWithTimestamp("Unexpected error starting RiakCS: exiting now.")
		panic(err)
	}

	// We sleep to ensure the PID exists when we query for it.
	timeout := 0
	for timeout < 60 {
		time.Sleep(1 * time.Second)
		pid, err := m.osHelper.RunCommand("pgrep", "-f", "beam.*riak-cs")

		timeout++

		if err == nil {
			pid = strings.TrimSpace(pid)
			err = m.osHelper.WriteStringToFile(m.riakCsPidFileLocation, pid)
			m.printWithTimestamp("Found the pid " + pid + " after " + strconv.Itoa(timeout) + " seconds.")
			return
		} else if err.Error() != "1" {
			panic(err)
		}
	}
	panic(errors.New("Timed out looking for RiakCS PID"))
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
