package sourcegraph

import (
	"strconv"
	"strings"

	"sourcegraph.com/sourcegraph/go-sourcegraph/router"
	"sourcegraph.com/sourcegraph/srclib/person"
)

// OrgsService communicates with the organizations-related endpoints in the
// Sourcegraph API.
type OrgsService interface {
	// Get fetches an organization.
	Get(org OrgSpec) (*Org, Response, error)

	// ListMembers lists members of an organization.
	ListMembers(org OrgSpec, opt *OrgListMembersOptions) ([]*Person, Response, error)

	// GetSettings fetches an org's configuration settings.
	GetSettings(org OrgSpec) (*OrgSettings, Response, error)

	// UpdateSettings updates an org's configuration settings.
	UpdateSettings(org OrgSpec, settings OrgSettings) (Response, error)
}

// orgsService implements OrgsService.
type orgsService struct {
	client *Client
}

var _ OrgsService = &orgsService{}

// OrgSpec specifies an organization. At least one of Email, Login, and UID must be
// nonempty.
type OrgSpec struct {
	Org string // name of organization
	UID int    // user ID of the "user" record for this organization
}

// PathComponent returns the URL path component that specifies the org.
func (s *OrgSpec) PathComponent() string {
	if s.Org != "" {
		return s.Org
	}
	if s.UID > 0 {
		return "$" + strconv.Itoa(s.UID)
	}
	panic("empty OrgSpec")
}

func (s *OrgSpec) RouteVars() map[string]string {
	return map[string]string{"OrgSpec": s.PathComponent()}
}

type Org struct {
	person.User
}

// ParseOrgSpec parses a string generated by (*OrgSpec).String() and
// returns the equivalent OrgSpec struct.
func ParseOrgSpec(pathComponent string) (OrgSpec, error) {
	if strings.HasPrefix(pathComponent, "$") {
		uid, err := strconv.Atoi(pathComponent[1:])
		return OrgSpec{UID: uid}, err
	}
	return OrgSpec{Org: pathComponent}, nil
}

func (s *orgsService) Get(org OrgSpec) (*Org, Response, error) {
	url, err := s.client.url(router.Org, org.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var org_ *Org
	resp, err := s.client.Do(req, &org_)
	if err != nil {
		return nil, resp, err
	}

	return org_, resp, nil
}

type OrgListMembersOptions struct {
	ListOptions
}

func (s *orgsService) ListMembers(org OrgSpec, opt *OrgListMembersOptions) ([]*Person, Response, error) {
	url, err := s.client.url(router.OrgMembers, org.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var members []*Person
	resp, err := s.client.Do(req, &members)
	if err != nil {
		return nil, resp, err
	}

	return members, resp, nil
}

// OrgSettings describes an org's configuration settings.
type OrgSettings struct {
	Plan *PlanSettings `json:",omitempty"`
}

func (s *orgsService) GetSettings(org OrgSpec) (*OrgSettings, Response, error) {
	url, err := s.client.url(router.OrgSettings, org.RouteVars(), nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var settings *OrgSettings
	resp, err := s.client.Do(req, &settings)
	if err != nil {
		return nil, resp, err
	}

	return settings, resp, nil
}

func (s *orgsService) UpdateSettings(org OrgSpec, settings OrgSettings) (Response, error) {
	url, err := s.client.url(router.OrgSettingsUpdate, org.RouteVars(), nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("PUT", url.String(), settings)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

type MockOrgsService struct {
	Get_            func(org OrgSpec) (*Org, Response, error)
	ListMembers_    func(org OrgSpec, opt *OrgListMembersOptions) ([]*Person, Response, error)
	GetSettings_    func(org OrgSpec) (*OrgSettings, Response, error)
	UpdateSettings_ func(org OrgSpec, settings OrgSettings) (Response, error)
}

var _ OrgsService = MockOrgsService{}

func (s MockOrgsService) Get(org OrgSpec) (*Org, Response, error) {
	if s.Get_ == nil {
		return nil, &HTTPResponse{}, nil
	}
	return s.Get_(org)
}

func (s MockOrgsService) ListMembers(org OrgSpec, opt *OrgListMembersOptions) ([]*Person, Response, error) {
	if s.ListMembers_ == nil {
		return nil, nil, nil
	}
	return s.ListMembers_(org, opt)
}

func (s MockOrgsService) GetSettings(org OrgSpec) (*OrgSettings, Response, error) {
	if s.GetSettings_ == nil {
		return nil, nil, nil
	}
	return s.GetSettings_(org)
}

func (s MockOrgsService) UpdateSettings(org OrgSpec, settings OrgSettings) (Response, error) {
	if s.UpdateSettings_ == nil {
		return nil, nil
	}
	return s.UpdateSettings_(org, settings)
}