<template>
  <div class="store-staff">
    <!-- 页面标题和操作按钮 -->
    <div class="flex justify-between items-center mb-6">
      <div class="flex items-center gap-4">
        <el-button @click="router.back()">
          <el-icon><ArrowLeft /></el-icon>返回
        </el-button>
        <h1 class="text-2xl font-bold">{{ storeName }} - 员工管理</h1>
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>新增员工
      </el-button>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="bg-white p-4 rounded-lg shadow mb-6">
      <el-form :inline="true" :model="searchForm" class="flex flex-wrap gap-4">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="员工姓名/编号" clearable />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="searchForm.role" placeholder="请选择角色" clearable>
            <el-option
              v-for="(label, value) in staffRoleOptions"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option
              v-for="(label, value) in staffStatusOptions"
              :key="value"
              :label="label"
              :value="Number(value)"
            />
          </el-select>
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

    <!-- 员工列表 -->
    <div class="bg-white rounded-lg shadow">
      <el-table
        v-loading="loading"
        :data="staffList"
        border
        style="width: 100%"
      >
        <el-table-column prop="userId" label="员工编号" width="120" />
        <el-table-column prop="role" label="角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.role)">
              {{ staffRoleOptions[row.role] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ staffStatusOptions[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="startDate" label="入职日期" width="120" />
        <el-table-column prop="endDate" label="离职日期" width="120" />
        <el-table-column prop="remark" label="备注" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" link @click="handleEdit(row)">
                <el-icon><Edit /></el-icon>编辑
              </el-button>
              <el-button type="danger" link @click="handleRemove(row)">
                <el-icon><Delete /></el-icon>移除
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

    <!-- 员工表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增员工' : '编辑员工'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="员工编号" prop="userId">
          <el-input v-model="form.userId" placeholder="请输入员工编号" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" placeholder="请选择角色">
            <el-option
              v-for="(label, value) in staffRoleOptions"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择状态">
            <el-option
              v-for="(label, value) in staffStatusOptions"
              :key="value"
              :label="label"
              :value="Number(value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="入职日期" prop="startDate">
          <el-date-picker
            v-model="form.startDate"
            type="date"
            placeholder="请选择入职日期"
          />
        </el-form-item>
        <el-form-item label="离职日期" prop="endDate">
          <el-date-picker
            v-model="form.endDate"
            type="date"
            placeholder="请选择离职日期"
          />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="3"
            placeholder="请输入备注信息"
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
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh, Edit, Delete, ArrowLeft } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { storeApi } from '@/api/store'
import type { StoreStaff, StoreStaffRole, StoreStaffStatus } from '@/types/store'

const router = useRouter()
const route = useRoute()
const storeId = Number(route.params.id)
const storeName = ref('')

// 加载店铺信息
const loadStoreInfo = async () => {
  try {
    const { data } = await storeApi.getStore(storeId)
    storeName.value = data.name
  } catch (error) {
    console.error('获取店铺信息失败:', error)
    ElMessage.error('获取店铺信息失败')
  }
}

const loading = ref(false)
const staffList = ref<StoreStaff[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  role: undefined as StoreStaffRole | undefined,
  status: undefined as number | undefined
})

// 员工角色选项
const staffRoleOptions: Record<StoreStaffRole, string> = {
  [StoreStaffRole.Manager]: '店长',
  [StoreStaffRole.Staff]: '店员',
  [StoreStaffRole.Technician]: '技师'
}

// 员工状态选项
const staffStatusOptions: Record<StoreStaffStatus, string> = {
  [StoreStaffStatus.Active]: '在职',
  [StoreStaffStatus.Inactive]: '离职',
  [StoreStaffStatus.OnLeave]: '请假'
}

// 表单相关
const dialogVisible = ref(false)
const dialogType = ref<'add' | 'edit'>('add')
const formRef = ref<FormInstance>()
const form = reactive<Partial<StoreStaff>>({
  userId: undefined,
  role: StoreStaffRole.Staff,
  status: StoreStaffStatus.Active,
  startDate: '',
  endDate: undefined,
  remark: ''
})

// 表单验证规则
const rules: FormRules = {
  userId: [{ required: true, message: '请输入员工编号', trigger: 'blur' }],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
  startDate: [{ required: true, message: '请选择入职日期', trigger: 'change' }]
}

// 获取员工列表
const getStaffList = async () => {
  loading.value = true
  try {
    const { data } = await storeApi.staff.getStaffList(storeId, {
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchForm.keyword,
      status: searchForm.status
    })
    staffList.value = data.items
    total.value = data.total
  } catch (error) {
    console.error('获取员工列表失败:', error)
    ElMessage.error('获取员工列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  getStaffList()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.role = undefined
  searchForm.status = undefined
  handleSearch()
}

// 新增员工
const handleAdd = () => {
  dialogType.value = 'add'
  Object.assign(form, {
    userId: undefined,
    role: StoreStaffRole.Staff,
    status: StoreStaffStatus.Active,
    startDate: '',
    endDate: undefined,
    remark: ''
  })
  dialogVisible.value = true
}

// 编辑员工
const handleEdit = (row: StoreStaff) => {
  dialogType.value = 'edit'
  Object.assign(form, row)
  dialogVisible.value = true
}

// 移除员工
const handleRemove = (row: StoreStaff) => {
  ElMessageBox.confirm('确定要移除该员工吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await storeApi.staff.removeStaff(storeId, row.id)
      ElMessage.success('移除成功')
      getStaffList()
    } catch (error) {
      console.error('移除员工失败:', error)
      ElMessage.error('移除员工失败')
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        if (dialogType.value === 'add') {
          await storeApi.staff.addStaff(storeId, form)
          ElMessage.success('添加成功')
        } else {
          await storeApi.staff.updateStaff(storeId, form.id!, form)
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        getStaffList()
      } catch (error) {
        console.error('保存员工信息失败:', error)
        ElMessage.error('保存员工信息失败')
      }
    }
  })
}

// 获取角色类型
const getRoleType = (role: StoreStaffRole) => {
  const types: Record<StoreStaffRole, string> = {
    [StoreStaffRole.Manager]: 'danger',
    [StoreStaffRole.Staff]: 'success',
    [StoreStaffRole.Technician]: 'warning'
  }
  return types[role]
}

// 获取状态类型
const getStatusType = (status: StoreStaffStatus) => {
  const types: Record<StoreStaffStatus, string> = {
    [StoreStaffStatus.Active]: 'success',
    [StoreStaffStatus.Inactive]: 'danger',
    [StoreStaffStatus.OnLeave]: 'warning'
  }
  return types[status]
}

// 分页相关
const handleSizeChange = (val: number) => {
  pageSize.value = val
  getStaffList()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  getStaffList()
}

onMounted(() => {
  loadStoreInfo()
  getStaffList()
})
</script>

<style scoped>
.store-staff {
  padding: 20px;
}
</style> 