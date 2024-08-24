import { useEffect, useState } from "react"
import { getUserQuestDetail } from "../../request/api/judgment";
import { Button, Input, InputNumber, message, Slider } from "antd";
import { download } from "../../utils/file/download";
import "./judg.scss";

const { TextArea } = Input;


export default function JudgReviewModal({uuid, address}) {
    
    const [data, setData] = useState({});
    const [questList, setQuestList] = useState([]);
    const [index, setIndex] = useState(0);
    const [selectOpenQs, setSelectOpenQs] = useState({});

    function changeIndex(index) {
        setIndex(index);
        setSelectOpenQs(questList[index]);
    }

    function init() {
        getUserQuestDetail({uuid, address})
        .then(res => {
            if (res.code === 0) {
                const quest = res.data.quest_data.questions || [];
                const all_score = quest.map(e => e.score).reduce((accumulator, currentValue) => accumulator + currentValue, 0);
                const data = {...res.data, all_score} || {};
                const answer = res.data.answer || [];
                
                const arr = quest.map((e, i) => {
                    let value = answer[i].value;
                    if (e.type === "multiple_choice") {
                        value = e.options[value];
                    }
                    if (e.type === "multiple_response") {
                        value = value.map(item => e.options[item]).join("\n");
                    }
                    if (e.type === "coding") {
                        value = answer[i].code;
                    }
                    return {
                        score: e.score,
                        description: e?.description,
                        title: e.title,
                        annex: answer[i].annex,
                        user_score: answer[i].score,
                        annotation: answer[i].annotation,
                        open_quest_review_time: answer[i].open_quest_review_time,
                        value: value,
                        type: e.type
                    }
                })
                setData(data);
                setQuestList(arr);
                setSelectOpenQs(arr[index]);
            }else{
                message.error(res.msg);
            }
        })
        .catch(err => {
            message.error(err.msg);
        })
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="judg-content">
            <h1>{data?.title}</h1>
                <div className="judg-info">
                    <div className="item">
                        <div className="item-title">挑战者: &nbsp;<a href={`${process.env.REACT_APP_LINK_URL || "https://decert.me"}/user/${data?.address}`} target="_blank" rel="noopener noreferrer">{data?.address}</a></div>
                    </div>

                    <div className="item">
                        <div className="item-title">提交时间: &nbsp;
                            <span className="item-content">{
                            data?.submit_time && data?.submit_time.indexOf("0001-01-01T") === -1
                            ? data?.submit_time.replace("T", " ").split(".")[0].split("+")[0]
                            : "-"}</span>
                        </div>
                    </div>

                    <div className="item">
                        <div className="item-title">题目: &nbsp;
                            <span className="item-content">{selectOpenQs?.title}</span>
                        </div>
                        {
                            selectOpenQs?.description &&
                            <div className="item-title">题干: &nbsp;
                                <span className="item-content">{selectOpenQs?.description}</span>
                            </div>
                        }
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
                            value={selectOpenQs?.value}
                        />
                    </div>
                    {
                        selectOpenQs.type === "open_quest" &&
                        <div className="item">
                            <p className="item-title">批注:</p>
                            <TextArea 
                                disabled
                                className="item-content box"
                                bordered={false} 
                                autoSize={{
                                    minRows: 3,
                                    maxRows: 5,
                                }}
                                value={selectOpenQs?.annotation}
                            />
                        </div>
                    }
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

                    <div className="item">
                        <div className="item-title">挑战总得分: &nbsp;
                            <span className="item-content">{data?.all_score}</span>
                        </div>
                    </div>

                    <div className="item">
                        <div className="item-title">挑战及格分: &nbsp;
                            <span className="item-content">{data?.quest_data?.passingScore}</span>
                        </div>
                    </div>

                    <div className="item">
                        <div className="item-title">本题评分: &nbsp;
                            <InputNumber
                                disabled
                                min={0}
                                max={selectOpenQs?.score}
                                style={{margin: '0 16px'}}
                                value={selectOpenQs?.user_score}
                            />
                        </div>
                        <div style={{width: "352px"}}>
                            <Slider
                                disabled
                                max={selectOpenQs?.score}
                                tooltip={{formatter: null}}
                                step={0.1}
                                value={selectOpenQs?.user_score}
                            />
                        </div>
                    </div>
                </div>
                <div className="flex-bte">
                    <div className="pagination">
                        <Button disabled={index === 0} onClick={() => changeIndex(index - 1)}>上一题</Button>
                        <p>{index+1}/<span style={{color: "#8B8D97"}}>{questList.length}</span></p>
                        <Button disabled={index + 1 === questList.length} onClick={() => changeIndex(index + 1)}>下一题</Button>
                    </div>
                </div>
        </div>
    )
}