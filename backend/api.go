package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"net/http"

	jwt "github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/na50r/gobank/backend/sse"
	"golang.org/x/crypto/bcrypt"
)

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
	broker     *sse.Broker
}

func NewAPIServer(listenAddr string, store Storage, broker *sse.Broker) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
		broker:     broker,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.Use(corsMiddleware)

	// SSE
	router.HandleFunc("/events", s.broker.SSEHandler)
	router.HandleFunc("/publish", s.broker.PublishEndpoint)

	// Endpoints
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleAccounts))
	router.HandleFunc("/account/{number}", withJWTAuth(makeHTTPHandleFunc(s.handleAccount)))
	router.HandleFunc("/transfer/{number}", withJWTAuth(makeHTTPHandleFunc(s.handleTransfer)))
	router.HandleFunc("/image/{number}", makeHTTPHandleFunc(s.handleImage))
	router.HandleFunc("/element", makeHTTPHandleFunc(s.handleGetElement))

	// Refresh
	router.HandleFunc("/refresh", makeHTTPHandleFunc(s.handleRefresh))

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// ChatGPT Aided
// Reference 1: https://stackhawkwpc.wpcomstaging.com/golang-cors-guide-what-it-is-and-how-to-enable-it/ (Only sets first header)
// Reference 2: https://stackoverflow.com/questions/61238680/access-to-fetch-at-from-origin-http-localhost3000-has-been-blocked-by-cors (Sets additional headers)
// Reference 3: https://medium.com/@gaurang.m/allowing-cross-site-requests-in-your-gin-app-golang-1332543d91ed (Implement something similar with Gin)
func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", CLIENT)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (s *APIServer) handleGetElement(w http.ResponseWriter, r *http.Request) error {
	req := new(ElementRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	resp := new(ElementResponse)
	result, err := s.store.GetElement(req.A, req.B)
	if err != nil {
		return err
	}
	resp.Result = *result
	return WriteJSON(w, http.StatusOK, resp)
}


func (s *APIServer) handleRefresh(w http.ResponseWriter, r *http.Request) error {
	req := new(RefreshRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	resp := new(RefreshResponse)

	tokenString := req.RefreshToken
	token, err := parseJWT(tokenString)
	if err != nil {
		return err
	}
	if !token.Valid {
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "unauthorized"})
	}

	acc, err := s.store.GetAccountByRefreshToken(tokenString)
	if err != nil {
		return err
	}

	at, err := createJWT(acc)
	if err != nil {
		return err
	}

	rt, err := NewRefreshToken(acc)
	if err != nil {
		return err
	}

	if err := s.store.UpdateRefreshToken(rt); err != nil {
		return err
	}
	resp.Token = at
	resp.RefreshToken = rt.Token

	re := RefreshEvent{
		Type:      "refresh",
		AccountNr: acc.Number,
	}
	msg := sse.Message{Data: re}
	s.broker.Publish(msg)

	return WriteJSON(w, http.StatusOK, resp)
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

	rt, err := s.store.GetRefreshToken(acc.ID)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token:        tokenString,
		RefreshToken: rt.Token,
	}

	return WriteJSON(w, http.StatusOK, resp)
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

func (s *APIServer) handleImage(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handlePostImage(w, r)
	}
	if r.Method == "GET" {
		return s.handleGetImage(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handlePostImage(w http.ResponseWriter, r *http.Request) error {
	num, err := getNumber(r)
	if err != nil {
		return err
	}
	req := new(ImageRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	acc, err := s.store.GetAccountByNumber(num)
	if err != nil {
		return err
	}
	if err := s.store.AddImage(req.Image, acc.ImageName); err != nil {
		return err
	}
	fmt.Println(req.Image)
	return WriteJSON(w, http.StatusOK, acc)
}

func (s *APIServer) handleGetImage(w http.ResponseWriter, r *http.Request) error {
	num, err := getNumber(r)
	if err != nil {
		return err
	}
	fmt.Println("Getting account")
	acc, err := s.store.GetAccountByNumber(num)
	if err != nil {
		return err
	}
	fmt.Println("Getting image")
	image, err := s.store.GetImage(acc.ImageName)
	if err != nil {
		return err
	}
	resp := new(ImageResponse)
	resp.Image = image
	return WriteJSON(w, http.StatusOK, resp)
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
	imageName := s.store.NewImageForAccount(acc.Number)
	acc.ImageName = imageName

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
	return WriteJSON(w, http.StatusCreated, acc)
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
	trn := TransactionEvent{
		Type:      "transaction",
		Sender:    sender.Number,
		Recipient: recipient.Number,
		Amount:    amnt,
	}
	s.broker.Publish(sse.Message{Data: trn})
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, req)
}

// Authentication Middleware Adapted from Anthony GG's tutorial
func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := parseJWT(tokenString)
		if err != nil && token != nil && !token.Valid {
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
		handlerFunc(w, r)
	}
}

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"exp":            time.Now().Add(45 * time.Second).Unix(),
		"account_number": account.Number,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(JWT_SECRET))
}

func parseJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})
}

// func getID(r *http.Request) (int, error) {
// 	idStr := mux.Vars(r)["id"]

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		return id, fmt.Errorf("invalid id given %s", idStr)
// 	}
// 	return id, nil
// }

func getNumber(r *http.Request) (int, error) {
	numberStr := mux.Vars(r)["number"]
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return number, fmt.Errorf("invalid number given %s", numberStr)
	}
	return number, nil
}

// ChatGPT Aided
// Reference: https://blog.devtrovert.com/p/reflection-in-go-everything-you-need (flexible update)
func updateAccountWithReflect(acc *Account, req *UpdateAccountRequest) {
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
