package public

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-07-22

type User struct {
	Id             int64  `json:"id,omitempty"`
	Username       string `json:"username,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

type RegisterRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type RegisterResponse struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type LoginResponse struct {
	Code  int32  `json:"code,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Token string `json:"token,omitempty"`
	User  *User  `json:"user,omitempty"`
}

type CheckTokenRequest struct {
	Token string `json:"token,omitempty"`
}

type CheckTokenResponse struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}

type GetProfileRequest struct {
	Token string `json:"token,omitempty"`
}

type GetProfileResponse struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	User *User  `json:"user,omitempty"`
}

type UpdateProfileRequest struct {
	Token          string `json:"token,omitempty"`
	Username       string `json:"username,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

type UpdateProfileResponse struct {
	Code int32  `json:"code,omitempty"`
	Msg  string ` json:"msg,omitempty"`
}

type LogoutRequest struct {
	Token string `json:"token,omitempty"`
}

type LogoutResponse struct {
	Code int32  `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
}
