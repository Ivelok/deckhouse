package synthetic

import (
	"context"
	"net/http"
	"time"

	"upmeter/pkg/checks"
	"upmeter/pkg/probes/util"
)

/*
CHECK:
It should be at least one smoke-mini Pod that response
with code 200 via "/neighbor" endpoint.

Get IPs from DNS, randomize list and sequentially
request endpoint until first success.

Period: 5 seconds.
Dns resolve timeout: 2 seconds.
Http response timeout: 4 seconds.
*/

func NewNeighborProbe() *checks.Probe {
	var nghProbeRef = checks.ProbeRef{
		Group: groupName,
		Probe: "neighbor",
	}
	const nghPeriod = 5 * time.Second
	const nghDnsTimeout = 2 * time.Second
	const nghTimeout = 4 * time.Second

	pr := &checks.Probe{
		Ref:    &nghProbeRef,
		Period: nghPeriod,
	}

	pr.RunFn = func() {
		log := pr.LogEntry()

		smokeIPs, found := LookupAndShuffleIPs(SmokeMiniAddr, nghDnsTimeout)
		if !found {
			pr.ResultCh <- pr.Result(checks.StatusUnknown)
			return
		}

		success := false

		util.SequentialDoWithTimer(
			context.Background(),
			nghTimeout,
			smokeIPs,
			func(ctx context.Context, idx int, item string) int {
				_, status, err := RequestSmokeMiniUrl(ctx, item, "/neighbor")
				if err != nil {
					log.Debugf("Request smoke mini '%s': %v", item, err)
					return 0
				}

				if status == http.StatusOK {
					success = true
					// Stop the loop
					return 1
				}
				return 0
			}, func(idx int, item string) {
				// The last smokeIp is timed out, send fail result.
				if idx == len(smokeIPs)-1 {
					pr.ResultCh <- pr.Result(checks.StatusFail)
				}
			})

		pr.ResultCh <- pr.Result(success)
	}

	return pr
}