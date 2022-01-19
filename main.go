package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/haozi4263/admission-resource/pkg"
	"k8s.io/klog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main()  {
	// webhook http server 需要和api-server交互需要是一个支持tls的webhook
	// 通过命令行参数传递证书
	var param pkg.WhSvrParam
	flag.IntVar(&param.Port, "port",443, "Webhook Server Port.")
	flag.StringVar(&param.CertFile, "tlsCertFile", "/etc/webhook/cert/tls.crt","x509 certification file")
	flag.StringVar(&param.KeyFile, "keyFile", "/etc/webhook/cert/tls.key","x509 private key file")
	flag.StringVar(&param.ConfigFile, "ConfigFile", "/etc/webhook/conf/mutate.yaml","mutate configFile")
	flag.Parse()

	cert, err :=tls.LoadX509KeyPair(param.CertFile, param.KeyFile)
	if err != nil {
		klog.Errorf("Failed to load key pair: %v", err)
		return
	}

	// 实例化一个Webhook Server
	whsrv := pkg.WebhookServer{
		Server: &http.Server{
			Addr: fmt.Sprintf(":%d", param.Port),
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		},
		RESOURCE_MULTIPLE: strings.Split(os.Getenv("RESOURCE_MULTIPLE"),":"),
	}

	// 定义http server handler
	mux := http.NewServeMux()
	mux.HandleFunc("/validate", whsrv.Handler)
	mux.HandleFunc("/mutate", whsrv.Handler)
	whsrv.Server.Handler = mux

	// 在一个新的goroutine里面启动 webhook server
	go func() {
		if err := whsrv.Server.ListenAndServeTLS("",""); err != nil {
			klog.Errorf("Failed to listen adn server webhook: %v", err)
		}
	}()
	klog.Info("Server started")
	// 监听OS的关闭新信号
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM )
	<- signalChan

	klog.Info("Got Os shutdown signal, gracefully shutting down...")
	if err := whsrv.Server.Shutdown(context.Background()); err != nil {
		klog.Errorf("HTTP Server Shutdown error: %v", err)
	}


}
