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

    // 上传图片
export const ipfsImg = (data) => {
    return serviceAxios({
        url: `/ipfs/uploadFile`,
        method: "post",
        data
    })
}

    // 上传教程
export const createTutorial = (data) => {
    return serviceAxios({
        url: `/tutorial/createTutorial`,
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