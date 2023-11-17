import { forwardRef, useEffect, useImperativeHandle, useState } from "react"
import { getQuest } from "../../request/api/quest";
import { Button, Radio, message } from "antd";
import { download } from "../../utils/file/download";
import { reviewOpenQuest } from "../../request/api/judgment";
import ReactMarkdown from 'react-markdown';


function ChallengeJudgPage({selectQuest, onFinish}, ref) {

    let [detail, setDetail] = useState();
    let [openQuest, setOpenQuest] = useState([]);
    let [selectOpenQs, setSelectOpenQs] = useState({});
    let [page, setPage] = useState(0);
    let [checked, setChecked] = useState();

    function confirm(params) {
        // 判断是否已经处理过
        if (selectQuest.open_quest_review_status === 2) {
            onFinish();
            return
        }
        let answer = selectQuest?.answer || [];
        for (let i = 0; i < openQuest.length; i++) {
            const quest = openQuest[i];
            if (quest.isPass === null) {
                message.warning(`请为第${i+1}道开放题打分!`)
                return
            }
            if (answer[quest.index]) {
                answer[quest.index].correct = quest.isPass;
            }
        }
        reviewOpenQuest({
            id: selectQuest.ID,
            answer
        })
        .then(res => {
            if (res.code === 0) {
                message.success("操作成功!");
            }else{
                message.error("操作失败!");
            }
            onFinish();
        })
    }

    function changePass(event) {
        checked = event.target.value;
        setChecked(checked);
        openQuest[page].isPass = event.target.value;
        setOpenQuest([...openQuest]);
    }

    async function init() {
        // 获取quest详情
        await getQuest({id: selectQuest.token_id})
        .then(res => {
            if (res.code === 0) {
                detail = res.data;
                setDetail({...detail});
            }
        })
        // 获取开放题列表
        const { quest_data } = detail;
        const arr = [];
        quest_data?.questions.forEach((quest, i) => {
            if (quest.type === "open_quest") {
                arr.push({
                    title: quest.title,
                    value: selectQuest.answer[i]?.value,
                    annex: selectQuest.answer[i]?.annex,
                    index: i,
                    isPass: null,
                    correct: selectQuest.answer[i]?.correct
                })
            }
        })

        openQuest = arr;
        setOpenQuest([...openQuest]);
        selectOpenQs = openQuest[page];
        setSelectOpenQs({...selectOpenQs});
        // 查看模式赋值
        if (selectQuest.open_quest_review_status === 2) {
            checked = openQuest[0].correct;
            setChecked(checked);
        }
    }

    function changePage(newPage) {
        page = newPage;
        setPage(page);
        checked = null;
        setChecked(checked);
        // 查看模式赋值
        if (selectQuest.open_quest_review_status === 2) {
            checked = openQuest[page].correct;
            setChecked(checked);
        }
    }

    useEffect(() => {
        selectQuest && init();
    },[selectQuest])

    useImperativeHandle(ref, () => ({
        confirm
    }))

    return (
        detail &&
        <div className="judg-content">
            <h1>{detail?.title}</h1>
                <div className="judg-info">
                    <div className="item">
                        <p className="item-title">题目:</p>
                        <div className="item-content">
                            <ReactMarkdown>{selectOpenQs.title}</ReactMarkdown>
                        </div>
                    </div>

                    <div className="item">
                        <p className="item-title">开放题答案:</p>
                        <p className="item-content box">{selectOpenQs.value}</p>
                    </div>

                    <div className="item">
                        <p className="item-title">附件:</p>
                        <div className="item-content">
                            {
                                selectOpenQs.annex.map(e => (
                                    <Button type="link" key={e.name} onClick={() => download(e.hash, e.name)}>{e.name}</Button>
                                ))
                            }
                        </div>
                    </div>

                    <div className="item">
                        <p className="item-title">判定结果:</p>
                        <Radio.Group 
                            onChange={changePass}
                            className="isPass"
                            value={checked}
                            disabled={selectQuest.open_quest_review_status === 2}
                        >
                            <Radio value={true}>通过</Radio>
                            <Radio value={false}>不通过</Radio>
                        </Radio.Group>
                    </div>
                </div>
            <div className="pagination">
                <Button disabled={page === 0} onClick={() => changePage(page - 1)}>上一题</Button>
                <p>{page + 1}/<span style={{color: "#8B8D97"}}>{openQuest.length}</span></p>
                <Button disabled={page+1 === openQuest.length} onClick={() => changePage(page + 1)}>下一题</Button>
            </div>
        </div>
    )
}

export default forwardRef(ChallengeJudgPage)