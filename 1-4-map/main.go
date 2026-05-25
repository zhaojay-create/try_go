package main

import (
	"fmt"
	"sync"
	"time"
)

// 1. 映射声明与初始化
// var m1 map[string]int        // 声明 nil 映射
// m2 := make(map[string]int)   // 使用 make 创建
// m3 := map[string]int{        // 字面量创建
//     "apple":  5,
//     "banana": 10,
// }

// 2. 操作
// m2["apple"] = 5              // 插入或更新
// value := m2["apple"]         // 查找
// delete(m2, "apple")          // 删除
// value, exists := m2["apple"] // 检查键是否存在

// 并发安全问题: Go的map在并发读写时不是线程安全的，需要配合sync包中的锁机制；或者使用sync.Map
//
// sync.RWMutex + map vs sync.Map：
//   sync.RWMutex + map：读写均衡或写多场景，手动加锁，性能更可控
//   sync.Map：读多写少，或不同 goroutine 操作不同 key，内置并发安全，无需手动加锁
//
// 场景选择：
//   读多写少                    → sync.Map
//   读写均衡 / 写多              → sync.RWMutex + map
//   每个 goroutine 操作独立的key → sync.Map

type User struct {
	UserID         int
	UserName       string
	LoginTime      time.Time
	LastActionTime time.Time
	IPAddress      string
}

// *User 要影响原始数据，所以用指针
type SessionManager struct {
	sessions map[string]*User
	mutex    sync.RWMutex // 解决并发安全
}

// 返回值是 *SessionManager（指针）
func NewSessionManager() *SessionManager {
	// & 取地址，返回指针
	return &SessionManager{
		sessions: make(map[string]*User),
	}
}

func (sm *SessionManager) AddSession(token string, userID int, userName string, ipAddress string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session := &User{
		UserID:         userID,
		UserName:       userName,
		LoginTime:      time.Now(),
		LastActionTime: time.Now(),
		IPAddress:      ipAddress,
	}

	sm.sessions[token] = session
}

func (sm *SessionManager) GetSession(token string) (*User, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	session, exists := sm.sessions[token]
	if exists {
		session.LastActionTime = time.Now()
	}
	return session, exists
}

func (sm *SessionManager) CleanUpSession(timeout time.Duration) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	now := time.Now()
	for token, session := range sm.sessions {
		if now.Sub(session.LastActionTime) > timeout {
			delete(sm.sessions, token)
		}
	}
}

func (sm *SessionManager) RemoveSession(token string) bool {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	_, exists := sm.sessions[token]
	if exists {
		delete(sm.sessions, token)
	}
	return exists
}

func main() {
	manager := NewSessionManager()

	manager.AddSession("token1", 1, "user1", "127.0.0.1")
	manager.AddSession("token2", 2, "user2", "127.0.0.2")

	removedToken1 := manager.RemoveSession("token1")
	fmt.Printf("移除 token1: %v\n", removedToken1)

	removed := manager.RemoveSession("token100")
	fmt.Printf("移除 token100: %v\n", removed)

	_, exists := manager.GetSession("token1")
	fmt.Printf("token1 是否还存在: %v\n", exists)

	session, exists := manager.GetSession("token1")
	if exists {
		fmt.Printf("用户ID: %d, 用户名: %s, 登录时间: %s, 最后活跃: %s, IP: %s\n",
			session.UserID, session.UserName,
			session.LoginTime.Format("2006-01-02 15:04:05"),
			session.LastActionTime.Format("2006-01-02 15:04:05"),
			session.IPAddress,
		)
	}

	manager.CleanUpSession(5 * time.Minute)

	fmt.Printf("------------------- \n")
	syncMapExample()
}

func syncMapExample() {
	var m sync.Map

	// 写入 key=uid, value=*User
	m.Store(1, &User{UserID: 1, UserName: "张三", IPAddress: "127.0.0.1", LoginTime: time.Now(), LastActionTime: time.Now()})
	m.Store(2, &User{UserID: 2, UserName: "李四", IPAddress: "127.0.0.2", LoginTime: time.Now(), LastActionTime: time.Now()})
	m.Store(3, &User{UserID: 3, UserName: "王五", IPAddress: "127.0.0.3", LoginTime: time.Now(), LastActionTime: time.Now()})

	// 读取
	if val, ok := m.Load(1); ok {
		user := val.(*User) // 类型断言，从 any 转回 *User
		fmt.Printf("uid=1: %s, IP: %s\n", user.UserName, user.IPAddress)
	}

	// 不存在时才写入
	_, loaded := m.LoadOrStore(1, &User{UserID: 1, UserName: "新张三"})
	fmt.Printf("uid=1 已存在: %v\n", loaded)

	// 删除
	m.Delete(2)

	// 遍历所有用户
	fmt.Println("所有在线用户:")
	m.Range(func(key, value any) bool {
		user := value.(*User)
		fmt.Printf("  uid: %v, 用户名: %s, IP: %s\n", key, user.UserName, user.IPAddress)
		return true
	})
}
