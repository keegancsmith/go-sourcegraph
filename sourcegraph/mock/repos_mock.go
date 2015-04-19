// generated by gen-mocks; DO NOT EDIT

package mock

import "sourcegraph.com/sourcegraph/go-sourcegraph/sourcegraph"

type ReposService struct {
	Get_       func(repo sourcegraph.RepoSpec, opt *sourcegraph.RepoGetOptions) (*sourcegraph.Repo, error)
	List_      func(opt *sourcegraph.RepoListOptions) ([]*sourcegraph.Repo, error)
	Create_    func(newRepo *sourcegraph.Repo) (*sourcegraph.Repo, error)
	GetReadme_ func(repo sourcegraph.RepoRevSpec) (*sourcegraph.Readme, error)
}

func (s *ReposService) Get(repo sourcegraph.RepoSpec, opt *sourcegraph.RepoGetOptions) (*sourcegraph.Repo, error) {
	return s.Get_(repo, opt)
}

func (s *ReposService) List(opt *sourcegraph.RepoListOptions) ([]*sourcegraph.Repo, error) {
	return s.List_(opt)
}

func (s *ReposService) Create(newRepo *sourcegraph.Repo) (*sourcegraph.Repo, error) {
	return s.Create_(newRepo)
}

func (s *ReposService) GetReadme(repo sourcegraph.RepoRevSpec) (*sourcegraph.Readme, error) {
	return s.GetReadme_(repo)
}

var _ sourcegraph.ReposService = (*ReposService)(nil)