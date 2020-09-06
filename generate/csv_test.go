package generate

import (
	"io/ioutil"
	"testing"
)

func TestCSV(t *testing.T) {
	expectedBytes, err := ioutil.ReadFile("../examples/csv/expected.csv")
	if err != nil {
		t.Fatal("read:", err)
	}

	actual, err := Generate("../examples", ".csv")
	if err != nil {
		t.Fatal("generate:", err)
	}

	expected := string(expectedBytes)
	if expected != actual {
		t.Errorf("Unexpected CSV. expected %v, actual %v", expected, actual)
	}
}
