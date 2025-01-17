package xbvr

import (
	"bytes"
	"net"
	"path/filepath"
	"time"

	"github.com/cld9x/xbvr/pkg/assets"
	"github.com/cld9x/xbvr/pkg/dms/dlna/dms"
)

type dmsConfig struct {
	Path                string
	IfName              string
	Http                string
	FriendlyName        string
	LogHeaders          bool
	NoTranscode         bool
	NoProbe             bool
	StallEventSubscribe bool
	NotifyInterval      time.Duration
	IgnoreHidden        bool
	IgnoreUnreadable    bool
}

func StartDMS() {
	var config = &dmsConfig{
		Path:             "",
		IfName:           "",
		Http:             ":1338",
		FriendlyName:     "",
		LogHeaders:       false,
		NotifyInterval:   30 * time.Second,
	}

	dmsServer := &dms.Server{
		Interfaces: func(ifName string) (ifs []net.Interface) {
			var err error
			if ifName == "" {
				ifs, err = net.Interfaces()
			} else {
				var if_ *net.Interface
				if_, err = net.InterfaceByName(ifName)
				if if_ != nil {
					ifs = append(ifs, *if_)
				}
			}
			if err != nil {
				log.Fatal(err)
			}
			var tmp []net.Interface
			for _, if_ := range ifs {
				if if_.Flags&net.FlagUp == 0 || if_.MTU <= 0 {
					continue
				}
				tmp = append(tmp, if_)
			}
			ifs = tmp
			return
		}(config.IfName),
		HTTPConn: func() net.Listener {
			conn, err := net.Listen("tcp", config.Http)
			if err != nil {
				log.Fatal(err)
			}
			return conn
		}(),
		FriendlyName:   config.FriendlyName,
		RootObjectPath: filepath.Clean(config.Path),
		LogHeaders:     config.LogHeaders,
		NoTranscode:    config.NoTranscode,
		NoProbe:        config.NoProbe,
		Icons: []dms.Icon{
			{
				Width:      32,
				Height:     32,
				Depth:      8,
				Mimetype:   "image/png",
				ReadSeeker: bytes.NewReader(assets.FileIconsXbvr32Png),
			},
			{
				Width:      128,
				Height:     128,
				Depth:      8,
				Mimetype:   "image/png",
				ReadSeeker: bytes.NewReader(assets.FileIconsXbvr128Png),
			},
		},
		StallEventSubscribe: config.StallEventSubscribe,
		NotifyInterval:      config.NotifyInterval,
		IgnoreHidden:        config.IgnoreHidden,
		IgnoreUnreadable:    config.IgnoreUnreadable,
	}
	go func() {
		log.Info("Starting DMS")
		if err := dmsServer.Serve(); err != nil {
			log.Fatal(err)
		}
	}()
	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	// <-sigs
	// err := dmsServer.Close()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := cache.save(config.FFprobeCachePath); err != nil {
	// 	log.Print(err)
	// }
}
