<template>
  <div class="device-check">
    <!-- 页面标题和操作按钮 -->
    <div class="flex justify-between items-center mb-6">
      <div class="flex items-center gap-4">
        <el-button @click="router.back()">
          <el-icon><ArrowLeft /></el-icon>返回
        </el-button>
        <h1 class="text-2xl font-bold">{{ storeName }} - {{ deviceName }} - 检查记录</h1>
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>新增检查
      </el-button>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="bg-white p-4 rounded-lg shadow mb-6">
      <el-form :inline="true" :model="searchForm" class="flex flex-wrap gap-4">
        <el-form-item label="检查类型">
          <el-select v-model="searchForm.checkType" placeholder="请选择检查类型" clearable>
            <el-option
              v-for="(label, value) in checkTypeOptions"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="检查状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option
              v-for="(label, value) in checkStatusOptions"
              :key="value"
              :label="label"
              :value="Number(value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="检查日期">
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

    <!-- 检查记录列表 -->
    <div class="bg-white rounded-lg shadow">
      <el-table
        v-loading="loading"
        :data="checkList"
        border
        style="width: 100%"
      >
        <el-table-column prop="checkType" label="检查类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getCheckTypeType(row.checkType)">
              {{ checkTypeOptions[row.checkType] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="检查状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ checkStatusOptions[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="检查内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="result" label="检查结果" min-width="200" show-overflow-tooltip />
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column prop="createdAt" label="检查时间" width="180" />
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

    <!-- 检查记录表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="新增检查记录"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="检查类型" prop="checkType">
          <el-select v-model="form.checkType" placeholder="请选择检查类型">
            <el-option
              v-for="(label, value) in checkTypeOptions"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="检查内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="3"
            placeholder="请输入检查内容"
          />
        </el-form-item>
        <el-form-item label="检查结果" prop="result">
          <el-input
            v-model="form.result"
            type="textarea"
            :rows="3"
            placeholder="请输入检查结果"
          />
        </el-form-item>
        <el-form-item label="检查状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择检查状态">
            <el-option
              v-for="(label, value) in checkStatusOptions"
              :key="value"
              :label="label"
              :value="Number(value)"
            />
          </el-select>
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
import { ElMessage } from 'element-plus'
import { Plus, Search, Refresh, ArrowLeft } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { storeApi } from '@/api/store'
import type { CheckType, CheckStatus, StoreDeviceCheck } from '@/types/store'

const router = useRouter()
const route = useRoute()
const storeId = Number(route.params.id)
const deviceId = Number(route.params.deviceId)
const storeName = ref('')
const deviceName = ref('')

// 加载店铺和设备信息
const loadInfo = async () => {
  try {
    const [storeRes, deviceRes] = await Promise.all([
      storeApi.getStore(storeId),
      storeApi.device.getDevice(storeId, deviceId)
    ])
    storeName.value = storeRes.data.name
    deviceName.value = deviceRes.data.name
  } catch (error) {
    console.error('获取信息失败:', error)
    ElMessage.error('获取信息失败')
  }
}

const loading = ref(false)
const checkList = ref<StoreDeviceCheck[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 搜索表单
const searchForm = reactive({
  checkType: undefined as CheckType | undefined,
  status: undefined as number | undefined,
  dateRange: [] as string[]
})

// 检查类型选项
const checkTypeOptions: Record<CheckType, string> = {
  [CheckType.Daily]: '日常检查',
  [CheckType.Weekly]: '周检查',
  [CheckType.Monthly]: '月检查',
  [CheckType.Quarterly]: '季度检查'
}

// 检查状态选项
const checkStatusOptions: Record<CheckStatus, string> = {
  [CheckStatus.Pass]: '通过',
  [CheckStatus.Fail]: '不通过',
  [CheckStatus.Pending]: '待处理'
}

// 表单相关
const dialogVisible = ref(false)
const formRef = ref<FormInstance>()
const form = reactive<Partial<StoreDeviceCheck>>({
  checkType: CheckType.Daily,
  content: '',
  result: '',
  status: CheckStatus.Pass
})

// 表单验证规则
const rules: FormRules = {
  checkType: [{ required: true, message: '请选择检查类型', trigger: 'change' }],
  content: [{ required: true, message: '请输入检查内容', trigger: 'blur' }],
  result: [{ required: true, message: '请输入检查结果', trigger: 'blur' }],
  status: [{ required: true, message: '请选择检查状态', trigger: 'change' }]
}

// 获取检查记录列表
const getCheckList = async () => {
  loading.value = true
  try {
    const { data } = await storeApi.device.check.getCheckList(storeId, deviceId, {
      page: currentPage.value,
      pageSize: pageSize.value,
      checkType: searchForm.checkType,
      status: searchForm.status,
      startDate: searchForm.dateRange[0],
      endDate: searchForm.dateRange[1]
    })
    checkList.value = data.items
    total.value = data.total
  } catch (error) {
    console.error('获取检查记录列表失败:', error)
    ElMessage.error('获取检查记录列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  getCheckList()
}

// 重置
const handleReset = () => {
  searchForm.checkType = undefined
  searchForm.status = undefined
  searchForm.dateRange = []
  handleSearch()
}

// 新增检查记录
const handleAdd = () => {
  Object.assign(form, {
    checkType: CheckType.Daily,
    content: '',
    result: '',
    status: CheckStatus.Pass
  })
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await storeApi.device.check.addCheck(storeId, deviceId, form)
        ElMessage.success('添加成功')
        dialogVisible.value = false
        getCheckList()
      } catch (error) {
        console.error('添加检查记录失败:', error)
        ElMessage.error('添加检查记录失败')
      }
    }
  })
}

// 获取检查类型样式
const getCheckTypeType = (type: CheckType) => {
  const types: Record<CheckType, string> = {
    [CheckType.Daily]: 'info',
    [CheckType.Weekly]: 'success',
    [CheckType.Monthly]: 'warning',
    [CheckType.Quarterly]: 'danger'
  }
  return types[type]
}

// 获取状态样式
const getStatusType = (status: CheckStatus) => {
  const types: Record<CheckStatus, string> = {
    [CheckStatus.Pass]: 'success',
    [CheckStatus.Fail]: 'danger',
    [CheckStatus.Pending]: 'warning'
  }
  return types[status]
}

// 分页相关
const handleSizeChange = (val: number) => {
  pageSize.value = val
  getCheckList()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  getCheckList()
}

onMounted(() => {
  loadInfo()
  getCheckList()
})
</script>

<style scoped>
.device-check {
  padding: 20px;
}
</style> 