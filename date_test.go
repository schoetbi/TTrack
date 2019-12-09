package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dateTest struct {
	in, out string
	err     bool
}

var testInputs = []dateTest{

	//  mm/dd/yyyy
	{in: "03/31/2014", out: "2014-03-31 00:00:00 +0000 UTC"},
	{in: "3/31/2014", out: "2014-03-31 00:00:00 +0000 UTC"},
	{in: "3/5/2014", out: "2014-03-05 00:00:00 +0000 UTC"},

	//  mm/dd/yy
	{in: "08/08/71", out: "1971-08-08 00:00:00 +0000 UTC"},
	{in: "8/8/71", out: "1971-08-08 00:00:00 +0000 UTC"},
	//  mm/dd/yy hh:mm:ss
	{in: "04/02/2014 04:08:09", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "4/2/2014 04:08:09", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "04/02/2014 4:08:09", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "04/02/2014 4:8:9", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "04/02/2014 04:08", out: "2014-04-02 04:08:00 +0000 UTC"},
	{in: "04/02/2014 4:8", out: "2014-04-02 04:08:00 +0000 UTC"},
	//  mm/dd/yy hh:mm:ss AM
	{in: "04/02/2014 04:08:09 AM", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "04/02/2014 04:08:09 PM", out: "2014-04-02 16:08:09 +0000 UTC"},
	{in: "04/02/2014 04:08 AM", out: "2014-04-02 04:08:00 +0000 UTC"},
	{in: "04/02/2014 04:08 PM", out: "2014-04-02 16:08:00 +0000 UTC"},
	{in: "04/02/2014 4:8 AM", out: "2014-04-02 04:08:00 +0000 UTC"},
	{in: "04/02/2014 4:8 PM", out: "2014-04-02 16:08:00 +0000 UTC"},
	//   yyyy/mm/dd
	{in: "2014/04/02", out: "2014-04-02 00:00:00 +0000 UTC"},
	{in: "2014/03/31", out: "2014-03-31 00:00:00 +0000 UTC"},
	{in: "2014/4/2", out: "2014-04-02 00:00:00 +0000 UTC"},
	//   yyyy/mm/dd hh:mm:ss AM
	{in: "2014/04/02 04:08", out: "2014-04-02 04:08:00 +0000 UTC"},
	{in: "2014/03/31 04:08", out: "2014-03-31 04:08:00 +0000 UTC"},
	{in: "2014/4/2 04:08", out: "2014-04-02 04:08:00 +0000 UTC"},
	{in: "2014/04/02 4:8", out: "2014-04-02 04:08:00 +0000 UTC"},
	{in: "2014/04/02 04:08:09", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "2014/03/31 04:08:09", out: "2014-03-31 04:08:09 +0000 UTC"},
	{in: "2014/4/2 04:08:09", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "2014/04/02 04:08:09 AM", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "2014/03/31 04:08:09 AM", out: "2014-03-31 04:08:09 +0000 UTC"},
	{in: "2014/4/2 04:08:09 AM", out: "2014-04-02 04:08:09 +0000 UTC"},
	//   yyyy-mm-dd
	{in: "2014-04-02", out: "2014-04-02 00:00:00 +0000 UTC"},
	{in: "2014-03-31", out: "2014-03-31 00:00:00 +0000 UTC"},
	{in: "2014-4-2", out: "2014-04-02 00:00:00 +0000 UTC"},

	{in: "2014-4-2 04:08:09", out: "2014-04-02 04:08:09 +0000 UTC"},
	//   yyyy-mm-dd hh:mm:ss AM
	{in: "2014-04-02 04:08", out: "2014-04-02 04:08:00 +0000 UTC"},
	{in: "2014-03-31 04:08", out: "2014-03-31 04:08:00 +0000 UTC"},
	{in: "2014-4-2 04:08", out: "2014-04-02 04:08:00 +0000 UTC"},
	{in: "2014-04-02 4:8", out: "2014-04-02 04:08:00 +0000 UTC"},
	{in: "2014-04-02 04:08:09", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "2014-03-31 04:08:09", out: "2014-03-31 04:08:09 +0000 UTC"},
	{in: "2014-4-2 04:08:09", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "2014-04-02 04:08:09 AM", out: "2014-04-02 04:08:09 +0000 UTC"},
	{in: "2014-03-31 04:08:09 AM", out: "2014-03-31 04:08:09 +0000 UTC"},
	{in: "2014-04-26 05:24:37 PM", out: "2014-04-26 17:24:37 +0000 UTC"},
	{in: "2014-4-2 04:08:09 AM", out: "2014-04-02 04:08:09 +0000 UTC"},
	// yyyy.mm.dd
	{in: "2018.09.30", out: "2018-09-30 00:00:00 +0000 UTC"},

	//   dd.mm.yyyy
	{in: "31.3.2014", out: "2014-03-31 00:00:00 +0000 UTC"},
	{in: "3.4.2014", out: "2014-04-03 00:00:00 +0000 UTC"},
	{in: "31.03.2014", out: "2014-03-31 00:00:00 +0000 UTC"},
	{in: "31.3.2014 2:15:09", out: "2014-03-31 02:15:09 +0000 UTC"},
	{in: "31.3.2014 2:15", out: "2014-03-31 02:15:00 +0000 UTC"},
	{in: "31.3.14 02:15", out: "2014-03-31 02:15:00 +0000 UTC"},
	{in: "31.3.2014 02:15", out: "2014-03-31 02:15:00 +0000 UTC"},
}

func TestParse(t *testing.T) {

	// Lets ensure we are operating on UTC
	time.Local = time.UTC

	for _, th := range testInputs {
		ts, _ := ParseDateTime(th.in)
		got := fmt.Sprintf("%v", ts.In(time.UTC))
		assert.Equal(t, th.out, got, "Expected %q but got %q from %q", th.out, got, th.in)
	}
}

/*
var testParseErrors = []dateTest{
	{in: "3", err: true},
	{in: `{"hello"}`, err: true},
	{in: "2009-15-12T22:15Z", err: true},
	{in: "5,000-9,999", err: true},
	{in: "xyzq-baad"},
	{in: "oct.-7-1970", err: true},
	{in: "septe. 7, 1970", err: true},
	{in: "SeptemberRR 7th, 1970", err: true},
	{in: "29-06-2016", err: true},
	// this is just testing the empty space up front
	{in: " 2018-01-02 17:08:09 -07:00", err: true},
}

func TestParseErrors(t *testing.T) {
	for _, th := range testParseErrors {
		v, err := ParseDateTime(th.in)
		assert.NotEqual(t, nil, err, "%v for %v", v, th.in)
	}
}

var testParseFormat = []dateTest{
	// errors
	{in: "3", err: true},
	{in: `{"hello"}`, err: true},
	{in: "2009-15-12T22:15Z", err: true},
	{in: "5,000-9,999", err: true},
	//
	{in: "oct 7, 1970", out: "Jan 2, 2006"},
	{in: "sept. 7, 1970", out: "Jan. 2, 2006"},
	{in: "May 05, 2015, 05:05:07", out: "Jan 02, 2006, 15:04:05"},
	// 03 February 2013
	{in: "03 February 2013", out: "02 January 2006"},
	// 13:31:51.999 -07:00 MST
	//   yyyy-mm-dd hh:mm:ss +00:00
	{in: "2012-08-03 18:31:59 +00:00", out: "2006-01-02 15:04:05 -07:00"},
	//   yyyy-mm-dd hh:mm:ss +0000 TZ
	// Golang Native Format
	{in: "2012-08-03 18:31:59 +0000 UTC", out: "2006-01-02 15:04:05 -0700 UTC"},
	//   yyyy-mm-dd hh:mm:ss TZ
	{in: "2012-08-03 18:31:59 UTC", out: "2006-01-02 15:04:05 UTC"},
	//   yyyy-mm-ddThh:mm:ss-07:00
	{in: "2009-08-12T22:15:09-07:00", out: "2006-01-02T15:04:05-07:00"},
	//   yyyy-mm-ddThh:mm:ss-0700
	{in: "2009-08-12T22:15:09-0700", out: "2006-01-02T15:04:05-0700"},
	//   yyyy-mm-ddThh:mm:ssZ
	{in: "2009-08-12T22:15Z", out: "2006-01-02T15:04Z"},
}

func TestParseLayout(t *testing.T) {
	for _, th := range testParseFormat {
		l, err := ParseDateTime(th.in)
		if th.err {
			assert.NotEqual(t, nil, err)
		} else {
			assert.Equal(t, nil, err)
			assert.Equal(t, th.out, l, "for in=%v", th.in)
		}
	}
}

var testParseStrict = []dateTest{
	//   dd-mon-yy  13-Feb-03
	{in: "03-03-14"},
	//   mm.dd.yyyy
	{in: "3.3.2014"},
	//   mm.dd.yy
	{in: "08.09.71"},
	//  mm/dd/yyyy
	{in: "3/5/2014"},
	//  mm/dd/yy
	{in: "08/08/71"},
	{in: "8/8/71"},
	//  mm/dd/yy hh:mm:ss
	{in: "04/02/2014 04:08:09"},
	{in: "4/2/2014 04:08:09"},
}
*/
