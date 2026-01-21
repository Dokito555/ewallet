package http

import (
	"net/http"

	"github.com/Dokito555/ewallet/ewallet-ums/constants"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/entity"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/models"
	"github.com/Dokito555/ewallet/ewallet-ums/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	Log      *logrus.Logger
	Service  *services.UserService
	Validate *validator.Validate
}

func NewUserController(log *logrus.Logger, svc *services.UserService, val *validator.Validate) *UserController {
	return &UserController{
		Log:      log,
		Service:  svc,
		Validate: val,
	}
}

func (c *UserController) Login(ctx *gin.Context) {
	var (
		req models.LoginRequest
		res *models.LoginResponse
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Info("Failed to parse request: ", err)
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		c.Log.Info("Failed to validate request: ", err)
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	res, err := c.Service.Login(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Info("Failed on login service: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, res)
}


func (c *UserController) Logout(ctx *gin.Context) {
	token := ctx.Request.Header.Get("Authorization")
	err := c.Service.Logout(ctx.Request.Context(), token)
	if err != nil {
		c.Log.Info("Failed on login service: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, constants.SUCCESSMessage)
	return
}


func (c *UserController) RefreshToken(ctx *gin.Context) {
	refreshToken := ctx.Request.Header.Get("Authorization")
	claim, ok := ctx.Get("auth")
	if !ok {
		c.Log.Info("Failed to get claim in context: ")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	tokenClaim, ok := claim.(*models.ClaimToken)
	if !ok {
		c.Log.Info("Failed on refresh token service: ")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	resp, err := c.Service.RefreshToken(ctx.Request.Context(), refreshToken, *tokenClaim)
	if err != nil {
		c.Log.Info("Failed to parse claim token: ")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, resp)
	return
}

func (c *UserController) Register(ctx *gin.Context) {
	req := entity.User{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Info("Failed to parse request: ", err)
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	if err := c.Validate.Struct(req); err != nil {
		c.Log.Info("Failed to validate request: ", err)
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	data, err := c.Service.Register(ctx.Request.Context(), req)
	if err != nil {
		c.Log.Error("Failed to register user: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, data)
	return
}

// func (s *UserController) ValidateToken(ctx context.Context, req *token_validation_proto.TokenRequest) (*token_validation_proto.TokenResponse, error) {
// 	var (
// 		token = req.GetToken()
// 	)

// 	if s.TokenValidationService == nil {
// 		err := fmt.Errorf("TokenValidationService is not initialized")
// 		return &token_validation_proto.TokenResponse{
// 			Message: err.Error(),
// 		}, nil
// 	}

// 	if token == "" {
// 		err := fmt.Errorf("token is empty")
// 		return &token_validation_proto.TokenResponse{
// 			Message: err.Error(),
// 		}, nil
// 	}

// 	claimToken, err := s.TokenValidationService.TokenValidation(ctx, token)
// 	if err != nil {
// 		return &token_validation_proto.TokenResponse{
// 			Message: err.Error(),
// 		}, nil
// 	}

// 	return &token_validation_proto.TokenResponse{
// 		Message: constants.SUCCESSMessage,
// 		Data: &token_validation_proto.UserData{
// 			UserId:   int64(claimToken.UserID),
// 			Username: claimToken.Username,
// 			FullName: claimToken.FullName,
// 			Email:    claimToken.Email,
// 		},
// 	}, nil
// }
