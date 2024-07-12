# leetgo

Cli APP for Leetcode.

![examples](example.gif)

> Only tested for Golang and Python until now.

```shell
USAGE:
   leetgo command [command options] [arguments...]

COMMANDS:
   config   init or config your leetgo project
   search   search questions by keywords
   view     view questions or solutions
   test     test your code locally and remotely
   submit   submit your codes
```

## Install

```shell
go install github.com/zrcoder/leetgo@latest
```

## How to auth?

You should login leetcode on browser firstly, then leetgo will search and read the browser cookie for auth.

> For security, leetcode doesn't store the cookies itself, just read from browsers' cache.

## Inspired by

<https://github.com/j178/leetgo>
