package matchertest

import (
	"fmt"
	ioutil "io/ioutil"
	"path/filepath"
	"strings"

	"github.com/rbarilani/eskip-match/matcher"
	"gopkg.in/yaml.v2"
)

type Scenario struct {
	Coverage    int            `yaml:"coverage"`
	File        string         `yaml:"file"`
	MockFilters []string       `yaml:"mock_filters"`
	Tests       []ScenarioTest `yaml:"tests"`
}

type ScenarioTest struct {
	RouteID    string                      `yaml:"route_id"`
	NoMatch    bool                        `yaml:"no_match"`
	Attributes []matcher.RequestAttributes `yaml:"attributes"`
}

func (s *Scenario) FromYAMLFile(file string) (*Scenario, error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		return nil, err
	}

	if s.File == "" {
		return nil, fmt.Errorf("You should define a file property")
	}
	return s, nil
}

type T interface {
	Fatal(...interface{})
	Errorf(format string, args ...interface{})
	Log(...interface{})
	Logf(format string, args ...interface{})
}

func Run(t T, file string) {

	path, _ := filepath.Abs(file)
	path = filepath.Dir(path)
	scenario := &Scenario{}
	scenario, err := scenario.FromYAMLFile(file)
	if err != nil {
		t.Fatal(err)
		return
	}

	hits := map[string]int{}
	tester, err := matcher.New(&matcher.Options{
		RoutesFile:  filepath.Join(path, scenario.File),
		MockFilters: scenario.MockFilters,
	})

	if err != nil {
		t.Fatal(err)
		return
	}

	for _, test := range scenario.Tests {
		for _, attrs := range test.Attributes {
			reqAttrs := &matcher.RequestAttributes{
				Path:    attrs.Path,
				Method:  attrs.Method,
				Headers: attrs.Headers,
			}
			res := tester.Test(reqAttrs)
			route := res.Route()
			req := res.Request()

			if route == nil {
				if test.NoMatch == false {
					t.Errorf("expected route id to be '%s' but no match\n request: %s %s", test.RouteID, req.Method, attrs.Path)
				}
			} else {
				hits[route.Id] = hits[route.Id] + 1
				if test.NoMatch == true {
					t.Errorf("request: %s %s shouldn't match but matches route id: %s", req.Method, attrs.Path, route.Id)
				} else if route.Id != test.RouteID {
					t.Errorf("expected route id to be '%s' but got '%s'\n request: %s %s", test.RouteID, route.Id, req.Method, attrs.Path)
				}
			}
		}
	}

	routes := tester.Routes()
	nohits := []string{}

	for _, route := range routes {
		if _, ok := hits[route.Id]; !ok {
			nohits = append(nohits, route.Id)
		}
	}

	hitslen := float64(len(hits))
	routeslen := float64(len(routes))
	coverage := ((hitslen / routeslen) * 100)
	coverageround := int(coverage + .5)
	t.Logf("Coverage: %d%%", coverageround)
	if coverageround < scenario.Coverage {
		t.Errorf("Expected coverage to be %d%% but got %d%%", scenario.Coverage, coverageround)
		t.Errorf("Missing hits for route ids: %s", strings.Join(nohits, ", "))
	}
}
