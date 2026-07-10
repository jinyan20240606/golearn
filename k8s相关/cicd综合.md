# `.cicd` CI/CD 流水线学习文档

> 本文档是对项目 `.cicd` 目录的完整分析，帮助你理解 CI/CD 实现逻辑，掌握后续独立部署能力。

- [`.cicd` CI/CD 流水线学习文档](#cicd-cicd-流水线学习文档)
  - [一、整体架构](#一整体架构)
    - [目录结构](#目录结构)
  - [二、CI 入口文件](#二ci-入口文件)
  - [三、公共定义](#三公共定义)
    - [3.1 变量定义](#31-变量定义)
    - [3.2 流水线阶段](#32-流水线阶段)
    - [3.3 核心 Shell 函数（YAML 锚点）](#33-核心-shell-函数yaml-锚点)
      - [ImageBuild — 构建 Docker 镜像](#imagebuild--构建-docker-镜像)
      - [ImagePushToHub — 推送镜像到 Harbor](#imagepushtohub--推送镜像到-harbor)
      - [ApplySecret — 将配置注入 K8s Secret](#applysecret--将配置注入-k8s-secret)
      - [Deploy — 完整的 Helm 部署流程](#deploy--完整的-helm-部署流程)
        - [helm相关概念](#helm相关概念)
    - [3.4 before\_script](#34-before_script)
    - [3.5 GitLab CI 解析 before\_script 的完整流程](#35-gitlab-ci-解析-before_script-的完整流程)
      - [阶段一：YAML 解析（GitLab Server 端）](#阶段一yaml-解析gitlab-server-端)
      - [关键要点](#关键要点)
  - [四、Job 定义](#四job-定义)
    - [Job 一览表](#job-一览表)
    - [代码编译](#代码编译)
    - [构建镜像](#构建镜像)
    - [部署测试环境](#部署测试环境)
    - [全量生产集群](#全量生产集群)
  - [五、Dockerfile](#五dockerfile)
  - [六、Helm Chart](#六helm-chart)
    - [6.1 Chart 元信息](#61-chart-元信息)
    - [6.2 辅助函数](#62-辅助函数)
    - [6.3 Deployment 模板](#63-deployment-模板)
    - [6.4 Service 模板](#64-service-模板)
  - [七、环境 Values](#七环境-values)
    - [测试环境](#测试环境)
    - [生产环境](#生产环境)
  - [八、完整部署流程](#八完整部署流程)
  - [九、架构设计要点总结](#九架构设计要点总结)
    - [9.1 配置与镜像分离](#91-配置与镜像分离)
    - [9.2 分支即环境](#92-分支即环境)
    - [9.3 YAML 锚点复用](#93-yaml-锚点复用)
    - [9.4 制品传递](#94-制品传递)
    - [9.5 安全门禁](#95-安全门禁)
    - [9.6 Helm values 分层](#96-helm-values-分层)
  - [十、要自己搭建这套能力，你需要准备](#十要自己搭建这套能力你需要准备)
    - [需要在 CI 中配置的变量](#需要在-ci-中配置的变量)

---

## 一、整体架构

这是一个 **GitLab CI + Docker + Helm + Kubernetes** 的完整 CI/CD 流水线，用于将 Go 服务自动化构建、打包镜像并部署到 K8s 集群。

### 目录结构

```
.cicd/
├── gitlab-ci.yaml              # CI 入口，引用子文件
├── gitlab-ci/
│   ├── common.yaml             # 公共变量、Shell 函数、阶段定义
│   └── build.yaml              # 各阶段 job 定义
├── dockerfile/
│   └── Dockerfile              # 构建生产镜像
├── helm/
│   ├── Chart.yaml              # Helm Chart 元信息
│   └── templates/
│       ├── _helpers.tpl        # Helm 模板辅助函数
│       ├── deployment.yaml     # K8s Deployment 模板
│       └── service.yaml        # K8s Service 模板
└── values/
    ├── test.yaml               # 测试环境 values
    └── prod.yaml               # 生产环境 values
```

---

## 二、CI 入口文件

**文件**: [`.cicd/gitlab-ci.yaml`](.cicd/gitlab-ci.yaml)

```yaml
include:
  - '.cicd/gitlab-ci/common.yaml'
  - '.cicd/gitlab-ci/build.yaml'
```

通过 `include` 将公共定义和 job 定义拆分为两个子文件，职责分离，便于维护。

---

## 三、公共定义

**文件**: [`.cicd/gitlab-ci/common.yaml`](.cicd/gitlab-ci/common.yaml)

### 3.1 变量定义

```yaml
variables:
  branch: ${CI_COMMIT_REF_NAME}          # 当前 Git 分支名
  version: ${CI_COMMIT_SHORT_SHA}        # Git commit 短 SHA
  image_name: harbor.qihoo.net/.../rpa-task:${CI_COMMIT_REF_NAME}-${CI_COMMIT_SHORT_SHA}
```

- `branch` — 用于 Helm release 命名，实现分支即环境
- `version` — 作为 Docker 镜像 tag 的一部分，确保每次构建唯一
- `image_name` — 完整镜像地址，格式 `{仓库}/{项目名}:{分支名}-{短SHA}`

### 3.2 流水线阶段

```yaml
stages:
  - build          # 编译 Go 二进制
  - imagebuild     # 构建 Docker 镜像并推送
  - deployTest     # 部署测试环境
  - deployProd     # 部署生产环境
```

### 3.3 核心 Shell 函数（YAML 锚点）

#### ImageBuild — 构建 Docker 镜像

```bash
function ImageBuild(){
    local dockerFile="${CI_PROJECT_DIR}/.cicd/dockerfile/Dockerfile"
    docker login harbor.qihoo.net -u $HABOR_USER -p $HABOR_PWD
    docker build -t ${image_name} -f ${dockerFile} ${CI_PROJECT_DIR}
}
```

#### ImagePushToHub — 推送镜像到 Harbor

```bash
function ImagePushToHub(){
    docker login harbor.qihoo.net -u $HABOR_USER -p $HABOR_PWD
    docker push ${image_name}
}
```

#### ApplySecret — 将配置注入 K8s Secret

```bash
function ApplySecret(){
    local env="${1}"
```
定义函数 `ApplySecret`。Shell 函数**没有形参列表**，参数通过位置变量 `$1`、`$2`... 在函数体内部获取。这里 `$1` 由调用方传入——`Deploy` 函数内调用 `ApplySecret ${env}` 时传入，而 `Deploy` 的 `$1` 来自 job script 中写的 `Deploy test` 或 `Deploy prod`。

```bash
    local kubectlConfFile="${K8S_SECRET_CONF_TEST}"
    local configYaml="${CONFIG_YAML_TEST}"
    local databaseYaml="${DATABASE_YAML_TEST}"
    local logYaml="${LOG_YAML_TEST}"
```
先默认赋值为**测试环境**的配置。`K8S_SECRET_CONF_TEST` 等是 GitLab CI 变量（在 CI/CD Settings 中配置），内容是测试环境 K8s 的 kubeconfig 文件和配置 YAML 内容。

```bash
    if [ "${env}" == "prod" ]; then
        kubectlConfFile="${K8S_SECRET_CONF_PROD}"
        configYaml="${CONFIG_YAML_PROD}"
        databaseYaml="${DATABASE_YAML_PROD}"
        logYaml="${LOG_YAML_PROD}"
    fi
```
如果传入的是 `prod`，则切换为**生产环境**的配置。

```bash
    local secretName="rpa-task-config-${env}"
```
Secret 名称 = `rpa-task-config-test` 或 `rpa-task-config-prod`，与环境绑定。

```bash
    /usr/local/bin/kubectl -n plat-arch --kubeconfig=${kubectlConfFile} create secret generic ${secretName} \
        --from-literal=config.yaml="${configYaml}" \
        --from-literal=database.yaml="${databaseYaml}" \
        --from-literal=log.yaml="${logYaml}" \
        --dry-run=client -o yaml \
        | /usr/local/bin/kubectl -n plat-arch --kubeconfig=${kubectlConfFile} apply -f -
}
```

这条核心命令逐段拆解：

| 部分 | 说明 |
|------|------|
| `/usr/local/bin/kubectl` | Runner 容器内预装的 kubectl |
| `-n plat-arch` | 指定 K8s 命名空间为 `plat-arch` |
| `--kubeconfig=${kubectlConfFile}` | 指定 K8s 连接凭证（区分测试/生产集群） |
| `create secret generic ${secretName}` | 创建名为 `rpa-task-config-{env}` 的 Secret |
| `--from-literal=config.yaml="${configYaml}"` | 将 `config.yaml` 的**文件内容**直接作为 Secret 的一个 key-value 存入。`${configYaml}` 是 CI 变量，值是整个 YAML 文件内容的字符串 |
| `--dry-run=client -o yaml` | **只生成 YAML 清单，不真正请求 K8s API**，结果输出到 stdout |
| `\| kubectl ... apply -f -` | 管道：上一步输出的 YAML 通过 stdin 传给 `apply -f -`，执行创建或更新 |

- --dry-run=client：客户端本地模拟执行，不发送任何请求到 K8s APIServer，不会真正创建资源；
- -o yaml：把即将生成的 Secret 资源输出为 YAML 格式，打印到标准输出 stdout。
- |：管道，将前一段输出的 Secret YAML 作为后一段命令的输入
- apply -f -：
  - -f - 代表「从标准输入读取资源文件」；
  - kubectl apply 特性：幂等操作
  - 集群不存在该 Secret → 新建；
  - 集群已存在同名 Secret → 对比内容，差异则更新覆盖

#### Deploy — 完整的 Helm 部署流程

```bash
function Deploy(){
    local env="${1}"
    # 根据环境选择 K8s 配置文件
    local kubectlConfFile="${K8S_SECRET_CONF_TEST}"
    if [ "${env}" == "prod" ]; then
        kubectlConfFile="${K8S_SECRET_CONF_PROD}"
    fi

    local valueFile=".cicd/values/${env}.yaml"

    # 1. 先注入配置 Secret
    ApplySecret ${env}

    # 2. Helm 安装/升级
    /usr/local/bin/helm -n plat-arch --kubeconfig=${kubectlConfFile} upgrade -i \
        --set-string branch=${branch},version=${version},env=${env} \
        -f ${valueFile} rpa-task-${branch} .cicd/helm

    # 3. 等待 Pod 就绪
    /usr/local/bin/kubectl -n plat-arch --kubeconfig=${kubectlConfFile} wait \
        --timeout=300s --for=condition=ready pod -l app.kubernetes.io/name=rpa-task-${branch}
}
```

- --kubeconfig=${kubectlConfFile}：使用指定集群凭证，多集群发布必备
- upgrade -i （--install的简写）
  - 集群不存在该 Release：执行安装（install）
  - 集群已存在同名 Release：执行滚动升级（upgrade）具备幂等性，流水线可重复执行。
  - --set-string branch=${branch},version=${version},env=${env}
    - 动态注入字符串类型变量到 Helm Chart，模板内可直接读取 .Values.branch、.Values.version、.Values.env。
    - --set-string 强制识别为字符串，避免数字版本号被 helm 转成数字类型。
  - -f ${valueFile}
    - 加载自定义 values yaml 文件，覆盖 Chart 默认配置。
  - rpa-task-${branch}
    - Release 名称，按分支区分，例如 rpa-task-main、rpa-task-feature，多分支环境隔离。
  - .cicd/helm
    - Helm Chart 本地目录路径，存放 Chart 模板。
- kubectl wait 等待 Pod 就绪
  - 阻塞脚本，直到所有匹配标签的 Pod 全部就绪，或超时 300 秒直接失败，保证发布后服务真正可用才往下走。

##### helm相关概念
- Chart：安装包（类似 npm 包、docker 镜像、安装程序模板），只存模板、配置、资源定义，是静态文件；
- Release：把 Chart 真正部署到集群后，生成的运行实例，包含：Deployment / Pod / Service / ConfigMap / Secret 等全套 k8s 资源，还有本次部署的版本、自定义参数、历史记录。

### 3.4 before_script

```yaml
before_script:
  - CI_COMMIT_REF_NAME=$(echo ${CI_COMMIT_REF_NAME//_/-})
  - *auto_devops_build
  - *auto_devops_deploy
```

`before_script` 中的**每一项都是 shell 命令**。`*auto_devops_build` 和 `*auto_devops_deploy` 是 YAML 锚点引用，但**锚点只是 YAML 层面的文本替换机制**，展开后仍然是 shell 命令。

| 行 | 本质 | 说明 |
|----|------|------|
| `CI_COMMIT_REF_NAME=$(echo ...)` | **Shell 命令** | 将分支名中的 `_` 替换为 `-`，避免 K8s/Helm 命名不合法 |
| `*auto_devops_build` | **YAML 锚点 → 展开为 Shell 命令** | `&auto_devops_build` 定义的是一段 shell 脚本字符串，`*` 引用后原地展开，变成 `function ImageBuild(){...}` 和 `function ImagePushToHub(){...}` 两条函数定义命令 |
| `*auto_devops_deploy` | **YAML 锚点 → 展开为 Shell 命令** | 同理，展开为 `function ApplySecret(){...}` 和 `function Deploy(){...}` |

YAML 锚点机制：`&` 定义锚点，`*` 引用锚点。展开后的效果等价于：

```yaml
before_script:
  - CI_COMMIT_REF_NAME=$(echo ${CI_COMMIT_REF_NAME//_/-})   # shell 命令1
  - |                                                        # shell 命令2（多行）
    function ImageBuild(){ ... }
    function ImagePushToHub(){ ... }
  - |                                                        # shell 命令3（多行）
    function ApplySecret(){ ... }
    function Deploy(){ ... }
```

所以 `before_script` 整体做的事情是：**先修改变量，再定义 4 个 Shell 函数**，这样后续 job 的 `script` 中才能直接调用 `ImageBuild`、`ImagePushToHub`、`ApplySecret`、`Deploy` 这些函数。

### 3.5 GitLab CI 解析 before_script 的完整流程

#### 阶段一：YAML 解析（GitLab Server 端）

当 push 代码后，GitLab Server 读取 `.gitlab-ci.yml` 并进行解析：

**Step 1: 合并 include**

`gitlab-ci.yaml` 中的 `include` 指令告诉 GitLab Server 将 common.yaml 和 build.yaml 的内容合并到主文件中，形成一个完整的 YAML 文档。

**Step 2: 解析 YAML 锚点**

`&auto_devops_build` 定义的代码块被 `*auto_devops_build` 引用处原地展开：

```yaml
# 展开前 (common.yaml)
.auto_devops_build: &auto_devops_build |   # 点号开头表示这是一个隐藏 job，不作为独立 job 执行
  function ImageBuild(){
    local dockerFile="${CI_PROJECT_DIR}/.cicd/dockerfile/Dockerfile"
    docker login harbor.qihoo.net -u $HABOR_USER -p $HABOR_PWD
    docker build -t ${image_name} -f ${dockerFile} ${CI_PROJECT_DIR}
  }
  function ImagePushToHub(){
    docker login harbor.qihoo.net -u $HABOR_USER -p $HABOR_PWD
    docker push ${image_name}
  }

.auto_devops_deploy: &auto_devops_deploy |   # 同理，隐藏 job + 锚点定义
  function ApplySecret(){...}
  function Deploy(){...}

before_script:
  - CI_COMMIT_REF_NAME=$(echo ${CI_COMMIT_REF_NAME//_/-})
  - *auto_devops_build    # 引用锚点，原地展开
  - *auto_devops_deploy   # 引用锚点，原地展开
```

**Step 3: 合并到每个 job 形成最终脚本**

GitLab CI 的 `before_script` 是**继承机制**——每个 job 的最终执行脚本 = `before_script` + `script`。Runner 拉取到的解析结果大致如下：

```bash
#!/bin/bash
set -e  # Runner 默认设置，任何命令失败立即退出 job

# === before_script === (Runner 拼接)
CI_COMMIT_REF_NAME=$(echo ${CI_COMMIT_REF_NAME//_/-})  # 修改变量名

function ImageBuild() { ... }      # 定义函数1
function ImagePushToHub() { ... }  # 定义函数2
function ApplySecret() { ... }     # 定义函数3
function Deploy() { ... }          # 定义函数4

# === script === (Runner 拼接，以"构建镜像" job 为例)
ImageBuild
ImagePushToHub
```

#### 关键要点

| 要点 | 说明 |
|------|------|
| **YAML 锚点解析时机** | GitLab Server 端，YAML 解析阶段完成，Runner 拿到的是展开后的结果 |
| **before_script 执行时机** | 每个 job 的容器启动后、script 执行前 |
| **函数作用域** | 函数在 before_script 中定义，与 script 在**同一个 shell 进程**内，所以 script 可以直接调用 |
| **每个 job 独立** | 每个 job 启动独立的容器，before_script 在每个 job 中各自执行一次 |
| **set -e** | Runner 默认 `set -e`，任何一行命令失败都会导致整个 job 立即退出 |

---

## 四、Job 定义

**文件**: [`.cicd/gitlab-ci/build.yaml`](.cicd/gitlab-ci/build.yaml)

### Job 一览表

| Job | Stage | 触发分支 | 说明 |
|-----|-------|---------|------|
| `代码编译` | `build` | 所有分支 | `go build` 编译，产物存为 artifacts |
| `构建镜像` | `imagebuild` | `test-1`, `master` | 依赖编译产物，构建并推送 Docker 镜像 |
| `部署测试环境` | `deployTest` | `test-1` | 自动部署到测试 K8s |
| `全量生产集群` | `deployProd` | `master` | **手动触发**，部署到生产 |

### 代码编译

```yaml
代码编译:
  stage: build
  retry: 2
  image: harbor.qihoo.net/.../golang:1.25.5
  script:
    - go mod tidy
    - CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/rpa-task cmd/main.go
  artifacts:
    paths:
      - bin/
    expire_in: 1 hour
```

- 使用 Go 官方镜像编译
- `CGO_ENABLED=0` 禁用 CGO，生成静态链接二进制
- `-ldflags="-w -s"` 去掉调试信息，减小体积
- 产物通过 `artifacts` 传递给后续 job，避免重复编译
- CI 任务执行结束后，把构建产物打包保存到 GitLab，供后续阶段下载、页面手动下载查看
- 指定要打包上传的目录：项目下 bin/ 整个文件夹；
构建脚本编译出的二进制程序、可执行文件都放在 bin 目录，流水线结束后作为产物归档；
不写 paths 不会保存任何文件
- 产物过期自动清理时长：仅保存 1 小时。超过 1 小时 GitLab 自动删除这份构建产物，节省存储。

### 构建镜像

```yaml
构建镜像:
  stage: imagebuild
  needs:
    - 代码编译
  image: harbor.qihoo.net/.../cicd-helm:1.0.5
  script:
    - ImageBuild
    - ImagePushToHub
  only:
    refs:
      - test-1
      - master
```

- `needs` 显式依赖编译 job
- 使用自定义的 `cicd-helm` 镜像（内含 docker + kubectl + helm 工具）

### 部署测试环境

```yaml
部署测试环境:
  stage: deployTest
  image: harbor.qihoo.net/.../cicd-helm:1.0.5
  script:
    - Deploy test
  environment:
    name: pre/$CI_COMMIT_REF_SLUG
  only:
    refs:
      - test-1
```

### 全量生产集群

```yaml
全量生产集群:
  stage: deployProd
  image: harbor.qihoo.net/.../cicd-helm:1.0.5
  script:
    - Deploy prod
  environment:
    name: release
  only:
    refs:
      - master
  when: manual
```

- `when: manual` — 生产部署需要人工确认，增加安全门禁

---

## 五、Dockerfile

**文件**: [`.cicd/dockerfile/Dockerfile`](.cicd/dockerfile/Dockerfile)

```dockerfile
FROM harbor.qihoo.net/commercial-platform-plat-arch/centos-sec-76-1810:1.0.1

WORKDIR /app/

COPY bin/rpa-task /app/bin/rpa-task
COPY config/dev/ /app/config/dev/

ENV TIME_ZONE Asia/Shanghai
RUN /bin/ln -sf /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime \
    && chmod +x /app/bin/rpa-task

ENTRYPOINT ["/app/bin/rpa-task", "api"]
```

- 基于 CentOS 安全基础镜像
- 只复制编译好的二进制和配置目录（配置会被 K8s Secret 覆盖）
- 设置时区为 Asia/Shanghai
- 入口为 `rpa-task api` 命令

---

## 六、Helm Chart

### 6.1 Chart 元信息

**文件**: [`.cicd/helm/Chart.yaml`](.cicd/helm/Chart.yaml)

```yaml
apiVersion: v2
name: rpa-task
description: rpa-task project
type: application
version: 0.1.0
appVersion: "0.1.0"
```
```yaml
# 说明
1. apiVersion: v2
Chart 规范版本，v2 代表 Helm 3 标准格式；
Helm2 使用 v1，现在集群全部用 Helm3，统一写 v2。
2. name: rpa-task
Chart 包名称，当前这套模板叫 rpa-task，对应你的 RPA 任务服务。
3. description: rpa-task project
Chart 描述，简单说明这个 Chart 用途。
4. type: application
Chart 分两种类型：
application：业务应用（后端服务、RPA 任务服务，你当前场景）
library：通用模板库，只供其他 Chart 引用，不能独立部署
这里代表这是一套可直接安装部署的业务服务模板。
5. version: 0.1.0
Chart 模板自身版本号
模板文件（deployment/service/values.yaml）修改时升级这个号；
只代表 Chart 资源定义的版本，和业务代码镜像版本无关。
6. appVersion: "0.1.0"
业务应用镜像版本
对应你项目代码打包的 Docker 镜像 tag，模板里一般会通过 .Chart.AppVersion 读取用来填充镜像地址。
和 version 区分：
version：Chart 模板版本
appVersion：业务程序镜像版本
7. fullname: rpa-task
自定义完整名称，Helm 渲染资源时会以此为基础拼接资源名、标签；
不自定义的话 Helm 会默认拼接 Release名-chart名，这里固定前缀为 rpa-task。
简单区分 version /appVersion
version：你改 k8s yaml 模板才升级
appVersion：你改代码、打包新镜像才升级
```

### 6.2 辅助函数

**文件**: [`.cicd/helm/templates/_helpers.tpl`](.cicd/helm/templates/_helpers.tpl)

- Helm 标准命名工具模板，定义两个可全局复用的命名函数：
  - rpa-task.name：生成基础短名称，用于标签、selector、小标识
  - rpa-task.fullname，用来统一生成资源名称，规避命名超长、分隔符多余横线问题

```yaml
{{- define "rpa-task.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{- define "rpa-task.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Chart.Name .Values.branch | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
```

> **关键设计**：`fullname` = `rpa-task-{branch}`。这意味着**不同分支可以共存于同一个 K8s 命名空间**，互不干扰。

### 6.3 Deployment 模板

**文件**: [`.cicd/helm/templates/deployment.yaml`](.cicd/helm/templates/deployment.yaml)

```yaml
# K8s API资源版本，Deployment稳定标准版本
apiVersion: apps/v1
# 资源类型：无状态应用部署控制器
kind: Deployment
# 资源元数据
metadata:
  # Deployment名称，调用helm模板生成 rpa-task-${branch}
  name: {{ include "rpa-task.fullname" . }}
# Deployment核心调度规格
spec:
  # 副本数量，从values.yaml读取配置
  replicas: {{ .Values.replicaCount }}
  # Pod标签选择器，用于关联管理下方模板创建的Pod
  selector:
    matchLabels:
      # 匹配Pod上同名完整应用标签
      app.kubernetes.io/name: {{ include "rpa-task.fullname" . }}
  # Pod模板，定义新建Pod的标准配置
  template:
    metadata:
      # Pod自身标签，和selector匹配，让Deployment管控该Pod
      labels:
        app.kubernetes.io/name: {{ include "rpa-task.fullname" . }}
    # Pod内部运行规格
    spec:
      # 容器列表
      containers:
        # 第一个容器名称
        - name: rpa-task
          # 镜像地址:镜像版本，仓库地址+流水线传入的version
          image: "{{ .Values.image.repository }}:{{ .Values.version }}"
          # 镜像拉取策略 Always/IfNotPresent/Never
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          # 容器暴露端口配置
          ports:
            # 容器内部监听端口，从values读取
            - containerPort: {{ .Values.service.targetPort }}
          # 容器内挂载数据卷配置
          volumeMounts:
            # 关联下方volumes定义的卷名 config
            - name: config
              # 容器内挂载路径，程序读取配置文件目录
              mountPath: /app/config/dev
      # Pod可用数据卷定义
      volumes:
        # 卷名称，和volumeMounts.name对应
        - name: config
          # 卷来源为Secret密钥资源
          secret:
            # 使用的Secret名称，按环境区分 rpa-task-config-dev/test/prod
            secretName: rpa-task-config-{{ .Values.env }}
```

- 镜像 tag 由 `version` 变量动态注入
- 配置通过 Secret 挂载到 `/app/config/dev`，覆盖镜像内的默认配置
- {{ include "rpa-task.fullname" . }}
  - 调用 _helpers.tpl 模板，自动拼接 rpa-task-${branch}，实现多分支资源隔离；
- Secret 挂载逻辑
  - 流水线脚本 ApplySecret 创建名为 rpa-task-config-${env} 的 Secret，内部存 config.yaml/database.yaml/log.yaml；
  - 容器内 /app/config/dev 会自动解压出三份配置文件给程序读取；
- 关联你发布命令
  - kubectl wait -l app.kubernetes.io/name=rpa-task-${branch} 匹配的就是此处 Pod 标签。

### 6.4 Service 模板

**文件**: [`.cicd/helm/templates/service.yaml`](.cicd/helm/templates/service.yaml)

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ include "rpa-task.fullname" . }}
spec:
  selector:
    app.kubernetes.io/name: {{ include "rpa-task.fullname" . }}
  ports:
    - protocol: TCP
      port: {{ .Values.service.port }}
      targetPort: {{ .Values.service.targetPort }}
```

标准 K8s Service，将 Pod 的 `targetPort`（8080）暴露为 `port`（80）。

---

## 七、环境 Values

### 测试环境

**文件**: [`.cicd/values/test.yaml`](.cicd/values/test.yaml)

```yaml
branch: test-1
version: latest
env: test
replicaCount: 1
image:
  repository: harbor.qihoo.net/commercial-platform-plat-arch/rpa-task
  pullPolicy: IfNotPresent
service:
  port: 80
  targetPort: 8080
```

### 生产环境

**文件**: [`.cicd/values/prod.yaml`](.cicd/values/prod.yaml)

```yaml
branch: master
version: latest
env: prod
replicaCount: 2
image:
  repository: harbor.qihoo.net/commercial-platform-plat-arch/rpa-task
  pullPolicy: IfNotPresent
service:
  port: 80
  targetPort: 8080
```

| 字段 | test.yaml | prod.yaml |
|------|-----------|-----------|
| `branch` | `test-1` | `master` |
| `env` | `test` | `prod` |
| `replicaCount` | 1 | 2 |
| `service.port` | 80 | 80 |
| `service.targetPort` | 8080 | 8080 |

生产环境 2 副本保证高可用，测试环境 1 副本节省资源。

---

## 八、完整部署流程

```
开发者 push 代码到 GitLab
    │
    ▼
┌─────────────────────────────────────────────────────────┐
│ GitLab CI 自动触发流水线                                  │
│                                                         │
│  Stage 1: build                                          │
│    └─ go mod tidy → go build → 产出 bin/rpa-task         │
│       └─ 产物保存为 artifacts，传递给后续 job              │
│                                                         │
│  Stage 2: imagebuild (仅 test-1 / master 分支)           │
│    └─ docker login → docker build → docker push          │
│       └─ 镜像推送到 Harbor 仓库                           │
│                                                         │
│  Stage 3: deployTest (仅 test-1 分支, 自动触发)           │
│    └─ ① ApplySecret: 创建/更新 K8s Secret (配置注入)      │
│    └─ ② helm upgrade -i: 部署到测试 K8s                  │
│    └─ ③ kubectl wait: 等待 Pod ready                     │
│                                                         │
│  Stage 4: deployProd (仅 master 分支, 需手动触发)          │
│    └─ ① ApplySecret: 创建/更新 K8s Secret (生产配置)      │
│    └─ ② helm upgrade -i: 部署到生产 K8s                  │
│    └─ ③ kubectl wait: 等待 Pod ready                     │
└─────────────────────────────────────────────────────────┘
```

---

## 九、架构设计要点总结

### 9.1 配置与镜像分离

配置通过 K8s Secret 注入，而非打包进镜像。不同环境（测试/生产）使用同一镜像，只需更换 Secret 中的配置内容。

### 9.2 分支即环境

Helm release 名 = `rpa-task-{branch}`，多分支可以共存于同一个 K8s 命名空间，互不干扰。方便同时部署多个 feature 分支进行测试。

### 9.3 YAML 锚点复用

用 `&auto_devops_build` / `&auto_devops_deploy` 将公共 Shell 函数抽取为 YAML 锚点，在 `before_script` 中引用，避免重复代码。

### 9.4 制品传递

编译产物通过 GitLab CI 的 `artifacts` 机制传递给后续 stage，避免重复编译，加快流水线速度。

### 9.5 安全门禁

生产部署使用 `when: manual`，需要人工在 GitLab CI 界面点击确认才能执行，增加安全控制。

### 9.6 Helm values 分层

按环境拆分 values 文件（`test.yaml` / `prod.yaml`），清晰管理不同环境的差异配置。

---

## 十、要自己搭建这套能力，你需要准备

| 组件 | 用途 | 可选方案 |
|------|------|---------|
| **CI Runner** | 执行流水线任务 | GitLab Runner / GitHub Actions Runner / Jenkins |
| **镜像仓库** | 存储 Docker 镜像 | Harbor / Docker Hub / AWS ECR |
| **K8s 集群** | 运行服务 | 自建 K8s / 云厂商 K8s（ACK / TKE / EKS） |
| **Helm** | K8s 包管理 | 安装到 K8s 集群即可 |
| **CI 变量** | 存储敏感信息 | 在 GitLab CI/CD Settings 中配置 |

### 需要在 CI 中配置的变量

| 变量名 | 说明 |
|--------|------|
| `HABOR_USER` | Harbor 镜像仓库用户名 |
| `HABOR_PWD` | Harbor 镜像仓库密码 |
| `K8S_SECRET_CONF_TEST` | 测试环境 K8s kubeconfig 文件内容 |
| `K8S_SECRET_CONF_PROD` | 生产环境 K8s kubeconfig 文件内容 |
| `CONFIG_YAML_TEST` | 测试环境 config.yaml 内容 |
| `DATABASE_YAML_TEST` | 测试环境 database.yaml 内容 |
| `LOG_YAML_TEST` | 测试环境 log.yaml 内容 |
| `CONFIG_YAML_PROD` | 生产环境 config.yaml 内容 |
| `DATABASE_YAML_PROD` | 生产环境 database.yaml 内容 |
| `LOG_YAML_PROD` | 生产环境 log.yaml 内容 |
