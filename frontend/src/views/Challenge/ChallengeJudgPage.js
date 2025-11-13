import { useEffect, useState } from "react"
import { Button, Input, InputNumber, Slider, message, Modal } from "antd";
import { download } from "../../utils/file/download";
import { getUserOpenQuestDetailList, reviewOpenQuest } from "../../request/api/judgment";
import { aiGrade } from "../../request/api/aiJudgeConfig";
import ReactMarkdown from 'react-markdown';
const { TextArea } = Input;


function ChallengeJudgPage({questDetail, reviewStatus, hideModal, updateList}) {

    const [index, setIndex] = useState(0);      // 第几题
    const [total, setTotal] = useState(0);
    const [isLoding, setTsLoding] = useState(false);
    const [aiGrading, setAiGrading] = useState(false); // AI判题中的loading状态
    const [debugModalOpen, setDebugModalOpen] = useState(false); // AI调试弹窗
    const [aiDebugInfo, setAiDebugInfo] = useState({ // AI对话信息
        request: '',
        response: '',
        systemPrompt: '',
        userPrompt: ''
    });

    let [reviewQuests, setReviewQuests] = useState([]);
    let [openQsList, setOpenQsList] = useState([]);
    let [selectOpenQs, setSelectOpenQs] = useState({});
    let [page, setPage] = useState(0);
    let [rateCache, setRateCache] = useState({
        rate: 0,
        annotation: ""
    });

    async function confirm() {
        // 没有改分直接退出
        const list = reviewQuests.filter(e => e);
        if (list.length === 0) {
            hideModal();
            return
        }
        await reviewOpenQuest(list)
        .then(res => {
            if (res.code === 0) {
                message.success(res.msg);
                hideModal();
                updateList();
            }
        })
        .catch(err => {
            message.error(err.msg);
        })
    }

    function setAnnotation(text) {
        rateCache.annotation = text;
        setRateCache({...rateCache});

        updateCache();
    }

    function setPercent(percent) {
        // 将rateCache写入数组中
        rateCache.score = percent;
        setRateCache({...rateCache});

        updateCache();
    }

    function updateCache() {
        reviewQuests[index - 1] = {
            id: selectOpenQs.ID,
            answer: {
                type: "open_quest",
                annex: selectOpenQs.answer.annex,
                value: selectOpenQs.answer.value,
                score: rateCache.score,
                annotation: rateCache.annotation,
                open_quest_review_time: new Date()
                    .toLocaleString()
                    .replace(/\//g, "-"),
            },
            index: selectOpenQs.index,
            updated_at: selectOpenQs.updated_at,
        }
        setReviewQuests([...reviewQuests]);
    }

    // AI判题
    async function handleAiGrade() {
        if (!selectOpenQs?.title) {
            message.warning('题目信息不完整');
            return;
        }

        // 检查是否有答案或附件
        const hasAnswer = selectOpenQs?.answer?.value && selectOpenQs.answer.value.trim();
        const hasAttachment = selectOpenQs?.answer?.annex && selectOpenQs.answer.annex.length > 0;

        if (!hasAnswer && !hasAttachment) {
            message.warning('用户未提供答案或附件，将判定为不通过');
        }

        setAiGrading(true);
        message.loading({ content: 'AI判题中...', key: 'aiGrading', duration: 0 });

        // 构建附件URL列表
        let attachmentUrls = [];
        if (selectOpenQs?.answer?.annex && selectOpenQs.answer.annex.length > 0) {
            attachmentUrls = selectOpenQs.answer.annex.map(annex => {
                // 根据hash构建附件URL
                const baseUrl = process.env.REACT_APP_IPFS_URL || 'https://ipfs.decert.me';
                return `${baseUrl}/ipfs/${annex.hash} (文件名: ${annex.name})`;
            });
        }

        try {
            const requestData = {
                question_title: selectOpenQs.title,
                question_score: selectOpenQs.score,
                pass_score: selectOpenQs.pass_score,
                user_answer: selectOpenQs.answer.value,
                attachment_urls: attachmentUrls
            };

            const res = await aiGrade(requestData);

            if (res.code === 0) {
                message.destroy('aiGrading');
                message.success('AI判题完成');

                // 设置分数和批注
                const { score, annotation, raw_result, system_prompt, user_prompt } = res.data;

                // 保存调试信息
                setAiDebugInfo({
                    request: JSON.stringify(requestData, null, 2),
                    response: raw_result || annotation,
                    systemPrompt: system_prompt || '',
                    userPrompt: user_prompt || ''
                });

                // 显示调试弹窗
                setDebugModalOpen(true);

                rateCache.score = score;
                rateCache.annotation = annotation;
                setRateCache({...rateCache});
                updateCache();
            } else {
                message.destroy('aiGrading');
                message.error(res.msg || 'AI判题失败');
            }
        } catch (error) {
            message.destroy('aiGrading');
            message.error('AI判题失败：' + (error.message || '未知错误'));
        } finally {
            setAiGrading(false);
        }
    }

    async function init() {
        rateCache = {
            rate: 0,
            annotation: ""
        }
        setRateCache({...rateCache});
        reviewQuests = [];
        setReviewQuests([...reviewQuests]);
        openQsList = [];
        setOpenQsList([...openQsList]);
        selectOpenQs = {};
        setSelectOpenQs({...selectOpenQs});
        setTotal(0);
        changePage(1, true);

    }

    // 切换上下题
    function changeIndex(index) {
        setTsLoding(true);
        // 评分模式从reviewlist读取缓存
        if (reviewStatus) {
            rateCache = {
                score: reviewQuests[index - 1]?.answer.score || 0,
                annotation: reviewQuests[index - 1]?.answer.annotation || ""
            }
            setRateCache({...rateCache});
        }
        if (index > openQsList.length) {
            changePage(page+1);
            return
        }
        setIndex(index);
        selectOpenQs = openQsList[index-1];
        setSelectOpenQs({...selectOpenQs});
        setTsLoding(false);
    }

    // 切换上下页
    async function changePage(newPage, isInit) {
        if (newPage) {
            page = newPage;
            setPage(page);
        }
        getUserOpenQuestDetailList({
            "page": page,
            "pageSize": 50,
            "open_quest_review_status": reviewStatus ? 1 : 2,
            ...questDetail
        })
        .then(res => {
            if (res.code === 0) {
                const list = res.data.list || [];
                openQsList = openQsList.concat(list);
                setOpenQsList([...openQsList]);
                if (openQsList.length === 0) {
                    setIndex(0);
                    return
                }
                changeIndex(isInit ? 1 : index+1);
                // 获取总页
                setTotal(res.data.total);
                // 初始化评分列表
                if (reviewQuests.length === 0) {
                    reviewQuests = new Array(res.data.total);
                    setReviewQuests([...reviewQuests]);
                } 
            }
        })
        .catch(err => {
            message.error(err.msg);
        })
    }

    useEffect(() => {
        questDetail && init();
    },[questDetail])

    return (
        <div className="judg-content">
            <h1>{selectOpenQs?.challenge_title}</h1>
                <div className="judg-info">
                    <div className="item">
                        <div className="item-title">第 <strong>{selectOpenQs?.submit_count}</strong> 次提交</div>
                    </div>
                    <div className="item">
                        <div className="item-title">挑战者: &nbsp;<a href={`${process.env.REACT_APP_LINK_URL || "https://decert.me"}/user/${selectOpenQs?.address}`} target="_blank" rel="noopener noreferrer">{selectOpenQs?.nickname}</a></div>
                    </div>

                    <div className="item">
                        <div className="item-title">提交时间: &nbsp;
                            <span className="item-content">{
                            selectOpenQs?.created_at && selectOpenQs?.created_at.indexOf("0001-01-01T") === -1
                            ? selectOpenQs?.created_at.replace("T", " ").split(".")[0].split("+")[0]
                            : "-"}</span>
                        </div>
                    </div>

                    <div className="item">
                        <div className="item-title">题目: &nbsp;
                            <span className="item-content">{selectOpenQs?.title}</span>
                        </div>
                        {/* <div className="item-content">
                            <ReactMarkdown>{selectOpenQs?.challenge_title}</ReactMarkdown>
                        </div> */}
                    </div>

                    <div className="item">
                        <p className="item-title">答案:</p>
                        <TextArea 
                            className="item-content box"
                            bordered={false} 
                            autoSize={{
                                minRows: 3,
                                maxRows: 5,
                            }}
                            readOnly
                            value={selectOpenQs?.answer?.value}
                        />
                    </div>

                    <div className="item">
                        <p className="item-title">批注:</p>
                        <TextArea 
                            disabled={!reviewStatus}
                            className="item-content box"
                            bordered={false} 
                            autoSize={{
                                minRows: 3,
                                maxRows: 5,
                            }}
                            onChange={(e) => setAnnotation(e.target.value)}
                            value={selectOpenQs?.answer?.annotation || rateCache.annotation}
                        />
                    </div>

                    <div className="item">
                        <p className="item-title">附件:</p>
                        <div className="item-content">
                            {
                                selectOpenQs?.answer?.annex && selectOpenQs?.answer?.annex.map(e => (
                                    <Button type="link" key={e.name} onClick={() => download(e.hash, e.name)}>{e.name}</Button>
                                ))
                            }
                        </div>
                    </div>

                    <div className="item">
                        <div className="item-title">挑战总得分: &nbsp;
                            {/* <span className="item-content">{selectOpenQs?.total_score}</span> */}
                            <span className="item-content">{(selectOpenQs?.answer?.score ? selectOpenQs?.user_score : rateCache?.score ? selectOpenQs?.user_score + rateCache.score : "")}</span>
                        </div>
                    </div>

                    <div className="item">
                        <div className="item-title">挑战及格分: &nbsp;
                            <span className="item-content">{selectOpenQs?.pass_score}</span>
                        </div>
                    </div>

                    <div className="item">
                        <div className="item-title">本题评分: &nbsp;
                            {/* <span className="item-content">{selectOpenQs?.total_score}</span> */}
                            <InputNumber
                                disabled={!reviewStatus}
                                min={0}
                                max={selectOpenQs?.score}
                                style={{margin: '0 16px'}}
                                step={1}
                                value={selectOpenQs?.answer?.score ? selectOpenQs.answer?.score : rateCache?.score ? rateCache.score : ""}
                                onChange={(value) => setPercent(value)}
                            />
                            {reviewStatus && (
                                <Button
                                    type="primary"
                                    size="small"
                                    onClick={handleAiGrade}
                                    loading={aiGrading}
                                    disabled={aiGrading}
                                    style={{marginLeft: 8}}
                                >
                                    AI判题
                                </Button>
                            )}
                        </div>
                        <div className="item-title">上次评分: &nbsp;
                            <span className="item-content">{selectOpenQs?.last_score}</span>
                        </div>
                        <div style={{width: "352px"}}>
                            <Slider
                                disabled={!reviewStatus}
                                max={selectOpenQs?.score}
                                step={1}
                                tooltip={{formatter: null}}
                                value={selectOpenQs?.answer?.score ? selectOpenQs.answer?.score : rateCache?.score ? rateCache.score : 0}
                                onChange={(percent) => setPercent(percent)}
                            />
                        </div>
                    </div>
                </div>
                <div className="flex-bte">
                    <div className="pagination">
                        <Button disabled={index <= 1} onClick={() => changeIndex(index - 1)}>上一题</Button>
                        <p>{index}/<span style={{color: "#8B8D97"}}>{total}</span></p>
                        <Button loading={isLoding} disabled={index === total} onClick={() => changeIndex(index + 1)}>下一题</Button>
                    </div>
                    {
                        reviewStatus &&
                        <Button className="submit" type="primary" size="large" onClick={confirm}>提交</Button>
                    }
                </div>

                {/* AI调试弹窗 */}
                <Modal
                    title="AI判题详情"
                    open={debugModalOpen}
                    onCancel={() => setDebugModalOpen(false)}
                    width={900}
                    footer={[
                        <Button key="close" type="primary" onClick={() => setDebugModalOpen(false)}>
                            关闭
                        </Button>
                    ]}
                >
                    <div style={{ maxHeight: '70vh', overflow: 'auto' }}>
                        <div style={{ marginBottom: 20 }}>
                            <h3 style={{ marginBottom: 10, color: '#1890ff' }}>📤 发送给AI的请求数据</h3>
                            <pre style={{
                                background: '#f5f5f5',
                                padding: 15,
                                borderRadius: 4,
                                fontSize: 12,
                                whiteSpace: 'pre-wrap',
                                wordWrap: 'break-word'
                            }}>
                                {aiDebugInfo.request}
                            </pre>
                        </div>

                        <div style={{ marginBottom: 20 }}>
                            <h3 style={{ marginBottom: 10, color: '#52c41a' }}>🤖 System Prompt（系统提示词）</h3>
                            <pre style={{
                                background: '#f6ffed',
                                padding: 15,
                                borderRadius: 4,
                                border: '1px solid #b7eb8f',
                                fontSize: 12,
                                whiteSpace: 'pre-wrap',
                                wordWrap: 'break-word'
                            }}>
                                {aiDebugInfo.systemPrompt}
                            </pre>
                        </div>

                        <div style={{ marginBottom: 20 }}>
                            <h3 style={{ marginBottom: 10, color: '#faad14' }}>💬 User Prompt（用户提示词）</h3>
                            <pre style={{
                                background: '#fffbe6',
                                padding: 15,
                                borderRadius: 4,
                                border: '1px solid #ffe58f',
                                fontSize: 12,
                                whiteSpace: 'pre-wrap',
                                wordWrap: 'break-word'
                            }}>
                                {aiDebugInfo.userPrompt}
                            </pre>
                        </div>

                        <div>
                            <h3 style={{ marginBottom: 10, color: '#f5222d' }}>📥 AI返回的原始结果</h3>
                            <pre style={{
                                background: '#fff1f0',
                                padding: 15,
                                borderRadius: 4,
                                border: '1px solid #ffccc7',
                                fontSize: 12,
                                whiteSpace: 'pre-wrap',
                                wordWrap: 'break-word'
                            }}>
                                {aiDebugInfo.response}
                            </pre>
                        </div>
                    </div>
                </Modal>
        </div>
    )
}

export default ChallengeJudgPage