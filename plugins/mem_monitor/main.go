package main

import (
	"os"
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
	ret := sdk.CallResponse{
		FuncName:	"getMEMUsage",
		Message:	"Memory usage warning!",
		Severity:	pluginpb.SEVERITY_UNKNOWN,
		State:		pluginpb.STATE_NONE,
		AlertTypes:	[]pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal().Str("module", "plugin").Msgf("Couldn't get hostname: %v", err)
	}

	v, _ := mem.VirtualMemory()
	log.Debug().Str("module", "plugin").Int("Memory Usage", int(v.UsedPercent)).Int("Urgent", urgent).Int("Warning", warning).Msg("mem_monitor")

	if int(v.UsedPercent) > urgent {
		var message string
		message = fmt.Sprint("[", hostname, "]\n", "Current Memory Usage is ", int(v.UsedPercent), "%, over urgent threshold ", urgent, "%")
		ret = sdk.CallResponse{
			FuncName:	"getMEMUsage",
			Message:	message,
			Severity:	pluginpb.SEVERITY_CRITICAL,
			State:		pluginpb.STATE_FAILURE,
			AlertTypes:	[]pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
		}

		log.Warn().Str("module", "plugin").Msg(message)
	}

	return ret, nil
}
