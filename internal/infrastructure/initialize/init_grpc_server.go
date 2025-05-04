package initialize

import (
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/global"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitCommentGrpcConn() *grpc.ClientConn {
	connection := fmt.Sprintf("%s:%d", global.Config.CommentCensorGrpcConn.Host, global.Config.CommentCensorGrpcConn.Port)
	fmt.Printf("Connecting to %s\n", connection)
	conn, err := grpc.NewClient(connection, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		global.Logger.Error("Failed to init grpc server", zap.Error(err))
		panic(err)
	}

	return conn
}
