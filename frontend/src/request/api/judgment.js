import serviceAxios from "../index";

// ============================ post ============================

    // 获取开放题列表
    export const getUserOpenQuestList = (data) => {
        return serviceAxios({
            url: `/challenge/getUserOpenQuestList`,
            method: "post",
            data
        })
    }

    // 开放题打分
    export const reviewOpenQuest = (data) => {
        return serviceAxios({
            url: `/challenge/reviewOpenQuest`,
            method: "post",
            data
        })
    }
    
