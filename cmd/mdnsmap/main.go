package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"mdnsmap/internal/config"
	"mdnsmap/internal/discovery"
	"mdnsmap/internal/netutil"
	"mdnsmap/internal/output"
)

func main() {
	var (
		cidrFlag     = flag.String("cidr", "", "目标IP网段，例如 192.168.1.0/24")
		portsFlag    = flag.String("ports", "", "端口范围，例如 1-1024 或 9,445,548,5000")
		servicesFlag = flag.String("services", "", "服务类型列表，默认内置")
		timeoutFlag  = flag.Duration("timeout", 5*time.Second, "mDNS浏览超时")
		prettyFlag   = flag.Bool("pretty", false, "格式化输出JSON")
		datasetFlag  = flag.Bool("dataset", false, "输出数据集格式")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "mDNS资产测绘工具\n")
		fmt.Fprintf(os.Stderr, "用法:\n")
		fmt.Fprintf(os.Stderr, "  %s -cidr <CIDR> -ports <PORTS> [选项]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *cidrFlag == "" || *portsFlag == "" {
		flag.Usage()
		os.Exit(1)
	}

	ipnet, err := netutil.ParseCIDR(*cidrFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "无效网段: %v\n", err)
		os.Exit(1)
	}

	allowedPorts, err := netutil.ParsePorts(*portsFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "无效端口范围: %v\n", err)
		os.Exit(1)
	}

	services := netutil.ParseServices(*servicesFlag)

	cfg := &config.Config{
		CIDR:        *cidrFlag,
		Ports:       *portsFlag,
		Services:    *servicesFlag,
		Timeout:     *timeoutFlag,
		Pretty:      *prettyFlag,
		DatasetMode: *datasetFlag,
	}

	ctx := context.Background()
	scanCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	assets, err := discovery.Scan(scanCtx, cfg, ipnet, allowedPorts, services)
	if err != nil {
		fmt.Fprintf(os.Stderr, "扫描错误: %v\n", err)
		os.Exit(1)
	}

	if cfg.DatasetMode {
		if err := output.PrintDataset(assets, cfg.Pretty); err != nil {
			fmt.Fprintf(os.Stderr, "输出错误: %v\n", err)
			os.Exit(1)
		}
	} else {
		if err := output.PrintJSON(assets, cfg.Pretty); err != nil {
			fmt.Fprintf(os.Stderr, "输出错误: %v\n", err)
			os.Exit(1)
		}
	}
}
