# 智慧海洋数据分析平台 API 接口文档

## 接口规范

### 基本信息
- 基础URL: `/api/v1`
- 请求/响应格式: JSON
- 认证方式: JWT (JSON Web Token)
- 状态码:
  - 200: 成功
  - 400: 请求错误
  - 401: 未授权
  - 403: 禁止访问
  - 404: 资源不存在
  - 500: 服务器错误

### 通用响应格式
```json
{
  "code": 200,          // 状态码
  "message": "成功",     // 状态描述
  "data": {},           // 响应数据
  "timestamp": 1634567890123 // 时间戳
}
```

## 1. 用户认证模块

### 1.1 用户注册

- **URL**: `/auth/register`
- **方法**: POST
- **描述**: 创建新用户账号
- **请求参数**:
  ```json
  {
    "username": "user001",
    "password": "Password123",
    "email": "user001@example.com",
    "fullName": "张三",
    "organization": "海洋研究所"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "注册成功",
    "data": {
      "userId": "u12345",
      "username": "user001",
      "email": "user001@example.com",
      "role": "researcher",
      "createTime": "2023-10-01T12:00:00Z"
    },
    "timestamp": 1634567890123
  }
  ```

### 1.2 用户登录

- **URL**: `/auth/login`
- **方法**: POST
- **描述**: 用户登录获取token
- **请求参数**:
  ```json
  {
    "username": "user001",
    "password": "Password123"
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "登录成功",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5...",
      "refreshToken": "eyJhbGciOiJIUzI1NiIsIn...",
      "expiresIn": 3600,
      "userId": "u12345",
      "username": "user001",
      "role": "researcher"
    },
    "timestamp": 1634567890123
  }
  ```

### 1.3 刷新Token

- **URL**: `/auth/refresh-token`
- **方法**: POST
- **描述**: 使用刷新令牌获取新的访问令牌
- **请求参数**:
  ```json
  {
    "refreshToken": "eyJhbGciOiJIUzI1NiIsIn..."
  }
  ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "刷新成功",
    "data": {
      "token": "eyJhbGciOiJIUzI1NiIsInR5...",
      "expiresIn": 3600
    },
    "timestamp": 1634567890123
  }
  ```

### 1.4 退出登录

- **URL**: `/auth/logout`
- **方法**: POST
- **描述**: 用户退出登录
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "登出成功",
    "data": null,
    "timestamp": 1634567890123
  }
  ```

### 1.5 获取当前用户信息

- **URL**: `/users/current`
- **方法**: GET
- **描述**: 获取当前登录用户信息
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "userId": "u12345",
      "username": "user001",
      "email": "user001@example.com",
      "fullName": "张三",
      "organization": "海洋研究所",
      "role": "researcher",
      "permissions": ["data:read", "data:write", "analysis:use"],
      "createTime": "2023-10-01T12:00:00Z",
      "lastLoginTime": "2023-10-15T08:30:00Z"
    },
    "timestamp": 1634567890123
  }
  ```

## 2. 数据管理模块

### 2.1 获取数据集列表

- **URL**: `/datasets`
- **方法**: GET
- **描述**: 获取数据集列表，支持分页和筛选
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**:
  - `page`: 页码，默认1
  - `size`: 每页条数，默认10
  - `type`: 数据类型，可选 ["temperature", "salinity", "wave", "current", "level"]
  - `startDate`: 开始日期，格式YYYY-MM-DD
  - `endDate`: 结束日期，格式YYYY-MM-DD
  - `region`: 区域，格式 "minLat,minLng,maxLat,maxLng"
  - `keyword`: 关键词搜索
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "total": 120,
      "page": 1,
      "size": 10,
      "datasets": [
        {
          "id": "ds001",
          "name": "北太平洋温盐数据集2023",
          "description": "2023年1月至6月北太平洋表层温度和盐度观测数据",
          "type": "temperature",
          "format": "netCDF",
          "region": {
            "name": "北太平洋",
            "bounds": [20.0, 120.0, 45.0, 170.0]
          },
          "timeRange": {
            "start": "2023-01-01T00:00:00Z",
            "end": "2023-06-30T23:59:59Z"
          },
          "resolution": {
            "spatial": "0.25度",
            "temporal": "日均值"
          },
          "size": 1240000000,
          "createdBy": "系统管理员",
          "createdAt": "2023-07-15T10:30:00Z",
          "tags": ["temperature", "salinity", "2023", "Pacific"]
        },
        // ... 更多数据集
      ]
    },
    "timestamp": 1634567890123
  }
  ```

### 2.2 获取数据集详情

- **URL**: `/datasets/{datasetId}`
- **方法**: GET
- **描述**: 获取指定数据集的详细信息
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "id": "ds001",
      "name": "北太平洋温盐数据集2023",
      "description": "2023年1月至6月北太平洋表层温度和盐度观测数据",
      "type": "temperature",
      "format": "netCDF",
      "region": {
        "name": "北太平洋",
        "bounds": [20.0, 120.0, 45.0, 170.0]
      },
      "timeRange": {
        "start": "2023-01-01T00:00:00Z",
        "end": "2023-06-30T23:59:59Z"
      },
      "resolution": {
        "spatial": "0.25度",
        "temporal": "日均值"
      },
      "size": 1240000000,
      "variables": [
        {
          "name": "sea_surface_temperature",
          "unit": "°C",
          "description": "海表温度",
          "range": [-2.0, 35.0]
        },
        {
          "name": "sea_surface_salinity",
          "unit": "PSU",
          "description": "海表盐度",
          "range": [30.0, 38.0]
        }
      ],
      "source": "全球海洋观测系统",
      "methodology": "卫星遥感结合浮标观测",
      "createdBy": "系统管理员",
      "createdAt": "2023-07-15T10:30:00Z",
      "updatedAt": "2023-07-16T08:45:00Z",
      "downloadCount": 45,
      "tags": ["temperature", "salinity", "2023", "Pacific"]
    },
    "timestamp": 1634567890123
  }
  ```

### 2.3 上传数据集

- **URL**: `/datasets/upload`
- **方法**: POST
- **描述**: 上传新的数据集文件
- **请求头**: 
  - `Authorization: Bearer {token}`
  - `Content-Type: multipart/form-data`
- **请求参数**:
  - `file`: 数据文件
  - `metadata`: 数据集元数据 (JSON字符串)
    ```json
    {
      "name": "南海海浪数据2023",
      "description": "2023年南海海浪观测数据",
      "type": "wave",
      "region": {
        "name": "南海",
        "bounds": [5.0, 105.0, 25.0, 125.0]
      },
      "timeRange": {
        "start": "2023-01-01T00:00:00Z",
        "end": "2023-12-31T23:59:59Z"
      },
      "tags": ["wave", "2023", "South China Sea"]
    }
    ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "上传成功",
    "data": {
      "datasetId": "ds123",
      "name": "南海海浪数据2023",
      "uploadTime": "2023-10-15T14:30:00Z",
      "size": 852000000,
      "status": "processing"
    },
    "timestamp": 1634567890123
  }
  ```

### 2.4 下载数据集

- **URL**: `/datasets/{datasetId}/download`
- **方法**: GET
- **描述**: 下载指定数据集
- **请求头**: `Authorization: Bearer {token}`
- **响应**: 文件流

## 3. 分析功能模块

### 3.1 温盐分析

#### 3.1.1 获取温盐数据时间序列

- **URL**: `/analysis/temperature-salinity/timeseries`
- **方法**: GET
- **描述**: 获取指定位置和时间范围的温盐数据时间序列
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**:
  - `datasetId`: 数据集ID
  - `lat`: 纬度
  - `lng`: 经度
  - `depth`: 深度(米)，可选
  - `startDate`: 开始时间
  - `endDate`: 结束时间
  - `interval`: 时间间隔，可选 ["hour", "day", "week", "month"]
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "location": {
        "lat": 22.5,
        "lng": 114.5,
        "depth": 0
      },
      "timeRange": {
        "start": "2023-01-01T00:00:00Z",
        "end": "2023-01-31T23:59:59Z"
      },
      "interval": "day",
      "series": [
        {
          "timestamp": "2023-01-01T00:00:00Z",
          "temperature": 25.4,
          "salinity": 33.2
        },
        {
          "timestamp": "2023-01-02T00:00:00Z",
          "temperature": 25.1,
          "salinity": 33.4
        },
        // ... 更多数据点
      ]
    },
    "timestamp": 1634567890123
  }
  ```

#### 3.1.2 获取温盐空间分布

- **URL**: `/analysis/temperature-salinity/spatial`
- **方法**: GET
- **描述**: 获取指定时间的温盐空间分布数据
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**:
  - `datasetId`: 数据集ID
  - `date`: 日期时间
  - `depth`: 深度(米)，可选
  - `bounds`: 边界范围，格式 "minLat,minLng,maxLat,maxLng"
  - `resolution`: 分辨率，可选 ["low", "medium", "high"]
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "time": "2023-01-15T00:00:00Z",
      "depth": 0,
      "bounds": [20.0, 110.0, 25.0, 120.0],
      "resolution": "medium",
      "grid": {
        "latCount": 20,
        "lngCount": 40,
        "latStep": 0.25,
        "lngStep": 0.25,
        "startLat": 20.0,
        "startLng": 110.0
      },
      "temperature": [
        [25.1, 25.2, 25.3, /* ... */],
        [24.9, 25.0, 25.1, /* ... */],
        // ... 更多数据行
      ],
      "salinity": [
        [33.2, 33.3, 33.4, /* ... */],
        [33.1, 33.2, 33.3, /* ... */],
        // ... 更多数据行
      ]
    },
    "timestamp": 1634567890123
  }
  ```

### 3.2 海标高度分析

#### 3.2.1 获取海标高度时间序列

- **URL**: `/analysis/sea-level/timeseries`
- **方法**: GET
- **描述**: 获取指定位置和时间范围的海平面高度时间序列
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**:
  - `datasetId`: 数据集ID
  - `lat`: 纬度
  - `lng`: 经度
  - `startDate`: 开始时间
  - `endDate`: 结束时间
  - `interval`: 时间间隔，可选 ["hour", "day", "week", "month"]
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "location": {
        "lat": 22.5,
        "lng": 114.5
      },
      "timeRange": {
        "start": "2023-01-01T00:00:00Z",
        "end": "2023-01-31T23:59:59Z"
      },
      "interval": "day",
      "series": [
        {
          "timestamp": "2023-01-01T00:00:00Z",
          "seaLevel": 0.45
        },
        {
          "timestamp": "2023-01-02T00:00:00Z",
          "seaLevel": 0.48
        },
        // ... 更多数据点
      ],
      "unit": "m",
      "reference": "平均海平面"
    },
    "timestamp": 1634567890123
  }
  ```

#### 3.2.2 获取海标高度空间分布

- **URL**: `/analysis/sea-level/spatial`
- **方法**: GET
- **描述**: 获取指定时间的海平面高度空间分布
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**:
  - `datasetId`: 数据集ID
  - `date`: 日期时间
  - `bounds`: 边界范围，格式 "minLat,minLng,maxLat,maxLng"
  - `resolution`: 分辨率，可选 ["low", "medium", "high"]
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "time": "2023-01-15T00:00:00Z",
      "bounds": [20.0, 110.0, 25.0, 120.0],
      "resolution": "medium",
      "grid": {
        "latCount": 20,
        "lngCount": 40,
        "latStep": 0.25,
        "lngStep": 0.25,
        "startLat": 20.0,
        "startLng": 110.0
      },
      "seaLevel": [
        [0.45, 0.46, 0.47, /* ... */],
        [0.44, 0.45, 0.46, /* ... */],
        // ... 更多数据行
      ],
      "unit": "m",
      "reference": "平均海平面"
    },
    "timestamp": 1634567890123
  }
  ```

### 3.3 温度预报分析

#### 3.3.1 获取温度预报数据

- **URL**: `/forecasts/temperature`
- **方法**: GET
- **描述**: 获取指定区域和时间范围的温度预报数据
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**:
  - `modelId`: 预报模型ID
  - `region`: 区域名称或边界范围
  - `depth`: 深度(米)，可选
  - `forecastDate`: 预报基准日期
  - `forecastDays`: 预报天数，默认7
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "modelInfo": {
        "id": "model001",
        "name": "全球海洋温度预报模型v2.1",
        "version": "2.1.0",
        "description": "基于机器学习的全球海洋温度预报系统"
      },
      "forecastInfo": {
        "baseTime": "2023-10-15T00:00:00Z",
        "forecastDays": 7,
        "updateInterval": "24h",
        "spatialResolution": "0.25度"
      },
      "region": {
        "name": "南海北部",
        "bounds": [18.0, 110.0, 23.0, 118.0]
      },
      "depth": 0,
      "forecasts": [
        {
          "forecastTime": "2023-10-16T00:00:00Z",
          "temperatureGrid": {
            "grid": {
              "latCount": 20,
              "lngCount": 32,
              "latStep": 0.25,
              "lngStep": 0.25,
              "startLat": 18.0,
              "startLng": 110.0
            },
            "values": [
              [27.1, 27.2, 27.3, /* ... */],
              [27.0, 27.1, 27.2, /* ... */],
              // ... 更多数据行
            ]
          }
        },
        // ... 更多预报时间点
      ],
      "accuracy": {
        "rmse": 0.42,
        "mae": 0.31
      }
    },
    "timestamp": 1634567890123
  }
  ```

### 3.4 海浪视频反演分析

#### 3.4.1 上传海浪视频

- **URL**: `/analysis/wave-inversion/upload`
- **方法**: POST
- **描述**: 上传海浪视频进行反演分析
- **请求头**: 
  - `Authorization: Bearer {token}`
  - `Content-Type: multipart/form-data`
- **请求参数**:
  - `videoFile`: 视频文件
  - `metadata`: 元数据 (JSON字符串)
    ```json
    {
      "location": {
        "lat": 22.5,
        "lng": 114.5,
        "name": "深圳湾"
      },
      "captureTime": "2023-10-14T10:30:00Z",
      "duration": 300,
      "cameraHeight": 15.5,
      "cameraAngle": 45,
      "description": "台风"海葵"过境期间拍摄"
    }
    ```
- **响应**:
  ```json
  {
    "code": 200,
    "message": "上传成功",
    "data": {
      "taskId": "task001",
      "status": "queued",
      "estimatedProcessingTime": 300,
      "videoInfo": {
        "name": "wave_shenzhen_20231014.mp4",
        "size": 45000000,
        "duration": 300,
        "resolution": "1920x1080"
      },
      "uploadTime": "2023-10-15T15:45:00Z"
    },
    "timestamp": 1634567890123
  }
  ```

#### 3.4.2 获取海浪反演任务状态

- **URL**: `/analysis/wave-inversion/tasks/{taskId}`
- **方法**: GET
- **描述**: 获取海浪视频反演任务的处理状态
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "taskId": "task001",
      "status": "processing",
      "progress": 65,
      "startTime": "2023-10-15T15:46:00Z",
      "estimatedCompletionTime": "2023-10-15T15:51:00Z",
      "currentStep": "波浪谱分析"
    },
    "timestamp": 1634567890123
  }
  ```

#### 3.4.3 获取海浪反演结果

- **URL**: `/analysis/wave-inversion/results/{taskId}`
- **方法**: GET
- **描述**: 获取海浪视频反演的分析结果
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "taskId": "task001",
      "status": "completed",
      "processingTime": 285,
      "videoInfo": {
        "name": "wave_shenzhen_20231014.mp4",
        "location": {
          "lat": 22.5,
          "lng": 114.5,
          "name": "深圳湾"
        },
        "captureTime": "2023-10-14T10:30:00Z"
      },
      "waveParameters": {
        "significantWaveHeight": 2.3,
        "peakPeriod": 8.5,
        "meanDirection": 135,
        "directionalSpread": 30
      },
      "waveSpectrum": {
        "frequencyRange": [0.05, 0.5],
        "directionRange": [0, 360],
        "spectrumMatrix": [
          [0.01, 0.02, 0.03, /* ... */],
          [0.02, 0.04, 0.05, /* ... */],
          // ... 更多数据行
        ],
        "frequencyResolution": 0.01,
        "directionResolution": 10
      },
      "timeSeries": {
        "time": [0, 1, 2, /* ... */],
        "surfaceElevation": [0.5, 0.6, 0.4, /* ... */]
      },
      "visualization": {
        "spectrumImageUrl": "/api/v1/files/images/spectrum_task001.png",
        "timeSeriesImageUrl": "/api/v1/files/images/timeseries_task001.png",
        "animationUrl": "/api/v1/files/animations/surface_task001.mp4"
      }
    },
    "timestamp": 1634567890123
  }
  ```

## 4. 系统管理模块

### 4.1 获取系统参数

- **URL**: `/system/config`
- **方法**: GET
- **描述**: 获取系统配置参数
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "version": "1.0.0",
      "dataStorage": {
        "totalSpace": 10995116277760,
        "usedSpace": 3298534883328,
        "availableSpace": 7696581394432
      },
      "processingUnits": {
        "maxConcurrentTasks": 10,
        "currentRunningTasks": 3
      },
      "dataRetentionPolicy": {
        "temporaryFilesRetention": 7,
        "analysisResultsRetention": 30,
        "userUploadedDataRetention": 365
      },
      "supportedDataFormats": ["netCDF", "GeoTIFF", "CSV", "JSON"],
      "supportedAnalysisTypes": ["temperature-salinity", "sea-level", "wave-inversion"]
    },
    "timestamp": 1634567890123
  }
  ```

### 4.2 获取系统日志

- **URL**: `/system/logs`
- **方法**: GET
- **描述**: 获取系统日志
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**:
  - `page`: 页码，默认1
  - `size`: 每页条数，默认20
  - `level`: 日志级别，可选 ["info", "warning", "error"]
  - `startDate`: 开始日期
  - `endDate`: 结束日期
  - `module`: 模块名称
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "total": 1243,
      "page": 1,
      "size": 20,
      "logs": [
        {
          "id": "log12345",
          "timestamp": "2023-10-15T14:30:45Z",
          "level": "error",
          "module": "dataProcessor",
          "message": "处理数据集DS002时发生错误：文件格式不支持",
          "details": "不支持的文件格式: .xyz, 支持的格式: [netCDF, GeoTIFF, CSV, JSON]",
          "userId": "u789",
          "ip": "192.168.1.105"
        },
        // ... 更多日志
      ]
    },
    "timestamp": 1634567890123
  }
  ```

### 4.3 获取访问统计

- **URL**: `/system/statistics`
- **方法**: GET
- **描述**: 获取系统访问和使用统计
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**:
  - `type`: 统计类型，可选 ["user", "data", "analysis", "resource"]
  - `period`: 统计周期，可选 ["day", "week", "month", "year"]
  - `startDate`: 开始日期
  - `endDate`: 结束日期
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "type": "data",
      "period": "month",
      "timeRange": {
        "start": "2023-09-01T00:00:00Z",
        "end": "2023-09-30T23:59:59Z"
      },
      "statistics": {
        "totalDataUploaded": 12500000000,
        "totalDataDownloaded": 78900000000,
        "datasetAccess": [
          {
            "datasetId": "ds001",
            "name": "北太平洋温盐数据集2023",
            "accessCount": 456,
            "downloadCount": 89
          },
          // ... 更多数据集统计
        ],
        "topDataTypes": [
          {
            "type": "temperature",
            "accessCount": 1245,
            "percentage": 42.5
          },
          {
            "type": "wave",
            "accessCount": 876,
            "percentage": 29.9
          },
          // ... 更多数据类型统计
        ],
        "dailyAccess": [
          {
            "date": "2023-09-01",
            "count": 87
          },
          {
            "date": "2023-09-02",
            "count": 93
          },
          // ... 更多日期统计
        ]
      }
    },
    "timestamp": 1634567890123
  }
  ```

## 5. 文件管理模块

### 5.1 上传文件

- **URL**: `/files/upload`
- **方法**: POST
- **描述**: 上传通用文件（图片、文档等）
- **请求头**: 
  - `Authorization: Bearer {token}`
  - `Content-Type: multipart/form-data`
- **请求参数**:
  - `file`: 文件
  - `type`: 文件类型，可选 ["image", "document", "other"]
  - `description`: 文件描述
- **响应**:
  ```json
  {
    "code": 200,
    "message": "上传成功",
    "data": {
      "fileId": "file001",
      "name": "sampling_locations.png",
      "size": 1540000,
      "type": "image",
      "mimeType": "image/png",
      "url": "/api/v1/files/file001",
      "uploadTime": "2023-10-15T16:20:00Z"
    },
    "timestamp": 1634567890123
  }
  ```

### 5.2 获取文件

- **URL**: `/files/{fileId}`
- **方法**: GET
- **描述**: 获取文件
- **请求头**: `Authorization: Bearer {token}`
- **响应**: 文件流

## 6. 通知模块

### 6.1 获取用户通知

- **URL**: `/notifications`
- **方法**: GET
- **描述**: 获取用户通知
- **请求头**: `Authorization: Bearer {token}`
- **请求参数**:
  - `page`: 页码，默认1
  - `size`: 每页条数，默认20
  - `unreadOnly`: 是否只获取未读通知，默认false
- **响应**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "total": 12,
      "unreadCount": 3,
      "page": 1,
      "size": 20,
      "notifications": [
        {
          "id": "n001",
          "type": "task_complete",
          "title": "数据分析任务完成",
          "content": "您的海浪反演任务 task001 已完成处理",
          "createdAt": "2023-10-15T16:00:00Z",
          "read": false,
          "data": {
            "taskId": "task001",
            "resultUrl": "/analysis/wave-inversion/results/task001"
          }
        },
        // ... 更多通知
      ]
    },
    "timestamp": 1634567890123
  }
  ```

### 6.2 标记通知为已读

- **URL**: `/notifications/{notificationId}/read`
- **方法**: PUT
- **描述**: 标记通知为已读
- **请求头**: `Authorization: Bearer {token}`
- **响应**:
  ```json
  {
    "code": 200,
    "message": "标记成功",
    "data": {
      "id": "n001",
      "read": true,
      "readAt": "2023-10-15T16:30:00Z"
    },
    "timestamp": 1634567890123
  }
  ```

## 7. 权限与角色

### 可用角色
- `admin`: 系统管理员
- `researcher`: 研究人员
- `student`: 学生用户
- `guest`: 访客

### 权限列表
- `user:read`: 读取用户信息
- `user:write`: 修改用户信息
- `data:read`: 读取数据集
- `data:write`: 上传/修改数据集
- `data:delete`: 删除数据集
- `analysis:use`: 使用分析功能
- `system:admin`: 系统管理功能权限

### 角色权限映射
- `admin`: 所有权限
- `researcher`: `user:read`, `data:read`, `data:write`, `analysis:use`
- `student`: `user:read`, `data:read`, `analysis:use`
- `guest`: `user:read`, `data:read`(部分) 