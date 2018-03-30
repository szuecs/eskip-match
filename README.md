# eskip-match (work in progress)

A package that helps you test [skipper](https://github.com/zalando/skipper) [`.eskip`](https://zalando.github.io/skipper/dataclients/eskip-file/) files routing matching logic.

## Install

```
go get  github.com/rbarilani/eskip-match/...
```

## Usage

Given an `.eskip` file:

*routes.eskip*
```
foo: Path("/foo") -> http://foo.com
bar: Path("/bar") -> http://bar.com
```

You can write a `go test` able to check if the matching logic is what you expect, using `eskip-match/matcher` package, for example something like so:

*main_test.go*
```go
package main

import (
	"github.com/rbarilani/eskip-match/matcher"
	"path/filepath"
	"testing"
)

type matcherTestScenario struct {
	expectedRouteID string
	reqAttributes   []*matcher.RequestAttributes
}

func TestRoutes(t *testing.T) {

	routesFile, err := filepath.Abs("/.routes.eskip")
	if err != nil {
		t.Error(err)
		return
	}
	tester, err := matcher.New(&matcher.Options{
		RoutesFile:    routesFile
	})

	if err != nil {
		t.Error(err)
		return
	}

	scenarios := []matcherTestScenario{
		{
			expectedRouteID: "foo",
			reqAttributes: []*matcher.RequestAttributes{
				{
					Path: "/foo",
				},
        {
					Path: "/foo/1",
				},
			},
		},
		{
			expectedRouteID: "bar",
			reqAttributes: []*matcher.RequestAttributes{
				{
					Path: "/bar",
				}
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.expectedRouteID, func(t *testing.T) {
			for _, a := range s.reqAttributes {
				route := tester.Test(a)
				if route == nil {
					t.Errorf("expected route id to be '%s' but no match\n request: %s", s.expectedRouteID, a.Path)
				} else if route.Id != s.expectedRouteID {
					t.Errorf("expected route id to be '%s' but got '%s'\n request: %s", s.expectedRouteID, route.Id, a.Path)
				}
			}
		})
	}
}

```

## CLI

The package provide a binary cli tool: `eskip-match`

```
NAME:
   eskip-match - A command line tool that helps you test .eskip files routing matching logic

USAGE:
   eskip-match [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     test, t  given a routes file and request attributes, checks if a route match
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load configuration from FILE
   --help, -h              show help
   --version, -v           print the version
```
### Commands

### Test

With `eskip-match test` command you can check if a route matches given specific request attributes.

#### Examples  

Test if path `/foo` matches a route

```bash
eskip-match test routes.eskip -p /foo
```

## License

Copyright 2018 Ruben Barilani

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.