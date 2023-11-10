import serviceAxios from "../index";

// ============================ get ============================

    // 获取登陆签名
    export const getLoginMessage = (data) => {
        return serviceAxios({
            url: `/user/getLoginMessage?address=${data.address}`,
            method: "get"
        })
    }


// ============================ post ============================

    // 校验签名
    export const authLoginSign = (data) => {
        return serviceAxios({
            url: `/user/authLoginSign`,
            method: "post",
            data
        })
    }