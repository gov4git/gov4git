package history

type User string

type JoinEvent struct {
	User User `json:"user"`
}
