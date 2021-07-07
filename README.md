# What is this?
I have run into the [mapstructure](https://github.com/mitchellh/mapstructure) package a few times reading online.  It looked useful for solving the problem where I have a common base struct with many different outer structs that embed it and I need to be able to deserialize a message to one of the outer struct types based on a key field.  In fact, this was the example mapstructure gave in its [But Why?!](https://pkg.go.dev/github.com/mitchellh/mapstructure#readme-but-why) section.  As they write, one solution to this is to
> [D]o two passes over the decoding of the JSON (reading the "type" first, and the rest later). However, it is much simpler to just decode this into a `map[string]interface{}` structure, read the `type` key, then use something like this library to decode it into the proper structure.

However, I saw that mapstructure uses reflection, which can be quite slow.  So I came up with these benchmark tests to compare the two approaches.


# Results
Doing two passes using JSON is faster, with less allocations, than using mapstructure.
```
goos: darwin
goarch: arm64
pkg: github.com/nate-trojian/mapstructure-benchmark
BenchmarkMapstructure/Intern-8         	  296448	      4071 ns/op	    2512 B/op	      60 allocs/op
BenchmarkMapstructure/Salary-8         	  288249	      4153 ns/op	    2512 B/op	      61 allocs/op
BenchmarkJSONwRegistry/Intern-8        	  666883	      1793 ns/op	     592 B/op	      15 allocs/op
BenchmarkJSONwRegistry/Salary-8        	  623422	      1883 ns/op	     584 B/op	      16 allocs/op
BenchmarkJSONwSwitch/Intern-8          	  667669	      1773 ns/op	     592 B/op	      15 allocs/op
BenchmarkJSONwSwitch/Salary-8          	  629593	      1898 ns/op	     584 B/op	      16 allocs/op
PASS
ok  	github.com/nate-trojian/mapstructure-benchmark	8.524s
```

## JSON Registry vs Switch
Out of my own curiosity, I include JSON tests using both a registry and a simple switch statement to determine what the correct type is.  The registry pattern is not one I see a lot in Golang codebases, but is one I find interesting.  The registry pattern does require more scaffolding code in order to setup, but it does allow each individual type to be completely encapsulated within a single file by using the `init` function to register themselves with their given key.  This makes it easier in the future for someone maintaining the codebase to change an existing type, or to add a new one.  While I was expecting a slight performance hit, the additional `39 ns/op` is negligable in my opinion.
