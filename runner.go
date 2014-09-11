package main

import (
	"flag"

	"github.com/cloudfoundry-incubator/riakcs_ctrl/os_helper"
	manager "github.com/cloudfoundry-incubator/riakcs_ctrl/riakcs_start_manager"
)

var vmArgsFileLocation = flag.String(
	"vmArgsFile",
	"",
	"Specifies the location of the vm args",
)

var appConfigFileLocation = flag.String(
	"appConfigFile",
	"",
	"Specifies the location of the app config",
)

func main() {
	flag.Parse()
	osHelper := os_helper.New()
	mgr := manager.New(
		*osHelper,
		*vmArgsFileLocation,
		*appConfigFileLocation,
	)
	mgr.Execute()

}
