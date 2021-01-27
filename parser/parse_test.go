package parser_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/parser"
	"github.com/jtarchie/jsyslog/payload"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

const validMessage = `<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su 12345 98765 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] 'su root' failed for lonvick on /dev/pts/8`

func ParseLogMessageTests(message string) {
	It("parses valid messages", func() {
		_, offset, err := parser.Parse(message)
		Expect(offset).ToNot(Equal(0))
		Expect(err).ToNot(HaveOccurred())
	})

	It("renders back to its original format", func() {
		log, _, _ := parser.Parse(message)
		Expect(log.String()).To(Equal(validMessage))
	})

	It("sets the version", func() {
		log, _, _ := parser.Parse(message)
		Expect(log.Version()).To(Equal(1))
	})

	Context("with the priority", func() {
		It("returns the severity", func() {
			log, _, _ := parser.Parse(message)
			Expect(log.Severity()).To(Equal(2))
		})

		It("returns the facility", func() {
			log, _, _ := parser.Parse(message)
			Expect(log.Facility()).To(Equal(4))
		})

		It("sets the priority", func() {
			log, _, _ := parser.Parse(message)
			Expect(log.Priority()).To(Equal(34))
		})
	})

	It("returns a valid date object", func() {
		log, _, _ := parser.Parse(message)
		Expect(log.Timestamp().String()).To(Equal("2003-10-11 22:14:15.003 +0000 UTC"))
	})

	It("sets the hostname", func() {
		log, _, _ := parser.Parse(message)
		Expect(log.Hostname()).To(BeEquivalentTo("mymachine.example.com"))
	})

	It("sets the app name", func() {
		log, _, _ := parser.Parse(message)
		Expect(log.Appname()).To(BeEquivalentTo("su"))
	})

	It("sets the proc id", func() {
		log, _, _ := parser.Parse(message)
		Expect(log.ProcID()).To(BeEquivalentTo("12345"))
	})

	It("sets the msg id", func() {
		log, _, _ := parser.Parse(message)
		Expect(log.MsgID()).To(BeEquivalentTo("98765"))
	})

	It("sets structure data", func() {
		log, _, _ := parser.Parse(message)
		Expect(log.StructureData()).To(HaveLen(1))
		data := log.StructureData()[0]
		Expect(data.ID()).To(BeEquivalentTo("exampleSDID@32473"))
		Expect(data.Properties()).To(BeEquivalentTo(payload.Properties{
			payload.NewProperty("iut", "3"),
			payload.NewProperty("eventSource", "Application"),
			payload.NewProperty("eventID", "1011"),
		}))
	})

	It("sets the message", func() {
		log, _, _ := parser.Parse(message)
		Expect(log.Message()).To(BeEquivalentTo("'su root' failed for lonvick on /dev/pts/8"))
	})
}

var _ = Describe("with a standard payload", func() {
	It("returns the correct offset", func() {
		_, offset, _ := parser.Parse(validMessage)
		Expect(offset).To(Equal(len(validMessage)))
	})

	ParseLogMessageTests(validMessage)

	Context("with the hostname", func() {
		It("is empty when '-'", func() {
			message := `<34>1 2003-10-11T22:14:15.003Z - su - - - 'su root' failed for lonvick on /dev/pts/8`
			log, _, _ := parser.Parse(message)
			Expect(log.Hostname()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(message))
		})
	})

	Context("with the app name", func() {
		It("is empty when '-'", func() {
			message := `<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8`
			log, _, _ := parser.Parse(message)
			Expect(log.Appname()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(message))
		})
	})

	Context("with the proc id", func() {
		It("is empty when '-'", func() {
			message := `<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8`
			log, _, _ := parser.Parse(message)
			Expect(log.ProcID()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(message))
		})
	})

	Context("with the log id", func() {
		It("is empty when '-'", func() {
			message := `<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8`
			log, _, _ := parser.Parse(message)
			Expect(log.MsgID()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(message))
		})
	})

	Context("with structure data", func() {
		It("is empty when '-'", func() {
			message := `<34>1 2003-10-11T22:14:15.003Z - - - - - 'su root' failed for lonvick on /dev/pts/8`
			log, _, _ := parser.Parse(message)
			Expect(log.StructureData()).To(HaveLen(0))
			Expect(log.String()).To(BeEquivalentTo(message))
		})

		It("allows character escaping", func() {
			message := `<29>50 2016-01-15T01:00:43Z hn S - - [my@id1 a="1" b="\"" c="\\" d="\]" e="\"There are \\many things here[1\]\""]`
			log, _, err := parser.Parse(message)
			Expect(err).ToNot(HaveOccurred())
			Expect(log.StructureData()).To(HaveLen(1))
			data := log.StructureData()[0]
			Expect(data.ID()).To(BeEquivalentTo("my@id1"))
			Expect(data.Properties()).To(BeEquivalentTo([]payload.Property{
				payload.NewProperty("a", "1"),
				payload.NewProperty("b", `"`),
				payload.NewProperty("c", `\`),
				payload.NewProperty("d", `]`),
				payload.NewProperty("e", `"There are \many things here[1]"`),
			}))
			Expect(log.String()).To(BeEquivalentTo(message))
		})
	})

	Context("with a message", func() {
		It("sets empty for no message", func() {
			message := `<34>1 2003-10-11T22:14:15.003Z - - - - -`
			log, _, _ := parser.Parse(message)
			Expect(log.Message()).To(BeEmpty())
			Expect(log.String()).To(BeEquivalentTo(message))
		})
	})

	Context("with multiple structure data elements", func() {
		It("parses them all", func() {
			message := `<29>50 2016-01-15T01:00:43Z hn S - - [my@id1 k="v"][my@id2 c="val"]`
			log, _, err := parser.Parse(message)
			Expect(err).ToNot(HaveOccurred())
			Expect(log.StructureData()).To(HaveLen(2))
			data := log.StructureData()[0]
			Expect(data.ID()).To(BeEquivalentTo("my@id1"))
			Expect(data.Properties()).To(BeEquivalentTo(payload.Properties{
				payload.NewProperty("k", "v"),
			}))
			data = log.StructureData()[1]
			Expect(data.ID()).To(BeEquivalentTo("my@id2"))
			Expect(data.Properties()).To(BeEquivalentTo(payload.Properties{
				payload.NewProperty("c", "val"),
			}))
		})
	})

	Context("with the example timestamps from the RFC", func() {
		DescribeTable("can parse them all", func(value, expected string) {
			message := fmt.Sprintf(`<34>1 %s mymachine.example.com su - - - 'su root' failed for lonvick on /dev/pts/8`, value)
			log, offset, err := parser.Parse(message)
			Expect(offset).To(Equal(len(message)))
			Expect(err).ToNot(HaveOccurred())
			Expect(log.Timestamp().String()).To(Equal(expected))
		},
			Entry("with nanoseconds", "2003-10-11T22:14:15.00003Z", "2003-10-11 22:14:15.00003 +0000 UTC"),
			Entry("with partial seconds", "1985-04-12T23:20:50.52Z", "1985-04-12 23:20:50.52 +0000 UTC"),
			Entry("with no timezone offset", "1985-04-12T23:20:50.52+00:00", "1985-04-12 23:20:50.52 +0000 UTC"),
			Entry("with a positive timezone offset", "1985-04-12T23:20:50.52+02:00", "1985-04-12 21:20:50.52 +0000 UTC"),
			Entry("with a negative timezone offset", "1985-04-12T18:20:50.52-02:00", "1985-04-12 20:20:50.52 +0000 UTC"),
		)

		It("fails parsing on unsupported formats", func() {
			message := `<34>1 2003-10-11T22:14:15.003Z07:00 mymachine.example.com su - - - 'su root' failed for lonvick on /dev/pts/8`
			_, n, err := parser.Parse(message)
			Expect(n).To(Equal(30), "the timestamp is length is incorrect")
			Expect(err).To(HaveOccurred())
		})

		It("is sets empty when '-'", func() {
			message := `<34>1 - - su - - - 'su root' failed for lonvick on /dev/pts/8`
			log, _, _ := parser.Parse(message)
			Expect(log.Timestamp().IsZero()).To(BeTrue())
			Expect(log.String()).To(BeEquivalentTo(message))
		})
	})
})

var _ = Describe("with a TCP payload", func() {
	Context("with a single message", func() {
		message := fmt.Sprintf("%d %s", len(validMessage), validMessage)

		It("returns the correct offset", func() {
			_, offset, _ := parser.Parse(message)
			Expect(offset).To(Equal(len(validMessage) + 4))
		})

		ParseLogMessageTests(message)
	})

	Context("with a multiple payload", func() {
		message := fmt.Sprintf("%d %s%d %s", len(validMessage), validMessage, len(validMessage), validMessage)

		It("returns the correct offset", func() {
			_, offset, _ := parser.Parse(message)
			Expect(offset).To(Equal(len(validMessage) + 4))
		})
		ParseLogMessageTests(message)
	})
})
