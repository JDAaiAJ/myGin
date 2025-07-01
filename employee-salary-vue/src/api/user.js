import request from '@/utils/request.js'

const API_PATH = {
    USER_LOGIN: '/login',
    USER_REGISTER: '/register',
    USER_LIST: '/userList',
    USER_DETAIL: '/userDetail',
    USER_UPDATE: '/userUpdate',
    USERS_BY_FACTORY: '/usersByFactory',
    USER_FACTORY_LIST: '/userFactoryList',
    USER_DELETE: '/userDelete'

};

//用户登录
export const userLoginService = (data) => {
    return request.post(API_PATH.USER_LOGIN, data)
}

//用户注册
export const userRegisterService = (data) => {
    return request.post(API_PATH.USER_REGISTER, data)
};

//用户列表查询
export const getUserListService = (params) => {
    return request.get(API_PATH.USER_LIST, {params: params})
}

//获取所有员工列表
export const getUsersBYFactoryService = (params) => {
    return request.get(API_PATH.USERS_BY_FACTORY, {params: params})
}

//获取所有员工列表
export const getUserFactoryService = (params) => {
    return request.get(API_PATH.USER_FACTORY_LIST, {params: params})
}

//获取用户具体详情
export const getUserDetailService = (params) => {
    return request.get(API_PATH.USER_DETAIL, {params: params})
}

//用户信息修改
export const updateUserService = (data) => {
    return request.post(API_PATH.USER_UPDATE,data)
}

//用户信息删除
export const deleteUserService = (id) => {
    return request.post(API_PATH.USER_DELETE, { id })
}


