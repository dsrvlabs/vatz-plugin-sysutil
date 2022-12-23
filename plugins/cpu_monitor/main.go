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
	severity := pluginpb.SEVERITY_INFO
	state := pluginpb.STATE_NONE
	var message string

	percent, err := cpu.Percent(defaultDuration * time.Second, false)
	if err != nil {
		ret := sdk.CallResponse{
			FuncName:	info["execute_method"].GetStringValue(),
			Message:	"failed to get cpu usage",
			Severity:	pluginpb.SEVERITY_CRITICAL,
			State:		pluginpb.STATE_FAILURE,
			AlertTypes:	[]pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
		}

		return ret, err
	}

	util := int(percent[0])
	log.Debug().Str("module", "plugin").Int("CPU Usage", util).Int("Urgent", urgent).Int("Warning", warning).Msg("cpu_monitor")

	if util < warning {
		message = fmt.Sprint("Current CPU Usage is ", util, "%, OK!")
		severity = pluginpb.SEVERITY_INFO
		state = pluginpb.STATE_SUCCESS
	} else if util < urgent {
		message = fmt.Sprint("Current CPU Usage is ", util, "%, over warning threshold ", warning, "%")
		severity = pluginpb.SEVERITY_WARNING
		state = pluginpb.STATE_SUCCESS
		log.Warn().Str("module", "plugin").Msg(message)
	} else {
		message = fmt.Sprint("Current CPU Usage is ", util, "%, over urgent threshold ", urgent, "%")
		severity = pluginpb.SEVERITY_CRITICAL
		state = pluginpb.STATE_SUCCESS
		log.Warn().Str("module", "plugin").Msg(message)
	}

	ret := sdk.CallResponse{
		FuncName:	info["execute_method"].GetStringValue(),
		Message:	message,
		Severity:	severity,
		State:		state,
		AlertTypes:	[]pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
