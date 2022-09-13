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
	"github.com/shirou/gopsutil/v3/disk"
)

type arrayFlags []string

const (
	defaultAddr = "127.0.0.1"
	defaultPort = 9096
	pluginName = "vatz-plugin-solana-disk-monitor"
	defaultUrgent = 95
	defaultWarning = 90
	defaultPath = "/"
)

var (
	urgent int
	warning int
	addr string
	port int
	mountPaths arrayFlags
)

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func init() {
	mountPaths.Set(defaultPath)
	flag.StringVar(&addr, "addr", defaultAddr, "Listening address")
	flag.IntVar(&port, "port", defaultPort, "Listening port")
	flag.IntVar(&urgent, "urgent", defaultUrgent, "Disk Usage Alert threshold")
	flag.IntVar(&warning, "warning", defaultWarning, "Disk Usage Warning threshold")
	flag.Var(&mountPaths, "paths", "mount paths to check available size (default \"/\")")

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
		FuncName:	"getDISKUsage",
		Message:	"Disk usage warning!",
		Severity:	pluginpb.SEVERITY_UNKNOWN,
		State:		pluginpb.STATE_NONE,
		AlertTypes:	[]pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal().Str("module", "plugin").Msgf("Couldn't get hostname: %v", err)
	}

	used := make([]int, len(mountPaths))
	for i := 0; i < len(mountPaths); i++ {
		temp, _ := disk.Usage(mountPaths[i])
		used[i] = int(temp.UsedPercent)
		log.Debug().Str("module", "plugin").Int(mountPaths[i], used[i]).Int("Urgent", urgent).Int("Warning", warning).Msg("disk_monitor: Disk Usage(%) of")
		if used[i] > urgent {
			var message string
			message = fmt.Sprint("[", hostname, "]\n", "Current Disk Usage of ", mountPaths[i], used[i], "%, over urgent threshold ", urgent, "%")
			ret = sdk.CallResponse{
				FuncName:	"getDISKUsgae",
				Message:	message,
				Severity:	pluginpb.SEVERITY_CRITICAL,
				State:		pluginpb.STATE_SUCCESS,
				AlertTypes:	[]pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
			}

			log.Warn().Str("module", "plugin").Msg(message)
		}
	}

	return ret, nil
}
