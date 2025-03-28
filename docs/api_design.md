# 连锁管理系统API接口设计文档

## 1. API概述

本文档详细描述了连锁管理系统的REST API接口设计，包括认证接口、会员管理、商品管理、订单管理等核心功能模块的接口规范。

### 1.1 接口规范

- 采用RESTful API设计风格
- 使用HTTP状态码表示请求结果
- 所有接口均采用JSON格式进行数据交换
- 接口返回统一格式：`{code: 状态码, message: "消息提示", data: 返回数据}`
- 所有需要鉴权的接口需在请求头中携带token: `Authorization: Bearer {token}`

### 1.2 接口域名及版本

- 基础URL: `/api/v1/`
- 测试环境: `https://test-api.example.com/api/v1/`
- 生产环境: `https://api.example.com/api/v1/`

### 1.3 通用错误码

- 200: 请求成功
- 400: 请求参数错误
- 401: 未授权，请先登录
- 403: 权限不足
- 404: 资源不存在
- 409: 资源冲突
- 429: 请求过于频繁
- 500: 服务器内部错误

## 2. API接口设计

### 2.1 会员管理接口
```

POST /api/v1/members # 创建会员
GET /api/v1/members # 获取会员列表(支持分页、筛选)
GET /api/v1/members/:id # 获取会员详情
PUT /api/v1/members/:id # 更新会员信息
DELETE /api/v1/members/:id # 逻辑删除会员

POST /api/v1/members/:id/recharge # 会员充值
POST /api/v1/members/:id/consume # 会员消费
GET /api/v1/members/:id/transactions # 获取会员交易记录(支持分页、筛选)

GET /api/v1/members/export # 导出会员数据(支持筛选条件)
POST /api/v1/members/import # 导入会员数据
GET /api/v1/members/statistics # 获取会员统计数据

GET /api/v1/member-tags # 获取所有会员标签
POST /api/v1/member-tags # 创建会员标签
DELETE /api/v1/member-tags/:id # 删除会员标签
POST /api/v1/members/:id/tags # 为会员添加标签
DELETE /api/v1/members/:id/tags/:tagId # 移除会员标签
```

#### 2.1.1 创建会员
- **URL**: `/api/v1/members`
- **方法**: POST
- **描述**: 创建新会员
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "name": "string",           // 姓名，必填
    "gender": number,           // 性别：0未知，1男，2女
    "phone": "string",          // 手机号，必填
    "card_id": "string",        // 会员卡号，必填
    "store_id": number,         // 所属店铺ID，必填
    "balance": number,          // 初始余额，可选，默认0
    "gift_balance": number,     // 初始赠送余额，可选，默认0
    "remark": "string",         // 备注，可选
    "tags": [number]            // 会员标签ID数组，可选
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "创建会员成功",
    "data": {
      "id": number,             // 会员ID
      "name": "string",         // 姓名
      "gender": number,         // 性别
      "phone": "string",        // 手机号
      "card_id": "string",      // 会员卡号
      "store_id": number,       // 所属店铺ID
      "store_name": "string",   // 所属店铺名称
      "balance": number,        // 余额
      "gift_balance": number,   // 赠送余额
      "status": number,         // 状态：1正常，0删除
      "created_at": "string",   // 创建时间
      "updated_at": "string"    // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 409: 会员卡号或手机号已存在
  - 500: 服务器内部错误

#### 2.1.2 获取会员列表
- **URL**: `/api/v1/members`
- **方法**: GET
- **描述**: 获取会员列表，支持分页和筛选
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  page: number           // 页码，从1开始，默认1
  page_size: number      // 每页数量，默认20
  keyword: string        // 关键词搜索(姓名/手机号/卡号)
  store_id: number       // 按店铺ID筛选
  gender: number         // 按性别筛选
  status: number         // 按状态筛选
  min_balance: number    // 最小余额筛选
  max_balance: number    // 最大余额筛选
  tag_id: number         // 按标签ID筛选
  start_date: string     // 创建开始日期
  end_date: string       // 创建结束日期
  sort_field: string     // 排序字段
  sort_order: string     // 排序方式(asc/desc)
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取会员列表成功",
    "data": {
      "total": number,          // 总记录数
      "page": number,           // 当前页码
      "page_size": number,      // 每页数量
      "total_pages": number,    // 总页数
      "items": [
        {
          "id": number,           // 会员ID
          "name": "string",       // 姓名
          "gender": number,       // 性别
          "phone": "string",      // 手机号
          "card_id": "string",    // 会员卡号
          "store_id": number,     // 所属店铺ID
          "store_name": "string", // 所属店铺名称
          "balance": number,      // 余额
          "gift_balance": number, // 赠送余额
          "total_recharge": number, // 总充值金额
          "total_consume": number,  // 总消费金额
          "last_active_time": "string", // 最后活跃时间
          "status": number,       // 状态
          "created_at": "string", // 创建时间
          "tags": [               // 会员标签
            {
              "id": number,
              "name": "string"
            }
          ]
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.1.3 获取会员详情
- **URL**: `/api/v1/members/:id`
- **方法**: GET
- **描述**: 获取会员详细信息
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取会员详情成功",
    "data": {
      "id": number,             // 会员ID
      "name": "string",         // 姓名
      "gender": number,         // 性别
      "phone": "string",        // 手机号
      "card_id": "string",      // 会员卡号
      "store_id": number,       // 所属店铺ID
      "store_name": "string",   // 所属店铺名称
      "balance": number,        // 余额
      "gift_balance": number,   // 赠送余额
      "total_recharge": number, // 总充值金额
      "total_consume": number,  // 总消费金额
      "first_active_time": "string", // 首次活跃时间
      "last_active_time": "string",  // 最后活跃时间
      "status": number,         // 状态
      "remark": "string",       // 备注
      "created_at": "string",   // 创建时间
      "updated_at": "string",   // 更新时间
      "tags": [                 // 会员标签
        {
          "id": number,
          "name": "string"
        }
      ],
      "recent_transactions": [  // 最近交易记录
        {
          "id": number,
          "type": "string",    // 交易类型(RECHARGE/CONSUME/GIFT)
          "amount": number,    // 交易金额
          "created_at": "string" // 交易时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 500: 服务器内部错误

#### 2.1.4 更新会员信息
- **URL**: `/api/v1/members/:id`
- **方法**: PUT
- **描述**: 更新会员信息
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "name": "string",         // 姓名，可选
    "gender": number,         // 性别，可选
    "phone": "string",        // 手机号，可选
    "card_id": "string",      // 会员卡号，可选
    "store_id": number,       // 所属店铺ID，可选
    "status": number,         // 状态，可选
    "remark": "string"        // 备注，可选
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "更新会员信息成功",
    "data": {
      "id": number,           // 会员ID
      "name": "string",       // 姓名
      "gender": number,       // 性别
      "phone": "string",      // 手机号
      "card_id": "string",    // 会员卡号
      "store_id": number,     // 所属店铺ID
      "balance": number,      // 余额
      "gift_balance": number, // 赠送余额
      "status": number,       // 状态
      "remark": "string",     // 备注
      "updated_at": "string"  // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 409: 会员卡号或手机号已存在
  - 500: 服务器内部错误

#### 2.1.5 删除会员
- **URL**: `/api/v1/members/:id`
- **方法**: DELETE
- **描述**: 逻辑删除会员（设置状态为删除）
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "删除会员成功",
    "data": null
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 500: 服务器内部错误

#### 2.1.6 会员充值
- **URL**: `/api/v1/members/:id/recharge`
- **方法**: POST
- **描述**: 为会员充值
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "amount": number,          // 充值金额，必填
    "gift_amount": number,     // 赠送金额，可选，默认0
    "payment_method": "string", // 支付方式，必填(CASH/WECHAT/ALIPAY/BANKCARD/OTHER)
    "promotion_id": number,     // 优惠活动ID，可选
    "remark": "string",         // 备注，可选
    "store_id": number          // 充值店铺ID，必填
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "会员充值成功",
    "data": {
      "id": number,             // 交易记录ID
      "member_id": number,      // 会员ID
      "type": "RECHARGE",       // 交易类型
      "amount": number,         // 充值金额
      "gift_amount": number,    // 赠送金额
      "balance": number,        // 充值后余额
      "gift_balance": number,   // 充值后赠送余额
      "payment_method": "string", // 支付方式
      "remark": "string",       // 备注
      "store_id": number,       // 充值店铺ID
      "store_name": "string",   // 充值店铺名称
      "operator_id": number,    // 操作人ID
      "operator_name": "string", // 操作人姓名
      "created_at": "string"    // 充值时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 500: 服务器内部错误

#### 2.1.7 会员消费
- **URL**: `/api/v1/members/:id/consume`
- **方法**: POST
- **描述**: 记录会员消费
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "amount": number,          // 消费金额，必填
    "payment_method": "string", // 支付方式(BALANCE/CASH/WECHAT/ALIPAY/BANKCARD/OTHER)，必填
    "item_type": "string",      // 消费项目类型(PRODUCT/VENUE/COACHING/SERVICE/OTHER)，可选
    "item_id": number,          // 消费项目ID，可选
    "description": "string",    // 消费描述，可选
    "remark": "string",         // 备注，可选
    "store_id": number          // 消费店铺ID，必填
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "会员消费成功",
    "data": {
      "id": number,             // 交易记录ID
      "member_id": number,      // 会员ID
      "type": "CONSUME",        // 交易类型
      "amount": number,         // 消费金额
      "balance": number,        // 消费后余额
      "gift_balance": number,   // 消费后赠送余额
      "payment_method": "string", // 支付方式
      "item_type": "string",    // 消费项目类型
      "item_id": number,        // 消费项目ID
      "description": "string",  // 消费描述
      "remark": "string",       // 备注
      "store_id": number,       // 消费店铺ID
      "store_name": "string",   // 消费店铺名称
      "operator_id": number,    // 操作人ID
      "operator_name": "string", // 操作人姓名
      "created_at": "string"    // 消费时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 409: 余额不足
  - 500: 服务器内部错误

#### 2.1.8 获取会员交易记录
- **URL**: `/api/v1/members/:id/transactions`
- **方法**: GET
- **描述**: 获取会员交易记录，支持分页和筛选
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  page: number           // 页码，从1开始，默认1
  page_size: number      // 每页数量，默认20
  type: string           // 交易类型(RECHARGE/CONSUME/GIFT)
  store_id: number       // 按店铺ID筛选
  start_date: string     // 交易开始日期
  end_date: string       // 交易结束日期
  min_amount: number     // 最小金额
  max_amount: number     // 最大金额
  payment_method: string // 支付方式
  sort_field: string     // 排序字段，默认created_at
  sort_order: string     // 排序方式(asc/desc)，默认desc
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取交易记录成功",
    "data": {
      "total": number,          // 总记录数
      "page": number,           // 当前页码
      "page_size": number,      // 每页数量
      "total_pages": number,    // 总页数
      "member": {               // 会员基本信息
        "id": number,
        "name": "string",
        "card_id": "string",
        "balance": number,
        "gift_balance": number
      },
      "items": [
        {
          "id": number,             // 交易记录ID
          "type": "string",         // 交易类型(RECHARGE/CONSUME/GIFT)
          "amount": number,         // 交易金额
          "balance": number,        // 交易后余额
          "gift_balance": number,   // 交易后赠送余额
          "payment_method": "string", // 支付方式
          "item_type": "string",    // 消费项目类型(仅消费时有)
          "item_id": number,        // 消费项目ID(仅消费时有)
          "description": "string",  // 交易描述
          "remark": "string",       // 备注
          "store_id": number,       // 店铺ID
          "store_name": "string",   // 店铺名称
          "operator_id": number,    // 操作人ID
          "operator_name": "string", // 操作人姓名
          "created_at": "string"    // 交易时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 500: 服务器内部错误

#### 2.1.9 导出会员数据
- **URL**: `/api/v1/members/export`
- **方法**: GET
- **描述**: 导出会员数据为Excel或CSV格式，支持筛选条件
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  format: string         // 导出格式(excel/csv)，默认excel
  keyword: string        // 关键词搜索(姓名/手机号/卡号)
  store_id: number       // 按店铺ID筛选
  gender: number         // 按性别筛选
  status: number         // 按状态筛选
  min_balance: number    // 最小余额筛选
  max_balance: number    // 最大余额筛选
  tag_id: number         // 按标签ID筛选
  start_date: string     // 创建开始日期
  end_date: string       // 创建结束日期
  ```
- **响应数据**: 文件下载流
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.1.10 导入会员数据
- **URL**: `/api/v1/members/import`
- **方法**: POST
- **描述**: 从Excel或CSV文件导入会员数据
- **请求头**: Authorization: Bearer {token}
- **请求参数**: multipart/form-data
  ```
  file: File          // 上传的文件
  store_id: number    // 导入的会员所属店铺ID
  update_existing: boolean // 是否更新已存在的会员，默认false
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "导入会员数据成功",
    "data": {
      "total": number,       // 总记录数
      "success": number,     // 成功导入数
      "fail": number,        // 失败数
      "errors": [            // 错误记录
        {
          "row": number,     // 行号
          "reason": "string" // 失败原因
        }
      ],
      "download_error_url": "string" // 下载错误记录文件URL
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误或文件格式不正确
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 413: 文件过大
  - 500: 服务器内部错误

#### 2.1.11 获取会员统计数据
- **URL**: `/api/v1/members/statistics`
- **方法**: GET
- **描述**: 获取会员相关统计数据
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  store_id: number     // 店铺ID，不传则查询所有店铺
  period: string       // 统计周期(day/week/month/year)，默认month
  start_date: string   // 开始日期
  end_date: string     // 结束日期
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取会员统计数据成功",
    "data": {
      "summary": {
        "total_members": number,        // 会员总数
        "active_members": number,       // 活跃会员数
        "new_members_today": number,    // 今日新增会员
        "new_members_week": number,     // 本周新增会员
        "new_members_month": number,    // 本月新增会员
        "total_balance": number,        // 会员总余额
        "total_gift_balance": number,   // 会员总赠送余额
        "average_balance": number,      // 平均余额
        "recharge_today": number,       // 今日充值金额
        "consume_today": number         // 今日消费金额
      },
      "trend": {                        // 趋势数据
        "dates": ["string"],            // 日期列表
        "new_members": [number],        // 新增会员数列表
        "recharge_amounts": [number],   // 充值金额列表
        "consume_amounts": [number]     // 消费金额列表
      },
      "gender_distribution": [          // 性别分布
        {
          "gender": number,
          "count": number,
          "percentage": number
        }
      ],
      "store_distribution": [           // 店铺分布
        {
          "store_id": number,
          "store_name": "string",
          "count": number,
          "percentage": number
        }
      ],
      "balance_distribution": [         // 余额分布
        {
          "range": "string",            // 余额区间
          "count": number,
          "percentage": number
        }
      ],
      "tag_distribution": [             // 标签分布
        {
          "tag_id": number,
          "tag_name": "string",
          "count": number,
          "percentage": number
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.1.12 获取所有会员标签
- **URL**: `/api/v1/member-tags`
- **方法**: GET
- **描述**: 获取所有会员标签
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取会员标签成功",
    "data": [
      {
        "id": number,
        "name": "string",       // 标签名称
        "color": "string",      // 标签颜色
        "member_count": number, // 使用该标签的会员数
        "created_at": "string"  // 创建时间
      }
    ]
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 500: 服务器内部错误

#### 2.1.13 创建会员标签
- **URL**: `/api/v1/member-tags`
- **方法**: POST
- **描述**: 创建新的会员标签
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "name": "string",  // 标签名称，必填
    "color": "string"  // 标签颜色，可选，默认随机颜色
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "创建会员标签成功",
    "data": {
      "id": number,
      "name": "string",
      "color": "string",
      "created_at": "string"
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 409: 标签名称已存在
  - 500: 服务器内部错误

#### 2.1.14 删除会员标签
- **URL**: `/api/v1/member-tags/:id`
- **方法**: DELETE
- **描述**: 删除会员标签
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "删除会员标签成功",
    "data": null
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 标签不存在
  - 500: 服务器内部错误

#### 2.1.15 为会员添加标签
- **URL**: `/api/v1/members/:id/tags`
- **方法**: POST
- **描述**: 为会员添加标签
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "tag_ids": [number]  // 标签ID数组，必填
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "添加会员标签成功",
    "data": {
      "member_id": number,
      "tags": [
        {
          "id": number,
          "name": "string",
          "color": "string"
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员或标签不存在
  - 500: 服务器内部错误

#### 2.1.16 移除会员标签
- **URL**: `/api/v1/members/:id/tags/:tagId`
- **方法**: DELETE
- **描述**: 移除会员的标签
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "移除会员标签成功",
    "data": null
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员或标签不存在
  - 500: 服务器内部错误

### 2.2 数据回收站接口
```
GET    /api/v1/recycle/members       # 获取已删除会员列表
POST   /api/v1/recycle/members/:id/restore  # 恢复已删除会员
DELETE /api/v1/recycle/members/:id   # 彻底删除会员
POST   /api/v1/recycle/members/batch-restore  # 批量恢复会员
DELETE /api/v1/recycle/members/batch-delete   # 批量彻底删除会员

GET    /api/v1/recycle/products      # 获取已删除商品列表
POST   /api/v1/recycle/products/:id/restore  # 恢复已删除商品
DELETE /api/v1/recycle/products/:id  # 彻底删除商品
POST   /api/v1/recycle/products/batch-restore  # 批量恢复商品
DELETE /api/v1/recycle/products/batch-delete   # 批量彻底删除商品

GET    /api/v1/recycle/activities    # 获取已删除活动列表
POST   /api/v1/recycle/activities/:id/restore  # 恢复已删除活动
DELETE /api/v1/recycle/activities/:id  # 彻底删除活动
POST   /api/v1/recycle/activities/batch-restore  # 批量恢复活动
DELETE /api/v1/recycle/activities/batch-delete   # 批量彻底删除活动
```

#### 2.2.1 获取已删除会员列表
- **URL**: `/api/v1/recycle/members`
- **方法**: GET
- **描述**: 获取已删除的会员列表，支持分页和筛选
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  page: number           // 页码，从1开始，默认1
  page_size: number      // 每页数量，默认20
  keyword: string        // 关键词搜索(姓名/手机号/卡号)
  store_id: number       // 按店铺ID筛选
  delete_start_date: string // 删除开始日期
  delete_end_date: string   // 删除结束日期
  operator_id: number    // 操作人ID
  sort_field: string     // 排序字段，默认delete_time
  sort_order: string     // 排序方式(asc/desc)，默认desc
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取已删除会员列表成功",
    "data": {
      "total": number,          // 总记录数
      "page": number,           // 当前页码
      "page_size": number,      // 每页数量
      "total_pages": number,    // 总页数
      "items": [
        {
          "id": number,             // 记录ID
          "original_id": number,    // 原会员ID
          "name": "string",         // 姓名
          "phone": "string",        // 手机号
          "card_id": "string",      // 会员卡号
          "store_id": number,       // 所属店铺ID
          "store_name": "string",   // 所属店铺名称
          "balance": number,        // 删除前余额
          "gift_balance": number,   // 删除前赠送余额
          "delete_time": "string",  // 删除时间
          "delete_reason": "string", // 删除原因
          "delete_operator_id": number, // 删除操作人ID
          "delete_operator_name": "string", // 删除操作人姓名
          "expiry_time": "string",  // 过期时间（超过将自动彻底删除）
          "is_recoverable": boolean // 是否可恢复
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.2.2 恢复已删除会员
- **URL**: `/api/v1/recycle/members/:id/restore`
- **方法**: POST
- **描述**: 恢复已删除的会员
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "恢复会员成功",
    "data": {
      "member_id": number,      // 恢复后的会员ID
      "name": "string",         // 姓名
      "phone": "string",        // 手机号
      "card_id": "string",      // 会员卡号
      "recovery_time": "string" // 恢复时间
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 已删除记录不存在
  - 409: 恢复冲突（卡号或手机号已被其他会员使用）
  - 410: 记录已过期无法恢复
  - 500: 服务器内部错误

#### 2.2.3 彻底删除会员
- **URL**: `/api/v1/recycle/members/:id`
- **方法**: DELETE
- **描述**: 彻底删除已删除的会员记录（不可恢复）
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "彻底删除会员成功",
    "data": null
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 已删除记录不存在
  - 500: 服务器内部错误

#### 2.2.4 批量恢复会员
- **URL**: `/api/v1/recycle/members/batch-restore`
- **方法**: POST
- **描述**: 批量恢复已删除的会员
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "ids": [number]    // 要恢复的记录ID数组，必填
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "批量恢复会员成功",
    "data": {
      "total": number,       // 总操作数
      "success": number,     // 成功数
      "failed": number,      // 失败数
      "results": [
        {
          "id": number,            // 记录ID
          "original_id": number,   // 原会员ID
          "success": boolean,      // 是否恢复成功
          "new_id": number,        // 恢复后的ID（成功时）
          "error": "string"        // 失败原因（失败时）
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.2.5 批量彻底删除会员
- **URL**: `/api/v1/recycle/members/batch-delete`
- **方法**: DELETE
- **描述**: 批量彻底删除已删除的会员记录（不可恢复）
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "ids": [number]    // 要彻底删除的记录ID数组，必填
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "批量彻底删除会员成功",
    "data": {
      "total": number,       // 总操作数
      "success": number,     // 成功数
      "failed": number       // 失败数
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.2.6 获取已删除商品列表
- **URL**: `/api/v1/recycle/products`
- **方法**: GET
- **描述**: 获取已删除的商品列表，支持分页和筛选
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  page: number           // 页码，从1开始，默认1
  page_size: number      // 每页数量，默认20
  keyword: string        // 关键词搜索(商品名称/编码)
  category_id: number    // 按分类ID筛选
  delete_start_date: string // 删除开始日期
  delete_end_date: string   // 删除结束日期
  operator_id: number    // 操作人ID
  sort_field: string     // 排序字段，默认delete_time
  sort_order: string     // 排序方式(asc/desc)，默认desc
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取已删除商品列表成功",
    "data": {
      "total": number,          // 总记录数
      "page": number,           // 当前页码
      "page_size": number,      // 每页数量
      "total_pages": number,    // 总页数
      "items": [
        {
          "id": number,             // 记录ID
          "original_id": number,    // 原商品ID
          "code": "string",         // 商品编码
          "name": "string",         // 商品名称
          "category_id": number,    // 分类ID
          "category_name": "string", // 分类名称
          "retail_price": number,   // 零售价
          "stock_quantity": number, // 删除前库存数量
          "delete_time": "string",  // 删除时间
          "delete_reason": "string", // 删除原因
          "delete_operator_id": number, // 删除操作人ID
          "delete_operator_name": "string", // 删除操作人姓名
          "expiry_time": "string",  // 过期时间（超过将自动彻底删除）
          "is_recoverable": boolean // 是否可恢复
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.2.7 恢复已删除商品
- **URL**: `/api/v1/recycle/products/:id/restore`
- **方法**: POST
- **描述**: 恢复已删除的商品
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "恢复商品成功",
    "data": {
      "product_id": number,     // 恢复后的商品ID
      "code": "string",         // 商品编码
      "name": "string",         // 商品名称
      "recovery_time": "string" // 恢复时间
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 已删除记录不存在
  - 409: 恢复冲突（商品编码已被使用）
  - 410: 记录已过期无法恢复
  - 500: 服务器内部错误

#### 2.2.8 彻底删除商品
- **URL**: `/api/v1/recycle/products/:id`
- **方法**: DELETE
- **描述**: 彻底删除已删除的商品记录（不可恢复）
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "彻底删除商品成功",
    "data": null
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 已删除记录不存在
  - 500: 服务器内部错误

#### 2.2.9 批量恢复商品
- **URL**: `/api/v1/recycle/products/batch-restore`
- **方法**: POST
- **描述**: 批量恢复已删除的商品
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "ids": [number]    // 要恢复的记录ID数组，必填
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "批量恢复商品成功",
    "data": {
      "total": number,       // 总操作数
      "success": number,     // 成功数
      "failed": number,      // 失败数
      "results": [
        {
          "id": number,            // 记录ID
          "original_id": number,   // 原商品ID
          "success": boolean,      // 是否恢复成功
          "new_id": number,        // 恢复后的ID（成功时）
          "error": "string"        // 失败原因（失败时）
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.2.10 批量彻底删除商品
- **URL**: `/api/v1/recycle/products/batch-delete`
- **方法**: DELETE
- **描述**: 批量彻底删除已删除的商品记录（不可恢复）
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "ids": [number]    // 要彻底删除的记录ID数组，必填
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "批量彻底删除商品成功",
    "data": {
      "total": number,       // 总操作数
      "success": number,     // 成功数
      "failed": number       // 失败数
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.2.11 获取已删除活动列表
- **URL**: `/api/v1/recycle/activities`
- **方法**: GET
- **描述**: 获取已删除的活动列表，支持分页和筛选
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  page: number           // 页码，从1开始，默认1
  page_size: number      // 每页数量，默认20
  keyword: string        // 关键词搜索(活动名称/编码)
  category: string       // 按活动类别筛选
  delete_start_date: string // 删除开始日期
  delete_end_date: string   // 删除结束日期
  operator_id: number    // 操作人ID
  sort_field: string     // 排序字段，默认delete_time
  sort_order: string     // 排序方式(asc/desc)，默认desc
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取已删除活动列表成功",
    "data": {
      "total": number,          // 总记录数
      "page": number,           // 当前页码
      "page_size": number,      // 每页数量
      "total_pages": number,    // 总页数
      "items": [
        {
          "id": number,             // 记录ID
          "original_id": number,    // 原活动ID
          "code": "string",         // 活动编码
          "title": "string",        // 活动标题
          "category": "string",     // 活动类别
          "start_time": "string",   // 活动开始时间
          "end_time": "string",     // 活动结束时间
          "delete_time": "string",  // 删除时间
          "delete_reason": "string", // 删除原因
          "delete_operator_id": number, // 删除操作人ID
          "delete_operator_name": "string", // 删除操作人姓名
          "expiry_time": "string",  // 过期时间（超过将自动彻底删除）
          "is_recoverable": boolean // 是否可恢复
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.2.12 恢复已删除活动
- **URL**: `/api/v1/recycle/activities/:id/restore`
- **方法**: POST
- **描述**: 恢复已删除的活动
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "恢复活动成功",
    "data": {
      "activity_id": number,    // 恢复后的活动ID
      "code": "string",         // 活动编码
      "title": "string",        // 活动标题
      "recovery_time": "string" // 恢复时间
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 已删除记录不存在
  - 409: 恢复冲突（活动编码已被使用）
  - 410: 记录已过期无法恢复
  - 500: 服务器内部错误

#### 2.2.13 彻底删除活动
- **URL**: `/api/v1/recycle/activities/:id`
- **方法**: DELETE
- **描述**: 彻底删除已删除的活动记录（不可恢复）
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "彻底删除活动成功",
    "data": null
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 已删除记录不存在
  - 500: 服务器内部错误

#### 2.2.14 批量恢复活动
- **URL**: `/api/v1/recycle/activities/batch-restore`
- **方法**: POST
- **描述**: 批量恢复已删除的活动
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "ids": [number]    // 要恢复的记录ID数组，必填
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "批量恢复活动成功",
    "data": {
      "total": number,       // 总操作数
      "success": number,     // 成功数
      "failed": number,      // 失败数
      "results": [
        {
          "id": number,            // 记录ID
          "original_id": number,   // 原活动ID
          "success": boolean,      // 是否恢复成功
          "new_id": number,        // 恢复后的ID（成功时）
          "error": "string"        // 失败原因（失败时）
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.2.15 批量彻底删除活动
- **URL**: `/api/v1/recycle/activities/batch-delete`
- **方法**: DELETE
- **描述**: 批量彻底删除已删除的活动记录（不可恢复）
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "ids": [number]    // 要彻底删除的记录ID数组，必填
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "批量彻底删除活动成功",
    "data": {
      "total": number,       // 总操作数
      "success": number,     // 成功数
      "failed": number       // 失败数
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

### 2.3 商品管理接口
```
POST   /api/v1/products              # 创建商品
GET    /api/v1/products              # 获取商品列表
GET    /api/v1/products/:id          # 获取商品详情
PUT    /api/v1/products/:id          # 更新商品信息
DELETE /api/v1/products/:id          # 删除商品

POST   /api/v1/products/:id/inventory # 获取商品库存记录
POST   /api/v1/products/:id/inventory/in # 商品入库
POST   /api/v1/products/:id/inventory/out # 商品出库
GET    /api/v1/products/inventory/alert # 获取库存预警商品列表

GET    /api/v1/products/export       # 导出商品数据
POST   /api/v1/products/import       # 导入商品数据
GET    /api/v1/products/statistics   # 获取商品销售统计
```

#### 2.3.1 创建商品

- **URL**: `/api/v1/products`
- **方法**: POST
- **描述**: 创建一个新的商品信息
- **请求头**:
  - `Authorization: Bearer {token}`
  - `Content-Type: application/json`
- **请求参数**:
  ```json
  {
    "name": string,                // 商品名称（必填）
    "category_id": number,         // 商品类别ID（必填）
    "description": string,         // 商品描述
    "specifications": string,      // 规格参数
    "price": number,               // 售价（必填）
    "cost_price": number,          // 成本价
    "market_price": number,        // 市场价
    "stock": number,               // 初始库存数量（必填）
    "min_stock": number,           // 最小库存预警值
    "status": number,              // 状态：0-下架，1-上架（默认1）
    "images": string[],            // 商品图片URL数组
    "main_image": string,          // 主图URL
    "attributes": {                // 商品属性
      "color": string,             // 颜色
      "size": string,              // 尺寸
      "weight": number,            // 重量
      [key: string]: any           // 其他自定义属性
    },
    "tags": string[],              // 商品标签
    "barcode": string              // 商品条形码
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "创建商品成功",
    "data": {
      "id": number,                // 商品ID
      "name": string,              // 商品名称
      "category_id": number,       // 商品类别ID
      "description": string,       // 商品描述
      "specifications": string,    // 规格参数
      "price": number,             // 售价
      "cost_price": number,        // 成本价
      "market_price": number,      // 市场价
      "stock": number,             // 库存数量
      "min_stock": number,         // 最小库存预警值
      "status": number,            // 状态：0-下架，1-上架
      "images": string[],          // 商品图片URL数组
      "main_image": string,        // 主图URL
      "attributes": object,        // 商品属性
      "tags": string[],            // 商品标签
      "barcode": string,           // 商品条形码
      "created_at": string,        // 创建时间
      "updated_at": string         // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 409: 商品名称已存在
  - 500: 服务器内部错误

#### 2.3.2 获取商品列表

- **URL**: `/api/v1/products`
- **方法**: GET
- **描述**: 获取商品列表，支持分页、排序和筛选
- **请求头**:
  - `Authorization: Bearer {token}`
- **请求参数**:
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `keyword`: string, 搜索关键词（商品名称、描述）
  - `category_id`: number, 按类别筛选
  - `status`: number, 商品状态：0-下架，1-上架
  - `min_price`: number, 最低价格
  - `max_price`: number, 最高价格
  - `sort_by`: string, 排序字段：created_at/price/stock/sales
  - `sort_order`: string, 排序方式：asc/desc
  - `tag`: string, 按标签筛选
  - `stock_alert`: boolean, 是否只显示库存预警商品
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取商品列表成功",
    "data": {
      "total": number,             // 总记录数
      "page": number,              // 当前页码
      "page_size": number,         // 每页数量
      "list": [
        {
          "id": number,            // 商品ID
          "name": string,          // 商品名称
          "category_id": number,   // 商品类别ID
          "category_name": string, // 类别名称
          "price": number,         // 售价
          "market_price": number,  // 市场价
          "stock": number,         // 库存数量
          "sales": number,         // 销量
          "status": number,        // 状态
          "main_image": string,    // 主图URL
          "tags": string[],        // 商品标签
          "created_at": string,    // 创建时间
          "updated_at": string     // 更新时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 500: 服务器内部错误

#### 2.3.3 获取商品详情

- **URL**: `/api/v1/products/:id`
- **方法**: GET
- **描述**: 获取指定ID商品的详细信息
- **请求头**:
  - `Authorization: Bearer {token}`
- **请求参数**:
  - `id`: number, 商品ID（路径参数）
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取商品详情成功",
    "data": {
      "id": number,                // 商品ID
      "name": string,              // 商品名称
      "category_id": number,       // 商品类别ID
      "category_name": string,     // 类别名称
      "description": string,       // 商品描述
      "specifications": string,    // 规格参数
      "price": number,             // 售价
      "cost_price": number,        // 成本价
      "market_price": number,      // 市场价
      "stock": number,             // 库存数量
      "min_stock": number,         // 最小库存预警值
      "sales": number,             // 销量
      "status": number,            // 状态：0-下架，1-上架
      "images": string[],          // 商品图片URL数组
      "main_image": string,        // 主图URL
      "attributes": object,        // 商品属性
      "tags": string[],            // 商品标签
      "barcode": string,           // 商品条形码
      "created_at": string,        // 创建时间
      "updated_at": string,        // 更新时间
      "inventory_records": [       // 最近的库存变动记录
        {
          "id": number,
          "type": string,          // 类型：in-入库，out-出库
          "quantity": number,      // 数量
          "operator": string,      // 操作人
          "created_at": string     // 操作时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 404: 商品不存在
  - 500: 服务器内部错误

#### 2.3.4 更新商品信息

- **URL**: `/api/v1/products/:id`
- **方法**: PUT
- **描述**: 更新指定ID商品的信息
- **请求头**:
  - `Authorization: Bearer {token}`
  - `Content-Type: application/json`
- **请求参数**:
  - `id`: number, 商品ID（路径参数）
  ```json
  {
    "name": string,                // 商品名称
    "category_id": number,         // 商品类别ID
    "description": string,         // 商品描述
    "specifications": string,      // 规格参数
    "price": number,               // 售价
    "cost_price": number,          // 成本价
    "market_price": number,        // 市场价
    "min_stock": number,           // 最小库存预警值
    "status": number,              // 状态：0-下架，1-上架
    "images": string[],            // 商品图片URL数组
    "main_image": string,          // 主图URL
    "attributes": object,          // 商品属性
    "tags": string[],              // 商品标签
    "barcode": string              // 商品条形码
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "更新商品成功",
    "data": {
      "id": number,                // 商品ID
      "name": string,              // 商品名称
      "category_id": number,       // 商品类别ID
      "description": string,       // 商品描述
      "specifications": string,    // 规格参数
      "price": number,             // 售价
      "cost_price": number,        // 成本价
      "market_price": number,      // 市场价
      "stock": number,             // 库存数量
      "min_stock": number,         // 最小库存预警值
      "status": number,            // 状态：0-下架，1-上架
      "images": string[],          // 商品图片URL数组
      "main_image": string,        // 主图URL
      "attributes": object,        // 商品属性
      "tags": string[],            // 商品标签
      "barcode": string,           // 商品条形码
      "updated_at": string         // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 商品不存在
  - 409: 商品名称已存在
  - 500: 服务器内部错误

#### 2.3.5 删除商品

- **URL**: `/api/v1/products/:id`
- **方法**: DELETE
- **描述**: 删除指定ID的商品（软删除）
- **请求头**:
  - `Authorization: Bearer {token}`
- **请求参数**:
  - `id`: number, 商品ID（路径参数）
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "删除商品成功",
    "data": {
      "id": number                 // 被删除的商品ID
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 商品不存在
  - 409: 商品已关联订单，无法删除
  - 500: 服务器内部错误

#### 2.3.6 获取商品库存记录

- **URL**: `/api/v1/products/:id/inventory`
- **方法**: POST
- **描述**: 获取指定商品的库存变动记录
- **请求头**:
  - `Authorization: Bearer {token}`
  - `Content-Type: application/json`
- **请求参数**:
  - `id`: number, 商品ID（路径参数）
  ```json
  {
    "page": number,                // 页码，默认1
    "page_size": number,           // 每页数量，默认20
    "start_date": string,          // 开始日期，格式YYYY-MM-DD
    "end_date": string,            // 结束日期，格式YYYY-MM-DD
    "type": string                 // 类型筛选：in-入库，out-出库，all-全部（默认）
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取库存记录成功",
    "data": {
      "total": number,             // 总记录数
      "page": number,              // 当前页码
      "page_size": number,         // 每页数量
      "list": [
        {
          "id": number,            // 记录ID
          "product_id": number,    // 商品ID
          "product_name": string,  // 商品名称
          "type": string,          // 类型：in-入库，out-出库
          "quantity": number,      // 数量
          "before_stock": number,  // 操作前库存
          "after_stock": number,   // 操作后库存
          "operator_id": number,   // 操作人ID
          "operator_name": string, // 操作人姓名
          "remark": string,        // 备注
          "created_at": string     // 操作时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 404: 商品不存在
  - 500: 服务器内部错误

#### 2.3.7 商品入库

- **URL**: `/api/v1/products/:id/inventory/in`
- **方法**: POST
- **描述**: 为指定商品增加库存
- **请求头**:
  - `Authorization: Bearer {token}`
  - `Content-Type: application/json`
- **请求参数**:
  - `id`: number, 商品ID（路径参数）
  ```json
  {
    "quantity": number,            // 入库数量（必填，正整数）
    "batch_number": string,        // 批次号
    "supplier_id": number,         // 供应商ID
    "purchase_price": number,      // 本次采购单价
    "production_date": string,     // 生产日期
    "expiry_date": string,         // 过期日期
    "remark": string               // 备注信息
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "商品入库成功",
    "data": {
      "product_id": number,        // 商品ID
      "product_name": string,      // 商品名称
      "quantity": number,          // 入库数量
      "before_stock": number,      // 入库前库存
      "after_stock": number,       // 入库后库存
      "batch_number": string,      // 批次号
      "operator_id": number,       // 操作人ID
      "operator_name": string,     // 操作人姓名
      "created_at": string         // 操作时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误，入库数量必须为正整数
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 商品不存在
  - 500: 服务器内部错误

#### 2.3.8 商品出库

- **URL**: `/api/v1/products/:id/inventory/out`
- **方法**: POST
- **描述**: 为指定商品减少库存
- **请求头**:
  - `Authorization: Bearer {token}`
  - `Content-Type: application/json`
- **请求参数**:
  - `id`: number, 商品ID（路径参数）
  ```json
  {
    "quantity": number,            // 出库数量（必填，正整数）
    "reason": string,              // 出库原因（必填）：sale-销售，loss-损耗，return-退货，other-其他
    "order_id": number,            // 关联订单ID（如有）
    "remark": string               // 备注信息
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "商品出库成功",
    "data": {
      "product_id": number,        // 商品ID
      "product_name": string,      // 商品名称
      "quantity": number,          // 出库数量
      "reason": string,            // 出库原因
      "before_stock": number,      // 出库前库存
      "after_stock": number,       // 出库后库存
      "operator_id": number,       // 操作人ID
      "operator_name": string,     // 操作人姓名
      "created_at": string         // 操作时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误，出库数量必须为正整数
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 商品不存在
  - 409: 库存不足，无法出库
  - 500: 服务器内部错误

#### 2.3.9 获取库存预警商品列表

- **URL**: `/api/v1/products/inventory/alert`
- **方法**: GET
- **描述**: 获取库存低于预警值的商品列表
- **请求头**:
  - `Authorization: Bearer {token}`
- **请求参数**:
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `category_id`: number, 按类别筛选
  - `sort_by`: string, 排序字段：stock/stock_ratio（库存/预警值比率）
  - `sort_order`: string, 排序方式：asc/desc
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取库存预警商品列表成功",
    "data": {
      "total": number,             // 总记录数
      "page": number,              // 当前页码
      "page_size": number,         // 每页数量
      "list": [
        {
          "id": number,            // 商品ID
          "name": string,          // 商品名称
          "category_id": number,   // 商品类别ID
          "category_name": string, // 类别名称
          "stock": number,         // 当前库存
          "min_stock": number,     // 最小库存预警值
          "stock_ratio": number,   // 库存/预警值比率（百分比）
          "price": number,         // 售价
          "main_image": string,    // 主图URL
          "sales_30d": number,     // 近30天销量
          "created_at": string,    // 创建时间
          "updated_at": string     // 更新时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.3.10 导出商品数据

- **URL**: `/api/v1/products/export`
- **方法**: GET
- **描述**: 导出商品数据为Excel文件
- **请求头**:
  - `Authorization: Bearer {token}`
- **请求参数**:
  - `category_id`: number, 按类别筛选
  - `status`: number, 商品状态：0-下架，1-上架
  - `keyword`: string, 搜索关键词
  - `stock_alert`: boolean, 是否只导出库存预警商品
  - `export_type`: string, 导出类型：basic-基本信息，full-完整信息（含库存记录）
- **响应数据**:
  ```
  直接返回Excel文件流，Content-Type: application/vnd.ms-excel
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.3.11 导入商品数据

- **URL**: `/api/v1/products/import`
- **方法**: POST
- **描述**: 通过Excel文件导入商品数据
- **请求头**:
  - `Authorization: Bearer {token}`
  - `Content-Type: multipart/form-data`
- **请求参数**:
  - `file`: File, Excel文件（必填）
  - `update_existing`: boolean, 是否更新已存在的商品（默认false）
  - `category_mapping`: string, 类别映射JSON字符串（可选）
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "导入商品数据成功",
    "data": {
      "total": number,             // 总记录数
      "success": number,           // 成功导入数
      "failed": number,            // 失败数
      "updated": number,           // 更新数
      "created": number,           // 新建数
      "errors": [                  // 错误记录
        {
          "row": number,           // Excel行号
          "message": string        // 错误信息
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误，文件格式不正确
  - 401: 未授权，请先登录
  - 413: 文件过大
  - 500: 服务器内部错误

#### 2.3.12 获取商品销售统计

- **URL**: `/api/v1/products/statistics`
- **方法**: GET
- **描述**: 获取商品销售统计数据
- **请求头**:
  - `Authorization: Bearer {token}`
- **请求参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `category_id`: number, 按类别筛选
  - `product_id`: number, 指定商品ID
  - `sort_by`: string, 排序字段：sales_amount/sales_quantity/profit
  - `sort_order`: string, 排序方式：asc/desc
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取商品销售统计成功",
    "data": {
      "summary": {
        "total_sales_amount": number,    // 总销售额
        "total_sales_quantity": number,  // 总销售量
        "total_profit": number,          // 总利润
        "total_products": number         // 销售商品数量
      },
      "by_category": [                   // 按类别统计
        {
          "category_id": number,
          "category_name": string,
          "sales_amount": number,
          "sales_quantity": number,
          "profit": number
        }
      ],
      "by_time": [                       // 按时间统计
        {
          "date": string,                // 日期
          "sales_amount": number,
          "sales_quantity": number
        }
      ],
      "products": {                      // 商品销售排行
        "total": number,                 // 总记录数
        "page": number,                  // 当前页码
        "page_size": number,             // 每页数量
        "list": [
          {
            "product_id": number,        // 商品ID
            "product_name": string,      // 商品名称
            "category_id": number,       // 类别ID
            "category_name": string,     // 类别名称
            "sales_quantity": number,    // 销售量
            "sales_amount": number,      // 销售额
            "profit": number,            // 利润
            "avg_price": number          // 平均售价
          }
        ]
      }
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

### 2.4 营销活动接口

#### 2.4.1 POST /api/v1/activities - 创建营销活动
- **功能描述**: 创建新的营销活动
- **请求参数**:
  ```json
  {
    "name": string,                  // 活动名称
    "type": string,                  // 活动类型：折扣、满减、赠品等
    "start_time": string,            // 活动开始时间，ISO8601格式
    "end_time": string,              // 活动结束时间，ISO8601格式
    "status": string,                // 活动状态：draft(草稿)、active(生效)、inactive(未生效)、expired(已过期)
    "description": string,           // 活动描述
    "rule": {                        // 活动规则，根据不同活动类型有不同结构
      "discount_type": string,       // 折扣类型：percentage(百分比)、fixed(固定金额)
      "discount_value": number,      // 折扣值
      "min_purchase_amount": number, // 最低消费金额
      "max_discount_amount": number  // 最高折扣金额
    },
    "applicable_products": [         // 适用商品列表，为空表示全场通用
      {
        "product_id": number,
        "product_name": string
      }
    ],
    "applicable_categories": [       // 适用品类列表，为空表示全场通用
      {
        "category_id": number,
        "category_name": string
      }
    ],
    "member_levels": [               // 适用会员等级，为空表示对所有会员开放
      {
        "level_id": number,
        "level_name": string
      }
    ],
    "coupon_required": boolean,      // 是否需要优惠券
    "priority": number               // 活动优先级，数字越大优先级越高
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "创建成功",
    "data": {
      "activity_id": number,         // 活动ID
      "name": string,                // 活动名称
      "type": string,                // 活动类型
      "start_time": string,          // 开始时间
      "end_time": string,            // 结束时间
      "status": string,              // 活动状态
      "created_at": string,          // 创建时间
      "created_by": number           // 创建人ID
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.4.2 GET /api/v1/activities - 获取活动列表

- **功能描述**: 获取营销活动列表，支持分页和筛选
- **请求参数**:
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `status`: string, 活动状态筛选
  - `type`: string, 活动类型筛选
  - `name`: string, 活动名称关键词搜索
  - `start_time_begin`: string, 开始时间范围起点
  - `start_time_end`: string, 开始时间范围终点
  - `end_time_begin`: string, 结束时间范围起点
  - `end_time_end`: string, 结束时间范围终点
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": number,               // 总记录数
      "page": number,                // 当前页码
      "page_size": number,           // 每页数量
      "list": [
        {
          "activity_id": number,     // 活动ID
          "name": string,            // 活动名称
          "type": string,            // 活动类型
          "start_time": string,      // 开始时间
          "end_time": string,        // 结束时间
          "status": string,          // 活动状态
          "description": string,     // 活动描述
          "created_at": string,      // 创建时间
          "created_by": {            // 创建人信息
            "user_id": number,
            "username": string
          }
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 500: 服务器内部错误

#### 2.4.3 GET /api/v1/activities/:id - 获取活动详情

- **功能描述**: 获取指定ID的营销活动详细信息
- **路径参数**:
  - `id`: number, 活动ID
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "activity_id": number,         // 活动ID
      "name": string,                // 活动名称
      "type": string,                // 活动类型
      "start_time": string,          // 开始时间
      "end_time": string,            // 结束时间
      "status": string,              // 活动状态
      "description": string,         // 活动描述
      "rule": {                      // 活动规则
        "discount_type": string,
        "discount_value": number,
        "min_purchase_amount": number,
        "max_discount_amount": number
      },
      "applicable_products": [       // 适用商品
        {
          "product_id": number,
          "product_name": string,
          "product_code": string,
          "price": number
        }
      ],
      "applicable_categories": [     // 适用品类
        {
          "category_id": number,
          "category_name": string
        }
      ],
      "member_levels": [             // 适用会员等级
        {
          "level_id": number,
          "level_name": string
        }
      ],
      "coupon_required": boolean,    // 是否需要优惠券
      "priority": number,            // 活动优先级
      "created_at": string,          // 创建时间
      "updated_at": string,          // 更新时间
      "created_by": {                // 创建人信息
        "user_id": number,
        "username": string
      },
      "updated_by": {                // 更新人信息
        "user_id": number,
        "username": string
      }
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 活动不存在
  - 500: 服务器内部错误

#### 2.4.4 PUT /api/v1/activities/:id - 更新活动信息

- **功能描述**: 更新指定ID的营销活动信息
- **路径参数**:
  - `id`: number, 活动ID
- **请求参数**:
  ```json
  {
    "name": string,                  // 活动名称
    "type": string,                  // 活动类型
    "start_time": string,            // 活动开始时间
    "end_time": string,              // 活动结束时间
    "description": string,           // 活动描述
    "rule": {                        // 活动规则
      "discount_type": string,
      "discount_value": number,
      "min_purchase_amount": number,
      "max_discount_amount": number
    },
    "applicable_products": [         // 适用商品ID列表
      number
    ],
    "applicable_categories": [       // 适用品类ID列表
      number
    ],
    "member_levels": [               // 适用会员等级ID列表
      number
    ],
    "coupon_required": boolean,      // 是否需要优惠券
    "priority": number               // 活动优先级
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "更新成功",
    "data": {
      "activity_id": number,         // 活动ID
      "name": string,                // 活动名称
      "updated_at": string           // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 活动不存在
  - 500: 服务器内部错误

#### 2.4.5 DELETE /api/v1/activities/:id - 删除活动

- **功能描述**: 删除指定ID的营销活动
- **路径参数**:
  - `id`: number, 活动ID
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "删除成功",
    "data": null
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足，或活动已上线无法删除
  - 404: 活动不存在
  - 500: 服务器内部错误

#### 2.4.6 POST /api/v1/activities/:id/status - 修改活动状态

- **功能描述**: 修改指定ID的营销活动状态(上线/下线)
- **路径参数**:
  - `id`: number, 活动ID
- **请求参数**:
  ```json
  {
    "status": string                 // 目标状态：active(上线)、inactive(下线)
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "状态更新成功",
    "data": {
      "activity_id": number,         // 活动ID
      "name": string,                // 活动名称
      "status": string,              // 更新后的状态
      "updated_at": string           // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 活动不存在
  - 409: 状态冲突，如已过期活动无法上线
  - 500: 服务器内部错误

#### 2.4.7 POST /api/v1/coupons - 创建优惠券

- **功能描述**: 创建新的优惠券
- **请求参数**:
  ```json
  {
    "name": string,                  // 优惠券名称
    "type": string,                  // 优惠券类型：满减券、折扣券、代金券等
    "value": number,                 // 优惠券面值或折扣率
    "min_purchase_amount": number,   // 最低消费金额
    "start_time": string,            // 生效开始时间，ISO8601格式
    "end_time": string,              // 生效结束时间，ISO8601格式
    "quantity": number,              // 发行数量
    "per_limit": number,             // 每人限领数量
    "description": string,           // 使用说明
    "status": string,                // 状态：active(生效)、inactive(未生效)
    "applicable_products": [         // 适用商品ID列表，空数组表示全场通用
      number
    ],
    "applicable_categories": [       // 适用品类ID列表，空数组表示全场通用
      number
    ]
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "创建成功",
    "data": {
      "coupon_id": number,           // 优惠券ID
      "name": string,                // 优惠券名称
      "type": string,                // 优惠券类型
      "value": number,               // 优惠券面值或折扣率
      "min_purchase_amount": number,  // 最低消费金额
      "start_time": string,          // 生效开始时间
      "end_time": string,            // 生效结束时间
      "quantity": number,            // 发行数量
      "created_at": string           // 创建时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.4.8 GET /api/v1/coupons - 获取优惠券列表

- **功能描述**: 获取系统中的优惠券列表，支持分页和筛选
- **查询参数**:
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `status`: string, 优惠券状态筛选
  - `type`: string, 优惠券类型筛选
  - `name`: string, 优惠券名称关键词搜索
  - `start_time_from`: string, 开始时间筛选(起始)
  - `start_time_to`: string, 开始时间筛选(截止)
  - `end_time_from`: string, 结束时间筛选(起始)
  - `end_time_to`: string, 结束时间筛选(截止)
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": number,               // 总记录数
      "page": number,                // 当前页码
      "page_size": number,           // 每页数量
      "items": [
        {
          "coupon_id": number,       // 优惠券ID
          "name": string,            // 优惠券名称
          "type": string,            // 优惠券类型
          "value": number,           // 优惠券面值或折扣率
          "min_purchase_amount": number, // 最低消费金额
          "start_time": string,      // 生效开始时间
          "end_time": string,        // 生效结束时间
          "quantity": number,        // 发行数量
          "issued_count": number,    // 已发放数量
          "used_count": number,      // 已使用数量
          "status": string,          // 状态
          "created_at": string       // 创建时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.4.9 GET /api/v1/coupons/:id - 获取优惠券详情

- **功能描述**: 获取指定ID的优惠券详细信息
- **路径参数**:
  - `id`: number, 优惠券ID
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "coupon_id": number,           // 优惠券ID
      "name": string,                // 优惠券名称
      "type": string,                // 优惠券类型
      "value": number,               // 优惠券面值或折扣率
      "min_purchase_amount": number,  // 最低消费金额
      "start_time": string,          // 生效开始时间
      "end_time": string,            // 生效结束时间
      "quantity": number,            // 发行数量
      "per_limit": number,           // 每人限领数量
      "issued_count": number,        // 已发放数量
      "used_count": number,          // 已使用数量
      "description": string,         // 使用说明
      "status": string,              // 状态
      "applicable_products": [       // 适用商品ID列表
        {
          "product_id": number,
          "product_name": string
        }
      ],
      "applicable_categories": [     // 适用品类ID列表
        {
          "category_id": number,
          "category_name": string
        }
      ],
      "created_at": string,          // 创建时间
      "updated_at": string           // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 优惠券不存在
  - 500: 服务器内部错误

#### 2.4.10 PUT /api/v1/coupons/:id - 更新优惠券信息

- **功能描述**: 更新指定ID的优惠券信息
- **路径参数**:
  - `id`: number, 优惠券ID
- **请求参数**:
  ```json
  {
    "name": string,                  // 优惠券名称
    "description": string,           // 使用说明
    "status": string,                // 状态：active(生效)、inactive(未生效)
    "start_time": string,            // 生效开始时间，ISO8601格式
    "end_time": string,              // 生效结束时间，ISO8601格式
    "per_limit": number,             // 每人限领数量
    "applicable_products": [         // 适用商品ID列表
      number
    ],
    "applicable_categories": [       // 适用品类ID列表
      number
    ]
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "更新成功",
    "data": {
      "coupon_id": number,           // 优惠券ID
      "name": string,                // 优惠券名称
      "status": string,              // 状态
      "updated_at": string           // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 优惠券不存在
  - 409: 已发放的优惠券无法修改某些属性
  - 500: 服务器内部错误

#### 2.4.11 DELETE /api/v1/coupons/:id - 删除优惠券

- **功能描述**: 删除指定ID的优惠券
- **路径参数**:
  - `id`: number, 优惠券ID
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "删除成功",
    "data": null
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足，或优惠券已发放无法删除
  - 404: 优惠券不存在
  - 500: 服务器内部错误

#### 2.4.12 POST /api/v1/coupons/:id/issue - 发放优惠券给会员

- **功能描述**: 将指定优惠券发放给一个或多个会员
- **路径参数**:
  - `id`: number, 优惠券ID
- **请求参数**:
  ```json
  {
    "member_ids": [                  // 会员ID列表
      number
    ],
    "quantity": number,              // 每人发放数量，默认1
    "note": string                   // 发放备注
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "发放成功",
    "data": {
      "coupon_id": number,           // 优惠券ID
      "coupon_name": string,         // 优惠券名称
      "issued_count": number,        // 本次发放总数量
      "success_count": number,       // 成功发放数量
      "fail_count": number,          // 失败发放数量
      "fail_details": [              // 失败详情
        {
          "member_id": number,
          "reason": string
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 优惠券不存在
  - 409: 优惠券数量不足或已过期
  - 500: 服务器内部错误

#### 2.4.13 GET /api/v1/members/:id/coupons - 获取会员的优惠券列表

- **功能描述**: 获取指定会员拥有的优惠券列表
- **路径参数**:
  - `id`: number, 会员ID
- **查询参数**:
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `status`: string, 优惠券状态筛选：unused(未使用)、used(已使用)、expired(已过期)
  - `type`: string, 优惠券类型筛选
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": number,               // 总记录数
      "page": number,                // 当前页码
      "page_size": number,           // 每页数量
      "items": [
        {
          "member_coupon_id": number, // 会员优惠券ID
          "coupon_id": number,       // 优惠券ID
          "coupon_name": string,     // 优惠券名称
          "type": string,            // 优惠券类型
          "value": number,           // 优惠券面值或折扣率
          "min_purchase_amount": number, // 最低消费金额
          "start_time": string,      // 生效开始时间
          "end_time": string,        // 生效结束时间
          "status": string,          // 状态：unused(未使用)、used(已使用)、expired(已过期)
          "used_time": string,       // 使用时间
          "order_id": number,        // 使用订单ID
          "issued_time": string      // 获得时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 500: 服务器内部错误

#### 2.4.14 GET /api/v1/activities/:id/statistics - 获取活动效果统计

- **功能描述**: 获取指定营销活动的效果统计数据
- **路径参数**:
  - `id`: number, 活动ID
- **查询参数**:
  - `start_date`: string, 统计开始日期，格式YYYY-MM-DD
  - `end_date`: string, 统计结束日期，格式YYYY-MM-DD
  - `group_by`: string, 数据分组方式：day(按天)、week(按周)、month(按月)，默认day
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "activity_id": number,           // 活动ID
      "activity_name": string,         // 活动名称
      "total_views": number,           // 活动浏览次数
      "total_participants": number,    // 参与人数
      "conversion_rate": number,       // 转化率
      "total_orders": number,          // 相关订单数
      "total_sales": number,           // 销售总额
      "average_order_value": number,   // 平均订单金额
      "total_discount_amount": number, // 总优惠金额
      "roi": number,                   // 投资回报率
      "time_series_data": [            // 时间序列数据
        {
          "date": string,              // 日期
          "views": number,             // 浏览次数
          "participants": number,      // 参与人数
          "orders": number,            // 订单数
          "sales": number,             // 销售额
          "discount_amount": number    // 优惠金额
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 活动不存在
  - 500: 服务器内部错误

#### 2.4.15 GET /api/v1/activities/statistics - 获取所有活动统计数据

- **功能描述**: 获取所有营销活动的汇总统计数据
- **查询参数**:
  - `start_date`: string, 统计开始日期，格式YYYY-MM-DD
  - `end_date`: string, 统计结束日期，格式YYYY-MM-DD
  - `type`: string, 活动类型筛选
  - `status`: string, 活动状态筛选
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `sort_by`: string, 排序字段：roi(投资回报率)、sales(销售额)、participants(参与人数)等
  - `sort_order`: string, 排序方式：asc(升序)、desc(降序)，默认desc
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": number,                 // 总记录数
      "page": number,                  // 当前页码
      "page_size": number,             // 每页数量
      "summary": {
        "total_activities": number,    // 活动总数
        "active_activities": number,   // 进行中活动数
        "total_sales": number,         // 总销售额
        "total_discount_amount": number, // 总优惠金额
        "average_roi": number          // 平均ROI
      },
      "items": [
        {
          "activity_id": number,       // 活动ID
          "activity_name": string,     // 活动名称
          "type": string,              // 活动类型
          "status": string,            // 活动状态
          "start_time": string,        // 开始时间
          "end_time": string,          // 结束时间
          "total_views": number,       // 活动浏览次数
          "total_participants": number, // 参与人数
          "conversion_rate": number,    // 转化率
          "total_orders": number,       // 相关订单数
          "total_sales": number,        // 销售总额
          "total_discount_amount": number, // 总优惠金额
          "roi": number                // 投资回报率
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

### 2.5 订单管理接口
#### 2.5.1 POST /api/v1/orders - 创建订单

- **功能描述**: 创建新订单
- **请求参数**:
  ```json
  {
    "member_id": number,              // 会员ID，非会员可为null
    "store_id": number,               // 门店ID
    "cashier_id": number,             // 收银员ID
    "items": [                        // 订单商品列表
      {
        "product_id": number,         // 商品ID
        "quantity": number,           // 数量
        "unit_price": number,         // 单价
        "discount_amount": number     // 商品折扣金额
      }
    ],
    "payment_method": string,         // 支付方式：cash(现金)、wechat(微信)、alipay(支付宝)、card(会员卡)等
    "total_amount": number,           // 订单总金额
    "discount_amount": number,        // 订单折扣总金额
    "actual_amount": number,          // 实际支付金额
    "remark": string,                 // 订单备注
    "activity_ids": [number]          // 关联的营销活动ID列表
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "创建成功",
    "data": {
      "order_id": number,             // 订单ID
      "order_no": string,             // 订单编号
      "created_at": string,           // 创建时间
      "status": string                // 订单状态：pending(待支付)、paid(已支付)、completed(已完成)、cancelled(已取消)
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.5.2 GET /api/v1/orders - 获取订单列表

- **功能描述**: 获取订单列表，支持分页和多种筛选条件
- **请求参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `store_id`: number, 门店ID筛选
  - `member_id`: number, 会员ID筛选
  - `status`: string, 订单状态筛选
  - `payment_method`: string, 支付方式筛选
  - `min_amount`: number, 最小订单金额
  - `max_amount`: number, 最大订单金额
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `sort_by`: string, 排序字段：created_at(创建时间)、total_amount(订单金额)等
  - `sort_order`: string, 排序方式：asc(升序)、desc(降序)，默认desc
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "total": number,                // 总记录数
      "page": number,                 // 当前页码
      "page_size": number,            // 每页数量
      "items": [
        {
          "order_id": number,         // 订单ID
          "order_no": string,         // 订单编号
          "store_name": string,       // 门店名称
          "member_name": string,      // 会员姓名，非会员为null
          "total_amount": number,     // 订单总金额
          "discount_amount": number,  // 折扣金额
          "actual_amount": number,    // 实际支付金额
          "status": string,           // 订单状态
          "payment_method": string,   // 支付方式
          "created_at": string,       // 创建时间
          "cashier_name": string      // 收银员姓名
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.5.3 GET /api/v1/orders/:id - 获取订单详情

- **功能描述**: 获取指定订单的详细信息
- **请求参数**:
  - `id`: number, 订单ID (路径参数)
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "order_id": number,             // 订单ID
      "order_no": string,             // 订单编号
      "store_id": number,             // 门店ID
      "store_name": string,           // 门店名称
      "member_id": number,            // 会员ID，非会员为null
      "member_name": string,          // 会员姓名，非会员为null
      "cashier_id": number,           // 收银员ID
      "cashier_name": string,         // 收银员姓名
      "total_amount": number,         // 订单总金额
      "discount_amount": number,      // 折扣金额
      "actual_amount": number,        // 实际支付金额
      "payment_method": string,       // 支付方式
      "status": string,               // 订单状态
      "remark": string,               // 订单备注
      "created_at": string,           // 创建时间
      "paid_at": string,              // 支付时间
      "completed_at": string,         // 完成时间
      "cancelled_at": string,         // 取消时间
      "items": [                      // 订单商品列表
        {
          "product_id": number,       // 商品ID
          "product_name": string,     // 商品名称
          "product_code": string,     // 商品编码
          "quantity": number,         // 数量
          "unit": string,             // 单位
          "unit_price": number,       // 单价
          "discount_amount": number,  // 折扣金额
          "subtotal": number          // 小计金额
        }
      ],
      "activities": [                 // 关联的营销活动
        {
          "activity_id": number,      // 活动ID
          "activity_name": string,    // 活动名称
          "discount_amount": number   // 该活动提供的折扣金额
        }
      ],
      "payment_details": {            // 支付详情
        "transaction_id": string,     // 交易流水号
        "payment_time": string,       // 支付时间
        "payment_channel": string,    // 支付渠道
        "payment_status": string      // 支付状态
      }
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 订单不存在
  - 500: 服务器内部错误

#### 2.5.4 PUT /api/v1/orders/:id - 更新订单信息

- **功能描述**: 更新订单信息，仅限未支付状态的订单
- **请求参数**:
  - `id`: number, 订单ID (路径参数)
  ```json
  {
    "items": [                        // 订单商品列表
      {
        "product_id": number,         // 商品ID
        "quantity": number,           // 数量
        "unit_price": number,         // 单价
        "discount_amount": number     // 商品折扣金额
      }
    ],
    "payment_method": string,         // 支付方式
    "total_amount": number,           // 订单总金额
    "discount_amount": number,        // 订单折扣总金额
    "actual_amount": number,          // 实际支付金额
    "remark": string,                 // 订单备注
    "activity_ids": [number]          // 关联的营销活动ID列表
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "更新成功",
    "data": {
      "order_id": number,             // 订单ID
      "updated_at": string            // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 订单不存在
  - 409: 订单状态不允许修改
  - 500: 服务器内部错误

#### 2.5.5 DELETE /api/v1/orders/:id - 取消订单

- **功能描述**: 取消指定订单，仅限未支付或待处理状态的订单
- **请求参数**:
  - `id`: number, 订单ID (路径参数)
  ```json
  {
    "cancel_reason": string          // 取消原因
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "取消成功",
    "data": {
      "order_id": number,            // 订单ID
      "cancelled_at": string         // 取消时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 订单不存在
  - 409: 订单状态不允许取消
  - 500: 服务器内部错误

#### 2.5.6 POST /api/v1/orders/:id/status - 更新订单状态

- **功能描述**: 更新订单状态
- **请求参数**:
  - `id`: number, 订单ID (路径参数)
  ```json
  {
    "status": string,                // 目标状态：pending(待支付)、paid(已支付)、completed(已完成)、cancelled(已取消)
    "remark": string                 // 状态变更备注
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "状态更新成功",
    "data": {
      "order_id": number,            // 订单ID
      "previous_status": string,     // 之前的状态
      "current_status": string,      // 当前状态
      "updated_at": string           // 更新时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 订单不存在
  - 409: 状态变更不允许
  - 500: 服务器内部错误

#### 2.5.7 POST /api/v1/orders/:id/pay - 订单支付

- **功能描述**: 处理订单支付
- **请求参数**:
  - `id`: number, 订单ID (路径参数)
  ```json
  {
    "payment_method": string,        // 支付方式：cash(现金)、wechat(微信)、alipay(支付宝)、card(会员卡)等
    "payment_amount": number,        // 支付金额
    "transaction_id": string,        // 外部交易流水号(第三方支付平台)
    "payment_details": object        // 支付详情，根据不同支付方式有不同结构
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "支付成功",
    "data": {
      "order_id": number,            // 订单ID
      "order_no": string,            // 订单编号
      "payment_time": string,        // 支付时间
      "payment_status": string,      // 支付状态
      "transaction_id": string       // 交易流水号
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 订单不存在
  - 409: 订单状态不允许支付
  - 422: 支付金额不匹配
  - 500: 服务器内部错误

#### 2.5.8 POST /api/v1/orders/:id/refund - 订单退款

- **功能描述**: 处理订单退款
- **请求参数**:
  - `id`: number, 订单ID (路径参数)
  ```json
  {
    "refund_amount": number,         // 退款金额
    "refund_reason": string,         // 退款原因
    "refund_items": [                // 退款商品明细，部分退款时使用
      {
        "item_id": number,           // 订单商品ID
        "quantity": number,          // 退款数量
        "amount": number             // 退款金额
      }
    ],
    "operator_id": number,           // 操作人ID
    "remark": string                 // 退款备注
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "退款申请提交成功",
    "data": {
      "order_id": number,            // 订单ID
      "refund_id": number,           // 退款记录ID
      "refund_no": string,           // 退款单号
      "refund_amount": number,       // 退款金额
      "refund_status": string,       // 退款状态：processing(处理中)、success(成功)、failed(失败)
      "created_at": string           // 创建时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 订单不存在
  - 409: 订单状态不允许退款
  - 422: 退款金额超出可退金额
  - 500: 服务器内部错误

#### 2.5.9 GET /api/v1/orders/export - 导出订单数据

- **功能描述**: 导出订单数据为Excel或CSV格式
- **查询参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `store_id`: number, 门店ID筛选
  - `status`: string, 订单状态筛选
  - `payment_method`: string, 支付方式筛选
  - `min_amount`: number, 最小订单金额
  - `max_amount`: number, 最大订单金额
  - `format`: string, 导出格式：excel、csv，默认excel
- **响应结果**: 文件流，用于下载
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.5.10 GET /api/v1/orders/statistics - 获取订单统计数据

- **功能描述**: 获取订单相关统计数据
- **查询参数**:
  - `start_date`: string, 统计开始日期，格式YYYY-MM-DD
  - `end_date`: string, 统计结束日期，格式YYYY-MM-DD
  - `store_id`: number, 门店ID筛选
  - `group_by`: string, 分组方式：day(按天)、week(按周)、month(按月)、year(按年)，默认day
  - `payment_method`: string, 支付方式筛选
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "summary": {
        "total_orders": number,        // 订单总数
        "total_amount": number,        // 总金额
        "average_amount": number,      // 平均订单金额
        "total_discount": number,      // 总优惠金额
        "refund_orders": number,       // 退款订单数
        "refund_amount": number        // 退款总金额
      },
      "payment_distribution": [        // 支付方式分布
        {
          "payment_method": string,    // 支付方式
          "count": number,             // 订单数
          "amount": number,            // 金额
          "percentage": number         // 占比
        }
      ],
      "time_distribution": [           // 时间分布
        {
          "date": string,              // 日期
          "count": number,             // 订单数
          "amount": number,            // 金额
          "average": number            // 平均订单金额
        }
      ],
      "status_distribution": [         // 订单状态分布
        {
          "status": string,            // 订单状态
          "count": number,             // 订单数
          "amount": number,            // 金额
          "percentage": number         // 占比
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

### 2.6 数据统计分析接口
#### 2.6.1 GET /api/v1/statistics/sales - 获取销售统计数据

- **功能描述**: 获取销售统计数据，支持按日/周/月/年分组
- **请求参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `group_by`: string, 分组方式：day(按天)、week(按周)、month(按月)、year(按年)，默认day
  - `store_id`: number, 门店ID，不传则查询所有门店
  - `category_id`: number, 商品类别ID，不传则查询所有类别
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取销售统计数据成功",
    "data": {
      "summary": {
        "total_sales_amount": number,    // 总销售额
        "total_sales_quantity": number,  // 总销售量
        "total_orders": number,          // 总订单数
        "average_order_amount": number,  // 平均订单金额
        "gross_profit": number,          // 毛利润
        "gross_profit_rate": number      // 毛利率
      },
      "time_series": [                   // 时间序列数据
        {
          "date": string,                // 日期
          "sales_amount": number,        // 销售额
          "sales_quantity": number,      // 销售量
          "orders_count": number,        // 订单数
          "profit": number               // 利润
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.6.2 GET /api/v1/statistics/sales/products - 获取商品销售排行

- **功能描述**: 获取商品销售排行数据
- **请求参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `store_id`: number, 门店ID，不传则查询所有门店
  - `category_id`: number, 商品类别ID，不传则查询所有类别
  - `sort_by`: string, 排序字段：sales_amount(销售额)、sales_quantity(销售量)、profit(利润)，默认sales_amount
  - `sort_order`: string, 排序方式：asc(升序)、desc(降序)，默认desc
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取商品销售排行成功",
    "data": {
      "total": number,                 // 总记录数
      "page": number,                  // 当前页码
      "page_size": number,             // 每页数量
      "items": [
        {
          "product_id": number,        // 商品ID
          "product_name": string,      // 商品名称
          "category_id": number,       // 类别ID
          "category_name": string,     // 类别名称
          "sales_quantity": number,    // 销售量
          "sales_amount": number,      // 销售额
          "profit": number,            // 利润
          "profit_rate": number,       // 利润率
          "avg_price": number          // 平均售价
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.6.3 GET /api/v1/statistics/sales/categories - 获取品类销售分布

- **功能描述**: 获取商品类别销售分布数据
- **请求参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `store_id`: number, 门店ID，不传则查询所有门店
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取品类销售分布成功",
    "data": {
      "total_sales_amount": number,    // 总销售额
      "categories": [
        {
          "category_id": number,       // 类别ID
          "category_name": string,     // 类别名称
          "sales_amount": number,      // 销售额
          "sales_quantity": number,    // 销售量
          "percentage": number,        // 销售额占比
          "profit": number,            // 利润
          "profit_rate": number        // 利润率
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.6.4 GET /api/v1/statistics/members - 获取会员增长统计

- **功能描述**: 获取会员增长统计数据
- **请求参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `group_by`: string, 分组方式：day(按天)、week(按周)、month(按月)、year(按年)，默认day
  - `store_id`: number, 门店ID，不传则查询所有门店
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取会员增长统计成功",
    "data": {
      "summary": {
        "total_members": number,       // 会员总数
        "new_members": number,         // 新增会员数
        "active_members": number,      // 活跃会员数
        "inactive_members": number     // 不活跃会员数
      },
      "time_series": [                 // 时间序列数据
        {
          "date": string,              // 日期
          "new_members": number,       // 新增会员数
          "accumulated": number,       // 累计会员数
          "active_members": number     // 活跃会员数
        }
      ],
      "source_distribution": [         // 会员来源分布
        {
          "source": string,            // 来源
          "count": number,             // 数量
          "percentage": number         // 占比
        }
      ],
      "level_distribution": [          // 会员等级分布
        {
          "level": string,             // 等级
          "count": number,             // 数量
          "percentage": number         // 占比
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.6.5 GET /api/v1/statistics/members/consumption - 获取会员消费统计

- **功能描述**: 获取会员消费统计数据
- **请求参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `store_id`: number, 门店ID，不传则查询所有门店
  - `level_id`: number, 会员等级ID，不传则查询所有等级
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `sort_by`: string, 排序字段：consumption_amount(消费金额)、consumption_count(消费次数)、average_amount(平均消费)
  - `sort_order`: string, 排序方式：asc(升序)、desc(降序)，默认desc
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取会员消费统计成功",
    "data": {
      "summary": {
        "total_consumption_amount": number,  // 总消费金额
        "total_consumption_count": number,   // 总消费次数
        "average_consumption": number,       // 平均消费金额
        "consumption_members": number        // 消费会员数
      },
      "rfm_analysis": {                      // RFM分析
        "high_value": number,                // 高价值会员数
        "potential": number,                 // 潜力会员数
        "loyal": number,                     // 忠诚会员数
        "risk": number,                      // 流失风险会员数
        "sleep": number                      // 沉睡会员数
      },
      "consumption_distribution": [          // 消费金额分布
        {
          "range": string,                   // 金额范围
          "count": number,                   // 会员数
          "percentage": number               // 占比
        }
      ],
      "members": {                           // 会员消费排行
        "total": number,                     // 总记录数
        "page": number,                      // 当前页码
        "page_size": number,                 // 每页数量
        "items": [
          {
            "member_id": number,             // 会员ID
            "member_name": string,           // 会员姓名
            "phone": string,                 // 手机号
            "level": string,                 // 会员等级
            "consumption_amount": number,    // 消费金额
            "consumption_count": number,     // 消费次数
            "average_amount": number,        // 平均消费金额
            "last_consumption": string       // 最后消费时间
          }
        ]
      }
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.6.6 GET /api/v1/statistics/members/retention - 获取会员留存率

- **功能描述**: 获取会员留存率统计数据
- **请求参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `group_by`: string, 分组方式：week(按周)、month(按月)，默认month
  - `store_id`: number, 门店ID，不传则查询所有门店
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取会员留存率成功",
    "data": {
      "summary": {
        "average_retention_rate": number,  // 平均留存率
        "churn_rate": number               // 流失率
      },
      "retention_data": [                  // 留存数据
        {
          "date": string,                  // 日期
          "new_members": number,           // 新增会员数
          "retention_rates": [             // 留存率数组
            {
              "period": number,            // 周期(1表示1周后或1月后)
              "rate": number               // 留存率
            }
          ]
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.6.7 GET /api/v1/statistics/dashboard - 获取仪表盘数据概览

- **功能描述**: 获取系统仪表盘数据概览
- **请求参数**:
  - `date_range`: string, 日期范围：today(今日)、yesterday(昨日)、last7days(近7天)、last30days(近30天)、thisMonth(本月)、lastMonth(上月)，默认today
  - `store_id`: number, 门店ID，不传则查询所有门店
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取仪表盘数据成功",
    "data": {
      "sales": {
        "total_amount": number,          // 销售总额
        "comparison": number,            // 与上一周期比较，正数表示增长，负数表示下降
        "orders_count": number,          // 订单数
        "average_order_amount": number   // 平均订单金额
      },
      "members": {
        "total_count": number,           // 会员总数
        "new_count": number,             // 新增会员数
        "comparison": number,            // 与上一周期比较
        "active_count": number           // 活跃会员数
      },
      "products": {
        "total_count": number,           // 商品总数
        "sold_count": number,            // 已售商品数
        "stock_alert_count": number      // 库存预警商品数
      },
      "activities": {
        "active_count": number,          // 进行中活动数
        "upcoming_count": number,        // 即将开始活动数
        "ended_count": number            // 已结束活动数
      },
      "top_selling_products": [          // 热销商品
        {
          "product_id": number,
          "product_name": string,
          "sales_quantity": number,
          "sales_amount": number
        }
      ],
      "recent_orders": [                 // 最近订单
        {
          "order_id": number,
          "order_no": string,
          "created_at": string,
          "amount": number,
          "status": string
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.6.8 GET /api/v1/statistics/trend - 获取业务趋势数据

- **功能描述**: 获取业务趋势数据，包括销售额、订单量、会员增长等关键指标的趋势变化
- **请求参数**:
  - `start_date`: string, 开始日期，格式YYYY-MM-DD
  - `end_date`: string, 结束日期，格式YYYY-MM-DD
  - `group_by`: string, 分组方式：day(按天)、week(按周)、month(按月)、year(按年)，默认day
  - `metrics`: string, 需要获取的指标，多个指标用逗号分隔，可选值：sales_amount(销售额)、orders_count(订单量)、new_members(新增会员)、active_members(活跃会员)、profit(利润)
  - `store_id`: number, 门店ID，不传则查询所有门店
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取业务趋势数据成功",
    "data": {
      "time_range": {
        "start_date": string,          // 开始日期
        "end_date": string,            // 结束日期
        "group_by": string             // 分组方式
      },
      "trends": [
        {
          "date": string,              // 日期
          "sales_amount": number,       // 销售额
          "orders_count": number,       // 订单量
          "new_members": number,        // 新增会员
          "active_members": number,     // 活跃会员
          "profit": number              // 利润
        }
      ],
      "comparison": {                  // 与上一周期比较
        "sales_amount": number,         // 销售额环比变化率
        "orders_count": number,         // 订单量环比变化率
        "new_members": number,          // 新增会员环比变化率
        "active_members": number,       // 活跃会员环比变化率
        "profit": number                // 利润环比变化率
      }
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

### 2.7 系统管理接口
#### 2.7.1 POST /api/v1/users - 创建系统用户

- **功能描述**: 创建新的系统用户
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "username": string,        // 用户名，必填
    "password": string,        // 密码，必填
    "name": string,            // 姓名，必填
    "email": string,           // 电子邮箱，必填
    "mobile": string,          // 手机号码，必填
    "role_ids": [number],      // 角色ID数组，必填
    "store_ids": [number],     // 所属门店ID数组，可选
    "avatar": string,          // 头像URL，可选
    "status": number,          // 状态：1-启用，0-禁用，默认1
    "remark": string           // 备注，可选
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "创建用户成功",
    "data": {
      "id": number,
      "username": string,
      "name": string,
      "email": string,
      "mobile": string,
      "avatar": string,
      "status": number,
      "roles": [
        {
          "id": number,
          "name": string,
          "description": string
        }
      ],
      "stores": [
        {
          "id": number,
          "name": string
        }
      ],
      "created_at": string,
      "updated_at": string
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 409: 用户名已存在
  - 500: 服务器内部错误

#### 2.7.2 GET /api/v1/users - 获取系统用户列表

- **功能描述**: 获取系统用户列表，支持分页和筛选
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `keyword`: string, 关键词搜索(用户名/姓名/手机号/邮箱)
  - `status`: number, 状态筛选：1-启用，0-禁用
  - `role_id`: number, 按角色ID筛选
  - `store_id`: number, 按门店ID筛选
  - `sort_field`: string, 排序字段，默认created_at
  - `sort_order`: string, 排序方式(asc/desc)，默认desc
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取用户列表成功",
    "data": {
      "total": number,          // 总记录数
      "page": number,           // 当前页码
      "page_size": number,      // 每页数量
      "total_pages": number,    // 总页数
      "items": [
        {
          "id": number,
          "username": string,
          "name": string,
          "email": string,
          "mobile": string,
          "avatar": string,
          "status": number,
          "roles": [
            {
              "id": number,
              "name": string
            }
          ],
          "stores": [
            {
              "id": number,
              "name": string
            }
          ],
          "created_at": string,
          "updated_at": string
        }
      ]
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.7.3 GET /api/v1/users/:id - 获取用户详情

- **功能描述**: 获取指定用户的详细信息
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取用户详情成功",
    "data": {
      "id": number,
      "username": string,
      "name": string,
      "email": string,
      "mobile": string,
      "avatar": string,
      "status": number,
      "remark": string,
      "roles": [
        {
          "id": number,
          "name": string,
          "description": string
        }
      ],
      "stores": [
        {
          "id": number,
          "name": string,
          "address": string
        }
      ],
      "permissions": [string],  // 权限列表
      "last_login_time": string,
      "last_login_ip": string,
      "created_at": string,
      "updated_at": string
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 用户不存在
  - 500: 服务器内部错误

#### 2.7.4 PUT /api/v1/users/:id - 更新用户信息

- **功能描述**: 更新指定用户的信息
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "name": string,            // 姓名，可选
    "email": string,           // 电子邮箱，可选
    "mobile": string,          // 手机号码，可选
    "role_ids": [number],      // 角色ID数组，可选
    "store_ids": [number],     // 所属门店ID数组，可选
    "avatar": string,          // 头像URL，可选
    "status": number,          // 状态：1-启用，0-禁用，可选
    "remark": string,          // 备注，可选
    "password": string         // 密码，可选，不修改则不传
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "更新用户信息成功",
    "data": {
      "id": number,
      "username": string,
      "name": string,
      "email": string,
      "mobile": string,
      "avatar": string,
      "status": number,
      "roles": [
        {
          "id": number,
          "name": string
        }
      ],
      "stores": [
        {
          "id": number,
          "name": string
        }
      ],
      "updated_at": string
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 用户不存在
  - 500: 服务器内部错误

#### 2.7.5 DELETE /api/v1/users/:id - 删除用户

- **功能描述**: 删除指定用户
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "删除用户成功",
    "data": null
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 用户不存在
  - 500: 服务器内部错误

#### 2.7.6 POST /api/v1/users/:id/status - 更改用户状态

- **功能描述**: 启用或禁用指定用户
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "status": number  // 状态：1-启用，0-禁用，必填
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "更改用户状态成功",
    "data": {
      "id": number,
      "username": string,
      "status": number,
      "updated_at": string
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 用户不存在
  - 500: 服务器内部错误

#### 2.7.7 GET /api/v1/roles - 获取角色列表

- **功能描述**: 获取系统中所有角色列表
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  - `page`: number, 页码，默认1
  - `page_size`: number, 每页数量，默认20
  - `keyword`: string, 搜索关键词，可选
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取角色列表成功",
    "data": {
      "total": number,             // 总记录数
      "page": number,              // 当前页码
      "page_size": number,         // 每页数量
      "items": [
        {
          "id": number,            // 角色ID
          "name": string,          // 角色名称
          "description": string,   // 角色描述
          "users_count": number,   // 拥有该角色的用户数
          "permissions_count": number, // 该角色拥有的权限数
          "created_at": string,    // 创建时间
          "updated_at": string     // 更新时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.7.8 POST /api/v1/roles - 创建角色

- **功能描述**: 创建新的系统角色
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "name": string,              // 角色名称，必填
    "description": string,       // 角色描述，可选
    "permission_ids": [number]   // 权限ID数组，可选
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "创建角色成功",
    "data": {
      "id": number,
      "name": string,
      "description": string,
      "created_at": string,
      "updated_at": string,
      "permissions": [
        {
          "id": number,
          "name": string,
          "code": string,
          "description": string
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 409: 角色名称已存在
  - 500: 服务器内部错误

#### 2.7.9 PUT /api/v1/roles/:id - 更新角色信息

- **功能描述**: 更新指定角色的信息
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "name": string,              // 角色名称，可选
    "description": string,       // 角色描述，可选
    "permission_ids": [number]   // 权限ID数组，可选
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "更新角色成功",
    "data": {
      "id": number,
      "name": string,
      "description": string,
      "updated_at": string,
      "permissions": [
        {
          "id": number,
          "name": string,
          "code": string,
          "description": string
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 角色不存在
  - 409: 角色名称已存在
  - 500: 服务器内部错误

#### 2.7.10 DELETE /api/v1/roles/:id - 删除角色

- **功能描述**: 删除指定角色
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "删除角色成功",
    "data": null
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 角色不存在
  - 409: 角色正在使用中，无法删除
  - 500: 服务器内部错误

#### 2.7.11 GET /api/v1/permissions - 获取权限列表

- **功能描述**: 获取系统中所有权限列表
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  - `group_by`: string, 分组方式：module(按模块分组)，可选
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取权限列表成功",
    "data": {
      "total": number,           // 总记录数
      "items": [
        {
          "id": number,          // 权限ID
          "name": string,        // 权限名称
          "code": string,        // 权限代码
          "module": string,      // 所属模块
          "description": string  // 权限描述
        }
      ]
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.7.12 POST /api/v1/roles/:id/permissions - 分配权限给角色

- **功能描述**: 为指定角色分配权限
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "permission_ids": [number]   // 权限ID数组，必填
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "分配权限成功",
    "data": {
      "role_id": number,
      "role_name": string,
      "permissions": [
        {
          "id": number,
          "name": string,
          "code": string,
          "module": string,
          "description": string
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 角色不存在
  - 500: 服务器内部错误


#### 2.7.13 GET /api/v1/settings - 获取系统设置

- **功能描述**: 获取系统配置参数
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "category": string  // 可选，设置类别：基本信息/业务参数/安全参数/通知参数
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取系统设置成功",
    "data": {
      "settings": [
        {
          "id": number,
          "category": string,      // 设置类别
          "key": string,           // 设置键名
          "value": string,         // 设置值
          "description": string,   // 设置描述
          "updated_at": string,    // 最后更新时间
          "updated_by": string     // 最后更新人
        }
      ]
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.7.14 PUT /api/v1/settings - 更新系统设置

- **功能描述**: 更新系统配置参数
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "settings": [
      {
        "key": string,     // 设置键名，必填
        "value": string,   // 设置值，必填
        "category": string // 设置类别，必填
      }
    ]
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "更新系统设置成功",
    "data": {
      "updated_settings": [
        {
          "key": string,
          "value": string,
          "category": string,
          "updated_at": string
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.7.15 GET /api/v1/logs - 获取系统操作日志

- **功能描述**: 分页查询系统操作日志
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "page": number,           // 页码，默认1
    "page_size": number,      // 每页记录数，默认20
    "start_time": string,     // 开始时间，可选
    "end_time": string,       // 结束时间，可选
    "operator_id": number,    // 操作人ID，可选
    "operation_type": string, // 操作类型，可选
    "module": string,         // 操作模块，可选
    "ip_address": string,     // IP地址，可选
    "status": string          // 操作状态，可选
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 200,
    "message": "获取操作日志成功",
    "data": {
      "total": number,
      "page": number,
      "page_size": number,
      "logs": [
        {
          "id": number,
          "operator_id": number,
          "operator_name": string,
          "operation_type": string,
          "module": string,
          "operation_desc": string,
          "operation_time": string,
          "ip_address": string,
          "user_agent": string,
          "status": string,
          "request_params": object,
          "response_data": object,
          "error_message": string
        }
      ]
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.7.16 GET /api/v1/logs/export - 导出操作日志

- **功能描述**: 导出系统操作日志为Excel文件
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "start_time": string,     // 开始时间，可选
    "end_time": string,       // 结束时间，可选
    "operator_id": number,    // 操作人ID，可选
    "operation_type": string, // 操作类型，可选
    "module": string,         // 操作模块，可选
    "ip_address": string,     // IP地址，可选
    "status": string,         // 操作状态，可选
    "export_format": string   // 导出格式，默认"xlsx"
  }
  ```
- **响应结果**: 文件流，Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误


### 2.8 登录认证接口

#### 2.8.1 POST /api/v1/auth/login - 用户登录

- **功能描述**: 用户登录系统获取访问令牌
- **请求参数**:
  ```json
  {
    "username": string,       // 用户名/手机号/邮箱
    "password": string,       // 密码
    "captcha_code": string,   // 验证码，可选
    "captcha_id": string      // 验证码ID，可选
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "token": string,           // 访问令牌
      "refresh_token": string,   // 刷新令牌
      "expires_in": number,      // 过期时间(秒)
      "user_info": {
        "id": number,
        "username": string,
        "real_name": string,
        "avatar": string,
        "roles": [string],
        "permissions": [string]
      }
    }
  }
  ```
- **错误码**:
  - 400: 用户名或密码错误
  - 401: 验证码错误或已过期
  - 403: 账号已被锁定
  - 500: 服务器内部错误

#### 2.8.2 POST /api/v1/auth/logout - 用户登出

- **功能描述**: 用户退出登录，使当前令牌失效
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应结果**:
  ```json
  {
    "code": 0,
    "message": "已成功退出登录"
  }
  ```
- **错误码**:
  - 401: 未授权或令牌已失效
  - 500: 服务器内部错误

#### 2.8.3 GET /api/v1/auth/me - 获取当前用户信息

- **功能描述**: 获取当前登录用户的详细信息
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应结果**:
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": number,
      "username": string,
      "real_name": string,
      "avatar": string,
      "email": string,
      "phone": string,
      "gender": string,
      "status": string,
      "last_login_time": string,
      "last_login_ip": string,
      "roles": [
        {
          "id": number,
          "name": string,
          "code": string,
          "description": string
        }
      ],
      "permissions": [string],
      "stores": [
        {
          "id": number,
          "name": string,
          "is_primary": boolean
        }
      ]
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 500: 服务器内部错误

#### 2.8.4 POST /api/v1/auth/refresh-token - 刷新访问令牌

- **功能描述**: 使用刷新令牌获取新的访问令牌
- **请求参数**:
  ```json
  {
    "refresh_token": string  // 刷新令牌
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "token": string,         // 新的访问令牌
      "refresh_token": string, // 新的刷新令牌
      "expires_in": number     // 过期时间(秒)
    }
  }
  ```
- **错误码**:
  - 400: 无效的刷新令牌
  - 401: 刷新令牌已过期
  - 500: 服务器内部错误

#### 2.8.5 POST /api/v1/auth/change-password - 修改密码

- **功能描述**: 已登录用户修改自己的密码
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "old_password": string,  // 旧密码
    "new_password": string,  // 新密码
    "confirm_password": string  // 确认新密码
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 0,
    "message": "密码修改成功"
  }
  ```
- **错误码**:
  - 400: 新密码不符合安全要求/两次输入的新密码不一致
  - 401: 旧密码验证失败
  - 500: 服务器内部错误

#### 2.8.6 POST /api/v1/auth/forgot-password - 忘记密码

- **功能描述**: 用户忘记密码，申请重置密码
- **请求参数**:
  ```json
  {
    "username": string,     // 用户名/手机号/邮箱
    "captcha_code": string, // 验证码
    "captcha_id": string    // 验证码ID
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 0,
    "message": "重置密码链接已发送到您的邮箱/手机"
  }
  ```
- **错误码**:
  - 400: 用户不存在/验证码错误
  - 429: 请求过于频繁，请稍后再试
  - 500: 服务器内部错误

#### 2.8.7 POST /api/v1/auth/reset-password - 重置密码

- **功能描述**: 通过重置令牌设置新密码
- **请求参数**:
  ```json
  {
    "reset_token": string,      // 重置令牌
    "new_password": string,     // 新密码
    "confirm_password": string  // 确认新密码
  }
  ```
- **响应结果**:
  ```json
  {
    "code": 0,
    "message": "密码重置成功"
  }
  ```
- **错误码**:
  - 400: 重置令牌无效/两次输入的密码不一致/新密码不符合安全要求
  - 401: 重置令牌已过期
  - 500: 服务器内部错误

### 2.9 财务管理接口
```
GET    /api/v1/transactions                   # 获取交易记录列表
GET    /api/v1/transactions/:id               # 获取交易详情
POST   /api/v1/transactions                   # 创建交易记录
POST   /api/v1/transactions/:id/refund        # 退款操作

POST   /api/v1/member-accounts/recharge       # 会员余额充值
GET    /api/v1/member-accounts/:member_id/records # 获取会员账户明细
POST   /api/v1/member-accounts/:member_id/adjust  # 余额调整

GET    /api/v1/finance/revenue-overview        # 获取营收概览
GET    /api/v1/finance/member-consumption-report # 获取会员消费报表
GET    /api/v1/finance/reconciliation-report    # 获取财务对账报表
POST   /api/v1/finance/export                   # 导出财务数据
GET    /api/v1/finance/export/:task_id          # 获取导出任务状态
```

#### 2.5.1 获取交易记录列表
- **URL**: `/api/v1/transactions`
- **方法**: GET
- **描述**: 获取系统交易记录列表，支持分页和多条件筛选
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  page: number               // 页码，从1开始，默认1
  page_size: number          // 每页数量，默认20
  keyword: string            // 关键词搜索(交易单号、会员姓名等)
  transaction_type: string   // 交易类型：RECHARGE(充值), CONSUME(消费), REFUND(退款), COMMISSION(分佣) 等
  payment_method: string     // 支付方式：CASH(现金), CARD(银行卡), WECHAT(微信), ALIPAY(支付宝), BALANCE(余额)等
  start_date: string         // 开始日期筛选
  end_date: string           // 结束日期筛选
  min_amount: number         // 最小金额
  max_amount: number         // 最大金额
  status: string             // 状态：PENDING(待处理), SUCCESS(成功), FAILED(失败)
  store_id: number           // 门店ID
  member_id: number          // 会员ID
  operator_id: number        // 操作员ID
  sort_field: string         // 排序字段：created_at, amount
  sort_order: string         // 排序方式(asc/desc)，默认desc
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取交易记录成功",
    "data": {
      "total": number,             // 总记录数
      "page": number,              // 当前页码
      "page_size": number,         // 每页数量
      "total_pages": number,       // 总页数
      "items": [
        {
          "id": number,            // 交易ID
          "transaction_no": "string", // 交易单号
          "transaction_type": "string", // 交易类型
          "transaction_type_text": "string", // 交易类型文本
          "amount": number,        // 交易金额
          "payment_method": "string", // 支付方式
          "payment_method_text": "string", // 支付方式文本
          "status": "string",      // 交易状态
          "status_text": "string", // 交易状态文本
          "member": {              // 相关会员信息
            "id": number,
            "name": "string",
            "phone": "string"
          },
          "store": {               // 门店信息
            "id": number,
            "name": "string"
          },
          "operator": {            // 操作员信息
            "id": number,
            "name": "string"
          },
          "remark": "string",      // 备注
          "created_at": "string"   // 创建时间
        }
      ]
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.5.2 获取交易详情
- **URL**: `/api/v1/transactions/:id`
- **方法**: GET
- **描述**: 获取单个交易记录的详细信息
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取交易详情成功",
    "data": {
      "id": number,                // 交易ID
      "transaction_no": "string",  // 交易单号
      "external_order_no": "string", // 外部订单号（如有）
      "transaction_type": "string", // 交易类型
      "transaction_type_text": "string", // 交易类型文本
      "amount": number,            // 交易金额
      "actual_amount": number,     // 实际支付金额
      "discount_amount": number,   // 优惠金额
      "points_used": number,       // 使用积分
      "points_converted_amount": number, // 积分折算金额
      "payment_method": "string",  // 支付方式
      "payment_method_text": "string", // 支付方式文本
      "status": "string",          // 交易状态
      "status_text": "string",     // 交易状态文本
      "member": {                  // 相关会员信息
        "id": number,
        "name": "string",
        "phone": "string",
        "level": {
          "id": number,
          "name": "string"
        },
        "balance_before": number,  // 交易前余额
        "balance_after": number,   // 交易后余额
        "points_before": number,   // 交易前积分
        "points_after": number     // 交易后积分
      },
      "store": {                   // 门店信息
        "id": number,
        "name": "string",
        "address": "string"
      },
      "operator": {                // 操作员信息
        "id": number,
        "name": "string",
        "role": "string"
      },
      "related_order": {           // 关联订单信息（如有）
        "id": number,
        "order_no": "string",
        "type": "string"
      },
      "items": [                   // 交易项目明细（如消费明细）
        {
          "product_id": number,
          "product_name": "string",
          "quantity": number,
          "unit_price": number,
          "total_price": number,
          "discount_amount": number
        }
      ],
      "refund_info": {             // 退款信息（如有）
        "refund_no": "string",
        "refund_amount": number,
        "refund_reason": "string",
        "refund_time": "string",
        "operator": {
          "id": number,
          "name": "string"
        }
      },
      "coupon_info": [             // 使用的优惠券信息（如有）
        {
          "coupon_id": number,
          "coupon_code": "string",
          "coupon_name": "string",
          "discount_amount": number
        }
      ],
      "activity_info": {           // 参与的活动信息（如有）
        "activity_id": number,
        "activity_name": "string",
        "discount_amount": number
      },
      "payment_details": {         // 支付详情
        "payment_time": "string",
        "payment_channel": "string",
        "transaction_id": "string",
        "payment_status": "string"
      },
      "remark": "string",          // 备注
      "created_at": "string",      // 创建时间
      "updated_at": "string"       // 更新时间
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 交易记录不存在
  - 500: 服务器内部错误

#### 2.5.3 创建交易记录
- **URL**: `/api/v1/transactions`
- **方法**: POST
- **描述**: 手动创建交易记录（主要用于现金交易或后台调整）
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "transaction_type": "string",  // 交易类型，必填
    "amount": number,              // 交易金额，必填
    "payment_method": "string",    // 支付方式，必填
    "member_id": number,           // 会员ID，必填
    "store_id": number,            // 门店ID，必填
    "items": [                     // 交易项目明细（可选）
      {
        "product_id": number,
        "product_name": "string",
        "quantity": number,
        "unit_price": number
      }
    ],
    "remark": "string"             // 备注，可选
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "创建交易记录成功",
    "data": {
      "id": number,                // 交易ID
      "transaction_no": "string",  // 交易单号
      "transaction_type": "string", // 交易类型
      "amount": number,            // 交易金额
      "status": "string",          // 交易状态
      "created_at": "string",      // 创建时间
      "member": {
        "id": number,
        "name": "string",
        "balance": number          // 交易后余额
      }
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员或门店不存在
  - 500: 服务器内部错误

#### 2.5.4 退款操作
- **URL**: `/api/v1/transactions/:id/refund`
- **方法**: POST
- **描述**: 对已完成的交易进行退款操作
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "refund_amount": number,       // 退款金额，默认为全额退款
    "refund_reason": "string",     // 退款原因，必填
    "refund_to_balance": boolean,  // 是否退回到会员余额，默认false
    "items": [                     // 退款项目明细（部分退款时使用）
      {
        "item_id": number,
        "quantity": number,
        "amount": number
      }
    ]
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "退款操作成功",
    "data": {
      "transaction_id": number,    // 原交易ID
      "refund_id": number,         // 退款交易ID
      "refund_no": "string",       // 退款单号
      "refund_amount": number,     // 退款金额
      "refund_status": "string",   // 退款状态
      "member": {
        "id": number,
        "name": "string",
        "balance": number          // 退款后余额
      },
      "operator": {
        "id": number,
        "name": "string"
      },
      "refund_time": "string"      // 退款时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 交易记录不存在
  - 409: 该交易不支持退款或已超过退款期限
  - 500: 服务器内部错误

#### 2.5.5 会员余额充值
- **URL**: `/api/v1/member-accounts/recharge`
- **方法**: POST
- **描述**: 为会员账户充值
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "member_id": number,           // 会员ID，必填
    "amount": number,              // 充值金额，必填
    "payment_method": "string",    // 支付方式，必填
    "activity_id": number,         // 相关充值活动ID，可选
    "gift_amount": number,         // 赠送金额，可选
    "gift_points": number,         // 赠送积分，可选
    "store_id": number,            // 门店ID，必填
    "remark": "string"             // 备注，可选
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "充值成功",
    "data": {
      "transaction_id": number,    // 交易ID
      "transaction_no": "string",  // 交易单号
      "member": {
        "id": number,
        "name": "string",
        "phone": "string",
        "previous_balance": number, // 充值前余额
        "current_balance": number,  // 充值后余额
        "previous_points": number,  // 充值前积分
        "current_points": number    // 充值后积分
      },
      "amount": number,            // 充值金额
      "gift_amount": number,       // 赠送金额
      "gift_points": number,       // 赠送积分
      "total_amount": number,      // 总金额（充值+赠送）
      "created_at": "string"       // 充值时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 500: 服务器内部错误

#### 2.5.6 获取会员账户明细
- **URL**: `/api/v1/member-accounts/:member_id/records`
- **方法**: GET
- **描述**: 获取指定会员的账户交易记录
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  page: number               // 页码，从1开始，默认1
  page_size: number          // 每页数量，默认20
  record_type: string        // 记录类型：ALL, RECHARGE, CONSUME, REFUND, ADJUSTMENT
  start_date: string         // 开始日期筛选
  end_date: string           // 结束日期筛选
  sort_order: string         // 排序方式(asc/desc)，默认desc
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取会员账户明细成功",
    "data": {
      "member": {
        "id": number,
        "name": "string",
        "phone": "string",
        "level": {
          "id": number,
          "name": "string"
        },
        "current_balance": number, // 当前余额
        "current_points": number,  // 当前积分
        "total_recharge": number,  // 累计充值
        "total_consume": number,   // 累计消费
        "total_refund": number,    // 累计退款
        "first_recharge_time": "string", // 首次充值时间
        "last_recharge_time": "string",  // 最后充值时间
        "total_transactions": number     // 交易总次数
      },
      "records": {
        "total": number,           // 总记录数
        "page": number,            // 当前页码
        "page_size": number,       // 每页数量
        "total_pages": number,     // 总页数
        "items": [
          {
            "id": number,          // 记录ID
            "transaction_no": "string", // 交易单号
            "type": "string",      // 记录类型
            "type_text": "string", // 记录类型文本
            "amount": number,      // 变动金额
            "points": number,      // 变动积分
            "balance_after": number, // 操作后余额
            "points_after": number,  // 操作后积分
            "remark": "string",    // 备注
            "related_info": {      // 相关信息（如订单、活动等）
              "type": "string",
              "id": number,
              "name": "string"
            },
            "store": {             // 门店信息
              "id": number,
              "name": "string"
            },
            "operator": {          // 操作员信息
              "id": number,
              "name": "string"
            },
            "created_at": "string" // 创建时间
          }
        ]
      }
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 500: 服务器内部错误

#### 2.5.7 余额调整
- **URL**: `/api/v1/member-accounts/:member_id/adjust`
- **方法**: POST
- **描述**: 手动调整会员账户余额（增加或减少）
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "adjust_type": "string",     // 调整类型：INCREASE(增加), DECREASE(减少)
    "amount": number,            // 调整金额，必填
    "points": number,            // 调整积分，可选
    "reason": "string",          // 调整原因，必填
    "remark": "string"           // 备注，可选
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "账户调整成功",
    "data": {
      "member_id": number,
      "member_name": "string",
      "adjust_type": "string",
      "amount": number,
      "points": number,
      "previous_balance": number,
      "current_balance": number,
      "previous_points": number,
      "current_points": number,
      "operator": {
        "id": number,
        "name": "string"
      },
      "adjust_time": "string"
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 会员不存在
  - 409: 余额不足，无法减少
  - 500: 服务器内部错误

#### 2.5.8 获取营收概览
- **URL**: `/api/v1/finance/revenue-overview`
- **方法**: GET
- **描述**: 获取系统营收概览数据
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  time_range: string         // 时间范围：TODAY, YESTERDAY, THIS_WEEK, LAST_WEEK, THIS_MONTH, LAST_MONTH, CUSTOM
  start_date: string         // 自定义开始日期（当time_range为CUSTOM时必填）
  end_date: string           // 自定义结束日期（当time_range为CUSTOM时必填）
  store_id: number           // 门店ID，可选，不填则查询所有门店
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取营收概览成功",
    "data": {
      "time_range": {
        "start_date": "string",
        "end_date": "string",
        "days": number
      },
      "total_revenue": number,       // 总营收
      "total_revenue_growth": number, // 总营收增长率
      "total_transactions": number,   // 总交易笔数
      "total_transactions_growth": number, // 总交易笔数增长率
      "average_transaction": number,  // 平均交易金额
      "average_transaction_growth": number, // 平均交易金额增长率
      "revenue_by_type": [           // 按交易类型的营收
        {
          "type": "string",
          "type_text": "string",
          "amount": number,
          "percentage": number,
          "growth": number
        }
      ],
      "revenue_by_payment_method": [ // 按支付方式的营收
        {
          "method": "string",
          "method_text": "string",
          "amount": number,
          "percentage": number,
          "growth": number
        }
      ],
      "revenue_by_store": [          // 按门店的营收
        {
          "store_id": number,
          "store_name": "string",
          "amount": number,
          "percentage": number,
          "growth": number
        }
      ],
      "daily_revenue": [             // 每日营收数据
        {
          "date": "string",
          "amount": number,
          "transactions": number
        }
      ],
      "peak_periods": [              // 营收高峰时段
        {
          "hour": number,
          "amount": number,
          "transactions": number
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.5.9 获取会员消费报表
- **URL**: `/api/v1/finance/member-consumption-report`
- **方法**: GET
- **描述**: 获取会员消费情况分析报表
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  time_range: string         // 时间范围：THIS_MONTH, LAST_MONTH, THIS_QUARTER, LAST_QUARTER, THIS_YEAR, LAST_YEAR, CUSTOM
  start_date: string         // 自定义开始日期（当time_range为CUSTOM时必填）
  end_date: string           // 自定义结束日期（当time_range为CUSTOM时必填）
  store_id: number           // 门店ID，可选，不填则查询所有门店
  level_id: number           // 会员等级ID，可选，不填则查询所有等级
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取会员消费报表成功",
    "data": {
      "time_range": {
        "start_date": "string",
        "end_date": "string"
      },
      "total_consumption": number,     // 总消费金额
      "total_consumption_growth": number, // 总消费金额增长率
      "active_members": number,        // 活跃会员数量
      "active_members_growth": number, // 活跃会员数量增长率
      "new_members": number,           // 新增会员数量
      "average_consumption": number,   // 人均消费金额
      "consumption_by_level": [        // 按会员等级的消费
        {
          "level_id": number,
          "level_name": "string",
          "members_count": number,
          "total_amount": number,
          "average_amount": number,
          "percentage": number
        }
      ],
      "consumption_by_age": [          // 按年龄段的消费
        {
          "age_group": "string",
          "members_count": number,
          "total_amount": number,
          "average_amount": number,
          "percentage": number
        }
      ],
      "consumption_by_gender": [       // 按性别的消费
        {
          "gender": "string",
          "members_count": number,
          "total_amount": number,
          "average_amount": number,
          "percentage": number
        }
      ],
      "monthly_consumption": [         // 月度消费趋势
        {
          "month": "string",
          "amount": number,
          "members_count": number
        }
      ],
      "top_consumers": [               // 消费排行榜
        {
          "member_id": number,
          "member_name": "string",
          "level_name": "string",
          "consumption_amount": number,
          "consumption_count": number
        }
      ],
      "consumption_frequency": {       // 消费频次分布
        "once": {
          "members_count": number,
          "percentage": number
        },
        "twice": {
          "members_count": number,
          "percentage": number
        },
        "three_to_five": {
          "members_count": number,
          "percentage": number
        },
        "more_than_five": {
          "members_count": number,
          "percentage": number
        }
      }
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.5.10 获取财务对账报表
- **URL**: `/api/v1/finance/reconciliation-report`
- **方法**: GET
- **描述**: 获取每日财务对账报表
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```
  date: string              // 日期，格式YYYY-MM-DD，默认为昨日
  store_id: number          // 门店ID，可选，不填则查询所有门店
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取财务对账报表成功",
    "data": {
      "date": "string",
      "store_info": {                // 门店信息（如果指定了门店ID）
        "id": number,
        "name": "string"
      },
      "summary": {                   // 总结信息
        "total_income": number,      // 总收入
        "total_expense": number,     // 总支出
        "net_income": number,        // 净收入
        "total_transactions": number, // 总交易笔数
        "average_transaction": number // 平均交易金额
      },
      "income_details": {            // 收入明细
        "by_type": [                 // 按类型分类
          {
            "type": "string",
            "type_text": "string",
            "amount": number,
            "percentage": number,
            "transactions": number
          }
        ],
        "by_payment_method": [       // 按支付方式分类
          {
            "method": "string",
            "method_text": "string",
            "amount": number,
            "percentage": number,
            "transactions": number
          }
        ]
      },
      "expense_details": {           // 支出明细
        "by_type": [                 // 按类型分类
          {
            "type": "string",
            "type_text": "string",
            "amount": number,
            "percentage": number,
            "transactions": number
          }
        ]
      },
      "hourly_distribution": [       // 按小时分布
        {
          "hour": number,
          "income": number,
          "expense": number,
          "transactions": number
        }
      ],
      "operator_statistics": [       // 操作员统计
        {
          "operator_id": number,
          "operator_name": "string",
          "transactions": number,
          "income_amount": number,
          "expense_amount": number
        }
      ],
      "cash_flow": {                 // 现金流
        "opening_cash": number,      // 期初现金
        "cash_income": number,       // 现金收入
        "cash_expense": number,      // 现金支出
        "closing_cash": number,      // 期末现金
        "cash_variance": number      // 现金差异(如有)
      },
      "transaction_summary": [       // 交易汇总列表
        {
          "transaction_type": "string",
          "transaction_type_text": "string",
          "count": number,
          "total_amount": number,
          "percentage": number
        }
      ]
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.5.11 导出财务数据
- **URL**: `/api/v1/finance/export`
- **方法**: POST
- **描述**: 导出指定条件的财务数据
- **请求头**: Authorization: Bearer {token}
- **请求参数**:
  ```json
  {
    "report_type": "string",      // 报表类型：TRANSACTION, REVENUE, MEMBER_CONSUMPTION, RECONCILIATION
    "time_range": "string",       // 时间范围
    "start_date": "string",       // 开始日期
    "end_date": "string",         // 结束日期
    "store_id": number,           // 门店ID，可选
    "file_format": "string",      // 文件格式：EXCEL, CSV, PDF
    "include_details": boolean,   // 是否包含详细数据
    "filter_conditions": {        // 其他筛选条件，可选
      "transaction_type": "string",
      "payment_method": "string",
      "member_level": number
    }
  }
  ```
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "导出任务已创建",
    "data": {
      "task_id": "string",        // 导出任务ID
      "report_type": "string",    // 报表类型
      "file_format": "string",    // 文件格式
      "estimated_time": number,   // 预计完成时间(秒)
      "status": "string",         // 任务状态：PENDING, PROCESSING, COMPLETED, FAILED
      "created_at": "string"      // 创建时间
    }
  }
  ```
- **错误码**:
  - 400: 请求参数错误
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 500: 服务器内部错误

#### 2.5.12 获取导出任务状态
- **URL**: `/api/v1/finance/export/:task_id`
- **方法**: GET
- **描述**: 获取导出任务状态和下载链接
- **请求头**: Authorization: Bearer {token}
- **请求参数**: 无
- **响应数据**:
  ```json
  {
    "code": 200,
    "message": "获取导出任务成功",
    "data": {
      "task_id": "string",
      "status": "string",         // 任务状态：PENDING, PROCESSING, COMPLETED, FAILED
      "progress": number,         // 进度百分比
      "report_type": "string",    // 报表类型
      "file_format": "string",    // 文件格式
      "file_name": "string",      // 文件名
      "file_size": number,        // 文件大小(bytes)
      "download_url": "string",   // 下载链接(只有在状态为COMPLETED时才有)
      "expires_at": "string",     // 下载链接过期时间
      "error_message": "string",  // 错误信息(只有在状态为FAILED时才有)
      "created_at": "string",     // 创建时间
      "completed_at": "string"    // 完成时间
    }
  }
  ```
- **错误码**:
  - 401: 未授权，请先登录
  - 403: 权限不足
  - 404: 导出任务不存在
  - 500: 服务器内部错误
