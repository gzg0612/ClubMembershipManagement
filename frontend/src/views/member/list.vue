<template>
  <div class="app-container">
    <div class="filter-container">
      <el-input
        v-model="listQuery.keyword"
        placeholder="会员名/手机号"
        style="width: 200px;"
        class="filter-item"
        clearable
        @keyup.enter.native="handleFilter"
      />
      <el-select
        v-model="listQuery.level"
        placeholder="会员等级"
        clearable
        class="filter-item"
        style="width: 130px"
      >
        <el-option
          v-for="item in levelOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        />
      </el-select>
      <el-select
        v-model="listQuery.status"
        placeholder="状态"
        clearable
        class="filter-item"
        style="width: 130px"
      >
        <el-option
          v-for="item in statusOptions"
          :key="item.value"
          :label="item.label"
          :value="item.value"
        />
      </el-select>
      <el-button v-waves class="filter-item" type="primary" icon="el-icon-search" @click="handleFilter">
        搜索
      </el-button>
      <el-button class="filter-item" style="margin-left: 10px;" type="primary" icon="el-icon-plus" @click="handleCreate">
        添加会员
      </el-button>
      <el-button
        :loading="downloadLoading"
        class="filter-item"
        type="primary"
        icon="el-icon-download"
        @click="handleDownload"
      >
        导出
      </el-button>
    </div>

    <el-table
      v-loading="listLoading"
      :data="list"
      border
      fit
      highlight-current-row
      style="width: 100%;"
    >
      <el-table-column label="ID" prop="id" align="center" width="80">
        <template slot-scope="scope">
          <span>{{ scope.row.id }}</span>
        </template>
      </el-table-column>
      <el-table-column label="会员名" min-width="120px">
        <template slot-scope="{row}">
          <router-link :to="'/member/detail/'+row.id" class="link-type">
            <span>{{ row.name }}</span>
          </router-link>
        </template>
      </el-table-column>
      <el-table-column label="头像" width="100px" align="center">
        <template slot-scope="{row}">
          <el-avatar size="medium" :src="row.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'"></el-avatar>
        </template>
      </el-table-column>
      <el-table-column label="手机号" min-width="120px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.phone }}</span>
        </template>
      </el-table-column>
      <el-table-column label="会员等级" width="100px" align="center">
        <template slot-scope="{row}">
          <el-tag :type="row.level | levelFilter">{{ row.level }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="积分" width="80px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.points }}</span>
        </template>
      </el-table-column>
      <el-table-column label="消费金额" width="100px" align="center">
        <template slot-scope="{row}">
          <span>¥{{ row.amount }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" class-name="status-col" width="100">
        <template slot-scope="{row}">
          <el-tag :type="row.status | statusFilter">
            {{ row.status === 1 ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="注册时间" width="150px" align="center">
        <template slot-scope="{row}">
          <span>{{ row.createTime }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" align="center" width="200" class-name="small-padding fixed-width">
        <template slot-scope="{row,$index}">
          <el-button type="primary" size="mini" @click="handleUpdate(row)">
            编辑
          </el-button>
          <el-button size="mini" type="danger" @click="handleDelete(row,$index)">
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <pagination
      v-show="total>0"
      :total="total"
      :page.sync="listQuery.page"
      :limit.sync="listQuery.limit"
      @pagination="getList"
    />

    <el-dialog :title="textMap[dialogStatus]" :visible.sync="dialogFormVisible">
      <el-form
        ref="dataForm"
        :rules="rules"
        :model="temp"
        label-position="left"
        label-width="100px"
        style="width: 400px; margin-left:50px;"
      >
        <el-form-item label="会员名" prop="name">
          <el-input v-model="temp.name" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="temp.phone" />
        </el-form-item>
        <el-form-item label="会员等级" prop="level">
          <el-select v-model="temp.level" class="filter-item" placeholder="请选择">
            <el-option
              v-for="item in levelOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="积分">
          <el-input-number v-model="temp.points" :min="0" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="temp.status" class="filter-item" placeholder="请选择">
            <el-option
              v-for="item in statusOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="temp.remark"
            :autosize="{ minRows: 2, maxRows: 4}"
            type="textarea"
            placeholder="请输入内容"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogFormVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="dialogStatus==='create'?createData():updateData()">
          确认
        </el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import Pagination from '@/components/Pagination'

export default {
  name: 'MemberList',
  components: { Pagination },
  filters: {
    statusFilter(status) {
      const statusMap = {
        1: 'success',
        0: 'info'
      }
      return statusMap[status]
    },
    levelFilter(level) {
      const levelMap = {
        '普通会员': '',
        '白银会员': 'info',
        '黄金会员': 'warning',
        '铂金会员': 'success',
        '钻石会员': 'danger'
      }
      return levelMap[level]
    }
  },
  data() {
    return {
      list: [
        { id: 1001, name: '张三', phone: '13800138001', level: '黄金会员', points: 2500, amount: 8500, status: 1, createTime: '2025-03-15 10:15:20' },
        { id: 1002, name: '李四', phone: '13800138002', level: '白银会员', points: 1200, amount: 3500, status: 1, createTime: '2025-03-16 13:22:35' },
        { id: 1003, name: '王五', phone: '13800138003', level: '铂金会员', points: 5600, amount: 15200, status: 1, createTime: '2025-03-17 09:45:18' },
        { id: 1004, name: '赵六', phone: '13800138004', level: '普通会员', points: 500, amount: 1200, status: 0, createTime: '2025-03-18 16:33:42' },
        { id: 1005, name: '钱七', phone: '13800138005', level: '黄金会员', points: 3200, amount: 9800, status: 1, createTime: '2025-03-19 11:28:53' },
        { id: 1006, name: '孙八', phone: '13800138006', level: '钻石会员', points: 12000, amount: 32000, status: 1, createTime: '2025-03-20 14:15:30' }
      ],
      total: 6,
      listLoading: false,
      listQuery: {
        page: 1,
        limit: 10,
        keyword: '',
        level: '',
        status: ''
      },
      levelOptions: [
        { value: '普通会员', label: '普通会员' },
        { value: '白银会员', label: '白银会员' },
        { value: '黄金会员', label: '黄金会员' },
        { value: '铂金会员', label: '铂金会员' },
        { value: '钻石会员', label: '钻石会员' }
      ],
      statusOptions: [
        { value: 1, label: '正常' },
        { value: 0, label: '禁用' }
      ],
      temp: {
        id: undefined,
        name: '',
        phone: '',
        level: '',
        points: 0,
        status: 1,
        remark: ''
      },
      dialogFormVisible: false,
      dialogStatus: '',
      textMap: {
        update: '编辑会员',
        create: '添加会员'
      },
      rules: {
        name: [{ required: true, message: '会员名不能为空', trigger: 'blur' }],
        phone: [
          { required: true, message: '手机号不能为空', trigger: 'blur' },
          { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
        ],
        level: [{ required: true, message: '请选择会员等级', trigger: 'change' }]
      },
      downloadLoading: false
    }
  },
  created() {
    this.getList()
  },
  methods: {
    getList() {
      // 实际项目中这里应该调用API获取数据
      this.listLoading = true
      setTimeout(() => {
        this.listLoading = false
      }, 500)
    },
    handleFilter() {
      this.listQuery.page = 1
      this.getList()
    },
    resetTemp() {
      this.temp = {
        id: undefined,
        name: '',
        phone: '',
        level: '',
        points: 0,
        status: 1,
        remark: ''
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
          // 实际项目中这里应该调用API创建数据
          this.temp.id = parseInt(Math.random() * 1000) + 1000
          this.temp.createTime = new Date().toLocaleString()
          this.temp.amount = 0
          this.list.unshift(this.temp)
          this.total++
          this.dialogFormVisible = false
          this.$notify({
            title: '成功',
            message: '创建成功',
            type: 'success',
            duration: 2000
          })
        }
      })
    },
    handleUpdate(row) {
      this.temp = Object.assign({}, row) // copy obj
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
          // 实际项目中这里应该调用API更新数据
          const index = this.list.findIndex(v => v.id === this.temp.id)
          this.list.splice(index, 1, tempData)
          this.dialogFormVisible = false
          this.$notify({
            title: '成功',
            message: '更新成功',
            type: 'success',
            duration: 2000
          })
        }
      })
    },
    handleDelete(row, index) {
      this.$confirm('确认删除该会员吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        // 实际项目中这里应该调用API删除数据
        this.list.splice(index, 1)
        this.total--
        this.$message({
          type: 'success',
          message: '删除成功!'
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    handleDownload() {
      this.downloadLoading = true
      // 实际项目中这里应该调用API导出数据
      setTimeout(() => {
        this.downloadLoading = false
        this.$message({
          message: '导出成功',
          type: 'success'
        })
      }, 2000)
    }
  }
}
</script>

<style lang="scss" scoped>
.filter-container {
  padding-bottom: 10px;
  .filter-item {
    display: inline-block;
    vertical-align: middle;
    margin-bottom: 10px;
    margin-right: 10px;
  }
}

.link-type {
  color: #409EFF;
  text-decoration: none;
  &:hover {
    color: #2a8ce8;
  }
}
</style> 