package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const JWT_SECRET = "secret"

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleAccounts))
	router.HandleFunc("/account/{number}", withJWTAuth(makeHTTPHandleFunc(s.handleAccount), s.store))
	router.HandleFunc("/transfer/{number}", withJWTAuth(makeHTTPHandleFunc(s.handleTransfer), s.store))

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	req := new(LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByNumber(req.Number)
	if err != nil {
		return err
	}
	pw := req.Password
	encpw := acc.EncryptedPassword
	if err := bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)); err != nil {
		return err
	}

	tokenString, err := createJWT(acc)
	if err != nil {
		return err
	}
	w.Header().Add("x-jwt-token", tokenString)

	rt, err := s.store.GetRefreshToken(acc.ID)
	if err != nil {
		return err
	}
	if err := setRefreshToken(w, rt.Token); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, req)
}

func (s *APIServer) handleAccounts(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccounts(w)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter) error {
	accs, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accs)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	if r.Method == "PUT" {
		return s.handleUpdateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	num, err := getNumber(r)
	if err != nil {
		return err
	}
	acc, err := s.store.GetAccountByNumber(num)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	num, err := getNumber(r)
	if err != nil {
		return err
	}
	err = s.store.DeleteAccountByNumber(num)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, nil)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	//req := &CreateAccountRequest{}: Actual object, then you need to dereference it in Decode(req)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	acc, err := NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	rt, err := NewRefreshToken(acc)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(acc); err != nil {
		return err
	}
	if err := s.store.CreateRefreshToken(rt); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleUpdateAccount(w http.ResponseWriter, r *http.Request) error {
	num, err := getNumber(r)
	if err != nil {
		return err
	}
	req := new(UpdateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	acc, err := s.store.GetAccountByNumber(num)
	if err != nil {
		return err
	}
	updateAccountWithReflect(acc, req)
	if err := s.store.UpdateAccount(acc); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, req)
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	num, err := getNumber(r)
	if err != nil {
		return err
	}

	req := new(Transaction)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	amnt := req.Amount
	recipientNumber := req.Recipient

	sender, err := s.store.GetAccountByNumber(num)
	if err != nil {
		return err
	}

	recipient, err := s.store.GetAccountByNumber(recipientNumber)
	if err != nil {
		return err
	}

	if sender.Balance < amnt {
		return fmt.Errorf("insufficient funds")
	}

	sender.Balance -= amnt
	recipient.Balance += amnt

	if err := s.store.UpdateAccount(sender); err != nil {
		return err
	}
	if err := s.store.UpdateAccount(recipient); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, req)
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")
		token, err := parseJWT(tokenString)
		if err != nil && !token.Valid {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized"})
			return
		}
		accountNumber, err := getNumber(r)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		if accountNumber != int(claims["account_number"].(float64)) {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized"})
			return
		}
		printTimeLeft(token)
		handlerFunc(w, r)
	}
}

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"exp":            time.Now().Add(time.Hour).Unix(),
		"account_number": account.Number,
	}
	secret := JWT_SECRET
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func parseJWT(tokenString string) (*jwt.Token, error) {
	secret := JWT_SECRET
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}

func getNumber(r *http.Request) (int, error) {
	numberStr := mux.Vars(r)["number"]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return number, fmt.Errorf("invalid number given %s", numberStr)
	}
	return number, nil
}

func printTimeLeft(token *jwt.Token) {
	claims := token.Claims.(jwt.MapClaims)
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
	timeLeft := expirationTime.Sub(time.Now())
	fmt.Println("Time left: ", timeLeft)
}

func updateAccountWithReflect(acc *Account, req *UpdateAccountRequest) {
	//ChatGPT Aided
	//https://blog.devtrovert.com/p/reflection-in-go-everything-you-need
	//Loop through all fields
	//Update only if update request field is not nil
	accVal := reflect.ValueOf(acc).Elem()
	reqVal := reflect.ValueOf(req).Elem()
	for i := 0; i < reqVal.NumField(); i++ {
		reqField := reqVal.Field(i)
		if !reqField.IsNil() {
			fieldName := reqVal.Type().Field(i).Name
			accField := accVal.FieldByName(fieldName)
			if accField.IsValid() && accField.CanSet() {
				accField.Set(reqField.Elem())
			}
		}
	}
}

func setRefreshToken(w http.ResponseWriter, tokenString string) error {
	cookie := http.Cookie{
		Name:  "refresh_token",
		Value: tokenString,
		HttpOnly: true,
		Secure: false, // Should be changed
		SameSite: http.SameSiteLaxMode,
		Path: "/", // Should be changed
		Expires: time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, &cookie)
	return nil
}

