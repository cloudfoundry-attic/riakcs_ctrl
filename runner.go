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

var riakCsExecutableLocation = flag.String(
	"riakCsExecutable",
	"",
	"Specifies the location of the RiakCS executable",
)

var riakCsPidFileLocation = flag.String(
	"riakCsPidFile",
	"",
	"Specifies the location of the RiakCS PID file",
)

var ip = flag.String(
	"ip",
	"",
	"My ip - find binding to Riak and RiakCS",
)

func main() {
	flag.Parse()
	osHelper := os_helper.New()
	mgr := manager.New(
		*osHelper,
		*vmArgsFileLocation,
		*appConfigFileLocation,
		*riakCsExecutableLocation,
		*riakCsPidFileLocation,
		*ip,
	)
	mgr.Execute()

}
