package handlers

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"
	"context"
	"fmt"
)

// getFloatValue 安全获取浮点数值
func getFloatValue(ptr *float64) float64 {
	if ptr == nil {
		return 0.0
	}
	return *ptr
}

// getIntValue 安全获取整数值
func getIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// getDeviceStats 获取设备统计信息
func (h *WebHandlers) getDeviceStats(ctx context.Context) (*DeviceStats, error) {
	// 获取设备总数
	total, err := h.services.Device.CountDevices(ctx)
	if err != nil {
		return nil, err
	}

	// 获取在线设备数（这里简化处理，实际应该根据设备状态统计）
	onlineDevices := 0 // TODO: 实现在线设备统计
	offlineDevices := int(total) - onlineDevices
	activeDevices := onlineDevices // 假设在线设备都是活跃的

	return &DeviceStats{
		TotalDevices:   int(total),
		OnlineDevices:  onlineDevices,
		OfflineDevices: offlineDevices,
		ActiveDevices:  activeDevices,
	}, nil
}

// getLatestData 获取最新数据
func (h *WebHandlers) getLatestData(ctx context.Context) ([]AirQualityDataSummary, error) {
	// 获取所有设备
	devices, err := h.services.Device.ListDevices(ctx, 10, 0) // 获取前10个设备
	if err != nil {
		return nil, err
	}

	var summaries []AirQualityDataSummary
	for _, device := range devices {
		// 获取设备最新数据
		latestData, err := h.services.AirQuality.GetLatestData(ctx, device.ID)
		if err != nil {
			h.logger.Warn("获取设备最新数据失败", utils.String("device_id", device.ID), utils.ErrorField(err))
			continue
		}

		if latestData != nil {
			summary := AirQualityDataSummary{
				DeviceID:   device.ID,
				DeviceName: device.Name,
				PM25:       getFloatValue(latestData.PM25),
				PM10:       getFloatValue(latestData.PM10),
				Temp:       getFloatValue(latestData.Temperature),
				Humidity:   getFloatValue(latestData.Humidity),
				CreatedAt:  latestData.CreatedAt,
				Status:     string(device.Status),
			}
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

// getAlertStats 获取告警统计信息
func (h *WebHandlers) getAlertStats(ctx context.Context) (*AlertStats, error) {
	// 获取告警总数
	total, err := h.services.Alert.CountAlerts(ctx)
	if err != nil {
		return nil, err
	}

	// 获取未解决告警数
	unresolvedAlerts := 0 // TODO: 实现未解决告警统计
	criticalAlerts := 0   // TODO: 实现严重告警统计
	warningAlerts := 0    // TODO: 实现警告告警统计

	return &AlertStats{
		TotalAlerts:      int(total),
		UnresolvedAlerts: unresolvedAlerts,
		CriticalAlerts:   criticalAlerts,
		WarningAlerts:    warningAlerts,
	}, nil
}

// convertToChartData 将历史数据转换为图表数据格式
func (h *WebHandlers) convertToChartData(historyData []models.AirQualityData) *ChartData {
	var labels []string
	var pm25Data, pm10Data, tempData, humidityData []float64

	for _, data := range historyData {
		// 格式化时间标签
		label := data.CreatedAt.Format("15:04")
		labels = append(labels, label)

		// 添加数据点
		pm25Data = append(pm25Data, getFloatValue(data.PM25))
		pm10Data = append(pm10Data, getFloatValue(data.PM10))
		tempData = append(tempData, getFloatValue(data.Temperature))
		humidityData = append(humidityData, getFloatValue(data.Humidity))
	}

	return &ChartData{
		Labels: labels,
		Datasets: []Dataset{
			{
				Label:           "PM2.5",
				Data:            pm25Data,
				BorderColor:     "rgb(255, 99, 132)",
				BackgroundColor: "rgba(255, 99, 132, 0.2)",
				Fill:            false,
			},
			{
				Label:           "PM10",
				Data:            pm10Data,
				BorderColor:     "rgb(54, 162, 235)",
				BackgroundColor: "rgba(54, 162, 235, 0.2)",
				Fill:            false,
			},
			{
				Label:           "温度",
				Data:            tempData,
				BorderColor:     "rgb(255, 205, 86)",
				BackgroundColor: "rgba(255, 205, 86, 0.2)",
				Fill:            false,
			},
			{
				Label:           "湿度",
				Data:            humidityData,
				BorderColor:     "rgb(75, 192, 192)",
				BackgroundColor: "rgba(75, 192, 192, 0.2)",
				Fill:            false,
			},
		},
	}
}

// convertToCSV 将统一传感器数据转换为CSV格式
func (h *WebHandlers) convertToCSV(data []models.UnifiedSensorData) string {
	if len(data) == 0 {
		return ""
	}

	// CSV头部
	csv := "ID,设备ID,设备类型,传感器ID,传感器类型,时间戳,PM2.5,PM10,CO2,甲醛,温度,湿度,气压,电池,数据质量,纬度,经度,地址,质量评分,信号强度,创建时间\n"

	// 数据行
	for _, item := range data {
		csv += fmt.Sprintf("%d,%s,%s,%s,%s,%s,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%.2f,%s,%.8f,%.8f,%s,%.2f,%d,%s\n",
			item.ID,
			item.DeviceID,
			item.DeviceType,
			item.SensorID,
			item.SensorType,
			item.Timestamp.Format("2006-01-02 15:04:05"),
			getFloatValue(item.PM25),
			getFloatValue(item.PM10),
			getFloatValue(item.CO2),
			getFloatValue(item.Formaldehyde),
			getFloatValue(item.Temperature),
			getFloatValue(item.Humidity),
			getFloatValue(item.Pressure),
			float64(getIntValue(item.Battery)),
			item.DataQuality,
			getFloatValue(item.Latitude),
			getFloatValue(item.Longitude),
			"",  // Address字段不存在
			0.0, // QualityScore字段不存在
			getIntValue(item.SignalStrength),
			item.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	return csv
}
