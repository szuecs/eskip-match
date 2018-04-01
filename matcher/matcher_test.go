package matcher

import (
	"fmt"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRoutesScenario struct {
	expectedRouteID string
	expectNoMatch   bool
	reqAttributes   []*RequestAttributes
}

func TestMatcherError(t *testing.T) {
	_, err := New(&Options{
		RoutesFile: "",
	})
	assert.Error(t, err)
}

func TestMacherTestSetHeaders(t *testing.T) {
	routesFile, err := filepath.Abs("./testdata/routes.eskip")
	if err != nil {
		t.Error(err)
		return
	}
	tester, err := New(&Options{
		RoutesFile:  routesFile,
		MockFilters: []string{"customfilter"},
	})

	assert.NoError(t, err)

	res := tester.Test(&RequestAttributes{
		Method: "POST",
		Path:   "/foo",
		Headers: map[string]string{
			"X-Bar": "bar",
		},
	})

	assert.NotNil(t, res.Request())
	assert.Equal(t, "bar", res.Request().Header.Get("X-Bar"))
}

func TestRoutes(t *testing.T) {

	routesFile, err := filepath.Abs("./testdata/routes.eskip")
	if err != nil {
		t.Error(err)
		return
	}

	tester, err := New(&Options{
		RoutesFile:          routesFile,
		MockFilters:         []string{"customfilter"},
		Verbose:             true,
		IgnoreTrailingSlash: true,
	})

	if err != nil {
		t.Error(err)
		return
	}

	scenarios := []testRoutesScenario{
		{
			expectedRouteID: "foo",
			reqAttributes: []*RequestAttributes{
				{
					Method: "POST",
					Path:   "/foo",
				},
				{
					Method: "post",
					Path:   "/foo/1",
				},
			},
		},
		{
			expectedRouteID: "foo_get",
			reqAttributes: []*RequestAttributes{
				{
					Path: "/foo",
				},
				{
					Path: "foo",
				},
			},
		},
		{
			expectedRouteID: "query_param",
			reqAttributes: []*RequestAttributes{
				{
					Path: "/abdc?q=bar",
				},
			},
		},
		{
			expectedRouteID: "bar",
			reqAttributes: []*RequestAttributes{
				{
					Path: "/bar",
				},
			},
		},
		{
			expectedRouteID: "foo_header",
			reqAttributes: []*RequestAttributes{
				{
					Path: "/foo",
					Headers: map[string]string{
						"Accept": "application/json",
					},
				},
			},
		},
		{
			expectedRouteID: "customfilter",
			reqAttributes: []*RequestAttributes{
				{
					Path: "/customfilter",
				},
			},
		},
		{
			expectedRouteID: "no-match",
			expectNoMatch:   true,
			reqAttributes: []*RequestAttributes{
				{
					Path: "/blobblob",
				},
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.expectedRouteID, func(t *testing.T) {
			for _, a := range s.reqAttributes {
				result := tester.Test(a)

				route := result.Route()
				req := result.Request()
				attrs := result.Attributes()

				assert.NotNil(t, req)
				assert.NotNil(t, attrs)

				if route != nil {
					assert.Contains(t, result.PrettyPrint(), "matching")
				} else {
					assert.NotContains(t, result.PrettyPrint(), "matching")
				}

				if s.expectNoMatch == true && route != nil {
					t.Errorf("request: %s %s shouldn't match but matches route id: %s", req.Method, a.Path, route.Id)
					return
				}

				if s.expectNoMatch == true && route == nil {
					return
				}

				if route == nil {
					t.Errorf("expected route id to be '%s' but no match\n request: %s %s", s.expectedRouteID, req.Method, a.Path)
				} else if route.Id != s.expectedRouteID {
					t.Errorf("expected route id to be '%s' but got '%s'\n request: %s %s", s.expectedRouteID, route.Id, req.Method, a.Path)
				}

			}
		})
	}
}

func Example() {
	routesFile, err := filepath.Abs("./testdata/routes.eskip")
	if err != nil {
		log.Fatal(err)
		return
	}

	m, err := New(&Options{
		RoutesFile: routesFile,
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	result := m.Test(&RequestAttributes{
		Method: "GET",
		Path:   "/bar",
		Headers: map[string]string{
			"Accept": "application/json",
		},
	})

	route := result.Route()

	if route != nil {
		fmt.Println(route.Id)
		// Output: bar
	}
}
