package version

// BumpKind selects which core component a bump increments.
type BumpKind int

// Bump kinds.
const (
	BumpMajor BumpKind = iota
	BumpMinor
	BumpPatch
)

// IncMajor returns the next major version. Build metadata is cleared. Following the
// node-semver convention, a pre-release that already anchors the target (its minor
// and patch are zero) collapses to that stable version (e.g. v1.0.0-rc.1 -> v1.0.0);
// otherwise the major is incremented and the pre-release dropped.
func (v Version) IncMajor() Version {
	out := v.clone()
	out.Build = nil
	if v.IsPrerelease() && v.Minor == 0 && v.Patch == 0 {
		out.Pre = nil
		return out
	}
	out.Major, out.Minor, out.Patch, out.Pre = v.Major+1, 0, 0, nil
	return out
}

// IncMinor returns the next minor version. A pre-release whose patch is zero
// collapses to its stable minor (e.g. v1.2.0-rc.1 -> v1.2.0); otherwise the minor is
// incremented. See IncMajor for the general rule.
func (v Version) IncMinor() Version {
	out := v.clone()
	out.Build = nil
	if v.IsPrerelease() && v.Patch == 0 {
		out.Pre = nil
		return out
	}
	out.Minor, out.Patch, out.Pre = v.Minor+1, 0, nil
	return out
}

// IncPatch returns the next patch version. Any pre-release collapses to its stable
// patch (e.g. v1.2.3-rc.1 -> v1.2.3); otherwise the patch is incremented. See
// IncMajor for the general rule.
func (v Version) IncPatch() Version {
	out := v.clone()
	out.Build = nil
	if v.IsPrerelease() {
		out.Pre = nil
		return out
	}
	out.Patch, out.Pre = v.Patch+1, nil
	return out
}

func (v Version) bump(k BumpKind) Version {
	switch k {
	case BumpMajor:
		return v.IncMajor()
	case BumpMinor:
		return v.IncMinor()
	default:
		return v.IncPatch()
	}
}
