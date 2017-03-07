package github

import (
	"testing"
)

type versionTest struct {
	expected []string
	body     string
}

func TestExtractVersionFromBody(t *testing.T) {
	issues := []versionTest{
		{
			expected: []string{"Server:\n Version:      1.8.3", "1.8.3", ""},
			body: `
Client:
 Version:      1.8.3
 API version:  1.20
 Go version:   go1.4.2
 Git commit:   f4bf5c7
 Built:        Mon Oct 12 06:06:01 UTC 2015
 OS/Arch:      linux/amd64

Server:
 Version:      1.8.3
 API version:  1.20
 Go version:   go1.4.2
 Git commit:   f4bf5c7
 Built:        Mon Oct 12 06:06:01 UTC 2015
 OS/Arch:      linux/amd64
`,
		},
		{
			expected: []string{"Server:\n Version:         1.10.3-el7.centos", "1.10.3", "el7.centos"},
			body: `
Client:
 Version:         1.10.3-el7.centos
 API version:     1.22
 Package version: docker-1.10.3-10.el7.centos.x86_64
 Go version:      go1.4.2
 Git commit:      0b4a971-unsupported
 Built:           Tue Jun 21 17:51:37 2016
 OS/Arch:         linux/amd64

Server:
 Version:         1.10.3-el7.centos
 API version:     1.22
 Package version: docker-1.10.3-10.el7.centos.x86_64
 Go version:      go1.4.2
 Git commit:      0b4a971-unsupported
 Built:           Tue Jun 21 17:51:37 2016
 OS/Arch:         linux/amd64
`,
		},
		{
			expected: []string{"Server:\n Version:      1.11.2-cs5", "1.11.2", "cs5"},
			body: `
Client:
 Version:      1.11.2-cs5
 API version:  1.23
 Go version:   go1.5.4
 Git commit:   d364ea1
 Built:        Tue Sep 13 15:26:43 2016
 OS/Arch:      linux/amd64

Server:
 Version:      1.11.2-cs5
 API version:  1.23
 Go version:   go1.5.4
 Git commit:   d364ea1
 Built:        Tue Sep 13 15:26:43 2016
 OS/Arch:      linux/amd64
`,
		},
		{
			expected: []string{"Server:\n Version:      1.12.0-dev", "1.12.0", "dev"},
			body: `
Client:
 Version:      1.12.0-dev
 API version:  1.24
 Go version:   go1.5.4
 Git commit:   af60a9e-unsupported
 Built:        Tue May 17 02:04:00 2016
 OS/Arch:      linux/amd64

Server:
 Version:      1.12.0-dev
 API version:  1.24
 Go version:   go1.5.4
 Git commit:   af60a9e-unsupported
 Built:        Tue May 17 02:04:00 2016
 OS/Arch:      linux/amd64
`,
		},
		{
			expected: []string{"Server:\n Version:      1.13.0-rc4", "1.13.0", "rc4"},
			body: `
Client:
 Version:      1.13.0-rc4
 API version:  1.25
 Go version:   go1.7.3
 Git commit:   88862e7
 Built:        Sat Dec 17 01:34:17 2016
 OS/Arch:      darwin/amd64

Server:
 Version:      1.13.0-rc4
 API version:  1.25 (minimum version 1.12)
 Go version:   go1.7.3
 Git commit:   88862e7
 Built:        Sat Dec 17 01:34:17 2016
 OS/Arch:      linux/amd64
 Experimental: false
`,
		},
		{
			expected: []string{"Server:\n Version:      17.03.0-ce-rc1", "17.03.0", "ce-rc1"},
			body: `
Client:
 Version:      17.03.0-ce-rc1
 API version:  1.26
 Go version:   go1.7.5
 Git commit:   ce07fb6
 Built:        Mon Feb 20 10:12:38 2017
 OS/Arch:      darwin/amd64

Server:
 Version:      17.03.0-ce-rc1
 API version:  1.26 (minimum version 1.12)
 Go version:   go1.7.5
 Git commit:   ce07fb6
 Built:        Mon Feb 20 10:12:38 2017
 OS/Arch:      linux/amd64
 Experimental: true
`,
		},
		{
			expected: []string{"Server:\n Version:      17.03.0-ce", "17.03.0", "ce"},

			body: `
Client:
 Version:      17.03.0-ce
 API version:  1.26
 Go version:   go1.7.5
 Git commit:   60ccb22
 Built:        Thu Feb 23 10:40:59 2017
 OS/Arch:      windows/amd64

Server:
 Version:      17.03.0-ce
 API version:  1.26 (minimum version 1.12)
 Go version:   go1.7.5
 Git commit:   3a232c8
 Built:        Tue Feb 28 07:52:04 2017
 OS/Arch:      linux/amd64
 Experimental: true
`,
		},
	}

	for _, issue := range issues {
		versionMatch := extractVersionFromBody(issue.body)

		if len(versionMatch) != len(issue.expected) {
			t.Fatalf("expected %d matches, got %d (%q)\n", len(issue.expected), len(versionMatch), versionMatch)
		}

		for i, v := range issue.expected {
			if versionMatch[i] != v {
				t.Fatalf(`expected %q, got %q`, v, versionMatch[i])
			}
		}
	}
}

type versionLabelTest struct {
	version string
	suffix  string
	label   string
}

func TestLabelFromVersion(t *testing.T) {
	tests := []versionLabelTest{
		{
			version: "1.8.3",
			suffix:  "",
			label:   "version/1.8",
		},
		{
			version: "1.10.3",
			suffix:  "el7.centos",
			label:   "version/unsupported",
		},
		{
			version: "1.11.2",
			suffix:  "cs5",
			label:   "version/1.11",
		},
		{
			version: "1.12.0",
			suffix:  "dev",
			label:   "version/master",
		},
		{
			version: "1.13.0",
			suffix:  "rc4",
			label:   "version/1.13",
		},
		{
			version: "17.03.0",
			suffix:  "ce",
			label:   "version/17.03",
		},
		{
			version: "17.03.0",
			suffix:  "ce-rc1",
			label:   "version/17.03",
		},
	}

	for _, test := range tests {
		label := labelFromVersion(test.version, test.suffix)

		if label != test.label {
			t.Fatalf("expected %q matches, got %q\n", test.label, label)
		}
	}
}
