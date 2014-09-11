package riakcs_start_manager_test

import (
	"strings"
	"time"

	manager "."
	os_fakes "github.com/cloudfoundry-incubator/riakcs_ctrl/os_helper/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RiakCSStartManager", func() {

	var mgr *manager.RiakCSStartManager
	var fakeOsHelper *os_fakes.FakeOsHelper

	vmArgsFileLocation := "/some-unused-location"
	appConfigFileLocation := "/another-unused-location"
	riakCsExecutableLocation := "/riak-cs-location"
	riakCsPidFileLocation := "/riak-cs-pid-file-location"

	vmArgsFileContents := "VM ARGS: This is our IP address: 127.0.0.1. In case you missed it, it's 127.0.0.1"
	appConfigFileContents := "APP CONFIG: This is our IP address: 127.0.0.1. In case you missed it, it's 127.0.0.1"

	Context("during normal boot", func() {
		BeforeEach(func() {
			fakeOsHelper = new(os_fakes.FakeOsHelper)

			mgr = manager.New(
				fakeOsHelper,
				vmArgsFileLocation,
				appConfigFileLocation,
				riakCsExecutableLocation,
				riakCsPidFileLocation,
			)
		})

		It("replaces all instances of 127.0.0.1 with host ip in config files", func() {
			fakeOsHelper.ReadFileStub = func(filepath string) (string, error) {
				if filepath == vmArgsFileLocation {
					return vmArgsFileContents, nil
				} else if filepath == appConfigFileLocation {
					return appConfigFileContents, nil
				}
				panic("Unrecognized filepath - please update test.")
			}

			fakeIp := "1.2.3.4"
			fakeOsHelper.GetIpReturns(fakeIp, nil)

			mgr.Execute()

			filepath0, contents0 := fakeOsHelper.WriteStringToFileArgsForCall(0)
			Expect(filepath0).To(Equal(vmArgsFileLocation))
			Expect(contents0).To(Equal("VM ARGS: This is our IP address: " + fakeIp + ". In case you missed it, it's " + fakeIp))

			filepath1, contents1 := fakeOsHelper.WriteStringToFileArgsForCall(1)
			Expect(filepath1).To(Equal(appConfigFileLocation))
			Expect(contents1).To(Equal("APP CONFIG: This is our IP address: " + fakeIp + ". In case you missed it, it's " + fakeIp))
		})

		It("calls the RiakCS start script", func() {
			mgr.Execute()

			executable, args := fakeOsHelper.RunCommandArgsForCall(0)
			Expect(executable).To(Equal(riakCsExecutableLocation))
			Expect(args).To(Equal([]string{"start"}))
		})

		// There is no way to directly test that we called time.Sleep, nor exactly where in the function we called it.
		// This is the best we can do.
		It("takes between two and three seconds to execute", func() {
			startTime := time.Now()
			mgr.Execute()
			endTime := time.Now()

			Expect(endTime.After(startTime.Add(2 * time.Second))).To(BeTrue())
			Expect(endTime.After(startTime.Add(3 * time.Second))).To(BeFalse())
		})

		It("captures the pid of riakCS and writes that to a file", func() {
			fakeRiakPid := "12345"
			fakeOsHelper.RunCommandStub = func(executable string, args ...string) (string, error) {
				if executable == riakCsExecutableLocation {
					return "", nil
				} else if executable == "pgrep" {
					return fakeRiakPid, nil
				}
				panic("Unrecognized command: " + executable + strings.Join(args, " ") + " - please update test.")
			}

			mgr.Execute()

			executable, args := fakeOsHelper.RunCommandArgsForCall(1)
			Expect(executable).To(Equal("pgrep"))
			Expect(args).To(Equal([]string{"-f", "'beam.*riak-cs'"}))

			// The first two file writes are the config files. The third is the pidfile.
			filepath, contents := fakeOsHelper.WriteStringToFileArgsForCall(2)
			Expect(filepath).To(Equal(riakCsPidFileLocation))
			Expect(contents).To(Equal(fakeRiakPid))
		})
	})
})
