# 软工集市特供版企微群机器人

在.env文件中添加类似内容

```property
SSE_EMAIL=asd@qq.com
SSE_PASSWORD="asd"
SSE_TELEPHONE=asd
SSE_ID=1000 
```

这里的密码是密文，需要在开发者工具里查看

需要新建一个config.yaml，内容如下

```yaml
bot:
  webhook: 企微群机器人地址
```