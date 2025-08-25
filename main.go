package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		containers[f.Type] = append(containers[f.Type], f.Name)
	}

	totalStock := make(map[FruitType]int)
	for _, f := range fruits {
		totalStock[f.Type] += f.Stock
	}

	comment := generateComment(names, containers, totalStock)

	return Answer{
		AllFruitNames:         names,
		Containers:            containers,
		ContainerCount:        len(containers),
		TotalStockByContainer: totalStock,
		Comment:               comment,
	}
}

func generateComment(names []string, cont map[FruitType][]string, stock map[FruitType]int) string {
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
	b.WriteString("Tidak ada duplikasi tipe wadah yang diperlukan; hanya dua wadah (IMPORT & LOCAL) cukup.")
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
