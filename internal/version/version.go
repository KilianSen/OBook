package version

import (
	"strconv"
)

type Version struct {
	rsc   uint8
	patch uint8
	minor uint8
	major uint8

	backwardCompatibleUntil *Version
	hopId                   uint8
}

func (v Version) getWithoutBackwardCompatible() Version {
	return Version{
		rsc:                     v.rsc,
		patch:                   v.patch,
		minor:                   v.minor,
		major:                   v.major,
		backwardCompatibleUntil: nil,
	}
}

func (v Version) toString() string {
	var bcs string
	if v.backwardCompatibleUntil != nil {
		bcs = v.backwardCompatibleUntil.getWithoutBackwardCompatible().toString()
	}

	return strconv.FormatUint(uint64(v.major), 10) +
		"." + strconv.FormatUint(uint64(v.minor), 10) +
		"." + strconv.FormatUint(uint64(v.patch), 10) +
		"." + strconv.FormatUint(uint64(v.rsc), 10) +
		"-" + bcs +
		"-" + strconv.FormatUint(uint64(v.hopId), 10)
}

func (v Version) toFormatString() string {
	baseVersion := strconv.FormatUint(uint64(v.major), 10) +
		"." + strconv.FormatUint(uint64(v.minor), 10) +
		"." + strconv.FormatUint(uint64(v.patch), 10)

	var releaseChannel string
	if v.rsc == 0 {
		releaseChannel = "stable"
	}
	if v.rsc == 1 {
		releaseChannel = "rc"
	}
	if v.rsc == 2 {
		releaseChannel = "beta"
	}
	if v.rsc == 3 {
		releaseChannel = "alpha"
	}
	if v.rsc == 4 {
		releaseChannel = "nightly"
	}

	backwardCompatible := v.backwardCompatibleUntil.getWithoutBackwardCompatible().toString() +
		"-" + strconv.FormatUint(uint64(v.hopId), 10)

	return baseVersion + "-" + releaseChannel + "-" + backwardCompatible

}

// Compare two versions, return 1 if v1 > v2, -1 if v1 < v2, 0 if v1 == v2
func compareVersion(v1 Version, v2 Version) int {
	if v1.major > v2.major {
		return 1
	}
	if v1.major < v2.major {
		return -1
	}

	if v1.minor > v2.minor {
		return 1
	}
	if v1.minor < v2.minor {
		return -1
	}

	if v1.patch > v2.patch {
		return 1
	}
	if v1.patch < v2.patch {
		return -1
	}

	if v1.rsc > v2.rsc {
		return 1
	}
	if v1.rsc < v2.rsc {
		return -1
	}

	return 0
}

// Method is to be used to check if the current module version is able to update an external version to the current version
func (v Version) isCompatibleWith(external Version) bool {
	if compareVersion(v, external) == -1 {
		return false // The current version is older than the external version
	}

	if compareVersion(v, external) == 0 {
		return true // The current version is the same as the external version
	}

	if v.backwardCompatibleUntil == nil {
		return false // no backward compatibility defined
	}

	if compareVersion(v.backwardCompatibleUntil.getWithoutBackwardCompatible(), external) <= 0 {
		return true // The external version is within the backward compatibility range
	}

	return false // The external version is outside the backward compatibility range, and needs a compatibility hop.

}
