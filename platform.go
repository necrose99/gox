package main

import (
	"fmt"
	"log"
	"strings"

	version "github.com/hashicorp/go-version"
)

// Platform is a combination of OS/arch that can be built against.
type Platform struct {
	OS   string
	Arch string

	// Default, if true, will be included as a default build target
	// if no OS/arch is specified. We try to only set as a default popular
	// targets or targets that are generally useful. For example, Android
	// is not a default because it is quite rare that you're cross-compiling
	// something to Android AND something like Linux.
	Default bool
}

func (p *Platform) String() string {
	return fmt.Sprintf("%s/%s", p.OS, p.Arch)
}

/// Like `uname -s`
// Matches https://github.com/golang/go/blob/master/src/go/build/syslist.go
func (p *Platform) OSUname() string {
	return map[string]string{
		//"android":
		"darwin":    "Darwin",
		"dragonfly": "DragonFly",
		"freebsd":   "FreeBSD",
		"linux":     "Linux",
		//"nacl":
		"netbsd":    "NetBSD",
		"openbsd":   "OpenBSD",
		"plan9":     "Plan9",
		"solaris":   "SunOS",
		"windows":   "Windows",
		//"zos":
	}[p.OS]
}

/// Like `uname -m`
// Matches https://github.com/golang/go/blob/master/src/go/build/syslist.go
func (p *Platform) ArchUname() string {
	return map[string]string{
		"386":     "i386",
		"amd64":   "x86_64",
		//"amd64p32":
		"arm":     "arm",
		//"armbe":
		"arm64":   "aarch64",
		//"arm64be":
		"ppc64":   "ppc64",
		"ppc64le": "ppc64le",
		//"mips":
		//"mipsle":
		//"mips64":
		//"mips64p32":
		//"mips64p32le":
		//"ppc":
		//"s390":
		//"s390x":
		//"sparc":
		//"sparc64":
	}[p.Arch]
}

var (
	Platforms_1_0 = []Platform{
		{"darwin", "386", true},
		{"darwin", "amd64", true},
		{"linux", "386", true},
		{"linux", "amd64", true},
		{"linux", "arm", true},
		{"freebsd", "386", true},
		{"freebsd", "amd64", true},
		{"openbsd", "386", true},
		{"openbsd", "amd64", true},
		{"windows", "386", true},
		{"windows", "amd64", true},
	}

	Platforms_1_1 = append(Platforms_1_0, []Platform{
		{"freebsd", "arm", true},
		{"netbsd", "386", true},
		{"netbsd", "amd64", true},
		{"netbsd", "arm", true},
		{"plan9", "386", false},
	}...)

	Platforms_1_3 = append(Platforms_1_1, []Platform{
		{"dragonfly", "386", false},
		{"dragonfly", "amd64", false},
		{"nacl", "amd64", false},
		{"nacl", "amd64p32", false},
		{"nacl", "arm", false},
		{"solaris", "amd64", false},
	}...)

	Platforms_1_4 = append(Platforms_1_3, []Platform{
		{"android", "arm", false},
		{"plan9", "amd64", false},
	}...)

	Platforms_1_5 = append(Platforms_1_4, []Platform{
		{"darwin", "arm", false},
		{"darwin", "arm64", false},
		{"linux", "arm64", false},
		{"linux", "ppc64", false},
		{"linux", "ppc64le", false},
	}...)

	Platforms_1_6 = append(Platforms_1_5, []Platform{
		{"android", "386", false},
		{"linux", "mips64", false},
		{"linux", "mips64le", false},
	}...)

	Platforms_1_7 = append(Platforms_1_5, []Platform{
		// While not fully supported s390x is generally useful
		{"linux", "s390x", true},
		{"plan9", "arm", false},
		// Add the 1.6 Platforms, but reflect full support for mips64 and mips64le
		{"android", "386", false},
		{"linux", "mips64", true},
		{"linux", "mips64le", true},
	}...)

	Platforms_1_8 = append(Platforms_1_7, []Platform{
		{"linux", "mips", true},
		{"linux", "mipsle", true},
		{"linux", "arm64", true},
		//{"windows", "arm", true},  //not yet 
		//{"windows", "arm64", true}, //wont yet cross compile 

	}...)

	Platforms_1_9 = append(Platforms_1_8, []Platform{
		{"linux", "riscv64", true},
		{"freebsd", "riscv64", true},
		{"freebsd", "arm64", true},
		{"freebsd", "arm", true},
		{"openbsd", "arm64", true},
		{"openbsd", "arm", true},
		{"openbsd", "riscv64", true},
		{"windows", "arm", true},
		{"windows", "arm64", true},
		{"js", "wasm", true},
	}...)
	// no new platforms in 1.10
	Platforms_1_10 = Platforms_1_9

	Platforms_1_11 = append(Platforms_1_10, []Platform{
		{"js", "wasm", true},
		// Not sure arm64 was ported in 1.11 maybe before!
		{"linux", "arm64", true},
	}...)

	Platforms_1_12 = append(Platforms_1_11, []Platform{
		{"linux", "ppc64", true},
		{"windows", "arm", true},
		{"aix", "ppc64", true},
	}...)

	PlatformsLatest = Platforms_1_12
)

// SupportedPlatforms returns the full list of supported platforms for
// the version of Go that is
func SupportedPlatforms(v string) []Platform {
	// Use latest if we get an unexpected version string
	if !strings.HasPrefix(v, "go") {
		return PlatformsLatest
	}
	// go-version only cares about version numbers
	v = v[2:]

	current, err := version.NewVersion(v)
	if err != nil {
		log.Printf("Unable to parse current go version: %s\n%s", v, err.Error())

		// Default to latest
		return PlatformsLatest
	}

	var platforms = []struct {
		constraint string
		plat       []Platform
	}{
		{"<= 1.0", Platforms_1_0},
		{">= 1.1, < 1.3", Platforms_1_1},
		{">= 1.3, < 1.4", Platforms_1_3},
		{">= 1.4, < 1.5", Platforms_1_4},
		{">= 1.5, < 1.6", Platforms_1_5},
		{">= 1.6, < 1.7", Platforms_1_6},
		{">= 1.7, < 1.8", Platforms_1_7},
		{">= 1.8, < 1.9", Platforms_1_8},
		{">= 1.9, < 1.10", Platforms_1_9},
		{">=1.10, < 1.11", Platforms_1_10},
		{">=1.11, < 1.12", Platforms_1_11},
		{">=1.12, < 1.13", Platforms_1_12},
	}

	for _, p := range platforms {
		constraints, err := version.NewConstraint(p.constraint)
		if err != nil {
			panic(err)
		}
		if constraints.Check(current) {
			return p.plat
		}
	}

	// Assume latest
	return Platforms_1_12
}
