package common

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/krancour/go-parrot/protocols/arcommands"
	"github.com/krancour/go-parrot/ptr"
)

// Common state from product

// CommonState ...
// TODO: Document this
type CommonState interface {
	// RLock blocks until a read lock is obtained. This permits callers to procede
	// with querying any or all attributes of the common state without worry
	// that some attributes will be overwritten as others are read. i.e. It
	// permits the possibility of taking an atomic snapshop of common state.
	// Note that use of this function is not obligatory for applications that do
	// not require such guarantees. Callers MUST call RUnlock() or else piloting
	// state will never resume updating.
	RLock()
	// RUnlock releases a read lock on the common state. See RLock().
	RUnlock()
	// RSSI returns the relative signal stength between the client and the device
	// in dbm. A boolean value is also returned, indicating whether the first
	// value was reported by the device (true) or a default value (false). This
	// permits callers to distinguish real zero values from default zero values.
	RSSI() (int16, bool)
	// MassStorageID returns the mass storage ID. A boolean value is also
	// returned, indicating whether the first value was reported by the device
	// (true) or a default value (false). This permits callers to distinguish real
	// zero values from default zero values.
	MassStorageID() (uint8, bool)
	// PhotoCount returns the number of photos in mass storage. A boolean value is
	// also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	PhotoCount() (uint16, bool)
	// VideoCount returns the numnber of videos in mass storage. A boolean value
	// is also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	VideoCount() (uint16, bool)
	// PudCount returns the number of puds in mass storage. A boolean value is
	// also returned, indicating whether the first value was reported by the
	// device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	PudCount() (uint16, bool)
	// CrashLogCount returns the number of crash logs in mass storage. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	CrashLogCount() (uint16, bool)
	// RawPhotoCount returns the number of raw photos in mass stroage. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	RawPhotoCount() (uint16, bool)
	// CurrentRunMassStorageID returns the mass storage ID for the current run. A
	// boolean value is also returned, indicating whether the first value was
	// reported by the device (true) or a default value (false). This permits
	// callers to distinguish real zero values from default zero values.
	CurrentRunMassStorageID() (uint8, bool)
	// CurrentRunPhotoCount returns the number of photos in mass storage related
	// to the current run. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	CurrentRunPhotoCount() (uint16, bool)
	// CurrentRunVideoCount returns the number of photos in mass storage related
	// to the current run. A boolean value is also returned, indicating whether
	// the first value was reported by the device (true) or a default value
	// (false). This permits callers to distinguish real zero values from default
	// zero values.
	CurrentRunVideoCount() (uint16, bool)
	// CurrentRunRawPhotoCount returns the number of raw photos in mass storage
	// related to the current run. A boolean value is also returned, indicating
	// whether the first value was reported by the device (true) or a default
	// value (false). This permits callers to distinguish real zero values from
	// default zero values.
	CurrentRunRawPhotoCount() (uint16, bool)
	// BatteryPercent returns the percentage of battery life remaining. A boolean
	// value is also returned, indicating whether the first value was reported by
	// the device (true) or a default value (false). This permits callers to
	// distinguish real zero values from default zero values.
	BatteryPercent() (uint8, bool)
}

type commonState struct {
	// TODO: Is this right? I thought RSSI is a relative measure, while dbm
	// would seem to indicate an absolute measure.
	// rssi is the relative signal stength between the client and the device
	// in dbm
	rssi                    *int16
	massStorageID           *uint8
	photoCount              *uint16
	videoCount              *uint16
	pudCount                *uint16 // TODO: What is a pud?
	crashLogCount           *uint16
	rawPhotoCount           *uint16
	currentRunMassStorageID *uint8
	currentRunPhotoCount    *uint16
	currentRunVideoCount    *uint16
	currentRunRawPhotoCount *uint16
	batteryPercent          *uint8
	lock                    sync.RWMutex
}

func (c *commonState) ID() uint8 {
	return 5
}

func (c *commonState) Name() string {
	return "CommonState"
}

func (c *commonState) D2CCommands() []arcommands.D2CCommand {
	return []arcommands.D2CCommand{
		arcommands.NewD2CCommand(
			0,
			"AllStatesChanged",
			[]interface{}{},
			c.allStatesChanged,
		),
		arcommands.NewD2CCommand(
			1,
			"BatteryStateChanged",
			[]interface{}{
				uint8(0), // percent,
			},
			c.batteryStateChanged,
		),
		arcommands.NewD2CCommand(
			2,
			"MassStorageStateListChanged",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				string(0), // name,
			},
			c.massStorageStateListChanged,
		),
		arcommands.NewD2CCommand(
			3,
			"MassStorageInfoStateListChanged",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				uint32(0), // size,
				uint32(0), // used_size,
				uint8(0),  // plugged,
				uint8(0),  // full,
				uint8(0),  // internal,
			},
			c.massStorageInfoStateListChanged,
		),
		arcommands.NewD2CCommand(
			4,
			"CurrentDateChanged",
			[]interface{}{
				string(0), // date,
			},
			c.currentDateChanged,
		),
		arcommands.NewD2CCommand(
			5,
			"CurrentTimeChanged",
			[]interface{}{
				string(0), // time,
			},
			c.currentTimeChanged,
		),
		// arcommands.NewD2CCommand(
		// 	6,
		// 	"MassStorageInfoRemainingListChanged",
		// 	[]interface{}{
		// 		uint32(0), // free_space,
		// 		uint16(0), // rec_time,
		// 		uint32(0), // photo_remaining,
		// 	},
		// 	c.massStorageInfoRemainingListChanged,
		// ),
		arcommands.NewD2CCommand(
			7,
			"WifiSignalChanged",
			[]interface{}{
				int16(0), // rssi,
			},
			c.wifiSignalChanged,
		),
		arcommands.NewD2CCommand(
			8,
			"SensorsStatesListChanged",
			[]interface{}{
				int32(0), // sensorName,
				uint8(0), // sensorState,
			},
			c.sensorsStatesListChanged,
		),
		// arcommands.NewD2CCommand(
		// 	9,
		// 	"ProductModel",
		// 	[]interface{}{
		// 		int32(0), // model,
		// 	},
		// 	c.productModel,
		// ),
		// arcommands.NewD2CCommand(
		// 	10,
		// 	"CountryListKnown",
		// 	[]interface{}{
		// 		uint8(0),  // listFlags,
		// 		string(0), // countryCodes,
		// 	},
		// 	c.countryListKnown,
		// ),
		arcommands.NewD2CCommand(
			11,
			"DeprecatedMassStorageContentChanged",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				uint16(0), // nbPhotos,
				uint16(0), // nbVideos,
				uint16(0), // nbPuds,
				uint16(0), // nbCrashLogs,
			},
			c.deprecatedMassStorageContentChanged,
		),
		arcommands.NewD2CCommand(
			12,
			"MassStorageContent",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				uint16(0), // nbPhotos,
				uint16(0), // nbVideos,
				uint16(0), // nbPuds,
				uint16(0), // nbCrashLogs,
				uint16(0), // nbRawPhotos,
			},
			c.massStorageContent,
		),
		arcommands.NewD2CCommand(
			13,
			"MassStorageContentForCurrentRun",
			[]interface{}{
				uint8(0),  // mass_storage_id,
				uint16(0), // nbPhotos,
				uint16(0), // nbVideos,
				uint16(0), // nbRawPhotos,
			},
			c.massStorageContentForCurrentRun,
		),
		arcommands.NewD2CCommand(
			14,
			"VideoRecordingTimestamp",
			[]interface{}{
				uint64(0), // startTimestamp,
				uint64(0), // stopTimestamp,
			},
			c.videoRecordingTimestamp,
		),
	}
}

// TODO: Implement this
// Title: All states have been sent
// Description: All states have been sent.\n\n **Please note that you should not
//   care about this event if you are using the libARController API as this
//   library is handling the connection process for you.**
// Support: drones
// Triggered: when all states values have been sent.
// Result:
func (c *commonState) allStatesChanged(args []interface{}) error {
	log.Info("common.allStatesChanged() called")
	return nil
}

// TODO: Implement this
// Title: Battery state
// Description: Battery state.
// Support: drones
// Triggered: when the battery level changes.
// Result:
func (c *commonState) batteryStateChanged(args []interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.batteryPercent = ptr.ToUint8(args[0].(uint8))
	log.WithField(
		"batteryPercent", *c.batteryPercent,
	).Debug("battery state changed")
	return nil
}

// TODO: Implement this
// Title: Mass storage state list
// Description: Mass storage state list.
// Support: drones
// Triggered: when a mass storage is inserted or ejected.
// Result:
func (c *commonState) massStorageStateListChanged(args []interface{}) error {
	// mass_storage_id := args[0].(uint8)
	//   Mass storage id (unique)
	// name := args[1].(string)
	//   Mass storage name
	log.Info("common.massStorageStateListChanged() called")
	return nil
}

// TODO: Implement this
// Title: Mass storage info state list
// Description: Mass storage info state list.
// Support: drones
// Triggered: when a mass storage info changes.
// Result:
func (c *commonState) massStorageInfoStateListChanged(args []interface{}) error {
	// mass_storage_id := args[0].(uint8)
	//   Mass storage state id (unique)
	// size := args[1].(uint32)
	//   Mass storage size in MBytes
	// used_size := args[2].(uint32)
	//   Mass storage used size in MBytes
	// plugged := args[3].(uint8)
	//   Mass storage plugged (1 if mass storage is plugged, otherwise 0)
	// full := args[4].(uint8)
	//   Mass storage full information state (1 if mass storage full, 0
	//   otherwise).
	// internal := args[5].(uint8)
	//   Mass storage internal type state (1 if mass storage is internal, 0
	//   otherwise)
	log.Info("common.massStorageInfoStateListChanged() called")
	return nil
}

// TODO: Implement this
// Title: Date changed
// Description: Date changed.\n Corresponds to the latest date set on the drone.
//   \n\n **Please note that you should not care about this event if you are
//   using the libARController API as this library is handling the connection
//   process for you.**
// Support: drones
// Triggered: by [SetDate](#0-4-1).
// Result:
func (c *commonState) currentDateChanged(args []interface{}) error {
	// date := args[0].(string)
	//   Date with ISO-8601 format
	log.Info("common.currentDateChanged() called")
	return nil
}

// TODO: Implement this
// Title: Time changed
// Description: Time changed.\n Corresponds to the latest time set on the drone.
//   \n\n **Please note that you should not care about this event if you are
//   using the libARController API as this library is handling the connection
//   process for you.**
// Support: drones
// Triggered: by [SetTime](#0-4-2).
// Result:
func (c *commonState) currentTimeChanged(args []interface{}) error {
	// time := args[0].(string)
	//   Time with ISO-8601 format
	log.Info("common.currentTimeChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Mass storage remaining data list
// // Description: Mass storage remaining data list.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (c *commonState) massStorageInfoRemainingListChanged(
// 	args []interface{},
// ) error {
// 	// free_space := args[0].(uint32)
// 	//   Mass storage free space in MBytes
// 	// rec_time := args[1].(uint16)
// 	//   Mass storage record time reamining in minute
// 	// photo_remaining := args[2].(uint32)
// 	//   Mass storage photo remaining
// 	log.Info("common.massStorageInfoRemainingListChanged() called")
// 	return nil
// }

// wifiSignalChanged is invoked when the device reports relative wifi signal
// strength at regular intervals.
func (c *commonState) wifiSignalChanged(args []interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.rssi = ptr.ToInt16(args[0].(int16))
	log.WithField(
		"rssi", *c.rssi,
	).Debug("common state wifi signal strength updated")
	return nil
}

// TODO: Implement this
// Title: Sensors state list
// Description: Sensors state list.
// Support: 0901:2.0.3;0902;0905;0906;0907;0909;090a;090c;090e
// Triggered: at connection and when a sensor state changes.
// Result:
func (c *commonState) sensorsStatesListChanged(args []interface{}) error {
	// sensorName := args[0].(int32)
	//   Sensor name
	//   0: IMU: Inertial Measurement Unit sensor
	//   1: barometer: Barometer sensor
	//   2: ultrasound: Ultrasonic sensor
	//   3: GPS: GPS sensor
	//   4: magnetometer: Magnetometer sensor
	//   5: vertical_camera: Vertical Camera sensor
	// sensorState := args[1].(uint8)
	//   Sensor state (1 if the sensor is OK, 0 if the sensor is NOT OK)
	log.Info("common.sensorsStatesListChanged() called")
	return nil
}

// // TODO: Implement this
// // Title: Product sub-model
// // Description: Product sub-model.\n This can be used to customize the UI
// //   depending on the product.
// // Support: 0905;0906;0907;0909
// // Triggered: at connection.
// // Result:
// func (c *commonState) productModel(args []interface{}) error {
// 	// model := args[0].(int32)
// 	//   The Model of the product.
// 	//   0: RS_TRAVIS: Travis (RS taxi) model.
// 	//   1: RS_MARS: Mars (RS space) model
// 	//   2: RS_SWAT: SWAT (RS SWAT) model
// 	//   3: RS_MCLANE: Mc Lane (RS police) model
// 	//   4: RS_BLAZE: Blaze (RS fire) model
// 	//   5: RS_ORAK: Orak (RS carbon hydrofoil) model
// 	//   6: RS_NEWZ: New Z (RS wooden hydrofoil) model
// 	//   7: JS_MARSHALL: Marshall (JS fire) model
// 	//   8: JS_DIESEL: Diesel (JS SWAT) model
// 	//   9: JS_BUZZ: Buzz (JS space) model
// 	//   10: JS_MAX: Max (JS F1) model
// 	//   11: JS_JETT: Jett (JS flames) model
// 	//   12: JS_TUKTUK: Tuk-Tuk (JS taxi) model
// 	//   13: SW_BLACK: Swing black model
// 	//   14: SW_WHITE: Swing white model
// 	log.Info("common.productModel() called")
// 	return nil
// }

// // TODO: Implement this
// // Title: Country list
// // Description: List of countries known by the drone.
// // Support:
// // Triggered:
// // Result:
// // WARNING: Deprecated
// func (c *commonState) countryListKnown(args []interface{}) error {
// 	// listFlags := args[0].(uint8)
// 	//   List entry attribute Bitfield. 0x01: First: indicate it&#39;s the first
// 	//   element of the list. 0x02: Last: indicate it&#39;s the last element of
// 	//   the list. 0x04: Empty: indicate the list is empty (implies First/Last).
// 	//   All other arguments should be ignored.
// 	// countryCodes := args[1].(string)
// 	//   Following of country code with ISO 3166 format, separated by &#34;;&#34;.
// 	//   Be careful of the command size allowed by the network used. If necessary,
// 	//   split the list in several commands.
// 	log.Info("common.countryListKnown() called")
// 	return nil
// }

// deprecatedMassStorageContentChanged is deprecated in favor of
// massStorageContent, but since we can still see this command being invoked,
// we'll implement the command to avoid a warning, but the implementation will
// remain a no-op unless / until such time that it becomes clear that older
// versions of the firmware might require us to support both commands.
func (c *commonState) deprecatedMassStorageContentChanged(
	args []interface{},
) error {
	// mass_storage_id := args[0].(uint8)
	//   Mass storage id (unique)
	// nbPhotos := args[1].(uint16)
	//   Number of photos (does not include raw photos)
	// nbVideos := args[2].(uint16)
	//   Number of videos
	// nbPuds := args[3].(uint16)
	//   Number of puds
	// nbCrashLogs := args[4].(uint16)
	//   Number of crash logs
	log.Debug("mass storage content changed-- this is a no-op")
	return nil
}

// massStorageContent is invoked when the device reports that the content of the
// mass storage has changed.
func (c *commonState) massStorageContent(args []interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.massStorageID = ptr.ToUint8(args[0].(uint8))
	c.photoCount = ptr.ToUint16(args[1].(uint16))
	c.videoCount = ptr.ToUint16(args[2].(uint16))
	c.pudCount = ptr.ToUint16(args[3].(uint16))
	c.crashLogCount = ptr.ToUint16(args[4].(uint16))
	c.rawPhotoCount = ptr.ToUint16(args[5].(uint16))
	log.WithField(
		"massStorageID", *c.massStorageID,
	).WithField(
		"photoCount", *c.photoCount,
	).WithField(
		"videoCount", *c.videoCount,
	).WithField(
		"pudCount", *c.pudCount,
	).WithField(
		"crashLogCount", *c.crashLogCount,
	).WithField(
		"rawPhotoCount", *c.rawPhotoCount,
	).Debug("mass storage content changed")
	return nil
}

// massStorageContentForCurrentRun is invoked when the device reports that the
// content of the mass storage has changed and the content is related to the
// current run.s
func (c *commonState) massStorageContentForCurrentRun(
	args []interface{},
) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.currentRunMassStorageID = ptr.ToUint8(args[0].(uint8))
	c.currentRunPhotoCount = ptr.ToUint16(args[1].(uint16))
	c.currentRunVideoCount = ptr.ToUint16(args[2].(uint16))
	c.currentRunRawPhotoCount = ptr.ToUint16(args[3].(uint16))
	log.WithField(
		"currentRunMassStorageID", *c.currentRunMassStorageID,
	).WithField(
		"currentRunPhotoCount", *c.currentRunPhotoCount,
	).WithField(
		"currentRunVideoCount", *c.currentRunVideoCount,
	).WithField(
		"currentRunRawPhotoCount", *c.currentRunRawPhotoCount,
	).Debug("mass storage content for current run changed")
	return nil
}

// videoRecordingTimestamp is invoked by the device on video recording start and
// video recording stop or after that the date/time of the drone changed.
func (c *commonState) videoRecordingTimestamp(args []interface{}) error {
	// startTimestamp := args[0].(uint64)
	//   Timestamp in milliseconds since 00:00:00 UTC on 1 January 1970.
	// stopTimestamp := args[1].(uint64)
	//   Timestamp in milliseconds since 00:00:00 UTC on 1 January 1970. 0 mean
	//   that video is still recording.
	log.Info("common.videoRecordingTimestamp() called")
	return nil
}

func (c *commonState) RLock() {
	c.lock.RLock()
}

func (c *commonState) RUnlock() {
	c.lock.RUnlock()
}

func (c *commonState) RSSI() (int16, bool) {
	if c.rssi == nil {
		return 0, false
	}
	return *c.rssi, true
}

func (c *commonState) MassStorageID() (uint8, bool) {
	if c.massStorageID == nil {
		return 0, false
	}
	return *c.massStorageID, true
}

func (c *commonState) PhotoCount() (uint16, bool) {
	if c.photoCount == nil {
		return 0, false
	}
	return *c.photoCount, true
}

func (c *commonState) VideoCount() (uint16, bool) {
	if c.videoCount == nil {
		return 0, false
	}
	return *c.videoCount, true
}

func (c *commonState) PudCount() (uint16, bool) {
	if c.pudCount == nil {
		return 0, false
	}
	return *c.pudCount, true
}

func (c *commonState) CrashLogCount() (uint16, bool) {
	if c.crashLogCount == nil {
		return 0, false
	}
	return *c.crashLogCount, true
}

func (c *commonState) RawPhotoCount() (uint16, bool) {
	if c.rawPhotoCount == nil {
		return 0, false
	}
	return *c.rawPhotoCount, true
}

func (c *commonState) CurrentRunMassStorageID() (uint8, bool) {
	if c.currentRunMassStorageID == nil {
		return 0, false
	}
	return *c.currentRunMassStorageID, true
}

func (c *commonState) CurrentRunPhotoCount() (uint16, bool) {
	if c.currentRunPhotoCount == nil {
		return 0, false
	}
	return *c.currentRunPhotoCount, true
}

func (c *commonState) CurrentRunVideoCount() (uint16, bool) {
	if c.currentRunVideoCount == nil {
		return 0, false
	}
	return *c.currentRunVideoCount, true
}

func (c *commonState) CurrentRunRawPhotoCount() (uint16, bool) {
	if c.currentRunRawPhotoCount == nil {
		return 0, false
	}
	return *c.currentRunRawPhotoCount, true
}

func (c *commonState) BatteryPercent() (uint8, bool) {
	if c.batteryPercent == nil {
		return 0, false
	}
	return *c.batteryPercent, true
}
