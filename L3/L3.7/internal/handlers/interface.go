package handlers

import (
	"log/slog"
	"net/http"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/service"
)

type ItemHandler interface {
	CreateItem(w http.ResponseWriter, r *http.Request)
	GetItem(w http.ResponseWriter, r *http.Request)
	ListItems(w http.ResponseWriter, r *http.Request)
	UpdateItem(w http.ResponseWriter, r *http.Request)
	DeleteItem(w http.ResponseWriter, r *http.Request)
	GetItemHistory(w http.ResponseWriter, r *http.Request)
}

type AuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	GetCurrentUser(w http.ResponseWriter, r *http.Request)
	ListUsers(w http.ResponseWriter, r *http.Request)
}

type Config struct {
	JWTSecret string
}

type Handler struct {
	itemService service.ItemService
	userService service.UserService
	cfg         *Config
	log         *slog.Logger
}
