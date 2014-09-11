package riakcs_start_manager_test

import (
	os_fakes "github.com/cloudfoundry-incubator/riakcs_ctrl/os_helper/fakes"

	manager "."
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RiakCSStartManager", func() {

	var mgr *manager.RiakCSStartManager
	var fakeOsHelper *os_fakes.FakeOsHelper

	vmArgsFileLocation := "/some-unused-location"
	appConfigFileLocation := "/another-unused-location"

	vmArgsFileContents := "VM ARGS: This is our IP address: 127.0.0.1. In case you missed it, it's 127.0.0.1"
	appConfigFileContents := "APP CONFIG: This is our IP address: 127.0.0.1. In case you missed it, it's 127.0.0.1"

	Context("during normal boot", func() {
		BeforeEach(func() {
			fakeOsHelper = new(os_fakes.FakeOsHelper)

			mgr = manager.New(
				fakeOsHelper,
				vmArgsFileLocation,
				appConfigFileLocation,
			)

			fakeOsHelper.ReadFileStub = func(filepath string) (string, error) {
				if filepath == vmArgsFileLocation {
					return vmArgsFileContents, nil
				} else if filepath == appConfigFileLocation {
					return appConfigFileContents, nil
				}
				panic("Unrecognized filepath - please update test.")
			}
		})

		It("replaces all instances of 127.0.0.1 with host ip in config files", func() {
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
	})
})
