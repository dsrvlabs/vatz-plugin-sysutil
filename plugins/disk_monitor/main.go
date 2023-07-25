package main

import (
	"flag"
	"fmt"

	"github.com/rs/zerolog/log"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/shirou/gopsutil/v3/disk"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

type arrayFlags []string

type mountInfo struct {
	path  string
	usage int
}

type statusInfo struct {
	severity pluginpb.SEVERITY
	state    pluginpb.STATE
	message  string
}

const (
	defaultAddr    = "127.0.0.1"
	defaultPort    = 9003
	pluginName     = "vatz-plugin-sysutil-disk-monitor"
	defaultUrgent  = 95
	defaultWarning = 90
	defaultPath    = "/"
)

var (
	urgent     int
	warning    int
	addr       string
	port       int
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
	var statusInfos []statusInfo
	var mountPoints []mountInfo

	// Collect amount of disk usage data from all mountpoints
	for i, path := range mountPaths {
		severity := pluginpb.SEVERITY_INFO
		state := pluginpb.STATE_NONE
		var message string

		usageStat, err := disk.Usage(path)
		if err != nil {
			message = fmt.Sprint("failed to get disk usage on ", path)
			ret := sdk.CallResponse{
				FuncName:   info["execute_method"].GetStringValue(),
				Message:    message,
				Severity:   pluginpb.SEVERITY_CRITICAL,
				State:      pluginpb.STATE_FAILURE,
				AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
			}

			return ret, err
		}
		mountPoints = append(mountPoints, mountInfo{path, int(usageStat.UsedPercent)})
		log.Debug().Str("module", "plugin").Int(mountPoints[i].path, mountPoints[i].usage).Int("Urgent", urgent).Int("Warning", warning).Msg("disk_monitor: Disk Usage(%) of")

		if mountPoints[i].usage < warning {
			severity = pluginpb.SEVERITY_INFO
			state = pluginpb.STATE_SUCCESS
			message = fmt.Sprint("Current Disk Usage of ", mountPoints[i].path, " ", mountPoints[i].usage, "%, OK!")
		} else if mountPoints[i].usage < urgent {
			severity = pluginpb.SEVERITY_WARNING
			state = pluginpb.STATE_SUCCESS
			message = fmt.Sprint("Current Disk Usage of ", mountPoints[i].path, " ", mountPoints[i].usage, "%, over warning threshold ", warning, "%")
			log.Warn().Str("module", "plugin").Msg(message)
		} else {
			severity = pluginpb.SEVERITY_CRITICAL
			state = pluginpb.STATE_SUCCESS
			message = fmt.Sprint("Current Disk Usage of ", mountPoints[i].path, " ", mountPoints[i].usage, "%, over urgent threshold ", urgent, "%")
			log.Warn().Str("module", "plugin").Msg(message)
		}
		statusInfos = append(statusInfos, statusInfo{severity, state, message})
	}

	// If we can reached here, state will be STATE_SUCCESS.
	severity := pluginpb.SEVERITY_INFO
	state := pluginpb.STATE_SUCCESS
	var message string

	for i, status := range statusInfos {
		if severity.Number() > status.severity.Number() {
			severity = status.severity
		}
		log.Debug().Str("module", "plugin").Msgf("%d status.severity %s", i, severity.String())
		message += status.message + "\n"
	}

	log.Debug().Str("module", "plugin").Msgf("severity : %s, state : %s, message : %s", severity.String(), state.String(), message)

	ret := sdk.CallResponse{
		FuncName:   info["execute_method"].GetStringValue(),
		Message:    message,
		Severity:   severity,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
