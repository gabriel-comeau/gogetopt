# Gogetopt
## A cli option parser for Go

Gogetopt is package which parses CLI options.  It attempts to support all of the types and formats of options
that getopt() does.  This means both --long and -l forms, -lbx combined short options and multiple formats
for values attached to options like:

*  -l=foo
*  --long=foo
*  -l foo
*  --long foo

The order which the arguments appear in the command line invocation don't matter though if you have a value
attached, it needs to follow the appropriate option before another option appears.

To install the package:

```bash
go get github.com/gabriel-comeau/gogetopt
```

Here's a small example program which uses the package.  The program prints each argument it gets to its
own line.  If it gets the **allcaps** option, the output will be in all caps.  If it gets the
**appendstring** option, it will add a space and then the string to each argument.  If it gets the **help**
flag, it will print the generated usage text and then quit.

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gabriel-comeau/gogetopt"
)

// First register the three options and then run Parse()
func init() {
	var regErr error

	// Boolean opts
	regErr = gogetopt.RegisterOpt("allcaps", "caps", "c", true, false, "Switch output to ANGRY")
	if regErr != nil {
		fmt.Println("Gogetopt option registration error: " + regErr.Error())
		os.Exit(1)
	}

	regErr = gogetopt.RegisterOpt("help", "help", "h", true, false, "Prints help message")
	if regErr != nil {
		fmt.Println("Gogetopt option registration error: " + regErr.Error())
		os.Exit(1)
	}

	// String opts
	regErr = gogetopt.RegisterOpt("appendstring", "append", "a", false, false, "String to be appended to each line of ouput")
	if regErr != nil {
		fmt.Println("Gogetopt option registration error: " + regErr.Error())
		os.Exit(1)
	}

	// Parse the options
	gogetopt.Parse()

	if gogetopt.HasError() {
		printHelpAndDie(gogetopt.GetError())
	}

}

// Now, based on which options came in, run the program
func main() {

	if gogetopt.GetBool("help") {
		printHelpAndDie(nil)
	}

	extraArgs := gogetopt.GetArgs()
	for _, str := range extraArgs {
		out := str

		// gogetopt.GetString("key") returns "" if the option wasn't set
		if gogetopt.GetString("appendstring") != "" {
			out += " " + gogetopt.GetString("appendstring")
		}

		if gogetopt.GetBool("allcaps") {
			out = strings.ToUpper(out)
		}

		fmt.Println(out)
	}
}

func printHelpAndDie(err error) {
	if err != nil {
		fmt.Printf("Error during getopt parse: %v\n", err.Error())
	}
	fmt.Println(gogetopt.GetUsage())
	os.Exit(1)
}
```

The package is MIT licensed, details can be found in the included "COPYING" file.
