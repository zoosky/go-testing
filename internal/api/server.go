// Package api provides the HTTP API server for the application
// @title           Go Testing API
// @version         1.0
// @description     A sample API server demonstrating Go testing best practices
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.email   support@example.com
// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT
// @host            localhost:8080
// @BasePath        /
package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-testing/internal/calculator"
	"go-testing/internal/database"
	pkgcalculator "go-testing/pkg/calculator"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server represents our API server
type Server struct {
	userRepo   database.UserRepository
	calculator *calculator.Calculator
	pubCalc    *pkgcalculator.Calculator
}

// NewServer creates a new Server with the given dependencies
func NewServer(userRepo database.UserRepository, calc *calculator.Calculator) *Server {
	return &Server{
		userRepo:   userRepo,
		calculator: calc,
		pubCalc:    pkgcalculator.NewCalculator(),
	}
}

// Router returns the HTTP router for the server
func (s *Server) Router() http.Handler {
	mux := http.NewServeMux()
	
	// User endpoints
	mux.HandleFunc("GET /users", s.listUsers)
	mux.HandleFunc("GET /users/", s.getUser)
	mux.HandleFunc("POST /users", s.createUser)
	mux.HandleFunc("PUT /users/", s.updateUser)
	mux.HandleFunc("DELETE /users/", s.deleteUser)
	
	// Calculator endpoints
	mux.HandleFunc("GET /calculator/add", s.add)
	mux.HandleFunc("GET /calculator/subtract", s.subtract)
	mux.HandleFunc("GET /calculator/multiply", s.multiply)
	mux.HandleFunc("GET /calculator/divide", s.divide)
	
	// Swagger endpoint
	mux.HandleFunc("GET /swagger/*", func(w http.ResponseWriter, r *http.Request) {
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("list"),
			httpSwagger.DomID("swagger-ui"),
		).ServeHTTP(w, r)
	})
	
	return mux
}

// Helper function to respond with JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper function to respond with an error
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// User handlers

// listUsers godoc
// @Summary List all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} database.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (s *Server) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.userRepo.ListUsers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error retrieving users")
		return
	}
	
	respondJSON(w, http.StatusOK, users)
}

// getUser godoc
// @Summary Get a user by ID
// @Description Get a single user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} database.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	
	user, err := s.userRepo.GetUser(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	
	respondJSON(w, http.StatusOK, user)
}

// createUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body database.User true "User information"
// @Success 201 {object} database.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	if err := s.userRepo.CreateUser(&user); err != nil {
		respondError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	
	respondJSON(w, http.StatusCreated, user)
}

// updateUser godoc
// @Summary Update a user
// @Description Update an existing user's information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body database.User true "Updated user information"
// @Success 200 {object} database.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [put]
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	
	var user database.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	
	// Ensure ID in path matches ID in body
	user.ID = id
	
	if err := s.userRepo.UpdateUser(&user); err != nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	
	respondJSON(w, http.StatusOK, user)
}

// deleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [delete]
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path
	id, err := extractIDFromPath(r.URL.Path)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	
	if err := s.userRepo.DeleteUser(id); err != nil {
		respondError(w, http.StatusNotFound, "User not found")
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

// Calculator handlers

// add godoc
// @Summary Add two numbers
// @Description Add two numbers and return the result
// @Tags calculator
// @Accept json
// @Produce json
// @Param a query number true "First number"
// @Param b query number true "Second number"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} map[string]string
// @Router /calculator/add [get]
func (s *Server) add(w http.ResponseWriter, r *http.Request) {
	a, b, err := getOperands(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	result := s.pubCalc.Add(a, b)
	respondJSON(w, http.StatusOK, map[string]float64{"result": result})
}

// subtract godoc
// @Summary Subtract two numbers
// @Description Subtract the second number from the first and return the result
// @Tags calculator
// @Accept json
// @Produce json
// @Param a query number true "First number"
// @Param b query number true "Second number"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} map[string]string
// @Router /calculator/subtract [get]
func (s *Server) subtract(w http.ResponseWriter, r *http.Request) {
	a, b, err := getOperands(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	result := s.pubCalc.Subtract(a, b)
	respondJSON(w, http.StatusOK, map[string]float64{"result": result})
}

// multiply godoc
// @Summary Multiply two numbers
// @Description Multiply two numbers and return the result
// @Tags calculator
// @Accept json
// @Produce json
// @Param a query number true "First number"
// @Param b query number true "Second number"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} map[string]string
// @Router /calculator/multiply [get]
func (s *Server) multiply(w http.ResponseWriter, r *http.Request) {
	a, b, err := getOperands(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	result := s.pubCalc.Multiply(a, b)
	respondJSON(w, http.StatusOK, map[string]float64{"result": result})
}

// divide godoc
// @Summary Divide two numbers
// @Description Divide the first number by the second and return the result
// @Tags calculator
// @Accept json
// @Produce json
// @Param a query number true "First number (dividend)"
// @Param b query number true "Second number (divisor)"
// @Success 200 {object} map[string]float64
// @Failure 400 {object} map[string]string
// @Router /calculator/divide [get]
func (s *Server) divide(w http.ResponseWriter, r *http.Request) {
	a, b, err := getOperands(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	result, err := s.pubCalc.Divide(a, b)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Division by zero")
		return
	}
	
	respondJSON(w, http.StatusOK, map[string]float64{"result": result})
}

// Helper functions

func extractIDFromPath(path string) (int, error) {
	// Extract ID from path like "/users/123"
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0, strconv.ErrSyntax
	}
	
	return strconv.Atoi(parts[2])
}

func getOperands(r *http.Request) (float64, float64, error) {
	query := r.URL.Query()
	
	aStr := query.Get("a")
	bStr := query.Get("b")
	
	if aStr == "" || bStr == "" {
		return 0, 0, strconv.ErrSyntax
	}
	
	a, err := strconv.ParseFloat(aStr, 64)
	if err != nil {
		return 0, 0, err
	}
	
	b, err := strconv.ParseFloat(bStr, 64)
	if err != nil {
		return 0, 0, err
	}
	
	return a, b, nil
}