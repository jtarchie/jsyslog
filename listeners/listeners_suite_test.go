package listeners_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestListeners(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Listeners Suite")
}
