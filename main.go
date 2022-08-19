package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	locale "github.com/Xuanwo/go-locale"
	"github.com/bytedance/sonic"
	version "github.com/hashicorp/go-version"
	"github.com/inconshreveable/mousetrap"
	"github.com/logoove/cli"
	"github.com/logoove/go/php"
	archiver "github.com/mholt/archiver/v3"
	"github.com/olekukonko/tablewriter"
	progressbar "github.com/schollz/progressbar/v3"
	yaml "gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

var l map[string]any
var goDir, goDown, goInstall, goVer string
var pyDir, pyDown, pyInstall, pyVer string
var nodeDir, nodeDown, nodeInstall, nodeVer string
var flutterDir, flutterDown, flutterInstall, flutterVer string
var rs map[string]any
var data map[string]any

func main() {
	data = map[string]any{
		"gurl":   "https://studygolang.com/dl",
		"purl":   "https://www.python.org/ftp/python",
		"nurl":   "https://nodejs.org/download/release",
		"furl":   "https://storage.flutter-io.cn/flutter_infra_release/releases",
		"giturl": "https://ghproxy.com/",
		"ver":    "1.2.0",
		"cn": map[string]string{
			"name": "版本管理下载工具",
			"cmd":  "这是一个命令行工具\n您需要打开cmd.exe并从那里运行它",
			"desc": "一个支持go,python,nodejs,flutter的版本下载管理工具,命令需要管理员权限下执行",
		},
		"us": map[string]string{
			"name": "Download version management tools",
			"cmd":  "This is a command line tool\nYou need to open cmd.exe and run it from there.",
			"desc": "A support go, python, nodejs, flutter download version management tools.You need administrator rights command execution.",
		},
	}
	out, _ := yaml.Marshal(&data)
	if php.IsExist("xv.yaml") == false {
		os.WriteFile("xv.yaml", out, 0666)
	}
	str, _ := os.ReadFile("xv.yaml")
	yaml.Unmarshal(str, &rs)

	tag, _ := locale.Detect()
	if tag.String() == "zh-CN" {
		l = rs["cn"].(map[string]any)
	} else {
		l = rs["us"].(map[string]any)
	}
	if mousetrap.StartedByExplorer() {
		fmt.Print(l["cmd"])
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	//go根目录
	goDir = filepath.Join(php.Getosdir(), "app", "go")
	goInstall = filepath.Join(goDir, "go")
	goDown = filepath.Join(goDir, "download")
	goVer = filepath.Join(goDir, "version")
	php.CreateDir(pathx(goDown))
	php.CreateDir(pathx(goVer))

	//py根目录
	pyDir = filepath.Join(php.Getosdir(), "app", "python")
	pyInstall = filepath.Join(pyDir, "python")
	pyDown = filepath.Join(pyDir, "download")
	pyVer = filepath.Join(pyDir, "version")
	php.CreateDir(pathx(pyDown), pathx(pyVer))

	//nodejs
	nodeDir = filepath.Join(php.Getosdir(), "app", "nodejs")
	nodeInstall = filepath.Join(nodeDir, "nodejs")
	nodeDown = filepath.Join(nodeDir, "download")
	nodeVer = filepath.Join(nodeDir, "version")
	set1 := filepath.Join(nodeDir, "node-global") //全局文件夹
	set2 := filepath.Join(nodeDir, "node-cache")  //全局缓存
	php.CreateDir(nodeDown, nodeVer, set1, set2)
	//flutter
	flutterDir = filepath.Join(php.Getosdir(), "app", "flutter")
	flutterInstall = filepath.Join(flutterDir, "flutter")
	flutterDown = filepath.Join(flutterDir, "download")
	flutterVer = filepath.Join(flutterDir, "version")
	php.CreateDir(flutterDown, flutterVer)
	app := cli.NewApp()
	app.Name = "xv"
	app.Version = rs["ver"].(string)
	app.Authors = "Yoby\nWechat:logove\nUpdate:" + php.Date("Y-m-d", php.Timestr2Time("2022-08-18 00:00:00"))
	app.Description = l["desc"].(string)
	app.Usage = l["name"].(string)
	app.SeeAlso = "2020-" + php.Date("Y")
	app.Commands = []*cli.Command{
		{Name: "----", Usage: "-----Golang--------"},
		{
			Name:     "gls",
			Usage:    "Lists the installed version",
			Examples: "xv gls",
			Action:   gls,
		},
		{
			Name:     "gall",
			Usage:    "List all version",
			Examples: "xv gall",
			Action:   gall,
		},
		{
			Name:     "gi",
			Usage:    "Install a version",
			Examples: "xv gi 1.19",
			Action:   gi,
		},
		{
			Name:     "gu",
			Usage:    "Uninstall a version",
			Examples: "xv gu 1.19",
			Action:   gu,
		},
		{
			Name:     "guse",
			Usage:    "Switch version",
			Examples: "xv guse 1.19",
			Action:   guse,
		},
		{
			Name:     "gdel",
			Usage:    "Delete the package",
			Examples: "xv gdel",
			Action:   gdel,
		},
		{Name: "----", Usage: "------Python-------"},
		{
			Name:     "pall",
			Usage:    "List all version",
			Examples: "xv pall",
			Action:   pall,
		},
		{
			Name:     "pls",
			Usage:    "Lists the installed version",
			Examples: "xv pls",
			Action:   pls,
		},
		{
			Name:     "pi",
			Usage:    "Install a version",
			Examples: "xv pi 3.10.0",
			Action:   pi,
		},
		{
			Name:     "pu",
			Usage:    "Uninstall a version",
			Examples: "xv pu 3.10.0",
			Action:   pu,
		},
		{
			Name:     "puse",
			Usage:    "Switch version",
			Examples: "xv puse 3.10.0",
			Action:   puse,
		},
		{
			Name:     "pdel",
			Usage:    "Delete the package",
			Examples: "xv pdel",
			Action:   pdel,
		},
		{Name: "----", Usage: "------Nodejs-------"},
		{
			Name:     "nall",
			Usage:    "List all version",
			Examples: "xv nall",
			Action:   nall,
		},
		{
			Name:     "nls",
			Usage:    "Lists the installed version",
			Examples: "xv nls",
			Action:   nls,
		},
		{
			Name:     "ni",
			Usage:    "Install a version",
			Examples: "xv ni 16.0.0",
			Action:   ni,
		},
		{
			Name:     "nu",
			Usage:    "Uninstall a version",
			Examples: "xv nu 16.0.0",
			Action:   nu,
		},
		{
			Name:     "nuse",
			Usage:    "Switch version",
			Examples: "xv nuse 16.0.0",
			Action:   nuse,
		},
		{
			Name:     "ndel",
			Usage:    "Delete the package",
			Examples: "xv ndel",
			Action:   ndel,
		},
		{Name: "----", Usage: "-------Flutter------"},
		{
			Name:     "fall",
			Usage:    "List all version",
			Examples: "xv fall",
			Action:   fall,
		},
		{
			Name:     "fls",
			Usage:    "Lists the installed version",
			Examples: "xv fls",
			Action:   fls,
		},
		{
			Name:     "fi",
			Usage:    "Install a version",
			Examples: "xv fi 3.0.1",
			Action:   fi,
		},
		{
			Name:     "fu",
			Usage:    "Uninstall a version",
			Examples: "xv fu 3.0.1",
			Action:   fu,
		},
		{
			Name:     "fuse",
			Usage:    "Switch version",
			Examples: "xv fuse 3.0.1",
			Action:   fuse,
		},
		{
			Name:     "fdel",
			Usage:    "Delete the package",
			Examples: "xv fdel",
			Action:   fdel,
		},
		{Name: "----", Usage: "-------Other------"},
		{
			Name:     "get",
			Usage:    "Download github package",
			Examples: "xv get url",
			Action:   get,
		},
	}
	app.Run(os.Args)
}

func get(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	url := c.Args()[0]
	if strings.Contains(url, "https://github.com/") {
		url = rs["giturl"].(string) + url
	}
	f, _ := os.UserHomeDir()
	f = f + "\\Desktop" //桌面
	name := filepath.Base(url)
	wget(url, f+"\\"+name)
	php.Color("Download Ok", "green")
}
func gls(c *cli.Context) {
	infos, err := os.ReadDir(goVer)
	if err != nil || len(infos) <= 0 {
		php.Color("No installation", "")
	} else {
		php.Color("Is already installed\n", "green")
	}
	for i := range infos {
		if !infos[i].IsDir() {
			continue
		}
		vname := infos[i].Name()
		php.Color("√"+vname+"\n", "green")
	}
}
func gall(c *cli.Context) {
	var vers []string

	if rs["gover"] == nil {
		vers, _ = getGoAllVer()
		rs["gover"] = vers
		out, _ := yaml.Marshal(&rs)
		os.WriteFile("xv.yaml", out, 0666)
	} else {

		vers = any2str(rs["gover"])
	}
	vers = verSort(vers)
	ss := php.SliceSplit(vers, 5)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"A", "B", "C", "D", "E"})
	table.AppendBulk(ss)
	table.Render()
}
func gi(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                   //获取版本
	verDirs := filepath.Join(goVer, ver) //版本目录
	if php.IsDir(verDirs) {              //已安装直接退出
		php.Color("√"+ver, "green")
		os.Exit(0)
	}
	vers, _ := getGoAllVer()
	if php.InArray(ver, vers) { //判断版本是否在列表中
		urls, names := getUrlName(ver, rs["gurl"].(string))
		name := filepath.Join(goDown, names) //判断文件是否存在
		if !php.IsExist(name) {              //不存在下载
			downgo(urls, name)
			php.Color("Download successful\n", "green")
		} else {
			_, sha1s := getGoAllVer()
			sha1 := sha1s[ver]
			sha2, _ := SHA256File(name)
			if sha1 != sha2 {
				os.Remove(name) //删除不完整文件
				fmt.Println("File is incomplete, try again")
				os.Exit(0)
			}
			fmt.Println("Loading...\n")
		}

		os.RemoveAll(verDirs)                          //删除缓存的
		archiver.Unarchive(name, goVer)                //解压
		os.Rename(filepath.Join(goVer, "go"), verDirs) //重命名
		os.RemoveAll(goInstall)                        //建立软连接
		os.Symlink(verDirs, goInstall)
		fmt.Println("Successful installation", ver)
	} else {
		fmt.Println("Version does not exist")
	}

}
func gu(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                   //获取版本
	verDirs := filepath.Join(goVer, ver) //版本目录
	if !php.IsExist(verDirs) {
		php.Color("Version does not exist.\n", "red")
		return
	}
	os.RemoveAll(verDirs)
	php.Color("Uninstall the success.\n", "green")
}
func guse(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                          //获取版本
	verDirs := pathx(filepath.Join(goVer, ver)) //版本目录
	if !php.IsExist(verDirs) {
		php.Color("Version does not exist.\n", "red")
		return
	}
	goInstall = pathx(goInstall)
	os.RemoveAll(goInstall) //删除go文件
	os.Symlink(verDirs, goInstall)
	if output, err := exec.Command(filepath.Join(goInstall, "bin", "go"), "version").Output(); err == nil {
		fmt.Print("Switch success\n", string(output))
	}
}
func gdel(c *cli.Context) {
	list, _ := os.ReadDir(goDown)
	if len(list) == 0 {
		php.Color("No package\n", "")
		return
	}
	for i := range list {
		if err := os.RemoveAll(filepath.Join(goDown, list[i].Name())); err == nil {
			php.Color("Remove ok,"+list[i].Name()+"\n", "green")
		}
	}
}
func pall(c *cli.Context) {
	var vers []string
	if rs["pyver"] == nil {
		vers = getPyAllVer()
		rs["pyver"] = vers
		out, _ := yaml.Marshal(&rs)
		os.WriteFile("xv.yaml", out, 0666)
	} else {
		vers = any2str(rs["pyver"])
	}
	arr := make([][]string, 0)
	arr = php.SliceSplit(vers, 5)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"A", "B", "C", "D", "E"})
	table.AppendBulk(arr)
	table.Render()
}
func pls(c *cli.Context) {
	infos, err := os.ReadDir(pyVer)
	if err != nil || len(infos) <= 0 {
		php.Color("No installation", "")
	} else {
		php.Color("Is already installed\n", "green")
	}
	for i := range infos {
		if !infos[i].IsDir() {
			continue
		}
		vname := infos[i].Name()
		php.Color("√"+vname+"\n", "green")
	}
}
func pi(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	php.Color("下载太慢,推荐下载后放到download文件夹,然后安装,https://registry.npmmirror.com/binary.html?path=python/ \n", "")
	php.Color("2.7.1和3.9.5绿色版安装包,包含了pip,其它官方下载不包含哦,下载地址:https://gitee.com/yoby/xv/tree/master/zip \n", "")
	ver := c.Args()[0]                   //获取版本
	verDirs := filepath.Join(pyVer, ver) //版本目录
	if php.IsDir(verDirs) {              //已安装直接退出
		php.Color("√"+ver, "green")
		os.Exit(0)
	}
	vers := getPyAllVer() //获取版本列表
	if php.InArray(ver, vers) {
		urls, names := getUrlNamepy(ver, rs["purl"].(string))
		name := filepath.Join(pyDown, names) //判断文件是否存在
		if name == "" {                      //这个版本没有绿色版
			php.Color("Download successful\n", "green")
			os.Exit(0)
		}
		if !php.IsExist(name) { //不存在下载
			downpy(urls, name)
			php.Color("Download successful\n", "green")
		}

		os.RemoveAll(verDirs) //删除缓存的

		archiver.Unarchive(name, verDirs) //解压
		os.RemoveAll(pyInstall)
		os.Symlink(verDirs, pyInstall)
		php.Color("Set up OK\n", "green")
	} else {
		php.Color("No installation\n", "green")
	}
}
func pu(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                   //获取版本
	verDirs := filepath.Join(pyVer, ver) //版本目录
	if !php.IsExist(verDirs) {
		php.Color("Version does not exist.\n", "red")
		return
	}
	os.RemoveAll(verDirs)
	php.Color("Uninstall the success.\n", "green")
}
func puse(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                   //获取版本
	verDirs := filepath.Join(pyVer, ver) //版本目录
	if !php.IsExist(verDirs) {
		php.Color("Version does not exist.\n", "red")
		return
	}
	os.Remove(pyInstall)
	os.Symlink(verDirs, pyInstall)
	if output, err := exec.Command(filepath.Join(pyInstall, "python"), "-V").Output(); err == nil {
		fmt.Println("Switch success\n", string(output))
	}
}
func pdel(c *cli.Context) {
	list, _ := os.ReadDir(pyDown)
	if len(list) == 0 {
		php.Color("No package\n", "")
		return
	}
	for i := range list {
		if err := os.RemoveAll(filepath.Join(pyDown, list[i].Name())); err == nil {
			php.Color("Remove ok,"+list[i].Name()+"\n", "green")
		}
	}
}
func nall(c *cli.Context) {
	var vers []string
	if rs["nodever"] == nil {
		vers = getnodever()
		rs["nodever"] = vers
		out, _ := yaml.Marshal(&rs)
		os.WriteFile("xv.yaml", out, 0666)
	} else {
		vers = any2str(rs["nodever"])
	}
	ss := php.SliceSplit(vers, 5)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"A", "B", "C", "D", "E"})
	table.AppendBulk(ss)
	table.Render()
}
func nls(c *cli.Context) {
	infos, err := os.ReadDir(nodeVer)
	if err != nil || len(infos) <= 0 {
		php.Color("No installation", "")
	} else {
		php.Color("Is already installed\n", "green")
	}
	for i := range infos {
		if !infos[i].IsDir() {
			continue
		}
		vname := infos[i].Name()
		php.Color("√"+vname+"\n", "green")
	}
}
func ni(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                     //获取版本
	verDirs := filepath.Join(nodeVer, ver) //版本目录
	if php.IsDir(verDirs) {                //已安装直接退出
		php.Color("√"+ver, "green")
		os.Exit(0)
	}
	vers := getnodever() //获取版本列表
	if php.InArray(ver, vers) {
		urls, names := getNodeUrlName(ver, rs["nurl"].(string))
		name := filepath.Join(nodeDown, names) //判断文件是否存在
		if !php.IsExist(name) {                //不存在下载
			downnode(urls, name)
			php.Color("Download successful\n", "green")
		} else {
			sha2, _ := SHA256File(name)
			isok := getsha256(ver, sha2)
			if isok == false {
				os.Remove(name) //删除不完整文件
				php.Color("Download successful\n", "green")
				os.Exit(0)
			}
		}
		os.RemoveAll(verDirs)             //删除缓存的
		archiver.Unarchive(name, nodeVer) //解压
		file := filepath.Base(name)
		ext := path.Ext(name)
		os.Rename(filepath.Join(nodeVer, strings.TrimSuffix(file, ext)), verDirs)
		os.Remove(nodeInstall)
		os.Symlink(verDirs, nodeInstall)
		php.Color("Set up OK\n", "green")
	} else {
		php.Color("No installation\n", "green")
	}
}
func nu(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                     //获取版本
	verDirs := filepath.Join(nodeVer, ver) //版本目录
	if !php.IsExist(verDirs) {
		php.Color("Version does not exist.\n", "red")
		return
	}
	os.RemoveAll(verDirs)
	php.Color("Uninstall the success.\n", "green")
}
func nuse(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                     //获取版本
	verDirs := filepath.Join(nodeVer, ver) //版本目录
	if !php.IsExist(verDirs) {
		php.Color("Version does not exist.\n", "red")
		return
	}
	os.Remove(nodeInstall)
	os.Symlink(verDirs, nodeInstall)
	if output, err := exec.Command(filepath.Join(nodeInstall, "node"), "-v").Output(); err == nil {
		fmt.Println("Switch success\n", string(output))
	}

}
func ndel(c *cli.Context) {
	list, _ := os.ReadDir(nodeDown)
	if len(list) == 0 {
		php.Color("No package\n", "")
		return
	}
	for i := range list {
		if err := os.RemoveAll(filepath.Join(nodeDown, list[i].Name())); err == nil {
			php.Color("Remove ok,"+list[i].Name()+"\n", "green")
		}
	}
}
func fall(c *cli.Context) {
	arr := make([][]string, 0)
	var vers []string
	if rs["flutterver"] == nil {
		vers = getFlutterVer()
		rs["flutterver"] = vers
		out, _ := yaml.Marshal(&rs)
		os.WriteFile("xv.yaml", out, 0666)
	} else {

		vers = any2str(rs["flutterver"])
	}
	arr = php.SliceSplit(vers, 5)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"A", "B", "C", "D", "E"})
	table.AppendBulk(arr)
	table.Render()
}
func fls(c *cli.Context) {
	infos, err := os.ReadDir(flutterVer)
	if err != nil || len(infos) <= 0 {
		php.Color("No installation", "")
	} else {
		php.Color("Is already installed\n", "green")
	}
	for i := range infos {
		if !infos[i].IsDir() {
			continue
		}
		vname := infos[i].Name()
		php.Color("√"+vname+"\n", "green")
	}
}
func fi(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                        //获取版本
	verDirs := filepath.Join(flutterVer, ver) //版本目录
	if php.IsDir(verDirs) {                   //已安装直接退出
		php.Color("√"+ver, "green")
		os.Exit(0)
	}
	vers := getFlutterVer() //获取版本列表
	if php.InArray(ver, vers) {
		urls, names := getFlutterUrlName(ver, rs["furl"].(string))
		name := filepath.Join(flutterDown, names) //判断文件是否存在
		if !php.IsExist(name) {                   //不存在下载
			downflutter(urls, name)
			php.Color("Download successful\n", "green")
		} else {
			php.Color("Download successful\n", "green")
		}
		os.RemoveAll(verDirs)                //删除缓存的
		archiver.Unarchive(name, flutterVer) //解压
		os.Rename(filepath.Join(flutterVer, "flutter"), verDirs)
		os.Remove(flutterInstall)
		os.Symlink(verDirs, flutterInstall)
		php.Color("Set up OK\n", "green")
	} else {
		php.Color("No installation\n", "green")
	}
}
func fu(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                        //获取版本
	verDirs := filepath.Join(flutterVer, ver) //版本目录
	if !php.IsExist(verDirs) {
		php.Color("Version does not exist.\n", "red")
		return
	}
	os.RemoveAll(verDirs)
	php.Color("Uninstall the success.\n", "green")
}
func fuse(c *cli.Context) {
	if c.NArg() == 0 {
		php.Color("Input error", "red")
		os.Exit(0)
	}
	ver := c.Args()[0]                        //获取版本
	verDirs := filepath.Join(flutterVer, ver) //版本目录
	if !php.IsExist(verDirs) {
		php.Color("Version does not exist.\n", "red")
		return
	}
	os.Remove(flutterInstall)
	os.Symlink(verDirs, flutterInstall)
	if output, err := exec.Command(filepath.Join(flutterInstall, "bin", "flutter"), "--version").Output(); err == nil {
		fmt.Println("Switch success\n", string(output))
	}
}
func fdel(c *cli.Context) {
	list, _ := os.ReadDir(flutterDown)
	if len(list) == 0 {
		php.Color("No package\n", "")
		return
	}
	for i := range list {
		if err := os.RemoveAll(filepath.Join(flutterDown, list[i].Name())); err == nil {
			php.Color("Remove ok,"+list[i].Name()+"\n", "green")
		}
	}
}

// 路径替换
func pathx(str string) (s string) {
	if runtime.GOOS == "windows" {
		s = strings.Replace(str, "\\", "/", -1)
	} else {
		s = str
	}
	return s
}

// 版本排序
func verSort(vv []string) []string {
	versions := make([]*version.Version, len(vv))
	for i, raw := range vv {
		v, _ := version.NewVersion(raw)
		versions[i] = v
	}
	sort.Sort(version.Collection(versions))
	ver := make([]string, len(vv))
	for i, s := range versions {
		ver[i] = s.Original()

	}
	return ver
}

// any转换[]string
func any2str(a any) []string {
	strings, _ := a.([]interface{})
	array := make([]string, len(strings))
	for i, v := range strings {
		array[i] = v.(string)
	}
	return array
}

// 获取go所有版本
func getGoAllVer() ([]string, map[string]string) {
	var str string
	str = getdoc()
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
	var sc []string
	sha1x := make(map[string]string)
	doc.Find(`.toggleVisible,.toggle`).Each(func(i int, s *goquery.Selection) {
		ss, _ := s.Attr("id")
		if ss != "archive" && i < 30 { //只取最近30个版本
			ss = string([]byte(ss)[2:]) //去掉go两个字符
			sc = append(sc, ss)
			_, ver := getUrlName(ss, rs["gurl"].(string))
			sha1x[ss] = doc.Find(`a:contains("` + ver + `")`).Parent().Parent().Find("tt").Text()
		}
		i++
	})

	return sc, sha1x

}

// 获取go文档
func getdoc() string {
	res, _ := http.Get(rs["gurl"].(string))
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("status code error: %d %s", res.StatusCode, res.Status)
	}
	body, _ := io.ReadAll(res.Body)
	str := string(body)
	return str
}

// go下载链接
func getUrlName(ver, u string) (url, name string) {
	var ext string
	if runtime.GOOS == "windows" {
		ext = "zip"
	} else {
		ext = "tar.gz"
	}
	name = "go" + ver + "." + runtime.GOOS + "-" + runtime.GOARCH + "." + ext
	url = u + "/golang/" + name
	return url, name
}

// 下载go
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

// 验证go
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

// 获取py版本
func getPyAllVer() []string {
	var str string
	res, _ := http.Get(rs["purl"].(string))
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("Error:%d %s", res.StatusCode, res.Status)
	}
	body, _ := io.ReadAll(res.Body)
	str = string(body)

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(str))
	var sc []string
	var a1, b1 int
	doc.Find(`a`).Each(func(i int, s *goquery.Selection) {
		ss := strings.TrimSuffix(s.Text(), "/")
		if strings.Contains(ss, "binaries-1.1") {
			a1 = i
		}
		if ss == "3.0.1" {
			b1 = i
		}
		sc = append(sc, ss)
	})
	sc1 := sc[b1:a1]
	ver := verSort(sc1)
	ver = ver[len(ver)-29:]
	ver = append([]string{"2.7.1"}, ver...)
	return ver
}

// 获取py版本名称
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

// 下载py
func downpy(u, n string) {
	req, _ := http.NewRequest("GET", u, nil)
	fmt.Println("Download Url:" + u + "\n")
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	types := resp.Header.Get("Content-Type")
	if types == "text/html" {
		fmt.Println("Zip file does not exist!")
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

// 获得node版本
func getnodever() []string {
	var str string
	str = php.Get(rs["nurl"].(string) + "/index.json")
	var ver []map[string]string
	json.Unmarshal([]byte(str), &ver)
	var a []string
	for k, v := range ver {
		if k < 150 {
			a = append(a, strings.TrimLeft(v["version"], "v"))
		}
	}
	a = verSort(a)
	return a
}

// 获取node下载名称
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

// 下载node
func downnode(u, n string) {
	wget(u, n)
}

// 验证noide
func getsha256(ver string, s256 string) (s bool) {
	u := rs["nurl"].(string) + "/v" + ver + "/SHASUMS256.txt"
	str := php.Get(u)
	s = strings.Contains(str, s256)
	return s
}

// 获取flutter下载名称
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

// 下载flutter
func downflutter(u, n string) {
	wget(u, n)
}
func getFlutterVer() []string {
	var str string
	str = php.Get(rs["furl"].(string) + "/releases_windows.json")
	data := []byte(str)
	var a []string
	r, _ := sonic.Get(data, "releases")
	r1, _ := r.Array()
	for k, v := range r1 {
		v1 := v.(map[string]interface{})
		ch := v1["channel"]
		if ch == "stable" {
			ver, _ := v1["version"]
			if k > 43558 {
				return nil
			} else {
				a = append(a, ver.(string))
			}

		}
	}
	a = verSort(a)
	a = a[len(a)-30:]
	return a
}

// 通用下载
func wget(u, n string) {
	req, _ := http.NewRequest("GET", u, nil)
	php.Color("Download Url:"+u+"\n", "green")
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
