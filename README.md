# cmd

## step0:

```sh
git clone https://github.com/monaco-io/cmd.git
cd cmd
```

## [ascii_art](ascii_art) 命令行炫酷文案

### Installation

step1:

```sh
make ascii_art
```

step2

```sh
# add this to your .zshrc/.bashrc
ascii_art --name=bilibili.com
# ascii_art --name=tencent.com
```

### Usage

```
# monaco @ monacos-MacBook-Pro in ~/workspace/cmd on git:main o [17:32:27]
$ zsh
                        //
           \\          //
            \\        //
    ##DDDDDDDDDDDDDDDDDDDDDDDD##
    ## DDDDDDDDDDDDDDDDDDDDDD ##   ________   ___   ___        ___   ________   ___   ___        ___
    ## hh                  hh ##   |\   __  \ |\  \ |\  \      |\  \ |\   __  \ |\  \ |\  \      |\  \
    ## hh     //    \\     hh ##   \ \  \|\ /_\ \  \\ \  \     \ \  \\ \  \|\ /_\ \  \\ \  \     \ \  \
    ## hh    //      \\    hh ##    \ \   __  \\ \  \\ \  \     \ \  \\ \   __  \\ \  \\ \  \     \ \  \
    ## hh                  hh ##     \ \  \|\  \\ \  \\ \  \____ \ \  \\ \  \|\  \\ \  \\ \  \____ \ \  \
    ## hh       wwww       hh ##      \ \_______\\ \__\\ \_______\\ \__\\ \_______\\ \__\\ \_______\\ \__\
    ## hh                  hh ##       \|_______| \|__| \|_______| \|__| \|_______| \|__| \|_______| \|__|
    ## MMMMMMMMMMMMMMMMMMMMMM ##
    ##MMMMMMMMMMMMMMMMMMMMMMMM##                    live.bilibili.com  game.bilibili.com  www.bilibili.com
```

## [fanyi](fanyi) 中英翻译工具

### Installation

```sh
make fanyi
```

### Usage

```sh
➜  ~ fanyi i have an apple
我有一个苹果
```

```sh
➜  ~ fanyi 我有一个苹果
I have an apple
```

## [timestamp](timestamp) 时间戳转换工具

### Installation

```sh
make timestamp
```

### Usage

```sh
# monaco @ monacos-MacBook-Pro in ~/workspace/cmd on git:main o [17:42:21]
$ timestamp now
2022/09/09 17:42:24 now
unix timestamp now is: 1662716544

# monaco @ monacos-MacBook-Pro in ~/workspace/cmd on git:main o [17:42:24]
$ timestamp 1662716544
2022/09/09 17:42:33 1662716544
tims is: 2022-09-09T17:42:24+08:00
```
