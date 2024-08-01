import { useEffect, useState } from "react"
import { Button, Input, InputNumber, Slider, message } from "antd";
import { download } from "../../utils/file/download";
import { getUserOpenQuestDetailList, reviewOpenQuest } from "../../request/api/judgment";
import ReactMarkdown from 'react-markdown';
const { TextArea } = Input;


function ChallengeJudgPage({questDetail, reviewStatus, hideModal, updateList}) {

    const [index, setIndex] = useState(0);      // 第几题
    const [total, setTotal] = useState(0);
    const [isLoding, setTsLoding] = useState(false);
    
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
        </div>
    )
}

export default ChallengeJudgPage