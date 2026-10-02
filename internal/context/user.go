package context

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

var users = []User{
	{ID: 1, Username: "infamous", Role: "admin"},
	{ID: 2, Username: "syed", Role: "user"},
}

func CheckUserExists(username string) *User {
	var user *User
	for _, u := range users {
		if u.Username == username {
			user = &u
			return user
		}
	}
	return nil
}
