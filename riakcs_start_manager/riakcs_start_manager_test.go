package riakcs_start_manager_test

import (
	"fmt"
	"strings"
	"time"

	os_fakes "github.com/cloudfoundry-incubator/riakcs_ctrl/os_helper/fakes"
	manager "github.com/cloudfoundry-incubator/riakcs_ctrl/riakcs_start_manager"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RiakCSStartManager", func() {
	var (
		mgr                   *manager.RiakCSStartManager
		fakeOsHelper          *os_fakes.FakeOsHelper
		config                manager.Config
		vmArgsFileContents    string
		appConfigFileContents string
	)

	Context("during normal boot", func() {
		BeforeEach(func() {
			vmArgsFileContents = "VM ARGS: This is our IP address: 127.0.0.1.\n In case you missed it, it's 127.0.0.1"
			appConfigFileContents = "APP CONFIG: This is our IP address: 127.0.0.1.\n In case you missed it, it's 127.0.0.1"

			config = manager.Config{
				VmArgsFileLocation:       "/some-unused-location",
				AppConfigFileLocation:    "/another-unused-location",
				RiakCsExecutableLocation: "/riak-cs-location",
				RiakCsPidFileLocation:    "/riak-cs-pid-file-location",
				IP: "1.2.3.4",
			}

			fakeOsHelper = new(os_fakes.FakeOsHelper)

			mgr = manager.New(
				fakeOsHelper,
				config,
			)
		})

		It("replaces all instances of 127.0.0.1 with host ip in config files", func() {
			fakeOsHelper.ReadFileStub = func(filepath string) (string, error) {
				if filepath == config.VmArgsFileLocation {
					return vmArgsFileContents, nil
				} else if filepath == config.AppConfigFileLocation {
					return appConfigFileContents, nil
				}
				panic("Unrecognized filepath - please update test.")
			}

			mgr.Execute()

			filepath0, contents0 := fakeOsHelper.WriteStringToFileArgsForCall(0)
			Expect(filepath0).To(Equal(config.VmArgsFileLocation))
			Expect(contents0).To(Equal(
				fmt.Sprintf("VM ARGS: This is our IP address: %s.\n In case you missed it, it's %s", config.IP, config.IP),
			))

			filepath1, contents1 := fakeOsHelper.WriteStringToFileArgsForCall(1)
			Expect(filepath1).To(Equal(config.AppConfigFileLocation))
			Expect(contents1).To(Equal(
				fmt.Sprintf("APP CONFIG: This is our IP address: %s.\n In case you missed it, it's %s", config.IP, config.IP),
			))
		})

		It("calls the RiakCS start script", func() {
			mgr.Execute()

			executable, args := fakeOsHelper.RunCommandArgsForCall(0)
			Expect(executable).To(Equal(config.RiakCsExecutableLocation))
			Expect(args).To(Equal([]string{"start"}))
		})

		// There is no way to directly test that we called time.Sleep, nor exactly where in the function we called it.
		// This is the best we can do.
		It("takes between one and sixty-one seconds to execute", func() {
			startTime := time.Now()
			mgr.Execute()

			elapsedTime := time.Since(startTime)
			Expect(elapsedTime > 1*time.Second).To(BeTrue())
			Expect(elapsedTime < 61*time.Second).To(BeTrue())
		})

		It("captures the pid of riakCS and writes that to a file", func() {
			fakeRiakPid := "12345"
			fakeOsHelper.RunCommandStub = func(executable string, args ...string) (string, error) {
				if executable == config.RiakCsExecutableLocation {
					return "", nil
				} else if executable == "pgrep" {
					return fakeRiakPid, nil
				}
				panic("Unrecognized command: " + executable + strings.Join(args, " ") + " - please update test.")
			}

			mgr.Execute()

			executable, args := fakeOsHelper.RunCommandArgsForCall(1)
			Expect(executable).To(Equal("pgrep"))
			Expect(args).To(Equal([]string{"-f", "beam.*riak-cs"}))

			// The first two file writes are the config files. The third is the pidfile.
			filepath, contents := fakeOsHelper.WriteStringToFileArgsForCall(2)
			Expect(filepath).To(Equal(config.RiakCsPidFileLocation))
			Expect(contents).To(Equal(fakeRiakPid))
		})
	})
})
