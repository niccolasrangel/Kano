---
title: 命令行
description: 用 kano 命令行，从 shell 或脚本中操控 Kano 的会话和窗格。
---

# 命令行

`kano` 命令可以在任何 shell 中操控会话：你自己的、脚本的，或者智能体的。通过 **Kano › 安装命令行工具…** 安装它。

## 目标

大多数命令都接受一个**目标**：一个会话或一个窗格。目标可以是：

- 一个 id，或它任意一个唯一的前缀
- 会话的名称
- 一个会话，此时指它当前聚焦的窗格

在 Kano 中，`$KANO_SESSION` 和 `$KANO_PANE` 指向命令所在的会话和窗格，也是默认的目标。

## 会话和窗格

```sh
kano ls                          # 列出所有会话及其窗格
kano new -n api -d ~/code/api    # 新建会话，并打印它的 id
kano new -- npm run dev          # 新建一个运行某条命令的会话
kano split --down -- npm test    # 在当前窗格旁新建一个窗格
kano split --url localhost:3000  # 新建一个显示网页的浏览器窗格
kano rename api backend          # 重命名会话
kano focus backend               # 在 app 中显示它
kano close backend               # 结束会话（如果给的是窗格，则关闭该窗格）
kano restart <pane>              # 重新启动已退出窗格的 shell
```

`kano new` 的选项：

| 参数                        | 含义                                            |
| --------------------------- | ----------------------------------------------- |
| `-p`, `--project <name\|id>` | 放入哪个项目（默认按目录决定）                  |
| `-n`, `--name <name>`       | 为会话命名                                      |
| `-d`, `--dir <dir>`         | 在该目录启动（默认为当前目录）                  |
| `-f`, `--focus`             | 在 app 中显示它                                 |
| `-w`, `--wait`              | 等待它的命令结束，并以该命令的退出码退出        |
| `-u`, `--url <url>`         | 打开显示该网页的浏览器，而不是 shell            |

## 输入与读取

```sh
kano send api "npm run migrate" --enter   # 输入文本，然后按回车
kano send api -e 'y\n'                    # 解析 \n \r \t \e \xHH 转义
kano key api ctrl+c                       # 按键
kano read api -n 40                       # 屏幕上最后 40 行
kano read api -a                          # 屏幕内容连同回滚缓冲
```

`kano key` 支持 `enter`、`tab`、`escape`、`up`、`down`、`left`、`right`、`backspace`，以及 `ctrl+c`、`ctrl+d` 这类带修饰键的组合。

## 等待

```sh
kano wait api                             # 等到它不再工作
kano wait api --for waiting,error --timeout 10m
kano prompt review "summarize the diff"   # 发送并回车，等待，然后打印屏幕
```

`kano wait` 返回时会打印窗格的状态。`kano prompt` 发送文本并回车，等到窗格不再工作，再打印它的屏幕。这是向另一个智能体提问最快的方式。

### 等待命令结束

像 `make test` 这样的命令不会显示为*工作中*，所以要等它跑完，就用 `--wait` 启动它。`kano` 随后会以该命令自身的退出码退出：

```sh
if kano split --down --wait -- make test; then
  echo "tests passed"
fi
```

对于一个已经在运行启动时所带命令的窗格，`kano wait --exit` 起同样的作用：

```sh
pane=$(kano split --down -- go test ./...)
# … 做些别的事 …
kano wait "$pane" --exit --timeout 10m
echo "go test exited with $?"
```

命令结束后，窗格连同输出会保留下来，你可以用 `kano read` 读取。如果 shell 是 `cmd`，则无法等待命令结束，因为它判断不出命令何时结束。

<a id="browsing"></a>
## 浏览网页

`kano browser` 操控浏览器窗格中的网页，就像 `kano send` 和 `kano read` 操控终端一样：点击、输入、读取页面、截图。它的目标是一个浏览器窗格，或者一个会话（指其中的浏览器）。不指定目标时，它操控所在会话的浏览器。

```sh
kano split --url localhost:3000            # 打开网页，并打印窗格的 id
kano browser open localhost:3000/todos     # 前往某个网页，并等待加载完成
kano browser type '#new-todo' "Buy milk" --enter
kano browser wait '.todo-list li'          # 等到某个元素出现
kano browser wait --text "1 item left"     # 等到某段文本出现
kano browser text                          # 页面的文本
kano browser text h1                       # 某个元素的文本
kano browser click 'button[type=submit]'
kano browser screenshot -o page.png
kano browser eval 'document.title'
```

| 命令                            | 作用 |
| ------------------------------- | ---- |
| `open <url>`                    | 前往某个网页，并等待加载完成 |
| `back`, `forward`, `reload`     | 后退、前进或重新加载页面，并等待完成 |
| `text [selector]`               | 打印页面的可见文本，或第一个匹配元素的文本 |
| `html [selector]`               | 打印页面的 HTML，或该元素的 HTML |
| `click <selector>`              | 像鼠标一样点击该元素 |
| `type <selector> <text…>`       | 向输入框输入文本，替换原有内容；`--enter` 会在之后按回车，从而提交所在表单 |
| `wait [selector] [--text text]` | 等到元素或文本出现；两者都不给时，等到页面加载完成 |
| `eval <js…>`                    | 在页面中运行 JavaScript 并打印结果：字符串原样输出，其他值输出为 JSON |
| `screenshot [-o file.png]`      | 保存页面截图，并打印保存位置；`-o -` 会写到标准输出 |

选择器是 CSS 选择器，在 shell 中记得加引号。需要元素的命令会等它出现，`--timeout`（默认 `10s`）设定这些命令以及页面加载最多等多久。`type` 和 `eval` 接受 `-`，表示从标准输入读取文本。

点击和输入发送的是真实鼠标和键盘会产生的事件，所以 React 等框架构建的页面也能收到。截图需要窗格显示在屏幕上：先用 `kano focus` 显示它。其余命令对任何浏览器窗格都有效，无论是否显示。

和 `kano do` 一样，`kano browser` 需要开启 **设置 › 智能体** 中的 **远程控制**，并且 app 处于打开状态。

## 状态

```sh
kano status api --json                    # 状态、程序、目录
kano report working "Running migrations"  # 报告自己的状态
kano notify "Need your input"             # 将此窗格标记为需要你处理
```

各状态的含义见 [状态与通知](./status.md)。

## 改动

```sh
kano diff                                 # 窗格中未提交的改动，格式同 git diff
kano diff --turn                          # 它的智能体在上一轮中改了什么
kano diff review --branch                 # 自默认分支以来的改动
```

见 [改动](./changes.md)。

## 配置与设置

```sh
kano config path                 # 配置文件所在位置
kano config show                 # 每个配置项及其值
kano config show --defaults      # 每个配置项的默认值
kano config edit                 # 在 $EDITOR 中打开它
kano config check                # 报告未知的配置项和无效的值
kano projects                    # 列出项目
kano skill setup                 # 教编程智能体使用 kano
kano daemon status|stop|restart  # 管理 kanod
kano daemon install-service      # 让 kanod 随系统启动（systemd、launchd）
kano daemon serve                # 在前台运行 kanod，供服务管理器使用
kano clients                     # 查看谁连接着 kanod
```

## 其他机器

在命令前加 `-H <host>`，或设置 `$KANO_HOST`，可以让操作会话的命令指向另一台机器上的 kanod：可以是你添加过的主机，也可以是 `me@devbox`、`wsl://Ubuntu` 这样的地址。见 [远程主机](./remote.md)。

```sh
kano hosts                       # 你添加过的主机
kano hosts add devbox me@devbox  # 添加一台，app 中也会显示
kano hosts check devbox          # 尝试连接，并报告结果
kano hosts install devbox        # 通过 SSH 在那里安装 Kano
kano -H devbox ls
kano -H devbox new -n review -- claude
kano -H devbox daemon update     # 将那里的 Kano 更新到本机的版本
```

在窗格中，`kano` 总是与窗格所在机器的 kanod 通信，无论 app 在哪里。

## 执行 app 操作

`kano do` 可以执行 app 的任意操作，与菜单中的操作相同。它需要 **远程控制**，该选项在 **设置 › 智能体** 中默认开启。

```sh
kano actions                     # 列出所有操作
kano do pane.split.right
kano do pane.browser url=localhost:3000
kano do sidebar.toggle
```

## 退出码

| 退出码 | 含义                           |
| ------ | ------------------------------ |
| `0`    | 成功                           |
| `1`    | 失败                           |
| `2`    | 用法错误                       |
| `3`    | 目标不存在，或有歧义           |
| `4`    | 无法连接 kanod                 |

等待程序的命令，即 `kano wait --exit` 和 `--wait`，改为以该程序的退出码退出；当窗格的 shell 以错误退出时，`kano wait` 也是如此。等待超时则以 `1` 退出。

给 `ls`、`status` 和 `projects` 加上 `--json`，可得到便于程序处理的输出。
