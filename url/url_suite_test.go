package url_test

import (
	"github.com/jtarchie/jsyslog/url"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUrl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Url Suite")
}

var _ = Describe("Parse", func() {
	When("the scheme is file://", func() {
		It("supports windows paths", func() {
			uri, err := url.Parse("file://C:\\some\\path\\test.log")
			Expect(err).ToNot(HaveOccurred())

			Expect(uri.Path).To(Equal(`C:\some\path\test.log`))
			Expect(uri.Scheme).To(Equal("file"))
			Expect(uri.String()).To(Equal(`file://C:\some\path\test.log`))
		})

		It("supports unix paths", func() {
			uri, err := url.Parse("file:///var/log/test.log")
			Expect(err).ToNot(HaveOccurred())

			Expect(uri.Path).To(Equal(`/var/log/test.log`))
			Expect(uri.Scheme).To(Equal("file"))
			Expect(uri.String()).To(Equal(`file:///var/log/test.log`))
		})
	})

	When("any other scheme", func() {
		It("fallbacks to the url.Parse", func() {
			uri, err := url.Parse("https://user:password@google.com:1234/path")
			Expect(err).ToNot(HaveOccurred())

			Expect(uri.Path).To(Equal("/path"))
			Expect(uri.Scheme).To(Equal("https"))
			Expect(uri.User.Username()).To(Equal("user"))

			password, _ := uri.User.Password()
			Expect(password).To(Equal("password"))
			Expect(uri.Hostname()).To(Equal("google.com"))
			Expect(uri.Port()).To(Equal("1234"))
		})
	})
})
