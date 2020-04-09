package api

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/heiha/ssr2clashr/config"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// ClashRSSR struct
type ClashRSSR struct {
	Group         string `yaml:"group"`
	Name          string `yaml:"name"`
	Type          string `yaml:"type"`
	Server        string `yaml:"server"`
	Port          int    `yaml:"port"`
	Password      string `yaml:"password"`
	Cipher        string `yaml:"cipher"`
	Protocol      string `yaml:"protocol"`
	ProtocolParam string `yaml:"protocolparam"`
	OBFS          string `yaml:"obfs"`
	OBFSParam     string `yaml:"obfsparam"`
	UDP           bool   `yaml:"udp"`
}

// Clash struct
type Clash struct {
	Port               int                               `yaml:"port"`
	SocksPort          int                               `yaml:"socks-port"`
	RedirPort          int                               `yaml:"redir-port"`
	AllowLan           bool                              `yaml:"allow-lan"`
	BindAddress        string                            `yaml:"bind-address"`
	Mode               string                            `yaml:"mode"`
	LogLevel           string                            `yaml:"log-level"`
	ExternalController string                            `yaml:"external-controller"`
	ExternalUI         string                            `yaml:"external-ui"`
	Secret             string                            `yaml:"secret"`
	Experimental       map[string]interface{}            `yaml:"experimental"`
	Authentication     []string                          `yaml:"authentication"`
	HOSTS              map[string]string                 `yaml:"hosts"`
	DNS                map[string]interface{}            `yaml:"dns"`
	CFWByPass          []string                          `yaml:"cfw-bypass"`
	CFWLatencyTimeout  int                               `yaml:"cfw-latency-timeout"`
	Proxies            []map[string]interface{}          `yaml:"proxies"`
	ProxyProviders     map[string]map[string]interface{} `yaml:"proxy-providers"`
	ProxyGroups        []map[string]interface{}          `yaml:"proxy-groups"`
	Rules              []string                          `yaml:"rules"`
}

// ProxyGroup struct
type ProxyGroup struct {
	Name     string   `yaml:"name"`
	Type     string   `yaml:"type"`
	Proxies  []string `yaml:"proxies"`
	URL      string   `yaml:"url"`
	Interval int      `yaml:"interval"`
}

// ProxyGroupSelect struct
type ProxyGroupSelect struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
	Use     []string `yaml:"use"`
}

const (
	// SSRServer int
	SSRServer = iota
	// SSRPort int
	SSRPort
	// SSRProtocol int
	SSRProtocol
	// SSRCipher int
	SSRCipher
	// SSROBFS int
	SSROBFS
	// SSRSuffix int
	SSRSuffix
)

// GroupRules []string
var GroupRules []string

func validURL(url string) bool {
	if url != "" && strings.Index(url, "http") == 0 {
		return true
	}
	return false
}

func validByte(v []byte) bool {
	if len(v) > 0 && v != nil {
		return true
	}
	return false
}

func validFile(v string) bool {
	if _, err := os.Stat(v); err == nil {
		return true
	}
	return false
}

func includeRemarks(ssr ClashRSSR, t interface{}) bool {
	if t == nil || reflect.ValueOf(t).IsNil() {
		return false
	}
	switch reflect.TypeOf(t).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(t)
		for i := 0; i < s.Len(); i++ {
			if val := s.Index(i).Interface(); val != nil {
				for _, v := range strings.Split(val.(string), "|") {
					if strings.Contains(ssr.Name, v) {
						return true
					}
				}
			}
		}
	}
	return false
}

func unicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("(?i)\\\\u[0-9a-zA-Z]+")
	//提取emoji数据表达式
	reg := regexp.MustCompile("(?i)\\\\u")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

func base64DecodeStripped(s string) ([]byte, error) {
	if i := len(s) % 4; i != 0 {
		s += strings.Repeat("=", 4-i)
	}
	s = strings.ReplaceAll(s, " ", "+")
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(s)
	}
	return decoded, err
}

func getRuleVarURL(suf string, url string) []string {

	resp, err1 := http.Get(url)
	if err1 == nil && resp.StatusCode == http.StatusOK {
		buf, err2 := ioutil.ReadAll(resp.Body)
		if err2 == nil {
			defer resp.Body.Close()
			if viper.GetBool("debug") {
				fmt.Println("下载完成：", url)
			}
			return getRuleVarAsset(suf, buf)
		}
		fmt.Println(err2)
	}
	fmt.Println(err1)
	defer resp.Body.Close()
	return nil
}

func getRuleVarFile(suf string, path string) []string {
	buf, err := ioutil.ReadFile(path)
	if err == nil {
		if viper.GetBool("debug") {
			fmt.Println("读取文件完成：", path)
		}
		return getRuleVarAsset(suf, buf)
	}
	fmt.Println(err)
	return nil
}

func getRuleVarAsset(suf string, s []byte) []string {
	var t []string
	scanner := bufio.NewScanner(bytes.NewReader(s))
	for scanner.Scan() {
		if len(scanner.Text()) < 1 || strings.HasPrefix(scanner.Text(), "#") || strings.HasPrefix(scanner.Text(), "USER-AGENT") || strings.HasPrefix(scanner.Text(), "PROCESS-NAME") || strings.HasPrefix(scanner.Text(), "URL-REGEX") {
		} else {
			if strings.HasSuffix(scanner.Text(), "no-resolve") || strings.HasSuffix(scanner.Text(), "force-remote-dns") {
				n := strings.LastIndex(scanner.Text(), ",")
				t = append(t, fmt.Sprintf("%s,%s%s", scanner.Bytes()[:n], suf, scanner.Bytes()[n:]))
			} else {
				t = append(t, fmt.Sprintf("%s,%s", scanner.Text(), suf))
			}
		}
	}
	return t
}

func getTemplateBuffer() []byte {
	path := viper.GetString("template")
	if len(path) > 1 {
		if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
			resp, err := http.Get(path)
			if err == nil && resp.StatusCode == http.StatusOK {
				buf, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					defer resp.Body.Close()
					return buf
				}
			}
			defer resp.Body.Close()
			return nil
		} else if validFile(path) {
			buf, err := ioutil.ReadFile(path)
			if err == nil {
				return buf
			}
			return nil
		} else if tbytes, _ := config.Asset(path); tbytes != nil {
			return tbytes
		}
	}
	return nil
}

func loadTemplate(protos []ClashRSSR) Clash {

	clash := Clash{}
	err := yaml.Unmarshal(getTemplateBuffer(), &clash)
	if err != nil {
		fmt.Println(err)
	}

	//clash.Proxy
	clash.Proxies = nil
	var proxys []map[string]interface{}
	var proxies []string
	for _, proto := range protos {
		proxy := make(map[string]interface{})
		j, _ := yaml.Marshal(proto)
		yaml.Unmarshal(j, &proxy)
		proxys = append(proxys, proxy)
		clash.Proxies = append(clash.Proxies, proxy)
		proxies = append(proxies, proto.Name)
	}

	proxySet := viper.GetStringMap("proxy")
	proxyEnabled := proxySet["enabled"]
	proxyset := proxySet["proxy"]
	if proxysetValue := reflect.ValueOf(proxyset); proxyEnabled != nil && reflect.TypeOf(proxyEnabled).Kind() == reflect.Bool && reflect.ValueOf(proxyEnabled).Bool() && proxyset != nil && reflect.TypeOf(proxyset).Kind() == reflect.Slice && !proxysetValue.IsNil() {
		for i := 0; i < proxysetValue.Len(); i++ {
			val := proxysetValue.Index(i).Interface()
			proxy := make(map[string]interface{})
			j, _ := yaml.Marshal(val)
			yaml.Unmarshal(j, &proxy)
			proxys = append(proxys, proxy)
			proxies = append(proxies, proxy["name"].(string))
		}
	}

	clash.Proxies = proxys

	// clash.ProxyGroup
	groupSet := viper.GetStringMap("groupset")
	groupEnabled := groupSet["enabled"]
	groupset := groupSet["groupset"]

	if groupsetValue := reflect.ValueOf(groupset); groupEnabled != nil && reflect.TypeOf(groupEnabled).Kind() == reflect.Bool && reflect.ValueOf(groupEnabled).Bool() && groupset != nil && reflect.TypeOf(groupset).Kind() == reflect.Slice && !groupsetValue.IsNil() {
		clash.ProxyGroups = nil
		for i := 0; i < groupsetValue.Len(); i++ {

			var gname, gtype, gurl string
			var gproxies []string
			var ginterval int

			if val := groupsetValue.Index(i).Interface(); val != nil {
				groupsetSplit := strings.Split(val.(string), "`")
				gname = strings.TrimSpace(groupsetSplit[0])
				gtype = strings.TrimSpace(groupsetSplit[1])
				switch gtype {
				case "select":
					for _, v := range groupsetSplit[2:] {
						if strings.HasPrefix(v, "[]") {
							gproxies = append(gproxies, strings.TrimPrefix(v, "[]"))
						} else if v == ".*" {
							gproxies = append(gproxies, proxies...)
						}
					}
					tmpGroup := ProxyGroupSelect{
						gname,
						gtype,
						gproxies,
						[]string{},
					}
					tmpf := make(map[string]interface{})
					tmps, _ := yaml.Marshal(tmpGroup)
					yaml.Unmarshal(tmps, &tmpf)
					clash.ProxyGroups = append(clash.ProxyGroups, tmpf)

				default:
					for _, v := range groupsetSplit[2:] {
						if prefix := "[]"; strings.HasPrefix(v, prefix) {
							gproxies = append(gproxies, strings.TrimPrefix(v, prefix))
						} else if v == ".*" {
							gproxies = append(gproxies, proxies...)
						} else if prefix := "[URL]"; strings.HasPrefix(v, prefix) {
							gurl = strings.TrimPrefix(v, prefix)
						} else if prefix := "[INR]"; strings.HasPrefix(v, prefix) {
							ginterval, _ = strconv.Atoi(strings.TrimPrefix(v, prefix))
						}
					}
					tmpGroup := ProxyGroup{
						gname,
						gtype,
						gproxies,
						gurl,
						ginterval,
					}
					tmpf := make(map[string]interface{})
					tmps, _ := yaml.Marshal(tmpGroup)
					yaml.Unmarshal(tmps, &tmpf)
					clash.ProxyGroups = append(clash.ProxyGroups, tmpf)
				}
			}

		}
	} else {
		if len(clash.ProxyGroups) < 1 {
			tmpProxyGroup := ProxyGroupSelect{
				"Proxy",
				"select",
				proxies,
				[]string{},
			}
			tmpf := make(map[string]interface{})
			tmps, _ := yaml.Marshal(tmpProxyGroup)
			yaml.Unmarshal(tmps, &tmpf)
			clash.ProxyGroups = append(clash.ProxyGroups, tmpf)

		} else {
			for _, group := range clash.ProxyGroups {
				groupProxies := group["proxies"].([]interface{})
				for i, proxie := range groupProxies {
					if "1" == proxie {
						groupProxies = groupProxies[:i]
						var tmpGroupProxies []string
						for _, s := range groupProxies {
							tmpGroupProxies = append(tmpGroupProxies, s.(string))
						}
						tmpGroupProxies = append(tmpGroupProxies, proxies...)
						group["proxies"] = tmpGroupProxies
						break
					}
				}

			}
		}
	}

	// clash.Rule
	ruleSet := viper.GetStringMap("ruleset")
	ruledEnabled := ruleSet["enabled"]
	ruleset := ruleSet["ruleset"]
	if rulesetValue := reflect.ValueOf(ruleset); ruledEnabled != nil && reflect.TypeOf(ruledEnabled).Kind() == reflect.Bool && reflect.ValueOf(ruledEnabled).Bool() && ruleset != nil && reflect.TypeOf(ruleset).Kind() == reflect.Slice && !rulesetValue.IsNil() {
		clash.Rules = GroupRules
	} else if len(clash.Rules) < 1 {
		fmt.Println("Rules Nothing")
	}

	return clash
}

func ssr2clashR(body []byte) []byte {
	scanner := bufio.NewScanner(bytes.NewReader(body))
	var ssrs []ClashRSSR
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "ssr://") {
			continue
		}
		s := scanner.Text()[6:]
		s = strings.TrimSpace(s)
		rawSSRConfig, err := base64DecodeStripped(s)
		if err != nil {
			continue
		}
		params := strings.Split(string(rawSSRConfig), `:`)
		if 6 != len(params) {
			continue
		}
		ssr := ClashRSSR{}
		ssr.Type = "ssr"
		ssr.Server = params[SSRServer]
		ssr.Port, _ = strconv.Atoi(params[SSRPort])
		ssr.Protocol = params[SSRProtocol]
		ssr.Cipher = params[SSRCipher]
		ssr.OBFS = params[SSROBFS]

		// 如果兼容ss协议，就转换为clash的ss配置
		// https://github.com/Dreamacro/clash
		if "origin" == ssr.Protocol && "plain" == ssr.OBFS {
			switch ssr.Cipher {
			case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm",
				"aes-128-cfb", "aes-192-cfb", "aes-256-cfb",
				"aes-128-ctr", "aes-192-ctr", "aes-256-ctr",
				"rc4-md5", "chacha20", "chacha20-ietf", "xchacha20",
				"chacha20-ietf-poly1305", "xchacha20-ietf-poly1305":
				ssr.Type = "ss"
			}
		}

		suffix := strings.Split(params[SSRSuffix], "/?")
		if 2 != len(suffix) {
			continue
		}
		passwordBase64 := suffix[0]
		password, err := base64DecodeStripped(passwordBase64)
		if err != nil {
			continue
		}
		ssr.Password = string(password)

		m, err := url.ParseQuery(suffix[1])
		if err != nil {
			continue
		}

		for k, v := range m {
			de, err := base64DecodeStripped(v[0])
			if err != nil {
				continue
			}
			switch k {
			case "obfsparam":
				ssr.OBFSParam = string(de)
				continue
			case "protoparam":
				ssr.ProtocolParam = string(de)
				continue
			case "remarks":
				ssr.Name = strings.TrimSpace(string(de))
				continue
			case "group":
				ssr.Group = string(de)
				continue
			}
		}
		ssr.UDP = false

		node := viper.GetStringMap("node")
		if includeRemarks(ssr, node["include_remarks"]) || !includeRemarks(ssr, node["exclude_remarks"]) {
			ssrs = append(ssrs, ssr)
		}

	}

	d, err := yaml.Marshal(loadTemplate(ssrs))
	if err != nil {
		return nil
	}

	return []byte(unicodeEmojiDecode(string(d)))
}

func getSub(url string) []byte {
	var body []byte
	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return body
	}
	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return body
	}
	return body
}

// Execute func()
func Execute(url string) []byte {
	var body = []byte("未获取到订阅内容！")
	url = strings.ToLower(strings.TrimSpace(url))

	if !validURL(url) {
		url = viper.GetString("url")
		if !validURL(url) {
			return body
		}
	}

	b := getSub(url)

	if !validByte(b) {
		return body
	}

	b, err := base64DecodeStripped(string(b))
	if err != nil || !strings.HasPrefix(string(b), "ssr://") {
		return body
	}

	b = ssr2clashR(b)
	if !validByte(b) {
		return body
	}

	return b
}

// InitRules func
func InitRules() {
	var grules []string
	ruleSet := viper.GetStringMap("ruleset")
	ruledEnabled := ruleSet["enabled"]
	ruleset := ruleSet["ruleset"]
	if rulesetValue := reflect.ValueOf(ruleset); ruledEnabled != nil && reflect.TypeOf(ruledEnabled).Kind() == reflect.Bool && reflect.ValueOf(ruledEnabled).Bool() && ruleset != nil && reflect.TypeOf(ruleset).Kind() == reflect.Slice && !rulesetValue.IsNil() {
		for i := 0; i < rulesetValue.Len(); i++ {
			if val := rulesetValue.Index(i).Interface(); val != nil {
				rulesetSplit := strings.SplitN(val.(string), ",", 2)
				rname := rulesetSplit[0]
				rtmp := rulesetSplit[1]
				if len(rulesetSplit) > 1 {
					if prefix := "[]"; strings.HasPrefix(rtmp, prefix) {
						grules = append(grules, fmt.Sprintf("%s,%s", strings.TrimPrefix(rtmp, prefix), rname))
					} else if strings.HasPrefix(rtmp, "https://") || strings.HasPrefix(rtmp, "http://") {
						grules = append(grules, getRuleVarURL(rname, rtmp)...)
					} else if validFile(rtmp) {
						grules = append(grules, getRuleVarFile(rname, rtmp)...)
					} else if rbytes, _ := config.Asset(rtmp); rbytes != nil {
						grules = append(grules, getRuleVarAsset(rname, rbytes)...)
					}
				}
			}
		}
		GroupRules = grules
	}
}
