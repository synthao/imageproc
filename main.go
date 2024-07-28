package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	pb "github.com/synthao/imageproc/gen/go/imgproc"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
)

type server struct {
	pb.UnimplementedImageProcessingServiceServer
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env", err)
	}

	logger, err := newLogger(os.Getenv("LOG_LEVEL"))
	if err != nil {
		log.Fatal("failed to create logger", err)
	}

	s := grpc.NewServer()
	pb.RegisterImageProcessingServiceServer(s, &server{})

	lis, err := net.Listen("tcp", os.Getenv("GRPC_SERVER_ADDRESS"))
	if err != nil {
		logger.Error(err.Error())
	}

	go func() {
		http.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

		http.ListenAndServe(net.JoinHostPort(os.Getenv("HTTP_SERVER_HOST"), os.Getenv("HTTP_SERVER_PORT")), nil)
	}()

	log.Printf("imgproc GRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		logger.Error("failed to serve", zap.Error(err))
	}

	logger.Info("response from imgproc") // TODO show images
}

func (s *server) ProcessImage(ctx context.Context, req *pb.ProcessImageRequest) (*pb.ProcessImageResponse, error) {
	fmt.Println(">> REQ:", req.W, req.H) // TODO
	// Реализация логики обработки изображений
	return &pb.ProcessImageResponse{
		Small:  "small_" + req.Path,
		Medium: "medium_" + req.Path,
		Large:  "large_" + req.Path,
	}, nil
}

func newLogger(lvl string) (*zap.Logger, error) {
	atomicLogLevel, err := zap.ParseAtomicLevel(lvl)
	if err != nil {
		return nil, err
	}

	atom := zap.NewAtomicLevelAt(atomicLogLevel.Level())
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		),
		zap.WithCaller(true),
		zap.AddStacktrace(zap.ErrorLevel),
	), nil
}
