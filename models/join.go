package models

type GifsTags struct {
	tableName struct{} `sql:"gifs_tags"`

	TagID int `sql:",pk"`
	Tag   *Tag
	GifID int `sql:",pk"`
	Gif   *Gif
}

type UsersAchievements struct {
	tableName struct{} `sql:"users_achievements"`

	AchievementID int `sql:",pk"`
	Achievement   *Achievement
	UserID        int `sql:",pk"`
	User          *User
	Score         int
	Unlocked      bool
}

type UsersFriends struct {
	tableName struct{} `sql:"users_friends"`

	UserID   int `sql:",pk"`
	User     *User
	FriendID int `sql:",pk"`
	Friend   *User
}

type UsersRoles struct {
	tableName struct{} `sql:"users_roles"`

	RoleID int `sql:",pk"`
	Role   *Role
	UserID int `sql:",pk"`
	User   *User
}
