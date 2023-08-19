package usecase

import (
	"context"
	"fmt"

	"github.com/V1taly5/APIYDisk/internal/entity"
	"go.mongodb.org/mongo-driver/mongo"
)

type userUseCase struct {
	userRepository entity.UserRepository
}

func NewUserUseCase(u entity.UserRepository) entity.UserUseCase {
	return &userUseCase{u}
}

func (useCase *userUseCase) CreateUser(ctx context.Context, chatID int) error {
	const op = "usecase.user.CreateUser"
	_, err := useCase.userRepository.FindByChatID(ctx, chatID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			var user entity.User
			user.ChatID = chatID
			user.State = ""
			useCase.userRepository.Create(ctx, &user)
			return nil
		}
		return err
	}

	// указвать: пользователь уже существует
	return fmt.Errorf("user already exists")

}

func (useCase *userUseCase) InsertDiskToken(ctx context.Context, user entity.User) error {
	const op = "usecase.user.InsertDiskToken"
	_, err := useCase.userRepository.FindByChatID(ctx, user.ChatID)
	if err != nil {
		return err
	}

	err = useCase.userRepository.Update(ctx, &user)
	if err != nil {
		return err
	}
	return nil
}
