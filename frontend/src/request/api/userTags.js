import serviceAxios from "../index";

// ============================ get ============================

    // 获取登陆签名
    // export const getLoginMessage = (data) => {
    //     return serviceAxios({
    //         url: `/tag/tagAdd`,
    //         method: "get"
    //     })
    // }


// ============================ post ============================

    // 获取用户标签列表
    export const getTagList = (data) => {
        return serviceAxios({
            url: `/tag/getTagList`,
            method: "post",
            data
        })
    }

    // 获取用户标签详情
    export const getTagInfo = (data) => {
        return serviceAxios({
            url: `/tag/getTagInfo`,
            method: "post",
            data
        })
    }

    // 添加用户标签
    export const addUserTag = (data) => {
        return serviceAxios({
            url: `/tag/tagAdd`,
            method: "post",
            data
        })
    }

    // 修改用户标签
    export const tagUpdate = (data) => {
        return serviceAxios({
            url: `/tag/tagUpdate`,
            method: "post",
            data
        })
    }

    // 获取标签用户列表
    export const getTagUserList = (data) => {
        return serviceAxios({
            url: `/tag/getTagUserList`,
            method: "post",
            data
        })
    }
    

    // 批量更新用户标签
    export const updateUsersInfo = (data) => {
        return serviceAxios({
            url: `/users/updateUsersInfo`,
            method: "post",
            data
        })
    }

    // 删除用户标签
    export const tagDeleteBatch = (data) => {
        return serviceAxios({
            url: `/tag/tagDeleteBatch`,
            method: "post",
            data
        })
    }

    // 批量删除用户标签
    export const tagUserDeleteBatch = (data) => {
        return serviceAxios({
            url: `/tag/tagUserDeleteBatch`,
            method: "post",
            data
        })
    }

    // 动态搜索地址
    export const getUsersList = (data) => {
        return serviceAxios({
            url: `/users/getUsersList`,
            method: "post",
            data
        })
    }

    // 用户添加标签
    export const tagUserUpdate = (data) => {
        return serviceAxios({
            url: `/tag/tagUserUpdate`,
            method: "post",
            data
        })
    }

    // 获取用户所有标签
    export const getUsersInfo = (data) => {
        return serviceAxios({
            url: `/users/getUsersInfo`,
            method: "post",
            data
        })
    }