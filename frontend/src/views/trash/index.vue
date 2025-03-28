<template>
  <div class="trash-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>数据回收站</span>
          <el-select v-model="recordType" placeholder="选择记录类型" clearable>
            <el-option label="会员" :value="RecordType.MEMBER" />
            <el-option label="商品" :value="RecordType.PRODUCT" />
            <el-option label="活动" :value="RecordType.ACTIVITY" />
          </el-select>
        </div>
      </template>

      <el-table :data="records" v-loading="loading" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="recordType" label="记录类型" width="100">
          <template #default="{ row }">
            {{ getRecordTypeLabel(row.recordType) }}
          </template>
        </el-table-column>
        <el-table-column prop="recordId" label="记录ID" width="100" />
        <el-table-column prop="deleteReason" label="删除原因" />
        <el-table-column prop="operatorName" label="操作人" width="120" />
        <el-table-column prop="deletedAt" label="删除时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.deletedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              查看
            </el-button>
            <el-button type="success" link @click="handleRestore(row)">
              恢复
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">
              永久删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 查看记录详情对话框 -->
    <el-dialog
      v-model="viewDialogVisible"
      title="记录详情"
      width="800px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="记录类型">
          {{ getRecordTypeLabel(currentRecord?.recordType) }}
        </el-descriptions-item>
        <el-descriptions-item label="记录ID">
          {{ currentRecord?.recordId }}
        </el-descriptions-item>
        <el-descriptions-item label="删除原因">
          {{ currentRecord?.deleteReason }}
        </el-descriptions-item>
        <el-descriptions-item label="操作人">
          {{ currentRecord?.operatorName }}
        </el-descriptions-item>
        <el-descriptions-item label="删除时间">
          {{ formatDateTime(currentRecord?.deletedAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">
          {{ formatDateTime(currentRecord?.createdAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="更新时间">
          {{ formatDateTime(currentRecord?.updatedAt) }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="record-data">
        <h3>删除前数据</h3>
        <pre>{{ formatJson(currentRecord?.recordData) }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { DeletedRecord } from '@/types/trash';
import { RecordType } from '@/types/trash';
import {
  getDeletedRecords,
  getDeletedRecord,
  restoreRecord,
  permanentlyDelete
} from '@/api/trash';
import { formatDateTime } from '@/utils/format';

// 数据列表
const records = ref<DeletedRecord[]>([]);
const loading = ref(false);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);
const recordType = ref<RecordType | ''>('');

// 查看详情
const viewDialogVisible = ref(false);
const currentRecord = ref<DeletedRecord | null>(null);

// 获取记录类型标签
const getRecordTypeLabel = (type: RecordType | undefined) => {
  if (!type) return '';
  const map: Record<RecordType, string> = {
    [RecordType.MEMBER]: '会员',
    [RecordType.PRODUCT]: '商品',
    [RecordType.ACTIVITY]: '活动'
  };
  return map[type] || type;
};

// 格式化JSON数据
const formatJson = (json: string | undefined) => {
  if (!json) return '';
  try {
    return JSON.stringify(JSON.parse(json), null, 2);
  } catch {
    return json;
  }
};

// 获取删除记录列表
const fetchRecords = async () => {
  loading.value = true;
  try {
    const res = await getDeletedRecords({
      type: recordType.value || undefined,
      page: currentPage.value,
      pageSize: pageSize.value
    });
    records.value = res.list;
    total.value = res.total;
  } catch (error) {
    console.error('获取删除记录列表失败:', error);
    ElMessage.error('获取删除记录列表失败');
  } finally {
    loading.value = false;
  }
};

// 查看记录详情
const handleView = async (row: DeletedRecord) => {
  try {
    const res = await getDeletedRecord(row.id);
    currentRecord.value = res.data;
    viewDialogVisible.value = true;
  } catch (error) {
    console.error('获取记录详情失败:', error);
    ElMessage.error('获取记录详情失败');
  }
};

// 恢复记录
const handleRestore = async (row: DeletedRecord) => {
  try {
    await ElMessageBox.confirm(
      '确定要恢复该记录吗？',
      '提示',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );
    await restoreRecord(row.id);
    ElMessage.success('恢复记录成功');
    fetchRecords();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('恢复记录失败:', error);
      ElMessage.error('恢复记录失败');
    }
  }
};

// 永久删除记录
const handleDelete = async (row: DeletedRecord) => {
  try {
    await ElMessageBox.confirm(
      '确定要永久删除该记录吗？此操作不可恢复！',
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    );
    await permanentlyDelete(row.id);
    ElMessage.success('永久删除记录成功');
    fetchRecords();
  } catch (error) {
    if (error !== 'cancel') {
      console.error('永久删除记录失败:', error);
      ElMessage.error('永久删除记录失败');
    }
  }
};

// 分页
const handleSizeChange = (val: number) => {
  pageSize.value = val;
  fetchRecords();
};

const handleCurrentChange = (val: number) => {
  currentPage.value = val;
  fetchRecords();
};

// 监听记录类型变化
watch(recordType, () => {
  currentPage.value = 1;
  fetchRecords();
});

onMounted(() => {
  fetchRecords();
});
</script>

<style scoped>
.trash-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.record-data {
  margin-top: 20px;
}

.record-data h3 {
  margin-bottom: 10px;
}

.record-data pre {
  background-color: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  overflow: auto;
  max-height: 400px;
}
</style> 