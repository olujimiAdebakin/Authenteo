package constants

type Type string

const (
    Type2FA           Type = "2fa"
    TypePasswordReset Type = "password_reset"
    TypeEmailVerify   Type = "email_verify"
)