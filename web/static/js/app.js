// 空气质量监测系统 - 前端应用脚本

// 全局配置
const AppConfig = {
    apiBaseUrl: '/api/v1',
    webApiBaseUrl: '/web/api',
    refreshInterval: 30000, // 30秒
    chartColors: {
        pm25: 'rgb(255, 99, 132)',
        pm10: 'rgb(54, 162, 235)',
        temp: 'rgb(255, 205, 86)',
        humidity: 'rgb(75, 192, 192)'
    }
};

// 工具函数
const Utils = {
    // 格式化时间
    formatTime: function(timestamp) {
        const date = new Date(timestamp * 1000);
        return date.toLocaleString('zh-CN');
    },
    
    // 格式化持续时间
    formatDuration: function(seconds) {
        if (seconds < 60) return '刚刚';
        if (seconds < 3600) return Math.floor(seconds / 60) + '分钟前';
        if (seconds < 86400) return Math.floor(seconds / 3600) + '小时前';
        return Math.floor(seconds / 86400) + '天前';
    },
    
    // 获取空气质量等级
    getAirQualityLevel: function(pm25) {
        if (pm25 <= 35) return { level: '优', class: 'air-quality-excellent' };
        if (pm25 <= 75) return { level: '良', class: 'air-quality-good' };
        if (pm25 <= 115) return { level: '轻度污染', class: 'air-quality-moderate' };
        if (pm25 <= 150) return { level: '中度污染', class: 'air-quality-unhealthy' };
        return { level: '重度污染', class: 'air-quality-hazardous' };
    },
    
    // 显示通知
    showNotification: function(message, type = 'info') {
        const alertClass = {
            'success': 'alert-success',
            'error': 'alert-danger',
            'warning': 'alert-warning',
            'info': 'alert-info'
        }[type] || 'alert-info';
        
        const notification = `
            <div class="alert ${alertClass} alert-dismissible fade show" role="alert">
                ${message}
                <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
            </div>
        `;
        
        // 在页面顶部显示通知
        const container = document.querySelector('.container-fluid');
        if (container) {
            container.insertAdjacentHTML('afterbegin', notification);
            
            // 5秒后自动隐藏
            setTimeout(() => {
                const alert = container.querySelector('.alert');
                if (alert) {
                    alert.remove();
                }
            }, 5000);
        }
    }
};

// 页面初始化
document.addEventListener('DOMContentLoaded', function() {
    // 初始化工具提示
    const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });
    
    // 初始化弹出框
    const popoverTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="popover"]'));
    popoverTriggerList.map(function (popoverTriggerEl) {
        return new bootstrap.Popover(popoverTriggerEl);
    });
    
    // 设置当前时间
    updateCurrentTime();
    setInterval(updateCurrentTime, 1000);
});

// 更新当前时间
function updateCurrentTime() {
    const timeElements = document.querySelectorAll('.current-time');
    const now = new Date();
    const timeString = now.toLocaleString('zh-CN');
    
    timeElements.forEach(element => {
        element.textContent = timeString;
    });
}

// 导出全局函数
window.Utils = Utils;
