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

type Config struct {
	VmArgsFileLocation       string
	AppConfigFileLocation    string
	RiakCsExecutableLocation string
	RiakCsPidFileLocation    string
	IP                       string
}

type RiakCSStartManager struct {
	osHelper os_helper.OsHelper
	config   Config
}

func New(osHelper os_helper.OsHelper, config Config) *RiakCSStartManager {
	return &RiakCSStartManager{
		osHelper: osHelper,
		config:   config,
	}
}

func (m *RiakCSStartManager) Execute() {
	err := m.replaceAllIpInFiles(m.config.IP)
	if err != nil {
		panic(err)
	}

	out, err := m.osHelper.RunCommand(m.config.RiakCsExecutableLocation, "start")
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
			err = m.osHelper.WriteStringToFile(m.config.RiakCsPidFileLocation, pid)
			m.printWithTimestamp("Found the pid " + pid + " after " + strconv.Itoa(timeout) + " seconds.")
			return
		} else if err.Error() != "1" {
			panic(err)
		}
	}
	panic(errors.New("Timed out looking for RiakCS PID"))
}

func (m *RiakCSStartManager) replaceAllIpInFiles(newIp string) error {
	vmArgsFileContents, err := m.osHelper.ReadFile(m.config.VmArgsFileLocation)
	if err != nil {
		return err
	}

	var appConfigFileContents string
	appConfigFileContents, err = m.osHelper.ReadFile(m.config.AppConfigFileLocation)
	if err != nil {
		return err
	}

	newVmArgsFileContents := strings.Replace(vmArgsFileContents, "127.0.0.1", newIp, -1)
	newAppConfigFileContents := strings.Replace(appConfigFileContents, "127.0.0.1", newIp, -1)

	err = m.osHelper.WriteStringToFile(m.config.VmArgsFileLocation, newVmArgsFileContents)
	if err != nil {
		return err
	}

	err = m.osHelper.WriteStringToFile(m.config.AppConfigFileLocation, newAppConfigFileContents)
	if err != nil {
		return err
	}

	return nil
}
