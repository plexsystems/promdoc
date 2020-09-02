package rendering

import (
	"io/ioutil"
	"testing"
)

func testCSV(t *testing.T, filename string) {
	expected, err := ioutil.ReadFile("../../examples/csv/expected.csv")
	if err != nil {
		t.Fatal("read expected file")
	}

	actual, err := Render(filename, ".csv")
	if err != nil {
		t.Fatal("getting rule groups:", err)
	}

	if string(expected) != actual {
		t.Error("rendered markdown did not match expected output")
	}
}
func TestCSVWithDescription(t *testing.T) {
	testCSV(t, "../../examples/rule.yaml")
}

func TestCSVWithMessage(t *testing.T) {
	testCSV(t, "../../examples/rule2.yaml")
}
