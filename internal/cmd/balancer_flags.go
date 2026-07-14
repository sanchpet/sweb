package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sanchpet/sweb-go-sdk/balancer"
	"github.com/sanchpet/sweb-go-sdk/flex"
)

// balancerConfig aliases balancer.Config so `balancer config` can embed the
// catalog next to the createEnabled flag on the json output path.
type balancerConfig = balancer.Config

// parseServers turns repeatable --server values into balancer.Server entries.
// Format: "ip[,weight]" — weight (1..5) applies only to the roundrobin type and
// defaults to 0 (absent) when omitted.
func parseServers(vals []string) ([]balancer.Server, error) {
	servers := make([]balancer.Server, 0, len(vals))
	for _, v := range vals {
		fields := strings.Split(v, ",")
		s := balancer.Server{IP: strings.TrimSpace(fields[0])}
		if s.IP == "" {
			return nil, fmt.Errorf("--server %q: missing IP", v)
		}
		if len(fields) > 1 && strings.TrimSpace(fields[1]) != "" {
			weight, err := strconv.Atoi(strings.TrimSpace(fields[1]))
			if err != nil {
				return nil, fmt.Errorf("--server %q: bad weight: %w", v, err)
			}
			s.Weight = flex.Int(weight)
		}
		if len(fields) > 2 {
			s.VPSName = strings.TrimSpace(fields[2])
		}
		servers = append(servers, s)
	}
	return servers, nil
}

// parseRules turns repeatable --rule values into balancer.Rule entries.
// Format: "protoBalancer:portBalancer:protoServer:portServer".
func parseRules(vals []string) ([]balancer.Rule, error) {
	rules := make([]balancer.Rule, 0, len(vals))
	for _, v := range vals {
		fields := strings.Split(v, ":")
		if len(fields) != 4 {
			return nil, fmt.Errorf("--rule %q: want protoBalancer:portBalancer:protoServer:portServer", v)
		}
		rules = append(rules, balancer.Rule{
			ProtocolBalancer: strings.TrimSpace(fields[0]),
			PortBalancer:     strings.TrimSpace(fields[1]),
			ProtocolServer:   strings.TrimSpace(fields[2]),
			PortServer:       strings.TrimSpace(fields[3]),
		})
	}
	return rules, nil
}
