package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kardianos/service"
	hostfile "github.com/pkf/hosts/internal/hostFile"
	iplookup "github.com/pkf/hosts/internal/ipLookup"
)

var domains = []string{
	"alive.github.com",
	"live.github.com",
	"github.githubassets.com",
	"central.github.com",
	"desktop.githubusercontent.com",
	"assets-cdn.github.com",
	"camo.githubusercontent.com",
	"github.map.fastly.net",
	"github.global.ssl.fastly.net",
	"gist.github.com",
	"github.io",
	"github.com",
	"github.blog",
	"api.github.com",
	"raw.githubusercontent.com",
	"user-images.githubusercontent.com",
	"favicons.githubusercontent.com",
	"avatars5.githubusercontent.com",
	"avatars4.githubusercontent.com",
	"avatars3.githubusercontent.com",
	"avatars2.githubusercontent.com",
	"avatars1.githubusercontent.com",
	"avatars0.githubusercontent.com",
	"avatars.githubusercontent.com",
	"codeload.github.com",
	"github-cloud.s3.amazonaws.com",
	"github-com.s3.amazonaws.com",
	"github-production-release-asset-2e65be.s3.amazonaws.com",
	"github-production-user-asset-6210df.s3.amazonaws.com",
	"github-production-repository-file-5c1aeb.s3.amazonaws.com",
	"githubstatus.com",
	"github.community",
	"github.dev",
	"media.githubusercontent.com",
}

func main() {
	srvConfig := &service.Config{
		Name:        "AutoUpdateHosts",
		DisplayName: "Auto Update Hosts Service",
		Description: "自动更新Hosts文件服务",
	}
	prg := &program{}
	s, err := service.New(prg, srvConfig)
	if err != nil {
		fmt.Println(err)
	}
	if len(os.Args) > 1 {
		serviceAction := os.Args[1]
		switch serviceAction {
		case "install":
			err := s.Install()
			if err != nil {
				fmt.Println("安装服务失败: ", err.Error())
			} else {
				fmt.Println("安装服务成功")
			}
			return
		case "uninstall":
			err := s.Uninstall()
			if err != nil {
				fmt.Println("卸载服务失败: ", err.Error())
			} else {
				fmt.Println("卸载服务成功")
			}
			return
		case "start":
			err := s.Start()
			if err != nil {
				fmt.Println("运行服务失败: ", err.Error())
			} else {
				fmt.Println("运行服务成功")
			}
			return
		case "stop":
			err := s.Stop()
			if err != nil {
				fmt.Println("停止服务失败: ", err.Error())
			} else {
				fmt.Println("停止服务成功")
			}
			return
		}
	}

	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}
}

type program struct {
	timer    *time.Timer
	stopChan chan struct{}
}

func (p *program) Start(s service.Service) error {
	p.timer = time.NewTimer(4 * time.Hour)
	p.stopChan = make(chan struct{})
	go p.run()
	return nil
}
func (p *program) run() {
	hosts := iplookup.LookupDomains(domains)
	hostfile.UpdateHostFile(hosts)

	for {
		p.timer.Reset(4 * time.Hour)
		select {
		case <-p.timer.C:
			hosts := iplookup.LookupDomains(domains)
			hostfile.UpdateHostFile(hosts)
		case <-p.stopChan:
			p.timer.Stop()
			break
		}
	}

}
func (p *program) Stop(s service.Service) error {
	close(p.stopChan)
	return nil
}
