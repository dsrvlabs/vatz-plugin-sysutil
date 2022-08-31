package main

import (
    "fmt"
    "github.com/rs/zerolog/log"
    "flag"
    "time"

    pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
    "github.com/dsrvlabs/vatz/sdk"
    "golang.org/x/net/context"
    "google.golang.org/protobuf/types/known/structpb"
    "github.com/shirou/gopsutil/v3/cpu"
)

const (
    defaultAddr = "127.0.0.1"
    defaultPort = 9094
    pluginName = "vatz-plugin-solana-cpu-monitor"
    defaultUrgent = 95
    defaultWarning = 90
    defaultDuration = 0
)

var (
    urgent int
    warning int
    addr string
    port int
)

func init() {
    flag.StringVar(&addr, "addr", defaultAddr, "Listening address")
    flag.IntVar(&port, "port", defaultPort, "Listening port")
    flag.IntVar(&urgent, "urgent", defaultUrgent, "CPU Usage Alert threshold")
    flag.IntVar(&warning, "warning", defaultWarning, "CPU Usage Warning threshold")

    flag.Parse()
}

func main() {
    p := sdk.NewPlugin(pluginName)
    p.Register(pluginFeature)

    ctx := context.Background()
    if err := p.Start(ctx, addr, port); err != nil {
        log.Info().Str("module", "plugin").Msg("exit")
    }
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
    // TODO: Fill here.
    ret := sdk.CallResponse{
        FuncName:   "getCPUUsage",
        Message:    "CPU usage warning!",
        Severity:   pluginpb.SEVERITY_UNKNOWN,
        State:      pluginpb.STATE_NONE,
        AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
    }

    util, _ := cpu.Percent(defaultDuration * time.Second, false)
    log.Debug().Str("module", "plugin").Int("CPU Usage", int(util[0])).Int("Urgent", urgent).Int("Warning", warning).Msg("cpu_monitor")

    if int(util[0]) > urgent {
        var message string
        message = fmt.Sprint("Current CPU Usage is ", int(util[0]), "%, over urgent threshold ", urgent, "%")
        ret = sdk.CallResponse{
            Message:    message,
            Severity:	pluginpb.SEVERITY_CRITICAL,
        }

        log.Warn().Str("module", "plugin").Msg(message)
    }

    return ret, nil
}
