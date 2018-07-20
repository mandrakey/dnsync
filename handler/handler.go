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
    switch handler.Type {
    case HANDLER_BIND:
        return handleMessageBind(handler, msg, raddr)

    default:
        return fmt.Errorf("No such handler type: %s", handler.Type)
    }
}

func handleMessageBind(handler *config.Handler, msg *dns.Msg, raddr *net.UDPAddr) error {
    domain := strings.TrimSuffix(msg.Answer[0].Header().Name, ".")
    zone := bind.Zone{
        Name: domain,
        Masters: []string{raddr.IP.String()},
        File: fmt.Sprintf("%s/%s.host", handler.BindZonefilesPath, domain),
    }

    cfg := config.AppConfigInstance()

    bc := bind.NewBindConfig()
    bc.Load(handler.BindConfigFile)
    if cfg.Verbose {
        fmt.Printf("Current slave zones:%s\n", bc.String())
    }

    bc.AddZone(&zone)
    if cfg.Verbose {
        fmt.Printf("New slave zones:\n%s", bc.String())
    }

    bc.Save(handler.BindConfigFile)

    return nil
}
