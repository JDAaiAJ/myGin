<template>
  <el-sub-menu
      v-for="item in routes"
      :key="item.path"
      :index="item.path"
      v-if="hasChildren(item)"
  >
    <template #title>
      <span>{{ item.meta.title }}</span>
    </template>
    <!-- 递归调用自身 -->
    <SidebarMenu :routes="item.children" @menu-click="$emit('menu-click')" />
  </el-sub-menu>

  <el-menu-item
      v-for="item in routes"
      :key="item.path"
      :index="item.path"
      v-if="!hasChildren(item)"
      @click="$emit('menu-click', item.path)"
  >
    <span>{{ item.meta.title }}</span>
  </el-menu-item>
</template>

<script setup>
import { defineProps } from 'vue'

const props = defineProps({
  routes: {
    type: Array,
    required: true
  }
})

const hasChildren = (route) => {
  return route && route.children && route.children.length > 0
}
</script>
