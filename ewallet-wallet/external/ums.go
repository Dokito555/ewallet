package external

import (
	"context"
	"fmt"

	"github.com/Dokito555/ewallet/ewallet-ewallet/constants"
	token_validation_proto "github.com/Dokito555/ewallet/ewallet-ewallet/internal/delivery/grpc/proto/token_validation"
	"github.com/Dokito555/ewallet/ewallet-ewallet/internal/models"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ExtUMS struct {
	Log *logrus.Logger
}

func NewExtUMS(log *logrus.Logger) *ExtUMS {
	return &ExtUMS{
		Log: log,
	}
}

func (e *ExtUMS) ValidateToken(ctx context.Context, token string) (models.TokenData, error) {
	var (
		resp models.TokenData
	)

	conn, err := grpc.Dial("UMS_GRPC_HOST", grpc.WithInsecure())
	if err != nil {
		e.Log.Warnf("failed to dial ums grpc: ", err)
		return resp, fmt.Errorf("failed to dial ums grpc: ", err)
	}
	defer conn.Close()

	client := token_validation_proto.NewTokenValidationClient(conn)

	req := &token_validation_proto.TokenRequest{
		Token: token,
	}

	response, err := client.ValidateToken(ctx, req)
	if err != nil {
		e.Log.Warnf("failed to validate token: ", err)
		return resp, fmt.Errorf("failed to validate token")
	}

	if response.Message != constants.SUCCESSMessage {
		e.Log.Warnf("got response from ums: ", err)
		return resp, fmt.Errorf("got response error from ums: %s", response.Message)
	}

	resp.UserID = response.Data.UserId
	resp.Username = response.Data.Username
	resp.FullName = response.Data.FullName
	resp.Email = response.Data.Email

	return resp, nil
}
