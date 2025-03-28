<template>
  <div class="settings-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统设置</span>
          <el-button type="primary" @click="handleSave">保存设置</el-button>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <!-- 基本设置 -->
        <el-tab-pane label="基本设置" name="basic">
          <el-form ref="basicFormRef" :model="settings" label-width="120px">
            <el-form-item label="系统名称" prop="siteName">
              <el-input v-model="settings.siteName" placeholder="请输入系统名称" />
            </el-form-item>
            <el-form-item label="系统Logo" prop="siteLogo">
              <el-upload
                class="avatar-uploader"
                action="/api/v1/upload"
                :show-file-list="false"
                :on-success="handleLogoSuccess"
                :before-upload="beforeLogoUpload"
              >
                <img v-if="settings.siteLogo" :src="settings.siteLogo" class="avatar" />
                <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
              </el-upload>
            </el-form-item>
            <el-form-item label="系统描述" prop="siteDescription">
              <el-input
                v-model="settings.siteDescription"
                type="textarea"
                placeholder="请输入系统描述"
              />
            </el-form-item>
            <el-form-item label="联系电话" prop="contactPhone">
              <el-input v-model="settings.contactPhone" placeholder="请输入联系电话" />
            </el-form-item>
            <el-form-item label="联系邮箱" prop="contactEmail">
              <el-input v-model="settings.contactEmail" placeholder="请输入联系邮箱" />
            </el-form-item>
            <el-form-item label="地址" prop="address">
              <el-input v-model="settings.address" placeholder="请输入地址" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 会员设置 -->
        <el-tab-pane label="会员设置" name="member">
          <el-form ref="memberFormRef" :model="settings" label-width="120px">
            <el-form-item label="默认会员等级" prop="defaultMemberLevel">
              <el-input-number v-model="settings.defaultMemberLevel" :min="1" :max="10" />
            </el-form-item>
            <el-form-item label="自动升级" prop="enableAutoUpgrade">
              <el-switch v-model="settings.enableAutoUpgrade" />
            </el-form-item>
            <el-form-item label="积分系统" prop="enablePoints">
              <el-switch v-model="settings.enablePoints" />
            </el-form-item>
            <el-form-item
              v-if="settings.enablePoints"
              label="积分比例"
              prop="pointsPerYuan"
            >
              <el-input-number
                v-model="settings.pointsPerYuan"
                :min="0"
                :max="100"
                :precision="2"
              />
              <span class="ml-2">积分/元</span>
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 通知设置 -->
        <el-tab-pane label="通知设置" name="notification">
          <el-form ref="notificationFormRef" :model="settings" label-width="120px">
            <el-form-item label="短信通知" prop="enableSMS">
              <el-switch v-model="settings.enableSMS" />
            </el-form-item>
            <el-form-item label="邮件通知" prop="enableEmail">
              <el-switch v-model="settings.enableEmail" />
            </el-form-item>
            <el-form-item label="微信通知" prop="enableWechat">
              <el-switch v-model="settings.enableWechat" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 安全设置 -->
        <el-tab-pane label="安全设置" name="security">
          <el-form ref="securityFormRef" :model="settings" label-width="120px">
            <el-form-item label="密码最小长度" prop="passwordMinLength">
              <el-input-number v-model="settings.passwordMinLength" :min="6" :max="20" />
            </el-form-item>
            <el-form-item label="特殊字符" prop="requireSpecialChar">
              <el-switch v-model="settings.requireSpecialChar" />
            </el-form-item>
            <el-form-item label="数字" prop="requireNumber">
              <el-switch v-model="settings.requireNumber" />
            </el-form-item>
            <el-form-item label="大写字母" prop="requireUppercase">
              <el-switch v-model="settings.requireUppercase" />
            </el-form-item>
            <el-form-item label="登录尝试次数" prop="loginAttempts">
              <el-input-number v-model="settings.loginAttempts" :min="1" :max="10" />
            </el-form-item>
            <el-form-item label="锁定时间(分钟)" prop="lockoutDuration">
              <el-input-number v-model="settings.lockoutDuration" :min="1" :max="1440" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 备份设置 -->
        <el-tab-pane label="备份设置" name="backup">
          <el-form ref="backupFormRef" :model="settings" label-width="120px">
            <el-form-item label="自动备份" prop="autoBackup">
              <el-switch v-model="settings.autoBackup" />
            </el-form-item>
            <el-form-item v-if="settings.autoBackup" label="备份频率" prop="backupFrequency">
              <el-select v-model="settings.backupFrequency">
                <el-option label="每天" value="daily" />
                <el-option label="每周" value="weekly" />
                <el-option label="每月" value="monthly" />
              </el-select>
            </el-form-item>
            <el-form-item v-if="settings.autoBackup" label="备份时间" prop="backupTime">
              <el-time-picker v-model="settings.backupTime" format="HH:mm" />
            </el-form-item>
            <el-form-item v-if="settings.autoBackup" label="保留天数" prop="backupRetention">
              <el-input-number v-model="settings.backupRetention" :min="1" :max="365" />
            </el-form-item>
          </el-form>
        </el-tab-pane>

        <!-- 其他设置 -->
        <el-tab-pane label="其他设置" name="other">
          <el-form ref="otherFormRef" :model="settings" label-width="120px">
            <el-form-item label="时区" prop="timezone">
              <el-select v-model="settings.timezone">
                <el-option label="(GMT+08:00) 北京" value="Asia/Shanghai" />
                <el-option label="(GMT+00:00) 伦敦" value="Europe/London" />
                <el-option label="(GMT-05:00) 纽约" value="America/New_York" />
              </el-select>
            </el-form-item>
            <el-form-item label="日期格式" prop="dateFormat">
              <el-select v-model="settings.dateFormat">
                <el-option label="YYYY-MM-DD" value="YYYY-MM-DD" />
                <el-option label="MM/DD/YYYY" value="MM/DD/YYYY" />
                <el-option label="DD/MM/YYYY" value="DD/MM/YYYY" />
              </el-select>
            </el-form-item>
            <el-form-item label="时间格式" prop="timeFormat">
              <el-select v-model="settings.timeFormat">
                <el-option label="24小时制" value="24" />
                <el-option label="12小时制" value="12" />
              </el-select>
            </el-form-item>
            <el-form-item label="货币单位" prop="currency">
              <el-select v-model="settings.currency">
                <el-option label="人民币" value="CNY" />
                <el-option label="美元" value="USD" />
                <el-option label="欧元" value="EUR" />
              </el-select>
            </el-form-item>
            <el-form-item label="系统语言" prop="language">
              <el-select v-model="settings.language">
                <el-option label="简体中文" value="zh-CN" />
                <el-option label="English" value="en-US" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { ElMessage, ElMessageBox } from 'element-plus';
import { Plus } from '@element-plus/icons-vue';
import type { SystemSettings } from '@/types/settings';
import { getSettings, updateSettings } from '@/api/settings';

// 当前激活的标签页
const activeTab = ref('basic');

// 系统设置数据
const settings = ref<SystemSettings>({
  // 基本设置
  siteName: '',
  siteLogo: '',
  siteDescription: '',
  contactPhone: '',
  contactEmail: '',
  address: '',

  // 会员设置
  defaultMemberLevel: 1,
  enableAutoUpgrade: false,
  pointsPerYuan: 1,
  enablePoints: true,

  // 通知设置
  enableSMS: false,
  enableEmail: false,
  enableWechat: false,

  // 安全设置
  passwordMinLength: 8,
  requireSpecialChar: true,
  requireNumber: true,
  requireUppercase: true,
  loginAttempts: 5,
  lockoutDuration: 30,

  // 备份设置
  autoBackup: false,
  backupFrequency: 'daily',
  backupTime: '00:00',
  backupRetention: 30,

  // 其他设置
  timezone: 'Asia/Shanghai',
  dateFormat: 'YYYY-MM-DD',
  timeFormat: '24',
  currency: 'CNY',
  language: 'zh-CN'
});

// 获取系统设置
const fetchSettings = async () => {
  try {
    const res = await getSettings();
    if (res.code === 200) {
      settings.value = res.data;
    } else {
      ElMessage.error(res.message || '获取系统设置失败');
    }
  } catch (error) {
    console.error('获取系统设置失败:', error);
    ElMessage.error('获取系统设置失败');
  }
};

// 保存系统设置
const handleSave = async () => {
  try {
    await ElMessageBox.confirm('确定要保存系统设置吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    });
    const res = await updateSettings({ settings: settings.value });
    if (res.code === 200) {
      ElMessage.success('保存系统设置成功');
      fetchSettings();
    } else {
      ElMessage.error(res.message || '保存系统设置失败');
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('保存系统设置失败:', error);
      ElMessage.error('保存系统设置失败');
    }
  }
};

// Logo上传相关
const handleLogoSuccess = (response: any) => {
  if (response.code === 200) {
    settings.value.siteLogo = response.data.url;
    ElMessage.success('上传Logo成功');
  } else {
    ElMessage.error(response.message || '上传Logo失败');
  }
};

const beforeLogoUpload = (file: File) => {
  const isImage = file.type.startsWith('image/');
  const isLt2M = file.size / 1024 / 1024 < 2;

  if (!isImage) {
    ElMessage.error('只能上传图片文件！');
    return false;
  }
  if (!isLt2M) {
    ElMessage.error('图片大小不能超过 2MB！');
    return false;
  }
  return true;
};

onMounted(() => {
  fetchSettings();
});
</script>

<style scoped>
.settings-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.avatar-uploader {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  width: 178px;
  height: 178px;
}

.avatar-uploader:hover {
  border-color: #409eff;
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
  line-height: 178px;
}

.avatar {
  width: 178px;
  height: 178px;
  display: block;
}
</style> 