// generated by gen-mocks; DO NOT EDIT

package mock

import (
	"golang.org/x/net/context"
	"sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"
	"sourcegraph.com/sourcegraph/go-vcs/vcs"
	"sourcegraph.com/sqs/pbtypes"
)

type RepoBadgesServer struct {
	ListBadges_   func(v0 context.Context, v1 *sourcegraph.RepoSpec) (*sourcegraph.BadgeList, error)
	ListCounters_ func(v0 context.Context, v1 *sourcegraph.RepoSpec) (*sourcegraph.CounterList, error)
	RecordHit_    func(v0 context.Context, v1 *sourcegraph.RepoSpec) (*pbtypes.Void, error)
	CountHits_    func(v0 context.Context, v1 *sourcegraph.RepoBadgesCountHitsOp) (*sourcegraph.RepoBadgesCountHitsResult, error)
}

func (s *RepoBadgesServer) ListBadges(v0 context.Context, v1 *sourcegraph.RepoSpec) (*sourcegraph.BadgeList, error) {
	return s.ListBadges_(v0, v1)
}

func (s *RepoBadgesServer) ListCounters(v0 context.Context, v1 *sourcegraph.RepoSpec) (*sourcegraph.CounterList, error) {
	return s.ListCounters_(v0, v1)
}

func (s *RepoBadgesServer) RecordHit(v0 context.Context, v1 *sourcegraph.RepoSpec) (*pbtypes.Void, error) {
	return s.RecordHit_(v0, v1)
}

func (s *RepoBadgesServer) CountHits(v0 context.Context, v1 *sourcegraph.RepoBadgesCountHitsOp) (*sourcegraph.RepoBadgesCountHitsResult, error) {
	return s.CountHits_(v0, v1)
}

var _ sourcegraph.RepoBadgesServer = (*RepoBadgesServer)(nil)

type RepoStatusesServer struct {
	Create_      func(v0 context.Context, v1 *sourcegraph.RepoStatusesCreateOp) (*sourcegraph.RepoStatus, error)
	GetCombined_ func(v0 context.Context, v1 *sourcegraph.RepoRevSpec) (*sourcegraph.CombinedStatus, error)
}

func (s *RepoStatusesServer) Create(v0 context.Context, v1 *sourcegraph.RepoStatusesCreateOp) (*sourcegraph.RepoStatus, error) {
	return s.Create_(v0, v1)
}

func (s *RepoStatusesServer) GetCombined(v0 context.Context, v1 *sourcegraph.RepoRevSpec) (*sourcegraph.CombinedStatus, error) {
	return s.GetCombined_(v0, v1)
}

var _ sourcegraph.RepoStatusesServer = (*RepoStatusesServer)(nil)

type ReposServer struct {
	Get_          func(v0 context.Context, v1 *sourcegraph.RepoSpec) (*sourcegraph.Repo, error)
	List_         func(v0 context.Context, v1 *sourcegraph.RepoListOptions) (*sourcegraph.RepoList, error)
	GetReadme_    func(v0 context.Context, v1 *sourcegraph.RepoRevSpec) (*sourcegraph.Readme, error)
	Enable_       func(v0 context.Context, v1 *sourcegraph.RepoSpec) (*pbtypes.Void, error)
	Disable_      func(v0 context.Context, v1 *sourcegraph.RepoSpec) (*pbtypes.Void, error)
	GetCommit_    func(v0 context.Context, v1 *sourcegraph.RepoRevSpec) (*vcs.Commit, error)
	ListCommits_  func(v0 context.Context, v1 *sourcegraph.ReposListCommitsOp) (*sourcegraph.CommitList, error)
	ListBranches_ func(v0 context.Context, v1 *sourcegraph.ReposListBranchesOp) (*sourcegraph.BranchList, error)
	ListTags_     func(v0 context.Context, v1 *sourcegraph.ReposListTagsOp) (*sourcegraph.TagList, error)
}

func (s *ReposServer) Get(v0 context.Context, v1 *sourcegraph.RepoSpec) (*sourcegraph.Repo, error) {
	return s.Get_(v0, v1)
}

func (s *ReposServer) List(v0 context.Context, v1 *sourcegraph.RepoListOptions) (*sourcegraph.RepoList, error) {
	return s.List_(v0, v1)
}

func (s *ReposServer) GetReadme(v0 context.Context, v1 *sourcegraph.RepoRevSpec) (*sourcegraph.Readme, error) {
	return s.GetReadme_(v0, v1)
}

func (s *ReposServer) Enable(v0 context.Context, v1 *sourcegraph.RepoSpec) (*pbtypes.Void, error) {
	return s.Enable_(v0, v1)
}

func (s *ReposServer) Disable(v0 context.Context, v1 *sourcegraph.RepoSpec) (*pbtypes.Void, error) {
	return s.Disable_(v0, v1)
}

func (s *ReposServer) GetCommit(v0 context.Context, v1 *sourcegraph.RepoRevSpec) (*vcs.Commit, error) {
	return s.GetCommit_(v0, v1)
}

func (s *ReposServer) ListCommits(v0 context.Context, v1 *sourcegraph.ReposListCommitsOp) (*sourcegraph.CommitList, error) {
	return s.ListCommits_(v0, v1)
}

func (s *ReposServer) ListBranches(v0 context.Context, v1 *sourcegraph.ReposListBranchesOp) (*sourcegraph.BranchList, error) {
	return s.ListBranches_(v0, v1)
}

func (s *ReposServer) ListTags(v0 context.Context, v1 *sourcegraph.ReposListTagsOp) (*sourcegraph.TagList, error) {
	return s.ListTags_(v0, v1)
}

var _ sourcegraph.ReposServer = (*ReposServer)(nil)

type MirrorReposServer struct {
	RefreshVCS_ func(v0 context.Context, v1 *sourcegraph.RepoSpec) (*pbtypes.Void, error)
}

func (s *MirrorReposServer) RefreshVCS(v0 context.Context, v1 *sourcegraph.RepoSpec) (*pbtypes.Void, error) {
	return s.RefreshVCS_(v0, v1)
}

var _ sourcegraph.MirrorReposServer = (*MirrorReposServer)(nil)