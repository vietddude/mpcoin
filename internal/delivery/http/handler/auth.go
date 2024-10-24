package handler

import (
	_ "mpc/docs"
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

// Login godoc
// @Summary User Login
// @Description Authenticates a user and returns access and refresh tokens along with user details.
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body domain.LoginRequest true "Login Request containing email and password"
// @Success 200 {object} docs.LoginResponse "Successful login response with user details, access token, and refresh token"
// @Failure 400 {string} string "Bad request error due to invalid input"
// @Failure 401 {string} string "Unauthorized error due to incorrect email or password"
// @Router /auth/login [post]
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

// Signup godoc
// @Summary User Signup
// @Description Registers a new user and returns user details, wallet details, access token, and refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Param signupRequest body domain.SignupRequest true "Signup Request containing email and password"
// @Success 201 {object} docs.SignupResponse "Successful signup response with user details, wallet details, access token, and refresh token"
// @Failure 400 {string} string "Bad request error due to invalid input"
// @Failure 500 {string} string "Internal server error"
// @Router /auth/signup [post]
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

// Logout godoc
// @Summary User Logout
// @Description Logs out a user by invalidating the refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} string "Logged out successfully"
// @Failure 401 {string} string "Unauthorized error due to invalid token"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := utils.GetAuthToken(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	h.authUC.Logout(c, token)
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// Refresh godoc
// @Summary Refresh Token
// @Description Refreshes the access token using the refresh token.
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} docs.RefreshResponse "Successful response with new access token and refresh token"
// @Failure 401 {string} string "Unauthorized error due to invalid token"
// @Router /auth/refresh [post]
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
