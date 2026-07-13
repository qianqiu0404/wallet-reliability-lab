package lab

func BuiltinScenarios() []Scenario {
	return []Scenario{
		{
			ID: "normal-withdrawal", Index: "01", ShortTitle: "正常提现基线", Title: "一笔提现如何安全走完状态机",
			Summary: "从审核、风控、构建、签名和广播推进到链上确认与通知。",
			Steps: []Step{
				{ID: "normal-accept", Service: "wallet-service", Phase: "REQUEST", State: "review_pending", Tone: "active", Title: "接受提现并先落库", VirtualTime: "09:42:11.000"},
				{ID: "normal-risk", Service: "risk-service", Phase: "RISK", State: "create_unsign", Tone: "success", Title: "风控放行并生成审批凭证", VirtualTime: "09:42:11.860"},
				{ID: "normal-build", Service: "wallet-api", Phase: "BUILD", State: "create_unsign", Tone: "active", Title: "构建稳定的待签名交易", VirtualTime: "09:42:13.120"},
				{ID: "normal-sign", Service: "wallet-sign", Phase: "SIGN", State: "create_unsign", Tone: "active", Title: "签名机校验双重审批后签名", VirtualTime: "09:42:14.430"},
				{ID: "normal-store", Service: "wallet-service", Phase: "PERSIST", State: "signed", Tone: "success", Title: "保存 signed raw transaction", VirtualTime: "09:42:15.210"},
				{ID: "normal-broadcast", Service: "wallet-api", Phase: "BROADCAST", State: "broadcasted", Tone: "success", Title: "广播成功，进入链上等待", VirtualTime: "09:42:16.090"},
				{ID: "normal-confirm", Service: "chain", Phase: "CONFIRM", State: "wallet_done", Tone: "recovered", Title: "扫链确认交易完成", VirtualTime: "09:42:27.320"},
				{ID: "normal-notify", Service: "notifier", Phase: "NOTIFY", State: "notified", Tone: "success", Title: "通知业务方并闭环", VirtualTime: "09:42:28.010"},
			},
		},
		{
			ID: "retryable-broadcast", Index: "02", ShortTitle: "可重试广播失败", Title: "广播失败为什么仍然保持 signed",
			Summary: "节点短暂不可用时保留原始交易，不重新签名，等待安全重试。",
			Steps: []Step{
				{ID: "retry-ready", Service: "wallet-service", Phase: "PREPARE", State: "signed", Tone: "success", Title: "交易已经签名并安全落库", VirtualTime: "10:08:31.000"},
				{ID: "retry-attempt", Service: "wallet-api", Phase: "BROADCAST", State: "signed", Tone: "active", Title: "发送已签名原始交易", VirtualTime: "10:08:31.380"},
				{ID: "retry-failure", Service: "wallet-api", Phase: "FAULT", State: "signed", Tone: "warning", Title: "节点返回可重试错误", VirtualTime: "10:08:36.380"},
				{ID: "retry-preserve", Service: "wallet-service", Phase: "CLASSIFY", State: "signed", Tone: "warning", Title: "保持 signed，不产生余额副作用", VirtualTime: "10:08:36.520"},
				{ID: "retry-backoff", Service: "wallet-service", Phase: "RETRY", State: "signed", Tone: "active", Title: "退避后重新加载原交易", VirtualTime: "10:08:46.520"},
				{ID: "retry-success", Service: "wallet-api", Phase: "BROADCAST", State: "broadcasted", Tone: "recovered", Title: "第二次广播成功", VirtualTime: "10:08:47.060"},
				{ID: "retry-confirm", Service: "chain", Phase: "CONFIRM", State: "wallet_done", Tone: "success", Title: "链上确认完成", VirtualTime: "10:08:57.830"},
			},
		},
		{
			ID: "compensation-recovery", Index: "03", ShortTitle: "广播后补偿恢复", Title: "广播后长期未完成如何自动补偿",
			Summary: "超过阈值的 broadcasted 记录通过链上查询恢复本地事实。",
			Steps: []Step{
				{ID: "comp-broadcasted", Service: "wallet-service", Phase: "OBSERVE", State: "broadcasted", Tone: "success", Title: "提现已经记录为 broadcasted", VirtualTime: "11:20:00.000"},
				{ID: "comp-wait", Service: "chain", Phase: "WAIT", State: "broadcasted", Tone: "warning", Title: "正常确认路径仍未推进", VirtualTime: "11:29:59.999"},
				{ID: "comp-stale", Service: "wallet-service", Phase: "SELECT", State: "broadcasted", Tone: "active", Title: "超过阈值，进入补偿候选", VirtualTime: "11:30:00.001"},
				{ID: "comp-query", Service: "wallet-api", Phase: "RECONCILE", State: "broadcasted", Tone: "active", Title: "按交易哈希查询链上状态", VirtualTime: "11:30:00.280"},
				{ID: "comp-found", Service: "chain", Phase: "FACT", State: "broadcasted", Tone: "recovered", Title: "链上交易已经确认", VirtualTime: "11:30:00.740"},
				{ID: "comp-recovered", Service: "wallet-service", Phase: "RECOVER", State: "wallet_done", Tone: "recovered", Title: "补偿推进到 wallet_done", VirtualTime: "11:30:01.030"},
				{ID: "comp-notify", Service: "notifier", Phase: "NOTIFY", State: "notified", Tone: "success", Title: "恢复后的结果继续正常通知", VirtualTime: "11:30:01.420"},
			},
		},
	}
}

func FindScenario(id string) (Scenario, bool) {
	for _, scenario := range BuiltinScenarios() {
		if scenario.ID == id {
			return scenario, true
		}
	}
	return Scenario{}, false
}
