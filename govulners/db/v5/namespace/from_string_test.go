package namespace

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nextlinux/govulners/govulners/db/v5/namespace/cpe"
	"github.com/nextlinux/govulners/govulners/db/v5/namespace/distro"
	"github.com/nextlinux/govulners/govulners/db/v5/namespace/language"
	govulnersDistro "github.com/nextlinux/govulners/govulners/distro"
	syftPkg "github.com/anchore/syft/syft/pkg"
)

func TestFromString(t *testing.T) {
	tests := []struct {
		namespaceString string
		result          Namespace
	}{
		{
			namespaceString: "github:language:python",
			result:          language.NewNamespace("github", syftPkg.Python, ""),
		},
		{
			namespaceString: "github:language:python:python",
			result:          language.NewNamespace("github", syftPkg.Python, syftPkg.PythonPkg),
		},
		{
			namespaceString: "debian:distro:debian:8",
			result:          distro.NewNamespace("debian", govulnersDistro.Debian, "8"),
		},
		{
			namespaceString: "unknown:distro:amazonlinux:2022.15",
			result:          distro.NewNamespace("unknown", govulnersDistro.AmazonLinux, "2022.15"),
		},
		{
			namespaceString: "ns-1:distro:unknowndistro:abcdefg~~~",
			result:          distro.NewNamespace("ns-1", govulnersDistro.Type("unknowndistro"), "abcdefg~~~"),
		},
		{
			namespaceString: "abc.xyz:cpe",
			result:          cpe.NewNamespace("abc.xyz"),
		},
	}

	for _, test := range tests {
		result, _ := FromString(test.namespaceString)
		assert.Equal(t, result, test.result)
	}
}
