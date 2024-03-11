## webrtc_learn

### 这个项目做了什么？

这个项目融合了example里我觉得最有助于理解这个框架的两个例子

1. 浏览器使用摄像头捕捉信息->传到后端->后端再次传到浏览器
2. 浏览器和后端的data交互(本例互相发了几个message)

### 怎么把这个项目跑起来？

克隆仓库

```
git clone https://github.com/jiuxia211/webrtc_learn
```

打开这个网页

[Edit fiddle - JSFiddle - Code Playground](https://jsfiddle.net/nsm8ovjt/)

左边的`html`和`JavaScript`的内容我暂存在文件夹里了(我只能看懂大概做了什么，需要前端同学研究深一点),不过用空白的jsfiddle cv进去跑好像要不到权限

如果你是第一次跑，浏览器应该会向你请求摄像头的权限，点击同意

如果右边Logs下没有报错信息，正常你的`Browser base64 Session Description`内会有一个**信令**

（e开头的一长串）（这个东西不是很稳定，有时候会拿不到，可以考虑换浏览器试试）

在你的项目文件夹下创建一个名叫file的文件,将信令copy进去

然后执行

```
go run main.go < file
```

正常你的终端会输出一大串,这也是**信令**，把他copy到浏览器的`Golang base64 Session Description`内

点击`Start Session`

你就可以看到你的狗头啦！

你还可以看到Logs内看到后端传来的消息，你可以在main.go里找到对应的实现

点击`Send message`你的后端就可以接收到这条消息

### 跑成功之后

这时，你再去看main.go里的代码，就会清楚许多，我尽量做了详细的注释，但知识盲区还是很多的

然后前端同学主要研究一下浏览器是怎么做的。

