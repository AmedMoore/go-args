Package `args` implements command-line arguments parsing.

## Usage

Construct a new parser using `args.NewParser()`.

This constructs a new `Parser` struct and returns a pointer to it:

```go
import "github.com/amedmoore/go-args"
parser := args.NewParser()
```

You can pass the arguments slice to be parsed to this function.

```go
parser := args.NewParser(os.Args[1:])
```

Or you can pass it later to the `parser.Parse()` function.

> **IMPORTANT NOTE** is that you have to pass the args slice to one of these
> function, otherwise the parser will act as if you passed an empty slice.

After you construct the parser struct, call

```go
parser.Parse()
```

to parse the arguments into three categories _positional_, _options_,
and _arguments_.

Parsed arguments may then be accessed using one of the following accessors:

To access positional argument use

```go
parser.At(index)
```

To ask if an option exists use

```go
parser.HasOption("--my-option")
```

To access value of an argument use

```go
parser.GetString("--my-arg")
```

if multiple values was passed with the same name, `parser.GetString()`
will return the first value **FROM THE RIGHT SIDE**,
meaning if your command was

```console
$ myapp --my-arg val0 --my-arg val1 --my-arg val2
```

`parser.GetString("--my-arg")` will return `"val2"`

If your program allow passing multiple values to the same argument use

```go
parser.Get("--my-arg")
```

this will return a slice that includes all values for the argument `"--my-arg"`

Command line arguments syntax

The following forms are permitted:

    arg
    -o
    --option
    -a val
    --arg val

where `arg` represents a positional argument,
`-o`, and `--option` represents options,
and `-a val`, and `--arg val` represents arguments with value.

One or two minus signs (hyphens) may be used; they are equivalent.

## Testing

- On Windows:

  ```console
  $ go test -bench .
  goos: windows
  goarch: amd64
  pkg: github.com/amedmoore/go-args
  cpu: Intel(R) Core(TM) i7-10870H CPU @ 2.20GHz
  Benchmark_Parser_Parse-16           1000000000
  PASS
  ok      github.com/amedmoore/go-args        0.025s
  ```

- On Mac:

  ```console
  $ go test -bench .
  goos: darwin
  goarch: arm64
  pkg: github.com/amedmoore/go-args
  Benchmark_Parser_Parse-8      1000000000      0.0000110 ns/op
  PASS
  ok      github.com/amedmoore/go-args        0.729s
  ```

## TODO

- [ ] Better error handling.
- [ ] Support for the (`--arg=val`) syntax.
- [ ] Support for the Windows (`/opt`, `/arg val`, and `/arg=val`) syntax.
- [ ] Maybe add auto-cast for argument values? like `GetString() string`, `GetInt() int`, and `GetBool() bool`, etc...
- [x] Allow alternative names lookup in argument accessor functions (i.e. `GetString("-h", "--help")`)
- [x] ~~Add type alias `Parser` for the struct `ArgsParser`.~~

## License

This package is licensed under the [MIT License][license] feel free to use it as you want!

[license]: https://github.com/amedmoore/go-args/blob/main/LICENSE
