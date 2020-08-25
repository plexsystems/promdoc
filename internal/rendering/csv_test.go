package rendering

import (
	"io/ioutil"
	"testing"
)

func TestCSV(t *testing.T) {
	expected, err := ioutil.ReadFile("../../examples/csv/expected.csv")
	if err != nil {
		t.Fatal("read expected file")
	}

	actual, err := Render("../../examples/rule.yaml", ".csv")
	if err != nil {
		t.Fatal("getting rule groups:", err)
	}

	if string(expected) != actual {
		t.Error("rendered markdown did not match expected output")
	}
}
