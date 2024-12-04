package handler

import (
	"context"
	"workshop-1/internal/infrastructure/metrics"
	"workshop-1/internal/usecase/auth"
	"workshop-1/pkg/pvz/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *PVZHandler) RegisterUser(ctx context.Context, req *pvz.RegisterRequest) (*pvz.RegisterResponse, error) {
	handlerName := "RegisterUser"
	var errCode codes.Code

	defer func() {
		if errCode == codes.OK {
			metrics.IncOkResponseTotal(handlerName)
		} else {
			metrics.IncErrResponseTotal(handlerName, errCode.String())
		}
	}()

	if err := req.ValidateAll(); err != nil {
		errCode = codes.InvalidArgument
		return nil, status.Error(errCode, err.Error())
	}

	token, err := auth.CreateAccessToken(int(req.GetUserId()))
	if err != nil {
		errCode = codes.Internal
		return nil, status.Error(errCode, err.Error())
	}

	return &pvz.RegisterResponse{ApiToken: token}, nil
}
