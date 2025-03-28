<template>
  <div class="import-preview">
    <el-dialog
      v-model="dialogVisible"
      title="导入预览"
      width="80%"
      :close-on-click-modal="false"
      :before-close="handleClose"
    >
      <div class="preview-content">
        <!-- 步骤条 -->
        <el-steps :active="currentStep" finish-status="success" class="mb-4">
          <el-step title="选择文件" />
          <el-step title="字段映射" />
          <el-step title="数据预览" />
          <el-step title="导入确认" />
        </el-steps>

        <!-- 步骤1：选择文件 -->
        <div v-if="currentStep === 0" class="step-content">
          <div class="template-download mb-4">
            <span class="mr-2">下载导入模板：</span>
            <el-button type="primary" link @click="downloadTemplate(ImportTemplateType.BASIC)">基础模板</el-button>
            <el-button type="primary" link @click="downloadTemplate(ImportTemplateType.DETAILED)">详细模板</el-button>
          </div>
          <el-upload
            class="upload-demo"
            drag
            action="#"
            :auto-upload="false"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            accept=".xlsx,.xls"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                只能上传 xlsx/xls 文件，且文件大小不超过 10MB
              </div>
            </template>
          </el-upload>
        </div>

        <!-- 步骤2：字段映射 -->
        <div v-if="currentStep === 1" class="step-content">
          <el-form :model="mappingForm" label-width="120px">
            <el-form-item
              v-for="field in requiredFields"
              :key="field.key"
              :label="field.label"
              :required="field.required"
            >
              <el-select
                v-model="mappingForm[field.key]"
                placeholder="请选择对应的列"
                clearable
              >
                <el-option
                  v-for="header in fileHeaders"
                  :key="header"
                  :label="header"
                  :value="header"
                />
              </el-select>
            </el-form-item>
          </el-form>
        </div>

        <!-- 步骤3：数据预览 -->
        <div v-if="currentStep === 2" class="step-content">
          <el-alert
            v-if="previewErrors.length > 0"
            type="warning"
            :closable="false"
            class="mb-4"
          >
            <p>发现 {{ previewErrors.length }} 个潜在问题：</p>
            <ul>
              <li v-for="(error, index) in previewErrors" :key="index">
                第 {{ error.row }} 行：{{ error.message }}
              </li>
            </ul>
          </el-alert>

          <el-table
            :data="previewData"
            border
            height="400"
            style="width: 100%"
          >
            <el-table-column
              v-for="field in requiredFields"
              :key="field.key"
              :prop="field.key"
              :label="field.label"
              min-width="120"
            />
            <el-table-column
              fixed="right"
              label="状态"
              width="100"
            >
              <template #default="{ row }">
                <el-tag
                  :type="row.valid ? 'success' : 'danger'"
                  size="small"
                >
                  {{ row.valid ? '有效' : '无效' }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 步骤4：导入确认 -->
        <div v-if="currentStep === 3" class="step-content">
          <el-result
            :icon="importStatus.icon"
            :title="importStatus.title"
            :sub-title="importStatus.subTitle"
          >
            <template #extra>
              <div v-if="importStatus.showRetry">
                <el-button type="primary" @click="handleRetry">重试</el-button>
                <el-button @click="handleClose">关闭</el-button>
              </div>
              <div v-else>
                <el-button type="primary" @click="handleClose">完成</el-button>
              </div>
            </template>
          </el-result>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleClose">取消</el-button>
          <el-button
            v-if="currentStep > 0 && currentStep < 3"
            @click="handlePrevStep"
          >
            上一步
          </el-button>
          <el-button
            v-if="currentStep < 3"
            type="primary"
            @click="handleNextStep"
            :loading="loading"
            :disabled="!canProceed"
          >
            {{ currentStep === 2 ? '开始导入' : '下一步' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue';
import { ElMessage } from 'element-plus';
import { UploadFilled } from '@element-plus/icons-vue';
import type { UploadFile } from 'element-plus';
import { downloadImportTemplate, previewImport, confirmImport } from '@/api/finance';
import { ImportTemplateType, type ImportPreviewResponse, type ImportConfirmResponse } from '@/types/finance';

// Props & Emits
const props = defineProps<{
  modelValue: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void;
  (e: 'success'): void;
}>();

// 状态
const currentStep = ref(0);
const loading = ref(false);
const selectedFile = ref<File | null>(null);
const fileHeaders = ref<string[]>([]);
const previewData = ref<Array<Record<string, any>>>([]);
const previewErrors = ref<Array<{ row: number; message: string }>>([]);
const fileId = ref('');

// 字段映射表单
const mappingForm = reactive<Record<string, string>>({});

// 必填字段定义
const requiredFields = [
  { key: 'type', label: '交易类型', required: true },
  { key: 'amount', label: '金额', required: true },
  { key: 'paymentMethod', label: '支付方式', required: true },
  { key: 'status', label: '状态', required: false },
  { key: 'remark', label: '备注', required: false }
];

// 导入状态
const importStatus = reactive({
  icon: 'success' as 'success' | 'error',
  title: '',
  subTitle: '',
  showRetry: false
});

// 计算属性
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const canProceed = computed(() => {
  switch (currentStep.value) {
    case 0:
      return !!selectedFile.value;
    case 1:
      return requiredFields
        .filter(field => field.required)
        .every(field => !!mappingForm[field.key]);
    case 2:
      return previewData.value.length > 0;
    default:
      return true;
  }
});

// 方法
const handleClose = () => {
  dialogVisible.value = false;
  currentStep.value = 0;
  selectedFile.value = null;
  fileHeaders.value = [];
  previewData.value = [];
  previewErrors.value = [];
  fileId.value = '';
  Object.keys(mappingForm).forEach(key => {
    mappingForm[key] = '';
  });
};

const downloadTemplate = async (type: ImportTemplateType) => {
  try {
    const blob = await downloadImportTemplate(type);
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `交易记录导入模板_${type}.xlsx`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  } catch (error) {
    console.error('下载模板失败:', error);
    ElMessage.error('下载模板失败');
  }
};

const handleFileChange = async (file: UploadFile) => {
  if (!file.raw) return;
  
  selectedFile.value = file.raw;
  loading.value = true;
  
  try {
    const res = await previewImport(file.raw) as ImportPreviewResponse;
    fileHeaders.value = res.headers;
    fileId.value = res.fileId;
    
    // 自动映射字段
    Object.entries(res.mapping).forEach(([key, value]) => {
      if (requiredFields.some(field => field.key === key)) {
        mappingForm[key] = value;
      }
    });
  } catch (error) {
    console.error('预览文件失败:', error);
    ElMessage.error('预览文件失败');
    selectedFile.value = null;
  } finally {
    loading.value = false;
  }
};

const handleFileRemove = () => {
  selectedFile.value = null;
};

const handlePrevStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--;
  }
};

const handleNextStep = async () => {
  if (!canProceed.value) return;
  
  loading.value = true;
  
  try {
    switch (currentStep.value) {
      case 0:
        currentStep.value++;
        break;
        
      case 1:
        // 预览数据
        if (!selectedFile.value) return;
        const previewRes = await previewImport(selectedFile.value) as ImportPreviewResponse;
        previewData.value = previewRes.preview;
        previewErrors.value = previewRes.errors || [];
        currentStep.value++;
        break;
        
      case 2:
        // 确认导入
        const importRes = await confirmImport({
          fileId: fileId.value,
          mapping: mappingForm,
          skipErrors: previewErrors.value.length > 0
        }) as ImportConfirmResponse;
        
        importStatus.icon = 'success';
        importStatus.title = '导入成功';
        importStatus.subTitle = `成功导入 ${importRes.success} 条记录，失败 ${importRes.failed} 条`;
        importStatus.showRetry = false;
        
        currentStep.value++;
        emit('success');
        break;
    }
  } catch (error) {
    console.error('操作失败:', error);
    if (currentStep.value === 2) {
      importStatus.icon = 'error';
      importStatus.title = '导入失败';
      importStatus.subTitle = '请检查数据后重试';
      importStatus.showRetry = true;
      currentStep.value++;
    } else {
      ElMessage.error('操作失败');
    }
  } finally {
    loading.value = false;
  }
};

const handleRetry = () => {
  currentStep.value = 0;
  selectedFile.value = null;
  fileHeaders.value = [];
  previewData.value = [];
  previewErrors.value = [];
  fileId.value = '';
  Object.keys(mappingForm).forEach(key => {
    mappingForm[key] = '';
  });
};
</script>

<style scoped>
.import-preview {
  .preview-content {
    padding: 20px 0;
  }
  
  .step-content {
    min-height: 300px;
  }
  
  .template-download {
    display: flex;
    align-items: center;
  }
  
  .mr-2 {
    margin-right: 8px;
  }
  
  .mb-4 {
    margin-bottom: 16px;
  }
  
  .upload-demo {
    text-align: center;
  }
  
  .el-upload__tip {
    margin-top: 10px;
    color: #909399;
  }
}
</style> 