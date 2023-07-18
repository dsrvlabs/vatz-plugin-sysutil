# vatz-plugin-sysutil
vatz plugin for system utilization monitoring
> Please, be advised that VATZ plugin supports only Linux operating systems, other operating systems may not operate properly.
 
## Plugins
- cpu_monitor : monitor CPU utilization
- disk_monitor : monitor disk space utilization
- mem_monitor : monitor memory utilization
- net_monitor : monitor network utilization

## Installation and Usage
> Please make sure [`Vatz`](https://github.com/dsrvlabs/vatz) is running with proper configuration. [Vatz Installation Guide](https://github.com/dsrvlabs/vatz/blob/main/docs/installation.md).


### Install Plugins
- Install with source
```
$ git clone https://github.com/dsrvlabs/vatz-plugin-sysutil.git
$ cd vatz-plugin-sysutil
$ make install
```
- Install with Vatz CLI command
```
$ vatz plugin install --help
Install new plugin

Usage:
   plugin install [flags]

Examples:
vatz plugin install github.com/dsrvlabs/<somewhere> name

Flags:
  -h, --help   help for install
```
> Please make sure install path for the plugins repository URL.
```
$ vatz plugin install github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor vatz_cpu_monitor
$ vatz plugin install github.com/dsrvlabs/vatz-plugin-sysutil/plugins/mem_monitor vatz_mem_monitor
$ vatz plugin install github.com/dsrvlabs/vatz-plugin-sysutil/plugins/disk_monitor vatz_disk_monitor
$ vatz plugin install github.com/dsrvlabs/vatz-plugin-sysutil/plugins/net_monitor vatz_net_monitor
```
- Check plugins list with Vatz CLI command
```
$ vatz plugin list
2023-05-31T01:57:35Z INF List plugins module=plugin
2023-05-31T01:57:35Z INF Create DB Instance module=db
2023-05-31T01:57:35Z INF List module=plugin
2023-05-31T01:57:35Z INF List module=db
+-------------------+------------+---------------------+--------------------------------------------------------------+---------+
| NAME              | IS ENABLED | INSTALL DATE        | REPOSITORY                                                   | VERSION |
+-------------------+------------+---------------------+--------------------------------------------------------------+---------+
| vatz_net_monitor  | true       | 2023-05-24 06:14:01 | github.com/dsrvlabs/vatz-plugin-sysutil/plugins/net_monitor  | latest  |
| vatz_cpu_monitor  | true       | 2023-05-31 01:57:07 | github.com/dsrvlabs/vatz-plugin-sysutil/plugins/cpu_monitor  | latest  |
| vatz_mem_monitor  | true       | 2023-05-31 01:57:21 | github.com/dsrvlabs/vatz-plugin-sysutil/plugins/mem_monitor  | latest  |
| vatz_disk_monitor | true       | 2023-05-31 01:57:32 | github.com/dsrvlabs/vatz-plugin-sysutil/plugins/disk_monitor | latest  |
+-------------------+------------+---------------------+--------------------------------------------------------------+---------+
```

### Run
> Run as default config
```
$ ./cpu_monitor
2023-05-31T07:30:21Z INF Register module=grpc
2023-05-31T07:30:21Z INF Start 127.0.0.1 9001 module=sdk
2023-05-31T07:30:21Z INF Start module=grpc
2023-05-31T07:30:54Z INF Execute module=grpc
2023-05-31T07:30:54Z DBG cpu_monitor CPU Usage=1 Urgent=95 Warning=90 module=plugin
```
```
$ ./mem_monitor
2023-05-31T07:31:04Z INF Register module=grpc
2023-05-31T07:31:04Z INF Start 127.0.0.1 9002 module=sdk
2023-05-31T07:31:04Z INF Start module=grpc
2023-05-31T07:32:24Z INF Execute module=grpc
2023-05-31T07:32:24Z DBG mem_monitor Memory Usage=3 Urgent=95 Warning=90 module=plugin
```
```
$ ./disk_monitor
2023-05-31T07:33:37Z INF Register module=grpc
2023-05-31T07:33:37Z INF Start 127.0.0.1 9003 module=sdk
2023-05-31T07:33:37Z INF Start module=grpc
2023-05-31T07:33:54Z INF Execute module=grpc
2023-05-31T07:33:54Z DBG disk_monitor: Disk Usage(%) of /=37 Urgent=95 Warning=90 module=plugin
2023-05-31T07:33:54Z DBG 0 status.severity INFO module=plugin
2023-05-31T07:33:54Z DBG severity : INFO, state : SUCCESS, message : Current Disk Usage of / 37%, OK!
```
```
$ ./net_monitor
2023-05-31T07:45:19Z INF Register module=grpc
2023-05-31T07:45:19Z INF Start 127.0.0.1 9004 module=sdk
2023-05-31T07:45:19Z INF Start module=grpc
2023-05-31T07:45:46Z INF Execute module=grpc
2023-05-31T07:45:46Z DBG  module=plugin
2023-05-31T07:46:16Z INF Execute module=grpc
2023-05-31T07:46:16Z DBG  module=plugin
2023-05-31T07:46:46Z INF Execute module=grpc
2023-05-31T07:46:46Z DBG  module=plugin
2023-05-31T07:47:16Z INF Execute module=grpc
2023-05-31T07:47:16Z DBG  module=plugin
2023-05-31T07:47:46Z INF Execute module=grpc
2023-05-31T07:47:46Z DBG lo: NORMAL
lo: received 0 Mbytes, sent 0 Mbytes
enp4s0: NORMAL
enp4s0: received 10 Mbytes, sent 0 Mbytes
wlp3s0: NORMAL
wlp3s0: received 0 Mbytes, sent 0 Mbytes
docker0: NORMAL
docker0: received 0 Mbytes, sent 0 Mbytes
mpqemubr0: NORMAL
mpqemubr0: received 0 Mbytes, sent 0 Mbytes
 module=plugin
```
## Command line arguments
- cpu_monitor
```
Usage of cpu_monitor:
  -addr string
    	Listening address (default "127.0.0.1")
  -port int
    	Listening port (default 9001)
  -urgent int
    	CPU Usage Alert threshold (default 95)
  -warning int
    	CPU Usage Warning threshold (default 90)
```
- mem_monitor
```
Usage of mem_monitor:
  -addr string
    	Listening address (default "127.0.0.1")
  -port int
    	Listening port (default 9002)
  -urgent int
    	Mem Usage Alert threshold (default 95)
  -warning int
    	Mem Usage Warning threshold (default 90)
```
- disk_monitor
```
Usage of disk_monitor:
  -addr string
    	Listening address (default "127.0.0.1")
  -paths value
    	mount paths to check available size (default "/")
  -port int
    	Listening port (default 9003)
  -urgent int
    	Disk Usage Alert threshold (default 95)
  -warning int
    	Disk Usage Warning threshold (default 90)
```
- net_monitor
```
Usage of ./net_monitor:
  -addr string
    	Listening address (default "127.0.0.1")
  -port int
    	Listening port (default 9004)
  -urgent int
    	Network Traffic Alert threshold (in MBps) (default 1000)
  -warning int
    	Network Traffic Warning threshold (in MBps) (default 100)
```

## License

`vatz-plugin-sysutil` is licensed under the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included in our repository in the `LICENSE` file.
