package main

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"log"
	pb "my_grpc/api"
	"my_grpc/internal/db"
	"net"
	"os"
	"os/signal"
)

type Po interface {
	QueryUserInfo(id string) (*db.UserInfo, error)
}

type server struct{}

type userInfoDto struct {
	name    string
	age     int
	address string
}

func convertPoUserInfoToDto(userInfo *db.UserInfo) (*userInfoDto, error) {
	if userInfo == nil {
		return nil, errors.New("the db userinfo is nil")
	}
	return &userInfoDto{
		name:    userInfo.Name,
		age:     userInfo.Age,
		address: userInfo.Address,
	}, nil
}

var po Po

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if in == nil {
		return &pb.HelloReply{}, errors.New("invalid request")
	}
	userInfo, err := po.QueryUserInfo(in.Id)
	if err != nil {
		return &pb.HelloReply{}, errors.Wrap(err, "access userInfo from storage failed")
	}
	ud, err := convertPoUserInfoToDto(userInfo)
	if err != nil {
		return &pb.HelloReply{}, errors.Wrap(err, "convert userinfo to dto failed")
	}
	return &pb.HelloReply{
		Name:    ud.name,
		Age:     uint32(ud.age),
		Address: ud.address,
	}, nil

}

func NewApp(d *sql.DB) *db.App {
	return &db.App{
		Db: d,
	}
}

func InitConf() {
	app, err := InitApp()
	if err != nil {
		log.Fatal(err)
	}
	db.GlobalApp = app
}

func InitPo() Po {
	dao, err := db.NewDao()
	if err != nil {
		log.Fatal(err)
	}
	return dao
}

func main() {
	InitConf()
	po = InitPo()
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("failed to listen tcp")
	}
	g := new(errgroup.Group)
	grpcServer := grpc.NewServer()
	g.Go(func() error {
		pb.RegisterGreeterServer(grpcServer, &server{})
		if err = grpcServer.Serve(lis); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		<-s
		grpcServer.GracefulStop()
		return errors.New("the grpc server will down soon because receive the interrupt signal")
	})
	if err = g.Wait(); err != nil {
		log.Printf("grpc server exit...")
	}
}
