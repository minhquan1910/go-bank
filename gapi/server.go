package gapi

import (
	db "bank/db/sqlc"
	"bank/pb"
	"bank/token"
	"bank/util"
	"bank/worker"
)

// Server serve Http request for our services
type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	tokenMaker      token.Maker
	config          util.Config
	taskDistributor worker.TaskDistributor
}

// Create new Http server and setup routing
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
