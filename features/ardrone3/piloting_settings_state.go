package ardrone3

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/lock"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Piloting Settings state from product

// PilotingSettingsState ...
// TODO: Document this
type PilotingSettingsState interface {
	lock.ReadLockable
	// MotionDetectionEnabled returns a boolean indicating whether motion
	// detection is enabled. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	MotionDetectionEnabled() (bool, bool)
	// MaxAltitude returns the currently configured maximum altitude of the device
	// in meters. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	MaxAltitude() (float32, bool)
	// MaxAltitudeRangeMin returns the minimum altitude, in meters, that the
	// device's maximum altitude can be configured to. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MaxAltitudeRangeMin() (float32, bool)
	// MaxAltitudeRangeMax returns the maximum altitude, in meters, that the
	// device's maximum altitude can be configured to. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MaxAltitudeRangeMax() (float32, bool)
	// MaxDistance returns the currently configured maximum altitude of the device
	// in meters. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	MaxDistance() (float32, bool)
	// MaxDistanceRangeMin returns the minimum distance, in meters, that the
	// device's maximum distance can be configured to. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MaxDistanceRangeMin() (float32, bool)
	// MaxDistanceRangeMax returns the maximum distance, in meters, that the
	// device's maximum distance can be configured to. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MaxDistanceRangeMax() (float32, bool)
	// MaxTilt returns the configured maximum tilt (pitch and roll) of the device
	// in degrees. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	MaxTilt() (float32, bool)
	// MaxTiltRangeMin returns the minimum tilt (pitch and roll), in degrees, that
	// the device's maximum tilt can be configured to. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MaxTiltRangeMin() (float32, bool)
	// MaxTiltRangeMax is the maximum tilt (pitch and roll), in degrees, that the
	// device's maximum tilt can be configured to. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MaxTiltRangeMax() (float32, bool)
	// GeofencingEnabled indicates whether the drone should not fly beyond the
	// maximum configured distance (true) or may (false). A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	GeofencingEnabled() (bool, bool)
	// BankedTurningEnabled indicates whether banked turning is enabled. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	BankedTurningEnabled() (bool, bool)
}

type pilotingSettingsState struct {
	sync.RWMutex
	motionDetectionEnabled *bool
	// maxAltitude is the currently configured maximum altitude of the device in
	// meters.
	maxAltitude *float32
	// maxAltitudeRangeMin is the minimum altitude, in meters, that the device's
	// maximum altitude can be configured to.
	maxAltitudeRangeMin *float32
	// maxAltitudeRangeMax is the maximum altitude, in meters, that the device's
	// maximum altitude can be configured to.
	maxAltitudeRangeMax *float32
	// maxDistance is the configured maximum distance the device may fly from the
	// take off point in meters.
	maxDistance *float32
	// maxDistanceRangeMin is the minimum distance, in meters, that the device's
	// maximum distance can be configured to.
	maxDistanceRangeMin *float32
	// maxDistanceRangeMax is the maximum distance, in meters, that the device's
	// maximum distance can be configured to.
	maxDistanceRangeMax *float32
	// maxTilt is the configured maximum tilt (pitch and roll) of the device in
	// degrees.
	maxTilt *float32
	// maxTiltRangeMin is the minimum tilt (pitch and roll), in degrees, that the
	// device's maximum tilt can be configured to.
	maxTiltRangeMin *float32
	// maxTiltRangeMax is the maximum tilt (pitch and roll), in degrees, that the
	// device's maximum tilt can be configured to.
	maxTiltRangeMax *float32
	// geofencingEnabled indicates whether the drone should not fly beyond the
	// maximum configured distance
	geofencingEnabled *bool
	// bankedTurningEnabled indicates whether banked turning is enabled
	bankedTurningEnabled *bool
}

func newPilotingSettingsState() *pilotingSettingsState {
	return &pilotingSettingsState{}
}

func (p *pilotingSettingsState) ID() uint8 {
	return 6
}

func (p *pilotingSettingsState) Name() string {
	return "PilotingSettingsState"
}

func (p *pilotingSettingsState) D2CCommands(
	log *log.Entry,
) []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"MaxAltitudeChanged",
			[]interface{}{
				float32(0), // current,
				float32(0), // min,
				float32(0), // max,
			},
			p.maxAltitudeChanged,
			log,
		),
		arcommands.NewD2CCommand(
			1,
			"MaxTiltChanged",
			[]interface{}{
				float32(0), // current,
				float32(0), // min,
				float32(0), // max,
			},
			p.maxTiltChanged,
			log,
		),
		arcommands.NewD2CCommand(
			3,
			"MaxDistanceChanged",
			[]interface{}{
				float32(0), // current,
				float32(0), // min,
				float32(0), // max,
			},
			p.maxDistanceChanged,
			log,
		),
		arcommands.NewD2CCommand(
			4,
			"NoFlyOverMaxDistanceChanged",
			[]interface{}{
				uint8(0), // shouldNotFlyOver,
			},
			p.noFlyOverMaxDistanceChanged,
			log,
		),
		arcommands.NewD2CCommand(
			5,
			"AutonomousFlightMaxHorizontalSpeed",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxHorizontalSpeed,
			log,
		),
		arcommands.NewD2CCommand(
			6,
			"AutonomousFlightMaxVerticalSpeed",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxVerticalSpeed,
			log,
		),
		arcommands.NewD2CCommand(
			7,
			"AutonomousFlightMaxHorizontalAcceleration",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxHorizontalAcceleration,
			log,
		),
		arcommands.NewD2CCommand(
			8,
			"AutonomousFlightMaxVerticalAcceleration",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxVerticalAcceleration,
			log,
		),
		arcommands.NewD2CCommand(
			9,
			"AutonomousFlightMaxRotationSpeed",
			[]interface{}{
				float32(0), // value,
			},
			p.autonomousFlightMaxRotationSpeed,
			log,
		),
		arcommands.NewD2CCommand(
			10,
			"BankedTurnChanged",
			[]interface{}{
				uint8(0), // state,
			},
			p.bankedTurnChanged,
			log,
		),
		// arcommands.NewD2CCommand(
		// 	15,
		// 	"PitchModeChanged",
		// 	[]interface{}{
		// 		int32(0), // value,
		// 	},
		// 	p.pitchModeChanged,
		// ),
		arcommands.NewD2CCommand(
			16,
			"MotionDetection",
			[]interface{}{
				uint8(0), // enabled,
			},
			p.motionDetection,
			log,
		),
	}
}

// maxAltitudeChanged is invoked by the device when the maximum altitude setting
// is changed.
func (p *pilotingSettingsState) maxAltitudeChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.maxAltitude = ptr.ToFloat32(args[0].(float32))
	p.maxAltitudeRangeMin = ptr.ToFloat32(args[1].(float32))
	p.maxAltitudeRangeMax = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"maxAltitude", *p.maxAltitude,
	).WithField(
		"maxAltitudeRangeMin", *p.maxAltitudeRangeMin,
	).WithField(
		"maxAltitudeRangeMax", *p.maxAltitudeRangeMax,
	).Debug("max altitude changed")
	return nil
}

// maxTiltChanged is invoked by the device when the maximum tilt is changed.
func (p *pilotingSettingsState) maxTiltChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.maxTilt = ptr.ToFloat32(args[0].(float32))
	p.maxTiltRangeMin = ptr.ToFloat32(args[1].(float32))
	p.maxTiltRangeMax = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"maxTilt", *p.maxTilt,
	).WithField(
		"maxTiltRangeMin", *p.maxTiltRangeMin,
	).WithField(
		"maxTiltRangeMax", *p.maxTiltRangeMax,
	).Debug("max tilt changed")
	return nil
}

// maxDistanceChanged is invoked by the device when the max distance is changed.
func (p *pilotingSettingsState) maxDistanceChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.maxDistance = ptr.ToFloat32(args[0].(float32))
	p.maxDistanceRangeMin = ptr.ToFloat32(args[1].(float32))
	p.maxDistanceRangeMax = ptr.ToFloat32(args[2].(float32))
	log.WithField(
		"maxDistance", *p.maxDistance,
	).WithField(
		"maxDistanceRangeMin", *p.maxDistanceRangeMin,
	).WithField(
		"maxDistanceRangeMax", *p.maxDistanceRangeMax,
	).Debug("max distance changed")
	return nil
}

// noFlyOverMaxDistanceChanged is invoked by the device when geofencing is
// enabled or disabled.
func (p *pilotingSettingsState) noFlyOverMaxDistanceChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.geofencingEnabled = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"shouldNotFlyOverMaxDistance", args[0].(uint8),
	).Debug("geofencing enabled or disabled")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max horizontal speed
// Description: Autonomous flight max horizontal speed.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxHorizontalSpeed](#1-2-5).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxHorizontalSpeed(
	args []interface{},
	log *log.Entry,
) error {
	// value := args[0].(float32)
	//   maximum horizontal speed [m/s]
	log.Warn("command not implemented")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max vertical speed
// Description: Autonomous flight max vertical speed.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxVerticalSpeed](#1-2-6).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxVerticalSpeed(
	args []interface{},
	log *log.Entry,
) error {
	// value := args[0].(float32)
	//   maximum vertical speed [m/s]
	log.Warn("command not implemented")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max horizontal acceleration
// Description: Autonomous flight max horizontal acceleration.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxHorizontalAcceleration](#1-2-7).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxHorizontalAcceleration(
	args []interface{},
	log *log.Entry,
) error {
	// value := args[0].(float32)
	//   maximum horizontal acceleration [m/s2]
	log.Warn("command not implemented")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max vertical acceleration
// Description: Autonomous flight max vertical acceleration.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxVerticalAcceleration](#1-2-8).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxVerticalAcceleration(
	args []interface{},
	log *log.Entry,
) error {
	// value := args[0].(float32)
	//   maximum vertical acceleration [m/s2]
	log.Warn("command not implemented")
	return nil
}

// TODO: Implement this
// Title: Autonomous flight max rotation speed
// Description: Autonomous flight max rotation speed.
// Support: 0901:3.3.0;090c:3.3.0
// Triggered: by [SetAutonomousFlightMaxRotationSpeed](#1-2-9).
// Result:
func (p *pilotingSettingsState) autonomousFlightMaxRotationSpeed(
	args []interface{},
	log *log.Entry,
) error {
	// value := args[0].(float32)
	//   maximum yaw rotation speed [deg/s]
	log.Warn("command not implemented")
	return nil
}

// bankedTurnChanged is invoked by the device when banked turning is enabled or
// disabled
func (p *pilotingSettingsState) bankedTurnChanged(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.bankedTurningEnabled = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"state", args[0].(uint8),
	).Debug("banked turning enabled or disabled")
	return nil
}

// // TODO: Implement this
// // Title: Pitch mode
// // Description: Pitch mode.
// // Support: 090e
// // Triggered: by [SetPitchMode](#1-2-15).
// // Result:
// func (p *pilotingSettingsState) pitchModeChanged(args []interface{}) error {
// 	// value := args[0].(int32)
// 	//   The Pitch mode
// 	//   0: NORMAL: Positive pitch values will make the drone lower its nose.
// 	//      Negative pitch values will make the drone raise its nose.
// 	//   1: INVERTED: Pitch commands are inverted. Positive pitch values will
// 	//      make the drone raise its nose. Negative pitch values will make the
// 	//      drone lower its nose.
// 	log.Warn("command not implemented")
// 	return nil
// }

// motionDetection is invoked by the device to indicate whether motion
// detection is enabled.
func (p *pilotingSettingsState) motionDetection(
	args []interface{},
	log *log.Entry,
) error {
	p.Lock()
	defer p.Unlock()
	p.motionDetectionEnabled = ptr.ToBool(args[0].(uint8) == 1)
	log.WithField(
		"motionDetectionEnabled", *p.motionDetectionEnabled,
	).Debug("motion detection enabled changed")
	return nil
}

func (p *pilotingSettingsState) MotionDetectionEnabled() (bool, bool) {
	if p.motionDetectionEnabled == nil {
		return false, false
	}
	return *p.motionDetectionEnabled, true
}

func (p *pilotingSettingsState) MaxAltitude() (float32, bool) {
	if p.maxAltitude == nil {
		return 0, false
	}
	return *p.maxAltitude, true
}

func (p *pilotingSettingsState) MaxAltitudeRangeMin() (float32, bool) {
	if p.maxAltitudeRangeMin == nil {
		return 0, false
	}
	return *p.maxAltitudeRangeMin, true
}

func (p *pilotingSettingsState) MaxAltitudeRangeMax() (float32, bool) {
	if p.maxAltitudeRangeMax == nil {
		return 0, false
	}
	return *p.maxAltitudeRangeMax, true
}

func (p *pilotingSettingsState) MaxDistance() (float32, bool) {
	if p.maxDistance == nil {
		return 0, false
	}
	return *p.maxDistance, true
}

func (p *pilotingSettingsState) MaxDistanceRangeMin() (float32, bool) {
	if p.maxDistanceRangeMin == nil {
		return 0, false
	}
	return *p.maxDistanceRangeMin, true
}

func (p *pilotingSettingsState) MaxDistanceRangeMax() (float32, bool) {
	if p.maxDistanceRangeMax == nil {
		return 0, false
	}
	return *p.maxDistanceRangeMax, true
}

func (p *pilotingSettingsState) MaxTilt() (float32, bool) {
	if p.maxTilt == nil {
		return 0, false
	}
	return *p.maxTilt, true
}

func (p *pilotingSettingsState) MaxTiltRangeMin() (float32, bool) {
	if p.maxTiltRangeMin == nil {
		return 0, false
	}
	return *p.maxTiltRangeMin, true
}

func (p *pilotingSettingsState) MaxTiltRangeMax() (float32, bool) {
	if p.maxTiltRangeMax == nil {
		return 0, false
	}
	return *p.maxTiltRangeMax, true
}

func (p *pilotingSettingsState) GeofencingEnabled() (bool, bool) {
	if p.geofencingEnabled == nil {
		return false, false
	}
	return *p.geofencingEnabled, true
}

func (p *pilotingSettingsState) BankedTurningEnabled() (bool, bool) {
	if p.bankedTurningEnabled == nil {
		return false, false
	}
	return *p.bankedTurningEnabled, true
}
