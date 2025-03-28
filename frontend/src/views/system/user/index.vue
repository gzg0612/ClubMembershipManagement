<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.keyword"
        placeholder="用户名/姓名/手机号"
        style="width: 200px;"
        class="filter-item"
        @keyup.enter.native="handleFilter"
      />
      <el-select
        v-model="listQuery.status"
        placeholder="状态"
        clearable
        style="width: 120px"
        class="filter-item"
      >
        <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
      <el-select
        v-model="listQuery.role_id"
        placeholder="角色"
        clearable
        style="width: 120px"
        class="filter-item"
      >
        <el-option v-for="item in roleOptions" :key="item.id" :label="item.name" :value="item.id" />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        搜索
      </el-button>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-plus" @click="handleCreate">
        新增
      </el-button>
    </div>

    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="加载中..."
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="ID" align="center" width="80">
        <template slot-scope="{row}">
          <span>{{ row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="用户名" min-width="120px">
        <template slot-scope="{row}">
          <span>{{ row.username }}</span>
        </template>
      </el-table-column>
      <el-table-column label="姓名" min-width="120px">
        <template slot-scope="{row}">
          <span>{{ row.real_name }}</span>
        </template>
      </el-table-column>
      <el-table-column label="头像" width="100" align="center">
        <template slot-scope="{row}">
          <el-avatar
            v-if="row.avatar"
            size="small"
            :src="row.avatar"
          />
          <el-avatar
            v-else
            size="small"
            icon="el-icon-user-solid"
          />
        </template>
      </el-table-column>
      <el-table-column label="角色" min-width="150px">
        <template slot-scope="{row}">
          <el-tag
            v-for="(role, index) in row.roles"
            :key="role.id"
            size="mini"
            :type="index % 5 === 0 ? '' : ['success', 'info', 'warning', 'danger'][index % 4]"
            class="role-tag"
          >
            {{ role.name }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="手机号" min-width="120px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.phone || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100" align="center">
        <template slot-scope="{row}">
          <el-tag :type="row.status === 1 ? 'success' : row.status === 0 ? 'info' : 'danger'">
            {{ row.status === 1 ? '正常' : row.status === 0 ? '禁用' : '锁定' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最后登录时间" width="160" align="center">
        <template slot-scope="{row}">
          <span>{{ row.last_login_time || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="250" class-name="small-padding fixed-width">
        <template slot-scope="{row}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button
            v-if="row.status === 0"
            size="mini"
            type="success"
            @click="handleModifyStatus(row, 1)"
          >
            启用
          </el-button>
          <el-button
            v-else-if="row.status === 1"
            size="mini"
            type="info"
            @click="handleModifyStatus(row, 0)"
          >
            禁用
          </el-button>
          <el-button
            v-if="row.status === 2"
            size="mini"
            type="success"
            @click="handleModifyStatus(row, 1)"
          >
            解锁
          </el-button>
          <el-button size="mini" type="danger" @click="handleDelete(row)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      v-show="total > 0"
      :total="total"
      :page.sync="listQuery.page"
      :limit.sync="listQuery.page_size"
      @pagination="getList"
    />

    <!-- 用户编辑对话框 -->
    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="80px"
        style="width: 400px; margin-left:50px;"
      >
        <el-form-item label="用户名" prop="username">
          <el-input v-model="temp.username" :disabled="dialogStatus === 'update'" />
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="dialogStatus === 'create'">
          <el-input v-model="temp.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="姓名" prop="real_name">
          <el-input v-model="temp.real_name" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="temp.phone" />
        </el-form-item>
        <el-form-item label="邮箱">
          <el-input v-model="temp.email" />
        </el-form-item>
        <el-form-item label="角色" prop="role_ids">
          <el-select
            v-model="temp.role_ids"
            multiple
            placeholder="请选择"
            style="width: 100%"
          >
            <el-option
              v-for="item in roleOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="temp.status" class="filter-item" placeholder="请选择">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="dialogStatus === 'create' ? createData() : updateData()">
          确认
        </el-button>
      </div>
    </el-dialog>

    <!-- 重置密码对话框 -->
    <el-dialog title="重置密码" :visible.sync="passwordDialogVisible" width="30%">
      <el-form
        ref="passwordForm"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="新密码" prop="new_password">
          <el-input v-model="passwordForm.new_password" type="password" show-password />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="resetPassword">确定</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getUserList, createUser, updateUser, deleteUser, resetUserPassword } from '@/api/user'
import Pagination from '@/components/Pagination' // 分页组件

export default {
  name: 'UserList',
  components: { Pagination },
  data() {
    return {
      list: null,
      total: 0,
      listLoading: true,
      listQuery: {
        page: 1,
        page_size: 20,
        keyword: '',
        status: undefined,
        role_id: undefined
      },
      statusOptions: [
        { label: '正常', value: 1 },
        { label: '禁用', value: 0 },
        { label: '锁定', value: 2 }
      ],
      roleOptions: [], // 从后端获取角色列表
      temp: {
        id: undefined,
        username: '',
        password: '',
        real_name: '',
        phone: '',
        email: '',
        role_ids: [],
        status: 1
      },
      passwordForm: {
        new_password: '',
        user_id: null
      },
      passwordRules: {
        new_password: [
          { required: true, message: '请输入新密码', trigger: 'blur' },
          { min: 6, message: '密码长度不少于6位', trigger: 'blur' }
        ]
      },
      dialogFormVisible: false,
      passwordDialogVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑用户',
        create: '创建用户'
      },
      rules: {
        username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, message: '密码长度不少于6位', trigger: 'blur' }
        ],
        real_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
        role_ids: [{ required: true, message: '请选择角色', trigger: 'change' }]
      }
    }
  },
  created() {
    this.getList()
    // 获取角色下拉列表数据
    // 这里应该调用获取角色列表的API，为简化先使用模拟数据
    this.roleOptions = [
      { id: 1, name: '管理员' },
      { id: 2, name: '普通用户' }
    ]
  },
  methods: {
    getList() {
      this.listLoading = true
      getUserList(this.listQuery).then(response => {
        this.list = response.data.items
        this.total = response.data.total
        this.listLoading = false
      })
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        username: '',
        password: '',
        real_name: '',
        phone: '',
        email: '',
        role_ids: [],
        status: 1
      }
    },
    handleCreate() {
      this.resetTemp()
      this.dialogStatus = 'create'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    createData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          createUser(this.temp).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '创建用户成功',
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row)
      this.temp.role_ids = row.roles.map(role => role.id)
      this.dialogStatus = 'update'
      this.dialogFormVisible = true
      this.$nextTick(() => {
        this.$refs['dataForm'].clearValidate()
      })
    },
    updateData() {
      this.$refs['dataForm'].validate((valid) => {
        if (valid) {
          const tempData = Object.assign({}, this.temp)
          updateUser(tempData.id, tempData).then(() => {
            this.dialogFormVisible = false
            this.$notify({
              title: '成功',
              message: '更新用户成功',
              type: 'success',
              duration: 2000
            })
            this.getList()
          })
        }
      })
    },
    handleDelete(row) {
      this.$confirm('此操作将永久删除该用户, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        deleteUser(row.id).then(() => {
          this.$notify({
            title: '成功',
            message: '删除用户成功',
            type: 'success',
            duration: 2000
          })
          this.getList()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    handleModifyStatus(row, status) {
      const statusText = status === 1 ? '启用' : status === 0 ? '禁用' : '锁定'
      this.$confirm(`确认要${statusText}该用户吗?`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        updateUser(row.id, { status }).then(() => {
          this.$notify({
            title: '成功',
            message: `用户${statusText}成功`,
            type: 'success',
            duration: 2000
          })
          this.getList()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消操作'
        })
      })
    },
    showResetPassword(row) {
      this.passwordForm.user_id = row.id
      this.passwordForm.new_password = ''
      this.passwordDialogVisible = true
      this.$nextTick(() => {
        this.$refs['passwordForm'].clearValidate()
      })
    },
    resetPassword() {
      this.$refs['passwordForm'].validate((valid) => {
        if (valid) {
          resetUserPassword(this.passwordForm.user_id, {
            new_password: this.passwordForm.new_password
          }).then(() => {
            this.passwordDialogVisible = false
            this.$notify({
              title: '成功',
              message: '重置密码成功',
              type: 'success',
              duration: 2000
            })
          })
        }
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.filter-container {
  padding-bottom: 10px;
  
  .filter-item {
    margin-right: 10px;
  }
}

.role-tag {
  margin-right: 5px;
}
</style> 