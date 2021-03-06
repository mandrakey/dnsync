/* This file is part of DNSync.
 *
 * Copyright (C) 2018 Maurice Bleuel <mandrakey@bleuelmedia.com>
 * Licensed undert the simplified BSD license. For further details see COPYING.
 */

package main

import (
    "fmt"
    "os"
    "net"
    "time"
    "os/signal"
    "syscall"

    "github.com/mandrakey/dnsync/config"
    "github.com/mandrakey/dnsync/handler"

    "github.com/urfave/cli"
    "github.com/miekg/dns"
)

const AppVersion = "1.0.0"

var configFile string = "./dnsync.json"

func main() {
    app := cli.NewApp()
    app.Name = "dnsync"
    app.Usage = "DNS meta synchronizer"
    app.Version = AppVersion
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name: "config, c",
            Value: "./dnsync.json",
            Usage: "Load configuration from `FILE`",
            Destination: &configFile,
        },
    }
    app.Action = actionRun

    err := app.Run(os.Args)
    if err != nil {
        fmt.Printf("ERROR %s\n", err)
    }
}

// Run action that reads the config, starts the listening server and sets up signal catching.
func actionRun(c *cli.Context) error {
    // Load config
    cfg := config.AppConfigInstance()
    err := cfg.LoadFromFile(configFile); if err != nil {
        return err
    }
    cfg.ConfigFile = configFile

    // Setup logging
    config.SetupLogging(cfg.Logfile)
    log := config.Logger()
    log.Noticef("This is dnsync v.%s", AppVersion)
    fmt.Printf("This is dnsync v.%s\n", AppVersion)
    fmt.Println("Copyright (C) 2018 Maurice Bleuel")
    fmt.Println("Licensed under the MIT license.")

    if cfg.Verbose {
        log.Debugf("Loaded config:\n%s", cfg.String())
    }

    // Create UDP socket
    addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)); if err != nil {
        return err
    }
    conn, err := net.ListenUDP("udp", addr); if err != nil {
        return err
    }

    // Create signal catcher
    sigc := make(chan os.Signal, 2)
    signal.Notify(sigc, syscall.SIGINT)
    signal.Notify(sigc, syscall.SIGTERM)

    log.Infof("Listening on %s", addr.String())
    fmt.Printf("Listening on %s\n", addr.String())
    buf := make([]byte, 4096)

    loop:
    for {
        conn.SetReadDeadline(time.Now().Add(time.Second))
        n, raddr, _ := conn.ReadFromUDP(buf)

        if n > 0 {
            log.Debugf("Read %d bytes from %s", n, raddr.String())
            go handlePacket(buf, raddr)
        }

        // Check if we got a signal in the meantime
        select {
        case sig := <-sigc:
            if sig == syscall.SIGINT || sig == syscall.SIGTERM {
                log.Infof("Shutting down.")
                break loop
            }

        case <-time.After(time.Millisecond * 100):
            // pass to continue after 100ms if no signal in channel
        }
    }

    return nil
}

// Method to handle incoming DNS packets. Only packets with opcode NOTIFY and type SOA will be handled, everything
// else will be discarded. If a valid packet is found, it is sent to every registered handler to work with it. After
// all handlers have finished processing, a DNS reply packet will be sent to the client before the connection is closed.
func handlePacket(data []byte, raddr *net.UDPAddr) {
    log := config.Logger()

    if !validRemote(string(raddr.IP)) {
        log.Infof("Discard packet from invalid remote address %s", raddr.IP)
        return
    }

    msg := dns.Msg{}
    err := msg.Unpack(data); if err != nil {
        log.Errorf("Failed to unpack packet: %s", err)
        return
    }

    if msg.MsgHdr.Opcode != dns.OpcodeNotify || msg.Answer[0].Header().Rrtype != dns.TypeSOA {
        // invalid request, not a notify
        log.Info("Skip invalid notify")
        return
    }

    soa := msg.Answer[0].(*dns.SOA)
    log.Infof("Received notify for %s", soa.Hdr.Name)

    cfg := config.AppConfigInstance()
    for _, h := range(cfg.Handlers) {
        if cfg.Verbose {
            log.Debugf("Processing message for %s", h.Name)
        }
        err = handler.HandleMessage(&h, &msg, raddr); if err != nil {
            log.Error(err)
        }
    }

    // Send response
    res := dns.Msg{}
    res.SetReply(&msg)
    log.Debugf("Sending reply to %s:%d", raddr.IP, raddr.Port)

    c := dns.Client{}
    c.Exchange(&res, fmt.Sprintf("%s:%d", raddr.IP, raddr.Port))
}

// Checks whether or not a given ip address is in the list of configured remotes.
func validRemote(ip string) bool {
    cfg := config.AppConfigInstance()
    for _, r := range cfg.Remotes {
        if r == ip {
            return true
        }
    }
    return false
}
