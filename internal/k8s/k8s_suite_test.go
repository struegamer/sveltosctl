package k8s_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/cluster-api/util"
)

func TestShow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Show Suite")
}

func randomString() string {
	const length = 10
	return util.RandomString(length)
}
