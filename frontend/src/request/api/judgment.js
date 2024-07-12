import serviceAxios from "../index";

// ============================ post ============================

    // 获取开放题列表
    export const getUserOpenQuestList = (data) => {
        return serviceAxios({
            url: `/challenge/getUserOpenQuestListV2`,
            method: "post",
            data
        })
    }

    // 获取开放题详情
    export const getUserOpenQuestDetailList = (data) => {
        return serviceAxios({
            url: `/challenge/getUserOpenQuestDetailListV2`,
            method: "post",
            data
        })
    }

    // 开放题打分
    export const reviewOpenQuest = (data) => {
        return serviceAxios({
            url: `/challenge/reviewOpenQuestV2`,
            method: "post",
            data
        })
    }

    // 获取用户答题统计
    export const getUserQuestDetail = (data) => {
        return serviceAxios({
            url: `/challenge/getUserQuestDetail`,
            method: "post",
            data
        })
    }
