package rendering

import (
	"io/ioutil"
	"testing"
)

func TestCSV(t *testing.T) {
	expectedBytes, err := ioutil.ReadFile("../../examples/csv/expected.csv")
	if err != nil {
		t.Fatal("read:", err)
	}

	actual, err := Render("../../examples", ".csv")
	if err != nil {
		t.Fatal("render:", err)
	}

	expected := string(expectedBytes)
	if expected != actual {
		t.Errorf("Unexpected CSV. expected %v, actual %v", expected, actual)
	}
}
