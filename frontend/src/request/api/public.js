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

    // 获取youtube视频解析
    export const getYouTubePlayList = (data) => {
        return serviceAxios({
            url: `/video/getYouTubePlayList`,
            method: "post",
            data
        })
    }


// ============================ get ============================

    // 获取验证码
    export const userCaptcha = (data) => {
        return serviceAxios({
            url: `/user/captcha`,
            method: "get",
            data
        })
    }
