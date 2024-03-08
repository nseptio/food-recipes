package service

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nseptio/food-recipes/internal/config"
	"github.com/nseptio/food-recipes/internal/dto"
	"github.com/nseptio/food-recipes/internal/entity"
	"github.com/nseptio/food-recipes/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type UserService struct {
	DB             *gorm.DB
	UserRepository *repository.UserRepository
	Validate       validator.Validate
}

func NewUserService(DB *gorm.DB, userRepository *repository.UserRepository, validate validator.Validate) *UserService {
	return &UserService{
		DB:             DB,
		UserRepository: userRepository,
		Validate:       validate,
	}
}

func (s *UserService) FindUserById(ctx context.Context, id int) (*entity.User, error) {
	user := new(entity.User)
	if err := s.UserRepository.FindById(ctx, s.DB, user, id); err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}
	return user, nil
}

func (s *UserService) Register(ctx context.Context, request *dto.UserRegisterRequest) (*dto.UserResponse, error) {

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.Validate.Struct(request); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	total, err := s.UserRepository.CountByUsernameAndEmail(ctx, s.DB, request.Username, request.Email)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to count user from database")
	}

	if total > 0 {
		return nil, fiber.NewError(fiber.StatusConflict, "User already exist")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to generate bcrypt hash for password")
	}

	user := entity.User{
		Username:  request.Username,
		Email:     request.Email,
		Password:  string(password),
		FirstName: request.FirstName,
		LastName:  request.LastName,
	}

	if err := s.UserRepository.Create(s.DB, &user); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed create a user to database")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed commit transaction")
	}

	return dto.UserToResponse(user), nil
}

func (s *UserService) Login(ctx context.Context, request *dto.LoginRequest) (*dto.LoginResponse, error) {

	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()
	// TODO: validate request body

	//	Get email and password from user request
	user := new(entity.User)
	identifier := request.Identifier

	err := s.Validate.Var(identifier, "email")
	if err != nil {
		if err := s.UserRepository.FindByUsername(ctx, tx, user, identifier); err != nil {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Failed find user by username")
		}
	} else {
		if err := s.UserRepository.FindByEmail(ctx, tx, user, identifier); err != nil {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Failed find user by email")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed commit transaction")
	}

	// compare
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Password not match")
	}

	viper := config.NewViper()
	hmacSampleSecret := viper.GetString("app.jwt.key")
	// create jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(hmacSampleSecret))

	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to create token")
	}

	response := &dto.LoginResponse{Token: tokenString}
	return response, nil
}

func (s *UserService) LikeRecipe(ctx context.Context, user *entity.User, recipe *entity.Recipe) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.UserRepository.LikeRecipe(ctx, tx, user, recipe); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to like recipe")
	}

	if err := tx.Commit().Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed commit transaction")
	}

	return nil
}

func (s *UserService) UnlikeRecipe(ctx context.Context, user *entity.User, recipe *entity.Recipe) error {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := s.UserRepository.UnlikeRecipe(ctx, tx, user, recipe); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to unlike recipe")
	}

	if err := tx.Commit().Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed commit transaction")
	}

	return nil
}

func (s *UserService) GetLikedRecipes(ctx context.Context, user *entity.User) ([]entity.Recipe, error) {
	tx := s.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	recipes, err := s.UserRepository.GetLikedRecipes(ctx, tx, user)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to get liked recipes")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed commit transaction")
	}

	return recipes, nil
}
