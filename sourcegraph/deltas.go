package sourcegraph

import (
	"encoding/base64"
	"strings"

	"sourcegraph.com/sourcegraph/go-diff/diff"
	"sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"sourcegraph.com/sourcegraph/srclib/store"
	"sourcegraph.com/sourcegraph/srclib/unit"
)

// DeltasService interacts with the delta-related endpoints of the
// Sourcegraph API. A delta is all of the changes between two commits,
// possibly from two different repositories. It includes the usual
// file diffs as well as definition-level diffs, affected author/repo
// impact information, etc.
type DeltasService interface {
	// Get fetches a summary of a delta.
	Get(ds DeltaSpec, opt *DeltaGetOptions) (*Delta, Response, error)

	// ListUnits lists units added/changed/deleted in a delta.
	ListUnits(ds DeltaSpec, opt *DeltaListUnitsOptions) ([]*UnitDelta, Response, error)

	// ListDefs lists definitions added/changed/deleted in a delta.
	ListDefs(ds DeltaSpec, opt *DeltaListDefsOptions) (*DeltaDefs, Response, error)

	// ListFiles fetches the file diff for a delta.
	ListFiles(ds DeltaSpec, opt *DeltaListFilesOptions) (*DeltaFiles, Response, error)

	// ListAffectedAuthors lists authors whose code is added/deleted/changed
	// in a delta.
	ListAffectedAuthors(ds DeltaSpec, opt *DeltaListAffectedAuthorsOptions) ([]*DeltaAffectedPerson, Response, error)

	// ListAffectedClients lists clients whose code is affected by a delta.
	ListAffectedClients(ds DeltaSpec, opt *DeltaListAffectedClientsOptions) ([]*DeltaAffectedPerson, Response, error)
}

// deltasService implements DeltasService.
type deltasService struct {
	client *Client
}

var _ DeltasService = &deltasService{}

// A DeltaSpec specifies a delta.
type DeltaSpec struct {
	Base RepoRevSpec
	Head RepoRevSpec
}

// RouteVars returns the route variables for generating URLs to the
// delta specified by this DeltaSpec.
func (s DeltaSpec) RouteVars() map[string]string {
	m := s.Base.RouteVars()

	if s.Base.RepoSpec == s.Head.RepoSpec {
		m["DeltaHeadRev"] = s.Head.RevPathComponent()
	} else {
		m["DeltaHeadRev"] = encodeCrossRepoRevSpecForDeltaHeadRev(s.Head)
	}
	return m
}

func encodeCrossRepoRevSpecForDeltaHeadRev(rr RepoRevSpec) string {
	return base64.URLEncoding.EncodeToString([]byte(rr.RepoSpec.PathComponent())) + ":" + rr.RevPathComponent()
}

// UnmarshalDeltaSpec marshals a map containing route variables
// generated by (*DeltaSpec).RouteVars() and returns the
// equivalent DeltaSpec struct.
func UnmarshalDeltaSpec(routeVars map[string]string) (DeltaSpec, error) {
	s := DeltaSpec{}

	rr, err := UnmarshalRepoRevSpec(routeVars)
	if err != nil {
		return DeltaSpec{}, err
	}
	s.Base = rr

	dhr := routeVars["DeltaHeadRev"]
	if i := strings.Index(dhr, ":"); i != -1 {
		// base repo != head repo
		repoPCB64, revPC := dhr[:i], dhr[i+1:]

		repoPC, err := base64.URLEncoding.DecodeString(repoPCB64)
		if err != nil {
			return DeltaSpec{}, err
		}

		rr, err := UnmarshalRepoRevSpec(map[string]string{"RepoSpec": string(repoPC), "Rev": revPC})
		if err != nil {
			return DeltaSpec{}, err
		}

		s.Head = rr
	} else {
		rr, err := UnmarshalRepoRevSpec(map[string]string{"RepoSpec": routeVars["RepoSpec"], "Rev": dhr})
		if err != nil {
			return DeltaSpec{}, err
		}

		s.Head = rr
	}
	return s, nil
}

// Delta represents the difference between two commits (possibly in 2
// separate repositories).
type Delta struct {
	Base, Head             RepoRevSpec // base/head repo and revspec
	BaseCommit, HeadCommit *Commit     // base/head commits
	BaseRepo, HeadRepo     *Repo       // base/head repositories
	BaseBuild, HeadBuild   *Build      // base/head builds (or nil)
}

func (d *Delta) DeltaSpec() DeltaSpec {
	return DeltaSpec{
		Base: d.Base,
		Head: d.Head,
	}
}

// BaseAndHeadBuildsSuccessful returns true iff both the base and head
// builds are present and ended successfully.
func (d *Delta) BaseAndHeadBuildsSuccessful() bool {
	return d.BaseBuild != nil && d.BaseBuild.Success && d.HeadBuild != nil && d.HeadBuild.Success
}

// DeltaGetOptions specifies options for getting a delta.
type DeltaGetOptions struct{}

func (s *deltasService) Get(ds DeltaSpec, opt *DeltaGetOptions) (*Delta, Response, error) {
	url, err := s.client.URL(router.Delta, ds.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var delta *Delta
	resp, err := s.client.Do(req, &delta)
	if err != nil {
		return nil, resp, err
	}

	return delta, resp, nil
}

// A UnitDelta represents a single source unit that was changed. It
// has fields for the before (Base) and after (Head) versions. If both
// Base and Head are non-nil, then the unit was changed from base to
// head. Otherwise, one of the fields being nil means that the unit
// did not exist in that revision (e.g., it was added or deleted from
// base to head).
type UnitDelta struct {
	Base *unit.SourceUnit // the unit in the base commit (if nil, this unit was added in the head)
	Head *unit.SourceUnit // the unit in the head commit (if nil, this unit was deleted in the head)
}

// Added is whether this represents an added source unit (not present
// in base, present in head).
func (ud UnitDelta) Added() bool { return ud.Base == nil && ud.Head != nil }

// Changed is whether this represents a changed source unit (present
// in base, present in head).
func (ud UnitDelta) Changed() bool { return ud.Base != nil && ud.Head != nil }

// Deleted is whether this represents a deleted source unit (present
// in base, not present in head).
func (ud UnitDelta) Deleted() bool { return ud.Base != nil && ud.Head == nil }

type UnitDeltas []*UnitDelta

func (v UnitDeltas) Len() int      { return len(v) }
func (v UnitDeltas) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v UnitDeltas) Less(i, j int) bool {
	a, b := v[i], v[j]
	return (a.Added() && b.Added() && deltaUnitLess(a.Head, b.Head)) || (a.Changed() && b.Changed() && deltaUnitLess(a.Head, b.Head)) || (a.Deleted() && b.Deleted() && deltaUnitLess(a.Base, b.Base)) || (a.Added() && !b.Added()) || (a.Changed() && !b.Added() && !b.Changed())
}

func deltaUnitLess(a, b *unit.SourceUnit) bool {
	return a.Type < b.Type || (a.Type == b.Type && a.Name < b.Name)
}

// DeltaListUnitsOptions specifies options for ListUnits.
type DeltaListUnitsOptions struct{}

func (s *deltasService) ListUnits(ds DeltaSpec, opt *DeltaListUnitsOptions) ([]*UnitDelta, Response, error) {
	url, err := s.client.URL(router.DeltaUnits, ds.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var units []*UnitDelta
	resp, err := s.client.Do(req, &units)
	if err != nil {
		return nil, resp, err
	}

	return units, resp, nil
}

// DeltaFilter specifies criteria by which to filter results from
// DeltaListXxx methods.
type DeltaFilter struct {
	Unit     string `url:",omitempty"`
	UnitType string `url:",omitempty"`
}

func (f DeltaFilter) DefFilters() []store.DefFilter {
	if f.UnitType != "" && f.Unit != "" {
		return []store.DefFilter{store.ByUnits(unit.ID2{Type: f.UnitType, Name: f.Unit})}
	}
	return nil
}

// DeltaListDefsOptions specifies options for ListDefs.
type DeltaListDefsOptions struct {
	DeltaFilter
	ListOptions
}

// DeltaDefs describes definitions added/changed/deleted in a delta.
type DeltaDefs struct {
	Defs []*DefDelta // added/changed/deleted defs

	DiffStat diff.Stat // overall diffstat (not subject to pagination)
}

// A DefDelta represents a single definition that was changed. It has
// fields for the before (Base) and after (Head) versions. If both
// Base and Head are non-nil, then the def was changed from base to
// head. Otherwise, one of the fields being nil means that the def did
// not exist in that revision (e.g., it was added or deleted from base
// to head).
type DefDelta struct {
	Base *Def // the def in the base commit (if nil, this def was added in the head)
	Head *Def // the def in the head commit (if nil, this def was deleted in the head)
}

// Added is whether this represents an added def (not present in base,
// present in head).
func (dd DefDelta) Added() bool { return dd.Base == nil && dd.Head != nil }

// Changed is whether this represents a changed def (present in base,
// present in head).
func (dd DefDelta) Changed() bool { return dd.Base != nil && dd.Head != nil }

// Deleted is whether this represents a deleted def (present in base,
// not present in head).
func (dd DefDelta) Deleted() bool { return dd.Base != nil && dd.Head == nil }

func (v DeltaDefs) Len() int      { return len(v.Defs) }
func (v DeltaDefs) Swap(i, j int) { v.Defs[i], v.Defs[j] = v.Defs[j], v.Defs[i] }
func (v DeltaDefs) Less(i, j int) bool {
	a, b := v.Defs[i], v.Defs[j]
	return (a.Added() && b.Added() && deltaDefLess(a.Head, b.Head)) || (a.Changed() && b.Changed() && deltaDefLess(a.Head, b.Head)) || (a.Deleted() && b.Deleted() && deltaDefLess(a.Base, b.Base)) || (a.Added() && !b.Added()) || (a.Changed() && !b.Added() && !b.Changed())
}

func deltaDefLess(a, b *Def) bool {
	return a.UnitType < b.UnitType || (a.UnitType == b.UnitType && a.Unit < b.Unit) || (a.UnitType == b.UnitType && a.Unit == b.Unit && a.Path < b.Path)
}

func (s *deltasService) ListDefs(ds DeltaSpec, opt *DeltaListDefsOptions) (*DeltaDefs, Response, error) {
	url, err := s.client.URL(router.DeltaDefs, ds.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var defs *DeltaDefs
	resp, err := s.client.Do(req, &defs)
	if err != nil {
		return nil, resp, err
	}

	return defs, resp, nil
}

// DeltaListFilesOptions specifies options for
// ListFiles.
type DeltaListFilesOptions struct {
	// Formatted is whether the files should have their contents
	// code-formatted (syntax-highlighted and reference-linked) if
	// they contain code.
	Formatted bool `url:",omitempty"`

	// Filter filters the list of returned files to those whose name matches Filter.
	Filter string `url:",omitempty"`

	DeltaFilter
}

// DeltaFiles describes files added/changed/deleted in a delta.
type DeltaFiles struct {
	FileDiffs []*diff.FileDiff
}

// DiffStat returns a diffstat that is the sum of all of the files'
// diffstats.
func (d *DeltaFiles) DiffStat() diff.Stat {
	ds := diff.Stat{}
	for _, fd := range d.FileDiffs {
		st := fd.Stat()
		ds.Added += st.Added
		ds.Changed += st.Changed
		ds.Deleted += st.Deleted
	}
	return ds
}

func (s *deltasService) ListFiles(ds DeltaSpec, opt *DeltaListFilesOptions) (*DeltaFiles, Response, error) {
	url, err := s.client.URL(router.DeltaFiles, ds.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var files *DeltaFiles
	resp, err := s.client.Do(req, &files)
	if err != nil {
		return nil, resp, err
	}

	return files, resp, nil
}

// DeltaAffectedPerson describes a person (registered user or
// committer email address) that is affected by a delta. It includes
// fields for the person affected as well as the defs that are the
// reason why we consider them to be affected.
//
// The person's relationship to the Defs depends on what method
// returned this DeltaAffectedPerson. If it was returned by a method
// that lists authors, then the Defs are definitions that the Person
// committed. If it was returned by a method that lists clients (a.k.a
// users), then the Defs are definitions that the Person uses.
type DeltaAffectedPerson struct {
	Person // the affected person

	Defs []*Def // the defs they authored or use (the reason why they're affected)
}

// DeltaListAffectedAuthorsOptions specifies options for
// ListAffectedAuthors.
type DeltaListAffectedAuthorsOptions struct {
	DeltaFilter
	ListOptions
}

func (s *deltasService) ListAffectedAuthors(ds DeltaSpec, opt *DeltaListAffectedAuthorsOptions) ([]*DeltaAffectedPerson, Response, error) {
	url, err := s.client.URL(router.DeltaAffectedAuthors, ds.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var authors []*DeltaAffectedPerson
	resp, err := s.client.Do(req, &authors)
	if err != nil {
		return nil, resp, err
	}

	return authors, resp, nil
}

// DeltaListAffectedClientsOptions specifies options for
// ListAffectedClients.
type DeltaListAffectedClientsOptions struct {
	DeltaFilter
	ListOptions
}

func (s *deltasService) ListAffectedClients(ds DeltaSpec, opt *DeltaListAffectedClientsOptions) ([]*DeltaAffectedPerson, Response, error) {
	url, err := s.client.URL(router.DeltaAffectedClients, ds.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var clients []*DeltaAffectedPerson
	resp, err := s.client.Do(req, &clients)
	if err != nil {
		return nil, resp, err
	}

	return clients, resp, nil
}


