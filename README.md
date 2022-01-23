

# 一、设计要求

1. **语言**：熟悉Golang语言的基本语法，包括其反射、多协程并发、文件处理、接口等高级特性；了解前端三件套HTML、CSS、JavaScript的基本语法。
2. **网络**：熟悉计算机网络传输层的TCP/UDP协议、应用层的HTTP/WebSocket协议的相关API函数的功能与基本使用。
3. **前端**：熟悉渐进式前端框架Vue、了解前端组件库MUI的基本使用。
4. **后端**：熟悉后台web开发的基本原理与使用，如对象关系映射orm框架的使用、Gin路由框架的使用、缓存中间件的使用。
5. **数据库**：熟悉关系型数据库MySQL的底层原理、库表设计、sql语句的编写、性能优化等。
6. **协同开发**：熟悉版本控制工具Git的使用。
7. **设计模式**：熟悉单例模式、接口隔离原则、简单工厂模式、单一指责原则等软件设计模式的思想。
8. **架构思想**：了解的MVC、MVVM架构模式的思想与设计。



# 二、开发环境与工具

+ **开发语言**：Golang + HTML + CSS + JavaScript
+ **开发工具**：Goland(代码编写IDE) + Sequal Ace(数据库可视化软件)
+ **调试测试**：Edge(浏览器调试查看详细网络请求与响应)
+ **协同开发**：项目的开发采用git作为版本控制工具，所有代码实时同步到远程代码仓库Gitee



# 三、设计原理

+ **前端** 采用了渐进式框架Vue + 开源组件库MUI。之所以采用Vue框架的原因是为了实现前后端分离，前端后端并不需要在一个项目目录下相互耦合，可以独立开发独立部署，最后通过接口的形式从后端向前端传递Json数据即可。

+ **后端** 采用了websocket、udp相关的api进行实现，采用websocket协议的原因是因为http协是一个非持久化协议，一次请求与响应之后连接就关闭，若要进行下次的交互则需要重新建立连接，耗时耗资源；此外，http请求头占了很大一部分，而真正的数据部分只有很小一部分，这样也会浪费很多的带宽资源。虽然http1.1之后可以通过keepalive字段来实现半持久化连接，但是比较麻烦。本项目是即时通讯系统，对实时性的要求比较高，显然http协议并不适用，而websocket协议很好的解决了http协议出现的问题，它使得客户端和服务器之间的数据交换变得更加简单，允许服务端主动向客户端推送数据。在websocket api中，浏览器和服务器只需要完成一次握手，两者之间就直接可以创建持久性的连接，并进行双向数据传输。因此它只需要较少的控制开销，但是有更强的实时性。此外，为了实时性，在传输层方面本项目使用了udp协议，而并非tcp协议。

+ **数据库** 采用了MySQL，采用Xorm对象关系映射框架与其进行交互。MySQL是最常用的关系型数据库，Xorm是用go语言实现的一个orm框架，能够便捷的与数据库进行交互。

+ **代码层级设计**：采用了MVVM模式和MVC模式，其中MVC可以理解为整个项目的架构分层模式，而MVVM模式是前端采用的一套分层设计，两种模式的基本介绍如下所示：

    1. MVC是Model-View-Controller的缩写，它将应用程序划分为三个部分：

       Model模型层: 管理数据

       View视图层: 展示数据

       Controller控制层: 处理数据

       ![image-20220119215030731](https://gitee.com/bareth/images2/raw/master//img/image-20220119215030731.png)

    2. MVVM是Model-View-ViewModel的缩写，是为了解决前端的响应式编程而生，由于前端网页混合了HTML、CSS和JavaScript，而且页面众多，代码的组织和维护难度复杂，所以通过ViewModel实现View和Model的双向绑定。

       ![img](https://gitee.com/bareth/images2/raw/master//img/wps7B1rqz.jpg)

+ **代码层级目录介绍**：

  <img src="https://gitee.com/bareth/images2/raw/master//img/image-20220119215120035.png" alt="image-20220119215120035" style="zoom: 50%;" /> 

+ **项目的完整架构设计**如下图所示：

  ![img](https://gitee.com/bareth/images2/raw/master//img/wpsxqHqRm.jpg)



# 四、系统功能

如上述项目架构图所示，系统功能包括：

1. 用户注册
2. 用户登陆
3. 用户信息修改
4. 添加好友
5. 创建群聊
6. 添加群聊


# 五、系统实现

## 5.1、数据库设计

1. community群聊表

![img](https://gitee.com/bareth/images2/raw/master//img/wpsqiAXs3.jpg)

```go
// Community 群聊结构体
type Community struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"` // 群聊id
	Name     string    `xorm:"varchar(30)" form:"name" json:"name"`        // 群聊名称
	Ownerid  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid"`   // 群聊创建者
	Icon     string    `xorm:"varchar(250)" form:"icon" json:"icon"`       // 群聊头像
	Cate     int       `xorm:"nt(11)" form:"cate" json:"cate"`             // 群聊类型
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`       // 群聊描述
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`   // 群聊创建时间
} 
```

2. contact聊天会话表

![img](https://gitee.com/bareth/images2/raw/master//img/wps0dYH5i.jpg)

```go
// Contact 会话结构体
type Contact struct {
	Id       int64     `xorm:"pk autoincr bigint(20)" form:"id" json:"id"` // 聊天会话id
	Ownerid  int64     `xorm:"bigint(20)" form:"ownerid" json:"ownerid"`   // 聊天发起者id
	Dstid    int64     `xorm:"bigint(20)" form:"dstid" json:"dstid"`       // 聊天对话者id
	Cate     int       `xorm:"int(11)" form:"cate" json:"cate"`            // 聊天类型(好友/群聊)
	Memo     string    `xorm:"varchar(120)" form:"memo" json:"memo"`       // 聊天备注
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`   // 聊天会话创建时间
}
 
```

3. user用户表

![img](https://gitee.com/bareth/images2/raw/master//img/wpsR5Zqml.jpg)

```go
// User 用户结构体
type User struct {
	Id       int64     `xorm:"pk autoincr bigint(64)" form:"id" json:"id"`  // 用户id(主键自增)
	Mobile   string    `xorm:"varchar(20)" form:"mobile" json:"mobile"`     // 用户电话号码
	Password string    `xorm:"varchar(40)" form:"password" json:"-"`        // 用户密码
	Avatar   string    `xorm:"varchar(150)" form:"avatar" json:"avatar"`    // 用户头像
	Sex      string    `xorm:"varchar(2)" form:"sex" json:"sex"`            // 用户性别
	Nickname string    `xorm:"varchar(20)" form:"nickname" json:"nickname"` // 用户昵称
	Salt     string    `xorm:"varchar(10)" form:"salt" json:"-"`            // 用户盐度加密(首次注册需要)
	Online   int       `xorm:"int(10)" form:"online" json:"online"`         // 用户是否在线
	Token    string    `xorm:"varchar(40)" form:"token" json:"token"`       // 用户登陆token
	Memo     string    `xorm:"varchar(140)" form:"memo" json:"memo"`        // 用户备注
	Createat time.Time `xorm:"datetime" form:"createat" json:"createat"`    // 创建时间
}
```



## 5.2、后端实现

所有的路由函数如下图所示，每个函数的功能都有完整注释

![image-20220119215958620](https://gitee.com/bareth/images2/raw/master//img/image-20220119215958620.png)

### 5.2.1、用户登陆

用户登陆controller层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wps439Kxd.jpg)

用户登陆service层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsN92A15.jpg)

### 5.2.2、用户注册

用户注册controller层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpstvKd0P.jpg)

用户注册service层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsAmhUZL.jpg)

### 5.2.3、用户信息修改

用户信息修改controller层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpssOs06L.jpg)

用户信息修改service层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsrnCJWo.jpg)

### 5.2.4、查找指定用户

查找指定用户controller层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsThGvsX.jpg)

查找指定用户service层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wps8Nk2zf.jpg)

### 5.2.5、加载群聊信息

加载群聊信息controller层代码

![image-20220119220905716](https://gitee.com/bareth/images2/raw/master//img/image-20220119220905716.png)

加载群聊信息service层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsd9Ozd5.jpg)

### 5.2.1、创建群聊

创建群聊controller层代码：

![image-20220119220944377](https://gitee.com/bareth/images2/raw/master//img/image-20220119220944377.png)

创建群聊service层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsnNCYou.jpg)

### 5.2.7、加入群聊

加入群聊controller层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpscB4K5O.jpg)

加入群聊service层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsKUQcry.jpg)

### 5.2.8、加载所有好友

加载所有好友controller层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsAVlEff.jpg)

加载所有好友service层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsJ1jyC5.jpg)

### 5.2.9、添加好友

添加好友controller层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsXoRJ75.jpg)

添加好友service层代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsZipRoU.jpg)

### 5.2.10、聊天

**关键设计**：会话结构体Conversation、全局clientMap对象(存放userid与conversation的映射关系)

1. 会话conversation结构体

![img](https://gitee.com/bareth/images2/raw/master//img/wpswghfX1.jpg)

2. 全局变量

![img](https://gitee.com/bareth/images2/raw/master//img/wpskBEqOF.jpg)

聊天controller层代码：

![image-20220119221512237](https://gitee.com/bareth/images2/raw/master//img/image-20220119221512237.png)

发送逻辑代码：

![image-20220119221534215](https://gitee.com/bareth/images2/raw/master//img/image-20220119221534215.png)

接收逻辑代码：

![img](https://gitee.com/bareth/images2/raw/master//img/wpsmqgEFT.jpg)

发送消息代码：

![image-20220119221620600](https://gitee.com/bareth/images2/raw/master//img/image-20220119221620600.png)

后端通用调度逻辑代码：

![image-20220119222529101](https://gitee.com/bareth/images2/raw/master//img/image-20220119222529101.png)

## 5.3、功能截图

01_用户登陆

![image-20220119224457458](https://gitee.com/bareth/images2/raw/master//img/image-20220119224457458.png)

02_用户注册

![img](https://gitee.com/bareth/images2/raw/master//img/wpsAfL6MX.jpg)

03_个人中心

![2](https://gitee.com/bareth/images2/raw/master//img/2.png)

04_添加好友

![](https://gitee.com/bareth/images2/raw/master//img/image-20220119223951993.png)

05_加入群聊

![img](https://gitee.com/bareth/images2/raw/master//img/wpsxjyItT.jpg)

06_修改个人信息

![img](https://gitee.com/bareth/images2/raw/master//img/wpsVaeiYU.jpg)

07_好友列表

![img](https://gitee.com/bareth/images2/raw/master//img/wpsRYWBaK.jpg)

08_好友聊天(可以发送表情包和图片)

![image-20220119224134280](https://gitee.com/bareth/images2/raw/master//img/image-20220119224134280.png)

09_群聊列表

![image-20220119224213834](https://gitee.com/bareth/images2/raw/master//img/image-20220119224213834.png)

10_群聊页面

![img](https://gitee.com/bareth/images2/raw/master//img/wpsxuJqqg.jpg)



# 六、问题及其解决方案

**1. 用户接入与鉴权**

用户注册时密码经过md5加密，且登陆需要携带加密token，根据token校验结果来决定是否登陆

![img](https://gitee.com/bareth/images2/raw/master//img/wpsyZHCni.jpg)

**2. Conn的维护**

聊天分为私聊和群聊两种方式，一个用户可能与多个用户建立连接，因此通过一个会话结构体来	维持用户的conn，然后创建一个clientMap用来存放每个用户与其对应的会话

![img](https://gitee.com/bareth/images2/raw/master//img/wps5CUirH.jpg)

**3. 并发修改异常问题**

由于多个用户可能对同一个clientMap进行操作，可能会出现并发修改异常问题，所以用读写锁	来保证并发情况下的原子性

![img](https://gitee.com/bareth/images2/raw/master//img/wpsCkPpX2.jpg)

**4. 数据库的安全修改**

由于存在多用户对数据库的读写进行操作，为了避免修改冲突，对数据库对写操作采用事物进行控制

![1](https://gitee.com/bareth/images2/raw/master//img/1.png)
