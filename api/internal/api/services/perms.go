package services

import (
	"context"
	"fmt"
	ssov1 "github.com/Goose47/go-grpc-sso.protos/gen/go/sso"
	"google.golang.org/grpc"
	"log/slog"
)

type PermsService struct {
	log        *slog.Logger
	AuthClient ssov1.AuthClient
}

func NewPermsService(
	log *slog.Logger,
	authAddress string,
) (*PermsService, error) {
	//cc, err := grpc.NewClient(authAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	cc, err := grpc.NewClient(authAddress)
	if err != nil {
		return nil, fmt.Errorf("grpc connection failed: %s", err.Error())
	}

	authClient := ssov1.NewAuthClient(cc)

	return &PermsService{
		log:        log,
		AuthClient: authClient,
	}, nil
}

func (a *PermsService) IsAdmin(userID int64) (bool, error) {
	const op = "api.services.IsAdmin"

	log := a.log.With(slog.Int("user_id", int(userID)))

	res, err := a.AuthClient.IsAdmin(context.TODO(), &ssov1.IsAdminRequest{UserId: userID})
	if err != nil {
		log.Warn("failed to check whether user is admin")
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return res.IsAdmin, nil
}
