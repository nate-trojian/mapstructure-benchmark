# UPDATE
The results of these tests became the Go package [goregistry](https://github.com/nate-trojian/goregistry)

# What is this?
I have run into the [mapstructure](https://github.com/mitchellh/mapstructure) package a few times reading online.  It looked useful for solving the problem where I have a common base struct with many different outer structs that embed it and I need to be able to deserialize a message to one of the outer struct types based on a key field.  In fact, this was the example mapstructure gave in its [But Why?!](https://pkg.go.dev/github.com/mitchellh/mapstructure#readme-but-why) section.  As they write, one solution to this is to
> do two passes over the decoding of the JSON (reading the "type" first, and the rest later). However, it is much simpler to just decode this into a `map[string]interface{}` structure, read the `type` key, then use something like this library to decode it into the proper structure.

However, I saw that mapstructure uses reflection, which can be quite slow.  So I came up with these benchmark tests to compare the two approaches.


# Results
Doing two passes using JSON is faster, with less allocations, than using mapstructure.
```
goos: darwin
goarch: arm64
pkg: github.com/nate-trojian/mapstructure-benchmark
BenchmarkMapstructure/Intern-8                  	  293991	      4041 ns/op	    2512 B/op	      60 allocs/op
BenchmarkMapstructure/Salary-8                  	  286898	      4120 ns/op	    2512 B/op	      61 allocs/op
BenchmarkJSONwFunctionRegistry/Intern-8         	  688566	      1774 ns/op	     592 B/op	      15 allocs/op
BenchmarkJSONwFunctionRegistry/Salary-8         	  628483	      1870 ns/op	     584 B/op	      16 allocs/op
BenchmarkJSONwSwitch/Intern-8                   	  667989	      1775 ns/op	     592 B/op	      15 allocs/op
BenchmarkJSONwSwitch/Salary-8                   	  651626	      1862 ns/op	     584 B/op	      16 allocs/op
BenchmarkJSONwReflectRegistry/Intern-8          	  673431	      1741 ns/op	     544 B/op	      14 allocs/op
BenchmarkJSONwReflectRegistry/Salary-8          	  624691	      1884 ns/op	     536 B/op	      15 allocs/op
PASS
coverage: 75.0% of statements
ok  	github.com/nate-trojian/mapstructure-benchmark	10.871s
```

## JSON Registry vs Switch
Out of my own curiosity, I include JSON tests using both a registry and a simple switch statement to determine what the correct type is.  The registry pattern is not one I see a lot in Golang codebases, but is one I find interesting.  The registry pattern does require more scaffolding code in order to setup, but it does allow each individual type to be completely encapsulated within a single file by using the `init` function to register themselves with their given key.  This makes it easier in the future for someone maintaining the codebase to change an existing type, or to add a new one.  While I was expecting a slight performance hit, the additional `39 ns/op` is negligable in my opinion.
