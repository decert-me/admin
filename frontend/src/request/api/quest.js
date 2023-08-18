import questAxios from "../quest";


// ============================ get ============================
    //     获取挑战详情
    export const getQuest = ({id}) => {
        return questAxios({
            url: `/quests/${id}`,
            method: "get"
        })
    }
