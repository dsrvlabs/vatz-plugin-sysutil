package main

import (
	"bufio"
	"flag"
	"os"
	"strings"
	"time"
	"fmt"
	"strconv"

	pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
	"github.com/dsrvlabs/vatz/sdk"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	defaultAddr    = "127.0.0.1"
	defaultPort    = 9101
	pluginName     = "vatz-plugin-network-monitor"
	defaultUrgent  = 100
	defaultWarning = 95
)

var (
	urgent        int
	warning       int
	addr          string
	port          int
	prevRecvBytes	map[string]int64
	prevSentBytes	map[string]int64
	recvBytesDiff	map[string]int64
	sentBytesDiff	map[string]int64
	isFirstIteration bool
)

func init() {
	isFirstIteration = true
	prevRecvBytes = make(map[string]int64)
	prevSentBytes = make(map[string]int64)
	recvBytesDiff = make(map[string]int64)
	sentBytesDiff = make(map[string]int64)

	flag.StringVar(&addr, "addr", defaultAddr, "Listening address")
	flag.IntVar(&port, "port", defaultPort, "Listening port")
	flag.IntVar(&urgent, "urgent", defaultUrgent, "Network Traffic Alert threshold")
	flag.IntVar(&warning, "warning", defaultWarning, "Network Traffic Warning threshold")

	flag.Parse()
}

func main() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			calculateTraffic()
		}
	}()

	p := sdk.NewPlugin(pluginName)
	p.Register(pluginFeature)

	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		log.Info().Str("module", "plugin").Msg("exit")
	}
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
			log.Debug().Str("module", "plugin").Msgf("%s: received %d bytes, sent %d bytes\n", iface, recvBytesDiff[iface], sentBytesDiff[iface])

			// Update previous byte counts for this interface
			prevRecvBytes[iface] = currentRecvBytesTotal
			prevSentBytes[iface] = currentSentBytesTotal
		}
	}
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
	severity := pluginpb.SEVERITY_INFO
	state := pluginpb.STATE_NONE
	var message string

	for iface, recvBytes := range recvBytesDiff {
		if sentBytes, ok := sentBytesDiff[iface]; ok {
			message += fmt.Sprintf("%s: received %d bytes, sent %d bytes\n", iface, recvBytes, sentBytes) 
		}
	}
	// TODO: handle warning, error threshold
	log.Debug().Str("module", "plugin").Msg(message)

	ret := sdk.CallResponse{
		FuncName:	info["execute_method"].GetStringValue(),
		Message:	message,
		Severity:	severity,
		State:		state,
		AlertTypes:	[]pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
	}

	return ret, nil
}
