package player

import (
	context "context"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "player/proto"
)

func Run() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer func(l net.Listener) {
		if err := l.Close(); err != nil {
			log.Fatal(err)
		}
	}(listener)
	newServer := grpc.NewServer()
	pb.RegisterPlayerServiceServer(newServer, &Player{})
	defer newServer.GracefulStop()
}

type Player struct {
	pb.UnsafePlayerServiceServer
}

func (p Player) Play(ctx context.Context, song *pb.Song) (*pb.Song, error) {
	//TODO implement me
	panic("implement me")
}

func (p Player) Pause(ctx context.Context, song *pb.Song) (*pb.Song, error) {
	//TODO implement me
	panic("implement me")
}

func (p Player) Next(ctx context.Context, song *pb.Song) (*pb.Song, error) {
	//TODO implement me
	panic("implement me")
}

func (p Player) Prev(ctx context.Context, song *pb.Song) (*pb.Song, error) {
	//TODO implement me
	panic("implement me")
}

func (p Player) AddSong(ctx context.Context, song *pb.Song) (*pb.Song, error) {
	//TODO implement me
	panic("implement me")
}
