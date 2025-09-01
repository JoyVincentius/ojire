package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

	"ojire/db"
	"ojire/handlers"
	"ojire/middleware"

	"github.com/gin-gonic/gin"
)

// soal 1

type FruitType string

const (
	FruitTypeImport FruitType = "IMPORT"
	FruitTypeLocal  FruitType = "LOCAL"
)

type Fruit struct {
	ID    int       `json:"fruitId"`
	Name  string    `json:"fruitName"`
	Type  FruitType `json:"fruitType"`
	Stock int       `json:"stock"`
}

var rawJSON = `[
{
	"fruitId": 1,
	"fruitName": "Apel",
	"fruitType": "IMPORT",
	"stock": 10
},
{
	"fruitId": 2,
	"fruitName": "Kurma",
	"fruitType": "IMPORT",
	"stock": 20
},
{
	"fruitId": 3,
	"fruitName": "apel",
	"fruitType": "IMPORT",
	"stock": 50
},
{
	"fruitId": 4,
	"fruitName": "Manggis",
	"fruitType": "LOCAL",
	"stock": 100
},
{
	"fruitId": 5,
	"fruitName": "Jeruk Bali",
	"fruitType": "LOCAL",
	"stock": 10
},
{
	"fruitId": 5,
	"fruitName": "KURMA",
	"fruitType": "IMPORT",
	"stock": 20
},
{
	"fruitId": 5,
	"fruitName": "Salak",
	"fruitType": "LOCAL",
	"stock": 150
}
]`

func uniqueStrings(in []string) []string {
	seen := make(map[string]bool)
	out := []string{}
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

type Answer struct {
	AllFruitNames         []string               `json:"all_fruit_names"`
	Containers            map[FruitType][]string `json:"containers_by_type"`
	ContainerCount        int                    `json:"container_count"`
	TotalStockByContainer map[FruitType]int      `json:"total_stock_by_type"`
	Duplicates            []string               `json:"duplicates"`
	Comment               string                 `json:"comment"`
}

func solve(fruits []Fruit) Answer {
	var names []string
	for _, f := range fruits {
		names = append(names, f.Name)
	}
	names = uniqueStrings(names)

	containers := make(map[FruitType][]string)
	for _, f := range fruits {
		lcName := strings.ToLower(f.Name)
		containers[f.Type] = append(containers[f.Type], lcName)
	}

	for t, list := range containers {
		containers[t] = uniqueStrings(list)
	}

	totalStock := make(map[FruitType]int)
	for _, f := range fruits {
		totalStock[f.Type] += f.Stock
	}

	dupMap := make(map[string]int)
	for _, f := range fruits {
		lc := strings.ToLower(f.Name)
		dupMap[lc]++
	}
	var duplicates []string
	for name, cnt := range dupMap {
		if cnt > 1 {
			duplicates = append(duplicates, name) 
		}
	}
	duplicates = uniqueStrings(duplicates)

	comment := generateComment(names, containers, totalStock, duplicates)

	return Answer{
		AllFruitNames:         names,
		Containers:            containers,
		ContainerCount:        len(containers),
		TotalStockByContainer: totalStock,
		Duplicates:            duplicates,
		Comment:               comment,
	}
}

func generateComment(names []string, cont map[FruitType][]string, stock map[FruitType]int, dup []string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Andi memiliki %d jenis buah (%s). ", len(names), strings.Join(names, ", "))
	fmt.Fprintf(&b, "Buahnya terbagi menjadi %d wadah berdasarkan tipe: ", len(cont))
	types := []string{}
	for t := range cont {
		types = append(types, string(t))
	}
	fmt.Fprintf(&b, "%s. ", strings.Join(types, ", "))
	for t, s := range stock {
		fmt.Fprintf(&b, "Total stok untuk tipe %s adalah %d. ", t, s)
	}
	if len(dup) > 0 {
		fmt.Fprintf(&b, "Terdapat duplikasi nama buah: %s. ", strings.Join(dup, ", "))
	} else {
		b.WriteString("Tidak ada duplikasi nama buah. ")
	}
	return b.String()
}

// soal 2

type Comm struct {
	CommentID      int    `json:"commentId"`
	CommentContent string `json:"commentContent"`
	Replies        []Comm `json:"replies,omitempty"`
}

// hitungKomentar menghitung semua komentar termasuk balasan secara rekursif.
func hitungKomentar(comments []Comm) int {
	total := 0
	for _, c := range comments {
		total++ // menghitung komentar utama
		if len(c.Replies) > 0 {
			total += hitungKomentar(c.Replies) // menelusuri balasan
		}
	}
	return total
}

func main() {
	var fruits []Fruit
	if err := json.Unmarshal([]byte(rawJSON), &fruits); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	ans := solve(fruits)

	fmt.Println("=== Soal 1 ===")
	fmt.Printf("1. Buah yang dimiliki Andi: %v\n", ans.AllFruitNames)
	fmt.Printf("2. Jumlah wadah yang dibutuhkan: %d\n", ans.ContainerCount)
	fmt.Println("   Daftar buah per wadah:")
	for t, list := range ans.Containers {
		fmt.Printf("   - %s : %v\n", t, list)
	}
	fmt.Println("3. Total stock per wadah:")
	for t, tot := range ans.TotalStockByContainer {
		fmt.Printf("   - %s : %d\n", t, tot)
	}
	fmt.Printf("4. Komentar: %s\n", ans.Comment)

	fmt.Println("\n=== Soal 2 ===")
	data := `[
		{
			"commentId": 1,
			"commentContent": "Hai",
			"replies": [
				{
					"commentId": 11,
					"commentContent": "Hai juga",
					"replies": [
						{
							"commentId": 111,
							"commentContent": "Haai juga hai jugaa"
						},
						{
							"commentId": 112,
							"commentContent": "Haai juga hai jugaa"
						}
					]
				},
				{
					"commentId": 12,
					"commentContent": "Hai juga",
					"replies": [
						{
							"commentId": 121,
							"commentContent": "Haai juga hai jugaa"
						}
					]
				}
			]
		},
		{
			"commentId": 2,
			"commentContent": "Halooo"
		}
	]`

	var comments []Comm
	if err := json.Unmarshal([]byte(data), &comments); err != nil {
		fmt.Println("Gagal meng-parse JSON:", err)
	}

	total := hitungKomentar(comments)
	fmt.Printf("Total komentar (termasuk semua balasan): %d\n", total)

	fmt.Println("\n=== Soal tambahan ===")
	a := []int{3, 4, 2, 1, 3, 3}
	b := []int{4, 3, 5, 3, 9, 3}

	sort.Ints(a)
	sort.Ints(b)

	totalSelisih := 0
	for i := 0; i < len(a); i++ {
		selisih := a[i] - b[i]
		if selisih < 0 {
			selisih = -selisih
		}
		totalSelisih += selisih
		fmt.Printf("a: %d (kiri) - b: %d (kanan) - selisih = %d\n", a[i], b[i], selisih)
	}

	fmt.Printf("\nTotal selisih: %d\n", totalSelisih)

	db.Init()

	r := gin.Default()

	// Public routes
	r.POST("/login", handlers.LoginHandler)
	r.GET("/products", handlers.ListProductsHandler)

	// Protected group
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.GET("/cart", handlers.GetCartHandler)
	auth.POST("/cart/add", handlers.AddToCartHandler)
	auth.POST("/cart/remove", handlers.RemoveFromCartHandler)
	auth.POST("/cart/clear", handlers.ClearCartHandler)

	auth.POST("/checkout", handlers.CheckoutHandler)

	// start server
	r.Run(":8080") // http://localhost:8080
}
