package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHugoParser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "HugoParser Suite")
}
