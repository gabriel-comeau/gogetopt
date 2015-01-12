package gogetopt

import (
	"os"
	"regexp"
	"testing"
)

// Test out a valid, boolean option like -b
func TestShortBool(t *testing.T) {
	ClearAll()
	regErr := RegisterOpt("test", "", "l", true, false, "test usage")
	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-l"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}

	val := GetBool("test")
	if val != true {
		t.Error("Parse() test: Didn't get a positive boolean when expecting one")
	}
}

// Test out a valid, boolean option like --long
func TestLongBool(t *testing.T) {
	ClearAll()
	regErr := RegisterOpt("test", "test", "", true, false, "test usage")
	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "--test"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}

	val := GetBool("test")
	if val != true {
		t.Error("Parse() test: Didn't get a positive boolean when expecting one")
	}
}

// Test out several valid boolean shorts, not combined: -t -u -v
func TestNonCombinedShortBools(t *testing.T) {
	ClearAll()
	var regErr error
	regErr = RegisterOpt("test1", "", "t", true, false, "test usage")
	regErr = RegisterOpt("test2", "", "u", true, false, "test usage")
	regErr = RegisterOpt("test3", "", "v", true, false, "test usage")

	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-t", "foo", "-u", "-v"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}

	val1 := GetBool("test1")
	val2 := GetBool("test2")
	val3 := GetBool("test3")
	if !val1 || !val2 || !val3 {
		t.Error("Parse() test: Didn't get a positive boolean when expecting one")
	}
}

// Test out several combined valid boolean shorts: -tuv
func TestCombinedShortBools(t *testing.T) {
	ClearAll()
	var regErr error
	regErr = RegisterOpt("test1", "", "t", true, false, "test usage")
	regErr = RegisterOpt("test2", "", "u", true, false, "test usage")
	regErr = RegisterOpt("test3", "", "v", true, false, "test usage")

	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-tuv", "foo"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}

	val1 := GetBool("test1")
	val2 := GetBool("test2")
	val3 := GetBool("test3")
	if !val1 || !val2 || !val3 {
		t.Error("Parse() test: Didn't get a positive boolean when expecting one")
	}
}

// Mix long, short, and comb. short bools together, valid
func TestMixedLongShortBools(t *testing.T) {
	ClearAll()
	var regErr error
	regErr = RegisterOpt("test1", "", "t", true, false, "test usage")
	regErr = RegisterOpt("test2", "", "u", true, false, "test usage")
	regErr = RegisterOpt("test3", "", "v", true, false, "test usage")
	regErr = RegisterOpt("test4", "wow", "", true, false, "test usage")

	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-t", "--wow", "-uv", "foo"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}

	val1 := GetBool("test1")
	val2 := GetBool("test2")
	val3 := GetBool("test3")
	val4 := GetBool("test4")
	if !val1 || !val2 || !val3 || !val4 {
		t.Error("Parse() test: Didn't get a positive boolean when expecting one")
	}
}

//Get string vals from = signs for both short and long
func TestStringEqualsVals(t *testing.T) {
	const VAL1 = "val1"
	const VAL2 = "val2"

	ClearAll()
	var regErr error
	regErr = RegisterOpt("test1", "", "t", false, false, "test usage")
	regErr = RegisterOpt("test2", "wow", "", false, false, "test usage")

	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-t=" + VAL1, "--wow=" + VAL2, "foo"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}

	val1 := GetString("test1")
	val2 := GetString("test2")

	if val1 != VAL1 {
		t.Error("Parse() test: Didn't get a string val for short equals expr.  Got: " + val1 + ".  Expected: " + VAL1)
	}

	if val2 != VAL2 {
		t.Error("Parse() test: Didn't get a string val for long equals expr.  Got: " + val2 + ".  Expected: " + VAL2)
	}
}

//Get string vals from the lookahead form: -f foo (or --foo bar)
func TestLookaheadStrings(t *testing.T) {
	const VAL1 = "val1"
	const VAL2 = "val2"

	ClearAll()
	var regErr error
	regErr = RegisterOpt("test1", "", "t", false, false, "test usage")
	regErr = RegisterOpt("test2", "wow", "", false, false, "test usage")

	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-t", VAL1, "--wow", VAL2, "foo"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}

	val1 := GetString("test1")
	val2 := GetString("test2")

	if val1 != VAL1 {
		t.Error("Parse() test: Didn't get a string val for short lookahead expr.  Got: " + val1 + ".  Expected: " + VAL1)
	}

	if val2 != VAL2 {
		t.Error("Parse() test: Didn't get a string val for long lookahead expr.  Got: " + val2 + ".  Expected: " + VAL2)
	}
}

// Check the -tVAL form for short string vals.
func TestSmushedShortString(t *testing.T) {
	const VAL1 = "val1"

	ClearAll()
	var regErr error
	regErr = RegisterOpt("test1", "", "t", false, false, "test usage")

	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-t" + VAL1, "foo"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}

	val1 := GetString("test1")

	if val1 != VAL1 {
		t.Error("Parse() test: Didn't get a string val for short smushed expr.  Got: " + val1 + ".  Expected: " + VAL1)
	}
}

// Check for presence of non opt values insterspersed through the args.
func TestExtraVals(t *testing.T) {
	ClearAll()
	var regErr error
	regErr = RegisterOpt("test1", "", "a", false, false, "test usage")
	regErr = RegisterOpt("test2", "", "b", false, false, "test usage")
	regErr = RegisterOpt("test3", "", "c", true, false, "test usage")
	regErr = RegisterOpt("test4", "wow", "", true, false, "test usage")

	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-a=foo", "firstextra", "-b", "bar", "secondextra", "-c", "--wow", "lastextra"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}

	if len(GetArgs()) != 3 {
		t.Errorf("Parse() test: Didn't get 3 extra args, got: %+v\n", GetArgs())
	}
}

// Check for required options.  Make sure all the different types of string arguments
// which when set to registered actually count the arguments as included.
func TestReqAllPresent(t *testing.T) {
	ClearAll()
	var regErr error
	regErr = RegisterOpt("test1", "", "a", false, true, "test usage")
	regErr = RegisterOpt("test2", "wow", "", false, true, "test usage")
	regErr = RegisterOpt("test3", "", "b", false, true, "test usage")
	regErr = RegisterOpt("test4", "", "c", false, true, "test usage")
	regErr = RegisterOpt("test5", "hey", "", false, true, "test usage")

	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-a=foo", "--wow", "suchval", "-bbar", "-c", "baz", "--hey=yikes", "meh"}

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		t.Error("Parse() test: Got a parse error: " + GetError().Error())
	}
}

// Make sure when missing arguments are missing a parse error is thrown.
func TestReqMissing(t *testing.T) {
	ClearAll()
	var regErr error
	regErr = RegisterOpt("test1", "", "a", false, true, "test usage")
	regErr = RegisterOpt("test2", "wow", "", false, true, "test usage")
	regErr = RegisterOpt("test3", "", "b", false, true, "test usage")
	regErr = RegisterOpt("test4", "", "c", false, true, "test usage")
	regErr = RegisterOpt("test5", "hey", "", false, true, "test usage")

	if regErr != nil {
		t.Error("Parse() test: Testing opt parsing but got reg error: " + regErr.Error())
	}

	// Now to wipe out os.Args
	os.Args = []string{"ignoreme", "-a=foo", "--wow", "suchval"}

	expr := regexp.MustCompile(regexp.QuoteMeta(ERR_REQ) + ".+")

	// Parse the args, check for errors and correct val
	Parse()
	if HasError() {
		err := GetError()
		if err != nil {
			if !expr.MatchString(err.Error()) {
				t.Error("Parse() test: Wrong type of error given when missing required args: " + err.Error())
			}
		} else {
			t.Error("Parse() test: Severe error, HasError() returned true but error is nil!")
		}
	} else {
		t.Error("Parse() test: Didn't get a parse error when missing required options.")
	}
}
