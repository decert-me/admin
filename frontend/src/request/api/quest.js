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

    // 获取合辑列表
    export const getCollectionList = (data) => {
        return serviceAxios({
            url: `/collection/list`,
            method: "post",
            data
        })
    }

    // 获取合辑详情
    export const getCollectionDetail = (data) => {
        return serviceAxios({
            url: `/collection/detail`,
            method: "post",
            data
        })
    }

    // 获取合辑内的挑战列表
    export const getCollectionQuestList = (data) => {
        return serviceAxios({
            url: `/collection/collectionQuest`,
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

    // 删除挑战
    export const deleteQuest = (data) => {
        return serviceAxios({
            url: `/quest/delete`,
            method: "post",
            data
        })
    }

    // 更新挑战
    export const updateQuest = (data) => {
        return serviceAxios({
            url: `/quest/update`,
            method: "post",
            data
        })
    }
    
    // 创建合辑
    export const createCollection = (data) => {
        return serviceAxios({
            url: `/collection/create`,
            method: "post",
            data
        })
    }

    // 更新合辑上下架状态
    export const updateStatusCollection = (data) => {
        return serviceAxios({
            url: `/collection/updateStatus`,
            method: "post",
            data
        })
    }

    // 更新合辑
    export const updateCollection = (data) => {
        return serviceAxios({
            url: `/collection/update`,
            method: "post",
            data
        })
    }

    // 删除合辑
    export const deleteCollection = (data) => {
        return serviceAxios({
            url: `/collection/delete`,
            method: "post",
            data
        })
    }

    // 修改挑战合辑内的排序
    export const updateCollectionQuestSort = (data) => {
        return serviceAxios({
            url: `/collection/updateCollectionQuestSort`,
            method: "post",
            data
        })
    }