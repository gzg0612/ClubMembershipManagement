<template>
  <div class="activity-container">
    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="活动类型">
          <el-select v-model="searchForm.type" placeholder="请选择活动类型" clearable>
            <el-option
              v-for="type in activityTypes"
              :key="type.value"
              :label="type.label"
              :value="type.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="活动状态">
          <el-select v-model="searchForm.status" placeholder="请选择活动状态" clearable>
            <el-option
              v-for="status in activityStatuses"
              :key="status.value"
              :label="status.label"
              :value="status.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="门店">
          <el-select v-model="searchForm.storeId" placeholder="请选择门店" clearable>
            <el-option
              v-for="store in stores"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.timeRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card class="operation-card">
      <el-button type="primary" @click="handleCreate">新建活动</el-button>
    </el-card>

    <!-- 活动列表 -->
    <el-card class="list-card">
      <el-table
        v-loading="loading"
        :data="activityList"
        border
        style="width: 100%"
      >
        <el-table-column prop="name" label="活动名称" min-width="150" />
        <el-table-column prop="type" label="活动类型" width="120">
          <template #default="{ row }">
            {{ getActivityTypeLabel(row.type) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="活动状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getActivityStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="store.name" label="所属门店" width="150" />
        <el-table-column prop="startTime" label="开始时间" width="180" />
        <el-table-column prop="endTime" label="结束时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" link @click="handleView(row)">
                查看
              </el-button>
              <el-button type="primary" link @click="handleEdit(row)">
                编辑
              </el-button>
              <el-button
                type="primary"
                link
                @click="handleUpdateStatus(row)"
              >
                状态
              </el-button>
              <el-button type="danger" link @click="handleDelete(row)">
                删除
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
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

    <!-- 活动表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="dialogType === 'create' ? '新建活动' : '编辑活动'"
      width="60%"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="活动名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入活动名称" />
        </el-form-item>
        <el-form-item label="活动类型" prop="type">
          <el-select v-model="form.type" placeholder="请选择活动类型">
            <el-option
              v-for="type in activityTypes"
              :key="type.value"
              :label="type.label"
              :value="type.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="活动状态" prop="status">
          <el-select v-model="form.status" placeholder="请选择活动状态">
            <el-option
              v-for="status in activityStatuses"
              :key="status.value"
              :label="status.label"
              :value="status.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="所属门店" prop="storeId">
          <el-select v-model="form.storeId" placeholder="请选择门店">
            <el-option
              v-for="store in stores"
              :key="store.id"
              :label="store.name"
              :value="store.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="活动时间" prop="timeRange">
          <el-date-picker
            v-model="form.timeRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            value-format="YYYY-MM-DD HH:mm:ss"
          />
        </el-form-item>
        <el-form-item label="活动描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请输入活动描述"
          />
        </el-form-item>
        <el-form-item label="活动规则">
          <div v-for="(rule, index) in form.rules" :key="index" class="rule-item">
            <el-row :gutter="20">
              <el-col :span="8">
                <el-form-item
                  :prop="'rules.' + index + '.conditionType'"
                  :rules="rules.conditionType"
                >
                  <el-select
                    v-model="rule.conditionType"
                    placeholder="条件类型"
                  >
                    <el-option
                      v-for="type in conditionTypes"
                      :key="type.value"
                      :label="type.label"
                      :value="type.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item
                  :prop="'rules.' + index + '.conditionValue'"
                  :rules="rules.conditionValue"
                >
                  <el-input-number
                    v-model="rule.conditionValue"
                    :min="0"
                    placeholder="条件值"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item
                  :prop="'rules.' + index + '.rewardType'"
                  :rules="rules.rewardType"
                >
                  <el-select
                    v-model="rule.rewardType"
                    placeholder="奖励类型"
                  >
                    <el-option
                      v-for="type in rewardTypes"
                      :key="type.value"
                      :label="type.label"
                      :value="type.value"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item
                  :prop="'rules.' + index + '.rewardValue'"
                  :rules="rules.rewardValue"
                >
                  <el-input-number
                    v-model="rule.rewardValue"
                    :min="0"
                    placeholder="奖励值"
                  />
                </el-form-item>
              </el-col>
              <el-col :span="4">
                <el-button
                  type="danger"
                  link
                  @click="handleRemoveRule(index)"
                >
                  删除
                </el-button>
              </el-col>
            </el-row>
          </div>
          <el-button type="primary" link @click="handleAddRule">
            添加规则
          </el-button>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 活动状态更新对话框 -->
    <el-dialog
      v-model="statusDialogVisible"
      title="更新活动状态"
      width="30%"
    >
      <el-form
        ref="statusFormRef"
        :model="statusForm"
        :rules="statusRules"
        label-width="100px"
      >
        <el-form-item label="活动状态" prop="status">
          <el-select v-model="statusForm.status" placeholder="请选择活动状态">
            <el-option
              v-for="status in activityStatuses"
              :key="status.value"
              :label="status.label"
              :value="status.value"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="statusDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleStatusSubmit">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { FormInstance } from 'element-plus';
import {
  getActivityList,
  createActivity,
  updateActivity,
  deleteActivity,
  updateActivityStatus,
} from '@/api/activity';
import type {
  Activity,
  ActivityType,
  ActivityStatus,
  ActivityRuleConditionType,
  ActivityRuleRewardType,
  CreateActivityRequest,
  UpdateActivityRequest,
  UpdateActivityStatusRequest,
} from '@/types/activity';

// 活动类型选项
const activityTypes = [
  { label: '折扣活动', value: ActivityType.DISCOUNT },
  { label: '优惠券活动', value: ActivityType.COUPON },
  { label: '积分活动', value: ActivityType.POINTS },
  { label: '特价活动', value: ActivityType.SPECIAL },
];

// 活动状态选项
const activityStatuses = [
  { label: '草稿', value: ActivityStatus.DRAFT },
  { label: '待审核', value: ActivityStatus.PENDING },
  { label: '进行中', value: ActivityStatus.ACTIVE },
  { label: '已暂停', value: ActivityStatus.PAUSED },
  { label: '已结束', value: ActivityStatus.ENDED },
  { label: '已拒绝', value: ActivityStatus.REJECTED },
];

// 条件类型选项
const conditionTypes = [
  { label: '消费金额', value: ActivityRuleConditionType.AMOUNT },
  { label: '购买数量', value: ActivityRuleConditionType.QUANTITY },
  { label: '会员等级', value: ActivityRuleConditionType.MEMBER_LEVEL },
  { label: '会员积分', value: ActivityRuleConditionType.MEMBER_POINTS },
];

// 奖励类型选项
const rewardTypes = [
  { label: '折扣', value: ActivityRuleRewardType.DISCOUNT },
  { label: '金额', value: ActivityRuleRewardType.AMOUNT },
  { label: '积分', value: ActivityRuleRewardType.POINTS },
  { label: '优惠券', value: ActivityRuleRewardType.COUPON },
];

// 搜索表单
const searchForm = reactive({
  type: '',
  status: '',
  storeId: '',
  timeRange: [] as string[],
});

// 分页参数
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);

// 加载状态
const loading = ref(false);

// 活动列表
const activityList = ref<Activity[]>([]);

// 表单对话框
const dialogVisible = ref(false);
const dialogType = ref<'create' | 'edit'>('create');
const formRef = ref<FormInstance>();
const form = reactive<CreateActivityRequest>({
  name: '',
  type: ActivityType.DISCOUNT,
  status: ActivityStatus.DRAFT,
  startTime: '',
  endTime: '',
  description: '',
  storeId: 0,
  rules: [],
});

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入活动名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择活动类型', trigger: 'change' }],
  status: [{ required: true, message: '请选择活动状态', trigger: 'change' }],
  storeId: [{ required: true, message: '请选择门店', trigger: 'change' }],
  timeRange: [{ required: true, message: '请选择活动时间', trigger: 'change' }],
  conditionType: [{ required: true, message: '请选择条件类型', trigger: 'change' }],
  conditionValue: [{ required: true, message: '请输入条件值', trigger: 'blur' }],
  rewardType: [{ required: true, message: '请选择奖励类型', trigger: 'change' }],
  rewardValue: [{ required: true, message: '请输入奖励值', trigger: 'blur' }],
};

// 状态更新对话框
const statusDialogVisible = ref(false);
const statusFormRef = ref<FormInstance>();
const statusForm = reactive<UpdateActivityStatusRequest>({
  status: ActivityStatus.DRAFT,
});
const currentActivityId = ref<number>();

// 状态表单验证规则
const statusRules = {
  status: [{ required: true, message: '请选择活动状态', trigger: 'change' }],
};

// 获取活动列表
const fetchActivityList = async () => {
  loading.value = true;
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
      type: searchForm.type,
      status: searchForm.status,
      storeId: searchForm.storeId,
      startTime: searchForm.timeRange[0],
      endTime: searchForm.timeRange[1],
    };
    const res = await getActivityList(params);
    activityList.value = res.data.list;
    total.value = res.data.total;
  } catch (error) {
    console.error('获取活动列表失败:', error);
    ElMessage.error('获取活动列表失败');
  } finally {
    loading.value = false;
  }
};

// 搜索
const handleSearch = () => {
  currentPage.value = 1;
  fetchActivityList();
};

// 重置搜索
const handleReset = () => {
  searchForm.type = '';
  searchForm.status = '';
  searchForm.storeId = '';
  searchForm.timeRange = [];
  handleSearch();
};

// 新建活动
const handleCreate = () => {
  dialogType.value = 'create';
  Object.assign(form, {
    name: '',
    type: ActivityType.DISCOUNT,
    status: ActivityStatus.DRAFT,
    startTime: '',
    endTime: '',
    description: '',
    storeId: 0,
    rules: [],
  });
  dialogVisible.value = true;
};

// 编辑活动
const handleEdit = (row: Activity) => {
  dialogType.value = 'edit';
  Object.assign(form, {
    name: row.name,
    type: row.type,
    status: row.status,
    startTime: row.startTime,
    endTime: row.endTime,
    description: row.description,
    storeId: row.storeId,
    rules: row.rules.map(rule => ({
      conditionType: rule.conditionType,
      conditionValue: rule.conditionValue,
      rewardType: rule.rewardType,
      rewardValue: rule.rewardValue,
    })),
  });
  dialogVisible.value = true;
};

// 更新活动状态
const handleUpdateStatus = (row: Activity) => {
  currentActivityId.value = row.id;
  statusForm.status = row.status;
  statusDialogVisible.value = true;
};

// 删除活动
const handleDelete = (row: Activity) => {
  ElMessageBox.confirm(
    '确定要删除该活动吗？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
    try {
      await deleteActivity(row.id);
      ElMessage.success('删除成功');
      fetchActivityList();
    } catch (error) {
      console.error('删除活动失败:', error);
      ElMessage.error('删除活动失败');
    }
  });
};

// 添加规则
const handleAddRule = () => {
  form.rules.push({
    conditionType: ActivityRuleConditionType.AMOUNT,
    conditionValue: 0,
    rewardType: ActivityRuleRewardType.DISCOUNT,
    rewardValue: 0,
  });
};

// 删除规则
const handleRemoveRule = (index: number) => {
  form.rules.splice(index, 1);
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const [startTime, endTime] = form.timeRange;
        const data = {
          ...form,
          startTime,
          endTime,
        };
        if (dialogType.value === 'create') {
          await createActivity(data);
          ElMessage.success('创建成功');
        } else {
          await updateActivity(currentActivityId.value!, data);
          ElMessage.success('更新成功');
        }
        dialogVisible.value = false;
        fetchActivityList();
      } catch (error) {
        console.error('提交表单失败:', error);
        ElMessage.error('提交失败');
      }
    }
  });
};

// 提交状态更新
const handleStatusSubmit = async () => {
  if (!statusFormRef.value || !currentActivityId.value) return;
  await statusFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        await updateActivityStatus(currentActivityId.value, statusForm);
        ElMessage.success('更新成功');
        statusDialogVisible.value = false;
        fetchActivityList();
      } catch (error) {
        console.error('更新状态失败:', error);
        ElMessage.error('更新失败');
      }
    }
  });
};

// 分页大小改变
const handleSizeChange = (val: number) => {
  pageSize.value = val;
  fetchActivityList();
};

// 页码改变
const handleCurrentChange = (val: number) => {
  currentPage.value = val;
  fetchActivityList();
};

// 获取活动类型标签
const getActivityTypeLabel = (type: ActivityType) => {
  const option = activityTypes.find(opt => opt.value === type);
  return option ? option.label : '';
};

// 获取活动状态标签
const getActivityStatusLabel = (status: ActivityStatus) => {
  const option = activityStatuses.find(opt => opt.value === status);
  return option ? option.label : '';
};

// 获取状态类型
const getStatusType = (status: ActivityStatus) => {
  const statusMap: Record<ActivityStatus, string> = {
    [ActivityStatus.DRAFT]: 'info',
    [ActivityStatus.PENDING]: 'warning',
    [ActivityStatus.ACTIVE]: 'success',
    [ActivityStatus.PAUSED]: 'warning',
    [ActivityStatus.ENDED]: 'info',
    [ActivityStatus.REJECTED]: 'danger',
  };
  return statusMap[status];
};

// 页面加载时获取数据
onMounted(() => {
  fetchActivityList();
});
</script>

<style scoped>
.activity-container {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.search-form {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.operation-card {
  margin-bottom: 20px;
}

.list-card {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.rule-item {
  margin-bottom: 10px;
  padding: 10px;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 