<template>
  <div class="member-container">
    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm" class="demo-form-inline">
        <el-form-item label="会员姓名">
          <el-input v-model="searchForm.name" placeholder="请输入会员姓名"></el-input>
        </el-form-item>
        <el-form-item label="手机号码">
          <el-input v-model="searchForm.phone" placeholder="请输入手机号码"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="resetForm">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 操作按钮 -->
    <div class="operation-bar">
      <el-button type="primary" @click="handleAdd">新增会员</el-button>
      <el-button type="success" @click="handleBatchRecharge">批量充值</el-button>
    </div>

    <!-- 会员列表 -->
    <el-table :data="memberList" border style="width: 100%">
      <el-table-column prop="id" label="会员ID" width="80"></el-table-column>
      <el-table-column prop="name" label="会员姓名" width="120"></el-table-column>
      <el-table-column prop="phone" label="手机号码" width="120"></el-table-column>
      <el-table-column prop="balance" label="充值余额" width="120">
        <template slot-scope="scope">
          ¥{{ scope.row.balance }}
        </template>
      </el-table-column>
      <el-table-column prop="gift_balance" label="赠送余额" width="120">
        <template slot-scope="scope">
          ¥{{ scope.row.gift_balance }}
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="注册时间" width="180"></el-table-column>
      <el-table-column label="操作" width="300">
        <template slot-scope="scope">
          <el-button size="mini" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button size="mini" type="success" @click="handleRecharge(scope.row)">充值</el-button>
          <el-button size="mini" type="warning" @click="handleConsume(scope.row)">消费</el-button>
          <el-button size="mini" type="info" @click="handleTransaction(scope.row)">交易记录</el-button>
          <el-button size="mini" type="danger" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination">
      <el-pagination
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
        :current-page="currentPage"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="total">
      </el-pagination>
    </div>

    <!-- 新增/编辑会员对话框 -->
    <el-dialog :title="dialogTitle" :visible.sync="dialogVisible">
      <el-form :model="memberForm" :rules="rules" ref="memberForm" label-width="100px">
        <el-form-item label="会员姓名" prop="name">
          <el-input v-model="memberForm.name"></el-input>
        </el-form-item>
        <el-form-item label="手机号码" prop="phone">
          <el-input v-model="memberForm.phone"></el-input>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="submitForm">确 定</el-button>
      </div>
    </el-dialog>

    <!-- 充值对话框 -->
    <el-dialog title="会员充值" :visible.sync="rechargeVisible">
      <el-form :model="rechargeForm" :rules="rechargeRules" ref="rechargeForm" label-width="100px">
        <el-form-item label="充值金额" prop="amount">
          <el-input-number v-model="rechargeForm.amount" :min="0" :precision="2"></el-input-number>
        </el-form-item>
        <el-form-item label="赠送金额" prop="giftAmount">
          <el-input-number v-model="rechargeForm.giftAmount" :min="0" :precision="2"></el-input-number>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="rechargeVisible = false">取 消</el-button>
        <el-button type="primary" @click="submitRecharge">确 定</el-button>
      </div>
    </el-dialog>

    <!-- 交易记录对话框 -->
    <el-dialog title="交易记录" :visible.sync="transactionVisible" width="80%">
      <el-table :data="transactionList" border style="width: 100%">
        <el-table-column prop="created_at" label="交易时间" width="180"></el-table-column>
        <el-table-column prop="type" label="交易类型" width="120">
          <template slot-scope="scope">
            <el-tag :type="getTransactionTypeTag(scope.row.type)">
              {{ getTransactionTypeText(scope.row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="交易金额" width="120">
          <template slot-scope="scope">
            ¥{{ scope.row.amount }}
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="账户余额" width="120">
          <template slot-scope="scope">
            ¥{{ scope.row.balance }}
          </template>
        </el-table-column>
        <el-table-column prop="gift_balance" label="赠送余额" width="120">
          <template slot-scope="scope">
            ¥{{ scope.row.gift_balance }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="交易说明"></el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'MemberManagement',
  data() {
    return {
      // 搜索表单
      searchForm: {
        name: '',
        phone: ''
      },
      // 会员列表数据
      memberList: [],
      // 分页相关
      currentPage: 1,
      pageSize: 10,
      total: 0,
      // 对话框显示控制
      dialogVisible: false,
      rechargeVisible: false,
      transactionVisible: false,
      // 表单数据
      memberForm: {
        name: '',
        phone: ''
      },
      rechargeForm: {
        amount: 0,
        giftAmount: 0
      },
      // 交易记录列表
      transactionList: [],
      // 表单验证规则
      rules: {
        name: [
          { required: true, message: '请输入会员姓名', trigger: 'blur' }
        ],
        phone: [
          { required: true, message: '请输入手机号码', trigger: 'blur' },
          { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
        ]
      },
      rechargeRules: {
        amount: [
          { required: true, message: '请输入充值金额', trigger: 'blur' }
        ]
      }
    }
  },
  computed: {
    dialogTitle() {
      return this.memberForm.id ? '编辑会员' : '新增会员'
    }
  },
  methods: {
    // 搜索
    handleSearch() {
      this.fetchMemberList()
    },
    // 重置表单
    resetForm() {
      this.searchForm = {
        name: '',
        phone: ''
      }
      this.fetchMemberList()
    },
    // 获取会员列表
    async fetchMemberList() {
      // TODO: 调用后端API获取会员列表
    },
    // 新增会员
    handleAdd() {
      this.memberForm = {
        name: '',
        phone: ''
      }
      this.dialogVisible = true
    },
    // 编辑会员
    handleEdit(row) {
      this.memberForm = { ...row }
      this.dialogVisible = true
    },
    // 提交会员表单
    submitForm() {
      this.$refs.memberForm.validate(async (valid) => {
        if (valid) {
          // TODO: 调用后端API保存会员信息
          this.dialogVisible = false
          this.fetchMemberList()
        }
      })
    },
    // 充值
    handleRecharge(row) {
      this.rechargeForm = {
        amount: 0,
        giftAmount: 0
      }
      this.rechargeVisible = true
    },
    // 提交充值
    submitRecharge() {
      this.$refs.rechargeForm.validate(async (valid) => {
        if (valid) {
          // TODO: 调用后端API进行充值
          this.rechargeVisible = false
          this.fetchMemberList()
        }
      })
    },
    // 消费
    handleConsume(row) {
      // TODO: 实现消费功能
    },
    // 查看交易记录
    async handleTransaction(row) {
      // TODO: 调用后端API获取交易记录
      this.transactionVisible = true
    },
    // 删除会员
    handleDelete(row) {
      this.$confirm('确认删除该会员吗？', '提示', {
        type: 'warning'
      }).then(async () => {
        // TODO: 调用后端API删除会员
        this.fetchMemberList()
      })
    },
    // 获取交易类型标签样式
    getTransactionTypeTag(type) {
      const typeMap = {
        RECHARGE: 'success',
        CONSUME: 'warning',
        GIFT: 'info'
      }
      return typeMap[type] || 'info'
    },
    // 获取交易类型文本
    getTransactionTypeText(type) {
      const typeMap = {
        RECHARGE: '充值',
        CONSUME: '消费',
        GIFT: '赠送'
      }
      return typeMap[type] || type
    },
    // 分页相关方法
    handleSizeChange(val) {
      this.pageSize = val
      this.fetchMemberList()
    },
    handleCurrentChange(val) {
      this.currentPage = val
      this.fetchMemberList()
    }
  },
  created() {
    this.fetchMemberList()
  }
}
</script>

<style scoped>
.member-container {
  padding: 20px;
}
.search-bar {
  margin-bottom: 20px;
}
.operation-bar {
  margin-bottom: 20px;
}
.pagination {
  margin-top: 20px;
  text-align: right;
}
</style> 