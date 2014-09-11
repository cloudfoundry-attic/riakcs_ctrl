package riakcs_start_manager

import (
	"github.com/cloudfoundry-incubator/riakcs_ctrl/os_helper"
	"strings"
)

type RiakCSStartManager struct {
	vmArgsFileLocation    string
	appConfigFileLocation string
	osHelper              os_helper.OsHelper
}

func New(osHelper os_helper.OsHelper,
	vmArgsFileLocation string,
	appConfigFileLocation string,
) *RiakCSStartManager {
	return &RiakCSStartManager{
		vmArgsFileLocation:    vmArgsFileLocation,
		appConfigFileLocation: appConfigFileLocation,
		osHelper:              osHelper,
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
