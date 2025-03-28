# 会员管理系统前端

## 项目介绍
这是一个基于Vue 3 + Element Plus的会员管理系统前端项目。

## 技术栈
- Vue 3
- Vue Router
- Vuex
- Element Plus
- Axios
- Sass

## 开发环境
- Node.js >= 14.0.0
- npm >= 6.0.0

## 项目设置
```bash
# 安装依赖
npm install

# 启动开发服务器
npm run serve

# 构建生产环境
npm run build:prod

# 构建开发环境
npm run build:dev

# 代码检查
npm run lint
```

## 项目结构
```
frontend/
├── public/                 # 静态资源
├── src/                    # 源代码
│   ├── api/               # API接口
│   ├── assets/            # 资源文件
│   ├── components/        # 公共组件
│   ├── router/            # 路由配置
│   ├── store/             # Vuex状态管理
│   ├── styles/            # 全局样式
│   ├── utils/             # 工具函数
│   ├── views/             # 页面组件
│   ├── App.vue            # 根组件
│   └── main.js            # 入口文件
├── .env                   # 环境变量
├── .env.development       # 开发环境变量
├── .env.production        # 生产环境变量
├── .eslintrc.js          # ESLint配置
├── .gitignore            # Git忽略文件
├── babel.config.js       # Babel配置
├── package.json          # 项目配置
└── vue.config.js         # Vue配置
```

## 功能模块
- 用户认证
  - 登录
  - 注册
  - 找回密码
- 会员管理
  - 会员列表
  - 会员详情
  - 会员编辑
  - 会员删除
- 交易管理
  - 充值
  - 消费
  - 交易记录
- 数据统计
  - 会员统计
  - 交易统计
  - 趋势分析

## 开发规范
- 遵循Vue 3组合式API风格
- 使用ESLint进行代码检查
- 使用Sass进行样式开发
- 组件和文件使用PascalCase命名
- 变量和函数使用camelCase命名

## 部署说明
1. 修改生产环境配置文件
2. 执行构建命令
3. 将dist目录下的文件部署到Web服务器

## 注意事项
- 开发时请确保后端API服务已启动
- 提交代码前请执行代码检查
- 保持代码风格统一
- 及时更新依赖包版本 