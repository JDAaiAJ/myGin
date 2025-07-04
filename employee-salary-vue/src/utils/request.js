import axios from 'axios';
import {ElMessage} from "element-plus";
import router from "../router";

// 引入 Pinia Store
import { useUserStore } from '@/stores/useUserStore'

// 创建一个 axios 实例
const service = axios.create({
    baseURL: '/api/', // 所有请求都会带上 `/api` 前缀，由 vite 代理到后端
    timeout: 5000, // 请求超时时间
    headers: {
        'Content-Type': 'application/json;charset=utf-8',
    },
});

// 请求拦截器
service.interceptors.request.use(config => {
    // const token = localStorage.getItem('token')
    const userStore = useUserStore()
    const token = userStore.token // 从 Pinia 获取 token

    if (token) {
        config.headers['Authorization'] = 'Bearer ' + token
    }
    return config
}, error => {
    return Promise.reject(error)
})

// 响应拦截器
service.interceptors.response.use(
    (response) => {
        // 处理 HTTP 状态码 2xx
        return response.data;
    },
    (error) => {
        // 处理网络错误或 HTTP 错误（如 401, 404, 500）
        console.error('请求异常:', error);

        if (error.response) {
            // 服务器响应了但状态码不是 2xx
            const { status, data } = error.response;

            let message = '请求异常';

            switch (status) {
                case 400:
                    message = '请求参数错误';
                    break;
                case 401:
                    message = '认证失败，请重新登录，等待跳转登录页..';
                    //  清除 token 并跳转登录页 两秒后跳转
                    setTimeout(() => {
                        localStorage.removeItem('token');
                        router.push('/login');
                    }, 1200);
                    break;
                case 403:
                    message = '权限不足';
                    break;
                case 404:
                    message = '请求资源不存在';
                    break;
                case 500:
                    message = '服务器内部错误';
                    break;
                default:
                    message = data.msg || `服务器返回错误: ${status}`;
            }

            ElMessage.error(message);
        } else if (error.request) {
            // 请求已发出但没有收到响应（比如网络中断）
            ElMessage.error('网络连接异常，请检查网络');
        } else {
            // 其他错误
            ElMessage.error(error.message);
        }

        return Promise.reject(error);
    }
);

export default service;
