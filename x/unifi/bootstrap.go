package unifi

import (
	"context"
	"net/http"

	"golang.org/x/xerrors"
)

type Bootstrap struct {
	AuthUserID string  `json:"authUserId"`
	AccessKey  string  `json:"accessKey"`
	Cameras    Cameras `json:"cameras"`
	Users      []struct {
		Permissions         []string      `json:"permissions"`
		LastLoginIP         interface{}   `json:"lastLoginIp"`
		LastLoginTime       interface{}   `json:"lastLoginTime"`
		IsOwner             bool          `json:"isOwner"`
		LocalUsername       string        `json:"localUsername"`
		EnableNotifications bool          `json:"enableNotifications"`
		SyncSso             bool          `json:"syncSso"`
		Settings            interface{}   `json:"settings"`
		Groups              []string      `json:"groups"`
		CloudAccount        interface{}   `json:"cloudAccount"`
		AlertRules          []interface{} `json:"alertRules"`
		ID                  string        `json:"id"`
		HasAcceptedInvite   bool          `json:"hasAcceptedInvite"`
		Role                string        `json:"role"`
		AllPermissions      []string      `json:"allPermissions"`
		ModelKey            string        `json:"modelKey"`
	} `json:"users"`
	Groups []struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
		Type        string   `json:"type"`
		IsDefault   bool     `json:"isDefault"`
		ID          string   `json:"id"`
		ModelKey    string   `json:"modelKey"`
	} `json:"groups"`
	Liveviews []struct {
		Name      string `json:"name"`
		IsDefault bool   `json:"isDefault"`
		IsGlobal  bool   `json:"isGlobal"`
		Layout    int    `json:"layout"`
		Slots     []struct {
			Cameras       []string `json:"cameras"`
			CycleMode     string   `json:"cycleMode"`
			CycleInterval int      `json:"cycleInterval"`
		} `json:"slots"`
		Owner    string `json:"owner"`
		ID       string `json:"id"`
		ModelKey string `json:"modelKey"`
	} `json:"liveviews"`
	Nvr struct {
		Mac                     string      `json:"mac"`
		Host                    string      `json:"host"`
		Name                    string      `json:"name"`
		CanAutoUpdate           bool        `json:"canAutoUpdate"`
		IsStatsGatheringEnabled bool        `json:"isStatsGatheringEnabled"`
		Timezone                string      `json:"timezone"`
		Version                 string      `json:"version"`
		FirmwareVersion         string      `json:"firmwareVersion"`
		UIVersion               interface{} `json:"uiVersion"`
		HardwarePlatform        string      `json:"hardwarePlatform"`
		Ports                   struct {
			Ump             int `json:"ump"`
			HTTP            int `json:"http"`
			HTTPS           int `json:"https"`
			Rtsp            int `json:"rtsp"`
			Rtmp            int `json:"rtmp"`
			DevicesWss      int `json:"devicesWss"`
			CameraHTTPS     int `json:"cameraHttps"`
			CameraTCP       int `json:"cameraTcp"`
			LiveWs          int `json:"liveWs"`
			LiveWss         int `json:"liveWss"`
			TCPStreams      int `json:"tcpStreams"`
			EmsCLI          int `json:"emsCLI"`
			EmsLiveFLV      int `json:"emsLiveFLV"`
			CameraEvents    int `json:"cameraEvents"`
			DiscoveryClient int `json:"discoveryClient"`
		} `json:"ports"`
		SetupCode                    interface{} `json:"setupCode"`
		Uptime                       int         `json:"uptime"`
		LastSeen                     int64       `json:"lastSeen"`
		IsUpdating                   bool        `json:"isUpdating"`
		LastUpdateAt                 int64       `json:"lastUpdateAt"`
		IsConnectedToCloud           bool        `json:"isConnectedToCloud"`
		CloudConnectionError         interface{} `json:"cloudConnectionError"`
		IsStation                    bool        `json:"isStation"`
		EnableAutomaticBackups       bool        `json:"enableAutomaticBackups"`
		EnableStatsReporting         bool        `json:"enableStatsReporting"`
		IsSSHEnabled                 bool        `json:"isSshEnabled"`
		ErrorCode                    interface{} `json:"errorCode"`
		ReleaseChannel               string      `json:"releaseChannel"`
		AvailableUpdate              interface{} `json:"availableUpdate"`
		Hosts                        []string    `json:"hosts"`
		HardwareID                   string      `json:"hardwareId"`
		HardwareRevision             string      `json:"hardwareRevision"`
		HostType                     int         `json:"hostType"`
		HostShortname                string      `json:"hostShortname"`
		IsHardware                   bool        `json:"isHardware"`
		TimeFormat                   string      `json:"timeFormat"`
		RecordingRetentionDurationMs interface{} `json:"recordingRetentionDurationMs"`
		EnableCrashReporting         bool        `json:"enableCrashReporting"`
		DisableAudio                 interface{} `json:"disableAudio"`
		WifiSettings                 struct {
			UseThirdPartyWifi bool        `json:"useThirdPartyWifi"`
			Ssid              interface{} `json:"ssid"`
			Password          interface{} `json:"password"`
		} `json:"wifiSettings"`
		LocationSettings struct {
			IsAway              bool        `json:"isAway"`
			IsGeofencingEnabled bool        `json:"isGeofencingEnabled"`
			Latitude            interface{} `json:"latitude"`
			Longitude           interface{} `json:"longitude"`
			Radius              interface{} `json:"radius"`
		} `json:"locationSettings"`
		FeatureFlags struct {
			Beta bool `json:"beta"`
			Dev  bool `json:"dev"`
		} `json:"featureFlags"`
		StorageInfo struct {
			TotalSize          int64 `json:"totalSize"`
			TotalSpaceUsed     int64 `json:"totalSpaceUsed"`
			StorageUtilization []struct {
				Type      string `json:"type"`
				SpaceUsed int64  `json:"spaceUsed"`
			} `json:"storageUtilization"`
			HardDrives []struct {
				Status      string `json:"status"`
				Name        string `json:"name"`
				Serial      string `json:"serial"`
				Firmware    string `json:"firmware"`
				Size        int64  `json:"size"`
				RPM         int    `json:"RPM"`
				AtaVersion  string `json:"ataVersion"`
				SataVersion string `json:"sataVersion"`
				Health      string `json:"health"`
			} `json:"hardDrives"`
		} `json:"storageInfo"`
		DoorbellSettings struct {
			DefaultMessageText           string   `json:"defaultMessageText"`
			DefaultMessageResetTimeoutMs int      `json:"defaultMessageResetTimeoutMs"`
			CustomMessages               []string `json:"customMessages"`
			AllMessages                  []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"allMessages"`
		} `json:"doorbellSettings"`
		ID        string `json:"id"`
		IsAdopted bool   `json:"isAdopted"`
		IsAway    bool   `json:"isAway"`
		IsSetup   bool   `json:"isSetup"`
		Network   string `json:"network"`
		Type      string `json:"type"`
		UpSince   int64  `json:"upSince"`
		ModelKey  string `json:"modelKey"`
	} `json:"nvr"`
	LastUpdateID   string        `json:"lastUpdateId"`
	CloudPortalURL string        `json:"cloudPortalUrl"`
	Viewers        []interface{} `json:"viewers"`
	Lights         []interface{} `json:"lights"`
	Bridges        []interface{} `json:"bridges"`
	Sensors        []interface{} `json:"sensors"`
}

type Camera struct {
	IsDeleting            bool        `json:"isDeleting"`
	Mac                   string      `json:"mac"`
	Host                  string      `json:"host"`
	ConnectionHost        string      `json:"connectionHost"`
	Type                  string      `json:"type"`
	Name                  string      `json:"name"`
	UpSince               int64       `json:"upSince"`
	LastSeen              int64       `json:"lastSeen"`
	ConnectedSince        int64       `json:"connectedSince"`
	State                 string      `json:"state"`
	HardwareRevision      string      `json:"hardwareRevision"`
	FirmwareVersion       string      `json:"firmwareVersion"`
	FirmwareBuild         string      `json:"firmwareBuild"`
	IsUpdating            bool        `json:"isUpdating"`
	IsAdopting            bool        `json:"isAdopting"`
	IsAdopted             bool        `json:"isAdopted"`
	IsProvisioned         bool        `json:"isProvisioned"`
	IsRebooting           bool        `json:"isRebooting"`
	IsSSHEnabled          bool        `json:"isSshEnabled"`
	CanAdopt              bool        `json:"canAdopt"`
	IsAttemptingToConnect bool        `json:"isAttemptingToConnect"`
	IsHidden              interface{} `json:"isHidden"`
	LastMotion            int64       `json:"lastMotion"`
	MicVolume             int         `json:"micVolume"`
	IsMicEnabled          bool        `json:"isMicEnabled"`
	IsRecording           bool        `json:"isRecording"`
	IsMotionDetected      bool        `json:"isMotionDetected"`
	PhyRate               int         `json:"phyRate"`
	HdrMode               bool        `json:"hdrMode"`
	IsProbingForWifi      bool        `json:"isProbingForWifi"`
	ApMac                 interface{} `json:"apMac"`
	ApRssi                interface{} `json:"apRssi"`
	ElementInfo           interface{} `json:"elementInfo"`
	ChimeDuration         int         `json:"chimeDuration"`
	IsDark                bool        `json:"isDark"`
	LastRing              uint64      `json:"lastRing"`
	WiredConnectionState  struct {
		PhyRate int `json:"phyRate"`
	} `json:"wiredConnectionState"`
	Channels []struct {
		ID                       int         `json:"id"`
		VideoID                  string      `json:"videoId"`
		Name                     string      `json:"name"`
		Enabled                  bool        `json:"enabled"`
		IsRtspEnabled            bool        `json:"isRtspEnabled"`
		RtspAlias                interface{} `json:"rtspAlias"`
		Width                    int         `json:"width"`
		Height                   int         `json:"height"`
		Fps                      int         `json:"fps"`
		Bitrate                  int         `json:"bitrate"`
		MinBitrate               int         `json:"minBitrate"`
		MaxBitrate               int         `json:"maxBitrate"`
		MinClientAdaptiveBitRate int         `json:"minClientAdaptiveBitRate"`
		MinMotionAdaptiveBitRate int         `json:"minMotionAdaptiveBitRate"`
		FpsValues                []int       `json:"fpsValues"`
		IdrInterval              int         `json:"idrInterval"`
	} `json:"channels"`
	IspSettings struct {
		AeMode                         string `json:"aeMode"`
		IrLedMode                      string `json:"irLedMode"`
		IrLedLevel                     int    `json:"irLedLevel"`
		Wdr                            int    `json:"wdr"`
		IcrSensitivity                 int    `json:"icrSensitivity"`
		Brightness                     int    `json:"brightness"`
		Contrast                       int    `json:"contrast"`
		Hue                            int    `json:"hue"`
		Saturation                     int    `json:"saturation"`
		Sharpness                      int    `json:"sharpness"`
		Denoise                        int    `json:"denoise"`
		IsFlippedVertical              bool   `json:"isFlippedVertical"`
		IsFlippedHorizontal            bool   `json:"isFlippedHorizontal"`
		IsAutoRotateEnabled            bool   `json:"isAutoRotateEnabled"`
		IsLdcEnabled                   bool   `json:"isLdcEnabled"`
		Is3DnrEnabled                  bool   `json:"is3dnrEnabled"`
		IsExternalIrEnabled            bool   `json:"isExternalIrEnabled"`
		IsAggressiveAntiFlickerEnabled bool   `json:"isAggressiveAntiFlickerEnabled"`
		IsPauseMotionEnabled           bool   `json:"isPauseMotionEnabled"`
		DZoomCenterX                   int    `json:"dZoomCenterX"`
		DZoomCenterY                   int    `json:"dZoomCenterY"`
		DZoomScale                     int    `json:"dZoomScale"`
		DZoomStreamID                  int    `json:"dZoomStreamId"`
		FocusMode                      string `json:"focusMode"`
		FocusPosition                  int    `json:"focusPosition"`
		TouchFocusX                    int    `json:"touchFocusX"`
		TouchFocusY                    int    `json:"touchFocusY"`
		ZoomPosition                   int    `json:"zoomPosition"`
	} `json:"ispSettings"`
	TalkbackSettings struct {
		TypeFmt       string      `json:"typeFmt"`
		TypeIn        string      `json:"typeIn"`
		BindAddr      string      `json:"bindAddr"`
		BindPort      int         `json:"bindPort"`
		FilterAddr    interface{} `json:"filterAddr"`
		FilterPort    interface{} `json:"filterPort"`
		Channels      int         `json:"channels"`
		SamplingRate  int         `json:"samplingRate"`
		BitsPerSample int         `json:"bitsPerSample"`
		Quality       int         `json:"quality"`
	} `json:"talkbackSettings"`
	OsdSettings struct {
		IsNameEnabled  bool `json:"isNameEnabled"`
		IsDateEnabled  bool `json:"isDateEnabled"`
		IsLogoEnabled  bool `json:"isLogoEnabled"`
		IsDebugEnabled bool `json:"isDebugEnabled"`
	} `json:"osdSettings"`
	LedSettings struct {
		IsEnabled bool `json:"isEnabled"`
		BlinkRate int  `json:"blinkRate"`
	} `json:"ledSettings"`
	SpeakerSettings struct {
		IsEnabled              bool `json:"isEnabled"`
		AreSystemSoundsEnabled bool `json:"areSystemSoundsEnabled"`
		Volume                 int  `json:"volume"`
	} `json:"speakerSettings"`
	RecordingSettings struct {
		PrePaddingSecs            int         `json:"prePaddingSecs"`
		PostPaddingSecs           int         `json:"postPaddingSecs"`
		MinMotionEventTrigger     int         `json:"minMotionEventTrigger"`
		EndMotionEventDelay       int         `json:"endMotionEventDelay"`
		SuppressIlluminationSurge bool        `json:"suppressIlluminationSurge"`
		Mode                      string      `json:"mode"`
		Geofencing                string      `json:"geofencing"`
		RetentionDurationMs       interface{} `json:"retentionDurationMs"`
		UseNewMotionAlgorithm     bool        `json:"useNewMotionAlgorithm"`
		EnablePirTimelapse        bool        `json:"enablePirTimelapse"`
	} `json:"recordingSettings,omitempty"`
	RecordingSchedule interface{} `json:"recordingSchedule"`
	MotionZones       []struct {
		Name        string  `json:"name"`
		Color       string  `json:"color"`
		Points      [][]int `json:"points"`
		Sensitivity int     `json:"sensitivity"`
	} `json:"motionZones"`
	PrivacyZones []interface{} `json:"privacyZones"`
	Stats        struct {
		RxBytes int   `json:"rxBytes"`
		TxBytes int64 `json:"txBytes"`
		Wifi    struct {
			Channel        interface{} `json:"channel"`
			Frequency      interface{} `json:"frequency"`
			LinkSpeedMbps  interface{} `json:"linkSpeedMbps"`
			SignalQuality  int         `json:"signalQuality"`
			SignalStrength int         `json:"signalStrength"`
		} `json:"wifi"`
		Battery struct {
			Percentage interface{} `json:"percentage"`
			IsCharging bool        `json:"isCharging"`
			SleepState string      `json:"sleepState"`
		} `json:"battery"`
		Video struct {
			RecordingStart   int64 `json:"recordingStart"`
			RecordingEnd     int64 `json:"recordingEnd"`
			RecordingStartLQ int64 `json:"recordingStartLQ"`
			RecordingEndLQ   int64 `json:"recordingEndLQ"`
			TimelapseStart   int64 `json:"timelapseStart"`
			TimelapseEnd     int64 `json:"timelapseEnd"`
			TimelapseStartLQ int64 `json:"timelapseStartLQ"`
			TimelapseEndLQ   int64 `json:"timelapseEndLQ"`
		} `json:"video"`
		WifiQuality  int `json:"wifiQuality"`
		WifiStrength int `json:"wifiStrength"`
	} `json:"stats"`
	FeatureFlags struct {
		CanAdjustIrLedLevel   bool `json:"canAdjustIrLedLevel"`
		CanMagicZoom          bool `json:"canMagicZoom"`
		CanOpticalZoom        bool `json:"canOpticalZoom"`
		CanTouchFocus         bool `json:"canTouchFocus"`
		HasAccelerometer      bool `json:"hasAccelerometer"`
		HasAec                bool `json:"hasAec"`
		HasBattery            bool `json:"hasBattery"`
		HasBluetooth          bool `json:"hasBluetooth"`
		HasChime              bool `json:"hasChime"`
		HasExternalIr         bool `json:"hasExternalIr"`
		HasIcrSensitivity     bool `json:"hasIcrSensitivity"`
		HasLdc                bool `json:"hasLdc"`
		HasLedIr              bool `json:"hasLedIr"`
		HasLedStatus          bool `json:"hasLedStatus"`
		HasLineIn             bool `json:"hasLineIn"`
		HasMic                bool `json:"hasMic"`
		HasPrivacyMask        bool `json:"hasPrivacyMask"`
		HasRtc                bool `json:"hasRtc"`
		HasSdCard             bool `json:"hasSdCard"`
		HasSpeaker            bool `json:"hasSpeaker"`
		HasWifi               bool `json:"hasWifi"`
		HasHdr                bool `json:"hasHdr"`
		HasAutoICROnly        bool `json:"hasAutoICROnly"`
		HasMotionZones        bool `json:"hasMotionZones"`
		HasLcdScreen          bool `json:"hasLcdScreen"`
		HasNewMotionAlgorithm bool `json:"hasNewMotionAlgorithm"`
	} `json:"featureFlags"`
	PirSettings struct {
		PirSensitivity            int `json:"pirSensitivity"`
		PirMotionClipLength       int `json:"pirMotionClipLength"`
		TimelapseFrameInterval    int `json:"timelapseFrameInterval"`
		TimelapseTransferInterval int `json:"timelapseTransferInterval"`
	} `json:"pirSettings"`
	LcdMessage struct {
	} `json:"lcdMessage"`
	WifiConnectionState struct {
		Channel        interface{} `json:"channel"`
		Frequency      interface{} `json:"frequency"`
		PhyRate        interface{} `json:"phyRate"`
		SignalQuality  interface{} `json:"signalQuality"`
		SignalStrength interface{} `json:"signalStrength"`
	} `json:"wifiConnectionState"`
	ID           string `json:"id"`
	IsConnected  bool   `json:"isConnected"`
	Platform     string `json:"platform"`
	HasSpeaker   bool   `json:"hasSpeaker"`
	HasWifi      bool   `json:"hasWifi"`
	AudioBitrate int    `json:"audioBitrate"`
	CanManage    bool   `json:"canManage"`
	IsManaged    bool   `json:"isManaged"`
	ModelKey     string `json:"modelKey"`
}

func (c Camera) isDoorbell() bool {
	return c.Type == "UVC G4 Doorbell"
}

type Cameras []Camera

type Doorbell Camera
type Doorbells []Doorbell

func (c *Client) GetDoorbells(ctx context.Context) (Doorbells, error) {
	b, err := c.GetBootstrap(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to get doorbells: %w", err)
	}

	var ds Doorbells
	for _, c := range b.Cameras {
		if c.IsManaged && c.isDoorbell() {
			ds = append(ds, Doorbell(c))
		}
	}
	return ds, nil
}

func (d Doorbell) DoesRung(oldStates Doorbells) bool {
	for _, old := range oldStates {
		if d.Name == old.Name {
			return d.LastRing > old.LastRing
		}
	}
	return false
}

func (c *Client) GetBootstrap(ctx context.Context) (*Bootstrap, error) {
	u := c.baseURL()
	u.Path = "/api/bootstrap"

	var bootstrap Bootstrap
	if err := c.jsonRequest(ctx, http.MethodGet, u, nil, &bootstrap); err != nil {
		return nil, xerrors.Errorf("failed to get bootstrap: %w", err)
	}
	c.logger.Debugln("get bootstrap successfully")
	return &bootstrap, nil
}
