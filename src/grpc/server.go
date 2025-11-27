// Package grpc
package grpc

import (
	"context"
	pb "metar-service/src/interfaces/grpc"
	"metar-service/src/interfaces/metar"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"half-nothing.cn/service-core/interfaces/logger"
)

type MetarServer struct {
	pb.UnimplementedMetarServer
	logger       logger.Interface
	metarManager metar.ManagerInterface
	tafManager   metar.ManagerInterface
}

func NewMetarServer(
	lg logger.Interface,
	metarManager metar.ManagerInterface,
	tafManager metar.ManagerInterface,
) *MetarServer {
	return &MetarServer{
		logger:       logger.NewLoggerAdapter(lg, "GrpcMetarService"),
		metarManager: metarManager,
		tafManager:   tafManager,
	}
}

func (m MetarServer) GetMetar(_ context.Context, in *pb.MetarQuery) (*pb.MetarReply, error) {
	if len(in.Icao) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid ICAO")
	}
	return &pb.MetarReply{Metar: m.metarManager.BatchQuery(in.Icao)}, nil
}

func (m MetarServer) GetTaf(_ context.Context, in *pb.TafQuery) (*pb.TafReply, error) {
	if len(in.Icao) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid ICAO")
	}
	return &pb.TafReply{Taf: m.tafManager.BatchQuery(in.Icao)}, nil
}
