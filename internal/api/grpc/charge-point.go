package grpc

import (
	"context"

	chargePoint "github.com/ChargePi/ChargePi-go/internal/chargepoint"
	settings "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ChargePi-go/pkg/hardware"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/display"
	"github.com/ChargePi/ChargePi-go/pkg/proto/v1/grpc"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/core"
)

type ChargePointService struct {
	grpc.UnimplementedChargePointServer
	point           chargePoint.ChargePoint
	settingsManager settings.Manager
}

func NewChargePointService(point chargePoint.ChargePoint, settingsManager settings.Manager) *ChargePointService {
	return &ChargePointService{
		point:           point,
		settingsManager: settingsManager,
	}
}

func (s *ChargePointService) SetDisplaySettings(ctx context.Context, request *grpc.SetDisplaySettingsRequest) (*grpc.SetDisplaySettingsResponse, error) {
	response := &grpc.SetDisplaySettingsResponse{
		Status: grpc.ResponseStatus_Error,
	}

	displaySettings := toDisplay(request.GetDisplay())

	newDisplay, err := display.NewDisplay(displaySettings)
	if err != nil {
		return response, nil
	}

	err = s.point.SetDisplay(newDisplay)
	if err != nil {
		return response, nil
	}

	// todo set the display settings in the manager

	response.Status = grpc.ResponseStatus_Success
	return response, nil
}

func (s *ChargePointService) GetDisplaySettings(ctx context.Context, empty *empty.Empty) (*grpc.GetDisplaySettingsResponse, error) {
	response := &grpc.GetDisplaySettingsResponse{}

	displaySettings, err := s.settingsManager.GetChargePointSettings()
	if err != nil {
		return nil, err
	}

	response.Display = &grpc.Display{
		Type:     displaySettings.Hardware.Display.Driver,
		Enabled:  displaySettings.Hardware.Display.IsEnabled,
		Language: &displaySettings.Hardware.Display.Language,
		// I2C:      i2cSettings,
	}

	return response, nil
}

func (s *ChargePointService) SetReaderSettings(ctx context.Context, request *grpc.SetReaderSettingsRequest) (*grpc.SetReaderSettingsResponse, error) {
	response := &grpc.SetReaderSettingsResponse{
		Status: grpc.ResponseStatus_Error,
	}

	return response, nil
}

func (s *ChargePointService) GetReaderSettings(ctx context.Context, empty *empty.Empty) (*grpc.GetReaderSettingsResponse, error) {
	response := &grpc.GetReaderSettingsResponse{}

	readerSettings, err := s.settingsManager.GetChargePointSettings()
	if err != nil {
		return nil, err
	}

	response.Reader = &grpc.TagReader{
		Type:    readerSettings.Hardware.TagReader.ReaderModel,
		Enabled: readerSettings.Hardware.TagReader.IsEnabled,
		// DeviceAddress: readerSettings.Device,
	}

	return response, nil
}

func (s *ChargePointService) SetIndicatorSettings(ctx context.Context, request *grpc.SetIndicatorSettingsRequest) (*grpc.SetIndicatorSettingsResponse, error) {
	response := &grpc.SetIndicatorSettingsResponse{
		Status: grpc.ResponseStatus_Error,
	}

	// todo
	err := s.point.SetIndicator(nil)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *ChargePointService) GetIndicatorSettings(ctx context.Context, empty *empty.Empty) (*grpc.GetIndicatorSettingsResponse, error) {
	response := &grpc.GetIndicatorSettingsResponse{}

	indicatorSettings, err := s.settingsManager.GetChargePointSettings()
	if err != nil {
		return nil, err
	}

	response.Indicator = &grpc.Indicator{
		Type:             indicatorSettings.Hardware.Indicator.Type,
		Enabled:          indicatorSettings.Hardware.Indicator.Enabled,
		IndicateCardRead: &indicatorSettings.Hardware.Indicator.IndicateCardRead,
		// Invert:           indicatorSettings.Invert,
	}

	return response, nil
}

func (s *ChargePointService) Restart(ctx context.Context, request *grpc.RestartRequest) (*empty.Empty, error) {
	err := s.point.Reset(request.Type)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *ChargePointService) ChangeConnectionDetails(ctx context.Context, request *grpc.ChangeConnectionDetailsRequest) (*grpc.ChangeConnectionDetailsResponse, error) {
	response := &grpc.ChangeConnectionDetailsResponse{
		Status: grpc.ResponseStatus_Error,
	}

	return response, nil
}

func (s *ChargePointService) ChangeChargePointDetails(ctx context.Context, request *grpc.ChangeChargePointDetailsRequest) (*grpc.ChangeChargePointDetailsResponse, error) {
	response := &grpc.ChangeChargePointDetailsResponse{
		Status: grpc.ResponseStatus_Error,
	}

	return response, nil
}

func (s *ChargePointService) GetOCPPVariables(ctx context.Context, e *empty.Empty) (*grpc.GetVariablesResponse, error) {
	response := &grpc.GetVariablesResponse{
		Variables: []*grpc.OcppVariable{},
	}

	configuration, err := s.settingsManager.GetConfiguration()
	if err != nil {
		return nil, err
	}

	for _, config := range configuration {
		response.Variables = append(response.Variables, toConfiguration(config))
	}

	return response, nil
}

func (s *ChargePointService) GetVersion(ctx context.Context, e *empty.Empty) (*grpc.GetVersionResponse, error) {
	return &grpc.GetVersionResponse{
		Version: s.point.GetVersion(),
	}, nil
}

func (s *ChargePointService) GetStatus(ctx context.Context, e *empty.Empty) (*grpc.GetStatusResponse, error) {
	return &grpc.GetStatusResponse{
		Connected: s.point.IsConnected(),
		// Status:    s.point.GetStatus(),
	}, nil
}

// todo migrate
func (s *ChargePointService) SetOCPPVariables(ctx context.Context, request *grpc.SetVariablesRequest) (*grpc.SetVariablesResponse, error) {
	response := &grpc.SetVariablesResponse{}

	for _, variable := range request.GetVariables() {
		status := "Failed"

		err := s.settingsManager.UpdateKey(ocpp_v16.Key(variable.Key), variable.Value)
		if err == nil {
			status = "Success"
		}

		response.Statuses = append(response.Statuses, status)
	}

	return response, nil
}

func (s *ChargePointService) GetOCPPVariable(ctx context.Context, request *grpc.GetVariableRequest) (*grpc.OcppVariable, error) {
	value, err := s.settingsManager.GetConfigurationValue(ocpp_v16.Key(request.GetKey()))
	if err != nil {
		return nil, err
	}

	return toConfiguration(core.ConfigurationKey{
		Key:      request.Key,
		Readonly: false,
		Value:    value,
	}), nil
}

func (s *ChargePointService) mustEmbedUnimplementedChargePointServer() {
}

func toConfiguration(key core.ConfigurationKey) *grpc.OcppVariable {
	return &grpc.OcppVariable{
		Key:      key.Key,
		Readonly: key.Readonly,
		Value:    key.Value,
	}
}

func toDisplay(displayReq *grpc.Display) display.Settings {
	return display.Settings{
		IsEnabled: false,
		Driver:    displayReq.Type,
		Language:  *displayReq.Language,
		// I2C:       nil,
	}
}

func toI2c(i2c hardware.I2C) *grpc.I2C {
	return &grpc.I2C{
		Address: i2c.Address,
		Bus:     int32(i2c.Bus),
	}
}
