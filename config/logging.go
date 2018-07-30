/* This file is part of DNSync.
 *
 * Copyright (C) 2018 Maurice Bleuel <mandrakey@bleuelmedia.com>
 * Licensed undert the simplified BSD license. For further details see COPYING.
 */

package config

import (
    "os"

    "github.com/op/go-logging"
)

var (
    logger = logging.MustGetLogger("dnsync")
    logFormat = logging.MustStringFormatter(`[%{time:2006-01-02 15:04:05}] %{level} %{message}`)
)

// Set up the logging subsystem and have it log to the provided logfile.
func SetupLogging(logfile string) {
    var backend logging.Backend
    fp, fperr := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if fp != nil {
        backend = logging.NewLogBackend(fp, "", 0)
    } else {
        backend = logging.NewLogBackend(os.Stdout, "", 0)
    }

    cfg := AppConfigInstance()
    realBackend := logging.AddModuleLevel(logging.NewBackendFormatter(backend, logFormat))
    realBackend.SetLevel(stringToLoglevel(cfg.Loglevel), "")
    logging.SetBackend(realBackend)

    if fperr != nil {
        logger.Warningf("Failed to setup logging to file. Falling back to stdout.\n%s", fperr)
    }
}

// Get the global Logger instance.
func Logger() *logging.Logger {
    return logger
}

// Parse a given string s into a useable log level value.
func stringToLoglevel(s string) logging.Level {
    l, err := logging.LogLevel(s)
    if err == nil {
        return l
    } else {
        return logging.INFO
    }
}
