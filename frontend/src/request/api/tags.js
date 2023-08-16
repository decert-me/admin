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
    
    
// ============================ get ============================