package apierrors

type ApiError struct {
	Message string
}

var (
	ErrInternalServerError = ApiError{
		Message: "処理中に問題が発生しました。しばらく経ってから再度お試しください。",
	}
	ErrInvalidParameter = ApiError{
		Message: "無効なパラメータが検出されました。",
	}
	ErrValidation = ApiError{
		Message: "入力に誤りがあります。",
	}
	ErrTooManyRequest = ApiError{
		Message: "一定期間に多くのリクエストを検出しました。しばらく経ってから再度お試しください。",
	}
)

var (
	ErrUserAlreadyRegistered = ApiError{Message: "ユーザーはすでに登録されています。"}
	ErrUserNotFound          = ApiError{Message: "該当のユーザーが見つかりませんでした。"}

	ErrConfirmationTokenExpired = ApiError{Message: "登録用トークンの期限がすぎています。再度サインアップをお願いします。"}
)
