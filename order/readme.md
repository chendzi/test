#前提说明

1.电脑里面没有环境，所以从昨晚8点开始从头装的go语言安装包，配置go语言环境，liteide编译软件，mysql，mysql可视化程序，postman测试工具，openssl证书生成工具等，
家里网还不好，所以匆忙中，可能功能上实现的比较粗糙。请您谅解。

#功能说明：
※目前没有完成的是docker容器化的测试，windows下的docker使用的不太熟练，我只写了一个linux下的dockerfile，但是我还没有linux环境。。。
docker容器化命令:docker build -t ${IMAGE_NAME} -f Dockerfile .

目前我在windows下编译后生成的exe可以实现要求的两个接口，证书是用openssl生成的。

#测试说明：

测试需要安装mysql和postman，鉴于没有的话安装比较麻烦，可以先参考pic中的测试截图，我可以带电脑演示.
※如果电脑中恰巧有以上两个软件，可导入env文件夹下的文件用来测试。

测试方法：

*导入mysql与postman数据后，postman中会出现两个接口的post请求，mysql会出现database test，test下会有一个order表。
*双击启动order.exe，点击postman中相应接口的发送https请求，即可看到返回值。

#代码结构
1.代码分三个模块:servicecert(证书包)、handle(接口处理部分)、model(mysql数据库部分)。、
pic里面放了我用postman测试的截图
env里面放了mysql的导出表和postman的导出表
