package utils

import (
	"errors"
	"os"
	"reminder-server/internal/models"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var validPrioties = []string{models.PriorityLow, models.PriorityMedium, models.PriorityHigh}
var valiudStatuses = []string{models.StatusPending, models.StatusCompleted, models.StatusOverdue}

func GenerateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET is not set")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetCurrentTime() time.Time {
	return time.Now()
}

func IsValidPriority(priority string) bool {
	return slices.Contains(validPrioties, priority)
}

func IsValidStatus(status string) bool {
	return slices.Contains(valiudStatuses, status)
}
