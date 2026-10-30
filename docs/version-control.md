# Git Tag With Version

本实现复用 YoHub 的语义版本控制规则，将 GitHub/Gitea 的协议差异收敛到 `scripts/forge-api.mjs`。两个目标仓库使用相同脚本与工作流，不依赖 YoHub 前端、私有 Action 镜像、内网代理或固定仓库名。

## 自动打标规则

- `fix` 升 PATCH，`feat` 升 MINOR，提交标题的 `!` 或正文 `BREAKING CHANGE:` 升 MAJOR，其他提交默认不升级。发布范围取最高级别，提交数量只用于开发版身份。
- PR 可设置一个 `semver:major`、`semver:minor`、`semver:patch` 或 `semver:none` 标签，仅覆盖对应 PR 的提交。冲突、无效或重叠覆盖会阻止打标。
- 基线为发布分支可达的最高 `vMAJOR.MINOR.PATCH`，无标签时从 `v0.0.0` 起算。历史标签不改写。已有轻量正式标签可作历史基线，但正式构建与新标签必须使用附注标签。
- 所有分支 push、PR 和人工触发均先执行版本契约、Go 编译、vet、自有代码单测。仅 `main` 的 push 或人工流程在同次门禁通过后创建标签。
- MAJOR 不自动批准：日志给出候选版本和完整 SHA，在 Actions → **Git Tag With Version** → Run workflow 选择 `main`，填写 `approved_sha`。主线前进后，旧 SHA 不再有效。
- 打标作业串行执行，开始及创建前均重新获取远端主线与标签。同提交重试安全，冲突提交失败，禁止覆盖标签。远端引用创建负责最终并发仲裁。
- 浅克隆、脏工作区、错误分支或事件均拒绝打标。流程不创建 Release、不上传附件、不部署。

GitHub 的附注标签通过“创建标签对象 → 创建引用”两步完成，参考 [GitHub 官方文档](https://docs.github.com/en/rest/git/tags)。PR 查询使用 GitHub 的 `merged_at`、`per_page`，Gitea 使用 `merged`、`limit`，均完整读取分页。

## SDK 本地检查

本仓库当前为 Go SDK，不接入服务端可执行文件构建或数据库迁移版本。Go 模块版本由正式 Git 标签提供。要求 Node.js 24、Git 完整历史及 go.mod 指定的工具链。

```bash
node scripts/version.mjs
node --test scripts/version.test.mjs scripts/forge-api.test.mjs
go build ./...
go vet ./...
go test ./... -short
```

开发版示例为 `v0.2.0-dev.18+sha.3ce9a0d.dirty`；本地候选按提交推导，最终打标还会读取 PR 标签。SDK 消费者使用正式 Git 标签选择模块版本。

## 平台与复用

GitHub 入口为 `.github/workflows/version-tag.yml`；Gitea 入口为 `.gitea/workflows/version-tag.yml`。各自仅识别本平台目录。迁移其他仓库时复制 `scripts/` 中这五个版本相关文件、所需平台工作流和本文，即可；业务单测范围按项目调整。

GitHub 默认使用短期 `GITHUB_TOKEN`，只有打标作业申请 `contents: write` 与 `pull-requests: read`。Gitea 默认使用 `GITEA_TOKEN`，打标申请 `code: write` 与 `pull-requests: read`。保护规则要求专用机器人时可设置本仓库 `VERSION_BOT_TOKEN`；拒绝权限时流程失败，不尝试绕过规则。GitHub 默认令牌创建的标签通常不会触发另一轮工作流，本实现也不依赖标签事件。

脚本支持 `VERSION_PLATFORM=github|gitea`、`VERSION_BRANCH`（默认 `main`）、`VERSION_REMOTE`（默认 `origin`）。修改发布分支时，还须同步两个工作流中的分支条件、并发组及相关输入说明。平台地址、仓库、事件由工作流上下文注入。`GITHUB_API_URL` 可适配 GitHub Enterprise；Gitea 使用 `GITEA_SERVER_URL`。Gitea runner 须可获取标准 actions/checkout、setup-node、setup-go；内网镜像与代理按部署环境配置，不写入通用模板。

建议保护 `main` 和 `v*` 标签，并限制 `semver:*` 标签编辑权限。模板不会修改远端保护规则。本仓库 SDK 门禁使用 `go test ./... -short`。

## 启用检查

1. 将变更通过常规代码审查合入 `main`，确保仓库 Actions 已启用，runner 可用且具备相应令牌权限。
2. Go 工具链版本从当前 `go.mod` 读取，不在模板中硬编码；若模块要求的工具链无法下载，门禁会失败，修正工具链问题后重跑。
3. 首次无历史标签会分析全部历史提交，可能提出较大的初始版本；可通过合并 PR 的语义标签明确本次意图。MAJOR 始终需要精确 SHA 确认。
4. 本次迁移仅修改本地文件，未推送、未创建真实远端标签、Release 或部署。
