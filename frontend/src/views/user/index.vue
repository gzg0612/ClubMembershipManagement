<template>
  <div class="user-management">
    <!-- 页面标题和操作按钮 -->
    <div class="flex justify-between items-center mb-6">
      <h1 class="text-2xl font-bold">用户管理</h1>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>新增用户
      </el-button>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="bg-white p-4 rounded-lg shadow mb-6">
      <el-form :inline="true" :model="searchForm" class="flex flex-wrap gap-4">
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="用户名/姓名/手机号/邮箱"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="searchForm.role" placeholder="请选择角色" clearable>
            <el-option
              v-for="(label, value) in roleOptions"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option
              v-for="(label, value) in statusOptions"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="创建时间">
          <el-date-picker
            v-model="searchForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 用户列表 -->
    <div class="bg-white rounded-lg shadow">
      <el-table
        v-loading="loading"
        :data="userList"
        border
        style="width: 100%"
      >
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="name" label="姓名" width="120" />
        <el-table-column prop="phone" label="手机号" width="120" />
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }: { row: User }">
            <el-tag :type="getRoleType(row.role as typeof UserRole[keyof typeof UserRole])">
              {{ roleOptions[row.role as typeof UserRole[keyof typeof UserRole]] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }: { row: User }">
            <el-tag :type="getStatusType(row.status as typeof UserStatus[keyof typeof UserStatus])">
              {{ statusOptions[row.status as typeof UserStatus[keyof typeof UserStatus]] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="storeName" label="所属店铺" width="150" />
        <el-table-column prop="lastLoginAt" label="最后登录" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" link @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button type="primary" link @click="handleResetPassword(row)">
                重置密码
              </el-button>
              <el-button
                :type="row.status === UserStatus.Active ? 'danger' : 'success'"
                link
                @click="handleToggleStatus(row)"
              >
                {{ row.status === UserStatus.Active ? '禁用' : '启用' }}
              </el-button>
              <el-button type="danger" link @click="handleDelete(row)">
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="flex justify-end p-4">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>

    <!-- 用户表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增用户' : '编辑用户'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            :disabled="dialogType === 'edit'"
          />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input
            v-model="form.name"
            placeholder="请输入姓名"
          />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input
            v-model="form.phone"
            placeholder="请输入手机号"
          />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="form.email"
            placeholder="请输入邮箱"
          />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色">
            <el-option
              v-for="(label, value) in roleOptions"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="所属店铺" prop="storeId">
          <el-select v-model="form.storeId" placeholder="请选择店铺" clearable>
            <el-option
              v-for="store in storeList"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          v-if="dialogType === 'add'"
          label="密码"
          prop="password"
        >
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
        </el-form-item>
        <el-form-item
          v-if="dialogType === 'edit'"
          label="密码"
          prop="password"
        >
          <el-input
            v-model="form.password"
            type="password"
            placeholder="不修改请留空"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { userApi } from '@/api/user'
import { storeApi } from '@/api/store'
import {
  User,
  CreateUserParams,
  UpdateUserParams,
  UserQueryParams,
  UserListResponse,
  UserRole,
  UserStatus
} from '@/types/user'
import type { Store } from '@/types/store'

// 角色选项
const roleOptions: Record<typeof UserRole[keyof typeof UserRole], string> = {
  [UserRole.SuperAdmin]: '超级管理员',
  [UserRole.Admin]: '管理员',
  [UserRole.Manager]: '店长',
  [UserRole.Staff]: '店员',
  [UserRole.Technician]: '技师'
}

// 状态选项
const statusOptions: Record<typeof UserStatus[keyof typeof UserStatus], string> = {
  [UserStatus.Active]: '正常',
  [UserStatus.Inactive]: '禁用',
  [UserStatus.Pending]: '待审核'
}

// 列表数据
const loading = ref(false)
const userList = ref<User[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 搜索表单
const searchForm = reactive<UserQueryParams>({
  page: 1,
  pageSize: 20,
  keyword: '',
  role: undefined,
  status: undefined,
  dateRange: ['', '']
})

// 店铺列表
const storeList = ref<Store[]>([])

// 加载店铺列表
const loadStoreList = async () => {
  try {
    const { data } = await storeApi.getStoreList({
      page: 1,
      pageSize: 100
    })
    storeList.value = data.items
  } catch (error) {
    console.error('获取店铺列表失败:', error)
    ElMessage.error('获取店铺列表失败')
  }
}

// 获取用户列表
const getUserList = async () => {
  loading.value = true
  try {
    const params: UserQueryParams = {
      ...searchForm,
      page: currentPage.value,
      pageSize: pageSize.value,
      startDate: searchForm.dateRange?.[0],
      endDate: searchForm.dateRange?.[1]
    }
    const response = await userApi.getUserList(params)
    userList.value = response.items
    total.value = response.total
  } catch (error) {
    console.error('获取用户列表失败:', error)
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  getUserList()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.role = undefined
  searchForm.status = undefined
  searchForm.dateRange = ['', '']
  handleSearch()
}

// 分页相关
const handleSizeChange = (val: number) => {
  pageSize.value = val
  getUserList()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  getUserList()
}

// 表单相关
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const formRef = ref<FormInstance>()
const form = reactive<Partial<CreateUserParams> & { id?: number }>({
  username: '',
  password: '',
  name: '',
  phone: '',
  email: '',
  role: UserRole.Staff,
  storeId: undefined
})

// 表单验证规则
const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 新增用户
const handleAdd = () => {
  dialogType.value = 'add'
  Object.assign(form, {
    username: '',
    password: '',
    name: '',
    phone: '',
    email: '',
    role: UserRole.Staff,
    storeId: undefined
  })
  dialogVisible.value = true
}

// 编辑用户
const handleEdit = (row: User) => {
  dialogType.value = 'edit'
  Object.assign(form, {
    id: row.id,
    name: row.name,
    phone: row.phone,
    email: row.email,
    role: row.role,
    storeId: row.storeId,
    password: ''
  })
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          await userApi.createUser(form as CreateUserParams)
          ElMessage.success('创建成功')
        } else {
          const updateData: UpdateUserParams = {
            name: form.name,
            phone: form.phone,
            email: form.email,
            role: form.role,
            storeId: form.storeId
          }
          if (form.password) {
            updateData.password = form.password
          }
          await userApi.updateUser(form.id!, updateData)
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        getUserList()
      } catch (error) {
        console.error('保存用户失败:', error)
        ElMessage.error('保存用户失败')
      }
    }
  })
}

// 重置密码
const handleResetPassword = async (row: User) => {
  try {
    await ElMessageBox.confirm('确定要重置该用户的密码吗？', '提示', {
      type: 'warning'
    })
    await userApi.resetPassword(row.id)
    ElMessage.success('密码重置成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('重置密码失败:', error)
      ElMessage.error('重置密码失败')
    }
  }
}

// 启用/禁用用户
const handleToggleStatus = async (row: User) => {
  try {
    const newStatus = row.status === UserStatus.Active ? UserStatus.Inactive : UserStatus.Active
    const action = newStatus === UserStatus.Active ? '启用' : '禁用'
    await ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', {
      type: 'warning'
    })
    await userApi.toggleUserStatus(row.id, newStatus)
    ElMessage.success(`${action}成功`)
    getUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('操作失败:', error)
      ElMessage.error('操作失败')
    }
  }
}

// 删除用户
const handleDelete = async (row: User) => {
  try {
    await ElMessageBox.confirm('确定要删除该用户吗？', '提示', {
      type: 'warning'
    })
    await userApi.deleteUser(row.id)
    ElMessage.success('删除成功')
    getUserList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除用户失败:', error)
      ElMessage.error('删除用户失败')
    }
  }
}

// 获取角色样式
const getRoleType = (role: typeof UserRole[keyof typeof UserRole]) => {
  const types: Record<typeof UserRole[keyof typeof UserRole], string> = {
    [UserRole.SuperAdmin]: 'danger',
    [UserRole.Admin]: 'warning',
    [UserRole.Manager]: 'success',
    [UserRole.Staff]: 'info',
    [UserRole.Technician]: ''
  }
  return types[role]
}

// 获取状态样式
const getStatusType = (status: typeof UserStatus[keyof typeof UserStatus]) => {
  const types: Record<typeof UserStatus[keyof typeof UserStatus], string> = {
    [UserStatus.Active]: 'success',
    [UserStatus.Inactive]: 'danger',
    [UserStatus.Pending]: 'warning'
  }
  return types[status]
}

onMounted(() => {
  loadStoreList()
  getUserList()
})
</script>

<style scoped>
.user-management {
  padding: 20px;
}
</style> 