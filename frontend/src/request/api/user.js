import serviceAxios from "../index";

// ============================ post ============================
// 登录
export const userLogin = (data) => {
    return serviceAxios({
        url: `/user/login`,
        method: "post",
        data
    })
}

// ============================ get ============================
export const userCaptcha = (data) => {
    return serviceAxios({
        url: `/user/captcha`,
        method: "get",
        data
    })
}