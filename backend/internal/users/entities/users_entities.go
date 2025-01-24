package entities

type Follow struct {
	ID   int    `json:"id"`   // ID of the user being followed
	Name string `json:"name"` // Name of the user being followed
}

type UserWithStats struct {
	ID            int      `json:"id" db:"id"`
	Name          string   `json:"name" db:"name"`
	Bio           string   `json:"bio" db:"bio"`
	ProfileImage  string   `json:"profile_image" db:"profile_image"`
	FollowerCount int      `json:"follower_count" db:"follower_count"`
	Following     []Follow `json:"following" db:"following"`
}

type UserAddFollowingRequest struct {
	FollowerID  int `json:"follower_id"`
	FollowingID int `json:"following_id"`
}

type LoginRequest struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password_hash"`
}

type UserCredentials struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	Password string `db:"password_hash"`
	Role     string `db:"role"`
}

type User struct {
	ID           int    `db:"id"`
	Name         string `db:"name"`
	Email        string `db:"email"`
	Password     string `db:"password"`
	Role         string `db:"role"`
	Bio          string `db:"bio"`
	ProfileImage string `db:"profile_image"`
}

type RegisterRequest struct {
	Name         string `json:"name" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_image"`
}

// UserToken Struct
type UserToken struct {
	ID           string `json:"id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UserProfileResponse Struct
type UserProfileResponse struct {
	User  *User      `json:"user"`
	Token *UserToken `json:"token"`
}
