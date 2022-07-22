package public

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
	Register      func(request *RegisterRequest) *RegisterResponse
	Login         func(request *LoginRequest) *LoginResponse
	CheckToken    func(request *CheckTokenRequest) *CheckTokenResponse
	GetProfile    func(request *GetProfileRequest) *GetProfileResponse
	UpdateProfile func(request *UpdateProfileRequest) *UpdateProfileResponse
	Logout        func(request LogoutRequest) *LogoutResponse
)
