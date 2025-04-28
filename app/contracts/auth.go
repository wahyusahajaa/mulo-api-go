package contract

type AuthRepository interface {
	Store()
}

type AuthService interface {
	Register()
}
