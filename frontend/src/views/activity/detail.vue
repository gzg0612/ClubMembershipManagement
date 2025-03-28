<template>
  <div class="activity-detail-container">
    <el-card class="detail-card">
      <template #header>
        <div class="card-header">
          <span>活动详情</span>
          <el-button-group>
            <el-button type="primary" @click="handleEdit">编辑</el-button>
            <el-button type="primary" @click="handleUpdateStatus">更新状态</el-button>
            <el-button type="danger" @click="handleDelete">删除</el-button>
          </el-button-group>
        </div>
      </template>

      <el-descriptions :column="2" border>
        <el-descriptions-item label="活动名称">
          {{ activity.name }}
        </el-descriptions-item>
        <el-descriptions-item label="活动类型">
          {{ getActivityTypeLabel(activity.type) }}
        </el-descriptions-item>
        <el-descriptions-item label="活动状态">
          <el-tag :type="getStatusType(activity.status)">
            {{ getActivityStatusLabel(activity.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="所属门店">
          {{ activity.store?.name }}
        </el-descriptions-item>
        <el-descriptions-item label="开始时间">
          {{ activity.startTime }}
        </el-descriptions-item>
        <el-descriptions-item label="结束时间">
          {{ activity.endTime }}
        </el-descriptions-item>
        <el-descriptions-item label="活动描述" :span="2">
          {{ activity.description }}
        </el-descriptions-item>
      </el-descriptions>

      <!-- 活动规则 -->
      <div class="rules-section">
        <h3>活动规则</h3>
        <el-table :data="activity.rules" border style="width: 100%">
          <el-table-column prop="conditionType" label="条件类型" width="150">
            <template #default="{ row }">
              {{ getConditionTypeLabel(row.conditionType) }}
            </template>
          </el-table-column>
          <el-table-column prop="conditionValue" label="条件值" width="150" />
          <el-table-column prop="rewardType" label="奖励类型" width="150">
            <template #default="{ row }">
              {{ getRewardTypeLabel(row.rewardType) }}
            </template>
          </el-table-column>
          <el-table-column prop="rewardValue" label="奖励值" width="150" />
        </el-table>
      </div>

      <!-- 参与记录 -->
      <div class="participants-section">
        <h3>参与记录</h3>
        <el-table
          v-loading="loading"
          :data="participants"
          border
          style="width: 100%"
        >
          <el-table-column prop="member.name" label="会员姓名" width="150" />
          <el-table-column prop="member.phone" label="手机号码" width="150" />
          <el-table-column prop="participateTime" label="参与时间" width="180" />
          <el-table-column prop="amount" label="消费金额" width="150" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'success' ? 'success' : 'danger'">
                {{ row.status === 'success' ? '成功' : '失败' }}
              </el-tag>
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
      </div>
    </el-card>

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
import { useRoute, useRouter } from 'vue-router';
import { ElMessage, ElMessageBox } from 'element-plus';
import type { FormInstance } from 'element-plus';
import {
  getActivity,
  updateActivity,
  deleteActivity,
  updateActivityStatus,
  getActivityParticipants,
} from '@/api/activity';
import type {
  Activity,
  ActivityType,
  ActivityStatus,
  ActivityRuleConditionType,
  ActivityRuleRewardType,
  ActivityParticipant,
  UpdateActivityRequest,
  UpdateActivityStatusRequest,
} from '@/types/activity';

const route = useRoute();
const router = useRouter();

// 活动详情
const activity = ref<Activity>({} as Activity);

// 参与记录
const loading = ref(false);
const participants = ref<ActivityParticipant[]>([]);
const currentPage = ref(1);
const pageSize = ref(10);
const total = ref(0);

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

// 状态更新对话框
const statusDialogVisible = ref(false);
const statusFormRef = ref<FormInstance>();
const statusForm = reactive<UpdateActivityStatusRequest>({
  status: ActivityStatus.DRAFT,
});

// 状态表单验证规则
const statusRules = {
  status: [{ required: true, message: '请选择活动状态', trigger: 'change' }],
};

// 获取活动详情
const fetchActivityDetail = async () => {
  try {
    const id = Number(route.params.id);
    const res = await getActivity(id);
    activity.value = res.data;
  } catch (error) {
    console.error('获取活动详情失败:', error);
    ElMessage.error('获取活动详情失败');
  }
};

// 获取参与记录
const fetchParticipants = async () => {
  loading.value = true;
  try {
    const id = Number(route.params.id);
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
    };
    const res = await getActivityParticipants(id, params);
    participants.value = res.data.list;
    total.value = res.data.total;
  } catch (error) {
    console.error('获取参与记录失败:', error);
    ElMessage.error('获取参与记录失败');
  } finally {
    loading.value = false;
  }
};

// 编辑活动
const handleEdit = () => {
  router.push(`/activity/edit/${route.params.id}`);
};

// 更新活动状态
const handleUpdateStatus = () => {
  statusForm.status = activity.value.status;
  statusDialogVisible.value = true;
};

// 删除活动
const handleDelete = () => {
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
      const id = Number(route.params.id);
      await deleteActivity(id);
      ElMessage.success('删除成功');
      router.push('/activity');
    } catch (error) {
      console.error('删除活动失败:', error);
      ElMessage.error('删除活动失败');
    }
  });
};

// 提交状态更新
const handleStatusSubmit = async () => {
  if (!statusFormRef.value) return;
  await statusFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const id = Number(route.params.id);
        await updateActivityStatus(id, statusForm);
        ElMessage.success('更新成功');
        statusDialogVisible.value = false;
        fetchActivityDetail();
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
  fetchParticipants();
};

// 页码改变
const handleCurrentChange = (val: number) => {
  currentPage.value = val;
  fetchParticipants();
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

// 获取条件类型标签
const getConditionTypeLabel = (type: ActivityRuleConditionType) => {
  const option = conditionTypes.find(opt => opt.value === type);
  return option ? option.label : '';
};

// 获取奖励类型标签
const getRewardTypeLabel = (type: ActivityRuleRewardType) => {
  const option = rewardTypes.find(opt => opt.value === type);
  return option ? option.label : '';
};

// 页面加载时获取数据
onMounted(() => {
  fetchActivityDetail();
  fetchParticipants();
});
</script>

<style scoped>
.activity-detail-container {
  padding: 20px;
}

.detail-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.rules-section,
.participants-section {
  margin-top: 20px;
}

.rules-section h3,
.participants-section h3 {
  margin-bottom: 16px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 