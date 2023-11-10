import serviceAxios from "../index";

// ============================ get ============================

    // 获取用户列表
    export const getUserList = (data) => {
        return serviceAxios({
            url: `/user/list?page=${data.page}&pagesize=${data.pageSize}`,
            method: "get"
        })
    }

    // 获取用户详情
    export const getUserInfo = (data) => {
        return serviceAxios({
            url: `/user/info?id=${data.id}`,
            method: "get"
        })
    }

// ============================ post ============================

    // 上传头像
    export const uploadAvatar = (data) => {
        return serviceAxios({
            url: `/upload/avatar`,
            method: "post",
            data
        })
    }

    // 更新个人资料
    export const updatePersonalInfo = (data) => {
        return serviceAxios({
            url: `/user/update`,
            method: "post",
            data
        })
    }
    
    // 更新他人资料
    export const updateOtherInfo = (data) => {
        return serviceAxios({
            url: `/user/update`,
            method: "post",
            data
        })
    }

    // 创建管理员
    export const registerUser = (data) => {
        return serviceAxios({
            url: `/user/register`,
            method: "post",
            data
        })
    }

    // 删除管理员
    export const deleteUser = (data) => {
        return serviceAxios({
            url: `/user/delete`,
            method: "post",
            data
        })
    }