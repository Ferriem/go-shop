# Golang电商秒杀

## 秒杀系统

- [ ] 前台用户登陆，商品展示，抢购
- [ ] 后台订单管理

```mermaid
flowchart LR
	CDN --> 流量负载 --> 流量拦截 --> 服务器集群 --> RabbitMQ --> MySQL
```

