<script setup lang="ts">
import Request from '@/api/axios';
import { ElMessage } from 'element-plus';
import { ref, watch } from 'vue';
import { useUserStore } from '@/stores/user';
import { useDebounceFn } from '@vueuse/core';

interface Problem {
  id: string;
  title: string;
}

const userStore = useUserStore();
const { id } = userStore;
const TraingData = ref({
  title: '',
  introduction: '',
  user_id: '',
  password: '',
  is_private: false,// 是否私有，默认公开
  problems: [],// 已选择的题目
});
const handleClick = async () => {

  TraingData.value.user_id = id;
  TraingData.value.problems = AccpetedProblem.value;
  console.log(TraingData.value);
  ElMessage({
    message: JSON.stringify(TraingData.value),
    duration: 2000,
  })
  // await Request.post({ url: '/training/create', data: TraingData.value })
  //   .then(res => console.log(res))
  //   .catch(err => console.log(err));
}
const AccpetedProblem = ref([]);
const PendingProblem = ref([]);


// 创建防抖的搜索函数
const searchProblem = async (keyword: string) => {
  loading.value = true;
  if (keyword.trim() === '') {
    PendingProblem.value = [];
    loading.value = false;
    return;
  }

  await Request.get({ url: `/problem/search?keyword=${keyword}` })
    .then(res => PendingProblem.value = res.info.map((p: Problem) => ({ label: p.title, value: p.id })))
    .catch(err => console.log(err))
    .finally(() => loading.value = false);
};

// 使用 VueUse 的 useDebounceFn 创建防抖函数
const handleSearchProblem = useDebounceFn(searchProblem, 500);

const loading = ref(false);
</script>
<template>
  <el-form style="margin: 0 auto;max-width: 600px;" ref="form" :model="TraingData" label-width="auto">
    <el-form-item required label="标题" width="100px">
      <el-input v-model="TraingData.title"></el-input>
    </el-form-item>
    <el-form-item label="简介">
      <el-input type="textarea" autosize v-model="TraingData.introduction"></el-input>
    </el-form-item>
    <el-form-item label="密码">
      <el-input type="password" v-model="TraingData.password"></el-input>
    </el-form-item>
    <el-form-item label="是否私有">
      <el-switch v-model="TraingData.is_private"></el-switch>
    </el-form-item>
    <el-form-item label="题目选择">
      <el-select-v2 v-model="AccpetedProblem" filterable :options="PendingProblem" placeholder="选择题目" style="width: 240px"
        :loading="loading" multiple remote :remote-method="handleSearchProblem">
        </el-select-v2>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="handleClick">
        <el-icon type="plus">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1024 1024">
            <path fill="currentColor"
              d="m249.6 417.088 319.744 43.072 39.168 310.272L845.12 178.88 249.6 417.088zm-129.024 47.168a32 32 0 0 1-7.68-61.44l777.792-311.04a32 32 0 0 1 41.6 41.6l-310.336 775.68a32 32 0 0 1-61.44-7.808L512 516.992l-391.424-52.736z">
            </path>
          </svg>
        </el-icon>创建</el-button>
    </el-form-item>
  </el-form>

</template>
