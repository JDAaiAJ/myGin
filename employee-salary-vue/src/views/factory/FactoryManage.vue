<script setup>
import {nextTick, onMounted, ref} from 'vue'
import {ElMessage} from 'element-plus'



//工厂列表
const factoryList = ref([])

//服饰搜索框输入内容
const name = ref('')

//分页条数据模型
const pageNum = ref(1)//当前页
const total = ref(0)//总条数
const pageSize = ref(20)//每页条数

//获取所有用户记录数据
import {getClothingListService,addClothingService,deleteClothingService,getClothingDetailService,updateClothingService} from "@/api/clothing.js";
import {deleteClothingImageService} from "../../api/clothing";
import {addFactoryService, deleteFactoryService, getFactoryListService} from "../../api/factory";

const loadingMain = ref(false)

const getFactoryList = async () => {
  let params = {
    pageNum: pageNum.value,
    pageSize: pageSize.value,
    name: name.value ? name.value : ''
  }
  // console.log(params)
  loadingMain.value = true
  let result = await getFactoryListService(params);
  // console.log(result)
  if (result.code === 200) {
    factoryList.value = result.data.FactoryData;
    total.value = result.data.total;
    ElMessage.success(result.message ? result.message : '查询成功');
  } else {
    ElMessage.warning(result.message ? result.message : '查询失败');
  }
  loadingMain.value = false
}

getFactoryList()

//当每页条数发生了变化，调用此函数
const onSizeChange = (size) => {
  pageSize.value = size
  getFactoryList()
}
//当前页码发生变化，调用此函数
const onCurrentChange = (num) => {
  pageNum.value = num
  getFactoryList()
}

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

// 控制对话框显示隐藏
const dialogVisible = ref(false)

// 表单数据模型
const factoryForm = ref({
  id: '',
  name: '',
})

// 表单验证规则
const rules = {
  name: [
    {required: true, message: '厂不能为空', trigger: 'blur'}
  ],
}

// 重置表单
const addFactoryButton = () => {
  factoryForm.value.id = ''
  factoryForm.value.name = ''
  dialogVisible.value = true
}

// 提交函数
const submitForm = async () => {
  const result = await addFactoryService(factoryForm.value)
  console.log(result)
  if (result.code === 200) {
    ElMessage.success(result.message || '添加成功')
    dialogVisible.value = false
    await getFactoryList() // 刷新列表
  } else {
    ElMessage.error(result.message || '添加失败')
  }
}


//删除服饰
const deleteFactory = async (id) => {
  let result = await deleteFactoryService(id);
  console.log(result)
  if (result.code === 200){
    ElMessage.success(result.message? result.message : '删除成功');
    await getFactoryList() // 刷新列表
  }else{
    ElMessage.warning(result.message? result.message : '删除失败');
  }
}

const editDialogVisible = ref(false);

const fileList = ref([]); // 用于 el-upload 显示已有图片

const originalImage = ref(''); // 新增响应式变量

// 编辑服饰按钮
const editClothingBtn = async (id) => {
  let params = { id: id };
  let result = await getClothingDetailService(params);
  if (result.code === 200) {
    clothingForm.value = result.data;

    // 保留原始 image 路径用于提交更新
    const imagePath = clothingForm.value.image;
    originalImage.value = imagePath; // 记录原始图片路径

    // 设置 fileList 显示图片
    if (imagePath) {
      fileList.value = [
        {
          name: '图片',
          url: baseUrl + imagePath,
          uid: -1,
        },
      ];
    } else {
      fileList.value = [];
    }

    clothingForm.value.image = imagePath;

    await nextTick();
    editDialogVisible.value = true;
    ElMessage.success(result.message ? result.message : '获取服饰成功');
  } else {
    ElMessage.warning(result.message ? result.message : '获取服饰失败');
  }
};

// 编辑服饰确认按钮
const editFactoryConfirm = async () => {

  // 判断是否需要删除图片
  if (!clothingForm.value.image && originalImage.value) {
    // 用户删除了图片，需要调用删除接口
    const deleteResult = await deleteClothingImageService({
      filePath: originalImage.value
    });

    if (deleteResult.code !== 200) {
      ElMessage.warning(deleteResult.message || '图片删除失败');
    }else {
      ElMessage.success('图片已从服务器删除');
    }


  }

  // 提交服饰信息
  let result = await updateClothingService(clothingForm.value);

  if (result.code === 200){
    ElMessage.success(result.message? result.message : '编辑服饰成功');
    editDialogVisible.value = false;
    await getClothingList();
  } else {
    ElMessage.warning(result.message? result.message : '编辑服饰失败');
  }
}

</script>
<template>
  <el-card class="page-container">
    <template #header>
      <div class="header">
        <span>工厂管理</span>
        <el-button type="success" @click="addFactoryButton()" >添加工厂</el-button>
      </div>

    </template>
    <!-- 搜索表单 -->
    <el-form inline>
      <el-form-item label="工厂名：" size="large">
        <el-input v-model="name" placeholder="输入编号" size="large" style="width: 260px"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="getFactoryList();" size="large">搜索</el-button>
        <el-button @click="name = ''" size="large">重置</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="factoryList" border style="width: 100%;margin-top: 20px" v-loading="loadingMain">
      <el-table-column label="序号" type="index"></el-table-column>
      <el-table-column label="工厂名" prop="name"></el-table-column>
      <el-table-column label="入库时间" prop="insert_time">
        <template #default="scope">
          <el-tag>{{ formatDate(scope.row.insert_time) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="更新时间" prop="insert_time">
        <template #default="scope">
          <el-tag>{{ formatDate(scope.row.update_time) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="scope">
          <el-button type="primary" @click="editClothingBtn(scope.row.id)">编辑</el-button>
          <el-button type="danger" @click="deleteFactory(scope.row.id)">删除</el-button>
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

  <!-- 添加工厂的对话框 -->
  <el-dialog v-model="dialogVisible" title="添加工厂" width="30%">
    <el-form :model="factoryForm" label-width="100px" :rules="rules" ref="formRef">
      <el-form-item label="服饰编号：" prop="code">
        <el-input v-model="factoryForm.name" placeholder="请输入厂名"/>
      </el-form-item>

    </el-form>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="submitForm">提交</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="editDialogVisible" title="编辑工厂" width="30%">
    <el-form :model="factoryForm" label-width="100px" :rules="rules" ref="formRef">
      <el-form-item label="厂名：" prop="code">
        <el-input v-model="factoryForm.name" placeholder="请输入厂名"/>
      </el-form-item>

    </el-form>

    <template #footer>
      <el-button @click="editDialogVisible = false">取消</el-button>
      <el-button type="primary" @click="editClothingConfirm">提交</el-button>
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
