<script setup>
import { computed, ref, onMounted, nextTick, watch } from 'vue'
import { onBeforeRouteUpdate, useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/useUserStore'

const userStore = useUserStore()
const route = useRoute()
const router = useRouter()

// 获取菜单路由
const homeRoute = router.options.routes.find(route => route.path === '/home')
const menuRoutes = computed(() => {
  return homeRoute?.children || []
})

// 用户信息相关
const user = computed(() => userStore.user)
const username = computed(() => user.value?.name || '')
const usertype = computed(() => user.value?.type ?? null)

// 用户类型名称
const getUserTypeName = (type) => {
  switch (type) {
    case 0: return '管理员'
    case 1: return '厂长'
    case 2: return '车位'
    case 3: return '裁床'
    case 4: return '尾部'
    default: return '用户'
  }
}

// 当前激活的tab
let activeTab = ref("");
// 所有tab
let allTabs = ref([]);
// 当前激活的路由
let activePath = ref("");

// 缓存菜单路由
const menuRoutesMap = computed(() => {
  const map = new Map();
  menuRoutes.value.forEach(route => {
    map.set(route.path, route);
  });
  return map;
});

// 初始化
onMounted(() => {
  const storedTabs = sessionStorage.getItem("allTabs");
  const storedActiveTab = sessionStorage.getItem("activeTab");

  if (storedTabs) {
    allTabs.value = JSON.parse(storedTabs);
  } else {
    const defaultRoute = menuRoutes.value[0]?.path || '/user/manage';
    addTab(defaultRoute);
  }

  if (storedActiveTab) {
    activeTab.value = storedActiveTab;
    changeActiveRoute(storedActiveTab);
  } else {
    activeTab.value = allTabs.value[0]?.name || '';
    changeActiveRoute(activeTab.value);
  }
});

// 添加新tab
function addTab(path) {
  const routeItem = menuRoutesMap.value.get(path);
  if (!routeItem) return;

  // 如果tab已存在，直接激活
  const existingTabIndex = allTabs.value.findIndex(tab => tab.path === path);
  if (existingTabIndex > -1) {
    activeTab.value = path;
    changeActiveRoute(path);
    return;
  }

  // 添加新tab
  const newTab = {
    title: routeItem.meta.title,
    name: path,
    path: path
  };

  allTabs.value.push(newTab);
  activeTab.value = path;
  changeActiveRoute(path);

  // 保存到sessionStorage
  saveTabsToStorage();
}

// 点击 tab
const clickTab = (tab) => {
  let tabName = tab.paneName;
  const selectTab = allTabs.value.find(t => t.name === tabName);
  if (selectTab) {
    activeTab.value = tabName;
    changeActiveRoute(selectTab.path);
    saveTabsToStorage();
  }
};

// 移除 tab
const removeTab = (targetName) => {
  let tabs = allTabs.value;
  let activeName = activeTab.value;

  // 如果关闭的是当前激活的tab，需要切换到其他tab
  if (activeName === targetName) {
    const tabsLength = tabs.length;
    if (tabsLength === 1) {
      ElMessage.warning('最后一个标签无法删除');
      return;
    } else {
      // 找到当前tab的索引
      const currentIndex = tabs.findIndex(tab => tab.name === targetName);
      // 获取下一个或上一个tab
      let nextTabIndex = currentIndex === tabsLength - 1 ? currentIndex - 1 : currentIndex + 1;
      const nextTab = tabs[nextTabIndex];
      if (nextTab) {
        activeName = nextTab.name;
        changeActiveRoute(nextTab.path);
      }
    }
  }

  // 更新tab列表
  allTabs.value = tabs.filter(tab => tab.name !== targetName);
  activeTab.value = activeName;

  // 保存到sessionStorage
  saveTabsToStorage();
};

// 跳转路由
const changeActiveRoute = (path) => {
  if (path && path !== activePath.value) {
    activePath.value = path;
    router.push(path);
  }
};

// 保存tab状态到sessionStorage
const saveTabsToStorage = () => {
  sessionStorage.setItem("allTabs", JSON.stringify(allTabs.value));
  sessionStorage.setItem("activeTab", activeTab.value);
};

// 监听路由变化
watch(() => route.path, (newPath) => {
  // 检查路由是否在tab列表中
  const tabExists = allTabs.value.some(tab => tab.path === newPath);
  if (tabExists) {
    activeTab.value = newPath;
  } else {
    // 如果是新路由，添加到tab列表
    addTab(newPath);
  }
});

// 退出登录
const logout = () => {
  ElMessageBox.confirm('确认退出登录吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    userStore.logout();
    setTimeout(() => {
      router.push('/login');
      ElMessage.success('已退出登录');
    }, 500);
  }).catch(() => {});
};

const activeMenuIndex = computed(() => route.path);
</script>

<template>
  <el-container class="layout-container">
    <el-aside width="200px">
      <el-scrollbar>
        <div class="el-aside__logo"></div>
        <el-menu
            active-text-color="#ffd04b"
            background-color="#232323"
            text-color="#fff"
            :default-active="activeMenuIndex"
        >
          <el-menu-item
              v-for="item in menuRoutes"
              :key="item.path"
              :index="item.path"
              @click="addTab(item.path)"
          >
            <span>{{ item.meta.title }}</span>
          </el-menu-item>
        </el-menu>
      </el-scrollbar>
    </el-aside>

    <el-container>
      <el-header class="top-nav">
        <div class="welcome-message">
          欢迎 {{ getUserTypeName(usertype) }} : {{ username }} 进入本管理系统
        </div>
        <el-button type="primary" link @click="logout">退出登录</el-button>
      </el-header>

      <el-main>
        <!-- 动态tab标签页 -->
        <el-tabs
            v-model="activeTab"
            type="card"
            closable
            @tab-click="clickTab"
            @tab-remove="removeTab"
            style="height: 100%"
        >
          <el-tab-pane
              v-for="item in allTabs"
              :key="item.name"
              :label="item.title"
              :name="item.name"
          >

          </el-tab-pane>

          <RouterView v-slot="{ Component }">
            <KeepAlive>
              <component :is="Component"></component>
            </KeepAlive>
          </RouterView>

        </el-tabs>
      </el-main>

      <el-footer>开盘啦后台管理系统 ©2025 Created by jzh</el-footer>
    </el-container>
  </el-container>
</template>

<style lang="scss" scoped>
.layout-container {
  height: 100vh;

  .top-nav {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 20px;

    z-index: 100;
    position: relative;
    background: #ffffff;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
    border-bottom: 1px solid #dcdfe6;

    .welcome-message {
      font-size: 16px;
      color: #333;
    }
  }

  .el-aside {
    background-color: #232323;

    &__logo {
      height: 90px;
      width: 90px;
      background: url('@/assets/default.png') no-repeat center / cover;
      border-radius: 20px;
      overflow: hidden;
      margin: 10px auto 5px;
    }

    .el-menu {
      border-right: none;
    }
  }

  .el-main {
    margin-top: -10px;

  }

  .el-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    color: #666;
  }

  .el-tabs--card > .el-tabs__header .el-tabs__item.is-active {
    padding-left: 10px;
    padding-right: 10px;
    color: white;
    background-color: #409EFF;
  }
}
</style>