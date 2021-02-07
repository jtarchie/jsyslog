package parser_test

import (
	"fmt"
	"github.com/jtarchie/jsyslog/parser"
	"github.com/jtarchie/jsyslog/payload"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// This is here to avoid compiler optimizations that
// could remove the actual call we are benchmarking
// during benchmarks
var benchParseResult *payload.Message

type benchCase struct {
	input string
	label string
	valid bool
}

func rxpad(str string, lim int) string {
	str += strings.Repeat(" ", lim)
	return str[:lim]
}

var benchCases = []benchCase{
	{
		label: "[no] empty input",
		input: ``,
		valid: false,
	},
	{
		label: "[no] multiple syslog messages on multiple lines",
		input: "<1>1 - - - - - -\x0A<2>1 - - - - - -",
		valid: false,
	},
	{
		label: "[no] impossible timestamp",
		input: `<101>11 2003-09-31T22:14:15.003Z`,
		valid: false,
	},
	{
		label: "[no] malformed structured data",
		input: "<1>1 - - - - - X",
		valid: false,
	},
	// {
	// 	label: "[no] with duplicated structured data id",
	// 	input: "<165>3 2003-10-11T22:14:15.003Z example.com evnts - ID27 [id1][id1]",
	// 	valid: false,
	// },
	{
		label: "[ok] minimal",
		input: `<1>1 - - - - - -`,
		valid: true,
	},
	{
		label: "[ok] average message",
		input: `<29>1 2016-02-21T04:32:57Z web1 someservice - - [origin x-service="someservice"][meta sequenceId="14125553"] 127.0.0.1 - - 1456029177 "GET /v1/ok HTTP/1.1" 200 145 "-" "hacheck 0.9.0" 24306 127.0.0.1:40124 575`,
		valid: true,
	},
	{
		label: "[ok] complicated message",
		input: `<78>1 2016-01-15T00:04:01Z host1 CROND 10391 - [meta sequenceId="29" sequenceBlah="foo"][my key="value"] some_message`,
		valid: true,
	},
	{
		label: "[ok] very long message",
		input: `<190>1 2016-02-21T01:19:11Z batch6sj - - - [meta sequenceId="21881798" x-group="37051387"][origin x-service="tracking"] metascutellar conversationalist nephralgic exogenetic graphy streng outtaken acouasm amateurism prenotice Lyonese bedull antigrammatical diosphenol gastriloquial bayoneteer sweetener naggy roughhouser dighter addend sulphacid uneffectless ferroprussiate reveal Mazdaist plaudite Australasian distributival wiseman rumness Seidel topazine shahdom sinsion mesmerically pinguedinous ophthalmotonometer scuppler wound eciliate expectedly carriwitchet dictatorialism bindweb pyelitic idic atule kokoon poultryproof rusticial seedlip nitrosate splenadenoma holobenthic uneternal Phocaean epigenic doubtlessly indirection torticollar robomb adoptedly outspeak wappenschawing talalgia Goop domitic savola unstrafed carded unmagnified mythologically orchester obliteration imperialine undisobeyed galvanoplastical cycloplegia quinquennia foremean umbonal marcgraviaceous happenstance theoretical necropoles wayworn Igbira pseudoangelic raising unfrounced lamasary centaurial Japanolatry microlepidoptera`,
		valid: true,
	},
	{
		label: "[ok] all max length and complete",
		input: `<191>999 2019-01-01T23:58:59.999999Z abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabc abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdef abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzab abcdefghilmnopqrstuvzabcdefghilm [an@id key1="val1" key2="val2"][another@id key1="val1"] Some message "GET"`,
		valid: true,
	},
	{
		label: "[ok] all max length except structured data and message",
		input: `<191>999 2019-01-01T23:58:59.999999Z abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabc abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdef abcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzabcdefghilmnopqrstuvzab abcdefghilmnopqrstuvzabcdefghilm -`,
		valid: true,
	},
	{
		label: "[ok] minimal with message containing newline",
		input: "<1>1 - - - - - - x\x0Ay",
		valid: true,
	},
	{
		label: "[ok] w/o procid, w/o structured data, with message starting with BOM",
		input: "<34>1 2003-10-11T22:14:15.003Z mymachine.example.com su - ID47 - \xEF\xBB\xBF'su root' failed for lonvick on /dev/pts/8",
		valid: true,
	},
	{
		label: "[ok] minimal with UTF-8 message",
		input: "<0>1 - - - - - - ⠊⠀⠉⠁⠝⠀⠑⠁⠞⠀⠛⠇⠁⠎⠎⠀⠁⠝⠙⠀⠊⠞⠀⠙⠕⠑⠎⠝⠞⠀⠓⠥⠗⠞⠀⠍⠑",
		valid: true,
	},
	{
		label: "[ok] with structured data id, w/o structured data params",
		input: `<29>50 2016-01-15T01:00:43Z hn S - - [my@id]`,
		valid: true,
	},
	{
		label: "[ok] with multiple structured data",
		input: `<29>50 2016-01-15T01:00:43Z hn S - - [my@id1 k="v"][my@id2 c="val"]`,
		valid: true,
	},
	{
		label: "[ok] with escaped backslash within structured data param value, with message",
		input: `<29>50 2016-01-15T01:00:43Z hn S - - [meta es="\\valid"] 1452819643`,
		valid: true,
	},
	{
		label: "[ok] with UTF-8 structured data param value, with message",
		input: `<78>1 2016-01-15T00:04:01Z host1 CROND 10391 - [sdid x="⌘"] some_message`,
		valid: true,
	},
}

func BenchmarkParse(b *testing.B) {
	for _, tc := range benchCases {
		tc := tc
		b.Run(rxpad(tc.label, 50), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchParseResult, _, _ = parser.Parse(tc.input)
			}
		})
	}
}

var _ = Describe("Ensure benchmarks are valid", func() {
	for _, tc := range benchCases {
		tc := tc
		It(fmt.Sprintf("tries to parse a message '%s'", tc.label), func() {
			result, offset, err := parser.Parse(tc.input)
			if tc.valid {
				Expect(err).ToNot(HaveOccurred())
				Expect(result).ToNot(BeNil())
				Expect(offset).To(Equal(len(tc.input)))
				Expect(result.String()).To(BeEquivalentTo(tc.input))
			} else {
				Expect(err).To(HaveOccurred())
			}
		})
	}
})
