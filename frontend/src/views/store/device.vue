<template>
  <div class="store-device">
    <!-- 页面标题和操作按钮 -->
    <div class="flex justify-between items-center mb-6">
      <div class="flex items-center gap-4">
        <el-button @click="router.back()">
          <el-icon><ArrowLeft /></el-icon>返回
        </el-button>
        <h1 class="text-2xl font-bold">{{ storeName }} - 设备管理</h1>
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>新增设备
      </el-button>
    </div>

    <!-- 搜索和筛选区域 -->
    <div class="bg-white p-4 rounded-lg shadow mb-6">
      <el-form :inline="true" :model="searchForm" class="flex flex-wrap gap-4">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="设备名称/编号" clearable />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="请选择类型" clearable>
            <el-option
              v-for="(label, value) in deviceTypeOptions"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option
              v-for="(label, value) in deviceStatusOptions"
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

    <!-- 设备列表 -->
    <div class="bg-white rounded-lg shadow">
      <el-table
        v-loading="loading"
        :data="deviceList"
        border
        style="width: 100%"
      >
        <el-table-column prop="name" label="设备名称" min-width="150" />
        <el-table-column prop="code" label="设备编号" width="120" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeType(row.type)">
              {{ deviceTypeOptions[row.type] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ deviceStatusOptions[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastCheck" label="上次检查" width="120" />
        <el-table-column prop="nextCheck" label="下次检查" width="120" />
        <el-table-column prop="remark" label="备注" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" link @click="handleEdit(row)">
                <el-icon><Edit /></el-icon>编辑
              </el-button>
              <el-button type="success" link @click="handleCheck(row)">
                <el-icon><Check /></el-icon>检查
              </el-button>
              <el-button type="danger" link @click="handleDelete(row)">
                <el-icon><Delete /></el-icon>删除
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

    <!-- 设备表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'add' ? '新增设备' : '编辑设备'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="设备名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入设备名称" />
        </el-form-item>
        <el-form-item label="设备编号" prop="code">
          <el-input v-model="form.code" placeholder="请输入设备编号" />
        </el-form-item>
        <el-form-item label="设备类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择类型">
            <el-option
              v-for="(label, value) in deviceTypeOptions"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="设备状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择状态">
            <el-option
              v-for="(label, value) in deviceStatusOptions"
              :key="value"
              :label="label"
              :value="Number(value)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="上次检查" prop="lastCheck">
          <el-date-picker
            v-model="form.lastCheck"
            type="date"
            placeholder="请选择上次检查日期"
          />
        </el-form-item>
        <el-form-item label="下次检查" prop="nextCheck">
          <el-date-picker
            v-model="form.nextCheck"
            type="date"
            placeholder="请选择下次检查日期"
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

    <!-- 设备检查对话框 -->
    <el-dialog
      v-model="checkDialogVisible"
      title="设备检查"
      width="500px"
    >
      <el-form
        ref="checkFormRef"
        :model="checkForm"
        :rules="checkRules"
        label-width="100px"
      >
        <el-form-item label="检查类型" prop="checkType">
          <el-select v-model="checkForm.checkType" placeholder="请选择检查类型">
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
            v-model="checkForm.content"
            type="textarea"
            :rows="3"
            placeholder="请输入检查内容"
          />
        </el-form-item>
        <el-form-item label="检查结果" prop="result">
          <el-input
            v-model="checkForm.result"
            type="textarea"
            :rows="3"
            placeholder="请输入检查结果"
          />
        </el-form-item>
        <el-form-item label="检查状态" prop="status">
          <el-select v-model="checkForm.status" placeholder="请选择检查状态">
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
          <el-button @click="checkDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCheckSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh, Edit, Delete, ArrowLeft, Check } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { storeApi } from '@/api/store'
import type {
  StoreDevice,
  DeviceType,
  DeviceStatus,
  CheckType,
  CheckStatus,
  StoreDeviceCheck
} from '@/types/store'

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
const deviceList = ref<StoreDevice[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 搜索表单
const searchForm = reactive({
  keyword: '',
  type: undefined as DeviceType | undefined,
  status: undefined as number | undefined
})

// 设备类型选项
const deviceTypeOptions: Record<DeviceType, string> = {
  [DeviceType.Equipment]: '设备',
  [DeviceType.Tool]: '工具',
  [DeviceType.Furniture]: '家具'
}

// 设备状态选项
const deviceStatusOptions: Record<DeviceStatus, string> = {
  [DeviceStatus.Normal]: '正常',
  [DeviceStatus.Repair]: '维修中',
  [DeviceStatus.Scrapped]: '已报废'
}

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
const dialogType = ref<'add' | 'edit'>('add')
const formRef = ref<FormInstance>()
const form = reactive<Partial<StoreDevice>>({
  name: '',
  code: '',
  type: DeviceType.Equipment,
  status: DeviceStatus.Normal,
  lastCheck: '',
  nextCheck: '',
  remark: ''
})

// 表单验证规则
const rules: FormRules = {
  name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
  code: [{ required: true, message: '请输入设备编号', trigger: 'blur' }],
  type: [{ required: true, message: '请选择设备类型', trigger: 'change' }],
  status: [{ required: true, message: '请选择设备状态', trigger: 'change' }],
  lastCheck: [{ required: true, message: '请选择上次检查日期', trigger: 'change' }],
  nextCheck: [{ required: true, message: '请选择下次检查日期', trigger: 'change' }]
}

// 检查表单相关
const checkDialogVisible = ref(false)
const checkFormRef = ref<FormInstance>()
const currentDevice = ref<StoreDevice>()
const checkForm = reactive<Partial<StoreDeviceCheck>>({
  checkType: CheckType.Daily,
  content: '',
  result: '',
  status: CheckStatus.Pass
})

// 检查表单验证规则
const checkRules: FormRules = {
  checkType: [{ required: true, message: '请选择检查类型', trigger: 'change' }],
  content: [{ required: true, message: '请输入检查内容', trigger: 'blur' }],
  result: [{ required: true, message: '请输入检查结果', trigger: 'blur' }],
  status: [{ required: true, message: '请选择检查状态', trigger: 'change' }]
}

// 获取设备列表
const getDeviceList = async () => {
  loading.value = true
  try {
    const { data } = await storeApi.device.getDeviceList(storeId, {
      page: currentPage.value,
      pageSize: pageSize.value,
      keyword: searchForm.keyword,
      status: searchForm.status
    })
    deviceList.value = data.items
    total.value = data.total
  } catch (error) {
    console.error('获取设备列表失败:', error)
    ElMessage.error('获取设备列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  getDeviceList()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.type = undefined
  searchForm.status = undefined
  handleSearch()
}

// 新增设备
const handleAdd = () => {
  dialogType.value = 'add'
  Object.assign(form, {
    name: '',
    code: '',
    type: DeviceType.Equipment,
    status: DeviceStatus.Normal,
    lastCheck: '',
    nextCheck: '',
    remark: ''
  })
  dialogVisible.value = true
}

// 编辑设备
const handleEdit = (row: StoreDevice) => {
  dialogType.value = 'edit'
  Object.assign(form, row)
  dialogVisible.value = true
}

// 删除设备
const handleDelete = (row: StoreDevice) => {
  ElMessageBox.confirm('确定要删除该设备吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await storeApi.device.deleteDevice(storeId, row.id)
      ElMessage.success('删除成功')
      getDeviceList()
    } catch (error) {
      console.error('删除设备失败:', error)
      ElMessage.error('删除设备失败')
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
          await storeApi.device.addDevice(storeId, form)
          ElMessage.success('添加成功')
        } else {
          await storeApi.device.updateDevice(storeId, form.id!, form)
          ElMessage.success('更新成功')
        }
        dialogVisible.value = false
        getDeviceList()
      } catch (error) {
        console.error('保存设备信息失败:', error)
        ElMessage.error('保存设备信息失败')
      }
    }
  })
}

// 打开检查对话框
const handleCheck = (row: StoreDevice) => {
  currentDevice.value = row
  Object.assign(checkForm, {
    checkType: CheckType.Daily,
    content: '',
    result: '',
    status: CheckStatus.Pass
  })
  checkDialogVisible.value = true
}

// 提交检查表单
const handleCheckSubmit = async () => {
  if (!checkFormRef.value || !currentDevice.value) return
  await checkFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await storeApi.device.check.addCheck(storeId, currentDevice.value.id, checkForm)
        ElMessage.success('添加检查记录成功')
        checkDialogVisible.value = false
        getDeviceList()
      } catch (error) {
        console.error('添加检查记录失败:', error)
        ElMessage.error('添加检查记录失败')
      }
    }
  })
}

// 获取类型样式
const getTypeType = (type: DeviceType) => {
  const types: Record<DeviceType, string> = {
    [DeviceType.Equipment]: 'primary',
    [DeviceType.Tool]: 'success',
    [DeviceType.Furniture]: 'warning'
  }
  return types[type]
}

// 获取状态样式
const getStatusType = (status: DeviceStatus) => {
  const types: Record<DeviceStatus, string> = {
    [DeviceStatus.Normal]: 'success',
    [DeviceStatus.Repair]: 'warning',
    [DeviceStatus.Scrapped]: 'danger'
  }
  return types[status]
}

// 分页相关
const handleSizeChange = (val: number) => {
  pageSize.value = val
  getDeviceList()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  getDeviceList()
}

onMounted(() => {
  loadStoreInfo()
  getDeviceList()
})
</script>

<style scoped>
.store-device {
  padding: 20px;
}
</style> 