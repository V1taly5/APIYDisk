package entity

type User struct {
	ChatID          int
	YandexDiskToken string
	State           string
}

type UserRepository interface {
	Find(chatID int) (bool, error)
	Create(chat *User) error
	Update(chat *User) error
	GetUser(chatID int64) (User, error)
	FindByChatID(chatID int) (User, error)
}
