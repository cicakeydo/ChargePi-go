package grpc

import (
	"context"
	"time"

	"github.com/ChargePi/ChargePi-go/internal/evse"
	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	"github.com/ChargePi/ChargePi-go/pkg/proto/v1/grpc"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	grpc.UnimplementedEvseServer
	evseManager manager.Manager
}

func NewEvseService(manager manager.Manager) *Service {
	return &Service{
		evseManager: manager,
	}
}

func (s *Service) GetEVSEs(ctx context.Context, empty *empty.Empty) (*grpc.GetEvsesResponse, error) {
	response := &grpc.GetEvsesResponse{
		EVSEs: []*grpc.EVSE{},
	}

	for _, e := range s.evseManager.GetEVSEs() {
		evSe := toEvse(e)
		response.EVSEs = append(response.EVSEs, evSe)
	}

	return response, nil
}

func (s *Service) AddEVSE(ctx context.Context, request *grpc.AddEvseRequest) (*grpc.AddEvseResponse, error) {
	response := &grpc.AddEvseResponse{
		Status: grpc.ResponseStatus_Error,
	}

	// todo
	err := s.evseManager.AddEVSEFromSettings(ctx, evse.Settings{})
	if err != nil {
		return nil, err
	}

	response.Status = grpc.ResponseStatus_Success
	return response, nil
}

func (s *Service) GetEVSE(ctx context.Context, request *grpc.GetEvseRequest) (*grpc.GetEvseResponse, error) {
	res := &grpc.GetEvseResponse{}

	findEVSE, err := s.evseManager.GetEVSE(int(request.EvseId))
	if err != nil {
		return res, nil
	}

	res.EVSE = toEvse(findEVSE)
	return res, nil
}

func (s *Service) SetEVCC(ctx context.Context, request *grpc.SetEVCCRequest) (*grpc.SetEvccResponse, error) {
	// todo

	evse, err := s.evseManager.GetEVSE(int(request.EvseId))
	if err != nil {
		return nil, err
	}

	evse.SetEvcc(nil)

	return nil, nil
}

func (s *Service) SetPowerMeter(ctx context.Context, request *grpc.SetPowerMeterRequest) (*grpc.SetPowerMeterResponse, error) {
	// todo
	return nil, nil
}

func (s *Service) GetEVSEUsage(request *grpc.GetUsageForEVSERequest, server grpc.Evse_GetEVSEUsageServer) error {
	evseWithId, err := s.evseManager.GetEVSE(int(request.EvseId))
	if err != nil {
		return err
	}

	ctx := server.Context()

Loop:
	for {
		select {
		case <-ctx.Done():
			break Loop
		default:
			// Sample power meter
			samples, err := evseWithId.SamplePowerMeter([]types.Measurand{types.MeasurandEnergyActiveImportRegister})
			if err != nil {
				return err
			}

			// Convert to grpc samples
			var samplesToReturn []*grpc.Sample
			for _, sample := range samples {
				samplesToReturn = append(samplesToReturn, toSample(sample))
			}

			err = server.Send(&grpc.GetUsageForEVSEResponse{
				Samples: samplesToReturn,
			})
			if err != nil {
				return err
			}
			time.Sleep(time.Second * 10)
		}
	}

	return nil
}

func (s *Service) mustEmbedUnimplementedEvseServer() {
}

func toEvse(e evse.EVSE) *grpc.EVSE {
	return &grpc.EVSE{
		Id: int32(e.GetEvseId()),
		EVCC: &grpc.EVCC{
			Type:   e.GetEvcc().GetType(),
			Status: string(e.GetEvcc().GetState()),
		},
		PowerMeter: &grpc.PowerMeter{
			Type:    e.GetPowerMeter().GetType(),
			Enabled: false,
		},
		Status:  0,
		Session: &grpc.Session{},
	}
}

func toSample(sample types.SampledValue) *grpc.Sample {
	consumption := &grpc.Consumption{
		Unit: string(sample.Unit),
		// Value: sample.Value,
	}
	return &grpc.Sample{
		Consumption: consumption,
		Measureand:  string(sample.Measurand),
		Phase:       lo.ToPtr(string(sample.Phase)),
		Location:    lo.ToPtr(string(sample.Location)),
		Context:     lo.ToPtr(string(sample.Context)),
		Timestamp:   timestamppb.New(time.Now()),
	}
}
