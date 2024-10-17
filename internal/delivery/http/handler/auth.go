package handler

import (
	"mpc/internal/domain"
	"mpc/internal/usecase"
	"mpc/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC usecase.AuthUseCase
}

func NewAuthHandler(authUC usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUC: authUC,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	loginRequest, err := utils.ParseRequest[domain.LoginRequest](c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, accessToken, refreshToken, err := h.authUC.Login(c, loginRequest.Email, loginRequest.Password)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	loginResponse := domain.LoginResponse(user)

	utils.SuccessResponse(c, http.StatusOK, gin.H{"user": loginResponse, "access_token": accessToken, "refresh_token": refreshToken})
}

func (h *AuthHandler) Signup(c *gin.Context) {
	signupRequest, err := utils.ParseRequest[domain.SignupRequest](c)

	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user, wallet, accessToken, refreshToken, err := h.authUC.Signup(c, domain.CreateUserParams(signupRequest))

	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	signupResponse := domain.SignupResponse{
		ID:    user.ID,
		Email: user.Email,
	}

	createWalletResponse := domain.CreateWalletResponse{
		ID:      wallet.ID,
		UserID:  wallet.UserID,
		Address: wallet.Address,
	}

	utils.SuccessResponse(c, http.StatusCreated, gin.H{"user": signupResponse, "wallet": createWalletResponse, "access_token": accessToken, "refresh_token": refreshToken})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := utils.GetAuthToken(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.authUC.Logout(c, token)
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	token, err := utils.GetAuthToken(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	accessToken, refreshToken, err := h.authUC.RefreshToken(c, token)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"access_token": accessToken, "refresh_token": refreshToken})
}
