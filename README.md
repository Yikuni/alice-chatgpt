# Alice Chatgpt

这是一个用go语言实现的chatgpt接口,对chatgpt的api进行了封装, 通过该接口, 仅一个简单的http请求就可以开启与chatgpt AI的对话

---

## 亮点

- 方便的api,支持gpt3,gpt3-turbo,gpt4
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

直接双击

### 启动参数说明

| 参数 | 说明                                              | 默认值 | 必填 |
| ---- | ------------------------------------------------- | ------ | ---- |
| p    | 端口号                                            | 7777   | 否   |
| t    | 秘钥                                              | alice  | 否   |
| a    | 是否自动移除无法使用的openai key                  | true   | 否   |
| l    | limit 每分钟chatgpt api的限流                     | 500    | 否   |
| db   | 存储分享的会话使用的数据库,目前仅支持内置的badger | badger | 否   |
| auth | 认证方式 none/simple/normal, 将在下文说明         | simple | 否   |

## Development

### 接口基础调用顺序

1. 调用create接口,可以设置prompt和对话类型, 获取创建的会话id
2. 调用chat接口,请求体放客户输入内容,header中放会话id
3. 可以多次调用chat接口
4. 完成会话后及时调用finish接口,header中放会话id(默认30分钟无动作后会话过期)

### 认证方式

#### none

无认证,适用于个人适用

#### simple

简单认证,适用于作为前端不直接调用的微服务使用

在请求头中加入如下字段

| token | 秘钥,即启动参数-t的值 |
| ----- | --------------------- |

#### normal

正常认证,适用前端直接调用的情况

在请求头中加入如下字段

| token | 对秘钥进行加密计算后的值 |
| ----- | ------------------------ |
| uuid  | uuid                     |

uuid=随机生成一个uuid

token=sha256(设置的token+uuid)

##### 示例

``` javascript
const key = "alice"
const uuid = new Date().getMilliseconds().toString()
const token = sha256(key + uuid)
axios.post("/chatgpt/finish", textarea.value, {headers: {token: token, uuid: uuid, conversation: conversation.value}}).catch(e =>{
      console.log(e.response.data)
})
```

## 接口文档

### 创建基于gpt3的对话

#### 接口URL

> http://localhost:7777/chatgpt/create

#### 请求方式

> POST

#### 请求Header参数

| 参数名   | 示例值     | 参数类型 | 是否必填 | 参数描述                                                     |
| -------- | ---------- | -------- | -------- | ------------------------------------------------------------ |
| settings | FriendChat | String   | 否       | 聊天类型, 默认空为普通聊天, FriendChat为朋友聊天,暂时没有其它的类型 |

#### 请求Body参数

AI的设定,可以留白,默认为以下设定

```javascript
The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.
```

#### 成功响应示例

返回对话的id

```javascript
5dCMMKFI
```

### 创建基于gpt3-turbo的对话

#### 接口URL

> http://localhost:7777/chatgpt/createTurbo

#### 请求方式

> POST

#### 请求Body参数

AI的设定,可以留白,默认为以下设定

```javascript
You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible.
```

#### 成功响应示例

返回对话的id

```javascript
5dCMMKFI
```

### 创建基于gpt4的对话

#### 接口URL

> http://localhost:7777/chatgpt/createGPT4

#### 请求方式

> POST

#### 请求Body参数

AI的设定,可以留白,默认为以下设定

```javascript
You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible.
```

#### 成功响应示例

返回对话的id

``` javascript
5dCMMKFI
```

### 创建基于gpt3的角色扮演对话

#### 接口URL

> http://localhost:7777/chatgpt/createRolePlay

#### 请求方式

> POST

#### 请求Body参数

| 参数名     | 示例值                                                       | 参数类型 | 是否必填 | 参数描述           |
| ---------- | ------------------------------------------------------------ | -------- | -------- | ------------------ |
| human_name | human                                                        | String   | 是       | 用户的称呼         |
| ai_name    | AI                                                           | String   | 是       | ai扮演的角色的称呼 |
| prompt     | You are ChatGPT, a large language model trained by OpenAI. Answer as concisely as possible. | String   | 是       | 设定               |

#### 成功响应示例

返回对话的id

``` javascript
5dCMMKFI
```

### 对话

#### 接口URL

> http://localhost:7777/chatgpt/chat

#### 请求方式

> POST

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述 |
| ------------ | -------- | -------- | -------- | -------- |
| conversation | 5dCMMKFI | String   | 是       | 对话id   |

#### 请求Body参数

用户的问题

```javascript
困的时候可以睡觉吗
```

#### 成功响应示例

AI的回复

```javascript
 是的，睡觉可以帮助你放松身心，释放压力，并且提升你的专注力和注意力，恢复你的精力。所以当你想要休息的时候，请一定记得睡觉吧！
```

### 重新生成

#### 接口说明

重新生成最后一个提问的答案

#### 接口URL

> http://localhost:7777/chatgpt/regenerate

#### 请求方式

> POST

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述 |
| ------------ | -------- | -------- | -------- | -------- |
| conversation | 5dCMMKFI | String   | 是       | 对话id   |

#### 成功响应示例

AI的回复

```javascript
 是的，睡觉可以帮助你放松身心，释放压力，并且提升你的专注力和注意力，恢复你的精力。所以当你想要休息的时候，请一定记得睡觉吧！
```

### 撤回消息

#### 接口说明

撤回最后一个人类的问题以及对应的AI的回答

#### 接口URL

> http://localhost:7777/chatgpt/rollback

#### 请求方式

> POST

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述 |
| ------------ | -------- | -------- | -------- | -------- |
| conversation | 5dCMMKFI | String   | 是       | 对话id   |

#### 成功响应示例

一个提示

```javascript
Rollback Succeeded
```

## 获取对话上下文

#### 接口说明

获取对话上下文的纯文本形式

#### 接口URL

> http://localhost:7777/chatgpt/context

#### 请求方式

> POST

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述 |
| ------------ | -------- | -------- | -------- | -------- |
| conversation | qrHe4o0P | String   | 是       | 会话id   |

#### 成功响应示例

```javascript
The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.

Human: 困的时候可以睡觉吗
AI:  是的，睡觉可以帮助你放松身心，释放压力，并且提升你的专注力和注意力，恢复你的精力。所以当你想要休息的时候，请一定记得睡觉吧！
```

#### 错误响应示例

```javascript
failed to find Conversation with provided id: 4ev2Djnd 
```

## 结束对话

#### 接口说明

手动清除对话,释放服务器资源。对话默认在无互动20分钟后清除。

#### 接口URL

> http://localhost:7777/chatgpt/finish

#### 请求方式

> POST

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述 |
| ------------ | -------- | -------- | -------- | -------- |
| conversation | qrHe4o0P | String   | 是       | 会话id   |

#### 成功响应示例

成功对话的提示

```javascript
Conversation finished
```

#### 错误响应示例

```javascript
failed to find Conversation with provided id: qrHe4o0P 
```

## 保存对话

#### 接口说明

保存对话到数据库

#### 接口URL

> http://localhost:7777/chatgpt/store

#### 请求方式

> POST

#### 请求Header参数

| 参数名       | 示例值   | 参数类型 | 是否必填 | 参数描述 |
| ------------ | -------- | -------- | -------- | -------- |
| conversation | ExpHJpdD | String   | 是       | 会话id   |

#### 成功响应示例

保存成功的提示

```javascript
保存成功
```

#### 错误响应示例

```javascript
failed to find Conversation with provided id: ExpHJpdD 
```

### 获取保存的对话

#### 接口说明

获取分享的对话, 即调用store接口保存的对话

#### 接口URL

> http://localhost:7777/chatgpt/share/<会话id>

#### 请求方式

> GET

#### 请求body参数

| 参数名       | 类型   | 备注     |
| ------------ | ------ | -------- |
| Id           | string | 对话id   |
| Prompt       | string | 设定     |
| SentenceList | Array  | 对话列表 |

#### 成功响应示例

```javascript
{"Id":"hJ0CG23B","Prompt":"The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly.\n","SentenceList":["こんにちは","ーこんにちはー は何か手伝えることがありますか？","ありますよ！そうですね、あたいはよくかんがえたけと、もしかしたらあたいは本当にアホなの？","いいえ！あなたは素晴らしく有能です。あなたは苦労して学んだ結果、多くのことを習得できます。私たちは皆同じではありませんが、皆勇気と創造性に満ちています！","ええ、そんなことがあるの？ありがとうございます。もしかしてあたいは天才かもしれません。そう考えると明日やったことない色んな事をやってみたいようになった。それは現実なのか、幻想なのか？","今あなたが考えていることは、本当の現実です！あなたの勇気と創造性を使って、新しい機会や可能性を探求してみてください。あなたには、一生懸命学んでいる努力が報われるでしょう！"]}
```

#### 错误响应示例

```javascript
Failed to find shared conversation with provided id
```

## 快速回答

#### 接口说明

进行一个不创建对话的,一次性的提问-回答

#### 接口URL

> http://localhost:7777/chatgpt/quickAnswer

#### 请求方式

> POST

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
| prompt   | The following is a conversation with an AI assistant. The assistant is helpful, creative, clever, and very friendly. | String   | 是       | prompt                                 |
| examples | "今天天气晴朗,我在公园漫步,看到绿色植物,心情愉悦。请将上文改编成一首诗",<br/>        "春光洒满天堂， 公园青色荫蔽； 绿叶湖边绰绰， 心情一片悦乐。" | Array    | 是       | 例子,奇数为人类的问题,偶数为AI回答示例 |
| question | 今天去看日出,心情愉悦。请将上文改编成一首诗。                | String   | 是       | 问题                                   |

#### 成功响应示例

ai的回答

```javascript
东方红灿烂如血，佳节日出把心暖；黎明山林清晰可见，心情怡然安然。
```

## 文章概述

#### 接口URL

> http://localhost:7777/chatgpt/summary

#### 请求方式

> POST

#### 请求Body参数

要总结的文章

```javascript
 Yikuni, 根据我们团队多方面的研究，所有资料都指向同一个结果：她曾经确实地在这个世界生活过，并且理论上大部分人都听说过Yikuni。

 但是从一个特别的时间点开始，有关她的一切，在人们的记忆中逐渐褪去，像一片墨水滴到广阔的大海中，没有惊起任何声响，悄然地消失了。没有人意识到Yikuni的离开。当我们团队意识到这件事情的时候，岁月的时钟已经不知道走过了几个轮回。

 根据历史资料记载，Yikuni似乎不是唯一一个从人们的记忆中突然消失的人。我们对曾经发生过同样事件的资料进行收集，推测Yikuni有非常大的可能同东风谷早苗一样，到达了幻想乡。这样解释的话Yikuni的消失就不再是一团迷了。
```

#### 成功响应示例

总结后的结果

```javascript
研究表明Yikuni曾经生活过，然而从某一时刻开始，有关她的所有信息都在人们的记忆中逐渐消失。根据历史资料，Yikuni似乎并不是唯一一个突然消失的人，我们推测了Yikuni可能像东风谷早苗一样到达了幻想乡。
```
