# ojire – Simple Cart API

A lightweight Go REST API that demonstrates basic e‑commerce cart functionality:

* ✅ Add items to a cart
* ✅ List cart contents
* ✅ Remove/Update items
* ✅ JWT‑protected routes (user identification)

The project uses:

| Tool | Version |
|------|---------|
| Go   | 1.22+   |
| Gin  | v1.9    |
| sqlx | v1.7    |
| MySQL|
| JWT  | dgrijalva/jwt-go |

---

## 🚀 Getting Started

### Prerequisites

* **Go** (≥ 1.22)
* **Git**
* A SQL database (MySQL)  

### Clone the repo & install dependencies

```bash
git clone https://github.com/<your‑username>/ojire.git
cd ojire
go mod tidy
```

### Create DB and Tables

Create database and execute query in file schema.sql

### Start API Services

```bash
go run .
```

### Project Disclaimers

* Soal 1 & Soal 2 will be executed when running project
* For simple app, Stock only reduced when checkout
* All orders status is "pending"
* I've created postman collection, be sure to import it
* If you found any issue, please don't be shy to contact me.