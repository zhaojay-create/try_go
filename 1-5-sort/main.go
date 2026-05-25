package main

import (
	"fmt"
	"sort"
	"time"
)

// 1. 排序基本类型
// sort.Ints(slice []int)
// sort.Float64s(slice []float64)
// sort.Strings(slice []string)

// 2. 自定义排序（无需实现接口，推荐）
// sort.Slice(slice, func(i, j int) bool {
//     return slice[i] < slice[j] // 升序
// })

// 3. 实现 sort.Interface 接口排序（适合复杂类型）
// type Interface interface {
//     Len() int
//     Less(i, j int) bool
//     Swap(i, j int)
// }
// sort.Sort(slice)

type Product struct {
	Id        int
	Name      string
	Price     float64   // 价格
	Sales     int       //销量
	CreatedAt time.Time // 创建时间
}

type ProductManager struct {
	products []Product
}

func newProductManager() *ProductManager {
	return &ProductManager{
		products: []Product{
			{Id: 1, Name: "Product 1", Price: 100, Sales: 10, CreatedAt: time.Now().Add(-3 * time.Hour)},
			{Id: 2, Name: "Product 2", Price: 100, Sales: 20, CreatedAt: time.Now().Add(-1 * time.Hour)},
			{Id: 3, Name: "Product 3", Price: 200, Sales: 30, CreatedAt: time.Now().Add(-2 * time.Hour)},
			{Id: 4, Name: "Product 4", Price: 200, Sales: 40, CreatedAt: time.Now()},
		},
	}
}

func (pm *ProductManager) SortByPrice(ascending bool) {
	sort.Slice(pm.products, func(i, j int) bool {
		if ascending {
			return pm.products[i].Price < pm.products[j].Price
		}
		return pm.products[i].Price > pm.products[j].Price
	})
}

func (pm *ProductManager) SortByCreatedAt(ascending bool) {
	sort.Slice(pm.products, func(i, j int) bool {
		if ascending {
			return pm.products[i].CreatedAt.Before(pm.products[j].CreatedAt)
		}
		return pm.products[i].CreatedAt.After(pm.products[j].CreatedAt)
	})
}

// 按多个字段排序
func (pm *ProductManager) SortByPriceThenCreatedAt() {
	sort.Slice(pm.products, func(i, j int) bool {
		pi, pj := pm.products[i], pm.products[j]
		if pi.Price != pj.Price {
			return pi.Price < pj.Price // 价格升序
		}
		return pi.CreatedAt.After(pj.CreatedAt) // 价格相同，创建时间越晚越靠前
	})
}

func (pm *ProductManager) PrintProducts(title string) {
	fmt.Printf("--- %s ---\n", title)
	for _, p := range pm.products {
		fmt.Printf("  [%d] %-12s 价格: %6.2f  销量: %d\n", p.Id, p.Name, p.Price, p.Sales)
	}
}

func main() {
	pm := newProductManager()
	pm.PrintProducts("初始顺序")

	pm.SortByPrice(true)
	pm.PrintProducts("按价格升序")

	pm.SortByPrice(false)
	pm.PrintProducts("按价格降序")

	pm.SortByCreatedAt(true)
	pm.PrintProducts("按创建时间升序")

	pm.SortByCreatedAt(false)
	pm.PrintProducts("按创建时间降序")

	pm.SortByPriceThenCreatedAt()
	pm.PrintProducts("按价格升序+价格相同时创建时间越晚越靠前")

	name := []string{"Zig", "Alice", "Bob", "Cubit", "Dock"}
	sort.Strings(name)
	fmt.Println(name)

	index := sort.SearchStrings(name, "Cubit")
	fmt.Println(index)
}
