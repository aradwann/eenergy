package api

// import (
// 	"fmt"

// 	"github.com/aradwann/eenergy/pb"
// 	db "github.com/aradwann/eenergy/repository/store"
// 	"github.com/aradwann/eenergy/token"
// 	"github.com/aradwann/eenergy/util"
// 	"github.com/aradwann/eenergy/worker"
// )

// // Server serves gRPC requests for our eenergy service.
// type Server struct {
// 	pb.UnimplementedEenergyServiceServer
// 	config          util.Config
// 	store           db.Store
// 	tokenMaker      token.Maker
// 	taskDistributor worker.TaskDistributor
// }

// // NewServer creates a new gRPC server.
// func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
// 	tokenMaker, err := token.NewPASETOMaker(config.TokenSymmetricKey)
// 	if err != nil {
// 		return nil, fmt.Errorf("cannot create token maker: %w", err)
// 	}

// 	server := &Server{
// 		config:          config,
// 		store:           store,
// 		tokenMaker:      tokenMaker,
// 		taskDistributor: taskDistributor,
// 	}

// 	return server, nil
// }
