<template>
  <div class="activity-edit-container">
    <el-card class="edit-card">
      <template #header>
        <div class="card-header">
          <span>编辑活动</span>
          <el-button @click="handleCancel">取消</el-button>
        </div>
      </template>

      <ActivityForm
        ref="formRef"
        :form="form"
        @submit="handleSubmit"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import type { FormInstance } from 'element-plus';
import { getActivity, updateActivity } from '@/api/activity';
import type { Activity, ActivityType, ActivityStatus, UpdateActivityRequest } from '@/types/activity';
import ActivityForm from '@/components/activity/ActivityForm.vue';

const route = useRoute();
const router = useRouter();

// 表单引用
const formRef = ref<FormInstance>();

// 表单数据
const form = ref<{
  name: string;
  type: ActivityType;
  status: ActivityStatus;
  startTime: string;
  endTime: string;
  description: string;
  storeId: number;
  rules: Array<{
    conditionType: string;
    conditionValue: number;
    rewardType: string;
    rewardValue: number;
  }>;
  timeRange: [string, string];
}>({
  name: '',
  type: ActivityType.DISCOUNT,
  status: ActivityStatus.DRAFT,
  startTime: '',
  endTime: '',
  description: '',
  storeId: 0,
  rules: [],
  timeRange: ['', ''],
});

// 获取活动详情
const fetchActivityDetail = async () => {
  try {
    const id = Number(route.params.id);
    const res = await getActivity(id);
    const activity = res.data;
    form.value = {
      name: activity.name,
      type: activity.type,
      status: activity.status,
      startTime: activity.startTime,
      endTime: activity.endTime,
      description: activity.description,
      storeId: activity.storeId,
      rules: activity.rules.map((rule: {
        conditionType: string;
        conditionValue: number;
        rewardType: string;
        rewardValue: number;
      }) => ({
        conditionType: rule.conditionType,
        conditionValue: rule.conditionValue,
        rewardType: rule.rewardType,
        rewardValue: rule.rewardValue,
      })),
      timeRange: [activity.startTime, activity.endTime],
    };
  } catch (error) {
    console.error('获取活动详情失败:', error);
    ElMessage.error('获取活动详情失败');
  }
};

// 取消编辑
const handleCancel = () => {
  router.back();
};

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return;
  await formRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const id = Number(route.params.id);
        const [startTime, endTime] = form.value.timeRange;
        const data: UpdateActivityRequest = {
          name: form.value.name,
          type: form.value.type,
          status: form.value.status,
          startTime,
          endTime,
          description: form.value.description,
          storeId: form.value.storeId,
          rules: form.value.rules,
        };
        await updateActivity(id, data);
        ElMessage.success('更新成功');
        router.push(`/activity/detail/${id}`);
      } catch (error) {
        console.error('更新活动失败:', error);
        ElMessage.error('更新失败');
      }
    }
  });
};

// 页面加载时获取数据
onMounted(() => {
  fetchActivityDetail();
});
</script>

<style scoped>
.activity-edit-container {
  padding: 20px;
}

.edit-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style> 