import request from '@/utils/request.js'

const API_PATH = {
    FACTORY_LIST: '/factoryList',
    FACTORY_ADD: '/factoryAdd',
    FACTORY_UPDATE: '/factoryUpdate',
    FACTORY_DELETE: '/factoryDelete',
};


//factory列表查询
export const getFactoryListService = (params) => {
    return request.get(API_PATH.FACTORY_LIST,{params:params})
}

//factory添加
export const addFactoryService = (data) => {
    return request.post(API_PATH.FACTORY_ADD,data)
}


//factory修改
export const updateFactoryService = (data) => {
    return request.post(API_PATH.FACTORY_UPDATE,data)
}

//factory删除
export const deleteFactoryService = (id) => {
    return request.post(API_PATH.FACTORY_DELETE, { id })
}



