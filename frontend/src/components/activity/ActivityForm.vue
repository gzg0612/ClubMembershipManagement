<template>
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
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import type { FormInstance } from 'element-plus';

// 活动类型
const ActivityType = {
  DISCOUNT: 'DISCOUNT',
  COUPON: 'COUPON',
  POINTS: 'POINTS',
  SPECIAL: 'SPECIAL',
} as const;

// 活动状态
const ActivityStatus = {
  DRAFT: 'DRAFT',
  PENDING: 'PENDING',
  ACTIVE: 'ACTIVE',
  PAUSED: 'PAUSED',
  ENDED: 'ENDED',
  REJECTED: 'REJECTED',
} as const;

// 条件类型
const ActivityRuleConditionType = {
  AMOUNT: 'AMOUNT',
  QUANTITY: 'QUANTITY',
  MEMBER_LEVEL: 'MEMBER_LEVEL',
  MEMBER_POINTS: 'MEMBER_POINTS',
} as const;

// 奖励类型
const ActivityRuleRewardType = {
  DISCOUNT: 'DISCOUNT',
  AMOUNT: 'AMOUNT',
  POINTS: 'POINTS',
  COUPON: 'COUPON',
} as const;

// 表单引用
const formRef = ref<FormInstance>();

// 表单数据
const form = reactive({
  name: '',
  type: ActivityType.DISCOUNT,
  status: ActivityStatus.DRAFT,
  startTime: '',
  endTime: '',
  description: '',
  storeId: 0,
  rules: [] as Array<{
    conditionType: typeof ActivityRuleConditionType[keyof typeof ActivityRuleConditionType];
    conditionValue: number;
    rewardType: typeof ActivityRuleRewardType[keyof typeof ActivityRuleRewardType];
    rewardValue: number;
  }>,
  timeRange: ['', ''] as [string, string],
});

// 门店列表
const stores = ref<Array<{ id: number; name: string }>>([]);

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

// 暴露方法
defineExpose({
  formRef,
  form,
  validate: () => formRef.value?.validate(),
});
</script>

<style scoped>
.rule-item {
  margin-bottom: 10px;
  padding: 10px;
  border: 1px dashed #dcdfe6;
  border-radius: 4px;
}
</style> 