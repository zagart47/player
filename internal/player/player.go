package player

import (
	context "context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"os"
	"player/internal/music"
	"player/internal/repo"
	"player/internal/storage"
	pb "player/proto"
	"time"
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
	server := grpc.NewServer()
	pb.RegisterPlayerServiceServer(server, &Player{})
	defer server.GracefulStop()
	log.Fatal(server.Serve(listener))
}

type Player struct {
	pb.UnsafePlayerServiceServer
}

var r = repo.NewSQLiteRepository(repo.DB)

func (p Player) Play(request *pb.PlayRequest, stream pb.PlayerService_PlayServer) error {
	/*md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return myerror.Err.MetaData
	}
	mdFileName := md.Get("filename")[0]
	file := storage.NewFile(mdFileName)
	if err := r.CheckFileName(file.Name); err != nil {
		return err
	}*/

	pl := repo.NewPlaylist()

	for e := pl.Front(); e != nil; e = e.Next() {
		duration := music.TrimDuration(e.Value)
		for {
			resp := &pb.PlayResponse{
				Info: fmt.Sprintf("Now playing: %v. Elapsed time: %d.", e.Value, duration),
			}
			time.Sleep(time.Second * 1)
			duration--
			if err := stream.Send(resp); err != nil {
				return status.Error(codes.Internal, err.Error())
			}
		}
	}
	return nil
}

func (p Player) Pause(ctx context.Context, request *pb.PauseRequest) (*pb.PauseResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p Player) Next(ctx context.Context, request *pb.NextRequest) (*pb.NextResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p Player) Prev(ctx context.Context, request *pb.PrevRequest) (*pb.PrevResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p Player) DeleteSong(ctx context.Context, request *pb.DeleteSongRequest) (*pb.DeleteSongResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p Player) UpdateSong(ctx context.Context, request *pb.UpdateSongRequest) (*pb.UpdateSongResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (p Player) AddSong(stream pb.PlayerService_AddSongServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("error")
	}
	file := storage.NewFile(md.Get("filename")[0])
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			err := os.WriteFile(file.Path, file.Buffer.Bytes(), 0644)
			file.CheckDuration()
			if err != nil {
				return err
			}
			if err := r.CheckFileName(file.Name); err == errors.New("file already have") {
				return errors.New("file already have, please use the Update method")
			} else if errors.Is(err, nil) {
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
