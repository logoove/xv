package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	locale "github.com/Xuanwo/go-locale"
	"github.com/buger/jsonparser"
	"github.com/inconshreveable/mousetrap"
	"github.com/logoove/cli"
	"github.com/logoove/go/php"
	archiver "github.com/mholt/archiver/v3"
	"github.com/olekukonko/tablewriter"
	cache "github.com/patrickmn/go-cache"
	progressbar "github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var cn = php.MapStrStr{
	"name":      "版本管理下载工具",
	"cmd":       "这是一个命令行工具\n您需要打开cmd.exe并从那里运行它",
	"desc":      "一个支持go,python,nodejs,flutter的版本下载管理工具",
	"nover":     "暂无安装\n",
	"errorinfo": "参数错误,例如:",
	"okver":     "已安装版本\n",
	"downok":    "下载完成\n",
	"nook":      "下载文件不完整,请重试\n",
	"setupok":   "安装完成\n",
	"vererror":  "安装的版本不存在\n",
	"unok":      "卸载完成\n",
	"ok":        "设置成功\n",
	"gls":       "列出go已安装版本",
	"gall":      "列出go可安装版本",
	"gi":        "安装go指定版本",
	"gu":        "卸载go指定版本",
	"guse":      "设置go默认版本",
	"gset":      "设置go代理地址,不能下载时候设置,可选1,2,3",
	"gdel":      "删除go所有安装包",
	"pls":       "列出python已安装版本",
	"pall":      "列出python可安装版本",
	"pi":        "安装python指定版本",
	"pu":        "卸载python指定版本",
	"puse":      "设置python默认版本",
	"pset":      "设置python代理地址,不能下载时候设置,可选1,2,3",
	"pdel":      "删除python所有安装包",
	"nls":       "列出node已安装版本",
	"nall":      "列出node可安装版本",
	"ni":        "安装node指定版本",
	"nu":        "卸载node指定版本",
	"nuse":      "设置node默认版本",
	"nset":      "设置node代理地址,不能下载时候设置,可选1,2,3",
	"ndel":      "删除node所有安装包",
	"fls":       "列出flutter已安装版本",
	"fall":      "列出flutter可安装版本",
	"fi":        "安装flutter指定版本",
	"fu":        "卸载flutter指定版本",
	"fuse":      "设置flutter默认版本",
	"fset":      "设置flutter代理地址,不能下载时候设置,可选1,2",
	"fdel":      "删除flutter所有安装包",
	"urlerr":    "请求网址出错",
	"nozip":     "这个版本没有绿色版,换一个版本试试\n",
	"dirempty":  "下载目录为空\n",
}
var us = php.MapStrStr{
	"name":      "Download version management tools",
	"cmd":       "This is a command line tool\nYou need to open cmd.exe and run it from there.",
	"desc":      "A support go, python, nodejs, flutter download version management tools.",
	"nover":     "No install\n",
	"errorinfo": "Parameter error,For example:",
	"okver":     "Installed version\n",
	"downok":    "download finish\n",
	"nook":      "The download file is not complete, please try again\n",
	"setupok":   "The installation is complete\n",
	"vererror":  "The installed version does not exist\n",
	"unok":      "Uninstall complete\n",
	"ok":        "Setup success\n",
	"gls":       "List the go installed version",
	"gall":      "List all go version",
	"gi":        "Install the go specified version",
	"gu":        "Uninstall the go to specify version",
	"guse":      "Set the go version by default",
	"gset":      "Set go agent address, can't download time Settings, optional 1, 2, 3",
	"gdel":      "Delete all installation packages go",
	"pls":       "List the python installed version",
	"pall":      "List all python version",
	"pi":        "Install the python specified version",
	"pu":        "Uninstall the python to specify version",
	"puse":      "Set the python version by default",
	"pset":      "Set python agent address, can't download time Settings, optional 1, 2, 3",
	"pdel":      "Delete all installation packages python",
	"nls":       "List the node installed version",
	"nall":      "List all node version",
	"ni":        "Install the node specified version",
	"nu":        "Uninstall the node to specify version",
	"nuse":      "Set the node version by default",
	"nset":      "Set node agent address, can't download time Settings, optional 1, 2, 3",
	"ndel":      "Delete all installation packages node",
	"fls":       "List the flutter installed version",
	"fall":      "List all flutter version",
	"fi":        "Install the flutter specified version",
	"fu":        "Uninstall the flutter to specify version",
	"fuse":      "Set the flutter version by default",
	"fset":      "Set flutter agent address, can't download time Settings, optional 1, 2",
	"fdel":      "Delete all installation packages flutter",
	"urlerr":    "The request url error",
	"nozip":     "This version has no green version, try using a version\n",
	"dirempty":  "Download directory is empty\n",
}
var l php.MapStrStr
var gurl, purl, nurl, furl string
var ca *cache.Cache

//定义不同语言根目录,windows在C:\app\ Linux在/root/app/
var goDir, pyDir, nodeDir, flutterDir string

func init() {
	ca = cache.NewFrom(5*time.Minute, 60*time.Second, map[string]cache.Item{})
	ca.LoadFile("cache")
}
func main() {
	tag, _ := locale.Detect()
	if tag.String() == "zh-CN" {
		l = cn
	} else {
		l = us
	}
	if mousetrap.StartedByExplorer() {
		fmt.Print(l["cmd"])
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	//go
	gurls, isok1 := ca.Get("GURL")
	if isok1 == false {
		gurl = "https://studygolang.com/dl"
	} else {
		gurl = gurls.(string)
	}
	goDir = getGoDir()
	goDirs := filepath.Join(goDir, "go")       //安装目录
	goDown := filepath.Join(goDir, "download") //下载目录
	goVer := filepath.Join(goDir, "version")   //版本目录
	php.CreateDir(goDown, goVer)
	//python
	purls, isok2 := ca.Get("PURL")
	if isok2 == false {
		purl = "https://www.python.org/ftp/python"
	} else {
		purl = purls.(string)
	}
	pyDir = getPyDir()                         //py目录
	pyDirs := filepath.Join(pyDir, "python")   //python安装目录
	pyDown := filepath.Join(pyDir, "download") //下载目录
	pyVer := filepath.Join(pyDir, "version")   //版本目录
	php.CreateDir(pyDown, pyVer)
	//node
	nurls, isok3 := ca.Get("NURL")
	if isok3 == false {
		nurl = "https://npm.taobao.org/mirrors/node"
	} else {
		nurl = nurls.(string)
	}
	nodeDir = getNodeDir() //node根目录
	nodeDirs := filepath.Join(nodeDir, "nodejs")
	nodeDown := filepath.Join(nodeDir, "download") //下载目录
	nodeVer := filepath.Join(nodeDir, "version")   //版本目录
	set1 := filepath.Join(nodeDir, "node-global")  //全局文件夹
	set2 := filepath.Join(nodeDir, "node-cache")   //全局缓存
	php.CreateDir(nodeDown, nodeVer, set1, set2)
	//flutter
	furls, isok4 := ca.Get("FURL")
	if isok4 == false {
		furl = "https://storage.flutter-io.cn/flutter_infra_release/releases"
	} else {
		furl = furls.(string)
	}
	flutterDir = getFlutterDir() //flutter根目录
	flutterDirs := filepath.Join(flutterDir, "flutter")
	flutterDown := filepath.Join(flutterDir, "download") //下载目录
	flutterVer := filepath.Join(flutterDir, "version")   //版本目录
	php.CreateDir(flutterDown, flutterVer)

	app := cli.NewApp()
	app.Name = "xv"
	app.Version = "1.0.0"
	app.Authors = "Yoby\nWechat:logove\nUpdate:" + php.Date("Y-m-d H:i")
	app.Description = l["desc"]
	app.Usage = l["name"]
	app.SeeAlso = "2020-" + php.Date("Y")
	app.Commands = []*cli.Command{
		{
			Name:     "gls",
			Usage:    l["gls"],
			Examples: "xv gls",
			Action: func(c *cli.Context) {
				infos, err := os.ReadDir(goVer)
				if err != nil || len(infos) <= 0 {
					fmt.Printf(l["nover"])
				}
				for i := range infos {
					if !infos[i].IsDir() {
						continue
					}
					vname := infos[i].Name()
					php.Color(vname+"\n", "green")
				}
			},
		},
		{
			Name:     "gall",
			Usage:    l["gall"],
			Examples: "xv gall",
			Action: func(c *cli.Context) {
				vers, _ := getVer1()
				ss := php.SliceSplit(vers, 5)
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"A", "B", "C", "D", "E"})
				table.AppendBulk(ss)
				table.Render()
			},
		},
		{
			Name:     "gi",
			Usage:    l["gi"],
			Examples: "xv gi 1.17",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv gi 1.17", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                   //获取版本
				verDirs := filepath.Join(goVer, ver) //版本目录
				if php.IsDir(verDirs) {              //已安装直接退出
					php.Color(l["okver"]+ver, "green")
					os.Exit(0)
				}

				vers, _ := getVer1()
				if php.InArray(ver, vers) { //判断版本是否在列表中
					urls, names := getUrlName(ver, gurl)
					name := filepath.Join(goDown, names) //判断文件是否存在
					if !php.IsExist(name) {              //不存在下载
						downgo(urls, name)
						php.Color(l["downok"], "green")
					} else {
						_, sha1s := getVer1()
						sha1 := sha1s[ver]
						sha2, _ := SHA256File(name)
						if sha1 != sha2 {
							os.Remove(name) //删除不完整文件
							fmt.Println(l["nook"])
							os.Exit(0)
						}
					}

					os.RemoveAll(verDirs)                          //删除缓存的
					archiver.Unarchive(name, goVer)                //解压
					os.Rename(filepath.Join(goVer, "go"), verDirs) //重命名
					os.Remove(goDirs)                              //建立软连接
					os.Symlink(verDirs, goDirs)
					fmt.Println(l["setupok"], ver)
				} else {
					fmt.Println(l["vererror"])
				}
			},
		},
		{
			Name:     "gu",
			Usage:    l["gu"],
			Examples: "xv gu 1.17",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv gu 1.17", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                   //获取版本
				verDirs := filepath.Join(goVer, ver) //版本目录
				if !php.IsExist(verDirs) {
					php.Color(l["vererror"], "red")
					return
				}
				os.RemoveAll(verDirs)
				php.Color(l["unok"], "green")
			},
		},
		{
			Name:     "guse",
			Usage:    l["guse"],
			Examples: "xv guse 1.17",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv guse 1.17", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                   //获取版本
				verDirs := filepath.Join(goVer, ver) //版本目录
				if !php.IsExist(verDirs) {
					php.Color(l["vererror"], "red")
					return
				}
				os.Remove(goDirs) //建立软连接
				os.Symlink(verDirs, goDirs)
				if output, err := exec.Command(filepath.Join(goDirs, "bin", "go"), "version").Output(); err == nil {
					fmt.Print(l["ok"] + string(output))
				}
			},
		},
		{
			Name:     "gset",
			Usage:    l["gset"],
			Examples: "xv gset 1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv gset 1", "red")
					os.Exit(0)
				}

				if c.Args()[0] == "1" {
					ca.Set("GURL", "https://golang.google.cn/dl", 0)
				} else if c.Args()[0] == "2" {
					ca.Set("GURL", "https://dl.google.com", 0)
				} else {
					ca.Set("GURL", "https://studygolang.com/dl", 0)
				}
				vv, _ := ca.Get("GURL")
				ca.SaveFile("cache")
				php.Color(l["ok"]+":"+vv.(string), "green")
			},
		},
		{
			Name:     "gdel",
			Usage:    l["gdel"],
			Examples: "xv gdel",
			Action: func(c *cli.Context) {
				list, _ := os.ReadDir(goDown)
				if len(list) == 0 {
					fmt.Println(l["dirempty"])
					return
				}
				for i := range list {
					if err := os.RemoveAll(filepath.Join(goDown, list[i].Name())); err == nil {
						fmt.Println("Remove ok", list[i].Name())
					}
				}
			},
		},
		{
			Name:     "pls",
			Usage:    l["pls"],
			Examples: "xv pls",
			Action: func(c *cli.Context) {
				infos, err := os.ReadDir(pyVer)
				if err != nil || len(infos) <= 0 {
					fmt.Printf(l["nover"])
				}
				for i := range infos {
					if !infos[i].IsDir() {
						continue
					}
					vname := infos[i].Name()
					php.Color(vname+"\n", "green")
				}
			},
		},
		{
			Name:     "pall",
			Usage:    l["pall"],
			Examples: "xv pall",
			Action: func(c *cli.Context) {
				arr := make([][]string, 0)
				arr = php.SliceSplit(getver(), 5)
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"A", "B", "C", "D", "E"})
				table.AppendBulk(arr)
				table.Render()
			},
		},
		{
			Name:     "pi",
			Usage:    l["pi"],
			Examples: "xv pi 3.9.5",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv pi 3.9.5", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                   //获取版本
				verDirs := filepath.Join(pyVer, ver) //版本目录
				if php.IsDir(verDirs) {              //已安装直接退出
					php.Color(l["okver"]+ver, "green")
					os.Exit(0)
				}
				vers := getver() //获取版本列表
				if php.InArray(ver, vers) {

					urls, names := getUrlNamepy(ver, purl)
					name := filepath.Join(pyDown, names) //判断文件是否存在
					//fmt.Println(name)
					if name == "" { //这个版本没有绿色版
						fmt.Println(l["downok"])
						os.Exit(0)
					}
					if !php.IsExist(name) { //不存在下载
						downpy(urls, name)
						php.Color(l["downok"], "green")
					}

					os.RemoveAll(verDirs)             //删除缓存的
					archiver.Unarchive(name, verDirs) //解压
					os.Remove(pyDirs)
					os.Symlink(verDirs, pyDirs)
					fmt.Println(l["ok"], ver)
				} else {
					fmt.Println(l["nover"])
				}

			},
		},
		{
			Name:     "pu",
			Usage:    l["pu"],
			Examples: "xv pu 3.9.5",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv pu 3.9.5", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                   //获取版本
				verDirs := filepath.Join(pyVer, ver) //版本目录
				if !php.IsExist(verDirs) {
					php.Color(l["nover"], "red")
					return
				}
				os.RemoveAll(verDirs)
				php.Color(l["unok"], "green")
			},
		},
		{
			Name:     "puse",
			Usage:    l["puse"],
			Examples: "xv puse 12.18.1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorw"]+"xv puse 3.9.5", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                   //获取版本
				verDirs := filepath.Join(pyVer, ver) //版本目录
				if !php.IsExist(verDirs) {
					php.Color(l["nover"], "red")
					return
				}
				os.Remove(pyDirs)
				os.Symlink(verDirs, pyDirs)
				//fmt.Println(pyDirs)
				if output, err := exec.Command(filepath.Join(pyDirs, "python"), "-V").Output(); err == nil {
					fmt.Println(l["ok"] + string(output))
				}
			},
		},
		{
			Name:     "pset",
			Usage:    l["pset"],
			Examples: "xv pset 1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv pset 1", "red")
					os.Exit(0)
				}
				if c.Args()[0] == "1" {
					ca.Set("PURL", "https://www.python.org/ftp/python", 0)
				} else {

					ca.Set("PURL", "https://repo.huaweicloud.com/python", 0)
				}
				vv, _ := ca.Get("PURL")
				ca.SaveFile("cache")
				php.Color(l["ok"]+":"+vv.(string), "green")
			},
		},
		{
			Name:     "pdel",
			Usage:    l["pdel"],
			Examples: "xv pdel",
			Action: func(c *cli.Context) {
				list, _ := os.ReadDir(pyDown)
				if len(list) == 0 {
					fmt.Println(l["dirempty"])
					return
				}
				for i := range list {
					if err := os.RemoveAll(filepath.Join(pyDown, list[i].Name())); err == nil {
						fmt.Println("Remove ok", list[i].Name())
					}
				}
			},
		},
		{
			Name:     "nall",
			Usage:    l["nall"],
			Examples: "xv nall ",
			Action: func(c *cli.Context) {
				list := getnodever()
				ss := php.SliceSplit(list, 5)
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"A", "B", "C", "D", "E"})
				table.AppendBulk(ss)
				table.Render()
			},
		},
		{
			Name:     "nls",
			Usage:    l["nls"],
			Examples: "xv nls",
			Action: func(c *cli.Context) {
				infos, err := os.ReadDir(nodeVer)
				if err != nil || len(infos) <= 0 {
					fmt.Printf(l["nover"])
				}
				for i := range infos {
					if !infos[i].IsDir() {
						continue
					}
					vname := infos[i].Name()
					php.Color(vname+"\n", "green")
				}
			},
		},
		{
			Name:     "ni",
			Usage:    l["ni"],
			Examples: "xv ni 12.18.1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv ni 12.18.1", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                     //获取版本
				verDirs := filepath.Join(nodeVer, ver) //版本目录
				if php.IsDir(verDirs) {                //已安装直接退出
					php.Color(l["okver"]+ver, "green")
					os.Exit(0)
				}
				vers := getnodever() //获取版本列表
				if php.InArray(ver, vers) {
					urls, names := getNodeUrlName(ver, nurl)
					name := filepath.Join(nodeDown, names) //判断文件是否存在
					if !php.IsExist(name) {                //不存在下载
						downnode(urls, name)
						php.Color(l["downok"], "green")
					} else {
						sha2, _ := SHA256File(name)
						isok := getsha256(ver, sha2)
						if isok == false {
							os.Remove(name) //删除不完整文件
							fmt.Println(l["downok"])
							os.Exit(0)
						}
					}
					os.RemoveAll(verDirs)             //删除缓存的
					archiver.Unarchive(name, nodeVer) //解压
					file := filepath.Base(name)
					ext := path.Ext(name)
					os.Rename(filepath.Join(nodeVer, strings.TrimSuffix(file, ext)), verDirs)
					os.Remove(nodeDirs)
					os.Symlink(verDirs, nodeDirs)
					fmt.Println(l["ok"], ver)
				} else {
					fmt.Println(l["nover"])
				}

			},
		},
		{
			Name:     "nu",
			Usage:    l["nu"],
			Examples: "xv nu 12.18.1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv nu 12.18.1", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                     //获取版本
				verDirs := filepath.Join(nodeVer, ver) //版本目录
				if !php.IsExist(verDirs) {
					php.Color(l["nover"], "red")
					return
				}
				os.RemoveAll(verDirs)
				php.Color(l["unok"], "green")
			},
		},
		{
			Name:     "nuse",
			Usage:    l["nuse"],
			Examples: "xv nuse 12.18.1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv nuse 12.18.1", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                     //获取版本
				verDirs := filepath.Join(nodeVer, ver) //版本目录
				if !php.IsExist(verDirs) {
					php.Color(l["nover"], "red")
					return
				}
				os.Remove(nodeDirs)
				os.Symlink(verDirs, nodeDirs)
				if output, err := exec.Command(filepath.Join(nodeDirs, "node"), "-v").Output(); err == nil {
					fmt.Println(l["ok"] + string(output))
				}
			},
		},
		{
			Name:     "nset",
			Usage:    l["nset"],
			Examples: "xv nset 1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv nset 1", "red")
					os.Exit(0)
				}
				if c.Args()[0] == "1" {
					ca.Set("NURL", "https://npm.taobao.org/mirrors/node", 0)
				} else if c.Args()[0] == "2" {

					ca.Set("NURL", "https://mirror.tuna.tsinghua.edu.cn/nodejs-release", 0)
				} else {

					ca.Set("NURL", "https://repo.huaweicloud.com/nodejs", 0)
				}
				vv, _ := ca.Get("NURL")
				ca.SaveFile("cache")
				php.Color(l["ok"]+":"+vv.(string), "green")
			},
		},
		{
			Name:     "ndel",
			Usage:    l["ndel"],
			Examples: "xv ndel",
			Action: func(c *cli.Context) {
				list, _ := os.ReadDir(nodeDown)
				if len(list) == 0 {
					fmt.Println(l["dirempty"])
					return
				}
				for i := range list {
					if err := os.RemoveAll(filepath.Join(nodeDown, list[i].Name())); err == nil {
						fmt.Println("Remove ok ", list[i].Name())
					}
				}
			},
		},
		{
			Name:     "fall",
			Usage:    l["fall"],
			Examples: "xv fall",
			Action: func(c *cli.Context) {
				arr := make([][]string, 0)
				arr = php.SliceSplit(getFlutterVer(), 5)
				table := tablewriter.NewWriter(os.Stdout)
				table.SetHeader([]string{"A", "B", "C", "D", "E"})
				table.AppendBulk(arr)
				table.Render()
			},
		},
		{
			Name:     "fls",
			Usage:    l["fls"],
			Examples: "xv fls",
			Action: func(c *cli.Context) {
				infos, err := os.ReadDir(flutterVer)
				if err != nil || len(infos) <= 0 {
					fmt.Printf(l["nover"])
				}
				for i := range infos {
					if !infos[i].IsDir() {
						continue
					}
					vname := infos[i].Name()
					php.Color(vname+"\n", "green")
				}
			},
		},
		{
			Name:     "fi",
			Usage:    l["fi"],
			Examples: "xv fi 2.5.1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv fi 2.5.1", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                        //获取版本
				verDirs := filepath.Join(flutterVer, ver) //版本目录
				if php.IsDir(verDirs) {                   //已安装直接退出
					php.Color(l["okver"]+ver, "green")
					os.Exit(0)
				}
				vers := getFlutterVer() //获取版本列表
				if php.InArray(ver, vers) {
					urls, names := getFlutterUrlName(ver, furl)
					name := filepath.Join(flutterDown, names) //判断文件是否存在
					if !php.IsExist(name) {                   //不存在下载
						downflutter(urls, name)
						php.Color(l["downok"], "green")
					} else {
						fmt.Println(l["downok"])
					}
					os.RemoveAll(verDirs)                //删除缓存的
					archiver.Unarchive(name, flutterVer) //解压
					os.Rename(filepath.Join(flutterVer, "flutter"), verDirs)
					os.Remove(flutterDirs)
					os.Symlink(verDirs, flutterDirs)
					fmt.Println(l["ok"], ver)
				} else {
					fmt.Println(l["nover"])
				}

			},
		},
		{
			Name:     "fu",
			Usage:    l["fu"],
			Examples: "xv fu 2.5.1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv fu 2.5.1", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                        //获取版本
				verDirs := filepath.Join(flutterVer, ver) //版本目录
				if !php.IsExist(verDirs) {
					php.Color(l["nover"], "red")
					return
				}
				os.RemoveAll(verDirs)
				php.Color(l["unok"], "green")
			},
		},
		{
			Name:     "fuse",
			Usage:    l["fuse"],
			Examples: "xv fuse 2.5.1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv fuse 2.5.1", "red")
					os.Exit(0)
				}
				ver := c.Args()[0]                        //获取版本
				verDirs := filepath.Join(flutterVer, ver) //版本目录
				if !php.IsExist(verDirs) {
					php.Color(l["nover"], "red")
					return
				}
				os.Remove(flutterDirs)
				os.Symlink(verDirs, flutterDirs)
				if output, err := exec.Command(filepath.Join(flutterDirs, "bin", "flutter"), "--version").Output(); err == nil {
					fmt.Println(l["ok"] + string(output))
				}
			},
		},
		{
			Name:     "fset",
			Usage:    l["fset"],
			Examples: "xv fset 1",
			Action: func(c *cli.Context) {
				if c.NArg() == 0 {
					php.Color(l["errorinfo"]+"xv fset 1", "red")
					os.Exit(0)
				}
				if c.Args()[0] == "1" {
					ca.Set("FURL", "https://storage.flutter-io.cn/flutter_infra_release/releases", 0)
				} else {

					ca.Set("FURL", "https://storage.googleapis.com/flutter_infra_release/releases", 0)
				}
				vv, _ := ca.Get("FURL")
				ca.SaveFile("cache")
				php.Color(l["ok"]+":"+vv.(string), "green")
			},
		},
		{
			Name:     "fdel",
			Usage:    l["fdel"],
			Examples: "xv fdel",
			Action: func(c *cli.Context) {
				list, _ := os.ReadDir(flutterDown)
				if len(list) == 0 {
					fmt.Println(l["dirempty"])
					return
				}
				for i := range list {
					if err := os.RemoveAll(filepath.Join(flutterDown, list[i].Name())); err == nil {
						fmt.Println("Remove ok ", list[i].Name())
					}
				}
			},
		},
	}
	app.Run(os.Args)
}

//
//  getGoDir
//  @Description: 获取go安装目录
//  @return d
//
func getGoDir() (d string) {
	homeDir := php.Getosdir() + "/"
	homeDir = filepath.Join(homeDir, "app")
	return filepath.Join(homeDir, "go")
}

//
//  getPyDir
//  @Description: 获取py根目录
//  @return d
//
func getPyDir() (d string) {
	homeDir := php.Getosdir() + "/"
	homeDir = filepath.Join(homeDir, "app")
	return filepath.Join(homeDir, "python")
}

//
//  getNodeDir
//  @Description: 获取node目录
//  @return d
//
func getNodeDir() (d string) {
	homeDir := php.Getosdir() + "/"
	homeDir = filepath.Join(homeDir, "app")
	return filepath.Join(homeDir, "nodejs")
}

//
//  getFlutterDir
//  @Description: 获取flutter目录
//  @return d
//
func getFlutterDir() (d string) {
	homeDir := php.Getosdir() + "/"
	homeDir = filepath.Join(homeDir, "app")
	return filepath.Join(homeDir, "flutter")
}
func getdoc() string {
	res, _ := http.Get(gurl)
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, _ := io.ReadAll(res.Body)
	str := string(body)
	return str
}

//
//  getVer1
//  @Description: 获取go版本
//  @return []string
//  @return map[string]string
//
func getVer1() ([]string, map[string]string) {

	var str string
	html, isok1 := ca.Get("html")
	if isok1 == false {
		str = getdoc()
		ca.Set("html", str, 60*time.Second)
	} else {
		str = html.(string)
	}
	ca.SaveFile("cache")
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
	var sc []string
	sha1x := make(map[string]string)
	doc.Find(`.toggleVisible,.toggle`).Each(func(i int, s *goquery.Selection) {
		ss, _ := s.Attr("id")
		if ss != "archive" && i < 30 { //只取最近30个版本
			ss = string([]byte(ss)[2:]) //去掉go两个字符
			sc = append(sc, ss)
			_, ver := getUrlName(ss, gurl)
			sha1x[ss] = doc.Find(`a:contains("` + ver + `")`).Parent().Parent().Find("tt").Text()
		}
		i++
	})

	return sc, sha1x

}

//
//  getUrlName
//  @Description: 获取下载链接和名称
//  @param ver
//  @param u
//  @return url
//  @return name
//
func getUrlName(ver, u string) (url, name string) {
	var ext string
	if runtime.GOOS == "windows" {
		ext = "zip"
	} else {
		ext = "tar.gz"
	}
	name = "go" + ver + "." + runtime.GOOS + "-" + runtime.GOARCH + "." + ext
	if u == "https://studygolang.com/dl" {
		url = u + "/golang/" + name
	} else if u == "https://golang.google.cn/dl" {
		url = u + "/" + name
	} else { //国外
		url = u + "/go/" + name
	}
	return url, name
}

//
//  downgo
//  @Description: 下载go
//  @param url
//  @param name
//
func downgo(url, name string) {
	req, _ := http.NewRequest("GET", url, nil)
	fmt.Println("Download Url:" + url + "\n")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	f, _ := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stdout, "\n")
		}),
	)
	bar.RenderBlank()
	io.Copy(io.MultiWriter(f, bar), resp.Body)
}

//
//  SHA256File
//  @Description: 验证go文件
//  @param path
//  @return string
//  @return error
//
func SHA256File(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}
	h := sha256.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

//
//  getver
//  @Description: 获取python版本号
//  @return []string
//
func getver() []string {

	var str string
	html, isok1 := ca.Get("html1")
	if isok1 == false {
		res, _ := http.Get(purl)
		defer res.Body.Close()
		if res.StatusCode != 200 {
			fmt.Println(l["urlerr"]+":%d %s", res.StatusCode, res.Status)
		}
		body, _ := io.ReadAll(res.Body)
		str = string(body)
		ca.Set("html1", str, 60*time.Second)
	} else {
		str = html.(string)
	}
	ca.SaveFile("cache")

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
	var sc []string
	var a1, b1 int
	doc.Find(`a`).Each(func(i int, s *goquery.Selection) {
		ss := strings.TrimSuffix(s.Text(), "/")
		if strings.Contains(ss, "binaries-1.1") {
			a1 = i
		}
		if ss == "3.5.0" {
			b1 = i
		}
		sc = append(sc, ss)
	})
	sc1 := sc[b1:a1]
	sc1 = append([]string{"2.7.1"}, sc1...)
	return sc1
}

//
//  getFlutterVer
//  @Description: 获取flutter版本
//  @return []string
//
func getFlutterVer() []string {

	var str string
	html, isok1 := ca.Get("html2")
	if isok1 == false {
		str = php.Get(furl + "/releases_windows.json")
		ca.Set("html2", str, 60*time.Second)
	} else {
		str = html.(string)
	}
	ca.SaveFile("cache")

	data := []byte(str)
	var a []string
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		ch, _ := jsonparser.GetString(value, "channel")
		if ch == "stable" {
			ver, _ := jsonparser.GetString(value, "version")
			if offset > 43558 {
				return
			} else {
				a = append(a, ver)
			}

		}

	}, "releases")
	return a
}

//
//  getUrlNamepy
//  @Description: 获取py版本名称
//  @param ver
//  @param u
//  @return url1
//  @return name
//
func getUrlNamepy(ver, u string) (url1, name string) {
	if runtime.GOOS == "windows" {
		if ver == "2.7.1" {
			name = "python-" + ver + "-embed-" + runtime.GOARCH + ".zip"
			url1 = ""
		} else {
			name = "python-" + ver + "-embed-" + runtime.GOARCH + ".zip"
			url1 = u + "/" + ver + "/" + name
		}

	} else {
		name = ""
		url1 = ""
	}
	return url1, name
}

//
//  @Description: 下载py
//  @param u
//  @param n
//
func downpy(u, n string) {
	req, _ := http.NewRequest("GET", u, nil)
	fmt.Println("Download Url:" + u + "\n")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	types := resp.Header.Get("Content-Type")
	if types == "text/html" {
		fmt.Println(l["nozip"])
	} else {
		f, _ := os.OpenFile(n, os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		bar := progressbar.NewOptions64(
			resp.ContentLength,
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetDescription("Downloading"),
			progressbar.OptionSetWriter(os.Stdout),
			progressbar.OptionShowBytes(true),
			progressbar.OptionThrottle(65*time.Millisecond),
			progressbar.OptionShowCount(),
			progressbar.OptionOnCompletion(func() {
				fmt.Fprint(os.Stdout, "\n")
			}),
		)
		bar.RenderBlank()
		io.Copy(io.MultiWriter(f, bar), resp.Body)
	}
}

//
//  getnodever
//  @Description: 获取node版本
//  @return []string
//
func getnodever() []string {
	var str string
	html, isok1 := ca.Get("html3")
	if isok1 == false {
		str = php.Get(nurl + "/index.json")
		ca.Set("html3", str, 60*time.Second)
	} else {
		str = html.(string)
	}
	ca.SaveFile("cache")

	var ver []map[string]string
	json.Unmarshal([]byte(str), &ver)
	var a []string
	for k, v := range ver {
		if k < 150 {
			a = append(a, strings.TrimLeft(v["version"], "v"))
		}
	}
	return a
}

//
//  getNodeUrlName
//  @Description: 获取node下载名称
//  @param ver
//  @param u
//  @return url1
//  @return name
//
func getNodeUrlName(ver, u string) (url1, name string) {
	var ext, osx, nm string
	if runtime.GOOS == "windows" {
		ext = "zip"
		osx = "win"
		nm = "x64"
	} else if runtime.GOOS == "darwin" {
		ext = "tar.gz"
		osx = "darwin"
		nm = "x64"
	} else {
		ext = "tar.gz"
		osx = "linux"
		nm = "x64"
	}
	name = "node-v" + ver + "-" + osx + "-" + nm + "." + ext
	url1 = u + "/v" + ver + "/" + name
	return url1, name
}

//
//  downnode
//  @Description: 下载node
//  @param u
//  @param n
//
func downnode(u, n string) {
	req, _ := http.NewRequest("GET", u, nil)
	fmt.Println("Download Url:" + u + "\n")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	f, _ := os.OpenFile(n, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stdout, "\n")
		}),
	)
	bar.RenderBlank()
	io.Copy(io.MultiWriter(f, bar), resp.Body)
}

//
//  getsha256
//  @Description: 判断sha256是否在文件中
//  @param ver
//  @param s256
//  @return s
//
func getsha256(ver string, s256 string) (s bool) {
	u := nurl + "/v" + ver + "/SHASUMS256.txt"
	str := php.Get(u)
	s = strings.Contains(str, s256)
	return s
}

//
//  getFlutterUrlName
//  @Description: 获取flutter下载链接
//  @param ver
//  @param u
//  @return url1
//  @return name
//
func getFlutterUrlName(ver, u string) (url1, name string) {
	var ext, osx string
	if runtime.GOOS == "windows" {
		ext = "zip"
		osx = "windows"

	} else if runtime.GOOS == "darwin" {
		ext = "zip"
		osx = "macos"

	} else {
		ext = "tar.xz"
		osx = "linux"

	}
	name = "flutter_" + osx + "_" + ver + "-stable" + "." + ext
	url1 = u + "/stable/" + osx + "/" + name
	return url1, name
}

//
//  downflutter
//  @Description: 下载flutter
//  @param u
//  @param n
//
func downflutter(u, n string) {
	req, _ := http.NewRequest("GET", u, nil)
	fmt.Println("Download Url:" + u + "\n")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	f, _ := os.OpenFile(n, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetWriter(os.Stdout),
		progressbar.OptionShowBytes(true),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stdout, "\n")
		}),
	)
	bar.RenderBlank()
	io.Copy(io.MultiWriter(f, bar), resp.Body)
}
