<script setup>import {ref, computed} from "vue";
import {ElMessage} from 'element-plus';
import {getUserMonthSalaryListService, getUserSalaryListService} from "../../api/salary";



// 薪资列表
const monthSalaryList = ref([]);

// 搜索框内容
const name = ref('')
const month = ref('')

// 页面加载时设置默认月份为当前月份
const today = new Date()
const currentYear = today.getFullYear()
const currentMonth = String(today.getMonth() + 1).padStart(2, '0') // 补零
month.value = `${currentYear}-${currentMonth}`

const loadingMain = ref(false)

const getUserMonthSalaryList = async () => {
  let params = {
    month: month.value ? month.value : '',
    name: name.value ? name.value : ''
  }
  loadingMain.value = true
  let result = await getUserMonthSalaryListService(params);

  if (result.code === 200) {
    monthSalaryList.value = result.data.UserMonthSalaryList;
    ElMessage.success(result.message ? result.message : '查询成功');
  } else {
    ElMessage.warning(result.message ? result.message : '查询失败');
  }
  loadingMain.value = false
  // console.log(monthSalaryList.value)
}

getUserMonthSalaryList()

// 重置搜索条件
const resetSearch = () => {
  name.value = ''
  month.value = ''
}
</script>
<template>
  <el-card class="page-container">

    <!-- 上部分：搜索表单 -->
    <div class="search-form-container">
      <el-form inline>
        <el-form-item label="选择月份：" size="large">
          <el-date-picker
              v-model="month"
              type="month"
              placeholder="请选择月份"
              value-format="YYYY-MM" style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="员工姓名：" size="large">
          <el-input v-model="name" placeholder="输入姓名" size="large" style="width: 260px"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="getUserMonthSalaryList()" size="large">搜索</el-button>
          <el-button @click="resetSearch" size="large">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <!-- 卡片列表滚动区域 -->
    <div class="card-list-container" v-loading="loadingMain">
      <div v-if="monthSalaryList.length > 0" class="card-list">
        <el-card v-for="(item, index) in monthSalaryList" :key="index" shadow="hover" class="salary-card">
          <div class="card-header">
            <span><strong>月份：</strong> {{ item.month }}</span>
            <el-tag type="primary">{{ item.name }}</el-tag>
          </div>
          <div class="card-body">
            <p><strong>职位：</strong>
              <span v-if="item.type === 2">车位</span>
              <span v-else-if="item.type === 3">裁床</span>
              <span v-else-if="item.type === 4">尾部</span>
              <span v-else>未知职位</span>
            </p>
            <p><strong>薪资总额：</strong> {{ item.total }}</p>
          </div>
          <div class="card-footer">
            <el-button type="success" @click="">查看详情</el-button>
          </div>
        </el-card>
      </div>

      <!-- 空状态 -->
      <el-empty v-else description="没有数据"/>
    </div>
  </el-card>
</template>
<style lang="scss" scoped>
.page-container {
  min-height: 100%;
  box-sizing: border-box;
  background-color: #f0f0f0; // 整体背景颜色为灰色

  .search-form-container {
    width: 98%;
    background-color: #fff; // 纯白色背景
    border-radius: 5px;
    padding: 20px;
    margin-bottom: 15px;

    &:last-child {
      margin-bottom: 10px;
    }
  }

  .card-list-container, .el-empty {
    margin-top: 20px;
    max-height: 600px; // 固定最大高度
    overflow-y: auto; // 超出显示滚动条
    background-color: #fff; // 纯白色背景
    border-radius: 5px;
    padding: 20px;

    &::-webkit-scrollbar {
      width: 6px;
    }

    &::-webkit-scrollbar-thumb {
      background-color: #e4e7ed;
      border-radius: 3px;
    }
  }

  .card-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 20px;
  }

  .salary-card {
    background-color: #fff; // 卡片背景为纯白色
    border-radius: 5px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    transition: all 0.3s ease;
    display: flex;
    flex-direction: column;
    justify-content: space-between;

    &:hover {
      transform: translateY(-5px);
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      font-weight: bold;
      margin-bottom: 10px;

      span {
        font-size: 16px;
      }

      .el-tag {
        font-size: 14px;
      }
    }

    .card-body {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: center;

      p {
        margin: 8px 0;
        font-size: 14px;
      }
    }

    .card-footer {
      margin-top: 10px;
      text-align: right;

      .el-button {
        font-size: 14px;
      }
    }
  }
}
</style>

