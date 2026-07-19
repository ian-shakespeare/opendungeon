package providers

import (
	"context"
	"errors"
	"io"
	"net/http"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
	ErrMissingAvatar        = errors.New("missing avatar")
	ErrAvatarNotFound       = errors.New("avatar not found")
)

type UserInfo struct {
	ID        string
	Email     string
	Username  string
	AvatarUri *string
}

type Provider interface {
	AuthUrl(state string) string
	Exchange(ctx context.Context, code string) (UserInfo, error)
}

func GetAvatar(user UserInfo) (io.ReadCloser, error) {
	if user.AvatarUri == nil {
		return nil, ErrMissingAvatar
	}

	req, err := http.NewRequest(http.MethodGet, *user.AvatarUri, http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrAvatarNotFound
	}

	return res.Body, nil
}
