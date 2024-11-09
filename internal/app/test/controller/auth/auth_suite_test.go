package auth_test

import (
	"testing"

	test_helpers "github.com/Hari-Krishna-Moorthy/task-management-system/internal/app/test"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuthController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Suite")
}

var _ = AfterSuite(func() {
	test_helpers.TruncatesCollection()
})

var _ = BeforeSuite(func() {
	test_helpers.TruncatesCollection()
})
