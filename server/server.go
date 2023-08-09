package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("postgres", "postgres://chaeng:1234@localhost/go_todo?sslmode=disable")

	if err != nil {
		log.Fatal(err)
		return
	}
}

func toggleDone(c *gin.Context) {
	todoId := c.Param("id")
	i, err := strconv.Atoi(todoId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var todo Todo
	rows, err := db.Query("SELECT * from todo where id=$1", i)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	todo.ID = i
	for rows.Next() {
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Done); err != nil {
			fmt.Errorf("%w", err)
		}
	}

	todo.Done = !todo.Done

	_, err = db.Exec("UPDATE todo SET done = $1 where id = $2;", todo.Done, todo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, todo)
}

func getTodo(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM todo;")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		fmt.Println("Error SELECTING database")
	}
	defer rows.Close()

	var todos []Todo
	for rows.Next() {
		var todo Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Done); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("id s %d, and  title is %s\n", todo.ID, todo.Title)
		todos = append(todos, todo)
	}

	for i := 0; i < len(todos); i++ {
		log.Printf("id s %d, and  title is %s\n", todos[i].ID, todos[i].Title)
	}

	c.JSON(http.StatusOK, todos)
}

func createTodo(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	if contentType != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content-Type must be application/json"})
		return
	}

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
	}

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS todo(
        id SERIAL,
        title VARCHAR(100) NOT NULL,
        done BOOL
    );
    `)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Error creating table"})
		return
	}
	rows := db.QueryRow("INSERT INTO todo (title, done) VALUES ($1, $2) RETURNING id;", todo.Title, todo.Done).Scan(&todo.ID)
	if rows != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "This is error"})
	}
	c.JSON(http.StatusCreated, todo)
}

func deleteTodo(c *gin.Context) {
	idToDelete := c.Param("id")
	i, err := strconv.Atoi(idToDelete)
	if err != nil {
		log.Print(err)
		return
	}

	rows, err := db.Exec("DELETE FROM todo WHERE id = $1", i)
	fmt.Println(rows)
	if err != nil {
		log.Print(err)
		return
	}

	if err != nil {
		log.Print(err)
		return
	}
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/", getTodo)
	r.POST("/", createTodo)
	r.PUT("/:id", toggleDone)
	r.DELETE("/:id", deleteTodo)

	r.Run()
}
