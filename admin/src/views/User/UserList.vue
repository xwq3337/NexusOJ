<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { UserInfo } from '@/api/axios/type'
import Request from '@/api/axios'

const loading = ref(false)
const userList = ref<UserInfo[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)

const searchForm = reactive({
  keyword: '',
  role: undefined as number | undefined,
  status: 'all' as 'all' | 'banned' | 'normal'
})

// 封禁相关
const banDialogVisible = ref(false)
const banForm = reactive({
  userId: '',
  endTime: null as Date | null
})

// 角色相关
const roleDialogVisible = ref(false)
const roleForm = reactive({
  userId: '',
  role: 1
})

const fetchUserList = async () => {
  loading.value = true
    await Request.get({
      url: '/admin/user-list',
      params: {
        page: page.value,
        pageSize: pageSize.value,
        ...searchForm
      }
    }).then((res) => {
      userList.value = res.info
      total.value = res.count
    }).catch((err) => {
      ElMessage.error('获取用户列表失败')
      console.error(err)
    }).finally(() => {
      loading.value = false
    })
}

const handleSearch = () => {
  page.value = 1
  fetchUserList()
}

const resetSearch = () => {
  searchForm.keyword = ''
  searchForm.role = undefined
  searchForm.status = 'all'
  handleSearch()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchUserList()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  fetchUserList()
}

const getRoleTagType = (role: number) => {
  switch (role) {
    case 1:
      return ''
    case 2:
      return 'warning'
    case 3:
      return 'danger'
    default:
      return 'info'
  }
}

const getRoleName = (role: number) => {
  switch (role) {
    case 1:
      return '普通用户'
    case 2:
      return '管理员'
    case 3:
      return '超级管理员'
    default:
      return '未知'
  }
}

const handleBan = (user: UserInfo) => {
  banForm.userId = user.id
  banForm.endTime = null
  banDialogVisible.value = true
}

const confirmBan = async () => {
  if (!banForm.endTime) {
    ElMessage.warning('请选择封禁结束时间')
    return
  }

  try {
    await Request.post({
      url: `/user/${banForm.userId}/ban`,
      data: {
        endTime: banForm.endTime.toISOString()
      }
    })
    ElMessage.success('封禁成功')
    banDialogVisible.value = false
    fetchUserList()
  } catch (error) {
    ElMessage.error('封禁失败')
    console.error(error)
  }
}

const handleUnban = async (user: UserInfo) => {
  try {
    await ElMessageBox.confirm('确定要解封该用户吗？', '提示', {
      type: 'warning'
    })
    await Request.post({
      url: `/user/${user.id}/unban`
    })
    ElMessage.success('解封成功')
    fetchUserList()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('解封失败')
      console.error(error)
    }
  }
}

const handleUpdateRole = (user: UserInfo) => {
  roleForm.userId = user.id
  roleForm.role = user.userRole
  roleDialogVisible.value = true
}

const confirmUpdateRole = async () => {
  try {
    await Request.put({
      url: `/user/${roleForm.userId}/role`,
      data: {
        role: roleForm.role
      }
    })
    ElMessage.success('修改角色成功')
    roleDialogVisible.value = false
    fetchUserList()
  } catch (error) {
    ElMessage.error('修改角色失败')
    console.error(error)
  }
}

onMounted(() => {
  fetchUserList()
})
</script>

<template>
  <div class="user-list-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="用户名/邮箱/学校"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="用户角色">
          <el-select v-model="searchForm.role" style="width: 240px" placeholder="选择角色" clearable>
            <el-option label="普通用户" :value="1" />
            <el-option label="管理员" :value="2" />
            <el-option label="超级管理员" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" style="width: 240px" placeholder="选择状态" clearable>
            <el-option label="全部" value="all" />
            <el-option label="正常" value="normal" />
            <el-option label="封禁" value="banned" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table
      v-loading="loading"
      :data="userList"
      border
      style="width: 100%"
    >
      <el-table-column prop="username" label="用户名" width="70" show-overflow-tooltip />
      <el-table-column prop="email" label="邮箱" width="200" />
      <el-table-column label="姓名" min-width="120">
        <template #default="{ row }">
           {{ row.nickname }}
        </template>
      </el-table-column>
      <el-table-column prop="school" label="学校" width="150" />
      <el-table-column prop="rating" label="评分" width="60" align="center" />
      <el-table-column label="提交/通过" width="100" align="center">
        <template #default="{ row }">
          {{ row.submission }}/{{ row.accept }}
        </template>
      </el-table-column>
      <el-table-column label="用户角色" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getRoleTagType(row.user_role)">
            {{ getRoleName(row.user_role) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status ? 'danger' : 'success'">
            {{ row.status ? '已封禁' : '正常' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button
              v-if="!row.status"
              type="danger"
              size="small"
              @click="handleBan(row)"
            >
              封禁
            </el-button>
            <el-button
              v-else
              type="success"
              size="small"
              @click="handleUnban(row)"
            >
              解封
            </el-button>
            <el-button
              type="primary"
              size="small"
              @click="handleUpdateRole(row)"
            >
              修改角色
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 封禁对话框 -->
    <el-dialog
      v-model="banDialogVisible"
      title="封禁用户"
      width="400px"
    >
      <el-form :model="banForm" label-width="100px">
        <el-form-item label="封禁结束时间">
          <el-date-picker
            v-model="banForm.endTime"
            type="datetime"
            placeholder="选择日期时间"
            format="YYYY-MM-DD HH:mm:ss"
            :min-date="new Date()"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="banDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmBan">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 修改角色对话框 -->
    <el-dialog
      v-model="roleDialogVisible"
      title="修改用户角色"
      width="400px"
    >
      <el-form :model="roleForm" label-width="100px">
        <el-form-item label="用户角色">
          <el-select v-model="roleForm.role" style="width: 240px" placeholder="选择角色">
            <el-option label="普通用户" :value="1" />
            <el-option label="管理员" :value="2" />
            <el-option label="超级管理员" :value="3" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="roleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmUpdateRole">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>


<style scoped>
.user-list-container {
  padding: 20px;
}

.search-bar {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

:deep(.el-button-group) {
  .el-button {
    margin-left: 0 !important;
  }
}
</style>
