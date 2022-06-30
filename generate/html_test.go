package generate

import (
	"io/ioutil"
	"testing"
)

func TestHTML(t *testing.T) {
	expectedBytes, err := ioutil.ReadFile("../examples/html/expected.html")
	if err != nil {
		t.Fatal("read:", err)
	}

	actual, err := Generate("../examples", ".html", "kubernetes")
	if err != nil {
		t.Fatal("generate:", err)
	}

	expected := string(expectedBytes)
	if expected != actual {
		t.Errorf("Unexpected HTML. expected %v, actual %v", expected, actual)
	}
}
