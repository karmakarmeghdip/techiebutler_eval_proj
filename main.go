package main

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var db, err = sql.Open("sqlite3", "employees.db")

// var db, err = sql.Open("sqlite3", "file::memory:?cache=shared")

func main() {
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS employees (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, position TEXT, salary REAL)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	r := gin.Default()

	r.GET("/", indexHandler)
	r.GET("/employee/:id", getEmployeeByIdHandler)
	r.POST("/employee/create", createEmployeeHandler)
	r.PUT("/employee/update/:id", updateEmployeeHandler)
	r.DELETE("/employee/delete/:id", deleteEmployeeHandler)
	r.GET("/employees", getAllEmployeesHandler)

	r.Run(":8080")
}

// Model
type Employee struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Position string  `json:"position"`
	Salary   float64 `json:"salary"`
}

// Handler

func indexHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Server Online!!",
	})
}

func getEmployeeByIdHandler(c *gin.Context) {
	id := c.Param("id")
	var employee Employee
	err := db.QueryRow("SELECT id, name, position, salary FROM employees WHERE id = ?", id).Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to get employee",
		})
		return
	}
	c.JSON(200, employee)
}

func createEmployeeHandler(c *gin.Context) {
	var employee Employee
	err := c.BindJSON(&employee)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to bind JSON",
		})
		return
	}
	result, err := db.Exec("INSERT INTO employees (name, position, salary) VALUES (?, ?, ?)", employee.Name, employee.Position, employee.Salary)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to create employee",
		})
		return
	}
	id, _ := result.LastInsertId()
	employee.ID = int(id)
	c.JSON(200, employee)
}

func updateEmployeeHandler(c *gin.Context) {
	id := c.Param("id")
	var employee Employee
	err := c.BindJSON(&employee)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to bind JSON",
		})
		return
	}
	_, err = db.Exec("UPDATE employees SET name = ?, position = ?, salary = ? WHERE id = ?", employee.Name, employee.Position, employee.Salary, id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to update employee",
		})
		return
	}
	c.JSON(200, "Updated employee successfully")
}

func deleteEmployeeHandler(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to delete employee",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Employee deleted",
	})
}

func getAllEmployeesHandler(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("size", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "page must be a number",
		})
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "size must be a number",
		})
		return
	}
	offset := (page - 1) * pageSize
	rows, err := db.Query("SELECT id, name, position, salary FROM employees LIMIT ? OFFSET ?", pageSize, offset)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Failed to get employees",
		})
		return
	}
	var employees []Employee
	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary)
		if err != nil {
			c.JSON(500, gin.H{
				"message": "Failed to get employees",
			})
			return
		}
		employees = append(employees, employee)
	}
	if len(employees) == 0 {
		c.JSON(404, gin.H{
			"message": "No employees found",
		})
		return
	}
	c.JSON(200, employees)
}
