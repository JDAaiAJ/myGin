<script setup>
import {computed, ref, onMounted, watch} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import {ElMessage, ElMessageBox} from 'element-plus'
import {useUserStore} from '@/stores/useUserStore'

const userStore = useUserStore()
const route = useRoute()
const router = useRouter()

// 获取菜单路由
const homeRoute = router.options.routes.find(route => route.path === '/home')

// 用户信息相关
const user = computed(() => userStore.user)
const username = computed(() => user.value?.name || '')
const usertype = computed(() => user.value?.type ?? null)

// 用户类型名称
const getUserTypeName = (type) => {
  switch (type) {
    case 0:
      return '管理员'
    case 1:
      return '厂长'
    case 2:
      return '车位'
    case 3:
      return '裁床'
    case 4:
      return '尾部'
    default:
      return '用户'
  }
}

// 当前激活的tab
let activeTab = ref('')
// 所有tab
let allTabs = ref([])
// 当前激活的路由
let activePath = ref('')

// 工具函数：递归提取所有带 component 的路由
function getFlatRoutes(routes) {
  let result = []
  routes?.forEach(route => {
    if (route.component) {
      result.push(route)
    }
    if (route.children && route.children.length > 0) {
      result = result.concat(getFlatRoutes(route.children))
    }
  })
  return result
}

// 获取菜单路由（包含所有可点击的叶子节点路由）
const menuRoutes = computed(() => {
  return getFlatRoutes(homeRoute?.children || [])
})

// 缓存菜单路由
const menuRoutesMap = computed(() => {
  const map = new Map()
  menuRoutes.value.forEach(route => {
    map.set(route.path, route)
  })
  return map
})

// 初始化
onMounted(() => {
  const storedTabs = sessionStorage.getItem("allTabs")
  const storedActiveTab = sessionStorage.getItem("activeTab")

  if (storedTabs) {
    allTabs.value = JSON.parse(storedTabs)
  } else {
    // 获取第一个菜单项的实际跳转路径
    const firstMenuItem = menuRoutes.value[0]
    const defaultPath = firstMenuItem?.redirect || firstMenuItem?.children?.[0]?.path || '/user/list'
    addTab(defaultPath)
  }

  if (storedActiveTab) {
    activeTab.value = storedActiveTab
    changeActiveRoute(storedActiveTab)
  } else {
    activeTab.value = allTabs.value[0]?.name || ''
    changeActiveRoute(activeTab.value)
  }
})

// 添加新tab
function addTab(path) {
  const routeItem = menuRoutesMap.value.get(path)

  // 如果是父级菜单（没有 component），则不允许添加 tab
  if (!routeItem || !routeItem.component) return

  // 如果tab已存在，直接激活
  const existingTabIndex = allTabs.value.findIndex(tab => tab.path === path)
  if (existingTabIndex > -1) {
    activeTab.value = path
    changeActiveRoute(path)
    return
  }

  // 添加新tab
  const newTab = {
    title: routeItem.meta.title,
    name: path,
    path: path
  }

  allTabs.value.push(newTab)
  activeTab.value = path
  changeActiveRoute(path)

  // 保存到sessionStorage
  saveTabsToStorage()
}

// 点击 tab
const clickTab = (tab) => {
  let tabName = tab.paneName
  const selectTab = allTabs.value.find(t => t.name === tabName)
  if (selectTab) {
    activeTab.value = tabName
    changeActiveRoute(selectTab.path)
    saveTabsToStorage()
  }
}

// 移除 tab
const removeTab = (targetName) => {
  let tabs = allTabs.value
  let activeName = activeTab.value

  // 如果关闭的是当前激活的tab，需要切换到其他tab
  if (activeName === targetName) {
    const tabsLength = tabs.length
    if (tabsLength === 1) {
      ElMessage.warning('最后一个标签无法删除')
      return
    } else {
      const currentIndex = tabs.findIndex(tab => tab.name === targetName)
      let nextTabIndex = currentIndex === tabsLength - 1 ? currentIndex - 1 : currentIndex + 1
      const nextTab = tabs[nextTabIndex]
      if (nextTab) {
        activeName = nextTab.name
        changeActiveRoute(nextTab.path)
      }
    }
  }

  // 更新tab列表
  allTabs.value = tabs.filter(tab => tab.name !== targetName)
  activeTab.value = activeName

  // 保存到sessionStorage
  saveTabsToStorage()
}

// 跳转路由
const changeActiveRoute = (path) => {
  if (path && path !== activePath.value) {
    activePath.value = path
    router.push(path)
  }
}

// 保存tab状态到sessionStorage
const saveTabsToStorage = () => {
  sessionStorage.setItem("allTabs", JSON.stringify(allTabs.value))
  sessionStorage.setItem("activeTab", activeTab.value)
}

// 监听路由变化
watch(() => route.path, (newPath) => {
  // 检查路由是否在tab列表中
  const tabExists = allTabs.value.some(tab => tab.path === newPath)
  if (tabExists) {
    activeTab.value = newPath
  } else {
    // 如果是新路由，添加到tab列表
    addTab(newPath)
  }
})

// 退出登录
const logout = () => {
  ElMessageBox.confirm('确认退出登录吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    // 清除会话存储中的 allTabs 和 activeTab
    sessionStorage.removeItem('allTabs');
    sessionStorage.removeItem('activeTab');

    // 清除用户状态
    userStore.logout();

    // 跳转到登录页
    router.push('/login');

    // 提示用户已退出登录
    ElMessage.success('已退出登录');
  }).catch(() => {
    // 用户点击取消，不做任何操作
  });
}

const activeMenuIndex = computed(() => route.path)

const breadcrumbs = computed(() => {
  return route.matched.map(item => ({
    path: item.path,
    name: item.meta.title || item.name
  }))
})
</script>

<template>
  <el-container class="layout-container">
    <el-aside width="180px">
      <el-scrollbar>
        <div class="el-aside__logo"></div>

        <el-menu
            active-text-color="#ffd04b"
            background-color="#232323"
            text-color="#fff"
            :default-active="activeMenuIndex"
        >
          <template v-for="item in homeRoute?.children || []" :key="item.path">
            <!-- 如果有 children，则显示子菜单 -->
            <el-sub-menu v-if="item.children && item.children.length > 0" :index="item.path">
              <template #title>{{ item.meta.title }}</template>
              <el-menu-item
                  v-for="child in item.children"
                  :key="child.path"
                  :index="child.path"
                  @click="addTab(child.path)"
              >
                {{ child.meta.title }}
              </el-menu-item>
            </el-sub-menu>

            <!-- 否则显示普通菜单项 -->
            <el-menu-item v-else :index="item.path" @click="addTab(item.path)">
              {{ item.meta.title }}
            </el-menu-item>
          </template>
        </el-menu>
      </el-scrollbar>
    </el-aside>

    <el-container>
      <el-header class="top-nav">
        <!-- 左侧内容包裹容器 -->
        <div class="header-left">
          <!-- navigation图标 -->
          <img src="@/assets/icon/navigation.png" alt="navigation icon" class="nav-icon">

          <!-- 面包屑导航 -->
          <el-breadcrumb separator=">">
            <el-breadcrumb-item v-for="(item, index) in breadcrumbs" :key="index" :to="item.path"
                                style="font-size: 16px">
              {{ item.name }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>

        <!-- 右侧内容包裹容器 -->
        <div class="header-right">
          <!-- 退出登录按钮 -->
          <el-button type="primary" link @click="logout" style="font-size: 15px">退出登录</el-button>
        </div>
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
          ></el-tab-pane>

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
    align-items: center;
    padding: 0 20px;
    background: #fff;
    border-bottom: 1px solid #dcdfe6;

    .header-left {
      display: flex;
      align-items: center;

      .nav-icon {
        width: 24px;
        height: 24px;
        vertical-align: middle;
      }
    }

    .header-right {
      margin-left: auto;
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
    padding: 0;
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .el-tabs {
    height: 100%;

    .el-tabs__content {
      height: calc(100% - 40px);
      overflow: auto;
      display: flex;
      flex-direction: column;
    }
  }

  .el-footer {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    color: #666;
    background: #fff;
    border-top: 1px solid #dcdfe6;
  }

  .el-tabs--card > .el-tabs__header .el-tabs__item.is-active {
    padding-left: 10px;
    padding-right: 10px;
    color: white;
    background-color: #409EFF;
  }
}

.el-tabs--card > .el-tabs__header {
  border-bottom: 1px solid #E4E7ED;
  background-color: gainsboro;
}

:deep(.el-tabs__header) {
  margin-bottom: 0 !important;
}

:deep(.el-tabs--card > .el-tabs__header .el-tabs__item.is-active) {
  border-bottom-color: #409EFF !important;
}

.el-menu--dark .el-menu-item,
.el-menu--dark .el-sub-menu__title {
  color: #fff;
}
</style>
