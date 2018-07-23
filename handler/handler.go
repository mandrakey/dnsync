package handler

import (
    "fmt"
    "net"
    "strings"

    "github.com/miekg/dns"

    "mandrakey.cc/dnsync/bind"
    "mandrakey.cc/dnsync/config"
)

const (
    HANDLER_BIND = "bind"
)

func HandleMessage(handler *config.Handler, msg *dns.Msg, raddr *net.UDPAddr) error {
    log := config.Logger()

    switch handler.Type {
    case HANDLER_BIND:
        log.Debug("Handling BIND message")
        return handleMessageBind(handler, msg, raddr)

    default:
        return fmt.Errorf("No such handler type: %s", handler.Type)
    }
}

func handleMessageBind(handler *config.Handler, msg *dns.Msg, raddr *net.UDPAddr) error {
    log := config.Logger()

    domain := strings.TrimSuffix(msg.Answer[0].Header().Name, ".")
    zone := bind.Zone{
        Name: domain,
        Masters: []string{raddr.IP.String()},
        File: fmt.Sprintf("%s/%s.host", handler.BindZonefilesPath, domain),
    }
    log.Debugf("Handling BIND message for '%s':\n%s", domain, zone.String())

    bc := bind.NewBindConfig()
    bc.Load(handler.BindConfigFile)
    log.Debugf("Current slave zones: %s", bc.String())

    bc.AddZone(&zone)
    log.Debugf("New slave zones: %s", bc.String())

    bc.Save(handler.BindConfigFile)
    return nil
}
