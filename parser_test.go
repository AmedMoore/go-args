package args

import "testing"

var testArgs = []string{
	// positional arguments
	"pos0", "pos1", "pos2", "pos3", "pos4", "pos5",

	// option aliases
	"-a", "-b", "-c", "-d", "-e", "-f",

	// options
	"--opt0", "--opt1", "--opt2", "--opt3", "--opt4", "--opt5",

	// argument aliases
	"-g", "val0", "-h", "val1", "-i", "2",

	// arguments
	"--arg0", "val0", "--arg1", "val1", "--arg2", "val2",
	"--arg0", "val3", "--arg1", "val4", "--arg2", "5",
}

func Test_Parser_Parse(t *testing.T) {
	parser := NewParser(testArgs)
	err := parser.Parse()
	if err != nil {
		t.Errorf("Parse(): %s", err.Error())
	}

	// test length of positional arguments
	pos := len(parser.Positional())
	if pos != 6 {
		t.Errorf("len(parser.Positional()) = %d; want 6", pos)
	}

	// test length of options
	opt := len(parser.Options())
	if opt != 12 {
		t.Errorf("len(parser.Options()) = %d; want 12", opt)
	}

	// test length of value arguments
	args := len(parser.Args())
	if args != 9 {
		t.Errorf("len(parser.Args()) = %d; want 9", args)
	}

	// test positional arguments at 1
	at1, exist := parser.At(1)
	if !exist || at1 != "pos1" {
		t.Errorf("parser.At(1) = \"%s\"; want \"pos1\"", at1)
	}

	// test positional arguments at 4
	at4, exist := parser.At(4)
	if !exist || at4 != "pos4" {
		t.Errorf("parser.At(4) = \"%s\"; want \"pos4\"", at4)
	}

	// test options '-c' and '-e' exists
	if !parser.HasOption("-c") {
		t.Error("parser.HasOption(\"-c\") = false; want true")
	}
	if !parser.HasOption("-e") {
		t.Error("parser.HasOption(\"-e\") = false; want true")
	}

	// test options '--opt0' and '--opt5' exists
	if !parser.HasOption("--opt0") {
		t.Error("parser.HasOption(\"--opt0\") = false; want true")
	}
	if !parser.HasOption("--opt5") {
		t.Error("parser.HasOption(\"--opt5\") = false; want true")
	}

	// test value of argument '-h'
	h, exist := parser.LookupString("-h")
	if !exist || h != "val1" {
		t.Errorf("parser.LookupString(\"-h\") = \"%s\"; want \"val1\"", h)
	}

	// test value of argument '--arg1'
	arg1, exist := parser.LookupString("--arg1")
	if !exist || arg1 != "val4" {
		t.Errorf("parser.LookupString(\"--arg1\") = \"%s\"; want \"val4\"", arg1)
	}

	// test value of argument '-i'
	i, exist := parser.LookupInt("-i")
	if !exist || i != 2 {
		t.Errorf("parser.LookupInt(\"-i\") = \"%d\"; want \"2\"", i)
	}

	// test value of argument '--arg2'
	arg2, exist := parser.LookupInt("--arg2")
	if !exist || arg2 != 5 {
		t.Errorf("parser.LookupInt(\"--arg2\") = \"%d\"; want \"5\"", arg2)
	}

	// test multiple values of argument '--arg1'
	args1 := parser.Get("--arg1")
	if len(args1) != 2 || args1[0] != "val1" || args1[1] != "val4" {
		t.Errorf("parser.Get(\"--args1\") = \"%s\"; want [\"val1\", \"val4\"]", args1)
	}

	// test alternative names lookup
	alt, exists := parser.LookupString("--missing", "-g")
	if !exists || alt != "val0" {
		t.Errorf("parser.LookupString(\"--missing\", \"-g\") = \"%s\"; want \"val0\"", alt)
	}
}

func Benchmark_Parser_Parse(b *testing.B) {
	parser := NewParser(testArgs)
	err := parser.Parse()
	if err != nil {
		b.Errorf("Parse(): %s", err.Error())
	}
}
