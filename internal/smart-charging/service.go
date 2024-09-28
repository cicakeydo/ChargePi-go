package smartCharging

import (
	"errors"
	"strconv"
	"sync"

	"github.com/ChargePi/ChargePi-go/internal/evse/manager"
	manager2 "github.com/ChargePi/ChargePi-go/internal/pkg/configuration/manager"
	"github.com/ChargePi/ocppManager-go/ocpp_v16"
	"github.com/lorenzodonini/ocpp-go/ocpp1.6/types"
	log "github.com/sirupsen/logrus"
)

var (
	ErrProfileNotFound    = errors.New("profile not found")
	ErrCannotApplyProfile = errors.New("cannot apply profile")
)

type (
	Service interface {
		AddProfile(profile *types.ChargingProfile) error
		GetProfile(profileId int) (*types.ChargingProfile, error)
		GetProfiles() []types.ChargingProfile
		GetCompositeSchedule() []ScheduleInterval
		RemoveProfile(profileId int) error
	}

	Impl struct {
		numEvses           int
		maxCurrent         int
		evseManager        manager.Manager
		settingsManager    manager2.Manager
		logger             log.FieldLogger
		repository         ChargingScheduleRepository
		mu                 sync.Mutex
		compositeSchedules map[int][]ScheduleInterval
	}
)

func NewManager(maxCurrent int) *Impl {
	return &Impl{
		maxCurrent: maxCurrent,
		logger:     log.StandardLogger().WithField("component", "smart-charging"),
	}
}

func (m *Impl) AddProfile(profile *types.ChargingProfile) error {
	m.logger.WithField("profile", profile).Info("Adding profile")
	if profile == nil {
		return nil
	}

	if !m.canApplyProfile(profile) {
		return ErrCannotApplyProfile
	}

	profiles, err := m.GetProfiles()
	if err != nil {
		return err
	}

	// Check if profile exists
	for i, chargingProfile := range profiles {
		// If profile with id already exists, replace it
		if chargingProfile.ChargingProfileId == profile.ChargingProfileId {
			profiles[i] = *profile
			goto Store
		}
	}

	profiles = append(profiles, *profile)

Store:

	maxProfile := getProfileWithHighestStack(getValidProfiles(getProfilesWithPurpose(types.ChargingProfilePurposeChargePointMaxProfile, profiles)))
	txDefaultProfile := getProfileWithHighestStack(getValidProfiles(getProfilesWithPurpose(types.ChargingProfilePurposeTxDefaultProfile, profiles)))
	txProfile := getProfileWithHighestStack(getValidProfiles(getProfilesWithPurpose(types.ChargingProfilePurposeTxProfile, profiles)))

	compositeSchedule := CreateCompositeSchedule([]*types.ChargingProfile{txProfile, txDefaultProfile, maxProfile})
	m.validateCompositeSchedule(compositeSchedule)

	return m.repository.AddProfile(profile)
}

func (m *Impl) canApplyProfile(profile *types.ChargingProfile) bool {
	m.logger.WithField("profile", profile).Info("Checking if profile can be applied")

	// Get max stack level
	stackLevelStr, err := m.settingsManager.GetConfigurationValue(ocpp_v16.ChargeProfileMaxStackLevel)
	if err != nil {
		return false
	}

	stackLevel, _ := strconv.Atoi(*stackLevelStr)

	// Get max profiles
	maxProfilesStr, err := m.settingsManager.GetConfigurationValue(ocpp_v16.MaxChargingProfilesInstalled)
	if err != nil {
	}

	maxProfiles, _ := strconv.Atoi(*maxProfilesStr)

	// Get max charging schedule periods
	maxPeriodsStr, err := m.settingsManager.GetConfigurationValue(ocpp_v16.ChargingScheduleMaxPeriods)
	if err != nil {
	}

	maxPeriods, _ := strconv.Atoi(*maxPeriodsStr)

	// Check if the stack level is valid
	if profile.StackLevel > stackLevel {
		return false
	}

	// Check if MaxProfiles is reached
	profiles, err := m.GetProfiles()
	if err != nil {
		return false
	}

	if len(profiles)+1 >= maxProfiles {
		return false
	}

	// Check if the number of periods is valid
	if profile.ChargingSchedule != nil && len(profile.ChargingSchedule.ChargingSchedulePeriod) > maxPeriods {
		return false
	}

	return true
}

func (m *Impl) RemoveProfile(profileId int) error {
	m.logger.WithField("profile", profileId).Info("Removing profile")

	// Dont remove the profile if it is progress
	return m.repository.RemoveProfile(profileId)
}

func (m *Impl) GetProfile(profileId int) (*types.ChargingProfile, error) {
	m.logger.WithField("profile", profileId).Info("Getting a profile")

	return m.repository.GetProfile(profileId)
}

func (m *Impl) GetProfiles() ([]types.ChargingProfile, error) {
	m.logger.Info("Getting profiles")

	return m.repository.GetProfiles()
}

func (m *Impl) GetCompositeSchedule() []ScheduleInterval {
	m.logger.Info("Getting composite schedule")
	return nil
}

func (m *Impl) validateCompositeSchedule(newSchedule []ScheduleInterval) {

}
