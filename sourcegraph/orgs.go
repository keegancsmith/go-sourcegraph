package sourcegraph

import (
	"strconv"
	"strings"

	"sourcegraph.com/sourcegraph/go-sourcegraph/router"
)

// OrgsService communicates with the organizations-related endpoints in the
// Sourcegraph API.
type OrgsService interface {
	// Get fetches an organization.
	Get(org OrgSpec) (*Org, Response, error)

	// ListMembers lists members of an organization.
	ListMembers(org OrgSpec, opt *OrgListMembersOptions) ([]*User, Response, error)

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
	User
}

// OrgSpec returns the OrgSpec that specifies o.
func (o *Org) OrgSpec() OrgSpec { return OrgSpec{Org: o.Login, UID: int(o.UID)} }

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
	url, err := s.client.URL(router.Org, org.RouteVars(), nil)
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

func (s *orgsService) ListMembers(org OrgSpec, opt *OrgListMembersOptions) ([]*User, Response, error) {
	url, err := s.client.URL(router.OrgMembers, org.RouteVars(), opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	var members []*User
	resp, err := s.client.Do(req, &members)
	if err != nil {
		return nil, resp, err
	}

	return members, resp, nil
}

// OrgSettings describes an org's configuration settings.
type OrgSettings struct {
	PlanSettings `json:",omitempty"`
}

func (s *orgsService) GetSettings(org OrgSpec) (*OrgSettings, Response, error) {
	url, err := s.client.URL(router.OrgSettings, org.RouteVars(), nil)
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
	url, err := s.client.URL(router.OrgSettingsUpdate, org.RouteVars(), nil)
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


