<script setup>

import {nextTick, onActivated, onDeactivated, onMounted, ref, watch} from 'vue'
import {ElMessage} from 'element-plus'

//存储user列表数据
const userList = ref([])

const userInfo = ref([])

import {useUserStore} from '@/stores/useUserStore'

const userStore = useUserStore()

onMounted(() => {
  getUserFactoryList()
  userInfo.value = userStore.user
})

//user搜索框输入内容
const name = ref('')

//分页条数据模型
const pageNum = ref(1)//当前页
const total = ref(0)//总条数
const pageSize = ref(20)//每页条数

import {
  deleteUserService,
  getUserDetailService,
  getUserFactoryService,
  getUserListService,
  updateUserService,
  userRegisterService
} from "../../api/user";

const loadingMain = ref(false)

const getAllUserList = async () => {
  let params = {
    pageNum: pageNum.value,
    pageSize: pageSize.value,
    name: name.value ? name.value : ''
  }
  // loadingMain.value = true
  let result = await getUserListService(params);

  if (result.code === 200) {
    userList.value = result.data.UserData;
    total.value = result.data.total;
    ElMessage.success(result.message ? result.message : '查询成功');

  } else {
    ElMessage.warning(result.message ? result.message : '查询失败');
  }
  // loadingMain.value = false
}


getAllUserList()


//当每页条数发生了变化，调用此函数
const onSizeChange = (size) => {
  pageSize.value = size
  getAllUserList()
}
//当前页码发生变化，调用此函数
const onCurrentChange = (num) => {
  pageNum.value = num
  getAllUserList()
}

// 切换密码可见性
const togglePasswordVisibility = (row) => {
  if(userInfo.value.type === 0){
    if (row.showPassword === undefined) {
      // 如果没有定义 showPassword 属性，则使用 Vue.set 添加响应式属性
      row.showPassword = true;
    } else {
      row.showPassword = !row.showPassword;
    }
  }else {
    ElMessage.warning('暂无权限');
  }

};

//所有工厂列表
const factoryList = ref([])

//获取工厂列表
const getUserFactoryList = async () => {
  let result = await getUserFactoryService();
  if (result.code === 200) {
    factoryList.value = result.data;
  } else {
    ElMessage.warning(result.message ? result.message : '获取失败');
  }
}

// 控制对话框显示隐藏
const dialogVisible = ref(false)

const formRef = ref(null);

// 表单数据模型
const userForm = ref({
  id: '',
  name: '',
  username: '',
  password: '',
  f_id: 0,
  type: 0
})

const handleFIdChange = (newVal) => {
  if (newVal === 0) {
    userForm.value.type = 0
    isPositionDisabled.value = true
  } else {
    isPositionDisabled.value = false
  }
}

// 重置表单
const addUserButton = () => {
  userForm.value = {
    id: '',
    name: '',
    username: '',
    password: '',
    f_id: 0,
    type: 0
  }

  // 手动触发一次 watch 回调
  handleFIdChange(userForm.value.f_id)

  if (formRef.value) {
    formRef.value.resetFields()
  }

  dialogVisible.value = true
}

// 表单验证规则
const rules = {
  name: [
    { required: true, message: '真实姓名不能为空', trigger: 'blur' }
  ],
  username: [
    { required: true, message: '用户名不能为空', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '密码不能为空', trigger: 'blur' }
  ],
}

// 提交函数
const submitForm = async () => {
  // console.log(userForm.value)
  formRef.value.validate(async (valid) => {
    if (valid) {
      // 注册逻辑
      console.log(userForm.value)
      let result = await userRegisterService(userForm.value);
      console.log(result)
      if (result.code === 200) {
        ElMessage.success(result.message? result.message : '注册成功');
        dialogVisible.value = false
        await getAllUserList() // 刷新列表
      } else {
        ElMessage.error(result.message? result.message : '注册失败');
      }
    } else {
      console.log('表单验证失败');
      return false;
    }
  });

}

const editDialogVisible = ref(false);

// 编辑用户按钮
const editUserBtn = async (id) => {
  let params = { id: id };
  let result = await getUserDetailService(params);
  if (result.code === 200) {
    userForm.value = result.data;

    await nextTick();
    editDialogVisible.value = true;
    ElMessage.success(result.message ? result.message : '获取成功');
  } else {
    ElMessage.warning(result.message ? result.message : '获取失败');
  }
};

// 时间戳转换为完整日期时间格式（YYYY-MM-DD HH:mm:ss）
function formatDate(timestamp) {
  var date = new Date(timestamp * 1000); // 注意：如果你传入的是秒级时间戳，需要乘以 1000

  var year = date.getFullYear();
  var month = String(date.getMonth() + 1).padStart(2, '0'); // 补零
  var day = String(date.getDate()).padStart(2, '0');
  var hours = String(date.getHours()).padStart(2, '0');
  var minutes = String(date.getMinutes()).padStart(2, '0');
  var seconds = String(date.getSeconds()).padStart(2, '0');

  return `${year}年${month}月${day}日 ${hours}:${minutes}:${seconds}`;
}

// 编辑服饰确认按钮
const editUserConfirm = async () => {

  console.log(userForm.value)
  // 提交服饰信息
  let result = await updateUserService(userForm.value);
  if (result.code === 200){
    ElMessage.success(result.message? result.message : '编辑成功');
    editDialogVisible.value = false;
    await getAllUserList();
  } else {
    ElMessage.warning(result.message? result.message : '编辑失败');
  }
}

// 控制职位是否被禁用
const isPositionDisabled = ref(false)

watch(() => userForm.value.f_id, handleFIdChange)

//删除用户
const deleteUser = async (id) => {
  let result = await deleteUserService(id);
  console.log(result)
  if (result.code === 200){
    ElMessage.success(result.message? result.message : '删除成功');
    await getAllUserList() // 刷新列表
  }else{
    ElMessage.warning(result.message? result.message : '删除失败');
  }
}

</script>
<template>
  <el-card class="page-container">
    <template #header>
      <div class="header">
        <span>用户管理</span>
        <el-button v-if="userInfo.type === 1 || userInfo.u_id === 1" type="success" @click="addUserButton()" >添加用户</el-button>
      </div>
    </template>
    <!-- 搜索表单 -->
    <el-form inline>
      <el-form-item label="真实姓名：" size="large">
        <el-input v-model="name" placeholder="输入名称" size="large" style="width: 260px"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="getAllUserList();" size="large">搜索</el-button>
        <el-button @click="name = ''" size="large">重置</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="userList" border style="width: 100%;margin-top: 20px" v-loading="loadingMain">
      <el-table-column label="序号" type="index" ></el-table-column>
      <el-table-column label="用户ID" prop="id" width="100"></el-table-column>
      <el-table-column label="真实姓名" prop="name" width="120"></el-table-column>
      <el-table-column label="用户名" prop="username"></el-table-column>
      <el-table-column label="密码" prop="password" width="180">
        <template #default="scope">
          <span @click="togglePasswordVisibility(scope.row)" style="cursor: pointer">
          {{ scope.row.showPassword ? scope.row.password : '************' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="所属工厂" prop="f_name" width="260">
        <template #default="scope">
          <el-tag type="success" v-if="scope.row.f_id === 0">无</el-tag>
          <el-tag type="success" v-else >{{ scope.row.f_name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="职位" prop="type">
        <template #default="scope">
          <el-tag v-if="scope.row.type === 0 && scope.row.id === 1">管理员</el-tag>
          <el-tag v-if="scope.row.type === 0 && scope.row.id !== 1">暂无职位</el-tag>
          <el-tag v-if="scope.row.type === 1">厂长</el-tag>
          <el-tag v-if="scope.row.type === 2">车位</el-tag>
          <el-tag v-if="scope.row.type === 3">裁床</el-tag>
          <el-tag v-if="scope.row.type === 4">尾部</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" prop="status">
        <template #default="scope">
          <el-tag v-if="scope.row.status === 0">离职</el-tag>
          <el-tag v-if="scope.row.status === 1">在职</el-tag>

        </template>
      </el-table-column>
      <el-table-column label="入库时间" prop="insert_time" width="200">
        <template #default="scope">
          <el-tag>{{ formatDate(scope.row.insert_time) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="更新时间" prop="insert_time" width="200">
        <template #default="scope">
          <el-tag>{{ formatDate(scope.row.update_time) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button v-if="scope.row.id !== 1" type="primary" @click="editUserBtn(scope.row.id)">编辑</el-button>
          <el-button v-if="scope.row.id !== 1 && userStore.user.u_id === 1" type="danger" @click="deleteUser(scope.row.id)">删除</el-button>

        </template>
      </el-table-column>
      <template #empty>
        <el-empty description="没有数据"/>
      </template>
    </el-table>
  </el-card>
  <!-- 分页条 -->
  <el-pagination v-model:current-page="pageNum" v-model:page-size="pageSize" :page-sizes="[20, 50, 100, 150]"
                 layout="jumper, total, sizes, prev, pager, next" background :total="total" @size-change="onSizeChange"
                 @current-change="onCurrentChange" style="margin-top: 10px; justify-content: flex-end"/>

  <!-- 添加用户的对话框 -->
  <el-dialog v-model="dialogVisible" title="添加用户" width="26%">
    <el-form :model="userForm" label-width="100px" :rules="rules" ref="formRef">
      <el-form-item label="真实姓名：" style="font-weight: 600" prop="name">
        <el-input v-model="userForm.name" placeholder="请输入真实姓名" style="width: 300px"/>
      </el-form-item>

      <el-form-item label="用户名：" style="font-weight: 600" prop="username">
        <el-input v-model="userForm.username" placeholder="请输入用户名" style="width: 300px"/>
      </el-form-item>

      <el-form-item label="密码：" style="font-weight: 600" prop="password">
        <el-input v-model="userForm.password" placeholder="请输入密码" style="width: 300px"/>
      </el-form-item>

      <el-form-item label="所属工厂：" style="font-weight: 600">
        <el-select v-model="userForm.f_id" placeholder="请选择工厂" style="width: 200px" filterable>
          <el-option
              label="暂不选择"
              :value="0"
          />
          <el-option
              v-for="factory in factoryList"
              :key="factory.id"
              :label="factory.name"
              :value="factory.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="职位：" style="font-weight: 600">
        <el-select v-model="userForm.type" placeholder="请选择职位" style="width: 200px" filterable :disabled="isPositionDisabled">
          <el-option :value="0" label="暂无职位" />
          <el-option :value="1" label="厂长" />
          <el-option :value="2" label="车位" />
          <el-option :value="3" label="裁床" />
          <el-option :value="4" label="尾部" />
        </el-select>
      </el-form-item>

    </el-form>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="submitForm">提交</el-button>
    </template>
  </el-dialog>

  <!-- 编辑用户的对话框 -->
  <el-dialog v-model="editDialogVisible" title="添加用户" width="26%">
    <el-form :model="userForm" label-width="100px" :rules="rules" ref="formRef">
      <el-form-item label="真实姓名：" style="font-weight: 600" prop="name">
        <el-input v-model="userForm.name" placeholder="请输入真实姓名" style="width: 300px"/>
      </el-form-item>

      <el-form-item label="用户名：" style="font-weight: 600" prop="username">
        <el-input v-model="userForm.username" placeholder="请输入用户名" style="width: 300px"/>
      </el-form-item>

      <el-form-item label="密码：" style="font-weight: 600" prop="password">
        <el-input v-model="userForm.password" placeholder="请输入密码" style="width: 300px"/>
      </el-form-item>

      <el-form-item label="所属工厂：" style="font-weight: 600">
        <el-select v-model="userForm.f_id" placeholder="请选择工厂" style="width: 200px" filterable>
          <el-option
              label="暂不选择"
              :value="0"
          />
          <el-option
              v-for="factory in factoryList"
              :key="factory.id"
              :label="factory.name"
              :value="factory.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="职位：" style="font-weight: 600">
        <el-select v-model="userForm.type" placeholder="请选择职位" style="width: 200px" filterable :disabled="isPositionDisabled">
          <el-option :value="0" label="暂无职位" />
          <el-option :value="1" label="厂长" />
          <el-option :value="2" label="车位" />
          <el-option :value="3" label="裁床" />
          <el-option :value="4" label="尾部" />
        </el-select>
      </el-form-item>

    </el-form>

    <template #footer>
      <el-button @click="editDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="editUserConfirm">提交</el-button>
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
