<script setup>

import {ElMessage} from "element-plus";
import {computed, onMounted, ref, watch} from "vue";
import {
  ArrowLeftBold,
  ArrowRightBold
} from '@element-plus/icons-vue'
import {getSalaryClothingService} from "../../api/clothing";
import {addUserRecordService, addUserRecordSpecialService} from "../../api/record";
import {
  deleteDailySalaryService, deleteDailySalarySpecialService,
  getDailySalaryDetailService,
  getDailySalarySpecialDetailService, getUserSalaryListService,
  updateDailySalaryService, updateDailySalarySpecialService
} from "../../api/salary";
import {getUsersBYFactoryService} from "../../api/user";



//薪资列表
const userSalaryList = ref([])

// 薪资总额
const allTotal = ref('0.0');

//fs搜索框输入内容
const user_id = ref('')

const month = ref('')

//用户列表
const users = ref([])

//服饰列表
const clothingList = ref([])

// 页面加载时设置默认月份为当前月份
onMounted(() => {
  const today = new Date();
  let currentYear = today.getFullYear();
  let currentMonth = today.getMonth(); // 获取当前月份（0-11）

  // 如果当前是1月，则上个月是去年12月
  if (currentMonth === 0) {
    currentYear -= 1;
    currentMonth = 12;
  }

  const prevMonth = String(currentMonth).padStart(2, '0'); // 补零
  month.value = `${currentYear}-${prevMonth}`;
  setDefaultDate()
  getUsersByFactoryList()
  getSalaryClothingList()
})


const loadingMain = ref(false)

//获取用户列表
const getUsersByFactoryList = async () => {
  let result = await getUsersBYFactoryService();
  if (result.code === 200) {
    users.value = result.UsersData;
  } else {
    ElMessage.warning(result.message ? result.message : '获取失败');
  }
}

//获取服饰列表
const getSalaryClothingList = async () => {
  let result = await getSalaryClothingService();
  if (result.code === 200) {
    clothingList.value = result.data;
  } else {
    ElMessage.warning(result.message ? result.message : '获取失败');
  }
}

// 添加表单项
const dailyForm = ref({
  user_id: '',
  clothing_id: '',
  user_name: '',
  date: '',
  quantity: 1,
});

// 添加特殊表单项
const dailySpecialForm = ref({
  user_id: '',
  user_name: '',
  name: '',
  price: '',
  date: '',
  quantity: 1,
});


// 监听 month 的变化，自动更新日期
watch(() => month.value, () => {
  setDefaultDate()
})

// 判断是否禁用某一天（不在当前 month 范围内的日期）
const disabledDate = (date) => {
  if (!month.value) return false;


  const selectedMonth = new Date(month.value); // 当前选择的月份
  const dateToCheck = new Date(date); // 当前遍历的日期

  // 只允许选择相同年份和月份的日期
  return (
      dateToCheck.getFullYear() !== selectedMonth.getFullYear() ||
      dateToCheck.getMonth() !== selectedMonth.getMonth()
  );
};

// 设置默认 dailyForm.date 为 month 的第一天
const setDefaultDate = () => {
  if (month.value) {
    dailyForm.value.date = `${month.value}-01`
    dailySpecialForm.value.date = `${month.value}-01`
  }
}

// 判断是否是当月的第一天
const isFirstDay = computed(() => {
  if (!dailyForm.value.date) return true
  const date = new Date(dailyForm.value.date)
  return date.getDate() === 1
})

// 判断是否是当月的最后一天
const isLastDay = computed(() => {
  if (!dailyForm.value.date) return true
  const date = new Date(dailyForm.value.date)
  const year = date.getFullYear()
  const month = date.getMonth()

  // 获取该月的最后一天
  const lastDay = new Date(year, month + 1, 0).getDate()
  return date.getDate() === lastDay
})

// 改变日期函数
const changeDate = (delta) => {
  if (!dailyForm.value.date) return

  const date = new Date(dailyForm.value.date)
  date.setDate(date.getDate() + delta)

  // 更新日期
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  dailyForm.value.date = `${year}-${month}-${day}`
}

// 判断是否是当月的第一天
const isFirstDaySpecial = computed(() => {
  if (!dailySpecialForm.value.date) return true
  const date = new Date(dailySpecialForm.value.date)
  return date.getDate() === 1
})

// 判断是否是当月的最后一天
const isLastDaySpecial = computed(() => {
  if (!dailySpecialForm.value.date) return true
  const date = new Date(dailySpecialForm.value.date)
  const year = date.getFullYear()
  const month = date.getMonth()

  // 获取该月的最后一天
  const lastDay = new Date(year, month + 1, 0).getDate()
  return date.getDate() === lastDay
})

// 改变日期函数
const changeDateSpecial = (delta) => {
  if (!dailySpecialForm.value.date) return

  const date = new Date(dailySpecialForm.value.date)
  date.setDate(date.getDate() + delta)

  // 更新日期
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  dailySpecialForm.value.date = `${year}-${month}-${day}`
}

const getUserSalaryList = async () => {
  if (user_id.value !== '') {
    let params = {
      month: month.value ? month.value : '',
      user_id: user_id.value ? user_id.value : ''
    }
    loadingMain.value = true
    let result = await getUserSalaryListService(params);

    if (result.code === 200) {
      userSalaryList.value = result.data.UserSalaryList;
      allTotal.value = result.data.AllTotal || '0.0';
      ElMessage.success(result.message ? result.message : '查询成功');
    } else {
      ElMessage.warning(result.message ? result.message : '查询失败');
    }
    loadingMain.value = false
  } else {
    ElMessage.warning('请选择录入员工');
  }

}

// 监听 user_id 变化，同步更新 dailyForm.name
watch(() => user_id.value, (newVal) => {
  if (newVal) {
    const selectedUser = users.value.find(user => user.id === newVal);
    if (selectedUser) {
      dailyForm.value.user_name = selectedUser.name;
      dailyForm.value.user_id = selectedUser.id;
      dailySpecialForm.value.user_name = selectedUser.name;
      dailySpecialForm.value.user_id = selectedUser.id;
    }
  } else {
    dailyForm.value.name = '';
    dailyForm.value.user_id = '';
    dailySpecialForm.value.name = '';
    dailySpecialForm.value.user_id = '';
  }
})

//录入员工每日工资
const addUserRecord = async () => {
  console.log(dailyForm.value)
  if (dailyForm.value.user_id === '') {
    ElMessage.warning('请先选择录入员工');
    return;
  }

  if (dailyForm.value.clothing_id === '') {
    ElMessage.warning('请选择录入服饰');
    return;
  }
  const res = await addUserRecordService(dailyForm.value);
  if (res.code === 200) {
    ElMessage.success('录入成功');
    await getUserSalaryList()
  } else {
    ElMessage.error(res.message);
  }
}

//录入员工每日工资
const addUserSpecialRecord = async () => {
  console.log(dailyForm.value)
  if (dailySpecialForm.value.user_id === '') {
    ElMessage.warning('请先选择录入员工');
    return;
  }

  if (dailySpecialForm.value.name === '') {
    ElMessage.warning('请选择特殊名');
    return;
  }

  if (dailySpecialForm.value.price === '') {
    ElMessage.warning('请输入单价');
    return;
  }
  const res = await addUserRecordSpecialService(dailySpecialForm.value);
  if (res.code === 200) {
    ElMessage.success('录入成功');
    await getUserSalaryList()
  } else {
    ElMessage.error(res.message);
  }
}

//删除每日薪资
const deleteDailySalaryBtn = async (id,type) => {
  if (type === 1){
    let result = await deleteDailySalaryService(id);
    if (result.code === 200) {
      ElMessage.success(result.message ? result.message : '删除成功');
      await getUserSalaryList() // 刷新列表
    } else {
      ElMessage.warning(result.message ? result.message : '删除失败');
    }
  }else {
    let result = await deleteDailySalarySpecialService(id);
    if (result.code === 200) {
      ElMessage.success(result.message ? result.message : '删除成功');
      await getUserSalaryList() // 刷新列表
    } else {
      ElMessage.warning(result.message ? result.message : '删除失败');
    }
  }
}

const editDialogVisible = ref(false);

const editDialogVisibleSpecial = ref(false);

// 表单数据模型
const dailySalaryForm = ref({
  id: '',
  u_id: '',
  c_id: '',
  quantity: 1
})

// 特殊表单数据模型
const dailySalarySpecialForm = ref({
  id: '',
  u_id: '',
  name: '',
  price: '',
  quantity: 1
})

//编辑每日薪资
const editDailySalaryBtn = async (id,type) => {
  let params = {
    id: id
  }
  if (type === 1) {
    //服装编辑
    let result = await getDailySalaryDetailService(params);
    if (result.code === 200) {
      dailySalaryForm.value = result.data;
      editDialogVisible.value = true;
      ElMessage.success(result.message ? result.message : '获取成功');
    } else {
      ElMessage.warning(result.message ? result.message : '获取失败');
    }
  } else {
    //特殊编辑

    //服装编辑
    let result = await getDailySalarySpecialDetailService(params);
    if (result.code === 200) {
      dailySalarySpecialForm.value = result.data;
      editDialogVisibleSpecial.value = true;
      ElMessage.success(result.message ? result.message : '获取成功');
    } else {
      ElMessage.warning(result.message ? result.message : '获取失败');
    }
  }

}

// 判断是否是当月的第一天
const isFirstDayDetails = computed(() => {
  if (!dailySalaryForm.value.date) return true
  const date = new Date(dailySalaryForm.value.date)
  return date.getDate() === 1
})

// 判断是否是当月的最后一天
const isLastDayDetails = computed(() => {
  if (!dailySalaryForm.value.date) return true
  const date = new Date(dailySalaryForm.value.date)
  const year = date.getFullYear()
  const month = date.getMonth()

  // 获取该月的最后一天
  const lastDay = new Date(year, month + 1, 0).getDate()
  return date.getDate() === lastDay
})

// 改变日期函数（详情日期）
const changeDateDetails = (delta) => {
  if (!dailySalaryForm.value.date) return

  const date = new Date(dailySalaryForm.value.date)
  date.setDate(date.getDate() + delta)

  // 更新日期
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  dailySalaryForm.value.date = `${year}-${month}-${day}`
}

// 判断是否是当月的第一天
const isFirstDayDetailsSpecial = computed(() => {
  if (!dailySalarySpecialForm.value.date) return true
  const date = new Date(dailySalarySpecialForm.value.date)
  return date.getDate() === 1
})

// 判断是否是当月的最后一天
const isLastDayDetailsSpecial = computed(() => {
  if (!dailySalarySpecialForm.value.date) return true
  const date = new Date(dailySalarySpecialForm.value.date)
  const year = date.getFullYear()
  const month = date.getMonth()

  // 获取该月的最后一天
  const lastDay = new Date(year, month + 1, 0).getDate()
  return date.getDate() === lastDay
})

// 改变日期函数（详情日期）
const changeDateDetailsSpecial = (delta) => {
  if (!dailySalarySpecialForm.value.date) return

  const date = new Date(dailySalarySpecialForm.value.date)
  date.setDate(date.getDate() + delta)

  // 更新日期
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')

  dailySalarySpecialForm.value.date = `${year}-${month}-${day}`
}

//编辑服饰确认按钮
const editDailySalaryConfirm = async () => {
  console.log(dailySalaryForm.value)
  let result = await updateDailySalaryService(dailySalaryForm.value);
  if (result.code === 200) {
    ElMessage.success(result.message ? result.message : '编辑成功');
    editDialogVisible.value = false;
    await getUserSalaryList() // 刷新列表
  } else {
    ElMessage.warning(result.message ? result.message : '编辑失败');
  }
}

//编辑服饰确认按钮
const editDailySalarySpecialConfirm = async () => {
  console.log(dailySalarySpecialForm.value)
  let result = await updateDailySalarySpecialService(dailySalarySpecialForm.value);
  if (result.code === 200) {
    ElMessage.success(result.message ? result.message : '编辑成功');
    editDialogVisibleSpecial.value = false;
    await getUserSalaryList() // 刷新列表
  } else {
    ElMessage.warning(result.message ? result.message : '编辑失败');
  }
}

// 限制小数点后最多一位
const limitDecimal = () => {
  if (dailySpecialForm.value.price) {
    let value = dailySpecialForm.value.price.toString();
    // 使用正则表达式将小数部分限制为最多一位
    value = value.replace(/^(\d+(\.\d{0,1})?).*$/, '\$1');
    dailySpecialForm.value.price = value;
  }
};

const baseUrl = 'http://192.168.235.129:8080' // 或者从环境变量中获取

</script>

<template>
  <el-card class="page-container">
    <div class="header" style="font-size: 18px;margin-bottom: 10px">薪资录入</div>

    <!-- 搜索表单 -->
    <el-form inline>
      <el-form-item label="选择月份：">
        <el-date-picker
            v-model="month"
            type="month"
            placeholder="请选择月份"
            value-format="YYYY-MM"
            style="width: 200px"
        />
      </el-form-item>
      <el-form-item label="录入员工姓名：">
        <el-select v-model="user_id" placeholder="请选择姓名" style="width: 260px" clearable>
          <el-option
              v-for="user in users"
              :key="user.id"
              :label="user.name"
              :value="user.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="getUserSalaryList()">选择</el-button>
      </el-form-item>
    </el-form>

    <!-- 生产情况展示与新增组件布局 -->
    <el-row :gutter="20">
      <!-- 左侧：员工生产情况 -->
      <el-col :span="16">
        <el-card>
          <el-table :data="userSalaryList" border style="width: 100%;height: 620px">
            <el-table-column prop="date" label="日期"></el-table-column>
            <el-table-column prop="code" label="服饰编号/特殊">
              <template #default="scope">
                  {{ scope.row.type === 1 ? scope.row.code : scope.row.name }}
              </template>
            </el-table-column>
            <el-table-column label="图片" prop="image">
              <template #default="scope">
                <el-image :src="baseUrl + scope.row.image" style="width: 60px; height: 60px" :preview-src-list="[baseUrl + scope.row.image]" :preview-teleported="true"></el-image>
              </template>
            </el-table-column>
            <el-table-column prop="price" label="单价"></el-table-column>
            <el-table-column prop="quantity" label="数量"></el-table-column>
            <el-table-column prop="total" label="总价"></el-table-column>
            <el-table-column label="操作" :width="200">
              <template #default="scope" >
                <el-button type="primary" @click="editDailySalaryBtn(scope.row.id, scope.row.type)">编辑</el-button>
                <el-button type="danger" @click="deleteDailySalaryBtn(scope.row.id, scope.row.type)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
          <!--总薪资-->
          <div style="font-size: 18px; font-weight: bold; text-align: right;margin-top: 10px">
            本月总薪资：
            <span style="color: #F5222D;">￥{{ allTotal }}</span>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：添加每日生产情况 -->
      <el-col :span="8">
        <el-card>
          <div style="font-weight: 600;margin-bottom: 10px;text-align: center">添加每日生产情况</div>
          <el-form :model="dailyForm" label-width="100px" label-position="right">
            <el-form-item label="员工姓名：" style="font-weight: 600">
              <el-input v-model="dailyForm.user_name" placeholder="未选择" disabled style="width: 160px"/>
            </el-form-item>

            <el-form-item label="日期：" style="font-weight: 600">
              <!-- 上一日按钮 -->
              <el-button :disabled="isFirstDay" @click="changeDate(-1)" :icon="ArrowLeftBold"></el-button>
              <el-date-picker
                  v-model="dailyForm.date"
                  type="date"
                  placeholder="选择日期"
                  value-format="YYYY-MM-DD"
                  :disabledDate="disabledDate"
                  style="width: 160px"

              />
              <!-- 下一日按钮 -->
              <el-button :disabled="isLastDay" @click="changeDate(1)" :icon="ArrowRightBold"></el-button>
            </el-form-item>

            <el-form-item label="服饰编号：" style="font-weight: 600">
              <el-select v-model="dailyForm.clothing_id" placeholder="请选择服饰" style="width: 260px" filterable>
                <el-option
                    v-for="clothing in clothingList"
                    :key="clothing.id"
                    :label="clothing.code"
                    :value="clothing.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item label="数量：" style="font-weight: 600">
              <el-input-number v-model="dailyForm.quantity" :min="1"/>
            </el-form-item>

            <el-button type="success" @click="addUserRecord" style="width: 100%;margin-top:5px">添加</el-button>
          </el-form>
        </el-card>
        <el-card style="margin-top: 10px">
          <div style="font-weight: 600;margin-bottom: 10px;text-align: center">添加每日特殊生产情况</div>
          <el-form :model="dailySpecialForm" label-width="100px" label-position="right">
            <el-form-item label="员工姓名：" style="font-weight: 600">
              <el-input v-model="dailySpecialForm.user_name" placeholder="未选择" disabled style="width: 160px"/>
            </el-form-item>

            <el-form-item label="日期：" style="font-weight: 600">
              <!-- 上一日按钮 -->
              <el-button :disabled="isFirstDaySpecial" @click="changeDateSpecial(-1)" :icon="ArrowLeftBold"></el-button>
              <el-date-picker
                  v-model="dailySpecialForm.date"
                  type="date"
                  placeholder="选择日期"
                  value-format="YYYY-MM-DD"
                  :disabledDate="disabledDate"
                  style="width: 160px"

              />
              <!-- 下一日按钮 -->
              <el-button :disabled="isLastDaySpecial" @click="changeDateSpecial(1)" :icon="ArrowRightBold"></el-button>
            </el-form-item>

            <el-form-item label="特殊名：" style="font-weight: 600">
              <el-input v-model="dailySpecialForm.name" style="width: 200px" placeholder="请输入特殊名"/>
            </el-form-item>

            <el-form-item label="单价：" style="font-weight: 600">
              <el-input v-model="dailySpecialForm.price" style="width: 200px" placeholder="请输入单价"
                        @input="limitDecimal"/>
            </el-form-item>


            <el-form-item label="数量：" style="font-weight: 600">
              <el-input-number v-model="dailySpecialForm.quantity" :min="1"/>
            </el-form-item>

            <el-button type="primary" @click="addUserSpecialRecord" style="width: 100%;margin-top:5px">添加</el-button>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </el-card>
  <el-dialog v-model="editDialogVisible" title="编辑每日薪资" width="30%">
    <el-form :model="dailySalaryForm" label-width="100px">

      <el-form-item label="日期：" style="font-weight: 600" size="large">
        <!-- 上一日按钮 -->
        <el-button :disabled="isFirstDayDetails" @click="changeDateDetails(-1)" :icon="ArrowLeftBold"></el-button>
        <el-date-picker
            v-model="dailySalaryForm.date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            :disabledDate="disabledDate"
            style="width: 160px"

        />
        <!-- 下一日按钮 -->
        <el-button :disabled="isLastDayDetails" @click="changeDateDetails(1)" :icon="ArrowRightBold"></el-button>
      </el-form-item>
      <el-form-item label="服饰编号：" style="font-weight: 600" size="large">
        <el-select v-model="dailySalaryForm.c_id" placeholder="请选择服饰" style="width: 260px" clearable filterable>
          <el-option
              v-for="clothing in clothingList"
              :key="clothing.id"
              :label="clothing.code"
              :value="clothing.id"

          />
        </el-select>
      </el-form-item>
      <el-form-item label="数量：" style="font-weight: 600" size="large">
        <el-input-number v-model="dailySalaryForm.quantity" :min="1"/>
      </el-form-item>


    </el-form>

    <template #footer>
      <el-button @click="editDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="editDailySalaryConfirm">提交</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="editDialogVisibleSpecial" title="编辑特殊每日薪资" width="30%">
    <el-form :model="dailySalarySpecialForm" label-width="100px">

      <el-form-item label="日期：" style="font-weight: 600" size="large">
        <!-- 上一日按钮 -->
        <el-button :disabled="isFirstDayDetailsSpecial" @click="changeDateDetailsSpecial(-1)" :icon="ArrowLeftBold"></el-button>
        <el-date-picker
            v-model="dailySalarySpecialForm.date"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
            :disabledDate="disabledDate"
            style="width: 160px"

        />
        <!-- 下一日按钮 -->
        <el-button :disabled="isLastDayDetailsSpecial" @click="changeDateDetailsSpecial(1)" :icon="ArrowRightBold"></el-button>
      </el-form-item>

      <el-form-item label="特殊名：" style="font-weight: 600" size="large">
        <el-input v-model="dailySalarySpecialForm.name" style="width: 200px" placeholder="请输入特殊名"/>
      </el-form-item>

      <el-form-item label="单价：" style="font-weight: 600" size="large">
        <el-input v-model="dailySalarySpecialForm.price" style="width: 200px" placeholder="请输入单价"
                  @input="limitDecimal"/>
      </el-form-item>

      <el-form-item label="数量：" style="font-weight: 600" size="large">
        <el-input-number v-model="dailySalarySpecialForm.quantity" :min="1"/>
      </el-form-item>


    </el-form>

    <template #footer>
      <el-button @click="editDialogVisibleSpecial = false">取消</el-button>
      <el-button type="primary" @click="editDailySalarySpecialConfirm">提交</el-button>
    </template>
  </el-dialog>


</template>
<style lang="scss" scoped>
.page-container {
  min-height: 93%;
  box-sizing: border-box;

  .header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}


</style>