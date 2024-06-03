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

    // 删除用户标签
    export const tagDeleteBatch = (data) => {
        return serviceAxios({
            url: `/tag/tagDeleteBatch`,
            method: "post",
            data
        })
    }

