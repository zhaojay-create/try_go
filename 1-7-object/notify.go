package main

import "fmt"

type Notifier interface {
	Notify(message string) error
}

type EmailNotifier struct {
	SMTPHost string
	Port     string
}

type SmsNotifier struct {
	APIKey   string
	TmplCode string
}

func (e *EmailNotifier) Notify(message string) error {
	fmt.Printf("邮箱通知 [SMTP: %s:%s]: %s\n", e.SMTPHost, e.Port, message)
	return nil
}

func (s *SmsNotifier) Notify(message string) error {
	fmt.Printf("短信通知通知: APIKey: %s, TmplCode: %s", s.APIKey, s.TmplCode)
	return nil
}

type OrderService struct {
	notifier Notifier
}

func (o *OrderService) setNotifier(n Notifier) {
	o.notifier = n
}

func (o *OrderService) createOrder(productID int, quinity int) {
	fmt.Printf("创建订单")
	err := o.notifier.Notify(fmt.Sprintf("订单创建成功 商品ID:%d 数量: %d", productID, quinity))
	if err != nil {
		fmt.Println("邮箱发送失败")
	}
}

type BroadcaseNotify struct {
	notifiers []Notifier
}

func (b *BroadcaseNotify) Notify(message string) error {
	for _, notifier := range b.notifiers {
		notifier.Notify(message)
	}
	return nil
}

func notify() {
	orderService := &OrderService{}
	emailnot := &EmailNotifier{
		SMTPHost: "smtp.example.com",
		Port:     "587",
	}

	smsNot := &SmsNotifier{
		APIKey:   "api_key",
		TmplCode: "tmpl_code",
	}

	broadcase := &BroadcaseNotify{
		notifiers: []Notifier{emailnot, smsNot},
	}

	orderService.setNotifier(broadcase)
	orderService.createOrder(100, 1)

	orderService.setNotifier(emailnot)
	orderService.createOrder(200, 2)

	orderService.setNotifier(smsNot)
	orderService.createOrder(30, 3)
}
