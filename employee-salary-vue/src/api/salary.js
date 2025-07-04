import request from '@/utils/request.js'

const API_PATH = {
    USER_SALARY_LIST: '/userSalaryList',
    USER_MONTH_SALARY_LIST: '/userMonthSalaryList',
    DAILY_SALARY_DELETE: '/dailySalaryDelete',
    DAILY_SALARY_SPECIAL_DELETE: '/dailySalarySpecialDelete',
    DAILY_SALARY_DETAIL: '/dailySalaryDetail',
    DAILY_SALARY_SPECIAL_DETAIL: '/dailySalarySpecialDetail',
    DAILY_SALARY_UPDATE: '/dailySalaryUpdate',
    DAILY_SALARY_SPECIAL_UPDATE: '/dailySalarySpecialUpdate',
};


//salary列表查询
export const getUserSalaryListService = (params) => {
    return request.get(API_PATH.USER_SALARY_LIST,{params:params})
}

//monthSalary列表查询
export const getUserMonthSalaryListService = (params) => {
    return request.get(API_PATH.USER_MONTH_SALARY_LIST,{params:params})
}

//clothing删除
export const deleteDailySalaryService = (id) => {
    return request.post(API_PATH.DAILY_SALARY_DELETE, { id })
}

//clothing删除
export const deleteDailySalarySpecialService = (id) => {
    return request.post(API_PATH.DAILY_SALARY_SPECIAL_DELETE, { id })
}

//获取服饰具体详情
export const getDailySalaryDetailService = (params) => {
    return request.get(API_PATH.DAILY_SALARY_DETAIL, {params:params})
}

//获取特殊具体详情
export const getDailySalarySpecialDetailService = (params) => {
    return request.get(API_PATH.DAILY_SALARY_SPECIAL_DETAIL, {params:params})
}

//clothing修改
export const updateDailySalaryService = (data) => {
    return request.post(API_PATH.DAILY_SALARY_UPDATE,data)
}

//特殊修改
export const updateDailySalarySpecialService = (data) => {
    return request.post(API_PATH.DAILY_SALARY_SPECIAL_UPDATE,data)
}



