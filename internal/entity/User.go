package entity

import "context"

type User struct {
	ChatID          int
	YandexDiskToken string
	State           string
}

type UserRepository interface {
	Create(ctx context.Context, chat *User) error
	Update(ctx context.Context, chat *User) error
	FindByChatID(ctx context.Context, chatID int) (User, error)
}

type UserUseCase interface {
	CreateUser(ctx context.Context, chatID int) error
}
