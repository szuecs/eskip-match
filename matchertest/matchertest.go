package matchertest

import (
	ioutil "io/ioutil"
	"path/filepath"
	"testing"

	"github.com/rbarilani/eskip-match/matcher"
	"gopkg.in/yaml.v2"
)

type Scenarios struct {
	Scenarios []Scenario `yaml:"scenarios"`
}

type Scenario struct {
	File        string         `yaml:"file"`
	MockFilters []string       `yaml:"mock_filters"`
	Tests       []ScenarioTest `yaml:"tests"`
}

type ScenarioTest struct {
	RouteID    string                      `yaml:"route_id"`
	NoMatch    bool                        `yaml:"no_match"`
	Attributes []matcher.RequestAttributes `yaml:"attributes"`
}

func (s *Scenarios) FromYAMLFile(file string) (*Scenarios, error) {
	yamlFile, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func Run(t *testing.T, file string) {
	path, _ := filepath.Abs(file)
	path = filepath.Dir(path)
	scenarios := &Scenarios{}
	t.Log(path)
	scenarios, err := scenarios.FromYAMLFile(file)
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, scenario := range scenarios.Scenarios {
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
					if test.NoMatch == true {
						t.Errorf("request: %s %s shouldn't match but matches route id: %s", req.Method, attrs.Path, route.Id)
					} else if route.Id != test.RouteID {
						t.Errorf("expected route id to be '%s' but got '%s'\n request: %s %s", test.RouteID, route.Id, req.Method, attrs.Path)
					}
				}
			}
		}
	}
}
