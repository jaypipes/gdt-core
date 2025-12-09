// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.

package api

import (
	"encoding/xml"
	"time"
)

// XUnitResults is a wrapper struct allowing the output of well-formed
// JUnit/XUnit XML and JSON serialized scenario results
type XUnitResults struct {
	XMLName    xml.Name         `json:"-" xml:"testsuites"`
	TestSuites []XUnitTestSuite `json:"testsuites" xml:"testsuite"`
}

type XUnitTestSuite struct {
	Name       string          `json:"name" xml:"name,attr"`
	Time       string          `json:"time,omitempty" xml:"time,attr"`
	Skipped    bool            `json:"skipped,omitempty" xml:"skipped,omitempty"`
	Failures   int             `json:"failures" xml:"failures,attr"`
	Errors     int             `json:"errors" xml:"errors,attr"`
	Tests      int             `json:"tests" xml:"tests,attr"`
	Timestamp  time.Time       `json:"timestamp" xml:"timestamp,attr"`
	Properties []XUnitProperty `json:"properties,omitempty" xml:"properties,omitempty"`
	TestCases  []XUnitTestCase `json:"testcases,omitempty" xml:"testcases,omitempty"`
}

type XUnitProperty struct {
	Name  string `json:"name" xml:"name,attr"`
	Value string `json:"value" xml:"value,attr"`
}

type XUnitTestCase struct {
	Name       string          `json:"name" xml:"name,attr"`
	Time       string          `json:"time,omitempty" xml:"time,attr"`
	Status     string          `json:"status" xml:"status,attr"`
	Assertions int             `json:"assertions" xml:"assertions,attr,omitempty"`
	Skipped    bool            `json:"skipped,omitempty" xml:"skipped,omitempty"`
	Failure    *XUnitMessage   `json:"failure,omitempty" xml:"failure,omitempty"`
	Error      *XUnitMessage   `json:"error,omitempty" xml:"error,omitempty"`
	SystemOut  *XUnitTextBlock `json:"system-out,omitempty" xml:"system-out,omitempty"`
	SystemErr  *XUnitTextBlock `json:"system-err,omitempty" xml:"system-err,omitempty"`
}

type XUnitTextBlock struct {
	Text string `xml:",cdata"`
}

type XUnitMessage struct {
	Message string `json:"message,omitempty" xml:"message,attr"`
	Type    string `json:"type,omitempty" xml:"type,attr"`
}
