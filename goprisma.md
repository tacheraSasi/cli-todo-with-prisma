Here’s a detailed sheet on using the **Go Prisma package** for basic operations. Prisma is a great ORM for Go, providing type-safe database access and query building.

---

### **Setup Guide for Go Prisma**

1. **Install Prisma CLI:**
   ```bash
   npm install -g prisma
   ```

2. **Initialize Prisma in Your Project:**
   ```bash
   prisma init
   ```
   This creates a `prisma/schema.prisma` file and a `.env` file for your database configuration.

3. **Set Up Your Database:**
   Update the `.env` file with your database URL:
   ```dotenv
   DATABASE_URL="postgresql://user:password@localhost:5432/mydb?schema=public"
   ```

4. **Install Go Prisma Client:**
   In your Go project, install the Prisma client:
   ```bash
   go get github.com/prisma/prisma-client-go
   ```

5. **Define Your Prisma Schema:**
   Example of a `schema.prisma` file:
   ```prisma
   datasource db {
     provider = "postgresql"
     url      = env("DATABASE_URL")
   }

   generator client {
     provider = "prisma-client-go"
   }

   model User {
     id        Int      @id @default(autoincrement())
     name      String
     email     String   @unique
     createdAt DateTime @default(now())
   }
   ```

6. **Generate the Go Client:**
   ```bash
   prisma generate
   ```

---

### **Basic CRUD Operations in Go with Prisma**

#### **1. Connect to the Database**
```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/prisma/prisma-client-go/runtime/transaction"
    "your_project/prisma-client"
)

func main() {
    client := prisma.NewClient()
    if err := client.Connect(); err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    defer client.Disconnect()

    ctx := context.Background()

    // Perform operations here
    fmt.Println("Connected to the database successfully!")
}
```

---

#### **2. Create a Record**
```go
func createUser(client *prisma.Client) {
    ctx := context.Background()
    user, err := client.User.CreateOne(
        prisma.User.Name.Set("John Doe"),
        prisma.User.Email.Set("john.doe@example.com"),
    ).Exec(ctx)
    if err != nil {
        log.Fatalf("Error creating user: %v", err)
    }

    fmt.Printf("User created: %+v\n", user)
}
```

---

#### **3. Read Records**
- **Find All Users**
```go
func getAllUsers(client *prisma.Client) {
    ctx := context.Background()
    users, err := client.User.FindMany().Exec(ctx)
    if err != nil {
        log.Fatalf("Error fetching users: %v", err)
    }

    fmt.Printf("Users: %+v\n", users)
}
```

- **Find a User by Email**
```go
func getUserByEmail(client *prisma.Client, email string) {
    ctx := context.Background()
    user, err := client.User.FindUnique(prisma.User.Email.Equals(email)).Exec(ctx)
    if err != nil {
        log.Printf("Error fetching user: %v\n", err)
        return
    }

    fmt.Printf("User: %+v\n", user)
}
```

---

#### **4. Update a Record**
```go
func updateUserEmail(client *prisma.Client, id int, newEmail string) {
    ctx := context.Background()
    user, err := client.User.FindUnique(prisma.User.ID.Equals(id)).
        Update(
            prisma.User.Email.Set(newEmail),
        ).Exec(ctx)
    if err != nil {
        log.Fatalf("Error updating user: %v", err)
    }

    fmt.Printf("User updated: %+v\n", user)
}
```

---

#### **5. Delete a Record**
```go
func deleteUser(client *prisma.Client, id int) {
    ctx := context.Background()
    user, err := client.User.FindUnique(prisma.User.ID.Equals(id)).
        Delete().
        Exec(ctx)
    if err != nil {
        log.Fatalf("Error deleting user: %v", err)
    }

    fmt.Printf("User deleted: %+v\n", user)
}
```

---

#### **6. Transactions**
```go
func transactionalOperation(client *prisma.Client) {
    ctx := context.Background()
    err := client.Prisma.Transaction(ctx, func(ctx transaction.Context) error {
        // Create a user in a transaction
        _, err := client.User.CreateOne(
            prisma.User.Name.Set("Transactional User"),
            prisma.User.Email.Set("transactional@example.com"),
        ).Exec(ctx)
        if err != nil {
            return fmt.Errorf("Error in transaction: %v", err)
        }

        // Simulate an error to roll back
        return fmt.Errorf("Rolling back transaction")
    })

    if err != nil {
        log.Printf("Transaction rolled back: %v\n", err)
    } else {
        fmt.Println("Transaction committed successfully!")
    }
}
```

---

### **Complete Code Example**
```go
package main

import (
    "context"
    "fmt"
    "log"

    "your_project/prisma-client"
)

func main() {
    client := prisma.NewClient()
    if err := client.Connect(); err != nil {
        log.Fatalf("Error connecting to the database: %v", err)
    }
    defer client.Disconnect()

    ctx := context.Background()

    // Create a user
    user, err := client.User.CreateOne(
        prisma.User.Name.Set("Alice"),
        prisma.User.Email.Set("alice@example.com"),
    ).Exec(ctx)
    if err != nil {
        log.Fatalf("Error creating user: %v", err)
    }
    fmt.Printf("Created user: %+v\n", user)

    // Get all users
    users, err := client.User.FindMany().Exec(ctx)
    if err != nil {
        log.Fatalf("Error fetching users: %v", err)
    }
    fmt.Printf("All users: %+v\n", users)
}
```

---

### **Key Notes**
- Ensure the database is properly migrated using:
  ```bash
  prisma migrate dev
  ```
- Use type-safe query operations, which are auto-generated from your schema.

Prisma makes database operations intuitive and type-safe, leveraging Go's strong typing and simplicity!