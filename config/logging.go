package config

import (
    "os"

    "github.com/op/go-logging"
)

var (
    logger = logging.MustGetLogger("dnsync")
    logFormat = logging.MustStringFormatter(`[%{time:2006-01-02 15:04:05}] %{level} %{message}`)
)

func SetupLogging(logfile string) {
    var backend logging.Backend
    fp, fperr := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if fp != nil {
        backend = logging.NewLogBackend(fp, "", 0)
    } else {
        backend = logging.NewLogBackend(os.Stdout, "", 0)
    }

    realBackend := logging.AddModuleLevel(logging.NewBackendFormatter(backend, logFormat))
    realBackend.SetLevel(logging.DEBUG, "")
    logging.SetBackend(realBackend)

    if fperr != nil {
        logger.Warningf("Failed to setup logging to file. Falling back to stdout.\n%s", fperr)
    }
}

func Logger() *logging.Logger {
    return logger
}
