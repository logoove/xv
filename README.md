# xv

#### 介绍
一个支持go,python,nodejs,flutter的版本下载管理工具,不再需要安装多个工具,支持Linux,Windows,MacOS.

### 独有特点

- 使用缓存,程序初次运行释放配置文件xv.yaml在运行目录,所有版本号获取都缓存存到配置文件.
- 中英文支持,根据系统自动显示,一般本地win中文,linux服务器英文.
- 支持普通下载文件使用`xv get url`如果是github的zip会自动使用加速下载.
- 友好的提示,这是一个学习go作品,参考了其他类似工具,主要是解决go/py/node/flutter管理版本.
- win下运行需要管理员权限的cmd
- 命令都是每个语言首字母+命令组成,尽量简化到输入简单.
### 参考命令
~~~
gls    列出go已安装版本
gall   列出go可安装版本
gi     安装go指定版本
gu     卸载go指定版本
guse   设置go默认版本
gdel   删除go所有安装包

pls    列出python已安装版本
pall   列出python可安装版本
pi     安装python指定版本
pu     卸载python指定版本
puse   设置python默认版本
pdel   删除python所有安装包

nall   列出node可安装版本
nls    列出node已安装版本
ni     安装node指定版本
nu     卸载node指定版本
nuse   设置node默认版本
ndel   删除node所有安装包

fall   列出flutter可安装版本
fls    列出flutter已安装版本
fi     安装flutter指定版本
fu     卸载flutter指定版本
fuse   设置flutter默认版本
fdel   删除flutter所有安装包

get   下载安装包到桌面

xv gls //列出已安装go版本
xv gall //列出可安装go版本
xv gi 1.19 //安装go
xv gu 1.19 //卸载go
xv guse 1.19 //设置默认go版本
xv gdel //清空下载目录
xv -v //查看版本
xv gls -h 查看这个命令帮助和例子
~~~

### Go

- win 设置环境变量
- GOPATH=C:\www\go //这是工作目录,一般在里面建立src,bin,pkg三个文件夹,src下面放我们项目文件夹
- GOPROXY=https://goproxy.cn //设置代理,国内必须设置否则下载太慢了
- GOROOT=C:\app\go\go  //go语言的当前版本路径,非常重要
- GOENV=C:\app\go\env
- GOCACHE=C:\app\go\go-build
- GO111MODULE=on
- path里面加入 `C:\app\go\go\bin` `C:\www\go\bin` 这里只是个例子主要是能够任何目录下使用go 命令,当然我们也能加入xv.exe到这两个任意一个里面,cmd下都能执行.

- linux,设置环境变量
  vim /etc/profile 末尾添加
~~~
export GOROOT=/root/app/go/go
export PATH=$PATH:$GOROOT/bin
export GOPATH=/usr/gowork #项目目录般在里面建立src,bin,pkg三个文件夹
export PATH=$PATH:$GOPATH/bin
export GOPROXY=https://goproxy.cn,direct
~~~
- 上传xv到 `/usr/bin`,权限改成777,这样任何地方都能执行.,也可以建立软链接,都一样

### node
win 设置

- 输入npm config lis找到 .npmrc配置
~~~
  registry=http://registry.npmmirror.com
  prefix=C:\app\nodejs\node-global
  cache=C:\app\nodejs\node-cache
  python=C:\app\python\version\2.7.1\python.exe
~~~
常见模块
~~~
npm install -g cnpm --registry=http://registry.npmmirror.com
npm install -g  less
npm install -g  sass
npm install -g express
npm install webpack -g
npm install webpack-cli -g
npm i @vue/cli -g
npm i element-ui -g
~~~

### python

- win设置环境变量,python只需要设置path即可.
- path里面加入 `C:\app\python\python`,如果需要多版本共存,比如2.7.1和3.x那需要把`C:\app\python\version\2.7.1`也加入path.
- 在线是不能下载2.7.1的,这个版本官方没有绿色版,所以单独制作了2.7.1的绿色版,只需要把压缩包放到`C:\app\python\download`下面即可
- 已经自带pip了,需要设置环境变量路径`C:\app\python\python\Scripts`
- 推荐使用制作的3.9.5绿色版包含pip,官方安装的zip不含有.
- Linux因为系统自带,所以没必要使用此工具管理.

### flutter
win 设置环境变量
~~~
export PUB_HOSTED_URL=https://pub.flutter-io.cn
export FLUTTER_STORAGE_BASE_URL=https://storage.flutter-io.cn
CHROME_EXECUTABLE edge路径
~~~

安装Android studio/IDEA
下载android studio https://developer.android.google.cn/studio/#downloads 下载zip版本即可 解压到 C:\app\android-studio
接着安装flutter和Dart插件
1 . JAVA_HOME环境变量设置 C:\app\android-studio\jre
安装jdk1.8(jdk8),其他版本不支持

ANDROID_HOME环境变量 C:\app\Androidsdk
path设置 C:\app\flutter\bin 然后运行flutter doctor
检测显示 Android license status unknown
运行 flutter doctor --android-licenses 一路yes
安装Android SDK,连上安卓手机下载对应sdk,在真机上测试

### 感谢或使用的参考
- php,一个与php同名函数工具库<https://github.com/logoove/go/php>
- cli,本工具所使用的构建工具<https://github.com/logoove/cli>
- 另一个go版本管理工具<https://github.com/voidint/g>
- node 版本工具 <https://github.com/nvm-sh/nvm>
- 系统语言获取库<https://github.com/Xuanwo/go-locale>
- 判断是否双击打开,给与提示,而不是闪一下没了<https://github.com/inconshreveable/mousetrap>

### 其他作品
- weuiplus 一个移动端UI库,可开发公众号和一般移动端.<https://github.com/logoove/weui>
- php,一个函数工具库,不使用任何第三方库和扩展库只用基础库写成<https://github.com/logoove/go/php>
- cli 命令行工具构建工具<https://github.com/logoove/cli>
- xv 支持go,python,nodejs,flutter版本管理器 <https://github.com/logoove/xv>
- sqlite 一个不使用cgo的database/sql标准sqlite3驱动<https://github.com/logoove/sqlite>