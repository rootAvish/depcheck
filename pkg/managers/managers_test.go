package managers

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/rootAvish/depcheck/pkg/managers/version"
	"github.com/rootAvish/depcheck/pkg/payload"
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

func TestRepositoryWithManifest_GetVersions(t *testing.T) {
	type fields struct {
		Repository   *payload.Repository
		Manifests    map[string][]byte
		ManifestType ManifestType
		Dependencies []version.Dependency
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Test a glide manifest",
			fields: fields{
				Repository: &payload.Repository{
					URL:    "github.com/kubernetes/apimachinery",
					Name:   "apimachinery",
					Commit: "2ff2755528372bae36ede36ddccd1a77bd118887",
				},
				Manifests: map[string][]byte{
					"go.mod": []byte(`// This is a generated file. Do not edit directly.

				module k8s.io/apimachinery
				
				go 1.12
				
				require (
					github.com/davecgh/go-spew v1.1.1
					github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96
					github.com/elazarl/goproxy v0.0.0-20170405201442-c4fc26588b6e
					github.com/evanphx/json-patch v0.0.0-20190203023257-5858425f7550
					github.com/gogo/protobuf v0.0.0-20171007142547-342cbe0a0415
					github.com/golang/groupcache v0.0.0-20160516000752-02826c3e7903
					github.com/golang/protobuf v1.2.0
					github.com/google/go-cmp v0.3.0
					github.com/google/gofuzz v0.0.0-20170612174753-24818f796faf
					github.com/google/uuid v1.0.0
					github.com/googleapis/gnostic v0.0.0-20170729233727-0c5108395e2d
					github.com/hashicorp/golang-lru v0.5.0
					github.com/json-iterator/go v0.0.0-20180701071628-ab8a2e0c74be
					github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
					github.com/modern-go/reflect2 v1.0.1
					github.com/mxk/go-flowrate v0.0.0-20140419014527-cca7078d478f
					github.com/onsi/gomega v0.0.0-20190113212917-5533ce8a0da3 // indirect
					github.com/pmezard/go-difflib v1.0.0 // indirect
					github.com/spf13/pflag v1.0.1
					github.com/stretchr/testify v1.2.2
					golang.org/x/net v0.0.0-20190206173232-65e2d4e15006
					golang.org/x/sync v0.0.0-20181221193216-37e7f081c4d4 // indirect
					golang.org/x/sys v0.0.0-20190312061237-fead79001313 // indirect
					golang.org/x/text v0.3.1-0.20181227161524-e6919f6577db // indirect
					gopkg.in/inf.v0 v0.9.0
					gopkg.in/yaml.v2 v2.2.1
					k8s.io/klog v0.3.0
					k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
					sigs.k8s.io/yaml v1.1.0
				)
				
				replace (
					golang.org/x/sync => golang.org/x/sync v0.0.0-20181108010431-42b317875d0f
					golang.org/x/sys => golang.org/x/sys v0.0.0-20190209173611-3b5209105503
					golang.org/x/tools => golang.org/x/tools v0.0.0-20190313210603-aa82965741a9
				)
				`),
				},
				ManifestType: "vgo",
				Dependencies: []version.Dependency{},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RepositoryWithManifest{
				Repository:   tt.fields.Repository,
				Manifests:    tt.fields.Manifests,
				ManifestType: tt.fields.ManifestType,
				Dependencies: tt.fields.Dependencies,
			}
			if err := r.GetVersions(); (err != nil) != tt.wantErr {
				t.Errorf("RepositoryWithManifest.GetVersions() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Printf("%v\n", r.Dependencies)
		})
	}
}
