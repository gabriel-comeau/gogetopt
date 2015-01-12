package gogetopt

import "testing"

// Test all of the stateless functions (that don't depend on the package-wide opt maps / indexes)

// Test the lookaheadForVal() function
func TestLookAhead(t *testing.T) {

	const BS_ARGV0 = "ignoreme"
	const SHORT = "-f"
	const LONG = "--foo"
	const VAL_STR = "wow"

	// Correct vals
	var fakeArgs []string

	fakeArgs = []string{BS_ARGV0, SHORT, VAL_STR}
	shortWithVal := lookaheadForVal(fakeArgs, 1)
	if shortWithVal != VAL_STR {
		t.Error("lookaheadForVal() return invalid value for short with val: " + shortWithVal)
	}

	fakeArgs = []string{BS_ARGV0, LONG, VAL_STR}
	longWithVal := lookaheadForVal(fakeArgs, 1)
	if longWithVal != VAL_STR {
		t.Error("lookaheadForVal() return invalid value for long with val: " + longWithVal)
	}

	// Incorrect (missing)
	fakeArgs = []string{BS_ARGV0, SHORT}
	shortNoVal := lookaheadForVal(fakeArgs, 1)
	if shortNoVal != "" {
		t.Error("lookaheadForVal() returned value when there should not have been: " + shortNoVal)
	}

	fakeArgs = []string{BS_ARGV0, LONG}
	longNoVal := lookaheadForVal(fakeArgs, 1)
	if longNoVal != "" {
		t.Error("lookaheadForVal() returned value when there should not have been: " + longNoVal)
	}

}

// Test the stripDashes() function
func TestStripDashes(t *testing.T) {
	const SHORT_WITH_DASH = "-f"
	const SHORT_WITHOUT_DASH = "f"
	const LONG_WITH_DASHES = "--foo"
	const LONG_WITHOUT_DASHES = "foo"
	const NO_DASHES = "nodashes"

	strippedShortWithDashes := stripDashes(SHORT_WITH_DASH)
	if strippedShortWithDashes != SHORT_WITHOUT_DASH {
		t.Error("stripDashes() didn't strip the dash off of: " + SHORT_WITH_DASH + ".  Got: " + strippedShortWithDashes)
	}

	strippedLongWithDashes := stripDashes(LONG_WITH_DASHES)
	if strippedLongWithDashes != LONG_WITHOUT_DASHES {
		t.Error("stripDashes() didn't strip the dashes off of: " + LONG_WITH_DASHES + ".  Got: " + strippedLongWithDashes)
	}

	strippedNoDashes := stripDashes(NO_DASHES)
	if strippedNoDashes != NO_DASHES {
		t.Error("stripDashes() should have returned: " + NO_DASHES + ".  Got: " + strippedNoDashes)
	}
}

// Test the splitEqualsArg() function
func TestSplitEqualsArg(t *testing.T) {
	const SHORT_DASH_EQUALS_VAL = "-f=wow" //good
	const SHORT_EQUALS_VAL = "f=wow"       //good
	const SHORT_GOOD_PART_0 = "f"
	const SHORT_GOOD_PART_1 = "wow"

	const SHORT_DASH_EQUALS_NOVAL = "-f=" //bad
	const SHORT_EQUALS_NOVAL = "f="       //bad

	const LONG_DASHES_EQUALS_VAL = "--foo=wow" //good
	const LONG_EQUALS_VAL = "foo=wow"          //good
	const LONG_GOOD_PART_0 = "foo"
	const LONG_GOOD_PART_1 = "wow"

	const LONG_DASHES_EQUALS_NOVAL = "--foo=" //bad
	const LONG_EQUALS_NOVAL = "foo="          //bad

	const NOEQUALS = "wow" //bad

	// Check good ones first

	shortDashEqualsAndVal := splitEqualsArg(SHORT_DASH_EQUALS_VAL)
	if len(shortDashEqualsAndVal) != 2 {
		t.Error("splitEqualsArg() couldn't split: " + SHORT_DASH_EQUALS_VAL)
	}
	if shortDashEqualsAndVal[0] != SHORT_GOOD_PART_0 {
		t.Error("splitEqualsArg() returned incorrect value for part[0].  Expecting: " + SHORT_GOOD_PART_0 + ".  Got: " + shortDashEqualsAndVal[0])
	}
	if shortDashEqualsAndVal[1] != SHORT_GOOD_PART_1 {
		t.Error("splitEqualsArg() returned incorrect value for part[1].  Expecting: " + SHORT_GOOD_PART_1 + ".  Got: " + shortDashEqualsAndVal[1])
	}

	shortEqualsAndVal := splitEqualsArg(SHORT_EQUALS_VAL)
	if len(shortEqualsAndVal) != 2 {
		t.Error("splitEqualsArg() couldn't split: " + SHORT_EQUALS_VAL)
	}
	if shortEqualsAndVal[0] != SHORT_GOOD_PART_0 {
		t.Error("splitEqualsArg() returned incorrect value for part[0].  Expecting: " + SHORT_GOOD_PART_0 + ".  Got: " + shortEqualsAndVal[0])
	}
	if shortEqualsAndVal[1] != SHORT_GOOD_PART_1 {
		t.Error("splitEqualsArg() returned incorrect value for part[1].  Expecting: " + SHORT_GOOD_PART_1 + ".  Got: " + shortEqualsAndVal[1])
	}

	longDashEqualsAndVal := splitEqualsArg(LONG_DASHES_EQUALS_VAL)
	if len(longDashEqualsAndVal) != 2 {
		t.Error("splitEqualsArg() couldn't split: " + LONG_DASHES_EQUALS_VAL)
	}
	if longDashEqualsAndVal[0] != LONG_GOOD_PART_0 {
		t.Error("splitEqualsArg() returned incorrect value for part[0].  Expecting: " + LONG_GOOD_PART_0 + ".  Got: " + longDashEqualsAndVal[0])
	}
	if longDashEqualsAndVal[1] != LONG_GOOD_PART_1 {
		t.Error("splitEqualsArg() returned incorrect value for part[1].  Expecting: " + LONG_GOOD_PART_1 + ".  Got: " + longDashEqualsAndVal[1])
	}

	longEqualsAndVal := splitEqualsArg(LONG_EQUALS_VAL)
	if len(longEqualsAndVal) != 2 {
		t.Error("splitEqualsArg() couldn't split: " + LONG_EQUALS_VAL)
	}
	if longEqualsAndVal[0] != LONG_GOOD_PART_0 {
		t.Error("splitEqualsArg() returned incorrect value for part[0].  Expecting: " + LONG_GOOD_PART_0 + ".  Got: " + longEqualsAndVal[0])
	}
	if longEqualsAndVal[1] != LONG_GOOD_PART_1 {
		t.Error("splitEqualsArg() returned incorrect value for part[1].  Expecting: " + LONG_GOOD_PART_1 + ".  Got: " + longEqualsAndVal[1])
	}

	// Now various bad ones
	if splitEqualsArg(SHORT_EQUALS_NOVAL) != nil {
		t.Error("splitEqualsArg() didn't return nil for argument without value: " + SHORT_EQUALS_NOVAL)
	}

	if splitEqualsArg(SHORT_DASH_EQUALS_NOVAL) != nil {
		t.Error("splitEqualsArg() didn't return nil for argument without value: " + SHORT_DASH_EQUALS_NOVAL)
	}

	if splitEqualsArg(LONG_EQUALS_NOVAL) != nil {
		t.Error("splitEqualsArg() didn't return nil for argument without value: " + LONG_EQUALS_NOVAL)
	}

	if splitEqualsArg(LONG_DASHES_EQUALS_NOVAL) != nil {
		t.Error("splitEqualsArg() didn't return nil for argument without value: " + LONG_DASHES_EQUALS_NOVAL)
	}

	if splitEqualsArg(NOEQUALS) != nil {
		t.Error("splitEqualsArg() didn't return nil for argument without equals sign: " + NOEQUALS)
	}
}
