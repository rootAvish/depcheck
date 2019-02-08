package managers

import (
	"io/ioutil"
	"testing"

	"github.com/mfojtik/depcheck/pkg/payload"
)

func TestFetchManagerManifests(t *testing.T) {
	payloadBytes, err := ioutil.ReadFile("../payload/testdata/payload.json")
	if err != nil {
		t.Fatal(err)
	}
	p, err := payload.ReadPayloadJSON(payloadBytes)
	if err != nil {
		t.Fatal(err)
	}

	repositories := payload.ParseRepositoriesFromPayload(p)

	result := FetchManagerManifests(*repositories)

	for _, r := range result {
		if len(r.Manifests) == 0 {
			t.Errorf("%s: manifest is empty", r.Name)
		}

		if r.ManifestType == ManifestTypeGlide {
			if err := r.GetVersions(); err != nil {
				t.Errorf("%s: unable to get versions: %v", r.Name, err)
				continue
			}
			t.Logf("%s: %+v", r.Name, r.Dependencies)
		}
	}
}
