package payload

import (
	"io/ioutil"
	"testing"
)

func TestReadPayloadJSONOrDie(t *testing.T) {
	payloadBytes, err := ioutil.ReadFile("testdata/payload.json")
	if err != nil {
		t.Fatal(err)
	}
	payload, err := ReadPayloadJSON(payloadBytes)
	if err != nil {
		t.Fatal(err)
	}
	if len(payload.References.Spec.Tags) == 0 {
		t.Errorf("no tags")
	}
}
