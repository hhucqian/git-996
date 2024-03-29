# 一个统计git仓库每个人贡献的工具

## 使用简介：

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

## 使用方式

```shell
git-996 <git path>
```

如果一个人有多个邮箱，可以使用`merge-email`来合并数据：

```shell
git-996 --merge-mail=user-to@mail.com=user-from@mail.com <git path>
```

## 详细说明

由于我这边主要是做一些外包项目，一个项目三五个人，三五个月完成。每完成一个项目，都需要回答一个问题：到底谁干的活最多？为了回答这个问题，我就做了这个工具。

下表是一个真实项目的统计结果：

|  E-MAIL | NAME  |  DAYS |  + | - | 产出 | 留存率 | 贡献率 |
|:---| :---  | ---: |  ---: | ---: | ---: | ---: | ---: |
|user1@mail.com| 张三  | 97 |  142496 | 90325 | 48775 | 34.23% | 93.52% |
|user2@mail.com| 李四 | 33 | 11418 | 8172 | 1921 | 16.82% | 3.68% |
|user3@mail.com| 王五 | 53 | 15074 | 5689 | 1457 | 9.67% | 2.79% |
|FROM:2021-11-25| TO: 2023-01-19 | 117 | 168988 | 104186 | 52153 | 30.86% |  |

1. 贡献率（功劳）作为最重要的数据，并不是计算一个人总共写了多少代码，而是计算这个项目在最终完成状态，有多少代码是他写的。这里需要注意，不同类型的代码最好单独一个Git仓库，前端代码有大量的HTML，后端代码大量的增删查改，代码数量级不一样，分开计算更科学。
2. 留存率（苦劳）最为一个辅助指标，也能直观反应项目开发情况。比如这个项目整理留存率只有30%，说明这个项目开发过程中经常需要推到重来，可能是甲方需求不明，也可能是本身设计有缺陷。如果一个人的代码留存率显著低于整体留存率，那也需要关注，是他负责的部分由于需求不明导致留存率低，还是本身代码质量比较差导致的。
3. 总的来说，如果一个人代码总是删了重写（自己或他人），那么他的贡献率、留存率就会很低，基本上是等于在做无用功。如果一个人留存率很低但是贡献率尚可，说明这个人工作很幸苦，也有成效。如果一个人留存率、贡献率很高，那么这个人就是个高手。
4. 目前这个工具只适合统计这种短期项目，对于需要长期维护的项目，还需要设计一个适合分时间段统计的模式。