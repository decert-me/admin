import serviceAxios from "../index";

// ============================ post ============================
    // 新建标签
    export const createLabel = (data) => {
        return serviceAxios({
            url: `/label/createLabel`,
            method: "post",
            data
        })
    }
    // 删除标签
    export const deleteLabel = (data) => {
        return serviceAxios({
            url: `/label/deleteLabel`,
            method: "post",
            data
        })
    }
    // 修改标签
    export const updateLabel = (data) => {
        return serviceAxios({
            url: `/label/updateLabel`,
            method: "post",
            data
        })
    }
    
// ============================ get ============================
    //     获取标签列表
    export const getLabelList = (data) => {
        return serviceAxios({
            url: `/label/getLabelList`,
            method: "post",
            data
        })
    }