<h1 align="center">cc-sesh</h1>

<p align="center">
  <em>给 Claude Code 用户的 sesh fork —— 在 picker 里看见每个 tmux session 内的 Claude 实时状态。</em>
</p>

<p align="center">
  <a href="https://opensource.org/licenses/MIT">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="MIT License" />
  </a>
  <a href="https://github.com/joshmedeski/sesh">
    <img src="https://img.shields.io/badge/forked%20from-joshmedeski%2Fsesh-blue.svg" alt="forked from joshmedeski/sesh" />
  </a>
</p>

---

## 致敬原作者

cc-sesh 是 [**joshmedeski/sesh**](https://github.com/joshmedeski/sesh) 的下游 fork。

底层的 session 管理、tmux/zoxide/tmuxinator 集成、配置系统、命名策略、picker 的全部基础能力都来自 Josh Medeski 与 sesh 上游贡献者们多年的设计与打磨。**没有 sesh，就没有 cc-sesh。**

- 上游仓库：<https://github.com/joshmedeski/sesh>
- 上游作者：[Josh Medeski](https://github.com/joshmedeski)（[赞助 sesh 项目](https://github.com/sponsors/joshmedeski)）
- 完整功能列表与配置文档请以 [上游 README](https://github.com/joshmedeski/sesh#readme) 为准

LICENSE 沿用 MIT，版权署名 © 2023 Josh Medeski（参见 [LICENSE](LICENSE) 文件）。本 fork 的修改部分同样以 MIT 发布。

如果你不需要下文描述的 Claude Code 集成，**请优先使用上游 sesh** —— 它的功能更稳定、社区更活跃、生态扩展（Raycast / Ulauncher / Walker 等）也更全。

---

## 这个 fork 在做什么

cc-sesh 只在上游 sesh 之上做了**一件事**：让 picker 里能直接看到每个 tmux session 内 [Claude Code](https://docs.claude.com/en/docs/claude-code) 进程的实时状态，并且在「跑完一轮活」时打个粘性提醒，方便我在多个并行 session 之间分配注意力。

具体差异如下。

### 1. picker 内置 Claude 状态表格

原 sesh picker 的每行只有 src icon + session 名。cc-sesh 在中间加了一张 4 列的状态表：

```
  ^a all  ^t tmux  ^g configs  ^x zoxide  ^f find  ^d kill
  ──────────────────────────────────────────────────────────
       ATTN IDLE RUN  WAIT
  >   tm    1                 my-feature-branch
      tm                      bay-translate-extension
      tm    ●    1            ai-dev-kit          done 15m ago
      tm              2       long-running-task
      tm                  1   oauth-flow
      tm    1    1   1        mixed
      ze                      ~/Code/backend/athena
      ze                      ~/AI-Workspace/bay-translate
```

- **ATTN**（橙色 ●）：粘性提醒。**和 Claude 当前状态无关** —— 一旦检测到某个 session 完成了一轮 busy/subagent → idle 转换就亮起，直到我 attach 进去（或手动 dismiss / kill）才消失。用来回答「我刚才让它跑的活到底跑完没？」
- **IDLE**：当前空闲的 Claude 进程数
- **RUN**：busy + subagent 之和（实际在干活）
- **WAIT**：等用户授权（OAuth、permission prompt 等）的进程数

这些数据通过扫描 `~/.claude/sessions/*.json` 拿到，按 cwd 关联到 tmux pane，再聚合到 session。**不依赖任何 Claude Code 内部接口、不修改 Claude Code 任何配置**，纯只读。

### 2. picker 补齐的 hotkey

为了不再依赖 `fzf-tmux` 包一层，把上游 fzf 路径里那套 hotkey 直接补进了 picker：

| Hotkey | 行为 |
|---|---|
| `Ctrl-A` | all（默认 list） |
| `Ctrl-T` | 仅 tmux session |
| `Ctrl-G` | 仅 sesh.toml 配置 |
| `Ctrl-X` | 仅 zoxide 历史 |
| `Ctrl-F` | 在 `$HOME` 下深度 ≤ 2 列目录（替代 `fd`） |
| `Ctrl-D` | kill 当前光标所指的 tmux session |
| `Alt-D` | dismiss 当前行的 ATTN 标记（不 kill session） |

### 3. 重命名

为了不和已经装了上游 sesh 的环境冲突：

| 项目 | 上游 sesh | cc-sesh |
|---|---|---|
| 二进制 | `sesh` | `cc-sesh` |
| Go module | `github.com/joshmedeski/sesh/v2` | `github.com/Wingsdh/cc-sesh/v2` |
| 配置目录 | `$XDG_CONFIG_HOME/sesh/` | `$XDG_CONFIG_HOME/cc-sesh/` |
| 配置文件名 | `sesh.toml` | `sesh.toml`（沿用文件名） |
| 状态目录 | — | `$XDG_STATE_HOME/cc-sesh/`（仅本 fork 新增的 attention 状态） |

> 也就是说，你可以在同一台机器上同时装 `sesh` 和 `cc-sesh`，两套配置互不干扰。

---

## 安装

### Homebrew

```sh
brew install Wingsdh/cc-sesh/cc-sesh
```

通过我自维护的 [Homebrew tap](https://github.com/Wingsdh/homebrew-cc-sesh) 安装（不在 homebrew-core）。每次推 release tag 时由 GoReleaser 自动更新 formula。

### Go install

```sh
go install github.com/Wingsdh/cc-sesh/v2@latest
```

需要 Go 1.25+。

---

装完后二进制叫 `cc-sesh`，所有子命令名与上游 sesh 一致（`list / connect / picker / window / ...`）。

> 暂未提供 AUR / Conda / Nix 等打包 —— 这是一个为了我自己用而 fork 的项目。

---

## 用法

### 基本命令

完全同上游 sesh，把所有 `sesh` 替换成 `cc-sesh` 即可：

```sh
cc-sesh list           # 列出所有 session 来源（tmux + config + zoxide + tmuxinator）
cc-sesh connect <name> # connect 到 session（不存在则创建）
cc-sesh picker         # 打开内置 picker（推荐）
```

`list / connect / window / pane / clone / root / last` 等子命令的语义与 flag 与上游一致，**详见 [上游 README](https://github.com/joshmedeski/sesh#readme)**。

### 推荐绑定（tmux）

cc-sesh 的卖点是 picker，所以最佳用法是把它绑成 tmux popup：

```tmux
bind-key "K" display-popup -h 90% -w 60% -E "cc-sesh picker -i"
```

`-i` 显示 src icon（需要 Nerd Font）。

### Claude Code 集成的工作方式

picker 在 fetch 阶段会做这几件事：

1. 调上游 lister 拿 session 列表（tmux / zoxide / config / tmuxinator）
2. 扫 `~/.claude/sessions/*.json`，过滤出活进程
3. 用 `tmux list-panes` 拿每个 pane 的 cwd，把 Claude 进程的 cwd 匹配到 tmux session
4. 跟上一轮的 busy 状态对比，识别 busy/subagent → idle 转换并落入 `~/.local/state/cc-sesh/attention.json`
5. 把每个 session 的 LiveBadge + Attention 一起贴到行上渲染

整个链路对没有 Claude 的环境完全透明 —— 任何一步失败（`~/.claude/sessions/` 不存在、tmux 不在跑、JSON 解析失败……）都会降级回上游 sesh 行为，picker 该列就是空白。

### 配置

配置文件路径变成 `~/.config/cc-sesh/sesh.toml`，**配置 schema 与上游 sesh 完全一致**（`[default_session]` / `[[session]]` / `[[wildcard]]` / `[tui]` / `blacklist` / `dir_length` / `sort_order` / `cache` 等都原样可用）。

配置写法请直接看 [上游 README 的 Configuration 一节](https://github.com/joshmedeski/sesh#configuration)。

cc-sesh 本身没有新增任何配置项 —— Claude 集成是开箱即用、**不需要也不接受任何配置**。

### Attention 状态文件

为了让粘性提醒跨进程生效，cc-sesh 把 attention 状态持久化到：

```
$XDG_STATE_HOME/cc-sesh/attention.json
# 默认：~/.local/state/cc-sesh/attention.json
```

文件可以随时手动删，删完 attention 全清，下次 picker 打开会重新积累。

---

## 与上游的同步策略

- 仓库的 `main` 分支会不定期 rebase / merge 上游的 release tag
- 我自己的修改集中在：
  - `claude/` 包（live + attention）—— 全新增
  - `picker/` 包 —— UI 与 hotkey 改造
  - `seshcli/claude_wire.go` —— 把 lister + claude/live + claude/attention 串起来注入 picker
  - module path、二进制名、配置路径的全局重命名

如果你想把 Claude 集成提交回上游，请直接去和 [Josh Medeski](https://github.com/joshmedeski) 讨论 —— 我没有这个计划，因为这个改动对绝大多数 sesh 用户都是不必要的复杂度。

---

## License

MIT，沿袭上游 [sesh 的 LICENSE](LICENSE)，版权署名 © 2023 Josh Medeski。本 fork 新增部分同样以 MIT 发布。
