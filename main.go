package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	"github.com/lacion/iothub/config"
	"github.com/lacion/iothub/log"
	"github.com/lacion/iothub/middlewares"
	"github.com/olahol/melody"
)

func main() {

	versionFlag := flag.Bool("version", false, "Version")
	flag.Parse()

	if *versionFlag {
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Version:", Version)
		if VersionPrerelease != "" {
			fmt.Println("Version PreRelease:", VersionPrerelease)
		}
		return
	}

	log.WithFields(log.Fields{
		"EventName": "get_config_env_vars",
	}).Debug("Reading configuration from env vars")
	cfg := config.Config()

	log.WithFields(log.Fields{
		"EventName": "set_gin_mode",
		"Mode":      cfg.GetString("mode"),
	}).Debug("Setting gin mode to ", cfg.GetString("mode"))
	gin.SetMode(cfg.GetString("mode"))

	r := gin.New()
	m := melody.New()

	r.Use(ginrus.Ginrus(log.NewLogger(cfg), time.RFC3339, true))
	r.Use(gin.Recovery())

	authorized := r.Group("/")

	authorized.Use(middlewares.Auth())

	// Web Sockets

	authorized.GET("/channel/:name/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		log.WithFields(log.Fields{
			"EventName":     "ws_client_connect",
			"RemoteAddress": s.Request.RemoteAddr,
		}).Debug("new ws client connected ", s.Request.RemoteAddr)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.WithFields(log.Fields{
			"EventName":     "ws_client_disconnect",
			"RemoteAddress": s.Request.RemoteAddr,
		}).Debug("ws client disconnected ", s.Request.RemoteAddr)
	})

	m.HandleError(func(s *melody.Session, err error) {
		log.WithFields(log.Fields{
			"EventName":     "ws_error",
			"RemoteAddress": s.Request.RemoteAddr,
			"Error":         err.Error(),
		}).Error("error ocurred with ws client ", err.Error())
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			msgStr := string(msg[:])

			log.WithFields(log.Fields{
				"EventName": "ws_client_message",
				"Message":   msgStr,
				"Channel":   q.Request.URL.Path,
			}).Debug("got msg: ", msgStr)

			return q.Request.URL.Path == s.Request.URL.Path
		})
	})

	// Start Server

	log.WithFields(log.Fields{
		"EventName":         "start",
		"ListenAddress":     cfg.GetString("listen_address"),
		"GitCommit":         GitCommit,
		"Version":           Version,
		"VersionPrerelease": VersionPrerelease,
		"TLS": cfg.GetBool("secure"),
		"ReadTimeout": cfg.GetDuration("read_timeout"),
		"WriteTimeout": cfg.GetDuration("write_timeout"),
		"MaxHeaderBytes": cfg.GetInt("max_header_bytes"),
	}).Info("starting server and listening on ", cfg.GetString("listen_address"))

	// ssl labs recomendations
	tlscfg := &tls.Config{
		// the bellow 3 options would get a perfect score in ssl labs test, but you will losse client compatibility
		// MinVersion: tls.VersionTLS12,
		// CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		//CipherSuites: []uint16{
		//	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		//	tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		//	tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		//	tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		//},
		PreferServerCipherSuites: true,
	}

	server := &http.Server{
		Addr:           cfg.GetString("listen_address"),
		Handler:        r,
		ReadTimeout:    cfg.GetDuration("read_timeout"),
		WriteTimeout:   cfg.GetDuration("write_timeout"),
		MaxHeaderBytes: cfg.GetInt("max_header_bytes"),
		TLSConfig:    tlscfg,
		// the bellow line would get a perfect score in ssl labs test, but you will losse client compatibility
		// also no http/2
		// TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	switch cfg.GetBool("secure") {
	case true:
		// Key or DH parameter strength >= 4096 bits (e.g., 4096)
		log.Fatal(server.ListenAndServeTLS(cfg.GetString("cert_file"), cfg.GetString("key_file")))
	default:
		log.Fatal(server.ListenAndServe())
	}
}
