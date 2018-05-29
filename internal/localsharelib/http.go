package localsharelib

import "github.com/go-chi/chi"
import "io"
import "net"
import "net/http"
import "strconv"

func (instance *LocalshareInstance) startHttpServer() error {
	port, err := getFreePort()
	if err != nil {
		return err
	}
	router := instance.routes()
	go http.ListenAndServe(":"+strconv.Itoa(port), router)
	instance.port = port
	return nil
}

func getFreePort() (int, error) {
	ln, err := net.Listen("tcp", ":")
	if err != nil {
		return 0, err
	}

	defer ln.Close()
	return ln.Addr().(*net.TCPAddr).Port, nil
}

func (instance *LocalshareInstance) routes() http.Handler {
	r := chi.NewRouter()
	r.Get("/api/files/{name}", instance.serveFileRoute)
	r.Get("/api/files", instance.filesRoute)
	r.Get("/api/ping", instance.pingRoute)
	return r
}

func (instance *LocalshareInstance) filesRoute(w http.ResponseWriter, r *http.Request) {
	data := []map[string]string{}
	for _, file := range instance.files {
		data = append(data, map[string]string{"name": file.Name(), "size": strconv.FormatInt(file.Size(), 10)})
	}
	sendjson(w, data)
}

func (instance *LocalshareInstance) serveFileRoute(w http.ResponseWriter, r *http.Request) {
	reader, err := instance.files[chi.URLParam(r, "name")].Open()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("internal server error"))
	}
	io.Copy(w, reader)
}

func (instance *LocalshareInstance) pingRoute(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}
