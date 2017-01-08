package vpn

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/mholt/caddy"
)

func TestHeadersParse(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
		expected  []Rule
	}{
		{`header /foo Foo "Bar Baz"`,
			false, []Rule{
				{Path: "/foo", Headers: http.Header{
					"Foo": []string{"Bar Baz"},
				}},
			}},
		{`header /bar {
			Foo "Bar Baz"
			Baz Qux
			Foobar
		}`,
			false, []Rule{
				{Path: "/bar", Headers: http.Header{
					"Foo":    []string{"Bar Baz"},
					"Baz":    []string{"Qux"},
					"Foobar": []string{""},
				}},
			}},
		{`header /foo {
				Foo Bar Baz
			}`, true,
			[]Rule{}},
		{`header /foo {
				Test "max-age=1814400";
			}`, true, []Rule{}},
	}

	for i, test := range tests {
		actual, err := parse(caddy.NewTestController("http", test.input))

		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}

		if len(actual) != len(test.expected) {
			t.Fatalf("Test %d expected %d rules, but got %d",
				i, len(test.expected), len(actual))
		}

		for j, expectedRule := range test.expected {
			actualRule := actual[j]

			if actualRule.Path != expectedRule.Path {
				t.Errorf("Test %d, rule %d: Expected path %s, but got %s",
					i, j, expectedRule.Path, actualRule.Path)
			}

			expectedHeaders := fmt.Sprintf("%v", expectedRule.Headers)
			actualHeaders := fmt.Sprintf("%v", actualRule.Headers)

			if !reflect.DeepEqual(actualRule.Headers, expectedRule.Headers) {
				t.Errorf("Test %d, rule %d: Expected headers %s, but got %s",
					i, j, expectedHeaders, actualHeaders)
			}
		}
	}
}
