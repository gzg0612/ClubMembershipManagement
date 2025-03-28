# 连锁管理系统功能设计文档

## 1. 系统概述

### 1.1 项目背景
本系统旨在为射箭馆提供全面的管理解决方案，包括会员管理、商品管理、店铺管理、账号管理、财务管理、活动管理和数据回收站等功能。

### 1.2 技术栈
- 前端：Vue.js + Element UI
- 后端：Go + Gin + GORM
- 数据库：MySQL
- 缓存：Redis

## 2. 系统架构

### 2.1 整体架构
```
├── frontend/          # Vue.js前端项目
├── backend/           # Go后端项目
├── docs/              # 项目文档
└── prototype/         # 项目原型图
```

### 2.2 技术架构
- 采用前后端分离架构
- RESTful API设计
- JWT认证
- 微服务架构设计

## 3. 功能模块设计

### 3.1 控制台模块
- 关键数据统计展示（会员总数、营业额、订单数等）
- 营业额趋势图表
- 会员增长趋势图表
- 系统通知与提醒

### 3.2 会员管理模块
#### 3.2.1 核心功能
- 会员列表管理（展示所有会员信息）
- 会员查询功能（按条件筛选会员）
- 会员信息管理（增删改查）
- 会员充值（支持储值优惠活动）
  - 支持选择储值优惠活动
  - 储值优惠活动由活动管理模块配置
  - 自动计算优惠后实际充值金额
  - 区分记录储值金额和赠送金额
  - 满赠活动的赠送金额单独记入gift_balance
- 会员消费
  - 支持选择商品进行消费
  - 支持直接扣除费用
  - 消费时优先扣除赠送余额
  - 记录消费明细和消费金额
  - 自动更新会员余额
- 交易记录查询
- 会员余额管理（区分赠送金额和充值金额）
- 会员数据导出功能

#### 3.2.2 数据库设计
```sql
-- 会员表
CREATE TABLE members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL COMMENT '姓名',
    gender TINYINT NOT NULL DEFAULT 0 COMMENT '性别:0未知,1男,2女',
    phone VARCHAR(20) NOT NULL,
    card_id VARCHAR(50) NOT NULL COMMENT '会员卡号',
    store_id BIGINT NOT NULL COMMENT '所属店铺ID',
    balance DECIMAL(10,2) DEFAULT 0,
    gift_balance DECIMAL(10,2) DEFAULT 0,
    status TINYINT DEFAULT 1 COMMENT '状态:1正常,0删除',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id)
);

-- 交易记录表
CREATE TABLE transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    member_id BIGINT NOT NULL,
    store_id BIGINT NOT NULL COMMENT '消费店铺ID',
    type ENUM('RECHARGE', 'CONSUME', 'GIFT') NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    balance DECIMAL(10,2) NOT NULL,
    gift_balance DECIMAL(10,2) NOT NULL,
    item_type VARCHAR(50) COMMENT '消费项目类型',
    item_id BIGINT COMMENT '消费项目ID',
    description TEXT,
    operator_id BIGINT COMMENT '操作人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES members(id),
    FOREIGN KEY (store_id) REFERENCES stores(id)
);

-- 删除记录表
CREATE TABLE deleted_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    record_type ENUM('MEMBER', 'PRODUCT', 'ACTIVITY') NOT NULL,
    record_id BIGINT NOT NULL,
    record_data JSON NOT NULL COMMENT '删除前数据的JSON备份',
    delete_reason TEXT,
    deleted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    operator_id BIGINT COMMENT '删除操作人ID',
    operator_name VARCHAR(50) COMMENT '删除操作人姓名'
);
```

### 3.3 商品管理模块
#### 3.3.1 核心功能
- 商品信息管理
  - 商品基本信息录入（名称、描述、规格、品牌等）
  - 商品图片上传与管理
  - 商品编码管理（条形码/二维码）
  - 商品批量导入导出
- 商品分类管理
  - 多级分类结构设置
  - 分类添加/编辑/删除
  - 分类排序与显示控制
  - 分类关联商品管理
- 库存管理
  - 库存实时监控
  - 库存预警设置
  - 入库/出库记录
  - 库存盘点功能
  - 库存调拨管理（跨店铺）
- 价格管理
  - 基础售价设置
  - 会员价格设置
  - 批发价格设置
  - 促销价格管理
  - 价格变更历史记录
- 商品上下架管理
  - 商品状态控制（上架/下架/缺货）
  - 批量上下架操作
  - 定时上下架设置
  - 商品可见性控制

#### 3.3.2 数据库设计
```sql
-- 商品表
CREATE TABLE products (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(50) NOT NULL COMMENT '商品编码',
    barcode VARCHAR(50) COMMENT '条形码/二维码',
    name VARCHAR(100) NOT NULL COMMENT '商品名称',
    description TEXT COMMENT '商品描述',
    brand VARCHAR(50) COMMENT '品牌',
    specification VARCHAR(100) COMMENT '规格',
    unit VARCHAR(20) COMMENT '单位',
    cost_price DECIMAL(10,2) COMMENT '成本价',
    retail_price DECIMAL(10,2) NOT NULL COMMENT '零售价',
    member_price DECIMAL(10,2) COMMENT '会员价',
    wholesale_price DECIMAL(10,2) COMMENT '批发价',
    category_id BIGINT NOT NULL COMMENT '分类ID',
    stock_quantity INT DEFAULT 0 COMMENT '库存数量',
    stock_warning INT DEFAULT 10 COMMENT '库存预警值',
    status TINYINT DEFAULT 1 COMMENT '状态:1上架,0下架,2缺货',
    is_deleted TINYINT DEFAULT 0 COMMENT '是否删除:0否,1是',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code(code),
    INDEX idx_category(category_id),
    INDEX idx_status(status)
);

-- 商品分类表
CREATE TABLE product_categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL COMMENT '分类名称',
    parent_id BIGINT DEFAULT 0 COMMENT '父分类ID，0表示顶级分类',
    level INT DEFAULT 1 COMMENT '分类层级',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    is_visible TINYINT DEFAULT 1 COMMENT '是否可见:1是,0否',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent(parent_id)
);

-- 商品图片表
CREATE TABLE product_images (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT NOT NULL COMMENT '商品ID',
    image_url VARCHAR(255) NOT NULL COMMENT '图片URL',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    is_main TINYINT DEFAULT 0 COMMENT '是否主图:1是,0否',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id),
    INDEX idx_product(product_id)
);

-- 库存记录表
CREATE TABLE inventory_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT NOT NULL COMMENT '商品ID',
    type ENUM('IN', 'OUT', 'CHECK', 'TRANSFER') NOT NULL COMMENT '类型:入库,出库,盘点,调拨',
    quantity INT NOT NULL COMMENT '数量',
    before_quantity INT NOT NULL COMMENT '操作前库存',
    after_quantity INT NOT NULL COMMENT '操作后库存',
    reason VARCHAR(100) COMMENT '操作原因',
    remark TEXT COMMENT '备注',
    operator_id BIGINT COMMENT '操作人ID',
    store_id BIGINT COMMENT '店铺ID',
    target_store_id BIGINT COMMENT '目标店铺ID(调拨时使用)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id),
    INDEX idx_product(product_id),
    INDEX idx_store(store_id)
);

-- 价格变更历史表
CREATE TABLE price_change_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT NOT NULL COMMENT '商品ID',
    price_type ENUM('RETAIL', 'MEMBER', 'WHOLESALE', 'COST') NOT NULL COMMENT '价格类型',
    old_price DECIMAL(10,2) NOT NULL COMMENT '旧价格',
    new_price DECIMAL(10,2) NOT NULL COMMENT '新价格',
    reason VARCHAR(100) COMMENT '变更原因',
    operator_id BIGINT COMMENT '操作人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id),
    INDEX idx_product(product_id)
);

-- 促销活动表
CREATE TABLE promotions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '活动名称',
    description TEXT COMMENT '活动描述',
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME NOT NULL COMMENT '结束时间',
    status TINYINT DEFAULT 1 COMMENT '状态:1生效,0未生效',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_time(start_time, end_time)
);

-- 促销商品关联表
CREATE TABLE promotion_products (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    promotion_id BIGINT NOT NULL COMMENT '促销活动ID',
    product_id BIGINT NOT NULL COMMENT '商品ID',
    promotion_price DECIMAL(10,2) NOT NULL COMMENT '促销价格',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (promotion_id) REFERENCES promotions(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    UNIQUE KEY uk_prom_prod(promotion_id, product_id)
);

-- 批量导入导出记录表
CREATE TABLE product_batch_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    type ENUM('IMPORT', 'EXPORT') NOT NULL COMMENT '类型:导入,导出',
    file_name VARCHAR(255) COMMENT '文件名称',
    file_url VARCHAR(255) COMMENT '文件URL',
    total_count INT DEFAULT 0 COMMENT '总记录数',
    success_count INT DEFAULT 0 COMMENT '成功数',
    fail_count INT DEFAULT 0 COMMENT '失败数',
    status TINYINT DEFAULT 0 COMMENT '状态:0处理中,1成功,2失败',
    error_message TEXT COMMENT '错误信息',
    operator_id BIGINT COMMENT '操作人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 3.4 店铺管理模块
#### 3.4.1 总店管理
- 总店基本信息配置
  - 总店名称、地址、联系方式、营业时间设置
  - 总店LOGO和品牌形象配置
  - 总店负责人和管理团队设置
  - 总店经营许可证和资质管理
  - 总店财务账户配置
- 分店管理
  - 分店信息添加/编辑/删除
    - 分店基础信息维护（名称、地址、联系方式、营业时间）
    - 分店负责人和员工管理
    - 分店场地和设备配置管理
    - 分店经营状态监控
  - 分店状态管理
    - 分店营业状态切换（营业中、暂停营业、关闭等）
    - 分店紧急状态处理机制
    - 分店营业时间临时调整
    - 分店运营状态实时监控
  - 分店权限配置
    - 分店管理人员权限分配
    - 分店财务操作权限设置
    - 分店数据查看和导出权限控制
    - 分店系统功能模块访问控制
  - 分店数据统计
    - 分店营业额实时统计和对比
    - 分店会员数量和活跃度分析
    - 分店设备使用率和场地预约情况
    - 分店员工绩效数据统计
    - 分店运营成本分析
- 连锁店铺管理
  - 店铺资源调配
    - 人力资源跨店调配管理
    - 设备资源共享和调配机制
    - 教练资源优化分配
    - 培训资源统一调度
  - 店铺间库存调拨
    - 商品库存跨店调拨申请和审批
    - 库存调拨记录追踪和管理
    - 库存预警和自动调拨建议
    - 调拨物流状态跟踪
  - 跨店会员服务
    - 会员信息全店共享机制
    - 会员跨店消费和积分统一管理
    - 会员卡在不同分店的使用权限设置
    - 会员跨店预约和服务体验优化
  - 连锁店数据汇总
    - 全店经营数据实时汇总看板
    - 多维度数据分析和可视化展示
    - 连锁店经营趋势分析和预测
    - 门店间绩效对比和评估体系
    - 连锁店运营效率分析报告

#### 3.4.2 分店管理
- 店铺信息管理
  - 分店基本信息维护（名称、地址、电话、营业时间）
  - 分店负责人和联系方式管理
  - 分店营业状态实时更新
  - 分店环境照片和平面图管理
  - 分店特色服务和设施标记
- 场地管理
  - 射箭场地类型划分（初级区、中级区、高级区）
  - 场地容量和使用状态监控
  - 场地预约时段设置和管理
  - 场地维护计划和记录
  - 场地使用率统计和分析
- 设备管理
  - 弓箭装备库存管理（弓、箭、护具等）
  - 设备编号和二维码管理系统
  - 设备借用和归还流程管理
  - 设备损耗统计和更新计划
  - 设备分类和等级管理
- 设备维护记录
  - 定期设备检查计划制定
  - 设备维修记录和追踪
  - 设备维护成本统计
  - 设备使用寿命监控
  - 设备安全检查记录和合规性管理
- 预约管理
  - 场地在线预约系统
  - 教练预约和排班管理
  - 团体课程和活动预约
  - 预约提醒和确认机制
  - 预约数据分析和高峰期管理
  - 会员预约优先级设置

#### 3.4.3 数据库设计
```sql
-- 总店信息表
CREATE TABLE headquarters (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '总店名称',
    logo_url VARCHAR(255) COMMENT '总店LOGO',
    address VARCHAR(200) NOT NULL COMMENT '地址',
    phone VARCHAR(20) NOT NULL COMMENT '联系电话',
    email VARCHAR(100) COMMENT '邮箱',
    business_hours VARCHAR(100) COMMENT '营业时间',
    manager_name VARCHAR(50) COMMENT '负责人姓名',
    manager_phone VARCHAR(20) COMMENT '负责人电话',
    business_license VARCHAR(50) COMMENT '营业执照编号',
    license_image_url VARCHAR(255) COMMENT '营业执照图片',
    bank_account VARCHAR(50) COMMENT '财务账户',
    bank_name VARCHAR(100) COMMENT '开户行',
    tax_number VARCHAR(50) COMMENT '税号',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 分店信息表
CREATE TABLE stores (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(20) NOT NULL COMMENT '店铺编码',
    name VARCHAR(100) NOT NULL COMMENT '分店名称',
    address VARCHAR(200) NOT NULL COMMENT '地址',
    phone VARCHAR(20) NOT NULL COMMENT '联系电话',
    email VARCHAR(100) COMMENT '邮箱',
    business_hours VARCHAR(100) COMMENT '营业时间',
    manager_id BIGINT COMMENT '负责人ID',
    manager_name VARCHAR(50) COMMENT '负责人姓名',
    manager_phone VARCHAR(20) COMMENT '负责人电话',
    status ENUM('OPEN', 'CLOSED', 'SUSPENDED', 'EMERGENCY') DEFAULT 'OPEN' COMMENT '营业状态',
    area DECIMAL(10,2) COMMENT '店面面积(平方米)',
    open_date DATE COMMENT '开业日期',
    description TEXT COMMENT '店铺描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code(code),
    INDEX idx_status(status)
);

-- 店铺图片表
CREATE TABLE store_images (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    image_url VARCHAR(255) NOT NULL COMMENT '图片URL',
    image_type ENUM('ENVIRONMENT', 'FLOORPLAN', 'LICENSE', 'OTHER') COMMENT '图片类型',
    description VARCHAR(100) COMMENT '图片描述',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id)
);

-- 店铺特色服务表
CREATE TABLE store_services (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    service_name VARCHAR(50) NOT NULL COMMENT '服务名称',
    description TEXT COMMENT '服务描述',
    is_active TINYINT DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id)
);

-- 场地表
CREATE TABLE venues (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    name VARCHAR(50) NOT NULL COMMENT '场地名称',
    code VARCHAR(20) NOT NULL COMMENT '场地编码',
    type ENUM('BEGINNER', 'INTERMEDIATE', 'ADVANCED', 'MULTIPURPOSE') NOT NULL COMMENT '场地类型',
    capacity INT NOT NULL COMMENT '容纳人数',
    area DECIMAL(10,2) COMMENT '场地面积(平方米)',
    shooting_distance INT COMMENT '射箭距离(米)',
    target_count INT COMMENT '靶位数量',
    status ENUM('AVAILABLE', 'OCCUPIED', 'MAINTENANCE', 'CLOSED') DEFAULT 'AVAILABLE' COMMENT '场地状态',
    description TEXT COMMENT '场地描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_status(status)
);

-- 场地维护记录表
CREATE TABLE venue_maintenance (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    venue_id BIGINT NOT NULL COMMENT '场地ID',
    maintenance_type VARCHAR(50) NOT NULL COMMENT '维护类型',
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME COMMENT '结束时间',
    status ENUM('PLANNED', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED') DEFAULT 'PLANNED' COMMENT '状态',
    cost DECIMAL(10,2) COMMENT '维护成本',
    operator_id BIGINT COMMENT '执行人ID',
    operator_name VARCHAR(50) COMMENT '执行人姓名',
    description TEXT COMMENT '维护说明',
    result TEXT COMMENT '维护结果',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (venue_id) REFERENCES venues(id),
    INDEX idx_venue(venue_id),
    INDEX idx_time(start_time, end_time)
);

-- 设备表
CREATE TABLE equipment (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    code VARCHAR(20) NOT NULL COMMENT '设备编码',
    name VARCHAR(50) NOT NULL COMMENT '设备名称',
    type ENUM('BOW', 'ARROW', 'TARGET', 'PROTECTION', 'ELECTRONIC', 'OTHER') NOT NULL COMMENT '设备类型',
    category VARCHAR(50) COMMENT '分类',
    brand VARCHAR(50) COMMENT '品牌',
    model VARCHAR(50) COMMENT '型号',
    level ENUM('BEGINNER', 'INTERMEDIATE', 'ADVANCED', 'PROFESSIONAL') COMMENT '等级',
    purchase_date DATE COMMENT '购买日期',
    purchase_price DECIMAL(10,2) COMMENT '购买价格',
    expected_lifetime INT COMMENT '预计使用寿命(月)',
    status ENUM('AVAILABLE', 'IN_USE', 'MAINTENANCE', 'DAMAGED', 'RETIRED') DEFAULT 'AVAILABLE' COMMENT '设备状态',
    qrcode_url VARCHAR(255) COMMENT '二维码URL',
    location VARCHAR(50) COMMENT '存放位置',
    description TEXT COMMENT '设备描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_status(status),
    INDEX idx_code(code)
);

-- 设备借用记录表
CREATE TABLE equipment_lending (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    equipment_id BIGINT NOT NULL COMMENT '设备ID',
    member_id BIGINT COMMENT '会员ID',
    member_name VARCHAR(50) COMMENT '会员姓名',
    borrow_time DATETIME NOT NULL COMMENT '借出时间',
    expected_return_time DATETIME COMMENT '预计归还时间',
    actual_return_time DATETIME COMMENT '实际归还时间',
    status ENUM('BORROWED', 'RETURNED', 'OVERDUE', 'DAMAGED', 'LOST') DEFAULT 'BORROWED' COMMENT '状态',
    deposit DECIMAL(10,2) COMMENT '押金',
    damage_description TEXT COMMENT '损坏描述',
    damage_fee DECIMAL(10,2) COMMENT '损坏赔偿费',
    operator_id BIGINT COMMENT '操作员ID',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (equipment_id) REFERENCES equipment(id),
    INDEX idx_equipment(equipment_id),
    INDEX idx_member(member_id),
    INDEX idx_status(status)
);

-- 设备维护记录表
CREATE TABLE equipment_maintenance (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    equipment_id BIGINT NOT NULL COMMENT '设备ID',
    maintenance_type VARCHAR(50) NOT NULL COMMENT '维护类型',
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME COMMENT '结束时间',
    status ENUM('PLANNED', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED') DEFAULT 'PLANNED' COMMENT '状态',
    cost DECIMAL(10,2) COMMENT '维护成本',
    operator_id BIGINT COMMENT '执行人ID',
    operator_name VARCHAR(50) COMMENT '执行人姓名',
    description TEXT COMMENT '维护说明',
    result TEXT COMMENT '维护结果',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (equipment_id) REFERENCES equipment(id),
    INDEX idx_equipment(equipment_id),
    INDEX idx_time(start_time, end_time)
);

-- 设备调拨记录表
CREATE TABLE equipment_transfer (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    equipment_id BIGINT NOT NULL COMMENT '设备ID',
    from_store_id BIGINT NOT NULL COMMENT '源店铺ID',
    to_store_id BIGINT NOT NULL COMMENT '目标店铺ID',
    transfer_time DATETIME NOT NULL COMMENT '调拨时间',
    status ENUM('PENDING', 'IN_TRANSIT', 'COMPLETED', 'CANCELLED') DEFAULT 'PENDING' COMMENT '状态',
    reason TEXT COMMENT '调拨原因',
    applicant_id BIGINT COMMENT '申请人ID',
    approver_id BIGINT COMMENT '审批人ID',
    logistics_info VARCHAR(100) COMMENT '物流信息',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (equipment_id) REFERENCES equipment(id),
    FOREIGN KEY (from_store_id) REFERENCES stores(id),
    FOREIGN KEY (to_store_id) REFERENCES stores(id),
    INDEX idx_equipment(equipment_id),
    INDEX idx_store(from_store_id, to_store_id),
    INDEX idx_status(status)
);

-- 商品库存调拨记录表
CREATE TABLE product_transfer (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT NOT NULL COMMENT '商品ID',
    from_store_id BIGINT NOT NULL COMMENT '源店铺ID',
    to_store_id BIGINT NOT NULL COMMENT '目标店铺ID',
    quantity INT NOT NULL COMMENT '调拨数量',
    transfer_time DATETIME NOT NULL COMMENT '调拨时间',
    status ENUM('PENDING', 'IN_TRANSIT', 'COMPLETED', 'CANCELLED') DEFAULT 'PENDING' COMMENT '状态',
    reason TEXT COMMENT '调拨原因',
    applicant_id BIGINT COMMENT '申请人ID',
    approver_id BIGINT COMMENT '审批人ID',
    logistics_info VARCHAR(100) COMMENT '物流信息',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (from_store_id) REFERENCES stores(id),
    FOREIGN KEY (to_store_id) REFERENCES stores(id),
    INDEX idx_product(product_id),
    INDEX idx_store(from_store_id, to_store_id),
    INDEX idx_status(status)
);

-- 预约时段设置表
CREATE TABLE booking_time_slots (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    venue_id BIGINT COMMENT '场地ID，NULL表示适用所有场地',
    day_of_week TINYINT NOT NULL COMMENT '星期几(1-7)',
    start_time TIME NOT NULL COMMENT '开始时间',
    end_time TIME NOT NULL COMMENT '结束时间',
    max_capacity INT COMMENT '最大预约人数',
    price DECIMAL(10,2) COMMENT '标准价格',
    member_price DECIMAL(10,2) COMMENT '会员价格',
    is_active TINYINT DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (venue_id) REFERENCES venues(id),
    INDEX idx_store(store_id),
    INDEX idx_venue(venue_id),
    INDEX idx_day(day_of_week)
);

-- 预约记录表
CREATE TABLE bookings (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    booking_no VARCHAR(20) NOT NULL COMMENT '预约编号',
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    venue_id BIGINT NOT NULL COMMENT '场地ID',
    member_id BIGINT COMMENT '会员ID',
    member_name VARCHAR(50) NOT NULL COMMENT '预约人姓名',
    member_phone VARCHAR(20) NOT NULL COMMENT '预约人电话',
    booking_date DATE NOT NULL COMMENT '预约日期',
    start_time TIME NOT NULL COMMENT '开始时间',
    end_time TIME NOT NULL COMMENT '结束时间',
    people_count INT DEFAULT 1 COMMENT '人数',
    coach_id BIGINT COMMENT '教练ID',
    status ENUM('PENDING', 'CONFIRMED', 'CHECKED_IN', 'COMPLETED', 'CANCELLED', 'NO_SHOW') DEFAULT 'PENDING' COMMENT '状态',
    price DECIMAL(10,2) COMMENT '价格',
    payment_status ENUM('UNPAID', 'PAID', 'REFUNDED') DEFAULT 'UNPAID' COMMENT '支付状态',
    payment_time DATETIME COMMENT '支付时间',
    payment_method VARCHAR(20) COMMENT '支付方式',
    cancel_reason VARCHAR(100) COMMENT '取消原因',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    FOREIGN KEY (venue_id) REFERENCES venues(id),
    UNIQUE KEY uk_booking_no(booking_no),
    INDEX idx_store(store_id),
    INDEX idx_venue(venue_id),
    INDEX idx_member(member_id),
    INDEX idx_coach(coach_id),
    INDEX idx_date(booking_date),
    INDEX idx_status(status)
);

-- 教练排班表
CREATE TABLE coach_schedules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    coach_id BIGINT NOT NULL COMMENT '教练ID',
    schedule_date DATE NOT NULL COMMENT '排班日期',
    start_time TIME NOT NULL COMMENT '开始时间',
    end_time TIME NOT NULL COMMENT '结束时间',
    status ENUM('SCHEDULED', 'CANCELLED', 'COMPLETED') DEFAULT 'SCHEDULED' COMMENT '状态',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_coach(coach_id),
    INDEX idx_date(schedule_date)
);

-- 店铺营业数据统计表
CREATE TABLE store_statistics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    statistic_date DATE NOT NULL COMMENT '统计日期',
    visitor_count INT DEFAULT 0 COMMENT '访客数',
    new_member_count INT DEFAULT 0 COMMENT '新增会员数',
    booking_count INT DEFAULT 0 COMMENT '预约数',
    completed_booking_count INT DEFAULT 0 COMMENT '已完成预约数',
    cancelled_booking_count INT DEFAULT 0 COMMENT '取消预约数',
    total_revenue DECIMAL(12,2) DEFAULT 0 COMMENT '总营收',
    venue_revenue DECIMAL(12,2) DEFAULT 0 COMMENT '场地收入',
    coaching_revenue DECIMAL(12,2) DEFAULT 0 COMMENT '教练收入',
    retail_revenue DECIMAL(12,2) DEFAULT 0 COMMENT '零售收入',
    other_revenue DECIMAL(12,2) DEFAULT 0 COMMENT '其他收入',
    operating_cost DECIMAL(12,2) DEFAULT 0 COMMENT '运营成本',
    profit DECIMAL(12,2) DEFAULT 0 COMMENT '利润',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    UNIQUE KEY uk_store_date(store_id, statistic_date),
    INDEX idx_date(statistic_date)
);

-- 店铺状态变更记录表
CREATE TABLE store_status_changes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    previous_status VARCHAR(20) COMMENT '之前状态',
    new_status VARCHAR(20) NOT NULL COMMENT '新状态',
    change_reason TEXT COMMENT '变更原因',
    start_time DATETIME NOT NULL COMMENT '生效时间',
    end_time DATETIME COMMENT '结束时间',
    operator_id BIGINT COMMENT '操作人ID',
    operator_name VARCHAR(50) COMMENT '操作人姓名',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_time(start_time, end_time)
);

-- 人力资源调配记录表
CREATE TABLE staff_transfer (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    staff_id BIGINT NOT NULL COMMENT '员工ID',
    from_store_id BIGINT NOT NULL COMMENT '原店铺ID',
    to_store_id BIGINT NOT NULL COMMENT '目标店铺ID',
    transfer_type ENUM('TEMPORARY', 'PERMANENT') NOT NULL COMMENT '调配类型',
    start_date DATE NOT NULL COMMENT '开始日期',
    end_date DATE COMMENT '结束日期(临时调配需要)',
    reason TEXT COMMENT '调配原因',
    position VARCHAR(50) COMMENT '调配后职位',
    status ENUM('PENDING', 'APPROVED', 'COMPLETED', 'REJECTED', 'CANCELLED') DEFAULT 'PENDING' COMMENT '状态',
    applicant_id BIGINT COMMENT '申请人ID',
    approver_id BIGINT COMMENT '审批人ID',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (from_store_id) REFERENCES stores(id),
    FOREIGN KEY (to_store_id) REFERENCES stores(id),
    INDEX idx_staff(staff_id),
    INDEX idx_store(from_store_id, to_store_id),
    INDEX idx_status(status)
);

-- 店铺权限配置表
CREATE TABLE store_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    module_code VARCHAR(50) NOT NULL COMMENT '模块代码',
    permission_type ENUM('VIEW', 'EDIT', 'APPROVE', 'EXPORT', 'ALL') NOT NULL COMMENT '权限类型',
    is_granted TINYINT DEFAULT 1 COMMENT '是否授权',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    UNIQUE KEY uk_store_role_module(store_id, role_id, module_code)
);
```

### 3.5 账号管理模块
#### 3.5.1 用户管理
- 用户信息管理
  - 用户基本信息维护（姓名、手机号、邮箱等）
  - 用户状态管理（启用/禁用/锁定）
  - 用户与店铺关联管理
  - 用户头像上传与管理
  - 用户账号有效期设置
- 用户分组管理
  - 用户部门划分
  - 用户职位管理
  - 用户团队分组
  - 组织架构管理
- 用户权限分配
  - 用户角色分配
  - 用户特殊权限设置
  - 用户数据访问范围控制
  - 用户功能权限个性化配置

#### 3.5.2 角色管理
- 角色创建与配置
  - 角色基本信息设置
  - 角色权限矩阵配置
  - 角色层级关系设置
  - 角色模板管理
- 角色权限分配
  - 功能权限分配
  - 数据权限分配
  - 操作权限分配
  - 审批权限分配
- 角色用户管理
  - 角色下用户查看
  - 批量分配用户到角色
  - 角色用户移除
  - 角色用户权限继承关系

#### 3.5.3 权限管理
- 权限项配置
  - 系统功能权限定义
  - 数据访问权限定义
  - 操作权限定义
  - 权限依赖关系配置
- 权限分配管理
  - 基于角色的权限分配
  - 基于用户的权限分配
  - 基于部门的权限分配
  - 权限继承与覆盖规则
- 权限审计
  - 权限分配记录审计
  - 权限变更历史追踪
  - 权限使用情况分析
  - 敏感权限监控

#### 3.5.4 操作日志
- 用户登录日志
  - 登录时间记录
  - 登录IP与设备信息
  - 登录状态跟踪
  - 异常登录预警
- 操作行为日志
  - 关键操作记录
  - 数据变更记录
  - 操作时间与操作人记录
  - 操作内容详情记录
- 日志查询与分析
  - 多条件日志筛选
  - 日志导出功能
  - 操作行为分析
  - 异常操作预警

#### 3.5.5 密码安全
- 密码重置
  - 管理员密码重置
  - 用户自助密码重置
  - 密码重置验证流程
  - 密码重置通知机制
- 密码策略管理
  - 密码复杂度要求设置
  - 密码有效期设置
  - 密码历史记录检查
  - 密码锁定策略设置
- 安全认证管理
  - 双因素认证配置
  - 登录安全策略设置
  - 登录异常处理机制
  - 会话超时设置

#### 3.5.6 数据库设计

```sql
-- 用户表
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    password VARCHAR(100) NOT NULL COMMENT '密码(加密存储)',
    name VARCHAR(50) NOT NULL COMMENT '姓名',
    phone VARCHAR(20) COMMENT '手机号',
    email VARCHAR(100) COMMENT '邮箱',
    avatar_url VARCHAR(255) COMMENT '头像URL',
    status ENUM('ACTIVE', 'DISABLED', 'LOCKED') DEFAULT 'ACTIVE' COMMENT '状态',
    department_id BIGINT COMMENT '部门ID',
    position VARCHAR(50) COMMENT '职位',
    employee_id VARCHAR(50) COMMENT '员工编号',
    last_login_time DATETIME COMMENT '最后登录时间',
    last_login_ip VARCHAR(50) COMMENT '最后登录IP',
    account_expire_time DATETIME COMMENT '账号过期时间',
    password_expire_time DATETIME COMMENT '密码过期时间',
    password_update_time DATETIME COMMENT '密码更新时间',
    login_failed_count INT DEFAULT 0 COMMENT '登录失败次数',
    is_admin TINYINT DEFAULT 0 COMMENT '是否管理员',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    INDEX idx_username(username),
    INDEX idx_phone(phone),
    INDEX idx_email(email),
    INDEX idx_department(department_id),
    INDEX idx_status(status)
);

-- 部门表
CREATE TABLE departments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL COMMENT '部门名称',
    code VARCHAR(50) COMMENT '部门编码',
    parent_id BIGINT DEFAULT 0 COMMENT '父部门ID',
    level INT DEFAULT 1 COMMENT '层级',
    sort_order INT DEFAULT 0 COMMENT '排序',
    manager_id BIGINT COMMENT '部门负责人ID',
    description VARCHAR(200) COMMENT '部门描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_parent(parent_id),
    INDEX idx_manager(manager_id)
);

-- 用户团队表
CREATE TABLE teams (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL COMMENT '团队名称',
    leader_id BIGINT COMMENT '团队负责人ID',
    description VARCHAR(200) COMMENT '团队描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_leader(leader_id)
);

-- 用户团队关联表
CREATE TABLE user_team_relation (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    team_id BIGINT NOT NULL COMMENT '团队ID',
    join_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '加入时间',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (team_id) REFERENCES teams(id),
    UNIQUE KEY uk_user_team(user_id, team_id)
);

-- 用户店铺关联表
CREATE TABLE user_store_relation (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    is_primary TINYINT DEFAULT 0 COMMENT '是否主店铺',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE KEY uk_user_store(user_id, store_id)
);

-- 角色表
CREATE TABLE roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL COMMENT '角色名称',
    code VARCHAR(50) NOT NULL COMMENT '角色编码',
    description VARCHAR(200) COMMENT '角色描述',
    is_system TINYINT DEFAULT 0 COMMENT '是否系统角色',
    is_default TINYINT DEFAULT 0 COMMENT '是否默认角色',
    parent_id BIGINT DEFAULT 0 COMMENT '父角色ID',
    level INT DEFAULT 1 COMMENT '角色层级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    UNIQUE KEY uk_code(code),
    INDEX idx_parent(parent_id)
);

-- 用户角色关联表
CREATE TABLE user_role_relation (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    grant_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '授予时间',
    grant_by BIGINT COMMENT '授予人ID',
    expire_time DATETIME COMMENT '过期时间',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (role_id) REFERENCES roles(id),
    UNIQUE KEY uk_user_role(user_id, role_id)
);

-- 权限表
CREATE TABLE permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) NOT NULL COMMENT '权限编码',
    type ENUM('MENU', 'OPERATION', 'DATA', 'FIELD') NOT NULL COMMENT '权限类型',
    parent_id BIGINT DEFAULT 0 COMMENT '父权限ID',
    path VARCHAR(200) COMMENT '路径(菜单类型)',
    component VARCHAR(100) COMMENT '组件(菜单类型)',
    redirect VARCHAR(200) COMMENT '重定向(菜单类型)',
    icon VARCHAR(100) COMMENT '图标(菜单类型)',
    sort_order INT DEFAULT 0 COMMENT '排序(菜单类型)',
    is_hidden TINYINT DEFAULT 0 COMMENT '是否隐藏(菜单类型)',
    description VARCHAR(200) COMMENT '权限描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code(code),
    INDEX idx_parent(parent_id),
    INDEX idx_type(type)
);

-- 角色权限关联表
CREATE TABLE role_permission_relation (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL COMMENT '角色ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id),
    UNIQUE KEY uk_role_permission(role_id, permission_id)
);

-- 用户特殊权限表（用户级别的权限覆盖）
CREATE TABLE user_permission_relation (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    type ENUM('GRANT', 'DENY') NOT NULL COMMENT '授予类型',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    expire_time DATETIME COMMENT '过期时间',
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id),
    UNIQUE KEY uk_user_permission(user_id, permission_id)
);

-- 数据权限表
CREATE TABLE data_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT COMMENT '角色ID',
    user_id BIGINT COMMENT '用户ID',
    resource_type VARCHAR(50) NOT NULL COMMENT '资源类型',
    data_scope ENUM('ALL', 'DEPARTMENT', 'DEPARTMENT_AND_CHILD', 'SELF', 'CUSTOM') NOT NULL COMMENT '数据范围',
    custom_scope TEXT COMMENT '自定义范围(JSON格式)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    INDEX idx_role(role_id),
    INDEX idx_user(user_id),
    INDEX idx_resource_type(resource_type)
);

-- 登录日志表
CREATE TABLE login_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    login_type VARCHAR(20) COMMENT '登录类型',
    login_status ENUM('SUCCESS', 'FAILURE') NOT NULL COMMENT '登录状态',
    ip_address VARCHAR(50) COMMENT 'IP地址',
    location VARCHAR(100) COMMENT '登录地点',
    device VARCHAR(100) COMMENT '设备信息',
    browser VARCHAR(100) COMMENT '浏览器信息',
    os VARCHAR(100) COMMENT '操作系统',
    failure_reason VARCHAR(200) COMMENT '失败原因',
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    INDEX idx_user(user_id),
    INDEX idx_username(username),
    INDEX idx_login_time(login_time),
    INDEX idx_status(login_status)
);

-- 操作日志表
CREATE TABLE operation_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT COMMENT '用户ID',
    username VARCHAR(50) NOT NULL COMMENT '用户名',
    module VARCHAR(50) COMMENT '操作模块',
    operation VARCHAR(50) NOT NULL COMMENT '操作类型',
    method VARCHAR(100) COMMENT '请求方法',
    request_url VARCHAR(255) COMMENT '请求URL',
    request_method VARCHAR(10) COMMENT '请求方式',
    request_params TEXT COMMENT '请求参数',
    response_result TEXT COMMENT '返回结果',
    error_message TEXT COMMENT '错误消息',
    ip_address VARCHAR(50) COMMENT 'IP地址',
    execution_time BIGINT COMMENT '执行时长(毫秒)',
    status TINYINT DEFAULT 1 COMMENT '操作状态:1成功,0失败',
    operation_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '操作时间',
    INDEX idx_user(user_id),
    INDEX idx_username(username),
    INDEX idx_module(module),
    INDEX idx_operation(operation),
    INDEX idx_operation_time(operation_time),
    INDEX idx_status(status)
);

-- 密码历史表
CREATE TABLE password_history (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    password VARCHAR(100) NOT NULL COMMENT '历史密码(加密)',
    change_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
    FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX idx_user(user_id),
    INDEX idx_change_time(change_time)
);

-- 密码重置记录表
CREATE TABLE password_reset (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    reset_token VARCHAR(100) NOT NULL COMMENT '重置令牌',
    token_expire_time DATETIME NOT NULL COMMENT '令牌过期时间',
    reset_by BIGINT COMMENT '重置人ID',
    is_used TINYINT DEFAULT 0 COMMENT '是否已使用',
    used_time DATETIME COMMENT '使用时间',
    request_ip VARCHAR(50) COMMENT '请求IP',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    INDEX idx_user(user_id),
    INDEX idx_token(reset_token),
    INDEX idx_expire_time(token_expire_time)
);

-- 安全策略表
CREATE TABLE security_policies (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    policy_type VARCHAR(50) NOT NULL COMMENT '策略类型',
    policy_name VARCHAR(50) NOT NULL COMMENT '策略名称',
    policy_value TEXT NOT NULL COMMENT '策略值(JSON格式)',
    is_enabled TINYINT DEFAULT 1 COMMENT '是否启用',
    description VARCHAR(200) COMMENT '描述',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT COMMENT '更新人ID',
    UNIQUE KEY uk_type_name(policy_type, policy_name)
);

-- 权限变更记录表
CREATE TABLE permission_change_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    change_type ENUM('GRANT', 'REVOKE', 'MODIFY') NOT NULL COMMENT '变更类型',
    target_type ENUM('USER', 'ROLE') NOT NULL COMMENT '目标类型',
    target_id BIGINT NOT NULL COMMENT '目标ID',
    permission_id BIGINT COMMENT '权限ID',
    role_id BIGINT COMMENT '角色ID',
    before_value TEXT COMMENT '变更前值',
    after_value TEXT COMMENT '变更后值',
    operator_id BIGINT COMMENT '操作人ID',
    operator_name VARCHAR(50) COMMENT '操作人姓名',
    change_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '变更时间',
    remark VARCHAR(200) COMMENT '备注',
    INDEX idx_target(target_type, target_id),
    INDEX idx_permission(permission_id),
    INDEX idx_role(role_id),
    INDEX idx_change_time(change_time)
); 
```

### 3.6 财务管理模块
#### 3.6.1 总体财务管理
- 总营业额统计
- 总支出统计
- 总利润分析
- 现金流量表
- 资产负债表
- 财务趋势分析
- 财务预警指标

#### 3.6.2 店铺维度财务管理
- 店铺营收统计
  - 单店日/周/月/年营收报表
  - 店铺间营收对比
  - 营收构成分析
- 店铺支出统计
  - 人工成本统计
  - 运营成本统计
  - 设备折旧统计
- 店铺利润分析
  - 毛利率分析
  - 净利率分析
  - 成本收益分析
- 店铺财务报表
  - 店铺财务月报
  - 店铺财务季报
  - 店铺财务年报
#### 3.6.2 数据库设计

```sql
-- 收入记录表
CREATE TABLE income_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT COMMENT '店铺ID，NULL表示总部收入',
    income_type ENUM('MEMBER_RECHARGE', 'VENUE_RENTAL', 'PRODUCT_SALES', 'COACHING_FEE', 'ACTIVITY_FEE', 'OTHER') NOT NULL COMMENT '收入类型',
    amount DECIMAL(12,2) NOT NULL COMMENT '金额',
    related_id BIGINT COMMENT '关联ID(会员ID/订单ID等)',
    related_type VARCHAR(50) COMMENT '关联类型',
    payment_method ENUM('CASH', 'WECHAT', 'ALIPAY', 'BANKCARD', 'MEMBER_BALANCE', 'OTHER') NOT NULL COMMENT '支付方式',
    transaction_no VARCHAR(100) COMMENT '交易单号',
    remark TEXT COMMENT '备注',
    record_time DATETIME NOT NULL COMMENT '记录时间',
    operator_id BIGINT COMMENT '操作人ID',
    fiscal_year INT NOT NULL COMMENT '财年',
    fiscal_month INT NOT NULL COMMENT '财月',
    is_settled TINYINT DEFAULT 0 COMMENT '是否已结算',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_time(record_time),
    INDEX idx_type(income_type),
    INDEX idx_fiscal(fiscal_year, fiscal_month)
);

-- 支出记录表
CREATE TABLE expense_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT COMMENT '店铺ID，NULL表示总部支出',
    expense_type ENUM('RENT', 'SALARY', 'UTILITIES', 'EQUIPMENT', 'INVENTORY', 'MARKETING', 'MAINTENANCE', 'TAX', 'OTHER') NOT NULL COMMENT '支出类型',
    amount DECIMAL(12,2) NOT NULL COMMENT '金额',
    payment_method ENUM('CASH', 'WECHAT', 'ALIPAY', 'BANKCARD', 'OTHER') NOT NULL COMMENT '支付方式',
    transaction_no VARCHAR(100) COMMENT '交易单号',
    recipient VARCHAR(100) COMMENT '收款方',
    recipient_account VARCHAR(100) COMMENT '收款账号',
    expense_date DATE NOT NULL COMMENT '支出日期',
    invoice_no VARCHAR(100) COMMENT '发票号',
    has_invoice TINYINT DEFAULT 0 COMMENT '是否有发票',
    remark TEXT COMMENT '备注',
    operator_id BIGINT COMMENT '操作人ID',
    fiscal_year INT NOT NULL COMMENT '财年',
    fiscal_month INT NOT NULL COMMENT '财月',
    approval_status ENUM('PENDING', 'APPROVED', 'REJECTED') DEFAULT 'PENDING' COMMENT '审批状态',
    approver_id BIGINT COMMENT '审批人ID',
    approval_time DATETIME COMMENT '审批时间',
    approval_comment TEXT COMMENT '审批意见',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_date(expense_date),
    INDEX idx_type(expense_type),
    INDEX idx_fiscal(fiscal_year, fiscal_month),
    INDEX idx_status(approval_status)
);

-- 财务报表表
CREATE TABLE financial_statements (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT COMMENT '店铺ID，NULL表示总部财报',
    statement_type ENUM('DAILY', 'WEEKLY', 'MONTHLY', 'QUARTERLY', 'YEARLY') NOT NULL COMMENT '报表类型',
    period_start DATE NOT NULL COMMENT '开始日期',
    period_end DATE NOT NULL COMMENT '结束日期',
    total_income DECIMAL(12,2) NOT NULL COMMENT '总收入',
    total_expense DECIMAL(12,2) NOT NULL COMMENT '总支出',
    profit DECIMAL(12,2) NOT NULL COMMENT '利润',
    gross_profit_margin DECIMAL(5,2) COMMENT '毛利率(%)',
    net_profit_margin DECIMAL(5,2) COMMENT '净利率(%)',
    member_income DECIMAL(12,2) DEFAULT 0 COMMENT '会员收入',
    venue_income DECIMAL(12,2) DEFAULT 0 COMMENT '场地收入',
    product_income DECIMAL(12,2) DEFAULT 0 COMMENT '商品收入',
    coaching_income DECIMAL(12,2) DEFAULT 0 COMMENT '教练收入',
    activity_income DECIMAL(12,2) DEFAULT 0 COMMENT '活动收入',
    other_income DECIMAL(12,2) DEFAULT 0 COMMENT '其他收入',
    rent_expense DECIMAL(12,2) DEFAULT 0 COMMENT '租金支出',
    salary_expense DECIMAL(12,2) DEFAULT 0 COMMENT '工资支出',
    utilities_expense DECIMAL(12,2) DEFAULT 0 COMMENT '水电费支出',
    equipment_expense DECIMAL(12,2) DEFAULT 0 COMMENT '设备支出',
    inventory_expense DECIMAL(12,2) DEFAULT 0 COMMENT '库存支出',
    marketing_expense DECIMAL(12,2) DEFAULT 0 COMMENT '营销支出',
    maintenance_expense DECIMAL(12,2) DEFAULT 0 COMMENT '维护支出',
    tax_expense DECIMAL(12,2) DEFAULT 0 COMMENT '税费支出',
    other_expense DECIMAL(12,2) DEFAULT 0 COMMENT '其他支出',
    statement_status ENUM('DRAFT', 'CONFIRMED', 'AUDITED') DEFAULT 'DRAFT' COMMENT '报表状态',
    generator_id BIGINT COMMENT '生成人ID',
    auditor_id BIGINT COMMENT '审核人ID',
    audit_time DATETIME COMMENT '审核时间',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_period(period_start, period_end),
    INDEX idx_type(statement_type),
    UNIQUE KEY uk_store_period_type(store_id, period_start, period_end, statement_type)
);

-- 资产表
CREATE TABLE assets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT COMMENT '店铺ID，NULL表示总部资产',
    asset_name VARCHAR(100) NOT NULL COMMENT '资产名称',
    asset_type ENUM('REAL_ESTATE', 'EQUIPMENT', 'VEHICLE', 'FURNITURE', 'ELECTRONIC', 'INTANGIBLE', 'OTHER') NOT NULL COMMENT '资产类型',
    asset_code VARCHAR(50) COMMENT '资产编码',
    purchase_date DATE NOT NULL COMMENT '购入日期',
    purchase_price DECIMAL(12,2) NOT NULL COMMENT '购入价格',
    current_value DECIMAL(12,2) NOT NULL COMMENT '当前价值',
    depreciation_method ENUM('STRAIGHT_LINE', 'DECLINING_BALANCE', 'UNITS_OF_PRODUCTION', 'NONE') DEFAULT 'STRAIGHT_LINE' COMMENT '折旧方法',
    depreciation_period INT COMMENT '折旧年限(月)',
    monthly_depreciation DECIMAL(12,2) COMMENT '月折旧额',
    accumulated_depreciation DECIMAL(12,2) DEFAULT 0 COMMENT '累计折旧',
    asset_status ENUM('IN_USE', 'IDLE', 'MAINTENANCE', 'DISPOSED', 'SCRAPPED') DEFAULT 'IN_USE' COMMENT '资产状态',
    asset_location VARCHAR(100) COMMENT '资产位置',
    responsible_person VARCHAR(50) COMMENT '责任人',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_type(asset_type),
    INDEX idx_status(asset_status)
);

-- 财务预算表
CREATE TABLE financial_budgets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT COMMENT '店铺ID，NULL表示总部预算',
    budget_year INT NOT NULL COMMENT '预算年度',
    budget_month INT COMMENT '预算月份',
    budget_type ENUM('INCOME', 'EXPENSE') NOT NULL COMMENT '预算类型',
    budget_category VARCHAR(50) NOT NULL COMMENT '预算类别',
    budget_amount DECIMAL(12,2) NOT NULL COMMENT '预算金额',
    actual_amount DECIMAL(12,2) DEFAULT 0 COMMENT '实际金额',
    variance_amount DECIMAL(12,2) DEFAULT 0 COMMENT '差异金额',
    variance_percentage DECIMAL(5,2) DEFAULT 0 COMMENT '差异百分比',
    budget_status ENUM('PLANNING', 'APPROVED', 'IN_PROGRESS', 'COMPLETED') DEFAULT 'PLANNING' COMMENT '预算状态',
    creator_id BIGINT COMMENT '创建人ID',
    approver_id BIGINT COMMENT '审批人ID',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_period(budget_year, budget_month),
    INDEX idx_type(budget_type, budget_category),
    UNIQUE KEY uk_store_period_category(store_id, budget_year, budget_month, budget_type, budget_category)
);

-- 银行账户表
CREATE TABLE bank_accounts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT COMMENT '店铺ID，NULL表示总部账户',
    account_name VARCHAR(100) NOT NULL COMMENT '账户名称',
    account_number VARCHAR(50) NOT NULL COMMENT '账号',
    bank_name VARCHAR(100) NOT NULL COMMENT '开户行',
    account_type ENUM('CHECKING', 'SAVINGS', 'CREDIT', 'OTHER') NOT NULL COMMENT '账户类型',
    currency VARCHAR(10) DEFAULT 'CNY' COMMENT '币种',
    balance DECIMAL(12,2) DEFAULT 0 COMMENT '余额',
    is_primary TINYINT DEFAULT 0 COMMENT '是否主账户',
    account_status ENUM('ACTIVE', 'INACTIVE', 'FROZEN') DEFAULT 'ACTIVE' COMMENT '账户状态',
    open_date DATE COMMENT '开户日期',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_status(account_status),
    UNIQUE KEY uk_account_number(account_number)
);

-- 银行交易记录表
CREATE TABLE bank_transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL COMMENT '银行账户ID',
    transaction_type ENUM('DEPOSIT', 'WITHDRAWAL', 'TRANSFER', 'INTEREST', 'FEE', 'OTHER') NOT NULL COMMENT '交易类型',
    amount DECIMAL(12,2) NOT NULL COMMENT '交易金额',
    balance_before DECIMAL(12,2) NOT NULL COMMENT '交易前余额',
    balance_after DECIMAL(12,2) NOT NULL COMMENT '交易后余额',
    transaction_date DATETIME NOT NULL COMMENT '交易日期',
    transaction_no VARCHAR(50) COMMENT '交易流水号',
    counterparty VARCHAR(100) COMMENT '交易对手',
    counterparty_account VARCHAR(50) COMMENT '交易对手账号',
    reference_id BIGINT COMMENT '关联ID',
    reference_type VARCHAR(50) COMMENT '关联类型',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (account_id) REFERENCES bank_accounts(id),
    INDEX idx_account(account_id),
    INDEX idx_date(transaction_date),
    INDEX idx_type(transaction_type)
);

-- 应收账款表
CREATE TABLE accounts_receivable (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT COMMENT '店铺ID，NULL表示总部应收',
    customer_name VARCHAR(100) NOT NULL COMMENT '客户名称',
    customer_contact VARCHAR(50) COMMENT '客户联系方式',
    amount DECIMAL(12,2) NOT NULL COMMENT '应收金额',
    received_amount DECIMAL(12,2) DEFAULT 0 COMMENT '已收金额',
    due_date DATE NOT NULL COMMENT '截止日期',
    invoice_no VARCHAR(50) COMMENT '发票号',
    invoice_date DATE COMMENT '开票日期',
    status ENUM('PENDING', 'PARTIAL_PAID', 'PAID', 'OVERDUE', 'BAD_DEBT') DEFAULT 'PENDING' COMMENT '状态',
    related_id BIGINT COMMENT '关联ID',
    related_type VARCHAR(50) COMMENT '关联类型',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_date(due_date),
    INDEX idx_status(status)
);

-- 应付账款表
CREATE TABLE accounts_payable (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT COMMENT '店铺ID，NULL表示总部应付',
    vendor_name VARCHAR(100) NOT NULL COMMENT '供应商名称',
    vendor_contact VARCHAR(50) COMMENT '供应商联系方式',
    amount DECIMAL(12,2) NOT NULL COMMENT '应付金额',
    paid_amount DECIMAL(12,2) DEFAULT 0 COMMENT '已付金额',
    due_date DATE NOT NULL COMMENT '截止日期',
    invoice_no VARCHAR(50) COMMENT '发票号',
    invoice_date DATE COMMENT '收票日期',
    status ENUM('PENDING', 'PARTIAL_PAID', 'PAID', 'OVERDUE') DEFAULT 'PENDING' COMMENT '状态',
    related_id BIGINT COMMENT '关联ID',
    related_type VARCHAR(50) COMMENT '关联类型',
    expense_category VARCHAR(50) COMMENT '支出类别',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_date(due_date),
    INDEX idx_status(status)
);

-- 财务KPI指标表
CREATE TABLE financial_kpis (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT COMMENT '店铺ID，NULL表示总部KPI',
    period_year INT NOT NULL COMMENT 'KPI年度',
    period_month INT COMMENT 'KPI月份',
    kpi_name VARCHAR(50) NOT NULL COMMENT 'KPI名称',
    kpi_category ENUM('REVENUE', 'PROFIT', 'COST', 'EFFICIENCY', 'GROWTH', 'OTHER') NOT NULL COMMENT 'KPI类别',
    target_value DECIMAL(12,2) NOT NULL COMMENT '目标值',
    actual_value DECIMAL(12,2) DEFAULT 0 COMMENT '实际值',
    unit VARCHAR(20) COMMENT '单位',
    achievement_rate DECIMAL(5,2) DEFAULT 0 COMMENT '达成率(%)',
    weight DECIMAL(5,2) DEFAULT 1.00 COMMENT '权重',
    status ENUM('NOT_STARTED', 'IN_PROGRESS', 'ACHIEVED', 'FAILED', 'CANCELLED') DEFAULT 'NOT_STARTED' COMMENT '状态',
    responsible_person VARCHAR(50) COMMENT '责任人',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_period(period_year, period_month),
    INDEX idx_category(kpi_category),
    UNIQUE KEY uk_store_period_kpi(store_id, period_year, period_month, kpi_name)
);

-- 财务日结表
CREATE TABLE daily_settlements (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    settlement_date DATE NOT NULL COMMENT '结算日期',
    cash_income DECIMAL(12,2) DEFAULT 0 COMMENT '现金收入',
    online_income DECIMAL(12,2) DEFAULT 0 COMMENT '线上收入',
    pos_income DECIMAL(12,2) DEFAULT 0 COMMENT 'POS收入',
    other_income DECIMAL(12,2) DEFAULT 0 COMMENT '其他收入',
    total_income DECIMAL(12,2) DEFAULT 0 COMMENT '总收入',
    cash_expense DECIMAL(12,2) DEFAULT 0 COMMENT '现金支出',
    online_expense DECIMAL(12,2) DEFAULT 0 COMMENT '线上支出',
    other_expense DECIMAL(12,2) DEFAULT 0 COMMENT '其他支出',
    total_expense DECIMAL(12,2) DEFAULT 0 COMMENT '总支出',
    net_cash_flow DECIMAL(12,2) DEFAULT 0 COMMENT '净现金流',
    beginning_cash DECIMAL(12,2) DEFAULT 0 COMMENT '期初现金',
    ending_cash DECIMAL(12,2) DEFAULT 0 COMMENT '期末现金',
    cash_difference DECIMAL(12,2) DEFAULT 0 COMMENT '现金差额',
    difference_reason TEXT COMMENT '差额原因',
    settler_id BIGINT NOT NULL COMMENT '结算人ID',
    settle_time DATETIME NOT NULL COMMENT '结算时间',
    verifier_id BIGINT COMMENT '核验人ID',
    verify_time DATETIME COMMENT '核验时间',
    status ENUM('DRAFT', 'SETTLED', 'VERIFIED', 'ABNORMAL') DEFAULT 'DRAFT' COMMENT '状态',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_store(store_id),
    INDEX idx_date(settlement_date),
    INDEX idx_status(status),
    UNIQUE KEY uk_store_date(store_id, settlement_date)
);
```

### 3.7 活动管理模块
#### 3.7.1 活动策划
- 活动类型管理
  - 射箭比赛活动
  - 射箭培训活动
  - 会员专享活动
  - 节日主题活动
  - 团建活动
- 活动内容设计
  - 活动主题设置
  - 活动规则制定
  - 活动流程规划
  - 活动奖品设置
- 活动资源规划
  - 场地资源分配
  - 教练资源安排
  - 设备资源准备
  - 活动预算管理

#### 3.7.2 活动发布
- 活动信息管理
  - 活动基本信息设置
  - 活动图片上传
  - 活动详情编辑
  - 活动FAQ管理
- 活动发布渠道
  - 小程序发布
  - 公众号推文
  - 店内海报展示
  - 短信通知会员
- 活动状态管理
  - 草稿状态
  - 待发布状态
  - 进行中状态
  - 已结束状态
  - 已取消状态

#### 3.7.3 活动报名
- 报名渠道管理
  - 线上自助报名
  - 前台人工报名
  - 电话预约报名
  - 第三方平台报名
- 报名流程管理
  - 报名信息采集
  - 报名费用支付
  - 报名确认通知
  - 报名取消处理
- 报名名额管理
  - 活动容量设置
  - 候补名单管理
  - 名额动态调整
  - 分组名额控制

#### 3.7.4 活动统计
- 报名数据统计
  - 报名人数统计
  - 报名转化率分析
  - 报名来源分析
  - 报名时段分布
- 活动执行统计
  - 活动出席率统计
  - 活动满意度调查
  - 活动异常记录
  - 活动完成情况
- 活动效果分析
  - 活动收入统计
  - 新增会员转化率
  - 会员复购率分析
  - 活动ROI计算

#### 3.7.5 优惠券管理
- 优惠券类型设置
  - 满减券设置
  - 折扣券设置
  - 代金券设置
  - 体验券设置
- 优惠券发放管理
  - 批量发放功能
  - 定向发放功能
  - 自动触发发放
  - 活动关联发放
- 优惠券使用规则
  - 使用门槛设置
  - 有效期设置
  - 叠加规则设置
  - 适用范围限制
- 优惠券数据分析
  - 发放数量统计
  - 使用率分析
  - 带动消费分析
  - 失效原因分析

#### 3.7.6 储值活动管理
- 满赠活动
  - 满额赠送金额设置
    - 阶梯式满赠设置
    - 首次充值特殊满赠
    - 会员等级差异化满赠
    - 节日特殊满赠活动
  - 活动有效期设置
    - 固定日期区间设置
    - 灵活天数设置
    - 长期有效设置
    - 特殊节假日设置
  - 活动状态管理
    - 未开始状态管理
    - 进行中状态管理
    - 已暂停状态管理
    - 已结束状态管理
- 折扣活动  
  - 储值折扣比例设置
    - 阶梯式折扣设置
    - 会员等级差异化折扣
    - 特定时段特惠折扣
    - 指定项目折扣设置
  - 活动有效期设置
    - 固定时间段设置
    - 周期性时间设置
    - 限时折扣设置
    - 长期折扣设置
  - 活动状态管理
    - 活动预告管理
    - 活动上线管理
    - 活动变更管理
    - 活动终止管理

#### 3.7.7 数据库设计

```sql
-- 活动基本信息表
CREATE TABLE activities (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL COMMENT '活动标题',
    code VARCHAR(20) NOT NULL COMMENT '活动编码',
    store_id BIGINT COMMENT '所属店铺ID，NULL表示总部活动',
    category ENUM('COMPETITION', 'TRAINING', 'PROMOTION', 'SOCIAL', 'OTHER') NOT NULL COMMENT '活动类别',
    description TEXT COMMENT '活动描述',
    poster_url VARCHAR(255) COMMENT '活动海报URL',
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME NOT NULL COMMENT '结束时间',
    registration_start_time DATETIME COMMENT '报名开始时间',
    registration_end_time DATETIME COMMENT '报名截止时间',
    venue_id BIGINT COMMENT '活动场地ID',
    location VARCHAR(200) COMMENT '活动地点',
    capacity INT COMMENT '活动容量',
    fee DECIMAL(10,2) COMMENT '活动费用',
    member_fee DECIMAL(10,2) COMMENT '会员价格',
    fee_description VARCHAR(200) COMMENT '费用说明',
    organizer VARCHAR(50) COMMENT '主办方',
    contact_person VARCHAR(50) COMMENT '联系人',
    contact_phone VARCHAR(20) COMMENT '联系电话',
    status ENUM('DRAFT', 'PUBLISHED', 'ONGOING', 'ENDED', 'CANCELLED') DEFAULT 'DRAFT' COMMENT '活动状态',
    is_visible TINYINT DEFAULT 1 COMMENT '是否可见',
    is_deleted TINYINT DEFAULT 0 COMMENT '是否删除',
    creator_id BIGINT COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (store_id) REFERENCES stores(id),
    INDEX idx_code(code),
    INDEX idx_store(store_id),
    INDEX idx_status(status),
    INDEX idx_time(start_time, end_time)
);

-- 活动详情表（存储富文本内容）
CREATE TABLE activity_details (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    activity_id BIGINT NOT NULL COMMENT '活动ID',
    content LONGTEXT COMMENT '活动详情(富文本)',
    notice TEXT COMMENT '活动须知',
    agenda TEXT COMMENT '活动日程',
    awards TEXT COMMENT '奖项设置(比赛类活动)',
    rules TEXT COMMENT '活动规则',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    INDEX idx_activity(activity_id)
);

-- 活动图片表
CREATE TABLE activity_images (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    activity_id BIGINT NOT NULL COMMENT '活动ID',
    image_url VARCHAR(255) NOT NULL COMMENT '图片URL',
    image_type ENUM('POSTER', 'BANNER', 'DETAIL', 'RESULT', 'OTHER') COMMENT '图片类型',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    INDEX idx_activity(activity_id)
);

-- 活动报名表
CREATE TABLE activity_registrations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    activity_id BIGINT NOT NULL COMMENT '活动ID',
    member_id BIGINT COMMENT '会员ID',
    name VARCHAR(50) NOT NULL COMMENT '参与者姓名',
    phone VARCHAR(20) NOT NULL COMMENT '参与者电话',
    gender TINYINT COMMENT '性别:0未知,1男,2女',
    age INT COMMENT '年龄',
    email VARCHAR(100) COMMENT '电子邮箱',
    id_number VARCHAR(20) COMMENT '身份证号',
    team_name VARCHAR(50) COMMENT '团队名称',
    registration_time DATETIME NOT NULL COMMENT '报名时间',
    status ENUM('PENDING', 'CONFIRMED', 'CANCELLED', 'ATTENDED', 'ABSENT') DEFAULT 'PENDING' COMMENT '报名状态',
    fee_amount DECIMAL(10,2) COMMENT '报名费用',
    payment_status ENUM('UNPAID', 'PAID', 'REFUNDED') DEFAULT 'UNPAID' COMMENT '支付状态',
    payment_time DATETIME COMMENT '支付时间',
    payment_method VARCHAR(20) COMMENT '支付方式',
    transaction_id VARCHAR(50) COMMENT '交易ID',
    refund_status ENUM('NONE', 'REQUESTED', 'PROCESSED', 'REJECTED') DEFAULT 'NONE' COMMENT '退款状态',
    refund_time DATETIME COMMENT '退款时间',
    refund_amount DECIMAL(10,2) COMMENT '退款金额',
    extra_info JSON COMMENT '额外信息(JSON格式)',
    check_in_time DATETIME COMMENT '签到时间',
    check_in_operator_id BIGINT COMMENT '签到操作员ID',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    INDEX idx_activity(activity_id),
    INDEX idx_member(member_id),
    INDEX idx_status(status),
    INDEX idx_payment(payment_status)
);

-- 活动报名自定义字段表
CREATE TABLE activity_custom_fields (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    activity_id BIGINT NOT NULL COMMENT '活动ID',
    field_name VARCHAR(50) NOT NULL COMMENT '字段名称',
    field_label VARCHAR(50) NOT NULL COMMENT '字段标签',
    field_type ENUM('TEXT', 'NUMBER', 'DATE', 'SELECT', 'RADIO', 'CHECKBOX', 'TEXTAREA') NOT NULL COMMENT '字段类型',
    field_options TEXT COMMENT '字段选项(JSON格式,用于SELECT/RADIO/CHECKBOX)',
    is_required TINYINT DEFAULT 0 COMMENT '是否必填',
    placeholder VARCHAR(100) COMMENT '占位文本',
    default_value VARCHAR(100) COMMENT '默认值',
    validation_rule VARCHAR(200) COMMENT '验证规则',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    INDEX idx_activity(activity_id)
);

-- 活动评论表
CREATE TABLE activity_comments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    activity_id BIGINT NOT NULL COMMENT '活动ID',
    member_id BIGINT COMMENT '会员ID',
    name VARCHAR(50) NOT NULL COMMENT '评论者姓名',
    content TEXT NOT NULL COMMENT '评论内容',
    rating INT COMMENT '评分(1-5星)',
    image_urls TEXT COMMENT '图片URL列表(JSON格式)',
    status ENUM('PENDING', 'APPROVED', 'REJECTED') DEFAULT 'PENDING' COMMENT '状态',
    parent_id BIGINT DEFAULT 0 COMMENT '父评论ID，0表示顶级评论',
    is_hidden TINYINT DEFAULT 0 COMMENT '是否隐藏',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    INDEX idx_activity(activity_id),
    INDEX idx_member(member_id),
    INDEX idx_parent(parent_id)
);

-- 活动结果记录表（比赛类活动）
CREATE TABLE activity_results (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    activity_id BIGINT NOT NULL COMMENT '活动ID',
    registration_id BIGINT NOT NULL COMMENT '报名ID',
    ranking INT COMMENT '名次',
    score DECIMAL(10,2) COMMENT '得分',
    category VARCHAR(50) COMMENT '比赛类别',
    achievement VARCHAR(100) COMMENT '成绩',
    award VARCHAR(100) COMMENT '获得奖项',
    result_detail TEXT COMMENT '成绩细节',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    FOREIGN KEY (registration_id) REFERENCES activity_registrations(id),
    INDEX idx_activity(activity_id),
    INDEX idx_registration(registration_id)
);

-- 活动统计表
CREATE TABLE activity_statistics (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    activity_id BIGINT NOT NULL COMMENT '活动ID',
    view_count INT DEFAULT 0 COMMENT '浏览次数',
    registration_count INT DEFAULT 0 COMMENT '报名人数',
    confirmed_count INT DEFAULT 0 COMMENT '确认人数',
    attended_count INT DEFAULT 0 COMMENT '实际参与人数',
    revenue DECIMAL(12,2) DEFAULT 0 COMMENT '活动收入',
    cost DECIMAL(12,2) DEFAULT 0 COMMENT '活动成本',
    profit DECIMAL(12,2) DEFAULT 0 COMMENT '活动利润',
    avg_rating DECIMAL(3,2) DEFAULT 0 COMMENT '平均评分',
    comment_count INT DEFAULT 0 COMMENT '评论数量',
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后更新时间',
    FOREIGN KEY (activity_id) REFERENCES activities(id),
    UNIQUE KEY uk_activity(activity_id)
);

-- 优惠券表
CREATE TABLE coupons (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    code VARCHAR(20) NOT NULL COMMENT '优惠券代码',
    name VARCHAR(100) NOT NULL COMMENT '优惠券名称',
    type ENUM('FIXED', 'PERCENTAGE', 'FREE') NOT NULL COMMENT '优惠类型:固定金额,百分比,免费',
    value DECIMAL(10,2) NOT NULL COMMENT '优惠值(固定金额或百分比)',
    min_order_amount DECIMAL(10,2) DEFAULT 0 COMMENT '最低订单金额',
    max_discount_amount DECIMAL(10,2) COMMENT '最大优惠金额(百分比类型)',
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME NOT NULL COMMENT '结束时间',
    store_id BIGINT COMMENT '店铺ID，NULL表示全店通用',
    usage_scope ENUM('ALL', 'PRODUCT', 'CATEGORY', 'VENUE', 'ACTIVITY') DEFAULT 'ALL' COMMENT '使用范围',
    scope_values TEXT COMMENT '范围值(JSON格式)',
    total_quantity INT COMMENT '总数量',
    remaining_quantity INT COMMENT '剩余数量',
    per_user_limit INT DEFAULT 1 COMMENT '每人限领数量',
    user_level_limit VARCHAR(50) COMMENT '会员等级限制',
    is_stackable TINYINT DEFAULT 0 COMMENT '是否可叠加使用',
    description TEXT COMMENT '优惠券描述',
    status ENUM('ACTIVE', 'DISABLED', 'EXPIRED') DEFAULT 'ACTIVE' COMMENT '状态',
    creator_id BIGINT COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code(code),
    INDEX idx_store(store_id),
    INDEX idx_time(start_time, end_time),
    INDEX idx_status(status)
);

-- 用户优惠券表
CREATE TABLE user_coupons (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    coupon_id BIGINT NOT NULL COMMENT '优惠券ID',
    member_id BIGINT NOT NULL COMMENT '会员ID',
    coupon_code VARCHAR(50) NOT NULL COMMENT '优惠券兑换码',
    obtain_time DATETIME NOT NULL COMMENT '获取时间',
    status ENUM('UNUSED', 'USED', 'EXPIRED') DEFAULT 'UNUSED' COMMENT '状态',
    use_time DATETIME COMMENT '使用时间',
    use_order_id BIGINT COMMENT '使用订单ID',
    expire_time DATETIME NOT NULL COMMENT '过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (coupon_id) REFERENCES coupons(id),
    INDEX idx_coupon(coupon_id),
    INDEX idx_member(member_id),
    INDEX idx_code(coupon_code),
    INDEX idx_status(status)
);

-- 储值活动表
CREATE TABLE recharge_promotions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '活动名称',
    code VARCHAR(20) NOT NULL COMMENT '活动编码',
    type ENUM('AMOUNT_GIFT', 'DISCOUNT') NOT NULL COMMENT '活动类型:满额赠送,折扣优惠',
    store_id BIGINT COMMENT '店铺ID，NULL表示全店通用',
    start_time DATETIME NOT NULL COMMENT '开始时间',
    end_time DATETIME NOT NULL COMMENT '结束时间',
    description TEXT COMMENT '活动描述',
    member_level_limit VARCHAR(50) COMMENT '会员等级限制',
    status ENUM('ACTIVE', 'INACTIVE', 'EXPIRED') DEFAULT 'ACTIVE' COMMENT '状态',
    creator_id BIGINT COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code(code),
    INDEX idx_store(store_id),
    INDEX idx_time(start_time, end_time),
    INDEX idx_status(status)
);

-- 满额赠送规则表
CREATE TABLE amount_gift_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    promotion_id BIGINT NOT NULL COMMENT '储值活动ID',
    threshold_amount DECIMAL(10,2) NOT NULL COMMENT '满额阈值',
    gift_amount DECIMAL(10,2) NOT NULL COMMENT '赠送金额',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (promotion_id) REFERENCES recharge_promotions(id),
    INDEX idx_promotion(promotion_id)
);

-- 折扣优惠规则表
CREATE TABLE discount_rules (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    promotion_id BIGINT NOT NULL COMMENT '储值活动ID',
    threshold_amount DECIMAL(10,2) NOT NULL COMMENT '满额阈值',
    discount_rate DECIMAL(5,2) NOT NULL COMMENT '折扣比例(例如:0.85表示85折)',
    max_discount_amount DECIMAL(10,2) COMMENT '最大优惠金额',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (promotion_id) REFERENCES recharge_promotions(id),
    INDEX idx_promotion(promotion_id)
);

-- 活动参与记录表
CREATE TABLE promotion_participation_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    promotion_id BIGINT NOT NULL COMMENT '活动ID',
    promotion_type VARCHAR(50) NOT NULL COMMENT '活动类型',
    member_id BIGINT NOT NULL COMMENT '会员ID',
    transaction_id BIGINT NOT NULL COMMENT '交易ID',
    amount DECIMAL(10,2) NOT NULL COMMENT '交易金额',
    benefit_amount DECIMAL(10,2) NOT NULL COMMENT '优惠/赠送金额',
    participation_time DATETIME NOT NULL COMMENT '参与时间',
    store_id BIGINT NOT NULL COMMENT '店铺ID',
    operator_id BIGINT COMMENT '操作员ID',
    remark TEXT COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_promotion(promotion_id),
    INDEX idx_member(member_id),
    INDEX idx_transaction(transaction_id),
    INDEX idx_store(store_id),
    INDEX idx_time(participation_time)
);
```

### 3.8 数据回收站模块
#### 3.8.1 已删除会员数据管理
- 查询筛选功能
  - 按会员姓名筛选
  - 按会员手机号筛选
  - 按会员卡号筛选
  - 按删除时间范围筛选
  - 按删除操作人筛选
  - 按会员等级筛选
- 数据恢复功能
  - 单个会员数据恢复
  - 恢复确认机制
  - 恢复操作日志记录
  - 恢复后状态设置
  - 恢复冲突处理
- 彻底删除功能
  - 单个会员彻底删除
  - 删除确认机制
  - 删除操作日志记录
  - 删除权限控制

#### 3.8.2 已删除商品数据管理
- 查询筛选功能
  - 按商品名称筛选
  - 按商品编码筛选
  - 按商品分类筛选
  - 按删除时间范围筛选
  - 按删除操作人筛选
  - 按商品状态筛选
- 数据恢复功能
  - 单个商品数据恢复
  - 恢复确认机制
  - 恢复操作日志记录
  - 恢复后商品状态设置
  - 关联数据一并恢复
- 彻底删除功能
  - 单个商品彻底删除
  - 删除确认机制
  - 删除操作日志记录
  - 关联数据处理策略

#### 3.8.3 已删除活动数据管理
- 查询筛选功能
  - 按活动名称筛选
  - 按活动编码筛选
  - 按活动类型筛选
  - 按活动时间范围筛选
  - 按删除时间范围筛选
  - 按删除操作人筛选
- 数据恢复功能
  - 单个活动数据恢复
  - 恢复确认机制
  - 恢复操作日志记录
  - 恢复后活动状态设置
  - 关联报名数据处理
- 彻底删除功能
  - 单个活动彻底删除
  - 删除确认机制
  - 删除操作日志记录
  - 关联数据处理策略

#### 3.8.4 批量操作功能
- 批量恢复功能
  - 多选恢复机制
  - 批量恢复进度展示
  - 批量恢复结果反馈
  - 批量恢复操作日志
- 批量彻底删除功能
  - 多选删除机制
  - 批量删除安全确认
  - 批量删除权限控制
  - 批量删除操作日志
- 数据导出功能
  - 已删除数据导出
  - 导出格式选择
  - 导出内容选择
  - 导出权限控制

#### 3.8.5 数据库设计
```sql
-- 回收站记录主表
CREATE TABLE recycle_bin (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    record_type ENUM('MEMBER', 'PRODUCT', 'ACTIVITY', 'COUPON', 'OTHER') NOT NULL COMMENT '记录类型',
    record_id BIGINT NOT NULL COMMENT '原记录ID',
    record_code VARCHAR(50) COMMENT '记录编码(会员卡号/商品编码/活动编码等)',
    record_name VARCHAR(100) NOT NULL COMMENT '记录名称(会员姓名/商品名称/活动名称等)',
    store_id BIGINT COMMENT '所属店铺ID',
    delete_time DATETIME NOT NULL COMMENT '删除时间',
    delete_operator_id BIGINT COMMENT '删除操作人ID',
    delete_operator_name VARCHAR(50) COMMENT '删除操作人姓名',
    delete_reason TEXT COMMENT '删除原因',
    retention_days INT DEFAULT 90 COMMENT '保留天数',
    expiry_time DATETIME COMMENT '过期时间(超过将自动彻底删除)',
    is_recoverable TINYINT DEFAULT 1 COMMENT '是否可恢复',
    recovery_time DATETIME COMMENT '恢复时间',
    recovery_operator_id BIGINT COMMENT '恢复操作人ID',
    recovery_operator_name VARCHAR(50) COMMENT '恢复操作人姓名',
    status ENUM('DELETED', 'RECOVERED', 'PURGED') DEFAULT 'DELETED' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_type_id(record_type, record_id),
    INDEX idx_code(record_code),
    INDEX idx_name(record_name),
    INDEX idx_store(store_id),
    INDEX idx_delete_time(delete_time),
    INDEX idx_status(status),
    INDEX idx_expiry(expiry_time)
);

-- 已删除会员数据表
CREATE TABLE deleted_members (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    recycle_id BIGINT NOT NULL COMMENT '回收站记录ID',
    original_id BIGINT NOT NULL COMMENT '原会员ID',
    member_data JSON NOT NULL COMMENT '会员数据完整备份(JSON格式)',
    transaction_data JSON COMMENT '相关交易数据备份(JSON格式)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recycle_id) REFERENCES recycle_bin(id),
    INDEX idx_recycle(recycle_id),
    INDEX idx_original(original_id)
);

-- 已删除商品数据表
CREATE TABLE deleted_products (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    recycle_id BIGINT NOT NULL COMMENT '回收站记录ID',
    original_id BIGINT NOT NULL COMMENT '原商品ID',
    product_data JSON NOT NULL COMMENT '商品数据完整备份(JSON格式)',
    inventory_data JSON COMMENT '相关库存数据备份(JSON格式)',
    image_data JSON COMMENT '商品图片数据备份(JSON格式)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recycle_id) REFERENCES recycle_bin(id),
    INDEX idx_recycle(recycle_id),
    INDEX idx_original(original_id)
);

-- 已删除活动数据表
CREATE TABLE deleted_activities (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    recycle_id BIGINT NOT NULL COMMENT '回收站记录ID',
    original_id BIGINT NOT NULL COMMENT '原活动ID',
    activity_data JSON NOT NULL COMMENT '活动数据完整备份(JSON格式)',
    registration_data JSON COMMENT '活动报名数据备份(JSON格式)',
    image_data JSON COMMENT '活动图片数据备份(JSON格式)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recycle_id) REFERENCES recycle_bin(id),
    INDEX idx_recycle(recycle_id),
    INDEX idx_original(original_id)
);

-- 数据恢复日志表
CREATE TABLE recovery_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    recycle_id BIGINT NOT NULL COMMENT '回收站记录ID',
    record_type ENUM('MEMBER', 'PRODUCT', 'ACTIVITY', 'COUPON', 'OTHER') NOT NULL COMMENT '记录类型',
    original_id BIGINT NOT NULL COMMENT '原记录ID',
    new_id BIGINT COMMENT '恢复后的新记录ID',
    recovery_time DATETIME NOT NULL COMMENT '恢复时间',
    operator_id BIGINT COMMENT '操作人ID',
    operator_name VARCHAR(50) COMMENT '操作人姓名',
    recovery_result ENUM('SUCCESS', 'CONFLICT', 'FAILED') NOT NULL COMMENT '恢复结果',
    conflict_detail TEXT COMMENT '冲突详情',
    resolution TEXT COMMENT '解决方案',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recycle_id) REFERENCES recycle_bin(id),
    INDEX idx_recycle(recycle_id),
    INDEX idx_type(record_type),
    INDEX idx_original(original_id),
    INDEX idx_time(recovery_time),
    INDEX idx_result(recovery_result)
);

-- 彻底删除日志表
CREATE TABLE purge_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    recycle_id COMMENT '回收站记录ID',
    record_type ENUM('MEMBER', 'PRODUCT', 'ACTIVITY', 'COUPON', 'OTHER') NOT NULL COMMENT '记录类型',
    original_id BIGINT NOT NULL COMMENT '原记录ID',
    record_summary JSON COMMENT '记录摘要(保留少量关键信息)',
    purge_time DATETIME NOT NULL COMMENT '彻底删除时间',
    purge_type ENUM('MANUAL', 'AUTOMATIC', 'BATCH') NOT NULL COMMENT '删除类型:手动,自动过期,批量',
    operator_id BIGINT COMMENT '操作人ID',
    operator_name VARCHAR(50) COMMENT '操作人姓名',
    authorization_id BIGINT COMMENT '授权人ID(适用于需要高级权限的删除)',
    authorization_time DATETIME COMMENT '授权时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_type(record_type),
    INDEX idx_original(original_id),
    INDEX idx_time(purge_time),
    INDEX idx_purge_type(purge_type)
);

-- 批量操作记录表
CREATE TABLE batch_operation_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    operation_type ENUM('BATCH_RECOVER', 'BATCH_PURGE') NOT NULL COMMENT '操作类型',
    record_type ENUM('MEMBER', 'PRODUCT', 'ACTIVITY', 'MIXED') NOT NULL COMMENT '记录类型',
    record_count INT NOT NULL COMMENT '记录数量',
    success_count INT DEFAULT 0 COMMENT '成功数量',
    failed_count INT DEFAULT 0 COMMENT '失败数量',
    operation_time DATETIME NOT NULL COMMENT '操作时间',
    condition_criteria TEXT COMMENT '筛选条件(JSON格式)',
    target_ids TEXT COMMENT '目标ID列表(JSON格式)',
    operator_id BIGINT COMMENT '操作人ID',
    operator_name VARCHAR(50) COMMENT '操作人姓名',
    operation_status ENUM('PROCESSING', 'COMPLETED', 'FAILED', 'PARTIALLY_COMPLETED') DEFAULT 'PROCESSING' COMMENT '操作状态',
    error_message TEXT COMMENT '错误信息',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_type(operation_type, record_type),
    INDEX idx_time(operation_time),
    INDEX idx_status(operation_status)
);
```

### 3.9 系统设置模块
#### 3.9.1 基本信息设置
- 公司信息管理
  - 公司名称设置
  - 公司logo上传
  - 联系方式配置
  - 营业执照信息
  - 公司地址管理
- 系统主题设置
  - 系统颜色方案
  - 系统字体设置
  - 布局样式选择
  - 自定义主题配置
- 系统公告管理
  - 公告内容编辑
  - 公告发布设置
  - 公告显示位置
  - 公告有效期设置

#### 3.9.2 系统参数配置
- 业务参数设置
  - 会员积分规则配置
  - 储值优惠规则配置
  - 预约时间间隔设置
  - 取消预约规则设置
  - 自动结算时间设置
- 安全参数设置
  - 密码策略配置
  - 登录安全策略
  - 数据访问控制
  - 敏感操作配置
  - 会话超时设置
- 通知参数设置
  - 短信通知配置
  - 微信消息配置
  - 邮件通知配置
  - 站内消息配置
  - 消息模板管理

#### 3.9.3 数据备份与恢复
- 自动备份管理
  - 备份周期设置
  - 备份内容选择
  - 备份存储位置
  - 备份保留策略
  - 备份执行时间
- 手动备份功能
  - 全量数据备份
  - 增量数据备份
  - 指定表备份
  - 备份文件加密
  - 备份文件压缩
- 数据恢复管理
  - 备份文件列表
  - 恢复前验证
  - 选择性数据恢复
  - 恢复操作日志
  - 恢复后数据校验

#### 3.9.4 日志管理
- 系统日志查询
  - 操作日志查询
  - 登录日志查询
  - 错误日志查询
  - 安全审计日志
  - 性能监控日志
- 日志配置管理
  - 日志级别设置
  - 日志保留期限
  - 日志存储策略
  - 日志轮转配置
  - 日志导出格式
- 日志分析工具
  - 日志统计分析
  - 异常行为检测
  - 操作频率分析
  - 用户行为分析
  - 系统健康报告

#### 3.9.5 数据库设计
```sql
-- 公司信息表
CREATE TABLE company_info (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    company_name VARCHAR(100) NOT NULL COMMENT '公司名称',
    logo_url VARCHAR(255) COMMENT '公司LOGO URL',
    favicon_url VARCHAR(255) COMMENT '网站图标URL',
    contact_phone VARCHAR(20) COMMENT '联系电话',
    contact_email VARCHAR(100) COMMENT '联系邮箱',
    business_license VARCHAR(50) COMMENT '营业执照编号',
    license_image_url VARCHAR(255) COMMENT '营业执照图片URL',
    legal_representative VARCHAR(50) COMMENT '法定代表人',
    tax_id VARCHAR(50) COMMENT '税务登记号',
    address VARCHAR(200) COMMENT '公司地址',
    postal_code VARCHAR(20) COMMENT '邮政编码',
    website VARCHAR(100) COMMENT '公司网站',
    brief_introduction TEXT COMMENT '公司简介',
    setup_time DATETIME COMMENT '成立时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT COMMENT '更新人ID'
);

-- 系统主题表
CREATE TABLE system_themes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    theme_name VARCHAR(50) NOT NULL COMMENT '主题名称',
    is_default TINYINT DEFAULT 0 COMMENT '是否默认主题',
    is_active TINYINT DEFAULT 0 COMMENT '是否启用',
    primary_color VARCHAR(20) COMMENT '主色调',
    secondary_color VARCHAR(20) COMMENT '次色调',
    font_family VARCHAR(100) COMMENT '字体系列',
    font_size VARCHAR(20) COMMENT '默认字体大小',
    layout_type VARCHAR(20) COMMENT '布局类型',
    menu_style VARCHAR(20) COMMENT '菜单样式',
    header_style VARCHAR(20) COMMENT '页眉样式',
    theme_config JSON COMMENT '主题配置(JSON格式)',
    preview_image_url VARCHAR(255) COMMENT '预览图URL',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    INDEX idx_default(is_default),
    INDEX idx_active(is_active)
);

-- 用户主题偏好表
CREATE TABLE user_theme_preferences (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    theme_id BIGINT NOT NULL COMMENT '主题ID',
    is_active TINYINT DEFAULT 1 COMMENT '是否启用',
    custom_config JSON COMMENT '自定义配置(JSON格式)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (theme_id) REFERENCES system_themes(id),
    UNIQUE KEY uk_user(user_id)
);

-- 系统公告表
CREATE TABLE system_announcements (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL COMMENT '公告标题',
    content TEXT NOT NULL COMMENT '公告内容',
    announcement_type ENUM('SYSTEM', 'MAINTENANCE', 'PROMOTION', 'UPDATE', 'OTHER') NOT NULL COMMENT '公告类型',
    priority ENUM('HIGH', 'MEDIUM', 'LOW') DEFAULT 'MEDIUM' COMMENT '优先级',
    start_time DATETIME NOT NULL COMMENT '开始显示时间',
    end_time DATETIME COMMENT '结束显示时间',
    display_position VARCHAR(50) COMMENT '显示位置',
    is_popup TINYINT DEFAULT 0 COMMENT '是否弹窗显示',
    is_sticky TINYINT DEFAULT 0 COMMENT '是否置顶',
    status ENUM('DRAFT', 'PUBLISHED', 'EXPIRED', 'CANCELLED') DEFAULT 'DRAFT' COMMENT '状态',
    target_user_type VARCHAR(50) DEFAULT 'ALL' COMMENT '目标用户类型',
    target_store_ids TEXT COMMENT '目标店铺ID(JSON数组)',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    INDEX idx_time(start_time, end_time),
    INDEX idx_status(status),
    INDEX idx_type(announcement_type)
);

-- 公告阅读记录表
CREATE TABLE announcement_read_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    announcement_id BIGINT NOT NULL COMMENT '公告ID',
    user_id BIGINT NOT NULL COMMENT '用户ID',
    read_time DATETIME NOT NULL COMMENT '阅读时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (announcement_id) REFERENCES system_announcements(id),
    UNIQUE KEY uk_user_announcement(user_id, announcement_id)
);

-- 系统参数表
CREATE TABLE system_parameters (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    param_group VARCHAR(50) NOT NULL COMMENT '参数组',
    param_name VARCHAR(50) NOT NULL COMMENT '参数名称',
    param_key VARCHAR(50) NOT NULL COMMENT '参数键',
    param_value TEXT NOT NULL COMMENT '参数值',
    value_type VARCHAR(20) NOT NULL COMMENT '值类型(STRING/INT/FLOAT/BOOLEAN/JSON)',
    is_system TINYINT DEFAULT 0 COMMENT '是否系统参数',
    is_editable TINYINT DEFAULT 1 COMMENT '是否可编辑',
    is_sensitive TINYINT DEFAULT 0 COMMENT '是否敏感参数',
    description VARCHAR(200) COMMENT '参数描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT COMMENT '更新人ID',
    UNIQUE KEY uk_key(param_key),
    INDEX idx_group(param_group)
);

-- 数据字典表
CREATE TABLE system_dictionaries (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    dict_type VARCHAR(50) NOT NULL COMMENT '字典类型',
    dict_code VARCHAR(50) NOT NULL COMMENT '字典编码',
    dict_name VARCHAR(100) NOT NULL COMMENT '字典名称',
    dict_value VARCHAR(100) NOT NULL COMMENT '字典值',
    sort_order INT DEFAULT 0 COMMENT '排序顺序',
    is_default TINYINT DEFAULT 0 COMMENT '是否默认',
    status ENUM('ACTIVE', 'INACTIVE') DEFAULT 'ACTIVE' COMMENT '状态',
    remark VARCHAR(200) COMMENT '备注',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    UNIQUE KEY uk_type_code(dict_type, dict_code),
    INDEX idx_type(dict_type),
    INDEX idx_status(status)
);

-- 数据备份记录表
CREATE TABLE backup_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    backup_name VARCHAR(100) NOT NULL COMMENT '备份名称',
    backup_type ENUM('FULL', 'INCREMENTAL', 'PARTIAL') NOT NULL COMMENT '备份类型',
    backup_scope VARCHAR(100) COMMENT '备份范围',
    backup_time DATETIME NOT NULL COMMENT '备份时间',
    backup_size BIGINT COMMENT '备份大小(字节)',
    storage_path VARCHAR(255) NOT NULL COMMENT '存储路径',
    storage_type ENUM('LOCAL', 'CLOUD', 'FTP', 'NAS') DEFAULT 'LOCAL' COMMENT '存储类型',
    status ENUM('PROCESSING', 'COMPLETED', 'FAILED') NOT NULL COMMENT '状态',
    result_message TEXT COMMENT '结果信息',
    retention_days INT DEFAULT 30 COMMENT '保留天数',
    is_auto_backup TINYINT DEFAULT 0 COMMENT '是否自动备份',
    is_locked TINYINT DEFAULT 0 COMMENT '是否锁定(防止自动删除)',
    created_by BIGINT COMMENT '创建人ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_time(backup_time),
    INDEX idx_type(backup_type),
    INDEX idx_status(status)
);

-- 数据恢复记录表
CREATE TABLE restore_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    backup_id BIGINT COMMENT '备份ID',
    restore_name VARCHAR(100) NOT NULL COMMENT '恢复操作名称',
    restore_time DATETIME NOT NULL COMMENT '恢复时间',
    restore_scope VARCHAR(100) COMMENT '恢复范围',
    backup_info TEXT COMMENT '备份信息',
    status ENUM('PREPARING', 'PROCESSING', 'COMPLETED', 'FAILED') NOT NULL COMMENT '状态',
    result_message TEXT COMMENT '结果信息',
    operation_ip VARCHAR(50) COMMENT '操作IP',
    operator_id BIGINT NOT NULL COMMENT '操作人ID',
    operator_name VARCHAR(50) NOT NULL COMMENT '操作人姓名',
    approver_id BIGINT COMMENT '审批人ID',
    approval_time DATETIME COMMENT '审批时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_backup(backup_id),
    INDEX idx_time(restore_time),
    INDEX idx_status(status)
);

-- 系统日志清理策略表
CREATE TABLE log_cleanup_policies {
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    log_type VARCHAR(50) NOT NULL COMMENT '日志类型',
    retention_days INT NOT NULL COMMENT '保留天数',
    cleanup_frequency VARCHAR(20) DEFAULT 'DAILY' COMMENT '清理频率',
    cleanup_time TIME DEFAULT '03:00:00' COMMENT '清理时间',
    is_enabled TINYINT DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT COMMENT '更新人ID',
    UNIQUE KEY uk_log_type(log_type)
};

-- 系统更新记录表
CREATE TABLE system_update_logs {
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    version VARCHAR(20) NOT NULL COMMENT '版本号',
    update_type ENUM('MAJOR', 'MINOR', 'PATCH', 'HOTFIX') NOT NULL COMMENT '更新类型',
    update_time DATETIME NOT NULL COMMENT '更新时间',
    update_content TEXT NOT NULL COMMENT '更新内容',
    update_status ENUM('PENDING', 'PROCESSING', 'COMPLETED', 'FAILED', 'ROLLED_BACK') NOT NULL COMMENT '更新状态',
    operator_id BIGINT COMMENT '操作人ID',
    operator_name VARCHAR(50) COMMENT '操作人姓名',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_version(version),
    INDEX idx_time(update_time),
    INDEX idx_status(update_status)
};

-- 定时任务配置表
CREATE TABLE scheduled_tasks {
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_name VARCHAR(100) NOT NULL COMMENT '任务名称',
    task_key VARCHAR(50) NOT NULL COMMENT '任务键',
    task_group VARCHAR(50) DEFAULT 'DEFAULT' COMMENT '任务组',
    task_class VARCHAR(200) NOT NULL COMMENT '任务类',
    cron_expression VARCHAR(50) COMMENT 'CRON表达式',
    interval_seconds INT COMMENT '间隔秒数',
    task_data TEXT COMMENT '任务数据(JSON格式)',
    description VARCHAR(200) COMMENT '任务描述',
    is_enabled TINYINT DEFAULT 1 COMMENT '是否启用',
    concurrent_execution TINYINT DEFAULT 0 COMMENT '是否允许并发执行',
    last_execution_time DATETIME COMMENT '上次执行时间',
    next_execution_time DATETIME COMMENT '下次执行时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    UNIQUE KEY uk_key(task_key),
    INDEX idx_group(task_group),
    INDEX idx_enabled(is_enabled),
    INDEX idx_next_time(next_execution_time)
};

-- IP白名单表
CREATE TABLE ip_whitelist {
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ip_address VARCHAR(50) NOT NULL COMMENT 'IP地址',
    ip_description VARCHAR(100) COMMENT 'IP描述',
    access_type VARCHAR(20) DEFAULT 'ALL' COMMENT '访问类型',
    start_time DATETIME COMMENT '开始时间',
    end_time DATETIME COMMENT '结束时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    UNIQUE KEY uk_ip(ip_address)
};

-- 短信配置表
CREATE TABLE sms_settings {
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    provider_name VARCHAR(50) NOT NULL COMMENT '服务提供商',
    api_key VARCHAR(100) NOT NULL COMMENT 'API密钥',
    api_secret VARCHAR(100) COMMENT 'API密钥',
    access_token VARCHAR(200) COMMENT '访问令牌',
    sign_name VARCHAR(50) COMMENT '签名名称',
    daily_limit INT DEFAULT 1000 COMMENT '每日限额',
    monthly_limit INT DEFAULT 10000 COMMENT '每月限额',
    is_enabled TINYINT DEFAULT 1 COMMENT '是否启用',
    config_json TEXT COMMENT '配置JSON',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    updated_by BIGINT COMMENT '更新人ID',
    INDEX idx_provider(provider_name),
    INDEX idx_enabled(is_enabled)
};

-- 短信模板表
CREATE TABLE sms_templates {
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    template_code VARCHAR(50) NOT NULL COMMENT '模板编码',
    template_name VARCHAR(100) NOT NULL COMMENT '模板名称',
    template_content TEXT NOT NULL COMMENT '模板内容',
    template_type VARCHAR(20) NOT NULL COMMENT '模板类型',
    provider_template_id VARCHAR(50) COMMENT '服务商模板ID',
    provider_name VARCHAR(50) COMMENT '服务提供商',
    is_enabled TINYINT DEFAULT 1 COMMENT '是否启用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT COMMENT '创建人ID',
    updated_by BIGINT COMMENT '更新人ID',
    UNIQUE KEY uk_code(template_code),
    INDEX idx_type(template_type),
    INDEX idx_enabled(is_enabled)
};
```
