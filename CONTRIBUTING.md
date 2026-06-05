# 协作开发规范

这份文档用来约定 `Edu Admin` 后续的日常开发方式，目标只有两个：

- `main` 分支尽量保持稳定、随时可用
- 每次改动都有清楚来路，方便回看、协作和继续开发

## 1. 分支规则

默认只保留一个长期主分支：

- `main`：稳定分支，尽量保持可以直接启动和演示

日常开发一律从 `main` 拉新分支，不直接在 `main` 上改：

- `feat/...`：新功能
- `fix/...`：问题修复
- `docs/...`：文档调整
- `refactor/...`：重构整理
- `chore/...`：杂项维护
- `test/...`：测试补充

分支命名示例：

- `feat/student-attendance`
- `fix/class-list-empty-state`
- `docs/development-workflow`

约定：

- 一个分支只做一类事情
- 分支名尽量短，但要能看懂用途
- 做到一半发现需求变了，优先新开分支，不要把无关内容混在一起

## 2. 标准开发流程

每次开发尽量按下面流程走：

1. 先切回主分支并同步最新代码
2. 从 `main` 拉一个新分支
3. 在新分支里开发和自查
4. 提交到远程分支
5. 通过合并请求再合回 `main`

参考命令：

```bash
git checkout main
git pull origin main
git checkout -b feat/student-attendance
```

开发完成后：

```bash
git status
git add .
git commit -m "feat: add student attendance list"
git push -u origin feat/student-attendance
```

## 3. 提交说明规范

提交说明尽量采用统一前缀：

- `feat:` 新功能
- `fix:` 修复问题
- `docs:` 文档修改
- `refactor:` 结构整理
- `chore:` 工程维护
- `test:` 测试相关

示例：

- `feat: add class schedule list api`
- `fix: handle empty student list response`
- `docs: add development workflow guide`

约定：

- 一次提交尽量只表达一件事
- 不要使用看不出内容的说明，比如 `update`、`test`、`fix bug`
- 如果改动很多，优先拆成几次小提交

## 4. 合并规则

默认采用“分支开发，合并回主分支”的方式：

- 不直接把日常改动提交到 `main`
- 合并前至少自己过一遍改动内容
- 能跑的检查先跑一遍，再合并
- 默认优先使用 `Squash and merge`，让主分支历史更清楚

建议在合并前自查这些内容：

- 后端改动：跑 `make test`
- 前端改动：至少确认页面能打开，必要时跑 `npm run build`
- 接口改动：确认返回结构没有明显破坏已有页面
- 配置改动：确认 `.env.example` 和说明文档同步更新

## 5. 文档放置规则

为了避免文档混乱，后续统一这样放：

- 本地规划、草稿、讨论记录：放 `docs/`
- 对外公开、需要跟仓库一起版本管理的说明：放仓库根目录

常见例子：

- `README.md`：项目介绍、启动方式
- `CONTRIBUTING.md`：协作方式、提交流程
- `docs/`：产品规划、接口草案、结构草图等本地文档

补充说明：

- 当前项目把 `docs/` 作为本地规划区使用
- 所以要公开给协作者看的说明，尽量不要新增在 `docs/` 里

## 6. 以后和 AI 协作时的默认约定

后续如果继续让我帮你开发，默认按这个方式执行：

1. 先基于最新 `main` 开分支
2. 在分支里完成改动
3. 做基本自查
4. 说明改了什么
5. 再合回 `main`

这样做的好处是：

- 主分支更干净
- 每次改动更容易回看
- 出问题时更容易定位
- 后面如果有别人一起参与，也能直接接上流程

## 7. 当前阶段的简单落地建议

考虑到项目现在还在起步阶段，先执行这几条就够用：

- 所有开发都先开分支
- 所有改动都写清楚提交说明
- 所有合并都先自查再进 `main`
- 规划类文档继续本地保留
- 对外说明文档放根目录

等后面项目更稳定了，再补：

- 分支保护
- 自动检查
- 合并模板
- 发布版本记录
