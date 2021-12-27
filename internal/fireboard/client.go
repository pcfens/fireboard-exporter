package fireboard

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const namespace = "fireboard"

var (
	up = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "up"),
		"Returns 1 if the fireboard is checking in with a temperature probe attached",
		[]string{"fireboard_name"}, nil,
	)

	batteryVolts = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "battery_volts"),
		"Battery voltage of the Fireboard",
		[]string{"fireboard_name"}, nil,
	)

	probeTemperature = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "probe_temperature_degrees"),
		"Probe temperature in degrees. Units are dependent on Fireboard settings.",
		[]string{"fireboard_name", "channel"}, nil,
	)

	txPower = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "transmit_power_db"),
		"Transmitter power",
		[]string{"fireboard_name"}, nil,
	)

	signalLevel = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "signal_level_db"),
		"Wifi Signal Level",
		[]string{"fireboard_name"}, nil,
	)

	nightMode = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "night_mode"),
		"Is the fireboard in night mode",
		[]string{"fireboard_name"}, nil,
	)

	cpuUsage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cpu_usage_percent"),
		"Current CPU usage",
		[]string{"fireboard_name"}, nil,
	)
)

type Devices []Device
type Device struct {
	UUID        string    `json:"uuid"`
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	HardwareID  string    `json:"hardware_id"`
	Created     time.Time `json:"created"`
	LatestTemps []struct {
		Degreetype int       `json:"degreetype"`
		Temp       float64   `json:"temp"`
		Created    time.Time `json:"created"`
		Channel    int       `json:"channel"`
	} `json:"latest_temps"`
	DeviceLog struct {
		BleClientMAC   string  `json:"bleClientMAC"`
		DeviceID       string  `json:"deviceID"`
		Mode           string  `json:"mode"`
		BleSignalLevel int     `json:"bleSignalLevel"`
		InternalIP     string  `json:"internalIP"`
		VBattPerRaw    float64 `json:"vBattPerRaw"`
		AuxPort        string  `json:"auxPort"`
		Frequency      string  `json:"frequency"`
		VersionEspHal  string  `json:"versionEspHal"`
		TimeZoneBT     string  `json:"timeZoneBT"`
		Ssid           string  `json:"ssid"`
		DiskUsage      string  `json:"diskUsage"`
		VersionNode    string  `json:"versionNode"`
		Band           string  `json:"band"`
		Contrast       string  `json:"contrast"`
		Version        string  `json:"version"`
		VersionJava    string  `json:"versionJava"`
		VBattPer       float64 `json:"vBattPer"`
		VersionImage   string  `json:"versionImage"`
		VBatt          float64 `json:"vBatt"`
		Uptime         string  `json:"uptime"`
		Date           string  `json:"date"`
		MacNIC         string  `json:"macNIC"`
		Nightmode      bool    `json:"nightmode"`
		CommercialMode string  `json:"commercialMode"`
		Model          string  `json:"model"`
		VersionUtils   string  `json:"versionUtils"`
		Signallevel    int     `json:"signallevel"`
		MacAP          string  `json:"macAP"`
		PublicIP       string  `json:"publicIP"`
		TempFilter     string  `json:"tempFilter"`
		CPUUsage       string  `json:"cpuUsage"`
		BoardID        string  `json:"boardID"`
		Drivesettings  string  `json:"drivesettings"`
		Txpower        int     `json:"txpower"`
		OnboardTemp    float64 `json:"onboardTemp"`
		MemUsage       string  `json:"memUsage"`
		Linkquality    string  `json:"linkquality"`
	} `json:"device_log"`
	LastTemplog        time.Time `json:"last_templog"`
	LastBatteryReading float64   `json:"last_battery_reading"`
	Channels           []struct {
		LastTemplog struct {
			Degreetype int       `json:"degreetype"`
			Temp       float64   `json:"temp"`
			Created    time.Time `json:"created"`
			Channel    int       `json:"channel"`
		} `json:"last_templog,omitempty"`
		CurrentTemp      float64       `json:"current_temp,omitempty"`
		Enabled          bool          `json:"enabled"`
		Alerts           []interface{} `json:"alerts"`
		ChannelLabel     string        `json:"channel_label"`
		Channel          int           `json:"channel"`
		RangeAverageTemp float64       `json:"range_average_temp,omitempty"`
		RangeMinTemp     float64       `json:"range_min_temp,omitempty"`
		RangeMaxTemp     float64       `json:"range_max_temp,omitempty"`
		Created          time.Time     `json:"created"`
		Sessionid        int           `json:"sessionid"`
		ID               int           `json:"id"`
	} `json:"channels"`
	FbjVersion   string      `json:"fbj_version"`
	FbnVersion   string      `json:"fbn_version"`
	FbuVersion   string      `json:"fbu_version"`
	Version      string      `json:"version"`
	ProbeConfig  string      `json:"probe_config"`
	LastDrivelog interface{} `json:"last_drivelog"`
	ChannelCount int         `json:"channel_count"`
	Degreetype   int         `json:"degreetype"`
	Model        string      `json:"model"`
	Active       bool        `json:"active"`
}

type fireboard struct {
	token string
}

func New(token string) fireboard {
	fc := fireboard{token}
	return fc
}

func (fc fireboard) getDevices() (Devices, error) {
	devicesJson, err := fc.fireboardGet("https://fireboard.io/api/v1/devices.json")
	if err != nil {
		return Devices{}, err
	}
	var devices Devices
	json.Unmarshal([]byte(devicesJson), &devices)
	return devices, nil
}

func (fc fireboard) fireboardGet(url string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed constructing HTTP request")
	}
	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Token %s", fc.token)},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "Failed to make http request", err
	}
	defer resp.Body.Close()
	jsonString, err := io.ReadAll(resp.Body)
	if err != nil {
		return string(jsonString), err
	}
	return string(jsonString), nil
}

func (fc fireboard) Describe(ch chan<- *prometheus.Desc) {
	ch <- up
	ch <- batteryVolts
	ch <- probeTemperature
	ch <- txPower
	ch <- signalLevel
	ch <- nightMode
	ch <- cpuUsage
}

func (fc fireboard) Collect(ch chan<- prometheus.Metric) {
	devices, err := fc.getDevices()
	if err != nil {
		log.Fatal((err))
	}
	for i := 0; i < len(devices); i++ {

		if len(devices[i].LatestTemps) > 0 {
			ch <- prometheus.MustNewConstMetric(
				up, prometheus.GaugeValue, 1, devices[i].Title,
			)
			fc.updateMetrics(devices[i], ch)
		} else {
			ch <- prometheus.MustNewConstMetric(
				up, prometheus.GaugeValue, 0, devices[i].Title,
			)
		}
	}
}

func (fc fireboard) updateMetrics(device Device, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		batteryVolts, prometheus.GaugeValue, float64(device.DeviceLog.VBatt), device.Title,
	)

	ch <- prometheus.MustNewConstMetric(
		txPower, prometheus.GaugeValue, float64(device.DeviceLog.Txpower), device.Title,
	)

	ch <- prometheus.MustNewConstMetric(
		signalLevel, prometheus.GaugeValue, float64(device.DeviceLog.Signallevel), device.Title,
	)

	if device.DeviceLog.Nightmode {
		ch <- prometheus.MustNewConstMetric(
			nightMode, prometheus.GaugeValue, 1, device.Title,
		)
	} else {
		ch <- prometheus.MustNewConstMetric(
			nightMode, prometheus.GaugeValue, 0, device.Title,
		)
	}

	currentCpu, err := strconv.Atoi(strings.Split(device.DeviceLog.CPUUsage, "%")[0])
	if err != nil {
		log.Fatal(err)
	}
	ch <- prometheus.MustNewConstMetric(
		cpuUsage, prometheus.GaugeValue, float64(currentCpu), device.Title,
	)

	for j := 0; j < len(device.LatestTemps); j++ {
		channel := device.LatestTemps[j]
		ch <- prometheus.MustNewConstMetric(
			probeTemperature, prometheus.GaugeValue, channel.Temp, device.Title, strconv.Itoa(j+1),
		)
	}
}
