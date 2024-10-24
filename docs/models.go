package docs

import "mpc/internal/domain"

type LoginResponse struct {
	User         domain.LoginResponse `json:"user"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignupResponse struct {
	User         domain.SignupResponse `json:"user"`
	Wallet       domain.CreateWalletResponse `json:"wallet"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateTxnResponse struct {
	TransactionId   string `json:"txn_id"`
	Message    		string `json:"message"`
}

type SubmitTnxResponse struct {
	TransactionId   string `json:"tnx_hash"`
	Message    		string `json:"message"`
}


