package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestApi(t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS employees (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, position TEXT, salary REAL)")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	router := gin.Default()
	router.GET("/", indexHandler)
	router.GET("/employee/:id", getEmployeeByIdHandler)
	router.POST("/employee/create", createEmployeeHandler)
	router.PUT("/employee/update/:id", updateEmployeeHandler)
	router.DELETE("/employee/delete/:id", deleteEmployeeHandler)
	router.GET("/employees", getAllEmployeesHandler)

	// "/"
	// Test indexHandler
	index_r(router, t)

	// "/employee/create"
	// Test createEmployeeHandler
	create_employee(router, t, "1")
	create_employee(router, t, "2")
	create_employee(router, t, "3")
	create_employee(router, t, "4")
	create_employee(router, t, "5")
	create_employee(router, t, "6")
	create_employee(router, t, "7")
	create_employee(router, t, "8")
	create_employee(router, t, "9")
	create_employee(router, t, "10")

	// "/employee/:id"
	// Test getEmployeeByIdHandler
	get_by_id(router, t, 1)
	get_by_id(router, t, 4)
	get_by_id(router, t, 7)
	get_by_id(router, t, 10)

	// "/employees"
	// Test getAllEmployeesHandler
	get_all_employees(router, t)

	// "/employee/update/:id"
	// Test updateEmployeeHandler
	update_employee(router, t)

	// "/employee/delete/:id"
	// Test deleteEmployeeHandler
	delete_employee(router, t, 1)

}

func delete_employee(router *gin.Engine, t *testing.T, id int) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/employee/delete/%d", id), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
}

func update_employee(router *gin.Engine, t *testing.T) {
	req_body := []byte(`{
		"name": "Tester 1",
		"position": "QA",
		"salary": 200.56
	}`)
	req, _ := http.NewRequest("PUT", "/employee/update/1", bytes.NewBuffer(req_body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
}

func get_all_employees(router *gin.Engine, t *testing.T) {
	req, _ := http.NewRequest("GET", "/employees", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
	exp_res := `[{"id":1,"name":"Tester 1","position":"QA","salary":22.43},{"id":2,"name":"Tester 2","position":"QA","salary":22.43},{"id":3,"name":"Tester 3","position":"QA","salary":22.43},{"id":4,"name":"Tester 4","position":"QA","salary":22.43},{"id":5,"name":"Tester 5","position":"QA","salary":22.43},{"id":6,"name":"Tester 6","position":"QA","salary":22.43},{"id":7,"name":"Tester 7","position":"QA","salary":22.43},{"id":8,"name":"Tester 8","position":"QA","salary":22.43},{"id":9,"name":"Tester 9","position":"QA","salary":22.43},{"id":10,"name":"Tester 10","position":"QA","salary":22.43}]`
	if w.Body.String() != exp_res {
		t.Fatalf("Expected response body %s, but got %s", exp_res, w.Body.String())
	}
}

func get_by_id(router *gin.Engine, t *testing.T, n int) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/employee/%d", n), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
	exp_res := fmt.Sprintf(`{"id":%d,"name":"Tester %d","position":"QA","salary":22.43}`, n, n)
	if w.Body.String() != exp_res {
		t.Fatalf("Expected response body %s, but got %s", exp_res, w.Body.String())
	}
}

func index_r(router *gin.Engine, t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
	exp_res := "{\"message\":\"Server Online!!\"}"
	if w.Body.String() != exp_res {
		t.Fatalf("Expected response body %s, but got %s", exp_res, w.Body.String())
	}
}

func create_employee(router *gin.Engine, t *testing.T, name string) {
	req_body := []byte(fmt.Sprintf(`{
		"name": "Tester %s",
		"position": "QA",
		"salary": 22.43
	}`, name))
	req, _ := http.NewRequest("POST", "/employee/create", bytes.NewBuffer(req_body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected status code 200, but got %d", w.Code)
	}
}
