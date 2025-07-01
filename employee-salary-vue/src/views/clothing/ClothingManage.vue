<script setup>
import {nextTick, onMounted, ref} from 'vue'
import {ElMessage} from 'element-plus'



//服装列表
const clothingList = ref([])

//服饰搜索框输入内容
const code = ref('')

//分页条数据模型
const pageNum = ref(1)//当前页
const total = ref(0)//总条数
const pageSize = ref(20)//每页条数

//获取所有用户记录数据
import {
  addClothingService,
  deleteClothingImageService,
  deleteClothingService, getClothingDetailService,
  getClothingListService, updateClothingService
} from "../../api/clothing";

const loadingMain = ref(false)

const getClothingList = async () => {
  let params = {
    pageNum: pageNum.value,
    pageSize: pageSize.value,
    code: code.value ? code.value : ''
  }
  loadingMain.value = true
  let result = await getClothingListService(params);

  if (result.code === 200) {
    clothingList.value = result.data.ClothingData;
    total.value = result.data.total;
    ElMessage.success(result.message ? result.message : '查询成功');
  } else {
    ElMessage.warning(result.message ? result.message : '查询失败');
  }
  loadingMain.value = false
}

getClothingList()

//当每页条数发生了变化，调用此函数
const onSizeChange = (size) => {
  pageSize.value = size
  getClothingList()
}
//当前页码发生变化，调用此函数
const onCurrentChange = (num) => {
  pageNum.value = num
  getClothingList()
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
const clothingForm = ref({
  id: '',
  code: '',
  price: '',
  source: '',
  image: '',
})

// 表单验证规则
const rules = {
  code: [
    {required: true, message: '服饰编号不能为空', trigger: 'blur'}
  ],
  price: [
    {required: true, message: '单价不能为空', trigger: 'blur'},
    {pattern: /^\d+(\.\d{1})?$/, message: '最多一位小数', trigger: 'blur'}
  ],
  source: []
}

// 重置表单
const addClothingButton = () => {
  clothingForm.value.id = ''
  clothingForm.value.code = ''
  clothingForm.value.price = ''
  clothingForm.value.source = ''
  clothingForm.value.image = ''
  fileList.value = []
  dialogVisible.value = true
}

// 提交函数
const submitForm = async () => {
  console.log(clothingForm.value)
  const result = await addClothingService(clothingForm.value)
  console.log(result)
  if (result.code === 200) {
    ElMessage.success(result.message || '添加成功')
    dialogVisible.value = false
    await getClothingList() // 刷新列表
  } else {
    ElMessage.error(result.message || '添加失败')
  }
}


// 限制小数点后最多一位
const limitDecimal = () => {
  if (clothingForm.value.price) {
    let value = clothingForm.value.price.toString();
    // 使用正则表达式将小数部分限制为最多一位
    value = value.replace(/^(\d+(\.\d{0,1})?).*$/, '\$1');
    clothingForm.value.price = value;
  }
};

//删除服饰
const deleteClothing = async (id) => {
  let result = await deleteClothingService(id);
  console.log(result)
  if (result.code === 200){
    ElMessage.success(result.message? result.message : '删除成功');
    await getClothingList() // 刷新列表
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
const editClothingConfirm = async () => {

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

// 上传前校验（可选）
const beforeUpload = (file) => {
  const isValidType = ['image/jpeg', 'image/png', 'image/gif'].includes(file.type);
  if (!isValidType) {
    ElMessage.error('只能上传 JPG、PNG 或 GIF 格式的图片');
    return false;
  }
  const isValidSize = file.size / 1024 / 1024 < 2; // 小于 2MB
  if (!isValidSize) {
    ElMessage.error('图片大小不能超过 2MB');
    return false;
  }
  return isValidType && isValidSize;
};

// 上传成功回调
const handleSuccess = (response, file, fileList) => {
  ElMessage.success('图片上传成功');
  clothingForm.value.image = response.data.filePath; // 接收后端返回的图片路径
};

const handleExceed = () => {
  ElMessage.warning('只能上传一张图片');
};

// 删除图片时同步清空 clothingForm.image
const handleRemove = async (file, fileList) => {
  clothingForm.value.image = ''; // 清空图片路径
  ElMessage.success('图片已清空');
};

const baseUrl = 'http://192.168.235.129:8080' // 或者从环境变量中获取

</script>
<template>
  <el-card class="page-container">
    <template #header>
      <div class="header">
        <span>服饰管理</span>
        <el-button type="success" @click="addClothingButton()" >添加服饰</el-button>
      </div>

    </template>
    <!-- 搜索表单 -->
    <el-form inline>
      <el-form-item label="服饰编号：" size="large">
        <el-input v-model="code" placeholder="输入编号" size="large" style="width: 260px"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="getClothingList();" size="large">搜索</el-button>
        <el-button @click="code = ''" size="large">重置</el-button>
      </el-form-item>
    </el-form>
      <el-table :data="clothingList" border style="width: 100%;margin-top: 10px" v-loading="loadingMain" height="500">
        <el-table-column label="序号" type="index"></el-table-column>
        <el-table-column label="服饰编号" prop="code"></el-table-column>
        <el-table-column label="单价" prop="price" ></el-table-column>
        <el-table-column label="图片" prop="image">
          <template #default="scope">
            <el-image :src="baseUrl + scope.row.image" style="width: 80px; height: 80px" :preview-src-list="[baseUrl + scope.row.image]" :preview-teleported="true"></el-image>
          </template>
        </el-table-column>
        <el-table-column label="来源" prop="source"></el-table-column>
        <el-table-column label="添加人" prop="e_name"></el-table-column>
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
            <el-button type="danger" @click="deleteClothing(scope.row.id)">删除</el-button>
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

  <!-- 添加服饰的对话框 -->
  <el-dialog v-model="dialogVisible" title="添加服饰" width="30%">
    <el-form :model="clothingForm" label-width="100px" :rules="rules" ref="formRef">
      <el-form-item label="服饰编号：" prop="code">
        <el-input v-model="clothingForm.code" placeholder="请输入服饰编号"/>
      </el-form-item>

      <el-form-item label="单价：" prop="price">
        <el-input
            v-model="clothingForm.price"
            type="number"
            placeholder="请输入数字"
            @input="limitDecimal"
        />
      </el-form-item>

      <el-form-item label="来源：">
        <el-input v-model="clothingForm.source" placeholder="请输入来源"/>
      </el-form-item>

      <!-- 图片上传 -->
      <el-form-item label="图片：" prop="image">
        <div v-if="fileList.length < 1">
          <el-upload
              action="/api/clothingImageUpload"
              :on-success="handleSuccess"
              :before-upload="beforeUpload"
              :limit="1"
              list-type="picture-card"
              :file-list="fileList"
              accept="image/*"
              :on-exceed="handleExceed"
              :on-remove="handleRemove"
          >
            <el-button type="primary">点击上传</el-button>
          </el-upload>
        </div>
        <div v-else>
          <el-alert
              title="已上传一张图片"
              type="success"
              show-icon
              :closable="false"
          />
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="submitForm">提交</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="editDialogVisible" title="编辑服饰" width="30%">
    <el-form :model="clothingForm" label-width="100px" :rules="rules" ref="formRef">
      <el-form-item label="服饰编号：" prop="code">
        <el-input v-model="clothingForm.code" placeholder="请输入服饰编号"/>
      </el-form-item>

      <el-form-item label="单价：" prop="price">
        <el-input
            v-model="clothingForm.price"
            placeholder="请输入数字"
            @input="limitDecimal"
        />
      </el-form-item>

      <el-form-item label="来源：">
        <el-input v-model="clothingForm.source" placeholder="请输入来源"/>
      </el-form-item>

      <!--图片-->
      <el-form-item label="图片：" prop="image">
        <el-upload
            action="/api/clothingImageUpload"
            :on-success="handleSuccess"
            :before-upload="beforeUpload"
            :limit="1"
            list-type="picture-card"
            :file-list="fileList"
            accept="image/*"
            :on-exceed="handleExceed"
            :on-remove="handleRemove"
        >
          <el-button type="primary">点击上传</el-button>
        </el-upload>
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
