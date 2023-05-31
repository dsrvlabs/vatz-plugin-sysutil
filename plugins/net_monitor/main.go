package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	defaultAddr    = "127.0.0.1"
	defaultPort    = 9004
	pluginName     = "vatz-plugin-network-monitor"
	defaultUrgent  = 1000
	defaultWarning = 100
	INFO           = 0
	WARNING        = 1
	CRITICAL       = 2
)

var (
	urgent           int64
	warning          int64
	addr             string
	port             int
	prevRecvBytes    map[string]int64
	prevSentBytes    map[string]int64
	recvBytesDiff    map[string]int64
	sentBytesDiff    map[string]int64
	severity         map[string]int
	isFirstIteration bool
)

func init() {
	isFirstIteration = true
	prevRecvBytes = make(map[string]int64)
	prevSentBytes = make(map[string]int64)
	recvBytesDiff = make(map[string]int64)
	sentBytesDiff = make(map[string]int64)
	severity = make(map[string]int)

	flag.StringVar(&addr, "addr", defaultAddr, "Listening address")
	flag.IntVar(&port, "port", defaultPort, "Listening port")
	flag.Int64Var(&urgent, "urgent", defaultUrgent, "Network Traffic Alert threshold (in MBps)")
	flag.Int64Var(&warning, "warning", defaultWarning, "Network Traffic Warning threshold (in MBps)")

	flag.Parse()
}

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-exit:
				return
			case <-ticker.C:
				calculateTraffic()
			default:
			}
		}
	}()

	p := sdk.NewPlugin(pluginName)
	p.Register(pluginFeature)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := p.Start(ctx, addr, port); err != nil {
			log.Info().Str("module", "plugin").Msg("exit")
		}
	}()

	<-exit
	cancel()
	log.Info().Str("module", "plugin").Msg("exit!!")
}

func calculateTraffic() {
	// Open /proc/net/dev file
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		log.Error().Str("module", "plugin").Msgf("%v", err)
		return
	}
	defer file.Close()

	if isFirstIteration {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)

			// Check if the line contains network interface statistics
			if len(fields) > 1 && strings.Contains(fields[0], ":") {
				iface := strings.Trim(fields[0], ":")
				bytesRecv, _ := strconv.ParseInt(fields[1], 10, 64)
				bytesSent, _ := strconv.ParseInt(fields[9], 10, 64)
				// Initialize previous byte counts for this interface
				prevRecvBytes[iface] = bytesRecv
				prevSentBytes[iface] = bytesSent
			}
		}
		isFirstIteration = false

		return
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// Check if the line contains network interface statistics
		if len(fields) > 1 && strings.Contains(fields[0], ":") {
			iface := strings.Trim(fields[0], ":")
			bytesRecv, _ := strconv.ParseInt(fields[1], 10, 64)
			bytesSent, _ := strconv.ParseInt(fields[9], 10, 64)

			// Calculate the difference in bytes received and sent since the last iteration
			prevRecvBytesTotal := prevRecvBytes[iface]
			prevSentBytesTotal := prevSentBytes[iface]
			currentRecvBytesTotal := bytesRecv
			currentSentBytesTotal := bytesSent
			recvBytesDiff[iface] = currentRecvBytesTotal - prevRecvBytesTotal
			sentBytesDiff[iface] = currentSentBytesTotal - prevSentBytesTotal

			// Print statistics
			//log.Debug().Str("module", "plugin").Msgf("%s: received %d bytes, sent %d bytes\n", iface, recvBytesDiff[iface], sentBytesDiff[iface])

			// Update previous byte counts for this interface
			prevRecvBytes[iface] = currentRecvBytesTotal
			prevSentBytes[iface] = currentSentBytesTotal
		}
	}
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
	severityTotal := INFO
	severityToSend := pluginpb.SEVERITY_INFO
	state := pluginpb.STATE_NONE
	var message string

	// TODO: handle warning, error threshold
	for iface, recvBytes := range recvBytesDiff {
		if sentBytes, ok := sentBytesDiff[iface]; ok {
			recvMBytes := recvBytes / 1024 / 1024 / 60
			sentMBytes := sentBytes / 1024 / 1024 / 60
			if recvMBytes < warning && sentMBytes < warning {
				message += fmt.Sprintf("%s: NORMAL\n", iface)
				severity[iface] = INFO
			} else if recvMBytes > urgent || sentMBytes > urgent {
				message += fmt.Sprintf("%s: CRITICAL\n", iface)
				severity[iface] = CRITICAL
			} else {
				message += fmt.Sprintf("%s: WARNING\n", iface)
				severity[iface] = WARNING
			}
			message += fmt.Sprintf("%s: received %d Mbytes, sent %d Mbytes\n", iface, recvMBytes, sentMBytes)
			state = pluginpb.STATE_SUCCESS
		} else {
			message += fmt.Sprintf("%s: err to get data\n", iface)
			severity[iface] = CRITICAL
			state = pluginpb.STATE_FAILURE
		}

		if severityTotal < severity[iface] {
			severityTotal = severity[iface]
		}
	}
	log.Debug().Str("module", "plugin").Msg(message)

	if severityTotal == INFO {
		severityToSend = pluginpb.SEVERITY_INFO
	} else if severityTotal == WARNING {
		severityToSend = pluginpb.SEVERITY_WARNING
	} else {
		severityToSend = pluginpb.SEVERITY_CRITICAL
	}

	ret := sdk.CallResponse{
		FuncName:   info["execute_method"].GetStringValue(),
		Message:    message,
		Severity:   severityToSend,
		State:      state,
		AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
