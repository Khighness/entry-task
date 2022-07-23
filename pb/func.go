package pb

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-22

const (
	FuncRegister      = "register"
	FuncLogin         = "login"
	FuncCheckToken    = "checkToken"
	FuncGetProfile    = "getProfile"
	FuncUpdateProfile = "updateProfile"
	FuncLogout        = "logout"
)

var (
	Register      func(request RegisterRequest) (RegisterResponse, error)
	Login         func(request LoginRequest) (LoginResponse, error)
	CheckToken    func(request CheckTokenRequest) (CheckTokenResponse, error)
	GetProfile    func(request GetProfileRequest) (GetProfileResponse, error)
	UpdateProfile func(request UpdateProfileRequest) (UpdateProfileResponse, error)
	Logout        func(request LogoutRequest) (LogoutResponse, error)
)
