# mDNS资产测绘工具 (mdnsmap)

基于 Go 开发的 mDNS/DNS-SD 服务发现工具，支持 IP 网段与端口过滤，并提供 banner 深度识别功能和验证数据集输出。

## 功能特性

- 本地 mDNS/DNS-SD 服务浏览
- IP 网段范围过滤 (CIDR)
- 端口范围或列表过滤
- 内置常见服务类型探测
- banner 深度识别与评分
- 支持 JSON 输出与格式化
- 支持数据集模式输出，用于验证识别深度

## 安装与构建

```bash
git clone https://example.com/mdnsmap
cd mdnsmap
go mod tidy
go build -o mdnsmap ./cmd/mdnsmap
```

## CLI 参数说明

| 参数 | 说明 | 默认值 | 必填 |
|------|------|--------|------|
| `-cidr` | 目标 IP 网段，例如 192.168.1.0/24 | - | 是 |
| `-ports` | 端口范围或列表，例如 1-1024 或 9,445,548,5000 | - | 是 |
| `-services` | 服务类型列表，默认内置多种 | - | 否 |
| `-timeout` | mDNS 浏览超时 | 5s | 否 |
| `-format` | 输出格式，支持 json、services | json | 否 |
| `-pretty` | 是否格式化输出 JSON | false | 否 |
| `-dataset` | 是否输出数据集格式 | false | 否 |

## 使用示例

### 基本使用 (JSON 格式)

```bash
./mdnsmap -cidr 192.168.1.0/24 -ports 9,445,548,5000 -pretty
```

### services 格式输出

```bash
./mdnsmap -cidr 192.168.1.0/24 -ports 9,445,548,5000 -format services
```

输出示例：

```
services:
9/tcp workstation:
Name=slw-nas
IPv4=192.168.1.10
Hostname=slw-nas.local.
TTL=10

5000/tcp http:
Name=slw-nas
IPv4=192.168.1.10
Hostname=slw-nas.local.
TTL=10
path=/

445/tcp smb:
Name=slw-nas
IPv4=192.168.1.10
Hostname=slw-nas.local.
TTL=10

5000/tcp qdiscover:
Name=slw-nas
IPv4=192.168.1.10
Hostname=slw-nas.local.
TTL=10
accessType=https,accessPort=86,model=TS-X64,displayModel=TS-464C,fwVer=5.2.9,fwBuildNum=20260214

answers:
PTR:
_workstation._tcp.local.
_http._tcp.local.
_smb._tcp.local.
_qdiscover._tcp.local.
```

### 数据集验证模式

```bash
./mdnsmap -cidr 192.168.1.0/24 -ports 9,445,548,5000 -dataset -pretty
```

## JSON 输出示例

```json
[
  {
    "ip": "192.168.1.10",
    "port": 5000,
    "host": "slw-nas.local.",
    "service": "_qdiscover._tcp.local.",
    "name": "slw-nas",
    "ipv4": [
      "192.168.1.10"
    ],
    "hostname": "slw-nas.local.",
    "ttl": 10,
    "txt": [
      "accessType=https",
      "accessPort=86",
      "model=TS-X64",
      "displayModel=TS-464C",
      "fwVer=5.2.9",
      "fwBuildNum=20260214"
    ],
    "banner": {
      "raw": [
        "accessType=https",
        "accessPort=86",
        "model=TS-X64",
        "displayModel=TS-464C",
        "fwVer=5.2.9",
        "fwBuildNum=20260214"
      ],
      "fields": {
        "accessType": "https",
        "accessPort": "86",
        "model": "TS-X64",
        "displayModel": "TS-464C",
        "fwVer": "5.2.9",
        "fwBuildNum": "20260214"
      },
      "summary": "qnap-qdiscover / TS-X64 / QNAP",
      "fingerprint": {
        "vendor": "QNAP",
        "model": "TS-X64",
        "display_model": "TS-464C",
        "firmware_version": "5.2.9",
        "firmware_build": "20260214",
        "access_type": "https",
        "access_port": "86",
        "protocol_hint": "qnap-qdiscover"
      },
      "depth": {
        "level": "fingerprinted",
        "score": 100,
        "matched_fields": [
          "accessType",
          "accessPort",
          "model",
          "displayModel",
          "fwVer",
          "fwBuildNum"
        ],
        "evidence": [
          "TXT:accessType=https",
          "TXT:accessPort=86",
          "TXT:model=TS-X64",
          "TXT:displayModel=TS-464C",
          "TXT:fwVer=5.2.9",
          "TXT:fwBuildNum=20260214"
        ]
      }
    }
  }
]
```

## banner 深度识别说明

- **none**: 无任何 banner 信息
- **basic**: 仅有基本 TXT 记录
- **txt_fields**: 解析出结构化字段
- **fingerprinted**: 识别到厂商/型号/固件等指纹信息

## mDNS/DNS-SD 工作方式与限制

该工具基于 mDNS/DNS-SD 服务浏览，而非全端口 TCP 主动连接扫描：
- 仅能探测本地网络内开启 mDNS 广播的设备
- `-cidr` 与 `-ports` 为结果过滤条件，不会主动探测不在结果范围内的设备
- 默认仅对 IPv4 地址进行网段过滤，IPv6 会保留在输出中

## 常见问题

**Q: 扫描不到任何服务？**

请确保目标设备开启了 mDNS/DNS-SD 广播，且在同一本地网络内。

**Q: 如何自定义服务类型？**

使用 `-services` 参数，例如：`-services _http._tcp,_smb._tcp`
