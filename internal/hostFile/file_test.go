package hostfile

import (
	"testing"

	iplookup "github.com/pkf/hosts/internal/ipLookup"
)

func Test_update(t *testing.T) {
	var domains = []string{
		"alive.github.com",
		"live.github.com",
		"github.githubassets.com",
		"central.github.com",
		"desktop.githubusercontent.com",
		"assets-cdn.github.com",
		"camo.githubusercontent.com",
		"github.map.fastly.net",
		"github.global.ssl.fastly.net",
		"gist.github.com",
		"github.io",
		"github.com",
		"github.blog",
		"api.github.com",
		"raw.githubusercontent.com",
		"user-images.githubusercontent.com",
		"favicons.githubusercontent.com",
		"avatars5.githubusercontent.com",
		"avatars4.githubusercontent.com",
		"avatars3.githubusercontent.com",
		"avatars2.githubusercontent.com",
		"avatars1.githubusercontent.com",
		"avatars0.githubusercontent.com",
		"avatars.githubusercontent.com",
		"codeload.github.com",
		"github-cloud.s3.amazonaws.com",
		"github-com.s3.amazonaws.com",
		"github-production-release-asset-2e65be.s3.amazonaws.com",
		"github-production-user-asset-6210df.s3.amazonaws.com",
		"github-production-repository-file-5c1aeb.s3.amazonaws.com",
		"githubstatus.com",
		"github.community",
		"github.dev",
		"media.githubusercontent.com",
	}
	hosts := iplookup.LookupDomains(domains)

	update("./case/hosts", "./case/hosts.txt", hosts)
}
