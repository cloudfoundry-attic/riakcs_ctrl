package riakcs_start_manager_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRiakcs_start_manager(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "RiakCS Start Manager Suite")
}
