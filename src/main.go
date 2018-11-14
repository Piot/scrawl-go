/*

MIT License

Copyright (c) 2017 Peter Bjorklund

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

*/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/piot/scrawl-go/src/definition"
	"github.com/piot/scrawl-go/src/scrawl"
)

func parseOptions() (string, bool) {
	var commandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	protocolDefinitionFilename := commandLine.String("protocol", "protocol.txt", "Protocol definition")
	var flagForceColor = commandLine.Bool("color", false, "Enable color output")
	var flagVerbose = commandLine.Bool("verbose", false, "Verbose")

	commandLine.Parse(os.Args[1:])
	if *flagForceColor {
		color.NoColor = false
	}
	return *protocolDefinitionFilename, *flagVerbose
}

func printRoot(root *definition.Root) {
	fmt.Printf("--- Summary ---\n")
	for _, entity := range root.Entities() {
		fmt.Printf("%v\n", entity)
	}
}

func run() error {
	protocolDefinitionFilename, verbose := parseOptions()
	if protocolDefinitionFilename == "" {
		return fmt.Errorf("Must specify a protocol file")
	}
	root, rootErr := scrawl.ParseFile(protocolDefinitionFilename)
	if rootErr != nil {
		return rootErr
	}

	if verbose {
		printRoot(root)
	}
	return nil
}

func main() {
	color.New(color.FgCyan).Fprintf(os.Stderr, "scrawl protocol validator 0.2\n")
	err := run()
	if err != nil {
		color.New(color.FgRed).Fprintf(os.Stderr, "Validation Error: %v\n", err)
	} else {
		color.Green("Validation passed")
	}
}
