package executor

import (
	"fmt"
	"os"
	"sync"

	"github.com/BarryDeng/clash/adapter"
	"github.com/BarryDeng/clash/adapter/outboundgroup"
	"github.com/BarryDeng/clash/component/auth"
	"github.com/BarryDeng/clash/component/dialer"
	"github.com/BarryDeng/clash/component/iface"
	"github.com/BarryDeng/clash/component/profile"
	"github.com/BarryDeng/clash/component/profile/cachefile"
	"github.com/BarryDeng/clash/component/resolver"
	"github.com/BarryDeng/clash/component/trie"
	"github.com/BarryDeng/clash/config"
	C "github.com/BarryDeng/clash/constant"
	"github.com/BarryDeng/clash/constant/provider"
	"github.com/BarryDeng/clash/dns"
	P "github.com/BarryDeng/clash/listener"
	authStore "github.com/BarryDeng/clash/listener/auth"
	"github.com/BarryDeng/clash/log"
	"github.com/BarryDeng/clash/tunnel"
)

var mux sync.Mutex

func readConfig(path string) ([]byte, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("configuration file %s is empty", path)
	}

	return data, err
}

// Parse config with default config path
func Parse() (*config.Config, error) {
	return ParseWithPath(C.Path.Config())
}

// ParseWithPath parse config with custom config path
func ParseWithPath(path string) (*config.Config, error) {
	buf, err := readConfig(path)
	if err != nil {
		return nil, err
	}

	return ParseWithBytes(buf)
}

// ParseWithBytes config with buffer
func ParseWithBytes(buf []byte) (*config.Config, error) {
	return config.Parse(buf)
}

// ApplyConfig dispatch configure to all parts
func ApplyConfig(cfg *config.Config, force bool) {
	mux.Lock()
	defer mux.Unlock()

	updateUsers(cfg.Users)
	updateProxies(cfg.Proxies, cfg.Providers)
	updateRules(cfg.Rules)
	updateHosts(cfg.Hosts)
	updateProfile(cfg)
	updateGeneral(cfg.General, force)
	updateDNS(cfg.DNS)
	updateExperimental(cfg)
}

func GetGeneral() *config.General {
	ports := P.GetPorts()
	authenticator := []string{}
	if auth := authStore.Authenticator(); auth != nil {
		authenticator = auth.Users()
	}

	general := &config.General{
		Inbound: config.Inbound{
			Port:           ports.Port,
			SocksPort:      ports.SocksPort,
			RedirPort:      ports.RedirPort,
			TProxyPort:     ports.TProxyPort,
			MixedPort:      ports.MixedPort,
			Authentication: authenticator,
			AllowLan:       P.AllowLan(),
			BindAddress:    P.BindAddress(),
		},
		Mode:     tunnel.Mode(),
		LogLevel: log.Level(),
		IPv6:     !resolver.DisableIPv6,
	}

	return general
}

func updateExperimental(c *config.Config) {}

func updateDNS(c *config.DNS) {
	if !c.Enable {
		resolver.DefaultResolver = nil
		resolver.DefaultHostMapper = nil
		dns.ReCreateServer("", nil, nil)
		return
	}

	cfg := dns.Config{
		Main:         c.NameServer,
		Fallback:     c.Fallback,
		IPv6:         c.IPv6,
		EnhancedMode: c.EnhancedMode,
		Pool:         c.FakeIPRange,
		Hosts:        c.Hosts,
		FallbackFilter: dns.FallbackFilter{
			GeoIP:     c.FallbackFilter.GeoIP,
			GeoIPCode: c.FallbackFilter.GeoIPCode,
			IPCIDR:    c.FallbackFilter.IPCIDR,
			Domain:    c.FallbackFilter.Domain,
		},
		Default: c.DefaultNameserver,
		Policy:  c.NameServerPolicy,
	}

	r := dns.NewResolver(cfg)
	m := dns.NewEnhancer(cfg)

	// reuse cache of old host mapper
	if old := resolver.DefaultHostMapper; old != nil {
		m.PatchFrom(old.(*dns.ResolverEnhancer))
	}

	resolver.DefaultResolver = r
	resolver.DefaultHostMapper = m

	dns.ReCreateServer(c.Listen, r, m)
}

func updateHosts(tree *trie.DomainTrie) {
	resolver.DefaultHosts = tree
}

func updateProxies(proxies map[string]C.Proxy, providers map[string]provider.ProxyProvider) {
	tunnel.UpdateProxies(proxies, providers)
}

func updateRules(rules []C.Rule) {
	tunnel.UpdateRules(rules)
}

func updateGeneral(general *config.General, force bool) {
	log.SetLevel(general.LogLevel)
	tunnel.SetMode(general.Mode)
	resolver.DisableIPv6 = !general.IPv6

	dialer.DefaultInterface.Store(general.Interface)
	dialer.DefaultRoutingMark.Store(int32(general.RoutingMark))

	iface.FlushCache()

	if !force {
		return
	}

	allowLan := general.AllowLan
	P.SetAllowLan(allowLan)

	bindAddress := general.BindAddress
	P.SetBindAddress(bindAddress)

	tcpIn := tunnel.TCPIn()
	udpIn := tunnel.UDPIn()

	P.ReCreateHTTP(general.Port, tcpIn)
	P.ReCreateSocks(general.SocksPort, tcpIn, udpIn)
	P.ReCreateRedir(general.RedirPort, tcpIn, udpIn)
	P.ReCreateTProxy(general.TProxyPort, tcpIn, udpIn)
	P.ReCreateMixed(general.MixedPort, tcpIn, udpIn)
}

func updateUsers(users []auth.AuthUser) {
	authenticator := auth.NewAuthenticator(users)
	authStore.SetAuthenticator(authenticator)
	if authenticator != nil {
		log.Infoln("Authentication of local server updated")
	}
}

func updateProfile(cfg *config.Config) {
	profileCfg := cfg.Profile

	profile.StoreSelected.Store(profileCfg.StoreSelected)
	if profileCfg.StoreSelected {
		patchSelectGroup(cfg.Proxies)
	}
}

func patchSelectGroup(proxies map[string]C.Proxy) {
	mapping := cachefile.Cache().SelectedMap()
	if mapping == nil {
		return
	}

	for name, proxy := range proxies {
		outbound, ok := proxy.(*adapter.Proxy)
		if !ok {
			continue
		}

		selector, ok := outbound.ProxyAdapter.(*outboundgroup.Selector)
		if !ok {
			continue
		}

		selected, exist := mapping[name]
		if !exist {
			continue
		}

		selector.Set(selected)
	}
}
