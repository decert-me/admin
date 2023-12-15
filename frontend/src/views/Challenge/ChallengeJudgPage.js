import { forwardRef, useEffect, useImperativeHandle, useState } from "react"
import { Button, Input, Rate } from "antd";
import { download } from "../../utils/file/download";
import { reviewOpenQuest } from "../../request/api/judgment";
import ReactMarkdown from 'react-markdown';
import CustomIcon from "../../components/CustomIcon";
import { GetPercentScore } from "../../utils/int/bigInt";
const { TextArea } = Input;


function ChallengeJudgPage({data, isMobile, onFinish}, ref) {

    let [detail, setDetail] = useState();
    let [openQuest, setOpenQuest] = useState([]);
    let [reviewQuests, setReviewQuests] = useState([]);
    let [selectOpenQs, setSelectOpenQs] = useState({});
    let [page, setPage] = useState(0);

    // 比对当前已打分length 
    function isOver() {
        const flag = reviewQuests.length === detail.length;
        const remain = detail.length - reviewQuests.length;
        return  {flag, remain}
    }

    async function confirm(params) {
        // 没有改分直接退出
        if (reviewQuests.length === 0) {
            return
        }
        await reviewOpenQuest(reviewQuests)
        reviewQuests = [];
        setReviewQuests([...reviewQuests]);
        onFinish();
    }

    async function init() {
        page = 0;
        setPage(page);
        if (onFinish) {
            detail = data?.filter(e => e.open_quest_review_status === 1);
        }else{
            detail = data;
        }
        setDetail([...detail]);
        // 获取开放题列表
        const arr = [];
        detail.forEach((quest, i) => {
            let rate = 0;
            // 展示模式获取当前得分
            if (!onFinish) {
                rate = quest.answer?.correct ? quest.score : (quest.answer.score / quest.score * 5);
            }
            arr.push({
                index: i,
                isPass: null,
                rate: rate,
                title: quest.title,
                value: quest.answer.value,
                annex: quest.answer.annex,
                challenge_title: quest.challenge_title
            })
        })

        openQuest = arr;
        setOpenQuest([...openQuest]);
        selectOpenQs = openQuest[page];
        setSelectOpenQs({...selectOpenQs});
    }

    // 切换上下题
    function changePage(newPage) {
        page = newPage;
        setPage(page);
        selectOpenQs = openQuest[page];
        setSelectOpenQs({...selectOpenQs});
    }

    function getScore(percent) {
        // 记录rate
        openQuest[page].rate = percent;
        setOpenQuest([...openQuest]);
        selectOpenQs.rate = percent;
        setSelectOpenQs({...selectOpenQs});

        // 已打分列表
        const p = percent * 20 / 100;
        const info = detail[page];
        const score = GetPercentScore(info.score, p)
        const obj = {
            "id": info.ID,
            "answer": {
                "type": "open_quest",
                "annex": selectOpenQs.annex,
                "value": selectOpenQs.value,
                "score": score,
                "open_quest_review_time": new Date().toLocaleString().replace(/\//g, "-")
            },
            "index": info.index,
            "updated_at": info.updated_at
        };
        // 判断是否是新的
        const index = reviewQuests.findIndex(function(item) {
            return item.id === info.ID && item.index === info.index;
        });

        if (index === -1) {
            reviewQuests.push(obj)
        }else{
            reviewQuests[index] = obj;
        }
        setReviewQuests([...reviewQuests]);
    }

    useEffect(() => {
        data && init();
    },[data])

    useImperativeHandle(ref, () => ({
        confirm,
        isOver
    }))

    return (
        detail &&
        <div className="judg-content">
            <h1>{selectOpenQs?.challenge_title}</h1>
                <div className="judg-info">

                    <div className="item">
                        <p className="item-title">题目:</p>
                        <div className="item-content">
                            <ReactMarkdown>{selectOpenQs?.title}</ReactMarkdown>
                        </div>
                    </div>

                    <div className="item">
                        <p className="item-title">开放题答案:</p>
                        <TextArea 
                            className="item-content box"
                            bordered={false} 
                            maxLength={2000}
                            autoSize={{
                                minRows: 7,
                            }}
                            readOnly
                            value={selectOpenQs?.value}
                        />
                    </div>

                    <div className="item">
                        <p className="item-title">附件:</p>
                        <div className="item-content">
                            {
                                selectOpenQs?.annex && selectOpenQs?.annex.map(e => (
                                    <Button type="link" key={e.name} onClick={() => download(e.hash, e.name)}>{e.name}</Button>
                                ))
                            }
                        </div>
                    </div>

                    {/* <div className="item">
                        <p className="item-title">判定结果:&nbsp;<span style={{color: "#2B2F32"}}>{checked ? detail.quest_data.questions[page].score : 0}分</span></p>
                        <Radio.Group 
                            onChange={changePass}
                            className="isPass"
                            value={checked}
                            disabled={selectQuest.open_quest_review_status === 2}
                        >
                            <Radio value={true}>通过</Radio>
                            <Radio value={false}>不通过</Radio>
                        </Radio.Group>
                    </div> */}
                    <div className="item">
                        <div className="item-title flex">
                            <p className="item-title">评分: </p>
                                <Rate
                                    allowHalf 
                                    disabled={!onFinish}     //  预览模式不可选
                                    value={selectOpenQs?.rate}
                                    style={{color: "#DD8C53"}} 
                                    character={<CustomIcon type="icon-star" className="icon" />} 
                                    onChange={(percent) => getScore(percent)}
                                />
                        </div>
                    </div>
                </div>
            {
                onFinish &&
                <div className="pagination">
                    <Button disabled={page === 0} onClick={() => changePage(page - 1)}>上一题</Button>
                    <p>{page + 1}/<span style={{color: "#8B8D97"}}>{openQuest.length}</span></p>
                    <Button disabled={page+1 === openQuest.length} onClick={() => changePage(page + 1)}>下一题</Button>
                </div>
            }
        </div>
    )
}

export default forwardRef(ChallengeJudgPage)