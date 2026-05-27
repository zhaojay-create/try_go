package main

import (
	"fmt"
	"time"
)

// 普通用户
type BaseUser struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

func (u *BaseUser) getUserCreateAt() time.Time {
	return u.CreatedAt
}

func (u *BaseUser) displayUserInfo() {
	fmt.Printf("ID: %d, Name: %s, Email: %s, CreatedAt: %s\n", u.ID, u.Name, u.Email, u.CreatedAt)
}

// 地址
type Address struct {
	Province string
	City     string
	Street   string
}

func (a *Address) displayAddress() {
	fmt.Printf("Province: %s, City: %s, Street: %s\n", a.Province, a.City, a.Street)
}

type NormalUser struct {
	BaseUser
	Address []Address
}

func (v *NormalUser) addAddress(address Address) {
	v.Address = append(v.Address, address)
}

// VIP 用户信息
type VipUser struct {
	BaseUser
	Address    []Address
	VipLevel   int
	Discount   float64
	ExpireTime time.Time
}

func (v *VipUser) IsVipValid() bool {
	return time.Now().Before(v.ExpireTime)
}

func (v *VipUser) getDiscount() float64 {
	if v.IsVipValid() {
		return v.Discount
	}
	return 1.0
}

type UserService struct {
}

const (
	VipUserType    = "vip"
	NormalUserType = "normal"
)

func (v *UserService) createUser(userName string, email string, userType string) (interface{}, error) {

	base := BaseUser{
		ID:        time.Now().Unix(),
		Name:      userName,
		Email:     email,
		CreatedAt: time.Now(),
	}

	switch userType {
	case NormalUserType:
		normalUser := NormalUser{
			BaseUser: base,
			Address:  []Address{},
		}
		return normalUser, nil
	case VipUserType:
		vip := VipUser{
			BaseUser:   base,
			Address:    []Address{},
			VipLevel:   1,
			Discount:   0.8,
			ExpireTime: time.Now().Add(time.Hour * 24 * 365),
		}
		return vip, nil
	default:
		return 0, fmt.Errorf("无效用户类型")
	}
}

func main() {
	service := &UserService{}
	normalUser, err := service.createUser("John", "john@example.com", NormalUserType)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(normalUser)

	vipUser, err := service.createUser("Alice", "alicen@example.com", VipUserType)
	if err != nil {
		fmt.Println(err)
		return
	}

	nUser := normalUser.(NormalUser)
	nUser.addAddress(Address{
		Province: "A",
		City:     "C",
		Street:   "S",
	})
	nUser.displayUserInfo()

	users := []interface{}{normalUser, vipUser}

	for _, user := range users {
		switch u := user.(type) {
		case NormalUser:
			u.displayUserInfo()
		case VipUser:
			u.displayUserInfo()
			fmt.Printf("用户 vip: %t\n", u.IsVipValid())
			fmt.Printf("用户的折扣是: %f\n", u.getDiscount())
		}
	}

	fmt.Printf("\n ----------------- \n \n")

	notify()
}
