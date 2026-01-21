package repositories

import (
	"errors"

	"github.com/Dokito555/ewallet/ewallet-ums/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.User]
	Log *logrus.Logger
	DB *gorm.DB
}

func NewUserRepository(log *logrus.Logger, db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository: Repository[entity.User]{DB: db},
		Log: log,
		DB: db,
	}
}


func (r *UserRepository) GetUserByUsername(username string) (*entity.User, error) {
	var (
		user entity.User
		err error
	)

	err = r.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return &user, err
	}

	if user.ID == 0 {
		return &user, errors.New("user not found")
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var (
		user entity.User
		err error
	)

	err = r.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return &user, err
	}

	if user.ID == 0 {
		return &user, errors.New("user not found")
	}

	return &user, nil
}


func (r *UserRepository) InsertNewUserSession(userSession entity.UserSession) error {
	return r.DB.Create(userSession).Error
}

func (r *UserRepository) DeleteUserSession(token string) error {
	return r.DB.Exec("DELETE FROM user_sessions WHERE token = ?", token).Error
}

func (r *UserRepository) UpdateTokenByRefreshToken(token, refreshToken string) error {
	return r.DB.Exec("UPDATE user_sessions SET token = ? WHERE refresh_token", token, refreshToken).Error
}

func (r *UserRepository) GetUserSessionByToken(token string) (entity.UserSession, error) {
	var (
		session entity.UserSession
		err error
	)

	err = r.DB.Where("token = ?", token).First(&session).Error
	if err != nil {
		return session, err
	}

	if session.ID == 0 {
		return session, errors.New("user not found")
	}

	return session, nil
}

func (r *UserRepository) GetUserByRefreshToken(refreshToken string) (entity.UserSession, error) {
	var (
		session entity.UserSession
		err error
	)

	err = r.DB.Where("refresh_token = ?", refreshToken).First(&session).Error
	if err != nil {
		return session, err
	}

	if session.ID == 0 {
		return session, errors.New("user not found")
	}

	return session, nil
}