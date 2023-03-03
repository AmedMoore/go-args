package args

import (
	"strconv"
)

/*
	Package args implements command-line arguments parsing.

	Usage

	Construct a new parser using args.NewParser()

	This constructs a new Parser struct and returns a pointer to it:
		import "github.com/amedmoore/go-args"
		parser := args.NewParser()
	You can pass the arguments slice to be parsed to this function.
		parser := args.NewParser(os.Args[1:])
	Or you can pass it later to the parser.Parse() function.
	IMPORTANT NOTE is that you have to pass the args slice to one of these
	function, otherwise the parser will act as if you passed an empty slice.

	After you construct the parser struct, call
		parser.Parse()
	to parse the arguments into three categories positional, options,
	and arguments.

	Parsed arguments may then be accessed using one of the following accessors:

	To access positional argument use
	parser.At(index)

	To ask if an option exists use
	parser.HasOption("--my-option")

	To access value of an argument use
	parser.GetString("--my-arg")
	if multiple values was passed with the same name "--my-arg" parser.GetString()
	will return the first value FROM THE RIGHT SIDE,
	meaning if your command was "myapp --my-arg val0 --my-arg val1 --my-arg val2"
	parser.GetString("--my-arg") will return "val2"

	GetString() allow alternative names lookup, example:
	parser.GetString("--help", "-h")

	If your program allow passing multiple values to the same argument use
	parser.Get("--my-arg")
	this will return a slice that includes all values for the argument "--my-arg"

	Command line arguments syntax

	The following forms are permitted:

		arg
		-o
		--option
		-a val
		--arg val

	where "arg" represents a positional argument,
	"-o", and "--option" represents options,
	and "-a val", and "--arg val" represents arguments with value.

	One or two minus signs (hyphens) may be used; they are equivalent.
*/

// FlagPrefix defines the special character used to indicate
// that this arg is a flag, we use this instead of a hard coded
// value to make it easier if we want to support platform-specific
// syntax in the future (i.e. the /arg syntax in Windows).
const FlagPrefix = '-'

// Parser is the main struct that does
// the parsing and hold the results.
type Parser struct {
	rawArgs    []string
	positional []string
	options    []string
	args       []map[string]string
}

// NewParser constructs a new Parser struct
// and returns a pointer to it.
//
// Optionally you can pass the arguments slice
// to be parsed to this function, example:
//
//	parser := args.NewParser(os.Args[1:])
//
// Please note that you have to pass the arguments
// slice to either this or the `Parse()` function.
func NewParser(args ...[]string) *Parser {
	rawArgs := make([]string, 0)
	if len(args) > 0 {
		rawArgs = args[0]
	}
	return &Parser{
		rawArgs:    rawArgs,
		positional: make([]string, 0),
		options:    make([]string, 0),
		args:       make([]map[string]string, 0),
	}
}

// Parse parses the command-line arguments and store,
// the result into the owner struct.
//
// Optionally you can pass the arguments slice
// to be parsed to this function, example:
//
//	parser := args.NewParser()
//	parser.Parse(os.Args[1:])
//
// Please note that arguments passed to this function
// will be ignored if you declared the parser struct
// with initial arguments, like:
//
//	parser := args.NewParser(os.Args[1:])
func (p *Parser) Parse(args ...[]string) error {
	if len(args) > 0 && len(p.rawArgs) == 0 {
		p.rawArgs = args[0]
	}

	for i := 0; i < len(p.rawArgs); i++ {
		arg := p.rawArgs[i]
		ip1 := i + 1
		im1 := i - 1

		if arg == "" {
			continue
		}

		if (arg[0] != FlagPrefix && im1 <= 0) || (arg[0] != FlagPrefix && im1 > 0 && p.rawArgs[im1][0] != FlagPrefix) {
			p.positional = append(p.positional, arg)
			continue
		}

		if arg[0] == FlagPrefix && ip1 < len(p.rawArgs) && p.rawArgs[ip1][0] != FlagPrefix {
			p.args = append(p.args, map[string]string{arg: p.rawArgs[ip1]})
			i++
			continue
		}

		p.options = append(p.options, arg)
	}

	return nil
}

// Positional returns the parsed positional arguments in
// the form of string-slice.
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) Positional() []string {
	return p.positional
}

// At returns the positional argument in the specified index,
// and a bool value indicates weither the positional argument
// exists.
//
// If the specified index was not found (out-of-index) the first
// return value will be an empty string.
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) At(index int) (string, bool) {
	if index < len(p.positional) {
		return p.positional[index], true
	}
	return "", false
}

// Options returns the parsed options in the form of string-slice.
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) Options() []string {
	return p.options
}

// HasOption asks if a specific option was provided, example:
//
//	if parser.HasOption("-h") {
//		// display help message!
//	}
//
// HasOption also supports alias names lookup, example:
//
//	if parser.HasOption("--help", "-h") {
//		// display help message!
//	}
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) HasOption(option string, alts ...string) bool {
	names := append(alts, option)
	for _, opt := range p.options {
		for _, name := range names {
			if opt == name {
				return true
			}
		}
	}
	return false
}

// Args returns the parsed arguments (if has value) in
// the form of a slice of key-value pairs.
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) Args() []map[string]string {
	return p.args
}

// Get return all values provided with the given name.
// example:
//
//	$ myapp --name foo --name bar --name baz
//
//	names := parser.Get("--name")
//	for _, name := range names {
//		print(name)
//	}
//
// Get() allow alternative names lookup, example:
//
//	names := parser.Get("--name", "-n")
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) Get(name string, alts ...string) []string {
	names := append(alts, name)
	args := make([]string, 0)
	for _, n := range names {
		for _, arg := range p.args {
			k, v := p.firstPair(arg)
			if k == n {
				args = append(args, v)
			}
		}
	}
	return args
}

// GetString returns the value of the given argument name.
//
// If multiple values found with the same name, the first
// one from the right will be returned, example:
//
//	$ myapp --name foo --name bar --name baz
//
//	name := parser.GetString("--name")
//	print(name) // baz
//
// GetString() allow alternative name lookup, example:
//
//	names := parser.GetString("--name", "-n")
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) GetString(name string, alts ...string) string {
	names := append(alts, name)
	for i := len(p.args) - 1; i >= 0; i-- {
		for _, n := range names {
			k, v := p.firstPair(p.args[i])
			if k == n {
				return v
			}
		}
	}
	return ""
}

// LookupString returns the value of the given argument name,
// and a bool value indicates weither the argument exists.
//
// If multiple values found with the same name, the first
// one from the right will be returned, example:
//
//	$ myapp --name foo --name bar --name baz
//
//	name, exists := parser.LookupString("--name")
//	print(name)   // baz
//	print(exists) // true
//
// LookupString() allow alternative name lookup, example:
//
//	$ myapp --age 30 -n bar
//
//	name, exists := parser.LookupString("--name", "-n")
//	print(name)   // bar
//	print(exists) // true
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) LookupString(name string, alts ...string) (string, bool) {
	v := p.GetString(name, alts...)
	return v, v != ""
}

// GetInt returns the value of the given argument name.
//
// If multiple values found with the same name, the first
// one from the right will be returned, example:
//
//	$ myapp --age 18 --age 21 --age 30
//
//	age := parser.GetInt("--age")
//	print(age) // 30
//
// GetInt() allow alternative age lookup, example:
//
//	$ myapp --name foo -a 30
//
//	age := parser.LookupString("--age", "-a")
//	print(age) // 30
//
// GetInt() returns -1 If argument was not found.
//
// GetInt() panics if the argument was found but
//
//	was an invalid int value.
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) GetInt(name string, alts ...string) int64 {
	names := append(alts, name)
	for i := len(p.args) - 1; i >= 0; i-- {
		for _, n := range names {
			k, v := p.firstPair(p.args[i])
			if k == n {
				val, err := strconv.Atoi(v)
				if err != nil {
					panic(err)
				}
				return int64(val)
			}
		}
	}
	return -1
}

// LookupInt returns the value of the given argument name,
// and a bool value indicates weither the argument exists.
//
// If multiple values found with the same name, the first
// one from the right will be returned, example:
//
//	$ myapp --age 18 --age 21 --age 30
//
//	age, exists := parser.GetInt("--age")
//	print(age)    // 30
//	print(exists) // true
//
// LookupInt() allow alternative name lookup, example:
//
//	$ myapp --name foo -a 30
//
//	age, exists := parser.LookupInt("--age", "-a")
//	print(age)    // 30
//	print(exists) // true
//
// LookupInt() returns -1 If argument was not found.
//
// LookupInt() panics if the argument was found but
//
//	was an invalid int value.
//
// Make sure to call `Parse()` before using this function.
func (p *Parser) LookupInt(name string, alts ...string) (int64, bool) {
	v := p.GetInt(name, alts...)
	return v, v != -1
}

// firstPair is a helper function that returns the first key-value
// of the given key-value map.
func (p *Parser) firstPair(m map[string]string) (string, string) {
	for k, v := range m {
		return k, v
	}
	return "", ""
}
