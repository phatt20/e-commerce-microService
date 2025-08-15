package user 

import "time"

type (
	UserProfile struct {
		Id        string    `json:"_id"`
		Email     string    `json:"email"`
		Username  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserClaims struct {
		Id       string `json:"id"`
		RoleCode int    `json:"role_code"`
	}

	CreateUserReq struct {
		Email    string `json:"email" form:"email" validate:"required,email,max=255"`
		Password string `json:"password" form:"password" validate:"required,max=32"`
		Username string `json:"username" form:"username" validate:"required,max=64"`
	}

	CreateUserTransactionReq struct {
		UserId string  `json:"user_id" validate:"required,max=64"`
		Amount   float64 `json:"amount" validate:"required"`
	}

	RollbackUserTransactionReq struct {
		TransactionId string `json:"transaction_id"`
	}
)
