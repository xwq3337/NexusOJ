<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Request from '@/api/axios'
import dayjs from 'dayjs'
import type { BlogInfo } from '@/types/blog'
import { ElMessage, ElMessageBox } from 'element-plus'
const pageBlogs = ref<BlogInfo[]>([])
const blogList = ref<BlogInfo[]>([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const fetchBlogList = async () => {
  loading.value = true
  await Request.get({
    url: '/blog/all',
    params: {
      page: currentPage.value,
      pageSize: pageSize.value,
    },
  })
    .then((res) => {
      blogList.value = res.info
      // TODO: 处理分页
      handleCurrentChange(1)
      total.value = res.info.length
    })
    .catch((err) => {
      ElMessage.error('获取博客列表失败')
      console.error('获取博客列表失败:', err)
    })
    .finally(() => {
      loading.value = false
    })
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  pageBlogs.value = blogList.value.slice((currentPage.value - 1) * pageSize.value, currentPage.value * pageSize.value)
  // fetchBlogList()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  pageBlogs.value = blogList.value.slice((currentPage.value - 1) * pageSize.value, currentPage.value * pageSize.value)
  // fetchBlogList()
}


const handleEdit = (row: BlogInfo) => {
  console.log('编辑博客:', row.id)
  // TODO: 实现编辑博客功能
}

const handleDelete = async (row: BlogInfo) => {
  try {
    await ElMessageBox.confirm('确定要删除这篇博客吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    const { result } = await Request.post({
      url: '/blog/delete',
      data: { id: row.id },
    })
    if (result === null) {
      ElMessage.success('删除成功')
      fetchBlogList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
      console.error('删除博客失败:', error)
    }
  }
}

onMounted(() => {
  fetchBlogList()
})
</script>
<template>
  <div class="blog-info-container">
    <el-card class="blog-card">
      <template #header>
        <div class="card-header">
          <span>博客列表</span>
        </div>
      </template>
      <el-table :data="pageBlogs" :pagination="true" :border="true" style="width: 100%" v-loading="loading">
        <el-table-column prop="title" label="标题" width="150" show-overflow-tooltip />
        <el-table-column prop="username" label="作者" width="100" />
        <el-table-column prop="collection" label="收藏" width="60" />
        <el-table-column prop="like" label="点赞" width="60" />
        <el-table-column prop="is_private" label="隐私性" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_private ? 'danger' : 'success'">
              {{ row.is_private ? '私有' : '公开' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ dayjs(row.created_at).format("YYYY-MM-DD HH:mm:ss") }}
          </template>
        </el-table-column>
        <el-table-column prop="tags" label="标签" width="180">
          <template #default="{ row }">
            <el-tag v-for="tag in row.tags" :key="tag" type="success">{{ tag }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">
              <el-icon>
                <Edit />
              </el-icon>
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">
              <el-icon>
                <Delete />
              </el-icon>
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination v-model:current-page="currentPage" v-model:page-size="pageSize" :page-sizes="[10, 20, 30, 50]"
          layout="total, sizes, prev, pager, next" :total="total" @size-change="handleSizeChange"
          @current-change="handleCurrentChange" />
      </div>
    </el-card>
  </div>
</template>



<style scoped>
.blog-info-container {
  padding: 20px;
}

.blog-card {
  margin-bottom: 20px;
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
</style>
