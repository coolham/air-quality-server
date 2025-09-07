package handlers

import (
	"air-quality-server/internal/models"
	"context"
	"fmt"
	"sort"
	"time"
)

// getFloatValue 安全获取浮点数值
func getFloatValue(ptr *float64) float64 {
	if ptr == nil {
		return 0.0
	}
	return *ptr
}

// getIntValueFromPointer 安全获取整数值
func getIntValueFromPointer(ptr *int) int {
	if ptr == nil {
		return 0
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
	// 从UnifiedSensorData获取最新数据，按设备ID和传感器ID分组
	latestData, err := h.services.UnifiedSensorData.GetAllData(ctx, 50, 0) // 获取最新50条数据
	if err != nil {
		return nil, err
	}

	// 按设备ID和传感器ID分组，获取每个组合的最新数据
	latestByDeviceSensor := make(map[string]models.UnifiedSensorData)
	for _, data := range latestData {
		key := data.DeviceID + ":" + data.SensorID
		if existing, exists := latestByDeviceSensor[key]; !exists || data.Timestamp.After(existing.Timestamp) {
			latestByDeviceSensor[key] = data
		}
	}

	var summaries []AirQualityDataSummary
	for _, data := range latestByDeviceSensor {
		// 判断设备状态（基于数据时间戳）
		status := "online"
		if time.Since(data.Timestamp) > 5*time.Minute {
			status = "offline"
		}

		summary := AirQualityDataSummary{
			DeviceID:     data.DeviceID,
			DeviceName:   data.DeviceID, // 使用设备ID作为名称，后续可以从设备表获取真实名称
			DeviceType:   string(data.DeviceType),
			SensorID:     data.SensorID,
			SensorType:   data.SensorType,
			PM25:         getFloatValueFromPointer(data.PM25),
			PM10:         getFloatValueFromPointer(data.PM10),
			CO2:          getFloatValueFromPointer(data.CO2),
			Formaldehyde: getFloatValueFromPointer(data.Formaldehyde),
			Temp:         getFloatValueFromPointer(data.Temperature),
			Humidity:     getFloatValueFromPointer(data.Humidity),
			Pressure:     getFloatValueFromPointer(data.Pressure),
			Battery:      getIntValueFromPointer(data.Battery),
			DataQuality:  data.DataQuality,
			CreatedAt:    data.Timestamp,
			Status:       status,
		}
		summaries = append(summaries, summary)
	}

	// 按时间戳降序排序
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].CreatedAt.After(summaries[j].CreatedAt)
	})

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

// convertToChartDataFromUnified 将统一传感器数据转换为图表数据格式，支持指标选择和传感器过滤
func (h *WebHandlers) convertToChartDataFromUnified(historyData []models.UnifiedSensorData, metric string, sensorID string) *ChartData {
	var labels []string
	var pm25Data, formaldehydeData, tempData, humidityData []float64

	for _, data := range historyData {
		// 如果指定了传感器ID，则过滤数据
		if sensorID != "" && data.SensorID != sensorID {
			continue
		}

		// 格式化时间标签
		label := data.Timestamp.Format("15:04")
		labels = append(labels, label)

		// 添加数据点
		pm25Data = append(pm25Data, getFloatValueFromPointer(data.PM25))
		formaldehydeData = append(formaldehydeData, getFloatValueFromPointer(data.Formaldehyde))
		tempData = append(tempData, getFloatValueFromPointer(data.Temperature))
		humidityData = append(humidityData, getFloatValueFromPointer(data.Humidity))
	}

	var datasets []Dataset

	// 根据选择的指标添加数据集
	switch metric {
	case "pm25":
		datasets = append(datasets, Dataset{
			Label:           "PM2.5",
			Data:            pm25Data,
			BorderColor:     "rgb(255, 99, 132)",
			BackgroundColor: "rgba(255, 99, 132, 0.2)",
			Fill:            false,
		})
	case "formaldehyde":
		datasets = append(datasets, Dataset{
			Label:           "甲醛",
			Data:            formaldehydeData,
			BorderColor:     "rgb(220, 53, 69)",
			BackgroundColor: "rgba(220, 53, 69, 0.2)",
			Fill:            false,
		})
	case "temperature":
		datasets = append(datasets, Dataset{
			Label:           "温度",
			Data:            tempData,
			BorderColor:     "rgb(255, 205, 86)",
			BackgroundColor: "rgba(255, 205, 86, 0.2)",
			Fill:            false,
		})
	case "humidity":
		datasets = append(datasets, Dataset{
			Label:           "湿度",
			Data:            humidityData,
			BorderColor:     "rgb(75, 192, 192)",
			BackgroundColor: "rgba(75, 192, 192, 0.2)",
			Fill:            false,
		})
	default: // "all"
		datasets = []Dataset{
			{
				Label:           "PM2.5",
				Data:            pm25Data,
				BorderColor:     "rgb(255, 99, 132)",
				BackgroundColor: "rgba(255, 99, 132, 0.2)",
				Fill:            false,
			},
			{
				Label:           "甲醛",
				Data:            formaldehydeData,
				BorderColor:     "rgb(220, 53, 69)",
				BackgroundColor: "rgba(220, 53, 69, 0.2)",
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
		}
	}

	return &ChartData{
		Labels:   labels,
		Datasets: datasets,
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

// getFloatValueFromPointer 从指针获取float64值，处理nil指针
func getFloatValueFromPointer(ptr *float64) float64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
