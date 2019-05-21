package localsharelib

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
)

func getIps() string {
	var buf bytes.Buffer
	first := true

	ips, _ := net.InterfaceAddrs()
	for _, addr := range ips {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip.To4() != nil && ip.String() != "127.0.0.1" {
			if first {
				first = false
			} else {
				buf.WriteByte(',')
			}

			buf.WriteString(ip.String())
		}
	}

	return buf.String()
}

func getFirstIp() string {
	ips, _ := net.InterfaceAddrs()
	for _, addr := range ips {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip.To4() != nil && ip.String() != "127.0.0.1" {
			return ip.String()
		}
	}
	return ""
}

func (instance *LocalshareInstance) GetServerURL() string {
	return getFirstIp() + ":" + strconv.Itoa(instance.port) + "/api/files"
}

func sendjson(w http.ResponseWriter, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("internal server error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(b)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("internal server error"))
	}
}

func hash(v interface{}) string {
	sha_256 := sha256.New()
	sha_256.Write([]byte(fmt.Sprintf("%v", v)))
	return fmt.Sprintf("%x", sha_256.Sum(nil))
}

func firstOrEmpty(strs []string) string {
	if len(strs) >= 1 {
		return strs[0]
	}

	return ""
}
