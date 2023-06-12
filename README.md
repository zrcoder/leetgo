# leetgo

Cli APP for Leetcode.

![examples](example.gif)

```shell
USAGE:
   leetgo command [command options] [arguments...]

COMMANDS:
   config   show or set config of your leetgo project
   search   search questions by keywords or id
   view     view a question by id
   code     edit codes to solve the question
   test     test your codes locally and remotely
   submit   submit your codes
   book     view all questions you picked as a book
```

> Only tested for Golang until now.

## How to auth?

You should login leetcode on browser firstly, then leetgo will search and read the browser cookie for auth.

> For security, leetcode doesn't store the cookies itself, just read from browsers' cache.

## Inspired by

<https://github.com/j178/leetgo>
