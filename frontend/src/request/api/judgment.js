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
