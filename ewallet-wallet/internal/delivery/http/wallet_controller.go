package http

import (
	"net/http"
	"strconv"

	"github.com/Dokito555/ewallet/ewallet-ewallet/constants"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/entity"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/models"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type WalletController struct {
	Log *logrus.Logger
	Service *services.WalletService
	Validate *validator.Validate
}

func NewWalletController(log *logrus.Logger, svc *services.WalletService, val *validator.Validate) *WalletController {
	return &WalletController{
		Log: log,
		Service: svc,
		Validate: val,
	}
}

func (c *WalletController) Create(ctx *gin.Context) {
	var (
		req entity.Wallet
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Error("failed to parse request: ", err)
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	if req.UserID == 0 {
		c.Log.Error("user id is empty")
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	err := c.Service.Create(&req)
	if err != nil {
		c.Log.Error("failed to create wallet: ", err)
		ctx.JSON( http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON( http.StatusOK, req)
}

func (c *WalletController) CreditBalance(ctx *gin.Context) {
	var (
		req models.TransactionRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Error("failed to parse request: ", err)
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		c.Log.Error("failed to validate request: ", err)
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	token, ok := ctx.Get("token")
	if !ok {
		c.Log.Error("failed to get token")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}
	tokenData, ok := token.(models.TokenData)
	if !ok {
		c.Log.Error("failed to parse token")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	resp, err := c.Service.CreditBalance(int(tokenData.UserID), req)
	if err != nil {
		c.Log.Error("failed to credit balancet: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *WalletController) DebitBalance(ctx *gin.Context) {
	var (
		req models.TransactionRequest
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.Log.Error("failed to parse request: ", err)
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		c.Log.Error("failed to validate request: ", err)
		ctx.JSON(http.StatusBadRequest, constants.ErrFailedBadRequest)
		return
	}

	token, ok := ctx.Get("token")
	if !ok {
		c.Log.Error("failed to get token")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}
	tokenData, ok := token.(models.TokenData)
	if !ok {
		c.Log.Error("failed to parse token")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	resp, err := c.Service.DebitBalance(int(tokenData.UserID), req)
	if err != nil {
		c.Log.Error("failed to debit balance: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *WalletController) GetBalance(ctx *gin.Context) {
	token, ok := ctx.Get("token")
	if !ok {
		c.Log.Error("failed to get token")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}
	tokenData, ok := token.(models.TokenData)
	if !ok {
		c.Log.Error("failed to parse token")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	resp, err := c.Service.GetBalance(int(tokenData.UserID))
	if err != nil {
		c.Log.Error("failed to get balance: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *WalletController) GetWalletHistory(ctx *gin.Context) {
	var (
		param models.WalletHistoryParam
	)

	if err := ctx.ShouldBindQuery(&param); err != nil {
		c.Log.Error("failed to parse query param")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}

	if param.WalletTransactionType != "CREDIT" && param.WalletTransactionType != "DEBIT" || param.WalletTransactionType == "" {
		c.Log.Error("invalid wallet_transaction_type")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}

	token, ok := ctx.Get("token")
	if !ok {
		c.Log.Error("failed to get token")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}
	tokenData, ok := token.(models.TokenData)
	if !ok {
		c.Log.Error("failed to parse token")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	resp, err := c.Service.GetWalletHistory(int(tokenData.UserID), param)
	if err != nil {
		c.Log.Error("failed to get wallet history: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *WalletController) CreateWalletLink(ctx *gin.Context) {
	var (
		req entity.WalletLink
	)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.Log.Error("failed to parse req")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}

	clientID, ok := ctx.Get("client_id")
	if !ok {
		c.Log.Error("failed to get client id")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}
	clientSource, ok := clientID.(string)
	if !ok {
		c.Log.Error("failed to parse client id")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	resp, err := c.Service.CreateWalletLink(clientSource, &req)
	if err != nil {
		c.Log.Error("failed to create wallet link: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *WalletController) WalletLinkConfirmation(ctx *gin.Context) {
	var (
		req models.WalletStructOTP
	)

	if err := ctx.ShouldBindQuery(&req); err != nil {
		c.Log.Error("failed to parse req")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}

	walletIDs := ctx.Param("wallet_id")
	if walletIDs == "" {
		c.Log.Error("failed to get wallet id")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}
	walletID, err := strconv.Atoi(walletIDs)
	if err != nil {
		c.Log.Error("failed to parse wallet id to int")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}

	clientID, ok := ctx.Get("client_id")
	if !ok {
		c.Log.Error("failed to get client id")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}
	clientSource, ok := clientID.(string)
	if !ok {
		c.Log.Error("failed to parse client id")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	err = c.Service.WalletLinkConfirmation(walletID, clientSource, req.OTP)
	if err != nil {
		c.Log.Error("failed to confirm wallet link: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (c *WalletController) WalletUnlink(ctx *gin.Context) {
	walletIDs := ctx.Param("wallet_id")
	if walletIDs == "" {
		c.Log.Error("failed to get wallet id")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}
	walletID, err := strconv.Atoi(walletIDs)
	if err != nil {
		c.Log.Error("failed to parse wallet id to int")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}

	clientID, ok := ctx.Get("client_id")
	if !ok {
		c.Log.Error("failed to get client id")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}
	clientSource, ok := clientID.(string)
	if !ok {
		c.Log.Error("failed to parse client id")
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	err = c.Service.WalletUnlink(walletID, clientSource)
	if err != nil {
		c.Log.Error("failed to unlink wallet: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (c *WalletController) ExGetBalance(ctx *gin.Context) {
	walletIDs := ctx.Param("wallet_id")
	if walletIDs == "" {
		c.Log.Error("failed to get wallet id")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}
	walletID, err := strconv.Atoi(walletIDs)
	if err != nil {
		c.Log.Error("failed to parse wallet id to int")
		ctx.JSON(http.StatusInternalServerError, constants.ErrFailedBadRequest)
		return
	}

	resp, err := c.Service.ExGetBalance(walletID)
	if err != nil {
		c.Log.Error("failed to unlink wallet: ", err)
		ctx.JSON(http.StatusInternalServerError, constants.ErrServerError)
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
