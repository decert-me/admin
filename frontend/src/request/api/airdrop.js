import serviceAxios from "../index";

// ============================ post ============================

    // 获取空投列表
    export const getAirdropList = (data) => {
        return serviceAxios({
            url: `/airdrop/getAirdropList`,
            method: "post",
            data
        })
    }

    // 立即空投
    export const runAirdrop = (data) => {
        return serviceAxios({
            url: `/airdrop/runAirdrop`,
            method: "post",
            data
        })
    }