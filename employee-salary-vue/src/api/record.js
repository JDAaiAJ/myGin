import request from '@/utils/request.js'

const API_PATH = {
    USER_SALARY_ADD: '/dailySalaryAdd',
    USER_SALARY_SPECIAL_ADD: '/dailySalarySpecialAdd',
};



//添加用户每日薪资
export const addUserRecordService = (data) => {
    return request.post(API_PATH.USER_SALARY_ADD,data)
}

//添加用户每日特殊薪资
export const addUserRecordSpecialService = (data) => {
    return request.post(API_PATH.USER_SALARY_SPECIAL_ADD,data)
}


