package presenters

import (
	"bytes"
	"fmt"
)

func (p *Presenter) asString() {
	var response bytes.Buffer

	const (
		orgOverviewMsg               = "org %s is consuming %d MB of %d MB\n"
		spaceOverviewMsg             = "\tspace %s is consuming %d MB memory (%d%%) of org quota\n"
		spaceBillableAppInstancesMsg = "\t\tAIs billable: %d\n"
		spaceAppInstancesMsg         = "\t\tAIs canonical: %d (%d running, %d stopped)\n"
		spaceSCSMsg                  = "\t\tSCS instances: %d\n"
		reportSummaryMsg             = "across %d org(s), you have %d billable AIs, %d are canonical AIs (%d running, %d stopped), %d are SCS instances\n"
	)

	for _, OrgReport := range p.Report.OrgReports {
		response.WriteString(fmt.Sprintf(orgOverviewMsg, OrgReport.Name, OrgReport.MemoryUsage, OrgReport.MemoryQuota))
		for _, SpaceReport := range OrgReport.SpaceReport {
			if OrgReport.MemoryQuota > 0 {
				spaceMemoryConsumedPercentage := (100 * SpaceReport.ConsumedMemory / OrgReport.MemoryQuota)
				response.WriteString(fmt.Sprintf(spaceOverviewMsg, SpaceReport.Name, SpaceReport.ConsumedMemory, spaceMemoryConsumedPercentage))
			}
			response.WriteString(fmt.Sprintf(spaceBillableAppInstancesMsg, SpaceReport.BillableAppInstancesCount()))
			response.WriteString(fmt.Sprintf(spaceAppInstancesMsg, SpaceReport.AppInstancesCount, SpaceReport.RunningAppInstancesCount, SpaceReport.StoppedAppInstancesCount))
			response.WriteString(fmt.Sprintf(spaceSCSMsg, SpaceReport.SpringCloudServicesCount()))
		}
	}

	response.WriteString(
		fmt.Sprintf(
			reportSummaryMsg,
			len(p.Report.Orgs),
			p.Report.AggregateOrgReport.BillableAppInstancesCount,
			p.Report.AggregateOrgReport.AppInstancesCount,
			p.Report.AggregateOrgReport.RunningAppInstancesCount,
			p.Report.AggregateOrgReport.StoppedAppInstancesCount,
			p.Report.AggregateOrgReport.SpringCloudServicesCount,
		),
	)

	fmt.Println(response.String())
}
