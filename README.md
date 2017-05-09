# 通讯上位机服务

## 通讯协议

本系统通讯协议采用：`HJ 212-2017`

`HJ 212-2017` 协议注释：

污染物在线监控(监测)系统数据传输标准(Data transmission standard for online monitoring systems of pollutant)

协议起草单位：`西安交大长天软件股份有限公司`、`环境保护部信息中心`、`中国环境监测总站`。

此协议由`环境保护部` 2017 年 4 月 24 日批准，自 2017 年 5 月 1 日起实施。

### 通讯协议数据结构

1.通讯包结构组成

| 名称 | 类型 | 长度 | 描述 |
| :-: | :-: | :-: | :-- |
| 包头 | 字符 | 2 | 固定为## |
| 数据段长度 | 十进制整数 | 4 | 数据段的 ASCII 字符数，例如:长 255，则写为“0255” |
| 数据段 | 字符 | 0≤n≤1024 | 变长的数据，详见《数据段结构组成表》 |
| CRC 校验 | 十六进制整数 | 4 | 数据段的校验结果，CRC 校验算法见附录 A。接收到一条命令，如 果 CRC 错误，执行结束 |
| 包尾 | 字符 | 2 | 固定为<CR><LF>(回车、换行) |

2.数据段结构组成

表中 “长度” 包含字段名称、‘=’、字段内容三部分内容。

| 名称 | 类型 | 长度 | 描述 |
| :-: | :-: | :-: | :-- |
| 请求编码 QN | 字符 | 20 | 精确到毫秒的时间戳:QN=YYYYMMDDhhmmsszzz，用来唯一标识一次 命令交互 |
| 系统编码 ST | 字符 | 5 | ST=系统编码, 系统编码取值详见《系统编码表》 |
| 命令编码 CN | 字符 | 7 | CN=命令编码, 命令编码取值详见《命令编码表》 |
| 设备唯一标识 MN | 字符 | 27 | MN=设备唯一标识，这个标识固化在设备中，用于唯一标识一个设备。 MN 由 EPC-96 编码转化的字符串组成，即 MN 由 24 个 0~9，A~F 的字 符组成 |

