package managers

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/mfojtik/depcheck/pkg/managers/dep"
	"github.com/mfojtik/depcheck/pkg/managers/glide"
	"github.com/mfojtik/depcheck/pkg/managers/version"
	"github.com/mfojtik/depcheck/pkg/managers/vgo"
	"github.com/mfojtik/depcheck/pkg/payload"
)

type ManifestType string

var (
	ManifestTypeGlide  ManifestType = "glide"
	ManifestTypeGodeps ManifestType = "godeps"
	ManifestTypeDep    ManifestType = "dep"
	ManifestTypeVGo    ManifestType = "vgo"
)

type RepositoryWithManifest struct {
	*payload.Repository

	Manifests    map[string][]byte
	ManifestType ManifestType

	Dependencies []version.Dependency
}

func (r *RepositoryWithManifest) GetVersions() error {
	var err error
	switch r.ManifestType {
	case ManifestTypeGlide:
		r.Dependencies, err = glide.ParseManifest(r.Manifests)
		fmt.Printf("%s: fetched glide %d dependencies\n", r.Name, len(r.Dependencies))
	case ManifestTypeDep:
		r.Dependencies, err = dep.ParseManifest(r.Manifests)
		fmt.Printf("%s: fetched dep %d dependencies\n", r.Name, len(r.Dependencies))
	case ManifestTypeVGo:
		r.Dependencies, err = vgo.ParseManifest(r.Manifests)
		fmt.Printf("%s: fetched vgo %d dependencies\n", r.Name, len(r.Dependencies))
	default:
		r.Dependencies = []version.Dependency{}
		fmt.Printf("%s: unhandled package manager\n", r.Name)
	}
	return err
}

func FetchManagerManifests(repos payload.Repositories) []*RepositoryWithManifest {
	var wg sync.WaitGroup
	result := []*RepositoryWithManifest{}

	for _, repo := range repos {
		wg.Add(1)
		go func(r payload.Repository) {
			defer wg.Done()
			if bytes, err := fetch(buildManifestURLs(ManifestTypeGlide, r)); err == nil {
				result = append(result, &RepositoryWithManifest{Repository: &r, Manifests: bytes, ManifestType: ManifestTypeGlide})
				return
			}
			if bytes, err := fetch(buildManifestURLs(ManifestTypeDep, r)); err == nil {
				result = append(result, &RepositoryWithManifest{Repository: &r, Manifests: bytes, ManifestType: ManifestTypeDep})
				return
			}
			if bytes, err := fetch(buildManifestURLs(ManifestTypeVGo, r)); err == nil {
				result = append(result, &RepositoryWithManifest{Repository: &r, Manifests: bytes, ManifestType: ManifestTypeVGo})
				return
			}
			if bytes, err := fetch(buildManifestURLs(ManifestTypeGodeps, r)); err == nil {
				result = append(result, &RepositoryWithManifest{Repository: &r, Manifests: bytes, ManifestType: ManifestTypeGodeps})
				return
			}
		}(repo)
	}

	wg.Wait()
	return result
}

func fetch(urls []string) (map[string][]byte, error) {
	result := map[string][]byte{}
	for _, u := range urls {
		parts := strings.Split(u, "/")
		name := parts[len(parts)-1]
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
		response, err := client.Get(u)
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := response.Body.Close(); err != nil {
			}
		}()
		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get %q, HTTP code: %d", u, response.StatusCode)
		}
		result[name], err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func getGithubRawURL(repoURL string) string {
	repoName := strings.TrimPrefix(repoURL, "https://github.com/")
	return "https://raw.githubusercontent.com/" + repoName + "/master/"
}

func buildManifestURLs(manifestType ManifestType, repo payload.Repository) []string {
	switch manifestType {
	case ManifestTypeGlide:
		return []string{
			getGithubRawURL(repo.URL) + "glide.yaml",
			getGithubRawURL(repo.URL) + "glide.lock",
		}
	case ManifestTypeDep:
		return []string{
			getGithubRawURL(repo.URL) + "Gopkg.lock",
		}
	case ManifestTypeVGo:
		return []string{getGithubRawURL(repo.URL) + "go.mod"}
	case ManifestTypeGodeps:
		return []string{getGithubRawURL(repo.URL) + "Godeps/Godeps.json"}
	default:
		panic("unknown type")
	}
}
