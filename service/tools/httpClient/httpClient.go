package httpClient

import (
	"V2RayA/extra/proxyWithHttp"
	"V2RayA/model/v2ray"
	"V2RayA/persistence/configure"
	"V2RayA/tools"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
)

func GetHttpClientWithProxy(proxyURL string) (client *http.Client, err error) {
	u, err := url.Parse(proxyURL)
	if err != nil {
		return
	}
	dialer, err := proxyWithHttp.FromURL(u, proxyWithHttp.Direct)
	if err != nil {
		return
	}
	httpTransport := &http.Transport{}
	httpTransport.Dial = dialer.Dial
	client = &http.Client{Transport: httpTransport}
	return
}

func GetHttpClientWithV2RayAProxy() (client *http.Client, err error) {
	host := "127.0.0.1"
	//是否在docker环境
	if tools.IsInDocker() {
		//连接网关，即宿主机的端口，失败则用同网络下v2ray容器的
		out, err := exec.Command("sh", "-c", "ip route list default|head -n 1|awk '{print $3}'").Output()
		if err == nil {
			host = strings.TrimSpace(string(out))
		} else {
			host = "v2ray"
		}
	}
	return GetHttpClientWithProxy("socks5://" + host + ":20170")
}

func GetHttpClientWithV2RayAPac() (client *http.Client, err error) {
	host := "127.0.0.1"
	//是否在docker环境
	if tools.IsInDocker() {
		//连接网关，即宿主机的端口，失败则用同网络下v2ray容器的
		out, err := exec.Command("sh", "-c", "ip route|grep default|awk '{print $3}'").Output()
		if err == nil {
			host = strings.TrimSpace(string(out))
		} else {
			host = "v2ray"
		}
	}
	return GetHttpClientWithProxy("http://" + host + ":20172")
}

func GetHttpClientAutomatically() (c *http.Client, err error) {
	if s := configure.GetSettingNotNil(); !v2ray.IsV2RayRunning() || configure.GetConnectedServer() == nil || s.Transparent != configure.TransparentClose {
		return http.DefaultClient, nil
	}
	setting := configure.GetSettingNotNil()
	switch setting.ProxyModeWhenSubscribe {
	case configure.ProxyModePac:
		c, err = GetHttpClientWithV2RayAPac()
	case configure.ProxyModeProxy:
		c, err = GetHttpClientWithV2RayAProxy()
	default:
		c = http.DefaultClient
	}
	return
}
