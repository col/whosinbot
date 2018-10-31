package hangout_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHangout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hangout Suite")
}
