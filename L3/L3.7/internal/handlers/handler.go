package handlers

import (
	"log/slog"

	"github.com/MAPiryazev/Wildberries_L1/L3/L3.7/internal/service"
)

func New(
	itemService service.ItemService,
	userService service.UserService,
	jwtSecret string,
	log *slog.Logger,
) *Handler {
	return &Handler{
		itemService: itemService,
		userService: userService,
		cfg:         &Config{JWTSecret: jwtSecret},
		log:         log,
	}
}
