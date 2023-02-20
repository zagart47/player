package player

import (
	context "context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"os"
	"player/internal/myerror"
	"player/internal/repo"
	"player/internal/storage"
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

var r = repo.NewSQLiteRepository(repo.DB)

func (p Player) AddSong(stream pb.PlayerService_AddSongServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("error")
	}
	file := storage.NewFile(md.Get("song")[0])
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			err := os.WriteFile(file.Path, file.Buffer.Bytes(), 0644)
			file.CheckDuration()
			if err != nil {
				return err
			}
			if err := r.CheckFileName(file.Name); err == nil {
				if err := r.Update(file.Name, file.Duration); err != nil {
					return err
				}
			} else if errors.Is(err, myerror.Err.FileNotFound) {
				if err := r.Create(file.Name, file.Duration); err != nil {
					return err
				}
			} else if err != nil {
				return err
			}
			return stream.SendAndClose(&pb.AddSongResponse{})
		}
		file.Buffer.Write(req.GetChunk())
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}
