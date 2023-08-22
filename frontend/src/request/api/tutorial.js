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

    // 打包
    export const buildTutorial = (data) => {
        return serviceAxios({
            url: `/pack/pack`,
            method: "post",
            data
        })
    }

    // 修改教程上架状态
    export const updateTutorialStatus = (data) => {
        return serviceAxios({
            url: `/tutorial/updateTutorialStatus`,
            method: "post",
            data
        })
    }

    // 修改置顶状态
    export const topTutorial = (data) => {
        return serviceAxios({
            url: `/tutorial/topTutorial`,
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

    // 获取打包日志
    export const getPackLog = (data) => {
        return serviceAxios({
            url: `/pack/getPackLog`,
            method: "post",
            data
        })
    }