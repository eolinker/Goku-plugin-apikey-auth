# Goku Plugin：API Key

| 插件名称  | 文件名.so |  插件类型  | 错误处理方式 | 作用范围 |  优先级  |
| ------------ | ------------ | ------------ | ------------ | ------------ | ------------ |
| APIKey鉴权  | goku-apikey_auth | 访问策略 | 继续后续操作 | 转发前  |  1003  |

鉴权方式的一种，多用于OpenAPI，设置Apikey参数，Apikey默认支持在header、body、query中使用，不能通过认证的用户将无权访问接口。

# 目录
- [编译教程](#编译教程 "编译教程")
- [安装教程](#安装教程 "安装教程")
- [使用教程](#使用教程 "使用教程")
- [更新日志](#更新日志 "更新日志")

# 编译教程

#### 环境要求
* 系统：基于 Linux 内核（2.6.23+）的系统，CentOS、RedHat 等均可；

* golang版本号：12.x及其以上

* 环境变量设置：
	* GO111MODULE：on
	
	* GOPROXY：https://goproxy.io


#### 编译步骤

1.clone项目

2.进入项目文件夹，执行**build.sh**
```
cd goku-apikey_auth && chmod +x build.sh && ./build.sh
```

###### 注：build.sh为通用的插件编译脚本，自定义插件时可以拷贝直接使用。

3.执行第2步将会生成文件： **{插件名}.so**

将该文件上传到**节点服务器运行目录**下的**plugin**文件夹，然后在控制台安装插件即可使用。

# 安装插件
前往 Goku API Gateway 官方网站查看：[插件安装教程](url "https://help.eolinker.com/#/tutorial/?groupID=c-341&productID=19")

# 使用教程

#### 配置页面

进入控制台 >> 策略管理 >> 某策略 >> 策略插件 >> Apikey鉴权插件：

![](http://data.eolinker.com/course/U7UhASH318687437c2efcae19a6a1a8ac740c724f3b15b0)

#### 配置参数

| 参数名 | 说明   | 
| ------------ | ------------ |  
|  Apikey | Apikey的值 | 
| hildCredential  | 转发时是否隐藏Apikey |
| remark  | 备注 |

#### 配置示例

```
[
    {
        "Apikey": "key",
        "hideCredential": false
        "remark": ""
    },
    {
        "Apikey": "key2",
        "hideCredential": false
        "remark": ""
    }
]
```

#### API请求参数

| 参数名 | 说明  | 必填  |   值可能性   |  参数位置 |
| ------------ | ------------ |  ------------ |  ------------ |  ------------ |  
|  Strategy-Id | 策略id | 是 |   |  header  | 
|  Authorization-Type  | 鉴权方式 | 是 | Apikey  | header  |
| Authorization  |  鉴权值 |  是  |    | header |


#### 请求示例

###### 以下API测试页面来自于 **eoLinker API Studio** 接口管理平台

* Apikey在 **header** 中：

    ![](http://data.eolinker.com/course/V8i8viHdf6dc9648a19ec778ce5698af643cefa01d3fb2e)

* Apikey在 **body** 中：
头部填写 **Strategy-Id** 与 **Authorization-Type**:
![](http://data.eolinker.com/course/awEChej35aa995b7613f888ed9d0816cb3ff1eddd4cc3c2)

* Apikey在 **query** 中：
头部填写 **Strategy-Id** 与 **Authorization-Type**:
![](http://data.eolinker.com/course/M85um6g07177251dbdecd1593cc8b7306252e7c589f7653)
