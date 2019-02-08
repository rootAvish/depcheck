package payload

import (
	"encoding/json"

	"github.com/openshift/api/image/v1"
)

type Payload struct {
	Image  string `json:"image"`
	Digest string `json:"digest"`

	References v1.ImageStream `json:"references"`
}

type Repository struct {
	Name   string
	URL    string
	Commit string
}

type Repositories []Repository

func (r *Repositories) Add(name, url, commit string) {
	repositories := *r
	repositories = append(repositories, Repository{
		Name:   name,
		URL:    url,
		Commit: commit,
	})
	*r = repositories
}

func ReadPayloadJSON(payloadBytes []byte) (*Payload, error) {
	payload := Payload{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func ParseRepositoriesFromPayload(payload *Payload) *Repositories {
	repositories := &Repositories{}
	for _, tag := range payload.References.Spec.Tags {
		repositories.Add(
			tag.Name,
			tag.Annotations["io.openshift.build.source-location"],
			tag.Annotations["io.openshift.build.commit.id"],
		)
	}
	return repositories
}
