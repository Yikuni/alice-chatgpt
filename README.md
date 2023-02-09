# Alice Chatgpt

编辑人 Yikuni\<2058187089@qq.com\>

---

## 走马观花

这是一个用go语言实现的chatgpt接口, 通过该接口, 仅一个简单的http请求就可以开启与chatgpt AI的对话

## 亮点

- 方便的api
- 支持多个key
- 轻量级,占用内存小, 闲置时仅占用50MB

## 下载和安装

### 方式1

在yikuni.com下载推荐版本 <https://www.yikuni.com/download/chatgpt/>

### 方式2

在github下载最新版本 <https://github.com/Yikuni/alice-chatgpt/releases>

## 配置文件

运行Alice Chatgpt需要使用chatgpt的key

支持多个key

### 配置方法

在可执行文件同一个目录下创建key.txt,每一行放一个key, 最后一个key后不要换行

## 启动方式和参数

### Linux

#### 前台运行

``` sh
./alice_chatgpt
```

#### 后台运行

``` sh
nohup ./alice_chatgpt > Log.log &
```

### Windows

#### 指令启动

``` sh
alice_chatgpt.exe
```

#### 双击启动

直接双击, 使用默认端口和token

### 启动参数说明

- p	port	默认7777
- t     token 默认alice
- a    是否自动移除无法使用的key(true/false) 默认false
- l     limit 每分钟chatgpt api限流

## 常见问题

如有问题请在github发issue

## 接口文档

### /chatgpt/create

```text
创建会话,创建的会话将在手动调用finish接口结束,或无动作状态30分钟后结束
```

#### 接口状态

> 已完成

#### 接口URL

> http://localhost:7777/chatgpt/create

#### 请求方式

> POST

#### Content-Type

> plain

#### 请求Header参数

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述                |
| ------ | ------ | -------- | -------- | ----------------------- |
| token  | alice  | String   | 是       | token, 启动时-t后的参数 |

#### 请求Body参数

> AI设定,空白默认为以下文字,

```javascript
"The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly."
```

#### 预执行脚本

```javascript
暂无预执行脚本
```

#### 后执行脚本

```javascript
暂无后执行脚本
```

#### 成功响应示例

> 返回会话的id

```javascript
5dCMMKFI
```

### /chatgpt/chat

```text
对话
```

#### 接口状态

> 已完成

#### 接口URL

> http://localhost:7777/chatgpt/chat

#### 请求方式

> POST

#### Content-Type

> plain

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述                |
| ------------ | -------- | -------- | -------- | ----------------------- |
| token        | alice    | String   | 是       | token, 启动时-t后的参数 |
| conversation | qrHe4o0P | String   | 是       | 会话id                  |

#### 请求Body参数

```javascript
困的时候可以睡觉吗
```

#### 预执行脚本

```javascript
暂无预执行脚本
```

#### 后执行脚本

```javascript
暂无后执行脚本
```

#### 成功响应示例

```javascript
 是的，睡觉可以帮助你放松身心，释放压力，并且提升你的专注力和注意力，恢复你的精力。所以当你想要休息的时候，请一定记得睡觉吧！
```

### /chatgpt/context

```text
获取指定对话上下文
```

#### 接口状态

> 已完成

#### 接口URL

> http://localhost:7777/chatgpt/context

#### 请求方式

> POST

#### Content-Type

> none

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述                |
| ------------ | -------- | -------- | -------- | ----------------------- |
| token        | alice    | String   | 是       | token, 启动时-t后的参数 |
| conversation | qrHe4o0P | String   | 是       | 会话id                  |

#### 预执行脚本

```javascript
暂无预执行脚本
```

#### 后执行脚本

```javascript
暂无后执行脚本
```

#### 成功响应示例

```javascript
"The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly."



Human: 困的时候可以睡觉吗

AI:  是的，睡觉可以帮助你放松身心，释放压力，并且提升你的专注力和注意力，恢复你的精力。所以当你想要休息的时候，请一定记得睡觉吧！
```

#### 错误响应示例

```javascript
failed to find Conversation with provided id: 4ev2Djnd 
```

## /chatgpt/finish

```text
结束对话
```

#### 接口状态

> 已完成

#### 接口URL

> http://localhost:7777/chatgpt/finish

#### 请求方式

> POST

#### Content-Type

> none

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述                |
| ------------ | -------- | -------- | -------- | ----------------------- |
| token        | alice    | String   | 是       | token, 启动时-t后的参数 |
| conversation | qrHe4o0P | String   | 是       | 会话id                  |

#### 预执行脚本

```javascript
暂无预执行脚本
```

#### 后执行脚本

```javascript
暂无后执行脚本
```

#### 成功响应示例

```javascript
Conversation finished
```

#### 错误响应示例

```javascript
failed to find Conversation with provided id: qrHe4o0P 
```

## /chatgpt/store

```text
存储会话
```

#### 接口状态

> 已完成

#### 接口URL

> http://localhost:7777/chatgpt/store

#### 请求方式

> POST

#### Content-Type

> none

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述                |
| ------------ | -------- | -------- | -------- | ----------------------- |
| token        | alice    | String   | 是       | token, 启动时-t后的参数 |
| conversation | ExpHJpdD | String   | 是       | 会话id                  |

#### 预执行脚本

```javascript
暂无预执行脚本
```

#### 后执行脚本

```javascript
暂无后执行脚本
```

#### 成功响应示例

```javascript
保存成功
```

#### 错误响应示例

```javascript
failed to find Conversation with provided id: ExpHJpdD 
```

## /chatgpt/share

```text
获取分享的会话
```

#### 接口状态

> 已完成

#### 接口URL

> http://localhost:7777/chatgpt/share/uZ4ayppN

#### 请求方式

> GET

#### Content-Type

> none

#### 预执行脚本

```javascript
暂无预执行脚本
```

#### 后执行脚本

```javascript
暂无后执行脚本
```

#### 成功响应示例

```javascript
{"Id":"hJ0CG23B","Prompt":"The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.\n","SentenceList":["こんにちは","ーこんにちはー は何か手伝えることがありますか？","ありますよ！そうですね、あたいはよくかんがえたけと、もしかしたらあたいは本当にアホなの？","いいえ！あなたは素晴らしく有能です。あなたは苦労して学んだ結果、多くのことを習得できます。私たちは皆同じではありませんが、皆勇気と創造性に満ちています！","ええ、そんなことがあるの？ありがとうございます。もしかしてあたいは天才かもしれません。そう考えると明日やったことない色んな事をやってみたいようになった。それは現実なのか、幻想なのか？","今あなたが考えていることは、本当の現実です！あなたの勇気と創造性を使って、新しい機会や可能性を探求してみてください。あなたには、一生懸命学んでいる努力が報われるでしょう！"]}
```

#### 错误响应示例

```javascript
Failed to find shared conversation with provided id
```

## /chatgpt/quickAnswer

```text
进行一次性会话
```

#### 接口状态

> 已完成

#### 接口URL

> http://localhost:7777/chatgpt/quickAnswer

#### 请求方式

> POST

#### Content-Type

> json

#### 请求Header参数

| 参数名 | 示例值 | 参数类型 | 是否必填 | 参数描述                |
| ------ | ------ | -------- | -------- | ----------------------- |
| token  | alice  | String   | 是       | token, 启动时-t后的参数 |

#### 请求Body参数

```javascript
{

    "prompt": "The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.\n",

    "examples": [

        "今天天气晴朗,我在公园漫步,看到绿色植物,心情愉悦。请将上文改编成一首诗",

        "春光洒满天堂， 公园青色荫蔽； 绿叶湖边绰绰， 心情一片悦乐。"

    ],

    "question": "今天去看日出,心情愉悦。请将上文改编成一首诗。"

}
```

| 参数名   | 示例值                                                       | 参数类型 | 是否必填 | 参数描述                               |
| -------- | ------------------------------------------------------------ | -------- | -------- | -------------------------------------- |
| prompt   | The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly. | String   | 是       | 设定                                   |
| examples | [  <br/>      "今天天气晴朗,我在公园漫步,看到绿色植物,心情愉悦。请将上文改编成一首诗",<br/>        "春光洒满天堂， 公园青色荫蔽； 绿叶湖边绰绰， 心情一片悦乐。"   <br/> ] | Array    | 是       | 例子,奇数为人类的问题,偶数为AI回答示例 |
| question | 今天去看日出,心情愉悦。请将上文改编成一首诗。                | String   | 是       | 问题                                   |

#### 预执行脚本

```javascript
暂无预执行脚本
```

#### 后执行脚本

```javascript
暂无后执行脚本
```

#### 成功响应示例

```javascript
东方红灿烂如血，佳节日出把心暖；黎明山林清晰可见，心情怡然安然。
```
