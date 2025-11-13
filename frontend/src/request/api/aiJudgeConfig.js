import service from '../index'

// 获取AI判题配置列表
export function getAiJudgeConfigList() {
    return service({
        url: '/ai-judge-config/list',
        method: 'get',
    })
}

// 创建AI判题配置
export function createAiJudgeConfig(data) {
    return service({
        url: '/ai-judge-config/create',
        method: 'post',
        data
    })
}

// 更新AI判题配置
export function updateAiJudgeConfig(data) {
    return service({
        url: '/ai-judge-config/update',
        method: 'post',
        data
    })
}

// 删除AI判题配置
export function deleteAiJudgeConfig(data) {
    return service({
        url: '/ai-judge-config/delete',
        method: 'post',
        data
    })
}

// 切换AI判题配置启用状态
export function toggleAiJudgeConfig(data) {
    return service({
        url: '/ai-judge-config/toggle',
        method: 'post',
        data
    })
}

// 获取当前启用的AI判题配置
export function getEnabledAiJudgeConfig() {
    return service({
        url: '/ai-judge-config/enabled',
        method: 'get',
    })
}

// AI判题
export function aiGrade(data) {
    return service({
        url: '/ai-judge-config/grade',
        method: 'post',
        data
    })
}

// 切换自动判题状态
export function toggleAutoGrading(data) {
    return service({
        url: '/ai-judge-config/toggle-auto-grading',
        method: 'post',
        data
    })
}

// 获取待判题列表
export function getPendingGradeList(params) {
    return service({
        url: '/ai-judge-config/pending-list',
        method: 'get',
        params
    })
}

// 批量AI判题
export function batchGrade(data) {
    return service({
        url: '/ai-judge-config/batch-grade',
        method: 'post',
        data
    })
}

// 获取AI判题历史
export function getAiGradeHistory(params) {
    return service({
        url: '/ai-judge-config/history',
        method: 'get',
        params
    })
}

// 批量AI判题预览（不直接提交）
export function batchGradePreview(data) {
    return service({
        url: '/ai-judge-config/batch-grade-preview',
        method: 'post',
        data
    })
}

// 提交批量判题结果
export function submitBatchGrade(data) {
    return service({
        url: '/ai-judge-config/submit-batch-grade',
        method: 'post',
        data
    })
}
