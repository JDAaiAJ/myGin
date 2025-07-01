// 引入路由
import {createRouter,createWebHistory} from "vue-router";
import loginVue from "@/views/Login.vue";
import homeVue from "@/views/Home.vue";
import UserManageVue from "@/views/user/UserManage.vue";
import FactoryManageVue from "@/views/factory/FactoryManage.vue";
import ClothingManageVue from "@/views/clothing/ClothingManage.vue";
import SalaryManageVue from "@/views/salary/SalaryManage.vue";
import SalaryRecordVue from "@/views/record/SalaryRecord.vue";

// 路由配置
const routes = [
  {path: '/', redirect: '/login'},
  {path: '/login', component: loginVue},
  {path: '/home', component: homeVue, redirect: '/user/manage', children: [
          {
              path: '/user/manage',
              component: UserManageVue,
              meta: { title: '用户管理', keepAlive: true }
          },
          {
              path: '/factory/manage',
              component: FactoryManageVue,
              meta: { title: '工厂管理', keepAlive: true}
          },
          {
              path: '/clothing/manage',
              component: ClothingManageVue,
              meta: { title: '服装管理' , keepAlive: true}
          },
          {
              path: '/salary/record',
              component: SalaryRecordVue,
              meta: { title: '每日录入' , keepAlive: true}
          },
          {
              path: '/salary/manage',
              component: SalaryManageVue,
              meta: { title: '薪资管理', keepAlive: true }
          }

    ]},

];

//创建路由器
const router = createRouter({
  history: createWebHistory(),
  routes: routes
})
// 导出 router
export default router;