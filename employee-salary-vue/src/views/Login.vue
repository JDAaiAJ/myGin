<script setup>
import {ref} from 'vue';
import {userLoginService, userRegisterService} from "../api/user.js";
import {ElMessage} from "element-plus";

import { useUserStore } from '../stores/useUserStore'
import router from "../router";

const userStore = useUserStore()

const isSignUpActive = ref(false);

//注册表单
const registerForm = ref({
  username: '',
  name: '',
  password: '',
  confirmPassword: ''
})

//登录表单
const loginForm = ref({
  username: '',
  password: ''
})

function handleSignUp() {
  isSignUpActive.value = true;
  //清空注册表单
  registerForm.username = '';
  registerForm.name = '';
  registerForm.password = '';
  registerForm.confirmPassword = '';

  if (formRefRegister.value) {
    formRefRegister.value.resetFields(); // 重置表单验证状态
  }
}

function handleSignIn() {
  isSignUpActive.value = false;
  if (formRefLogin.value) {
    formRefLogin.value.resetFields();
  }

}

const formRefRegister = ref(null);

// 注册表单的验证规则
const registerRules = {
  username: [
      { required: true, message: '用户名不能为空', trigger: 'blur' }
  ],
  name: [
      { required: true, message: '真实姓名不能为空', trigger: 'blur' }
  ],
      password: [
    { required: true, message: '密码不能为空', trigger: 'blur' }
  ],
      confirmPassword: [
    { required: true, message: '确认密码不能为空', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== registerForm.value.password) {
          callback(new Error('确认密码与密码不一致'));
        } else {
          callback();
        }
      },
      trigger: 'blur'
    }
  ]
}

const formRefLogin = ref(null);

// 登录表单的验证规则
const loginRules = {
  username: [
    { required: true, message: '用户名不能为空', trigger: 'blur' }
  ],
      password: [
    { required: true, message: '密码不能为空', trigger: 'blur' }
  ]
}
// 注册
function SignUp() {
  formRefRegister.value.validate(async (valid) => {
    if (valid) {
      // 注册逻辑
      console.log(registerForm.value)
      let result = await userRegisterService(registerForm.value);
      console.log(result)
      if (result.code === 200) {
        ElMessage.success(result.message? result.message : '注册成功');
      } else {
        ElMessage.error(result.message? result.message : '注册失败');
      }
    } else {
      console.log('表单验证失败');
      return false;
    }
  });
}

// 登录
function SignIn() {
  formRefLogin.value.validate(async (valid) => {
    if (valid) {

      let result = await userLoginService(loginForm.value);
      if (result.code === 200) {
        // console.log(result)
        userStore.setUser({
          name: result.data.name,
          u_id: result.data.id,
          type: result.data.type
        })
        userStore.setToken(result.token)

        ElMessage.success(result.message? result.message : '登录成功');
        // UserManage.vue 或 login 页面中
        // localStorage.setItem('token', result.token);
        // localStorage.setItem('user', JSON.stringify({
        //   name: result.data.name, // 或者从 result.data 中获取更详细的用户信息
        //   u_id: result.data.id,
        //   type: result.data.type
        // }));


        // 跳转到主页
        // window.location.href = '/home';
        await router.push('/home')
      } else {
        ElMessage.error(result.message? result.message : '登录失败');
      }
    } else {
      console.log('表单验证失败');
      return false;
    }
  });
}

</script>

<template>
  <body>
    <div class="container" :class="{ 'right-panel-active': isSignUpActive }" id="login-box">
    <!-- 注册表单 -->
    <div class="form-container sign-up-container">
      <el-form :model="registerForm" label-position="top" label-width="120px" :rules="registerRules" ref="formRefRegister">
        <h1>注册</h1>
        <el-form-item label="用户名:" prop="username">
          <el-input v-model="registerForm.username"></el-input>
        </el-form-item>
        <el-form-item label="真实姓名:" prop="name">
          <el-input v-model="registerForm.name"></el-input>
        </el-form-item>
        <el-form-item label="密码:" prop="password">
          <el-input v-model="registerForm.password" show-password></el-input>
        </el-form-item>
        <el-form-item label="确认密码:" prop="confirmPassword">
          <el-input v-model="registerForm.confirmPassword" show-password></el-input>
        </el-form-item>
        <el-button type="primary" style="width: 100%" @click="SignUp">注册</el-button>
      </el-form>
    </div>

    <!-- 登录表单 -->
    <div class="form-container sign-in-container">
      <el-form :model="loginForm" label-position="top" label-width="120px" :rules="loginRules" ref="formRefLogin">
        <h1>登录</h1>
        <el-form-item label="用户名:" prop="username">
          <el-input v-model="loginForm.username" @keyup.enter="SignIn"></el-input>
        </el-form-item>
        <el-form-item label="密码:" prop="password">
          <el-input v-model="loginForm.password" show-password  @keyup.enter="SignIn"></el-input>
        </el-form-item>
        <a href="#">忘记密码？</a>
        <el-button type="primary" style="width: 100%" @click="SignIn">登录</el-button>
      </el-form>
    </div>

    <!-- 遮罩层 -->
    <div class="overlay-container">
      <div class="overlay">
        <div class="overlay-panel overlay-left">
          <h1>已有账号？</h1>
          <p>请使用您的账号进行登录</p>
          <button class="ghost" id="signIn" @click="handleSignIn">登录</button>
        </div>
        <div class="overlay-panel overlay-right">
          <h1>没有账号?</h1>
          <p>立即注册加入我们，和我们一起开始旅程吧</p>
          <button class="ghost" id="signUp" @click="handleSignUp">注册</button>
        </div>
      </div>
    </div>
  </div>
  </body>

</template>

<style scoped>
* {
  box-sizing: border-box;
}

body {
  font-family: 'Montserrat', sans-serif;
  //background-image: linear-gradient(120deg, #3498db, #8e44ad);;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100vh;
  margin: -20px 0 50px;
}

h1 {
  font-weight: bold;
  margin: 0;
}

p {
  font-size: 14px;
  line-height: 20px;
  letter-spacing: .5px;
  margin: 20px 0 30px;
}

span {
  font-size: 12px;
}

a {
  color: #333;
  font-size: 14px;
  text-decoration: none;
  margin: 15px 0;
}

.container {
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 14px 28px rgba(0, 0, 0, .25), 0 10px 10px rgba(0, 0, 0, .22);
  position: relative;
  overflow: hidden;
  width: 50%;
  max-width: 100%;
  min-height: 480px;

}

.form-container form {
  background: #fff;
  display: flex;
  flex-direction: column;
  padding: 0 50px;
  height: 100%;
  justify-content: center;
  text-align: center;
}

.social-container {
  margin: 20px 0;
}

.social-container a {
  border: 1px solid #ddd;
  border-radius: 50%;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  margin: 0 5px;
  height: 40px;
  width: 40px;
}

.social-container a:hover {
  background-color: #eee;

}

.txtb {
  border-bottom: 2px solid #adadad;
  position: relative;
  margin: 10px 0;
}

.txtb input {
  font-size: 15px;
  color: #333;
  border: none;
  width: 100%;
  outline: none;
  background: none;
  padding: 0 3px;
  height: 35px;
}

.txtb span::before {
  content: attr(data-placeholder);
  position: absolute;
  top: 50%;
  left: 5px;
  color: #adadad;
  transform: translateY(-50%);
  transition: .5s;
}

.txtb span::after {
  content: '';
  position: absolute;
  left: 0%;
  top: 100%;
  width: 0%;
  height: 2px;
  background-image: linear-gradient(120deg, #3498db, #8e44ad);
  transition: .5s;
}

.focus + span::before {
  top: -5px;
}

.focus + span::after {
  width: 100%;
}

button {
  border-radius: 20px;
  border: 1px solid #ff4b2b;
  background: #ff4b2b;
  color: #fff;
  font-size: 12px;
  font-weight: bold;
  padding: 12px 45px;
  letter-spacing: 1px;
  text-transform: uppercase;
  transition: transform 80ms ease-in;
  cursor: pointer;
}

button:active {
  transform: scale(.95);
}

button:focus {
  outline: none;
}

button.ghost {
  background: transparent;
  border-color: #fff;
}

.form-container {
  position: absolute;
  top: 0;
  height: 100%;
  transition: all .6s ease-in-out;
}

.form-container button {
  background: linear-gradient(120deg, #3498db, #8e44ad);
  border: none;
  background-size: 200%;
  color: #fff;
  transition: .5s;
}

.form-container button:hover {
  background-position: right;
}

.sign-in-container {
  left: 0;
  width: 50%;
  z-index: 2;
}

.sign-in-container form a {
  position: relative;
  left: 100px;
}

.sign-up-container {
  left: 0;
  width: 50%;
  z-index: 1;
  opacity: 0;
}

.sign-up-container button {
  margin-top: 20px;
}

.overlay-container {
  position: absolute;
  top: 0;
  left: 50%;
  width: 50%;
  height: 100%;
  overflow: hidden;
  transition: transform .6s ease-in-out;
  z-index: 100;
}

.overlay {
  background-image: linear-gradient(120deg, #3498db, #8e44ad);
  color: #fff;
  position: relative;
  left: -100%;
  height: 100%;
  width: 200%;
  transform: translateY(0);
  transition: transform .6s ease-in-out;
}

.overlay-panel {
  position: absolute;
  top: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 0 40px;
  height: 100%;
  width: 50%;
  text-align: center;
  transform: translateY(0);
  transition: transform .6s ease-in-out;
}

.overlay-right {
  right: 0;
  transform: translateY(0);

}

.overlay-left {
  transform: translateY(-20%);
}

.container.right-panel-active .sign-in-container {
  transform: translateY(100%);
}

.container.container.right-panel-active .overlay-container {
  transform: translateX(-100%);
}

.container.right-panel-active .sign-up-container {
  transform: translateX(100%);
  opacity: 1;
  z-index: 5;
}

.container.container.right-panel-active .overlay {
  transform: translateX(50%);
}

.container.container.right-panel-active .overlay-left {
  transform: translateY(0);
}

.container.container.right-panel-active .overlay-right {
  transform: translateY(20%);
}

</style>