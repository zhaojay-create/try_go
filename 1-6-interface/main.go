package main

import "fmt"

// 接口定义：
// 接口是Go语言中的一种类型，它定义了一组方法的集合。
// 一个类型如果实现了接口中定义的所有方法，那么这个类型就实现了这个接口。

// 空接口：
// 空接口可以存储任意类型的值

// 应用场景：
// 不同的支付方式（支付宝，微信支付）都有相同的核心操作，比如支付，查询订单，退款等
// 使用接口可以统一这些操作

type Payment interface {
	Pay(amount float64) (transactionID string, err error)
	QueryOrder(orderID string)
	Refund(amount float64)
}

func main() {
	alipay := &Alipay{AppId: "ali_appid123", AppSecret: "asd", MerchantId: "aapsjdiopj"}
	wechat := &WeChat{AppId: "wechat_appid456", AppSecret: "asasdd", MerchantId: "asdiopj"}

	transactions := []struct {
		Payment Payment
		amount  float64
		Name    string
	}{
		{alipay, 100, "支付宝"},
		{wechat, 200, "微信"},
	}

	for _, item := range transactions {
		transactionId, err := processPayment(item.Payment, item.amount)
		if err != nil {
			fmt.Printf("支付失败: %v\n", err)
		} else {
			fmt.Printf("支付成功，交易ID: %s\n", transactionId)
		}
	}

	fmt.Println("\n===== 接口组合示例 =====")
	// AlipayPro 实现了 FullPaymentService（组合接口）
	alipayPro := &AlipayPro{Alipay{AppId: "ali_pro_001", AppSecret: "xxx", MerchantId: "m001"}}
	processWithFullService(alipayPro, 500)
}

func processPayment(payment Payment, amount float64) (string, error) {
	fmt.Println("开始处理支付请求...")
	transactionID, err := payment.Pay(amount)
	if err != nil {
		fmt.Println("支付失败:", err)
		return "", err
	}
	fmt.Println("结束处理支付请求...")
	return transactionID, nil
}

type Alipay struct {
	AppId      string
	AppSecret  string
	MerchantId string
}

func (a *Alipay) Pay(amount float64) (string, error) {
	fmt.Printf("支付宝 %s 支付 %f\n", a.AppId, amount)
	return "alipay_transaction_123", nil
}

func (a *Alipay) QueryOrder(orderID string) {
	fmt.Printf("查询订单 %s\n", orderID)
}

func (a *Alipay) Refund(amount float64) {
	fmt.Printf("退款 %f\n", amount)
}

type WeChat struct {
	AppId      string
	AppSecret  string
	MerchantId string
}

func (w *WeChat) Pay(amount float64) (string, error) {
	fmt.Printf("微信支付 %s 支付 %f\n", w.AppId, amount)
	return "wx_transaction_123", nil
}

func (w *WeChat) QueryOrder(orderID string) {
	fmt.Printf("查询订单 %s\n", orderID)
}

func (w *WeChat) Refund(amount float64) {
	fmt.Printf("退款 %f\n", amount)
}

// ==================== 接口组合（接口嵌入）示例 ====================
// Go 支持接口组合：一个接口可以嵌入其他接口，继承它们的方法
// 这样可以把小接口组合成大接口，实现接口的复用

// 基础接口 1：日志功能
type Logger interface {
	Log(msg string)
}

// 基础接口 2：统计功能
type Statistic interface {
	Record(name string, value float64)
}

// 组合接口：同时包含 Payment、Logger、Statistic 的所有方法
// 实现 FullPaymentService 必须实现全部 5 个方法
type FullPaymentService interface {
	Payment   // 嵌入 Payment 接口（3 个方法）
	Logger    // 嵌入 Logger 接口（1 个方法）
	Statistic // 嵌入 Statistic 接口（1 个方法）
}

// AlipayPro 实现了完整的 FullPaymentService 接口
type AlipayPro struct {
	Alipay
}

func (a *AlipayPro) Log(msg string) {
	fmt.Printf("[日志] %s\n", msg)
}

func (a *AlipayPro) Record(name string, value float64) {
	fmt.Printf("[统计] %s: %f\n", name, value)
}

// 使用组合接口的函数
func processWithFullService(service FullPaymentService, amount float64) {
	service.Log("开始支付流程")
	transactionID, err := service.Pay(amount)
	if err != nil {
		service.Log("支付失败: " + err.Error())
		return
	}
	service.Record("支付金额", amount)
	service.Log("支付完成，交易ID: " + transactionID)
}
