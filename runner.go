package main

import (
	"flag"
	"os"

	"github.com/cloudfoundry-incubator/riakcs_ctrl/os_helper"
	manager "github.com/cloudfoundry-incubator/riakcs_ctrl/riakcs_start_manager"
	"github.com/pivotal-cf-experimental/service-config"
)

func main() {

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	serviceConfig := service_config.New()
	serviceConfig.AddFlags(flags)

	flags.Parse(os.Args[1:])

	var config manager.Config
	serviceConfig.Read(&config)

	osHelper := os_helper.New()
	mgr := manager.New(
		*osHelper,
		config,
	)
	mgr.Execute()
}
