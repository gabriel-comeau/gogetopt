// Gogetopt attempts to be a more in-depth cli option parser for go.  While it doesn't work exactly like
// GNU getopt(), it should be able to parse the same kinds of CLI args.
//
// The major difference is that there is no DSL for defining arguments the way getopt() does.  The colon
// syntax "f::" to define if an option takes a value and whether or not this is optional is gone.
//
// Instead options are created via RegisterOpt(), where you explicitly state whether or not an option
// is required and whether or not it takes a value or is just a boolean switch.  You can define both
// a long (--switch) and short (-s) for the same option.
//
// Once registered, call Parse() which will read the current command line args and compare them to
// which opts were registered.  If they are found, the calling program can get the values with GetBool()
// and GetString().
//
// If there were any errors which occurred during option registration, they will be returned at the time
// of reg.
//
// Parsing errors work a bit differently.  Parsing will halt on the first error encountered.  You can check
// for a parse error with HasError() and if it returns true use GetError() to get the error object.  Alternatively
// just call GetError() and check if it isn't nil.  Whichever you like better.
//
// Accepted forms of options:
//
// -l               (boolean shortopt)
// --long           (boolean longopt)
// -l=val           (string shortopt)
// --long=val       (string longopt)
// -l val           (string shortopt)
// --long val       (string longopt)
// -fbx             (combined boolean shortopts)
// -fVAL            (string shortopt, no space or = sign)
//
// Any non-boolean option can be set to required, which will result in a parse error state if the option isn't
// found
//
// Written by Gabriel Comeau
//
// See COPYING for license
package gogetopt

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

type opt struct {
	key      string
	long     string
	short    string
	isBool   bool
	required bool
	usage    string
}

var (

	// Master table and lookup tables
	opts         map[string]*opt
	shortKeys    map[string]*opt
	longKeys     map[string]*opt
	requiredOpts map[string]bool

	// Value holders
	boolVals   map[string]bool
	stringVals map[string]string
	extraArgs  []string

	// Regexes
	singleDash       *regexp.Regexp
	multiDash        *regexp.Regexp
	singleDashEquals *regexp.Regexp
	multiDashEquals  *regexp.Regexp

	// Error holder
	parseError string
)

const (
	ERR_MISSING_VAL            string = "Missing value for option: "
	ERR_NO_OPT                 string = "No such option: "
	ERR_BOOL_REQ               string = "An option can't be both boolean and required: "
	ERR_REQ                    string = "Required option(s) not provided: "
	ERR_BOOL_WITH_VAL          string = "Boolean options can't be passed values: "
	ERR_NONBOOL_MULTI          string = "Combined opts can't be non-boolean: "
	ERR_NO_KEY                 string = "An option must contain either a long or short key (or both): "
	ERR_SHORT_TOO_LONG         string = "A short option can be no longer one character: "
	ERR_LONG_TOO_SHORT         string = "A long option must be longer than one character: "
	ERR_OPT_KEY_ALREADY_EXISTS string = "An option was already registered with key: "
	ERR_SHORT_ALREADY_EXISTS   string = "An option was already registered with short key: "
	ERR_LONG_ALREADY_EXISTS    string = "An option was already registered with long key: "
)

func init() {
	opts = make(map[string]*opt)
	shortKeys = make(map[string]*opt)
	longKeys = make(map[string]*opt)
	requiredOpts = make(map[string]bool)

	boolVals = make(map[string]bool)
	stringVals = make(map[string]string)
	extraArgs = make([]string, 0)

	singleDash = regexp.MustCompile("^-.+")
	multiDash = regexp.MustCompile("^--.+")
	singleDashEquals = regexp.MustCompile("^-.+=")
	multiDashEquals = regexp.MustCompile("^--.+=")

	parseError = ""
}

//
// PUBLIC API
//

// Register an option to the list of options which will be parsed.  An option can have both
// a long and short val, and it will respond to either form on the command line.  If you only
// want one form to work, just push in an empty string.
func RegisterOpt(key, long, short string, isBool, isReq bool, usage string) error {
	o := new(opt)
	o.key = key
	o.long = stripDashes(long)
	o.short = stripDashes(short)
	o.isBool = isBool
	o.required = isReq
	o.usage = usage

	// Error condition - can't make a switch be both required and boolean
	if o.isBool && o.required {
		return errors.New(ERR_BOOL_REQ + o.key)
	}

	// Error condition - need to have at least either a short or long key for the opt
	if o.short == "" && o.long == "" {
		return errors.New(ERR_NO_KEY + o.key)
	}

	// Make sure lengths for short/longs are sane

	if o.short != "" && len(o.short) > 1 {
		return errors.New(ERR_SHORT_TOO_LONG + o.short)
	}

	if o.long != "" && len(o.long) < 2 {
		return errors.New(ERR_LONG_TOO_SHORT + o.long)
	}

	// Check for already existing keys registered (main key, short and long)

	_, oPres := opts[o.key]
	if oPres {
		return errors.New(ERR_OPT_KEY_ALREADY_EXISTS + o.key)
	}

	if o.short != "" {
		_, sPres := shortKeys[o.short]
		if sPres {
			return errors.New(ERR_SHORT_ALREADY_EXISTS + o.short)
		}
	}

	if o.long != "" {
		_, lPres := longKeys[o.long]
		if lPres {
			return errors.New(ERR_LONG_ALREADY_EXISTS + o.long)
		}
	}

	// Assign the option to the various maps as applicable
	opts[o.key] = o

	if o.short != "" {
		shortKeys[o.short] = o
	}

	if o.long != "" {
		longKeys[o.long] = o
	}

	if o.required {
		requiredOpts[o.key] = true
	}

	return nil
}

// Remove any registered options.  This is primarly to ease testing but could potentially be
// handy depending on execution context of a program?  This will also clear the list of "extra"
// arguments - to use any args at all from getopt, you'll need to re-run parse after running this.
func ClearAll() {
	for key, _ := range opts {
		Clear(key)
	}

	// Since we're clearing everything, wipe out the extra args too
	extraArgs = make([]string, 0)
	// Also any existing parse errors
	parseError = ""

}

// Remove a single option by key.  This will also remove it's bool/string val if parse has
// already been run
func Clear(key string) {
	opt, ok := opts[key]
	if ok {

		if opt.short != "" {
			delete(shortKeys, opt.short)
		}

		if opt.long != "" {
			delete(longKeys, opt.long)
		}

		if opt.required {
			delete(requiredOpts, opt.key)
		}

		// Delete any registered values for this option
		if opt.isBool {
			_, ok := boolVals[opt.key]
			if ok {
				delete(boolVals, opt.key)
			}
		} else {
			_, ok := stringVals[opt.key]
			if ok {
				delete(stringVals, opt.key)
			}
		}

		delete(opts, key)

	}
}

// Get a string value for an option key.  Only makes sense if Parse() has been called.
func GetString(key string) string {
	val, ok := stringVals[key]
	if ok {
		return val
	}
	return ""
}

// Get a bool value for an option key.  Only makes sense if Parse() has been called.
func GetBool(key string) bool {
	_, ok := boolVals[key]
	if ok {
		return true
	}
	return false
}

// Get any "extra" non-option arguments passed to the program.  This excludes argv[1] - the program
// name.  Only makes sense if Parse() has been called.
func GetArgs() []string {
	return extraArgs
}

// Check if there's a parse error.  Only makes sense if Parse() has been called.
func HasError() bool {
	if parseError == "" {
		return false
	}
	return true
}

// Get the parse error if present.  Only makes sense if Parse() has been called.
func GetError() error {
	if HasError() {
		return errors.New(parseError)
	}
	return nil
}

// Get the usage for each option
func GetUsage() string {
	useStr := ""
	for _, opt := range opts {
		// <option> <arg> <usage>

		if opt.short != "" {
			useStr += "-" + opt.short + " "
		}

		if opt.long != "" {
			useStr += "--" + opt.long + " "
		}

		if opt.required {
			useStr += "REQUIRED "
		}

		if !opt.isBool {
			useStr += "<value> "
		}

		if opt.usage != "" {
			useStr += opt.usage + "\n"
		}
	}
	return useStr
}

func Parse() {
	args := os.Args

	foundReqs := make(map[string]bool)

	// This isn't an error, it just doesn't need to parse any arguments.  Unless of course there are
	// required arguments.  Then it's totally an error.
	if len(args) < 2 {
		if len(requiredOpts) > 0 {
			parseError = getMissingReqOptsError(foundReqs, requiredOpts)
			return
		}

		return
	}

	// Main loop, iterating through each argument passed in to program.  Start at index 1 instead
	// of 0 because there's no sense in processing argv[1] (the program's name)
	for i := 1; i < len(args); i++ {

		arg := args[i]

		if singleDashEquals.MatchString(arg) || multiDashEquals.MatchString(arg) {

			// This is the case for -f=bar or --foo=bar

			key, val, err := getValForEqualsSignArg(arg)
			if err != nil {
				parseError = err.Error()
				return
			}

			// This should realistically never error out since getValForEqualsSignArg() should
			// have covered that possibility already.  Do the if checks anyway to prevent a run
			// time crash.  This may turn out to be a poor decision.
			opt, ok := opts[key]
			if ok {
				if opt.required {
					foundReqs[key] = true
				}
			}

			// All good
			stringVals[key] = val

		} else if multiDash.MatchString(arg) {
			// This is a --longopt formed option.  It can either be a boolean option or it can
			// have an argument trailing after it (as the next arg in argv[]).  It has already
			// been checked for the --foo=bar form and isn't that.
			//
			// TODO: is it a valid gnu-ism to do --fooVAL like with shortopts?

			stripped := stripDashes(arg)
			opt, ok := longKeys[stripped]

			if !ok {
				parseError = ERR_NO_OPT + arg
				return
			}

			if opt.isBool {
				// If it's a boolean value, set it and stop here
				boolVals[opt.key] = true

			} else {

				// Basically, assuming that --fooVAL is NOT a valid gnuism, do a lookahead which
				// attempts to get the next value in the args list and use that as a value.  If
				// that's a valid value (not another opt) set that value, otherwise it's an error.

				val := lookaheadForVal(args, i)
				if val == "" {
					parseError = ERR_MISSING_VAL + arg
					return
				}

				if opt.required {
					foundReqs[opt.key] = true
				}

				// All good - since a lookahead was done the loop counter MUST be incremented here
				// so an argument doesn't get double-processed
				i++
				stringVals[opt.key] = val
			}

		} else if singleDash.MatchString(arg) {

			// Saved the most complex possibility for last.  Shortopts can have many possible
			// outcomes (this ignores the -l=val form, already handled above):
			//
			// "-l" boolean switch
			// "-l val" value after the switch (needs lookahead)
			// "-lmx" multiopt (must all be boolean to be valid)
			// "-lVAL" value immediately after the switch, no lookahead

			// First strip the dashes
			stripped := stripDashes(arg)

			// Check length - if the len is 1, it's got to be a boolean switch or needs
			// lookahead to find the value
			if len(stripped) == 1 {
				opt, ok := shortKeys[stripped]
				if ok {

					if opt.isBool {
						boolVals[opt.key] = true
					} else {
						val := lookaheadForVal(args, i)
						if val == "" {
							parseError = ERR_MISSING_VAL + arg
							return
						}

						if opt.required {
							foundReqs[opt.key] = true
						}

						i++
						stringVals[opt.key] = val
					}

				} else {
					parseError = ERR_NO_OPT + arg
					return
				}
			} else {
				// Longer than 1 - this means either a multiopt or -lVAL format.
				// First check for multiopt
				multiOpts := getMultiOptKeys(arg)
				if multiOpts != nil {

					for _, k := range multiOpts {

						// Already did check for map presence in getMultiOptKeys()
						opt := shortKeys[k]
						if !opt.isBool {
							parseError = ERR_NONBOOL_MULTI + k
							return
						}
					}

					// Second iteration of the same slice.  If the code is here it means everything
					// was correct.
					for _, k := range multiOpts {
						// All good, so set these
						boolVals[shortKeys[k].key] = true
					}

				} else {
					// Not multiopt, so the last possibility is that the characters immediately following
					// the option are the string value.

					// Get the first char, make sure it's an actual option
					key := string(stripped[0])
					opt, ok := shortKeys[key]
					if !ok {
						parseError = ERR_NO_OPT + key
						return
					}

					// Make sure this isn't a boolean
					if opt.isBool {
						parseError = ERR_BOOL_WITH_VAL + arg
						return
					}

					// OK, all of the stuff that isn't the key in the string is the value
					// This counts as all good

					if opt.required {
						foundReqs[opt.key] = true
					}

					val := string(stripped[1:])
					stringVals[opt.key] = val
				}
			}

		} else {

			// Finally, this is just a "default" argument, no part of any option.  It goes into
			// its own slice of values, in the order provided to the script.
			extraArgs = append(extraArgs, arg)
		}
	}

	if len(requiredOpts) > 0 {
		parseError = getMissingReqOptsError(foundReqs, requiredOpts)
		if parseError != "" {
			return
		}
	}
}

//
//  Helper functions to make the parser more readable
//

// When provided with an argument with an equals sign in it, this will
// split the parts up and do checking on the option to make sure it
// both exists and isn't boolean
func getValForEqualsSignArg(arg string) (key, val string, err error) {

	// Defaults for the return values
	key = ""
	val = ""
	err = nil

	// Check to see if we can split the parts up properly
	parts := splitEqualsArg(arg)
	if parts == nil {
		err = errors.New(ERR_MISSING_VAL + arg)
		return
	}

	var opt *opt
	var ok bool = false

	if multiDash.MatchString(arg) {
		opt, ok = longKeys[parts[0]]
		if !ok {
			err = errors.New(ERR_NO_OPT + parts[0])
			return
		}
	} else if singleDash.MatchString(arg) {
		opt, ok = shortKeys[parts[0]]
		if !ok {
			err = errors.New(ERR_NO_OPT + parts[0])
			return
		}
	}

	// Make sure this isn't a boolean option
	if opt.isBool {
		err = errors.New(ERR_BOOL_WITH_VAL + arg)
		return
	}

	if len(parts[1]) < 1 {
		err = errors.New(ERR_MISSING_VAL + arg)
		return
	}

	// All good
	key = opt.key
	val = parts[1]
	return
}

// When passed in something in the form of -xyx, it could have one of two
// meanings:  -x -y -z or -x=yz.  This function checks to see if it's the latter
// and if so returns each opt shortval as a string slice
func getMultiOptKeys(arg string) []string {

	// strip the "-" from the front of the arg (in case)
	workingArg := stripDashes(arg)
	parts := strings.Split(workingArg, "")
	multiOptParts := make([]string, 0)
	isMultiOpt := true
	for _, part := range parts {
		_, ok := shortKeys[part]
		if !ok {
			isMultiOpt = false
			break
		} else {
			multiOptParts = append(multiOptParts, part)
		}
	}

	if isMultiOpt {
		return multiOptParts
	}

	return nil
}

// When passed a -f=bar or --foo=bar type argument where the value
// comes after the equals sign, this function will take them apart and
// return them as a 2 part slice of strings.  [0] is the key and [1] is the
// value.
func splitEqualsArg(arg string) []string {

	workingArg := stripDashes(arg) // copy to use original for errors

	parts := strings.Split(workingArg, "=")
	if len(parts) == 2 {
		if len(parts[1]) > 0 {
			return parts
		} else {
			return nil
		}
	}
	return nil
}

// Remove the - or -- from an option
func stripDashes(arg string) string {
	if multiDash.MatchString(arg) {
		return arg[2:]
	} else if singleDash.MatchString(arg) {
		return arg[1:]
	}
	return arg
}

// This looks for the final type of string value to an option - the kind which is actually
// the next argument in the args list (it is not directly attached or with an equals sign)
func lookaheadForVal(args []string, currentKey int) string {
	if len(args)-1 > currentKey {
		nextVal := args[currentKey+1]

		// The args shouldn't match --, -, -x=y or --x=y!
		if singleDashEquals.MatchString(nextVal) || singleDash.MatchString(nextVal) ||
			multiDashEquals.MatchString(nextVal) || multiDash.MatchString(nextVal) {

			return ""
		}
		return nextVal
	}
	return ""
}

// Check to see if any required options are missing and generate/return an error message if so.
func getMissingReqOptsError(foundReqOpts map[string]bool, requiredOpts map[string]bool) string {

	missingKeys := make([]string, 0)

	for reqName, _ := range requiredOpts {
		_, found := foundReqOpts[reqName]
		if !found {
			missingKeys = append(missingKeys, reqName)
		}
	}

	if len(missingKeys) > 0 {
		errorText := ERR_REQ
		for i, mk := range missingKeys {
			opt, ok := opts[mk]
			msgKey := mk
			if ok {

				if opt.short != "" {
					msgKey = "-" + opt.short
				}

				if opt.long != "" {
					msgKey = "--" + opt.long
				}

				if opt.short != "" && opt.long != "" {
					msgKey = "-" + opt.short + " or " + "--" + opt.long
				}
			}

			if i < len(missingKeys)-1 {
				errorText += msgKey + ", "
			} else {
				errorText += msgKey
			}

		}
		return errorText
	}

	return ""
}
