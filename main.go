/*****************************************************************
* Copyright©,2020-2022, email: 279197148@qq.com
* Version: 1.0.0
* @Author: yangtxiang
* @Date: 2020-08-19 16:49
* Description:
*****************************************************************/

package main

import (
	"flag"
	"fmt"
	"github.com/go-xe2/x/os/xfile"
	"github.com/go-xe2/x/type/t"
	"github.com/go-xe2/xthrift/pdl"
	"github.com/go-xe2/xthrift/registerCli"
)

var (
	bHelp = flag.Bool("h", false, "显示帮助")
	serverHost = flag.String("s", "", "服务地址")
	host = flag.String("host", "", "本地服务地址")
	port = flag.String("port", "", "本地服务监听端口")
	projPath = flag.String("d", "", "协议项目目录")
	pdlFile = flag.String("f", "", "要注册的协议项目pdl文件, 同时输入f和d参数时，f参数优先")
)

func main() {
	flag.Parse()

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	if *bHelp {
		flag.Usage()
		return
	}

	if *serverHost == "" {
		panic("请输入远程注册服务器的地址")
	}
	if *host == "" {
		panic("请输入本地服务地址")
	}
	nPort := t.Int(*port)
	if nPort <= 0 {
		panic("请输入本地服务监听端口号")
	}
	if *projPath == "" && *pdlFile == "" {
		panic("请输入协议项目目录-d参数或输入协议文件-f参数")
	}
	if *pdlFile != "" {
		if !xfile.Exists(*pdlFile) {
			panic(fmt.Errorf("协议文件%s不存在", *pdlFile))
		}
	}
	if *projPath != "" {
		if !xfile.Exists(*projPath) {
			panic(fmt.Errorf("协议工作目录%s不存在", *projPath))
		}
	}

	var project *pdl.FileProject

	if *pdlFile != "" {
		project = pdl.NewEmptyFileProject()
		if err := project.LoadFromFile(*pdlFile); err != nil {
			panic(err)
		}
	}
	if *projPath != "" {
		var err error
		project, err = pdl.NewFileProject(*projPath)
		if err != nil {
			panic(err)
		}
		if err := project.Load(); err != nil {
			panic(err)
		}
	}
	if project == nil {
		panic(fmt.Errorf("输入参数有误，使用-h查看帮助"))
	}
	if err := project.Check(); err != nil {
		panic(err)
	}
	client := registerCli.NewRegisterClient(*serverHost)
	err := client.Register(*host, nPort, project)
	if err != nil {
		panic(err)
	}
	fmt.Println("注册成功")
}

