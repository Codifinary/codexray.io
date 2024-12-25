package overview

import (
	"fmt"
	"sort"

	"codexray/model"
	"codexray/timeseries"
	"codexray/utils"

	"github.com/dustin/go-humanize/english"
)

func renderHealth(w *model.World) []*ApplicationStatus {
	var res []*ApplicationStatus
	for _, app := range w.Applications {
		switch {
		case app.IsK8s():
		case app.Id.Kind == model.ApplicationKindNomadJobGroup:
		case !app.IsStandalone():
		default:
			continue
		}

		a := &ApplicationStatus{
			Id:       app.Id,
			Category: app.Category,
			Type:     getApplicationType(app),
		}

		sloIsViolating := false
		for _, r := range app.Reports {
			for _, ch := range r.Checks {
				switch ch.Id {
				case model.Checks.SLOAvailability.Id:
					for _, sli := range app.AvailabilitySLIs {
						if ch.Status >= model.WARNING {
							a.Errors.Status = model.CRITICAL
							sloIsViolating = true
						}
						failed := sli.FailedRequests.Reduce(timeseries.NanSum)
						total := sli.TotalRequests.Reduce(timeseries.NanSum)
						if total > 0 && failed > 0 {
							a.Errors.Value = formatPercent(failed * 100 / total)
						}
						break
					}
				case model.Checks.SLOLatency.Id:
					for _, sli := range app.LatencySLIs {
						if ch.Status >= model.WARNING {
							a.Latency.Status = model.CRITICAL
							sloIsViolating = true
						}
						latency := quantile(sli.Histogram, sli.Config.ObjectivePercentage/100)
						if latency > 0 {
							a.Latency.Value = utils.FormatLatency(latency)
						}
						break
					}
				case model.Checks.InstanceAvailability.Id:
					a.Instances.Status = ch.Status
					if ch.Desired() > 0 {
						a.Instances.Value = fmt.Sprintf("%d/%d", ch.Count(), ch.Desired())
					}
				case model.Checks.InstanceRestarts.Id:
					a.Restarts.Status = ch.Status
					if ch.Count() > 0 {
						a.Restarts.Value = fmt.Sprintf("%d", ch.Count())
					}
				case model.Checks.CPUNode.Id:
					if ch.Status >= model.WARNING && sloIsViolating {
						a.CPU.Status = model.WARNING
						a.CPU.Value = "shortage"
					}
				case model.Checks.CPUContainer.Id:
					if ch.Status >= model.WARNING {
						a.CPU.Status = model.WARNING
						a.CPU.Value = "shortage"
					}
				case model.Checks.MemoryOOM.Id:
					if ch.Status >= model.WARNING {
						a.Memory.Status = model.WARNING
						a.Memory.Value = "OOM"
					}
				case model.Checks.MemoryLeakPercent.Id:
					if ch.Status >= model.WARNING && a.Memory.Status < model.WARNING {
						a.Memory.Status = model.WARNING
						a.Memory.Value = "leak"
					}
				case model.Checks.StorageIOLoad.Id:
					if ch.Status != model.UNKNOWN {
						a.DiskIOLoad.Status = ch.Status
					}
					if ch.Value() > 0 {
						a.DiskIOLoad.Value = utils.FormatFloat(ch.Value())
					}
				case model.Checks.StorageSpace.Id:
					a.DiskUsage.Status = ch.Status
					if ch.Value() > 0 {
						a.DiskUsage.Value = formatPercent(ch.Value())
					}
				case model.Checks.NetworkRTT.Id:
					if ch.Status != model.UNKNOWN {
						a.Network.Status = ch.Status
						if !sloIsViolating {
							a.Network.Status = model.OK
						}
					}
					if ch.Value() > 0 {
						a.Network.Value = utils.FormatLatency(ch.Value())
					}
				case model.Checks.NetworkConnectivity.Id:
					if ch.Status >= model.WARNING {
						a.Network.Status = ch.Status
						a.Network.Value = "packet loss"
					}
				case model.Checks.NetworkTCPConnections.Id:
					if ch.Status >= model.WARNING {
						a.Network.Status = ch.Status
						a.Network.Value = "failed conns"
					}
				case model.Checks.DnsLatency.Id:
					if ch.Status >= model.WARNING {
						a.DNS.Status = ch.Status
						if ch.Value() > 0 {
							a.DNS.Value = utils.FormatLatency(ch.Value())
						}
					}
				case model.Checks.DnsServerErrors.Id, model.Checks.DnsNxdomainErrors.Id:
					if ch.Status >= model.WARNING {
						a.DNS.Status = ch.Status
						a.DNS.Value = "errors"
					}
				case model.Checks.LogErrors.Id:
					if items := ch.Items(); items != nil && items.Len() > 0 {
						count := items.Len()
						a.Logs.Value = fmt.Sprintf("%d unique %s", count, english.PluralWord(count, "error", ""))
						a.Logs.Chart = ch.Values()
					}
					if ch.Status >= model.WARNING {
						a.Logs.Status = model.INFO
					}
				}
			}
		}

		upstreams := map[model.ApplicationId]bool{}
		for _, i := range app.Instances {
			for _, u := range i.Upstreams {
				upstream := u.RemoteApplication
				if upstream == nil || u.IsObsolete() {
					continue
				}
				if _, seen := upstreams[u.RemoteApplication.Id]; seen {
					continue
				}
				if app.Id == upstream.Id {
					continue
				}
				if !app.Category.Monitoring() && upstream.Category.Monitoring() {
					continue
				}
				for _, r := range upstream.Reports {
					if r.Name != model.AuditReportSLO {
						continue
					}
					for _, ch := range r.Checks {
						if ch.Status == model.UNKNOWN {
							continue
						}
						upstreams[upstream.Id] = upstreams[upstream.Id] || ch.Status >= model.WARNING
					}
				}
			}
		}
		if total := len(upstreams); total > 0 {
			var ok int
			for _, failed := range upstreams {
				if !failed {
					ok++
				}
			}
			if ok < total && sloIsViolating {
				a.Upstreams.Status = model.WARNING
			} else {
				a.Upstreams.Status = model.OK
			}
			a.Upstreams.Value = fmt.Sprintf("%d/%d", ok, total)
		}

		calcApplicationStatus(a)

		if a.Status == model.UNKNOWN {
			continue
		}

		if t := a.Type; t != nil {
			if t.Status > a.Status {
				a.Status = t.Status
			}
			if a.Status == model.OK && t.Status == model.UNKNOWN && t.Report != "" {
				a.Status = model.UNKNOWN
			}
		}

		res = append(res, a)
	}

	sort.Slice(res, func(i, j int) bool {
		if res[i].Status == res[j].Status {
			return res[i].Id.Name < res[j].Id.Name
		}
		if res[i].Status == model.OK {
			return false
		}
		if res[j].Status == model.OK {
			return true
		}
		return res[i].Status > res[j].Status
	})

	return res
}

func init() {
	calcApplicationStatus(&ApplicationStatus{}) // check for panics
}
