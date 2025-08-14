package service

import (
	"errors"
	"time"
	"kholid/ewallet/v2/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct { db *gorm.DB; secret string }
func NewAuthService(db *gorm.DB, secret string) *AuthService { return &AuthService{db: db, secret: secret} }

func (a *AuthService) Register(email, password string) (*models.User, error) {
	if email=="" || password=="" { return nil, errors.New("email and password required") }
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := &models.User{Email: email, Password: string(hash)}
	if err := a.db.Create(u).Error; err != nil { return nil, err }
	return u, nil
}
func (a *AuthService) Login(email, password string) (string, *models.User, error) {
	var u models.User
	if err := a.db.Where("email = ?", email).First(&u).Error; err != nil { return "", nil, errors.New("invalid credentials") }
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil { return "", nil, errors.New("invalid credentials") }
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": u.ID, "exp": time.Now().Add(24*time.Hour).Unix()})
	s, err := t.SignedString([]byte(a.secret)); return s, &u, err
}
