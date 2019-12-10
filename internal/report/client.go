package report

import (
	"github.com/aegershman/cf-report-usage-plugin/internal/v2client"
	"github.com/cloudfoundry/cli/plugin"
)

// Client orchestrates generation and aggregation of report data
type Client struct {
	client *v2client.Client
}

// NewClient -
func NewClient(cli plugin.CliConnection) (*Client, error) {
	client, err := v2client.NewClient(cli)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		client: client,
	}, nil
}

// GetSummaryReportByOrgNames -
func (r *Client) GetSummaryReportByOrgNames(orgNames ...string) (*SummaryReport, error) {
	populatedOrgs, err := r.getOrgs(orgNames...)
	if err != nil {
		return &SummaryReport{}, nil
	}

	var orgReports []OrgReport
	for _, org := range populatedOrgs {
		spaceReports := r.getSpaceReportsByOrg(org)
		orgQuota, _ := r.client.OrgQuotas.GetOrgQuota(org.QuotaURL)
		orgReport := *NewOrgReport(orgQuota, org, spaceReports...)
		orgReports = append(orgReports, orgReport)
	}

	return NewSummaryReport(orgReports...), nil
}

func (r *Client) getSpaceReportsByOrg(org v2client.Org) []SpaceReport {
	var spaceReports []SpaceReport
	for _, space := range org.Spaces {
		spaceReport := *NewSpaceReport(space)
		spaceReports = append(spaceReports, spaceReport)
	}
	return spaceReports
}

func (r *Client) getOrgs(orgNames ...string) ([]v2client.Org, error) {
	var rawOrgs []v2client.Org

	if len(orgNames) > 0 {
		for _, orgName := range orgNames {
			rawOrg, err := r.client.Orgs.GetOrg(orgName)
			if err != nil {
				return nil, err
			}
			rawOrgs = append(rawOrgs, rawOrg)
		}
	} else {
		extraRawOrgs, err := r.client.Orgs.GetOrgs()
		if err != nil {
			return nil, err
		}
		rawOrgs = extraRawOrgs
	}

	var orgs = []v2client.Org{}

	for _, o := range rawOrgs {
		orgDetails, err := r.getOrgDetails(o)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, orgDetails)
	}
	return orgs, nil
}

func (r *Client) getOrgDetails(o v2client.Org) (v2client.Org, error) {
	usage, err := r.client.Orgs.GetOrgMemoryUsage(o)
	if err != nil {
		return v2client.Org{}, err
	}

	quota, err := r.client.OrgQuotas.GetOrgQuota(o.QuotaURL)
	if err != nil {
		return v2client.Org{}, err
	}
	spaces, err := r.getSpaces(o.SpacesURL)
	if err != nil {
		return v2client.Org{}, err
	}

	return v2client.Org{
		GUID:        o.GUID,
		MemoryQuota: quota.MemoryLimit,
		MemoryUsage: int(usage),
		Name:        o.Name,
		QuotaURL:    o.QuotaURL,
		Spaces:      spaces,
		SpacesURL:   o.SpacesURL,
		URL:         o.URL,
	}, nil
}

func (r *Client) getSpaces(spaceURL string) ([]v2client.Space, error) {
	rawSpaces, err := r.client.Orgs.GetOrgSpaces(spaceURL)
	if err != nil {
		return nil, err
	}
	var spaces = []v2client.Space{}
	for _, s := range rawSpaces {
		apps, services, err := r.getAppsAndServices(s.GUID)
		if err != nil {
			return nil, err
		}
		spaces = append(spaces,
			v2client.Space{
				Apps:     apps,
				GUID:     s.GUID,
				Name:     s.Name,
				Services: services,
			},
		)
	}
	return spaces, nil
}

func (r *Client) getAppsAndServices(spaceGUID string) ([]v2client.App, []v2client.Service, error) {
	apps, services, err := r.client.Spaces.GetSpaceAppsAndServicesByGUID(spaceGUID)
	if err != nil {
		return nil, nil, err
	}
	return apps, services, nil
}
