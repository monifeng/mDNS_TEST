# mDNS 网站测绘 CLI 示例程序实施计划

## Summary

本计划用于从空目录初始化一个 Golang 网站测绘 CLI 示例程序，程序能力为：输入 IP 网段与端口范围，基于 DNS-SD/mDNS 服务发现结果，筛选并输出目标范围内的 mDNS 协议资产信息；输出字段至少包含 `ip`、`port`、`host` 与具备深度识别能力的 `banner`，并提供可用于验证 banner 识别深度的数据集输出。

用户已确认的关键决策：

- 当前目录为空，需要规划完成并获批后继续初始化并实现可运行 CLI 示例程序。
- mDNS 发现核心方式采用 DNS-SD 发现。
- CLI 默认输出格式采用 JSON。
- README 与 TODO 需要作为项目交付内容的一部分生成。
- 输出能力需覆盖用户示例中的服务信息结构，例如 `_workstation._tcp.local`、`_http._tcp.local`、`_smb._tcp.local`、`_qdiscover._tcp.local`、`_device-info._tcp.local`、`_afpovertcp._tcp.local`，以及 TXT/banner 深度字段如 `path=/`、`accessType=https`、`model=TS-X64`、`displayModel=TS-464C`、`fwVer=...` 等。

## Current State Analysis

当前工作目录为：

- `d:\develop\golang\test1`

探索结果：

- 当前目录为空。
- 未发现 `go.mod`。
- 未发现 Go 源码文件。
- 未发现 `README.md` 或 `TODO.md`。
- 未发现现有 CLI、mDNS、资产建模、测试或构建脚本。

因此本次实现不需要兼容已有代码风格，但应建立一个清晰、可维护、适合后续扩展的 Go CLI 项目结构。

## Proposed Changes

### 1. 初始化 Go 模块

新增文件：

- `go.mod`

建议模块名：

- `mdnsmap`

计划依赖：

- `github.com/grandcat/zeroconf`

选择原因：

- 该库可直接执行 DNS-SD/mDNS 服务浏览，能获取服务实例、HostName、AddrIPv4、AddrIPv6、Port、Text、TTL 等资产识别所需字段。
- 相比手写 mDNS 报文解析，更适合作为示例程序快速交付，同时便于 README 解释与后续扩展。

### 2. 建立 CLI 入口

新增文件：

- `cmd/mdnsmap/main.go`

职责：

- 解析 CLI 参数。
- 调用扫描逻辑。
- 输出 JSON 结果。
- 返回合理退出码。

计划支持参数：

- `-cidr`：必填，目标 IP 网段，例如 `192.168.1.0/24`。
- `-ports`：必填，端口范围或端口列表，例如 `1-1024`、`80,443,5000`、`9,445,548,5000`。
- `-services`：可选，服务类型列表，默认内置常见 DNS-SD 服务类型。
- `-timeout`：可选，mDNS 浏览超时时间，默认 `5s`。
- `-pretty`：可选，是否格式化 JSON 输出。
- `-dataset`：可选，输出 banner 深度验证数据集视图。

默认服务类型建议：

- `_workstation._tcp`
- `_http._tcp`
- `_https._tcp`
- `_smb._tcp`
- `_qdiscover._tcp`
- `_device-info._tcp`
- `_afpovertcp._tcp`
- `_ssh._tcp`
- `_ipp._tcp`

说明：

- DNS-SD/mDNS 服务发现天然是以服务类型为入口，不是传统 TCP connect 扫描。
- `-cidr` 与 `-ports` 在本程序中用于对发现结果进行范围过滤，即只输出 IP 命中网段且端口命中指定范围的服务实例。

### 3. 实现参数解析与范围匹配

新增文件：

- `internal/config/config.go`
- `internal/netutil/cidr.go`
- `internal/netutil/ports.go`

职责：

`internal/config/config.go`：

- 定义 `Config` 结构。
- 统一存放 CLI 参数解析结果。
- 校验必填项与默认值。

`internal/netutil/cidr.go`：

- 解析 CIDR。
- 判断 IPv4 是否在网段内。
- 对 IPv6 链路本地地址保留输出，但不作为 `-cidr` IPv4 网段过滤的主依据。

`internal/netutil/ports.go`：

- 解析端口列表与端口范围。
- 支持 `80`、`80,443`、`1-1024`、`9,445,548,5000`。
- 校验端口范围为 `1-65535`。

### 4. 实现 mDNS/DNS-SD 服务发现

新增文件：

- `internal/discovery/mdns.go`

职责：

- 对配置中的多个服务类型执行 DNS-SD 浏览。
- 收集 `zeroconf.ServiceEntry` 数据。
- 提取字段：
  - service type
  - service name
  - domain
  - host name
  - port
  - IPv4
  - IPv6
  - TXT records
  - TTL
- 使用超时上下文控制扫描生命周期。
- 对重复服务实例去重。

计划去重键：

- `service_type + name + host + port + ip`

并发策略：

- 每个服务类型启动一个 goroutine 浏览。
- 使用 channel 聚合结果。
- 扫描结束后统一过滤与排序，保证输出稳定。

### 5. 实现资产模型与 banner 深度识别

新增文件：

- `internal/asset/model.go`
- `internal/asset/banner.go`

`internal/asset/model.go` 定义输出结构：

```go
type Asset struct {
    IP          string            `json:"ip"`
    Port        int               `json:"port"`
    Host        string            `json:"host"`
    Service     string            `json:"service"`
    Name        string            `json:"name"`
    IPv4        []string          `json:"ipv4,omitempty"`
    IPv6        []string          `json:"ipv6,omitempty"`
    Hostname    string            `json:"hostname"`
    TTL         uint32            `json:"ttl,omitempty"`
    TXT         []string          `json:"txt,omitempty"`
    Banner      Banner            `json:"banner"`
}
```

`Banner` 结构建议：

```go
type Banner struct {
    Raw         []string          `json:"raw,omitempty"`
    Fields      map[string]string `json:"fields,omitempty"`
    Summary     string            `json:"summary"`
    Fingerprint Fingerprint       `json:"fingerprint"`
    Depth       BannerDepth       `json:"depth"`
}
```

`Fingerprint` 字段建议：

- `vendor`
- `product`
- `model`
- `display_model`
- `firmware_version`
- `firmware_build`
- `access_type`
- `access_port`
- `path`
- `protocol_hint`

`BannerDepth` 字段建议：

- `level`：例如 `none`、`basic`、`txt_fields`、`fingerprinted`。
- `score`：0-100。
- `matched_fields`：命中的关键字段列表。
- `evidence`：用于说明识别深度的数据来源，例如 TXT 记录中的 `model`、`fwVer`、`displayModel`。

深度识别策略：

- 解析 TXT 中的 `key=value` 项。
- 对用户示例中的 QNAP/Qdiscover 字段做结构化提取：
  - `accessType`
  - `accessPort`
  - `model`
  - `displayModel`
  - `fwVer`
  - `fwBuildNum`
- 对 `_http._tcp` 中的 `path=/` 做 HTTP 服务路径识别。
- 对 `_device-info._tcp` 中的 `model=Xserve` 做设备模型识别。
- 对 `_smb._tcp`、`_afpovertcp._tcp`、`_workstation._tcp` 提供协议级 fingerprint。

### 6. 实现输出与数据集视图

新增文件：

- `internal/output/json.go`
- `internal/output/dataset.go`

`internal/output/json.go`：

- 默认输出 JSON 数组。
- 支持 `-pretty` 输出缩进 JSON。

`internal/output/dataset.go`：

- 当用户传入 `-dataset` 时，输出面向验证的数据集结构。
- 数据集用于验证 banner 识别深度。

建议数据集结构：

```json
{
  "dataset_type": "mdns_banner_depth_validation",
  "assets": [
    {
      "ip": "192.168.1.10",
      "port": 5000,
      "host": "slw-nas.local",
      "service": "_qdiscover._tcp.local",
      "banner_raw": ["accessType=https", "accessPort=86", "model=TS-X64"],
      "expected_depth": "fingerprinted",
      "identified_fields": {
        "accessType": "https",
        "accessPort": "86",
        "model": "TS-X64"
      },
      "score": 100
    }
  ]
}
```

### 7. 增加测试

新增文件：

- `internal/netutil/ports_test.go`
- `internal/netutil/cidr_test.go`
- `internal/asset/banner_test.go`

测试重点：

- 端口范围解析。
- CIDR 命中判断。
- TXT banner 解析。
- QNAP/Qdiscover 深度识别。
- HTTP path、device-info model 等基础 fingerprint。

由于真实 mDNS 依赖局域网环境，不将网络发现作为默认单元测试强依赖。发现层可通过构造中间结构或后续接口抽象做集成测试。

### 8. 生成 README

新增文件：

- `README.md`

README 内容结构：

1. 项目简介。
2. 功能特性。
3. 安装与构建。
4. CLI 参数说明。
5. 使用示例。
6. JSON 输出示例。
7. banner 深度识别说明。
8. 数据集输出说明。
9. mDNS/DNS-SD 工作方式与限制。
10. 常见问题。

README 中需要明确：

- 该工具基于 mDNS/DNS-SD 服务发现，而非全端口 TCP 主动连接扫描。
- `-cidr` 与 `-ports` 对发现到的服务结果进行过滤。
- mDNS 通常只在本地链路/局域网内有效。
- 需要目标设备开启对应 mDNS/DNS-SD 广播。

### 9. 生成 TODO

新增文件：

- `TODO.md`

TODO 内容结构：

- MVP 已完成项。
- 后续增强项。
- 识别规则扩展。
- 输出格式扩展。
- 性能与并发优化。
- 测试与数据集增强。
- 安全与合规注意事项。

建议后续增强：

- 支持读取服务类型字典文件。
- 支持 JSONL/CSV 输出。
- 支持更完整的厂商 fingerprint 规则。
- 支持被动监听 mDNS 回包。
- 支持 pcap 导出或导入离线解析。
- 支持生成识别质量报告。

## Assumptions & Decisions

1. 项目从空目录开始初始化，不需要兼容已有代码。
2. 默认使用 Go 标准库 `flag` 实现 CLI，避免引入额外 CLI 框架；mDNS 发现使用 `github.com/grandcat/zeroconf`。
3. 默认输出 JSON 数组，`-pretty` 控制可读格式。
4. 用户要求的“输入 IP 网段和端口范围”在 DNS-SD 发现模式中解释为“对发现结果做 IP 与端口过滤”。
5. banner 深度识别优先基于 mDNS TXT 记录、服务类型、主机名与端口进行结构化 fingerprint，不主动进行 HTTP/SMB/AFP 协议深探测，以保持示例程序聚焦 mDNS 协议资产发现。
6. IPv4 CIDR 作为主要过滤依据；IPv6 地址会保留在输出中，但如果输入是 IPv4 CIDR，不要求 IPv6 参与范围匹配。
7. 单元测试不依赖真实局域网 mDNS 环境，避免测试不稳定。

## Verification Steps

计划实现完成后执行以下验证：

1. 格式化代码：

```powershell
gofmt -w .
```

2. 运行单元测试：

```powershell
go test ./...
```

3. 构建 CLI：

```powershell
go build ./cmd/mdnsmap
```

4. 本地运行帮助命令：

```powershell
.\mdnsmap.exe -h
```

5. 示例运行：

```powershell
.\mdnsmap.exe -cidr 192.168.1.0/24 -ports 9,445,548,5000 -pretty
```

6. 数据集输出验证：

```powershell
.\mdnsmap.exe -cidr 192.168.1.0/24 -ports 9,445,548,5000 -dataset -pretty
```

验收标准：

- 能成功构建可执行文件。
- CLI 参数校验清晰，缺少必填参数时给出错误。
- JSON 输出至少包含 `ip`、`port`、`host`、`banner`。
- 对用户示例中的 TXT banner 字段能结构化解析。
- `go test ./...` 通过。
- README 与 TODO 能清楚说明使用方式、实现范围与后续计划。
