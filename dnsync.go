package main

import (
    "fmt"
    "os"
    "net"
    "time"
    "os/signal"
    "syscall"

    "mandrakey.cc/dnsync/config"
    "mandrakey.cc/dnsync/handler"

    "github.com/urfave/cli"
    "github.com/miekg/dns"
)

var configFile string = "./dnsync.json"

func main() {
    // todo: Set up logging and replace all those printf-calls with debug/info/error output

    app := cli.NewApp()
    app.Name = "dnsync"
    app.Usage = "DNS meta synchronizer"
    app.Version = "1.0.0-alpha"
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

func actionRun(c *cli.Context) error {
    // Load config
    cfg := config.AppConfigInstance()
    err := cfg.LoadFromFile(configFile); if err != nil {
        return err
    }
    cfg.ConfigFile = configFile

    if cfg.Verbose {
        fmt.Printf("Loaded config:\n%s\n", cfg.String())
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

    fmt.Printf("Listening on %s\n", addr.String())
    buf := make([]byte, 4096)

    loop:
    for {
        conn.SetReadDeadline(time.Now().Add(time.Second))
        n, raddr, _ := conn.ReadFromUDP(buf)

        if n > 0 {
            fmt.Printf("Read %d bytes from %s:\n", n, raddr.String())
            go handlePacket(buf, raddr)
        }

        // Check if we got a signal in the meantime
        select {
        case sig := <-sigc:
            if sig == syscall.SIGINT || sig == syscall.SIGTERM {
                fmt.Println("Shutting down.")
                break loop
            }

        case <-time.After(time.Millisecond * 100):
            // pass to continue after 100ms if no signal in channel
        }
    }

    return nil
}

func handlePacket(data []byte, raddr *net.UDPAddr) {
    msg := dns.Msg{}
    err := msg.Unpack(data); if err != nil {
        fmt.Printf("Failed to unpack packet: %s\n", err)
        return
    }

    if msg.MsgHdr.Opcode != dns.OpcodeNotify || msg.Answer[0].Header().Rrtype != dns.TypeSOA {
        // invalid request, not a notify
        fmt.Printf("Skip invalid notify")
        return
    }

    soa := msg.Answer[0].(*dns.SOA)
    fmt.Printf("Received notify for %s\n", soa.Hdr.Name)

    cfg := config.AppConfigInstance()
    for _, h := range(cfg.Handlers) {
        if cfg.Verbose {
            fmt.Println("Processing message for", h.Name)
        }
        handler.HandleMessage(&h, &msg, raddr)
    }

    // Send response
    res := dns.Msg{}
    res.SetReply(&msg)
    fmt.Printf("Sending reply to %s:%d", raddr.IP, raddr.Port)

    c := dns.Client{}
    c.Exchange(&res, fmt.Sprintf("%s:%d", raddr.IP, raddr.Port))
}
