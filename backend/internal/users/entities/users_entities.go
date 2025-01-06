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
