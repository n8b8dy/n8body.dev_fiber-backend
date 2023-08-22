package utils

import (
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GetUserIDFromContextWithJWT(ctx *fiber.Ctx) (uuid.UUID, error) {
	user := ctx.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.MapClaims)

	id, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, err
}
