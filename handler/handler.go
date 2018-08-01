/* This file is part of DNSync.
 *
 * Copyright (C) 2018 Maurice Bleuel <mandrakey@bleuelmedia.com>
 * Licensed undert the simplified BSD license. For further details see COPYING.
 */

package handler

import (
    "fmt"
    "net"
    "strings"

    "github.com/miekg/dns"

    "github.com/mandrakey/dnsync/bind"
    "github.com/mandrakey/dnsync/config"
)

const (
    HANDLER_BIND = "bind"
)

// Takes a handler configuration, a DNS message packet and a UDP address struct denoting the client. The strategy
// for handling the packet will be determined using the Handler.Type field. Currently, only BIND is supported.
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

// Handles a DNS NOTIFY packet for a bind nameserver: The zone will be constructed and, if necessary, added to
// the bind dnsync configuration file.
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
