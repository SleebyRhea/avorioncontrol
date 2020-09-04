package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// https://stackoverflow.com/questions/43601359/how-do-i-serve-css-and-js-in-go
// Am thief. Credit to @RayfenWindspear :D
func serveHTTP(h *ConnHub, gs GameServer, out chan []byte) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		LogOutput(gs, "Received connection to /admin")
		t := template.Must(template.ParseFiles("templates/admin.html"))
		data := GameStatus(gs)

		if err := t.Execute(w, data); err != nil {
			log.Output(1, err.Error())
			LogHTTP(gs, 500, r)
		}

		LogHTTP(gs, 200, r)
	})

	http.HandleFunc("/api/ajax/fullstatus/", func(w http.ResponseWriter, r *http.Request) {
		LogInfo(gs, "Received fullstatus request: "+r.RequestURI)
		json, err := json.Marshal(GameStatus(gs))
		if err != nil {
			LogError(gs, err.Error())
			LogHTTP(gs, 500, r)
			return
		}

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
		LogHTTP(gs, 200, r)
	})

	http.HandleFunc("/api/player/kick/", func(w http.ResponseWriter, r *http.Request) {
		LogInfo(gs, "Received kick request: "+r.RequestURI)
		pn := strings.TrimPrefix(r.RequestURI, "/api/player/kick/")
		rc := 403

		if plr := gs.Player(pn); plr != nil {
			plr.Kick("Kicked by the internet")
			rc = 200
		} else {
			rc = 404
		}

		w.WriteHeader(rc)
		LogHTTP(gs, rc, r)
	})

	http.HandleFunc("/api/player/ban/", func(w http.ResponseWriter, r *http.Request) {
		pn := strings.TrimPrefix(r.RequestURI, "/api/player/ban/")
		var (
			rc  = 403
			msg = "Banned from the internet"
		)

		if plr := gs.Player(pn); plr != nil {
			rc = 200
			plr.Ban(msg)
		} else {
			rc = 404
		}

		LogHTTP(gs, rc, r)
		w.WriteHeader(rc)
	})

	http.HandleFunc("/api/server/password/", func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(r.RequestURI)
		p := strings.TrimPrefix(u.Path, "/api/server/password")
		p = strings.TrimPrefix(p, "/")

		if p == "" {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(200)
			w.Write([]byte(gs.Password()))
			SendCommand("password "+p, gs)
		}

		LogHTTP(gs, 200, r)
	})

	http.HandleFunc("/api/server/start/", func(w http.ResponseWriter, r *http.Request) {
		if gs.IsUp() {
			LogHTTP(gs, 403, r)
			w.WriteHeader(403)
			return
		}

		if err := gs.Start(); err != nil {
			log.Fatal(err)
		}

		LogHTTP(gs, 200, r)
		w.WriteHeader(200)
	})

	http.HandleFunc("/api/server/stop/", func(w http.ResponseWriter, r *http.Request) {
		if gs.IsUp() {
			LogHTTP(gs, 200, r)
			go func() { gs.Stop() }()
			w.WriteHeader(200)
			return
		}

		LogHTTP(gs, 400, r)
		w.WriteHeader(400)
	})

	http.HandleFunc("/api/server/status/", func(w http.ResponseWriter, r *http.Request) {
		if gs.IsUp() {
			LogHTTP(gs, 200, r)
			w.WriteHeader(200)
		}

		LogHTTP(gs, 400, r)
		w.WriteHeader(400)
	})

	http.HandleFunc("/api/server/restart/", func(w http.ResponseWriter, r *http.Request) {
		if err := gs.Restart(); err != nil {
			LogHTTP(gs, 500, r)
			w.WriteHeader(500)
			log.Fatal(err)
		}

		LogHTTP(gs, 200, r)
		w.WriteHeader(200)
	})

	http.HandleFunc("/api/server/say/", func(w http.ResponseWriter, r *http.Request) {
		LogOutput(gs, "Sending message: "+r.RequestURI)
		u, _ := url.Parse(r.RequestURI)
		SendCommand("say "+strings.TrimPrefix(u.Path, "/api/server/say/"), gs)
		LogHTTP(gs, 200, r)
	})

	http.HandleFunc("/api/server/motd/", func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(r.RequestURI)
		m := strings.TrimPrefix(u.Path, "/api/server/motd")
		m = strings.TrimPrefix(m, "/")

		w.WriteHeader(200)
		if m == "" {
			w.Write([]byte(gs.MOTD()))
		} else {
			SendCommand("motd "+m, gs)
			SendCommand("motd", gs)
		}

		LogHTTP(gs, 200, r)
	})

	http.HandleFunc("/api/server/time/", func(w http.ResponseWriter, r *http.Request) {
		LogOutput(gs, "Received time request: "+r.RequestURI)
		u, _ := url.Parse(r.RequestURI)
		t := strings.TrimPrefix(u.Path, "/api/server/time")
		set := ""
		switch t {
		case "/", "":
			SendCommand("time", gs)
			return
		case "/dawn":
			set = "dawn"
		case "/noon":
			set = "noon"
		case "/dusk":
			set = "dusk"
		case "/midnight":
			set = "midnight"
		default:
			LogHTTP(gs, 404, r)
			w.WriteHeader(404)
			return
		}

		if set != "" {
			SendCommand("say Setting time to "+set, gs)
			SendCommand(set, gs)
		}

		w.WriteHeader(200)
		LogHTTP(gs, 200, r)
	})

	http.HandleFunc("/api/server/settle/", func(w http.ResponseWriter, r *http.Request) {
		LogInfo(gs, "Settling liquids", out)
		SendCommand("settle", gs)
		w.WriteHeader(200)
		LogHTTP(gs, 200, r)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(h, w, r)
	})
}
