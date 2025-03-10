<h1 align="center">
  <img src="https://github.com/BarryDeng/clash/raw/master/docs/logo.png" alt="Clash" width="200">
  <br>Clash<br>
</h1>

<h4 align="center">A rule-based tunnel in Go.</h4>

<p align="center">
  <a href="https://github.com/BarryDeng/clash/actions">
    <img src="https://img.shields.io/github/workflow/status/BarryDeng/clash/Go?style=flat-square" alt="Github Actions">
  </a>
  <a href="https://goreportcard.com/report/github.com/BarryDeng/clash">
    <img src="https://goreportcard.com/badge/github.com/BarryDeng/clash?style=flat-square">
  </a>
  <img src="https://img.shields.io/github/go-mod/go-version/BarryDeng/clash?style=flat-square">
  <a href="https://github.com/BarryDeng/clash/releases">
    <img src="https://img.shields.io/github/release/BarryDeng/clash/all.svg?style=flat-square">
  </a>
  <a href="https://github.com/BarryDeng/clash/releases/tag/premium">
    <img src="https://img.shields.io/badge/release-Premium-00b4f0?style=flat-square">
  </a>
</p>

## Features

- Local HTTP/HTTPS/SOCKS server with authentication support
- VMess, Shadowsocks, Trojan, Snell protocol support for remote connections
- Built-in DNS server that aims to minimize DNS pollution attack impact, supports DoH/DoT upstream and fake IP.
- Rules based off domains, GEOIP, IPCIDR or Process to forward packets to different nodes
- Remote groups allow users to implement powerful rules. Supports automatic fallback, load balancing or auto select node based off latency
- Remote providers, allowing users to get node lists remotely instead of hardcoding in config
- Netfilter TCP redirecting. Deploy Clash on your Internet gateway with `iptables`.
- Comprehensive HTTP RESTful API controller

## Premium Features

- TUN mode on macOS, Linux and Windows. [Doc](https://github.com/BarryDeng/clash/wiki/premium-core-features#tun-device)
- Match your tunnel by [Script](https://github.com/BarryDeng/clash/wiki/premium-core-features#script)
- [Rule Provider](https://github.com/BarryDeng/clash/wiki/premium-core-features#rule-providers)

## Getting Started
Documentations are now moved to [GitHub Wiki](https://github.com/BarryDeng/clash/wiki).

## Premium Release
[Release](https://github.com/BarryDeng/clash/releases/tag/premium)

## Development
If you want to build an application that uses clash as a library, check out the the [GitHub Wiki](https://github.com/BarryDeng/clash/wiki/use-clash-as-a-library)

## Credits

* [riobard/go-shadowsocks2](https://github.com/riobard/go-shadowsocks2)
* [v2ray/v2ray-core](https://github.com/v2ray/v2ray-core)
* [WireGuard/wireguard-go](https://github.com/WireGuard/wireguard-go)

## License

This software is released under the GPL-3.0 license.

[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FBarryDeng%2Fclash.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FBarryDeng%2Fclash?ref=badge_large)
