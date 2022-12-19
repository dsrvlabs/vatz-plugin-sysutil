package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"flag"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	defaultAddr = "127.0.0.1"
	defaultPort = 9095
	pluginName = "vatz-plugin-solana-mem-monitor"
	defaultUrgent = 95
	defaultWarning = 90
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
	flag.IntVar(&urgent, "urgent", defaultUrgent, "Mem Usage Alert threshold")
	flag.IntVar(&warning, "warning", defaultWarning, "Mem Usage Warning threshold")

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

	v, err := mem.VirtualMemory()
	if err != nil {
		ret := sdk.CallResponse{
			FuncName:	info["execute_method"].GetStringValue(),
			Message:	"failed to get memory usage",
			Severity:	pluginpb.SEVERITY_CRITICAL,
			State:		pluginpb.STATE_FAILURE,
			AlertTypes:	[]pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
		}

		return ret, err
	}
	state = pluginpb.STATE_SUCCESS
	usage := int(v.UsedPercent)
	log.Debug().Str("module", "plugin").Int("Memory Usage", int(v.UsedPercent)).Int("Urgent", urgent).Int("Warning", warning).Msg("mem_monitor")

	if usage < warning {
		message = fmt.Sprint("Current Memory Usage is ", usage, "%, OK!")
		severity = pluginpb.SEVERITY_INFO
	} else if usage < urgent {
		message = fmt.Sprint("Current Memory Usage is ", usage, "%, over warning threshold ", warning, "%")
		severity = pluginpb.SEVERITY_WARNING
		log.Warn().Str("module", "plugin").Msg(message)
	} else {
		message = fmt.Sprint("Current Memory Usage is ", usage, "%, over urgent threshold ", urgent, "%")
		severity = pluginpb.SEVERITY_CRITICAL
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
