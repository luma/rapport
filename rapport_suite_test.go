package rapport_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRapport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rapport Suite")
}
