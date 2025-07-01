import request from '@/utils/request.js'

const API_PATH = {
    CLOTHING_LIST: '/clothingList',
    CLOTHING_ADD: '/clothingAdd',
    CLOTHING_UPDATE: '/clothingUpdate',
    CLOTHING_DELETE: '/clothingDelete',
    CLOTHING_DETAIL: '/clothingDetail',
    CLOTHING_SALARY_LIST: '/salaryClothingList',
    CLOTHING_IMAGE_UPLOAD: '/clothingImageUpload',
    CLOTHING_IMAGE_DELETE: '/clothingImageDelete'
};


//clothing列表查询
export const getClothingListService = (params) => {
    return request.get(API_PATH.CLOTHING_LIST,{params:params})
}

//clothing添加
export const addClothingService = (data) => {
    return request.post(API_PATH.CLOTHING_ADD,data)
}

//clothing修改
export const updateClothingService = (data) => {
    return request.post(API_PATH.CLOTHING_UPDATE,data)
}

//clothing删除
export const deleteClothingService = (id) => {
    return request.post(API_PATH.CLOTHING_DELETE, { id })
}

//获取服饰具体详情
export const getClothingDetailService = (params) => {
    return request.get(API_PATH.CLOTHING_DETAIL, {params:params})
}

//获取所有员工列表
export const getSalaryClothingService = (params) => {
    return request.get(API_PATH.CLOTHING_SALARY_LIST,{params:params})
}

//上传服饰图片
// export const uploadClothingService = (data) => {
//     return request.post(API_PATH.CLOTHING_IMAGE_UPLOAD,data)
// }

//删除服饰图片
export const deleteClothingImageService = (data) => {
    return request.post(API_PATH.CLOTHING_IMAGE_DELETE,data)
}


