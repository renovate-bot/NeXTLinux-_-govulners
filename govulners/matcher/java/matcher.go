package java

import (
	"fmt"
	"net/http"

	"github.com/nextlinux/govulners/govulners/distro"
	"github.com/nextlinux/govulners/govulners/match"
	"github.com/nextlinux/govulners/govulners/pkg"
	"github.com/nextlinux/govulners/govulners/search"
	"github.com/nextlinux/govulners/govulners/vulnerability"
	"github.com/nextlinux/govulners/internal/log"
	syftPkg "github.com/anchore/syft/syft/pkg"
)

const (
	sha1Query = `1:"%s"`
)

type Matcher struct {
	MavenSearcher
	cfg MatcherConfig
}

type ExternalSearchConfig struct {
	SearchMavenUpstream bool
	MavenBaseURL        string
}

type MatcherConfig struct {
	ExternalSearchConfig
	UseCPEs bool
}

func NewJavaMatcher(cfg MatcherConfig) *Matcher {
	return &Matcher{
		cfg: cfg,
		MavenSearcher: &mavenSearch{
			client:  http.DefaultClient,
			baseURL: cfg.MavenBaseURL,
		},
	}
}

func (m *Matcher) PackageTypes() []syftPkg.Type {
	return []syftPkg.Type{syftPkg.JavaPkg, syftPkg.JenkinsPluginPkg}
}

func (m *Matcher) Type() match.MatcherType {
	return match.JavaMatcher
}

func (m *Matcher) Match(store vulnerability.Provider, d *distro.Distro, p pkg.Package) ([]match.Match, error) {
	var matches []match.Match
	if m.cfg.SearchMavenUpstream {
		upstreamMatches, err := m.matchUpstreamMavenPackages(store, p)
		if err != nil {
			log.Debugf("failed to match against upstream data for %s: %v", p.Name, err)
		} else {
			matches = append(matches, upstreamMatches...)
		}
	}
	criteria := search.CommonCriteria
	if m.cfg.UseCPEs {
		criteria = append(criteria, search.ByCPE)
	}
	criteriaMatches, err := search.ByCriteria(store, d, p, m.Type(), criteria...)
	if err != nil {
		return nil, fmt.Errorf("failed to match by exact package: %w", err)
	}

	matches = append(matches, criteriaMatches...)
	return matches, nil
}

func (m *Matcher) matchUpstreamMavenPackages(store vulnerability.Provider, p pkg.Package) ([]match.Match, error) {
	var matches []match.Match

	if metadata, ok := p.Metadata.(pkg.JavaMetadata); ok {
		for _, d := range metadata.ArchiveDigests {
			if d.Algorithm == "sha1" {
				indirectPackage, err := m.GetMavenPackageBySha(d.Value)
				if err != nil {
					return nil, err
				}
				indirectMatches, err := search.ByPackageLanguage(store, *indirectPackage, m.Type())
				if err != nil {
					return nil, err
				}
				matches = append(matches, indirectMatches...)
			}
		}
	}

	match.ConvertToIndirectMatches(matches, p)

	return matches, nil
}
