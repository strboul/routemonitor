package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"strings"

	"github.com/google/gopacket/routing"
	netroute "github.com/libp2p/go-netroute"
)

type Router struct {
	config Config
	router routing.Router
	logger *slog.Logger
}

type CheckErrors []error

func parseIp(ip string) (net.IP, error) {
	parsedIp := net.ParseIP(ip).To4()
	if parsedIp == nil {
		return nil, fmt.Errorf("cannot parse IP")
	}
	return parsedIp, nil
}

func getCidrIps(ipCidr string) ([]netip.Addr, error) {
	var ips []netip.Addr
	p, err := netip.ParsePrefix(ipCidr)
	if err != nil {
		return ips, err
	}
	p = p.Masked()
	addr := p.Addr()
	for {
		if !p.Contains(addr) {
			break
		}
		ips = append(ips, addr)
		addr = addr.Next()
	}
	return ips, nil
}

func checkDevice(device string, expectDevice string) error {
	if device != expectDevice {
		// check if device interface exists
		_, err := net.InterfaceByName(expectDevice)
		if err != nil {
			// interface doesn't exist
			if errStr := err.Error(); strings.Contains(errStr, "no such network interface") {
				return fmt.Errorf("interface not exist expect=\"%s\"", expectDevice)
			}
		}
		return fmt.Errorf("mismatch device current=\"%s\" expect=\"%s\"", device, expectDevice)
	}
	return nil
}

func checkGateway(gateway string, expectGateway string) error {
	if gateway != expectGateway {
		return fmt.Errorf("mismatch gateway current=\"%s\" expect=\"%s\"", gateway, expectGateway)
	}
	return nil
}

func checkSource(source string, expectSource string) error {
	if source != expectSource {
		return fmt.Errorf("mismatch source current=\"%s\" expect=\"%s\"", source, expectSource)
	}
	return nil
}

func (r *Router) doExpect(iface *net.Interface, gateway net.IP, src net.IP, expects []RouteExpect) error {
	var checkedDevices []string

	numExpects := len(expects)

	for i, expect := range expects {
		if expect.When.Device != "" {
			err := checkDevice(iface.Name, expect.When.Device)
			checkedDevices = append(checkedDevices, expect.When.Device)
			if err != nil {
				if strings.Contains(err.Error(), "interface not exist") {
					// if this isn't the last, go to the next
					if i < numExpects-1 {
						continue
					}
				}
				if numExpects > 1 {
					return fmt.Errorf("checked all devices, not matching any current expect=\"%v\"", checkedDevices)
				}
				return err
			}
		}

		if expect.When.Gateway != "" {
			err := checkGateway(gateway.To4().String(), expect.When.Gateway)
			if err != nil {
				return err
			}
		}

		if expect.When.Source != "" {
			err := checkSource(src.To4().String(), expect.When.Source)
			if err != nil {
				return err
			}
		}

		break
	}

	return nil
}

func (r *Router) checkIpRoute(ipv4 net.IP, expects []RouteExpect) error {
	iface, gateway, src, err := r.router.Route(ipv4)
	if err != nil {
		return err
	}

	err = r.doExpect(iface, gateway, src, expects)
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) loopRoutes() (error, CheckErrors) {
	var checkErrs CheckErrors

	for _, route := range r.config.Route {

		ipAddr := strings.Split(route.IP, "/")

		if len(ipAddr) == 1 {
			route.IP = ipAddr[0] + "/32"
		}

		ips, err := getCidrIps(route.IP)
		if err != nil {
			return err, nil
		}

		r.logger.Info(
			"checking route",
			slog.String("name", route.Name),
			slog.String("ip", route.IP),
			slog.Int("num IPs", len(ips)),
			slog.Any("expects", route.Expect),
		)

		for _, ip := range ips {
			ipv4, err := parseIp(ip.String())
			if err != nil {
				return err, nil
			}

			err = r.checkIpRoute(ipv4, route.Expect)
			if err != nil {
				errMsg := fmt.Errorf(`name="%s" ip="%s" %s`, route.Name, route.IP, err)
				checkErrs = append(checkErrs, errMsg)
				if r.config.FailFast {
					return nil, checkErrs
				}
			}
		}
	}

	if len(checkErrs) > 0 {
		return nil, checkErrs
	}

	return nil, nil
}

func newRouter(config Config, logger *slog.Logger) (*Router, error) {
	router, err := netroute.New()
	if err != nil {
		return nil, err
	}
	return &Router{
		config: config,
		router: router,
		logger: logger,
	}, nil
}

func CheckRoutes(config Config, logger *slog.Logger) error {
	router, err := newRouter(config, logger)
	if err != nil {
		return err
	}
	err, checkErrs := router.loopRoutes()
	if err != nil {
		return err
	}
	if checkErrs != nil {
		for _, err := range checkErrs {
			logger.Error(err.Error())
		}
		return errors.New("There are mismatch in routes")
	}

	fmt.Println("All routes are as expected.")

	return nil
}
