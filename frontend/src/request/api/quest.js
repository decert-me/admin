import questAxios from "../quest";
import serviceAxios from "../index";


// ============================ get ============================

    //     获取挑战详情
    export const getQuest = ({id}) => {
        return questAxios({
            url: `/quests/${id}`,
            method: "get"
        })
    }

    // 获取挑战列表
    export const getQuestList = (data) => {
        return serviceAxios({
            url: `/quest/list`,
            method: "post",
            data
        })
    }

// ============================ post ============================

    // 更改挑战上架状态
    export const updateQuestStatus = (data) => {
        return serviceAxios({
            url: `/quest/updateQuestStatus`,
            method: "post",
            data
        })
    }

    // 置顶、取消置顶挑战
    export const topQuest = (data) => {
        return serviceAxios({
            url: `/quest/topQuest`,
            method: "post",
            data
        })
    }