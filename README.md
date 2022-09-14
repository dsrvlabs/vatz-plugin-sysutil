# vatz-plugin-sysutil
vatz plugin for system utilization monitoring

## Plugins
- cpu_monitor : monitor CPU utilization
- disk_monitor : monitor disk space utilization
- mem_monitor : monitor memory utilization

## Installation and Usage
> Please make sure [Vatz](https://github.com/dsrvlabs/vatz) is running with proper configuration. [Vatz Installation Guide](https://github.com/dsrvlabs/vatz/blob/main/docs/installation.md)

### Install
```
$ make install
```

### Run
> Run as default config
```
$ cpu_monitor
2022-09-14T08:17:33+02:00 INF Register module=grpc
2022-09-14T08:17:33+02:00 INF Start 127.0.0.1 9094 module=sdk
2022-09-14T08:17:33+02:00 INF Start module=grpc
2022-09-14T08:17:48+02:00 INF Execute module=grpc
2022-09-14T08:17:48+02:00 DBG cpu_monitor CPU Usage=26 Urgent=95 Warning=90 module=plugin
```
```
$ disk_monitor
2022-09-14T08:19:19+02:00 INF Register module=grpc
2022-09-14T08:19:19+02:00 INF Start 127.0.0.1 9096 module=sdk
2022-09-14T08:19:19+02:00 INF Start module=grpc
2022-09-14T08:19:33+02:00 INF Execute module=grpc
2022-09-14T08:19:33+02:00 DBG disk_monitor: Disk Usage(%) of /=34 Urgent=95 Warning=90 module=plugin
```
```
$ mem_monitor
2022-09-14T08:20:12+02:00 INF Register module=grpc
2022-09-14T08:20:12+02:00 INF Start 127.0.0.1 9095 module=sdk
2022-09-14T08:20:12+02:00 INF Start module=grpc
2022-09-14T08:20:13+02:00 INF Execute module=grpc
2022-09-14T08:20:13+02:00 DBG mem_monitor Memory Usage=72 Urgent=95 Warning=90 module=plugin 
```

## Command line arguments
- cpu_monitor
```
Usage of cpu_monitor:
  -addr string
    	Listening address (default "127.0.0.1")
  -port int
    	Listening port (default 9094)
  -urgent int
    	CPU Usage Alert threshold (default 95)
  -warning int
    	CPU Usage Warning threshold (default 90)
```
- disk_monitor
```
Usage of disk_monitor:
  -addr string
    	Listening address (default "127.0.0.1")
  -paths value
    	mount paths to check available size (default "/")
  -port int
    	Listening port (default 9096)
  -urgent int
    	Disk Usage Alert threshold (default 95)
  -warning int
    	Disk Usage Warning threshold (default 90)
```
- mem_monitor
```
Usage of mem_monitor:
  -addr string
    	Listening address (default "127.0.0.1")
  -port int
    	Listening port (default 9095)
  -urgent int
    	Mem Usage Alert threshold (default 95)
  -warning int
    	Mem Usage Warning threshold (default 90)
```
