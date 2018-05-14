package localsharelib

import "bytes"
import "encoding/json"
import "net"
import "net/http"

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
