一个统计git仓库每个人贡献的工具

使用帮助：

```
git-996 是一个统计代码提交的工具

Usage:
  git-996 [flags]

Examples:
git-996 <git path>

Flags:
  -f, --format string         输出格式, table | json (default "table")
  -h, --help                  help for git-996
      --merge-email strings   合并人员 例如：user-to@mail.com=user-from@mail.com
  -r, --revert                逆序 (default true)
  -s, --sort string           排序方式： i | increase | d | decrease | l | left (default "l")
```

使用方式

```shell
git-996 <git path>
```

如果一个人有多个邮箱：

```shell
git-996 --merge-mail=user-to@mail.com=user-from@mail.com <git path>
```