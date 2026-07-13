package lab

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type Connector interface {
	Probe(context.Context) map[string]ServiceStatus
}

type LoopbackConnector struct {
	Targets map[string]string
	Timeout time.Duration
}

func IsLoopbackTarget(target string) bool {
	host, _, err := net.SplitHostPort(target)
	if err != nil {
		return false
	}
	ip := net.ParseIP(strings.Trim(host, "[]"))
	return host == "localhost" || (ip != nil && ip.IsLoopback())
}

func (c LoopbackConnector) Probe(ctx context.Context) map[string]ServiceStatus {
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = 350 * time.Millisecond
	}
	statuses := make(map[string]ServiceStatus, len(c.Targets))
	for name, target := range c.Targets {
		status := ServiceStatus{Target: target}
		if !IsLoopbackTarget(target) {
			status.Message = "target rejected: connected mode only accepts loopback"
			statuses[name] = status
			continue
		}
		dialer := net.Dialer{Timeout: timeout}
		connection, err := dialer.DialContext(ctx, "tcp", target)
		if err != nil {
			status.Message = "local service is not reachable"
		} else {
			status.Reachable = true
			status.Message = "local service accepted TCP connection"
			_ = connection.Close()
		}
		statuses[name] = status
	}
	return statuses
}

func ParseTargets(raw string) (map[string]string, error) {
	targets := map[string]string{}
	if strings.TrimSpace(raw) == "" {
		return targets, nil
	}
	for _, item := range strings.Split(raw, ",") {
		parts := strings.SplitN(strings.TrimSpace(item), "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid target %q; expected service=host:port", item)
		}
		targets[parts[0]] = parts[1]
	}
	return targets, nil
}
