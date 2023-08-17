import serviceAxios from "../index";

// 操作 ==================>

// 上传教程
export const createTutorial = (data) => {
    return serviceAxios({
        url: `/tutorial/createTutorial`,
        method: "post",
        data
    })
}

// 修改教程
export const updateTutorial = (data) => {
    return serviceAxios({
        url: `/tutorial/updateTutorial`,
        method: "post",
        data
    })
}

// 删除教程
export const deleteTutorial = (data) => {
    return serviceAxios({
        url: `/tutorial/deleteTutorial`,
        method: "post",
        data
    })
}

// 获取 ==================>

//  获取教程列表
export const getTutorialList = (data) => {
    return serviceAxios({
        url: `/tutorial/getTutorialList`,
        method: "post",
        data
    })
}

// 获取教程详情
export const getTutorial = (data) => {
    return serviceAxios({
        url: `/tutorial/getTutorial`,
        method: "post",
        data
    })
}

// 获取打包列表
export const getPackList = (data) => {
    return serviceAxios({
        url: `/pack/getPackList`,
        method: "post",
        data
    })
}