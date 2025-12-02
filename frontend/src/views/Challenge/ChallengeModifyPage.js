import { Link, useNavigate, useParams } from "react-router-dom";
import {
    ArrowLeftOutlined,
    PlusOutlined,
    DeleteOutlined,
  } from '@ant-design/icons';
import { Button, Form, Input, InputNumber, Select, message, Card, Space, Radio } from "antd";
import { useEffect, useState } from "react";
import { getCollectionList, getQuest, updateQuest } from "../../request/api/quest";
import { getLabelList } from "../../request/api/tags";
import ReactQuill from 'react-quill';
import 'react-quill/dist/quill.snow.css';
const { TextArea } = Input;

export default function ChallengeModifyPage(params) {

    const [form] = Form.useForm();
    const { id, tokenId } = useParams();
    const navigateTo = useNavigate();

    
    let [data, setData] = useState();
    let [fields, setFields] = useState([]);
    let [collection, setCollection] = useState([]);
    const [loading, setLoading] = useState(false);
    const [categoryOption, setCategoryOption] = useState([]);
    const [category, setCategory] = useState([]);

    // 新增状态：题目内容和题目列表
    const [questContent, setQuestContent] = useState('');
    const [questions, setQuestions] = useState([]);
    const [canEditQuestData, setCanEditQuestData] = useState(true); // 是否可以编辑题目数据
    
    function onFinish({difficulty, estimateTime, collection_id, type, sort, description}) {
        setLoading(true);
        const obj = {
            id: Number(id),
            difficulty,
            estimate_time: estimateTime && estimateTime !== 0 ? estimateTime * 60 : null,
            sort,
            category,
            collection_id: collection_id ? [collection_id] : [],
            description
        }

        // 如果可以编辑题目数据，则添加quest_data字段
        if (canEditQuestData) {
            obj.quest_data = {
                ...data.quest_data,
                content: questContent,
                questions: questions
            };
        }

        if (data.metadata?.description) {
            delete obj.description
        }
        updateQuest(obj)
        .then(res => {
            if (res.code === 0) {
                message.success(res.msg);
                setTimeout(() => {
                    navigateTo("/dashboard/challenge/list");
                }, 500);
            }else{
                setLoading(false);
            }
        })
        .catch(err => {
            setLoading(false);
            message.error(err)
        })
    }

    function changeCategory(value) {
        if (value.length > 5) {
            return
        }
        setCategory([...value]);
    }

    // 新增：添加题目
    function addQuestion() {
        const newQuestion = {
            type: 'multiple_choice',
            title: '',
            score: 10,
            options: ['']
        };
        setQuestions([...questions, newQuestion]);
    }

    // 新增：删除题目
    function deleteQuestion(index) {
        const newQuestions = questions.filter((_, i) => i !== index);
        setQuestions(newQuestions);
    }

    // 新增：更新题目
    function updateQuestion(index, field, value) {
        const newQuestions = [...questions];
        newQuestions[index][field] = value;
        setQuestions(newQuestions);
    }

    // 新增：添加选项
    function addOption(questionIndex) {
        const newQuestions = [...questions];
        if (!newQuestions[questionIndex].options) {
            newQuestions[questionIndex].options = [];
        }
        newQuestions[questionIndex].options.push('');
        setQuestions(newQuestions);
    }

    // 新增：删除选项
    function deleteOption(questionIndex, optionIndex) {
        const newQuestions = [...questions];
        newQuestions[questionIndex].options = newQuestions[questionIndex].options.filter((_, i) => i !== optionIndex);
        setQuestions(newQuestions);
    }

    // 新增：更新选项
    function updateOption(questionIndex, optionIndex, value) {
        const newQuestions = [...questions];
        newQuestions[questionIndex].options[optionIndex] = value;
        setQuestions(newQuestions);
    }

    function init(params) {
        getLabelList({type: "category"})
        .then(res => {
            if (res.code === 0) {
            const list = res.data;
            const data = list ? list : [];
            // 添加key
            data.forEach(ele => {
                ele.value = ele.ID
                ele.label = ele.Chinese
            })
            setCategoryOption([...data]);
            }else{
                message.success(res.msg);
            }
        })
        .catch(err => {
            message.error(err)
        })
        getCollectionList()
        .then(res => {
            if (res.code === 0) {
                const list = res.data.list;
                const arr = list ? list : [];
                collection = [];
                arr.forEach(e => {
                    collection.push({ label: e.title, value: e.id })
                })
                setCollection([...collection]);
            }
        })
        .catch(err => {
            console.error('getCollectionList error:', err);
            message.error('获取合辑列表失败: ' + err);
        })
        getQuest({id: tokenId})
        .then(res => {
            console.log('getQuest response:', res);
            if (res.code === 0) {
                data = res.data;
                console.log('Quest data:', data);
                setData({...data});
                setCategory([...data.category||[]]);

                // 初始化题目内容和题目列表
                setQuestContent(data.quest_data?.content || '');
                setQuestions(data.quest_data?.questions || []);

                // 允许编辑题目数据
                setCanEditQuestData(true);

                fields = [
                    {name: ["difficulty"], value: data.metadata?.attributes?.difficulty || 0},
                    {name: ["estimateTime"], value: data.quest_data?.estimateTime ? data.quest_data.estimateTime / 60 : 0},
                    {name: ["sort"], value: Number(data.sort || 0)},
                    {name: ["type"], value: data.collection_id?.length === 0 ? "default" : "compilation"},
                    {name: ["collection_id"], value: data.collection_id?.[0]},
                    {name: ["description"], value: data.description || ''}
                ];
                console.log('Fields:', fields);
                setFields([...fields]);
            }
        })
        .catch(err => {
            console.error('getQuest error:', err);
            message.error('获取挑战数据失败: ' + err);
        })
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="challenge">
            <Link to={`/dashboard/challenge/list`}>
                <ArrowLeftOutlined />
            </Link>
            <h2>编辑</h2>
            {
                data &&
                <Form
                    name="basic"
                    labelCol={{ span: 6 }}
                    wrapperCol={{ span: 18 }}
                    style={{ maxWidth: 800 }}
                    onFinish={onFinish}
                    autoComplete="off"
                    fields={fields}
                    form={form}
                >
                    <Form.Item
                        label="NFT(不可编辑)"
                        name="nft"
                    >
                        {data.metadata?.image && <img src={data.metadata.image.replace("ipfs://", "https://ipfs.decert.me/")} alt="" style={{height: "100px"}} />}
                    </Form.Item>
                    <Form.Item
                        label="标题(不可编辑)"
                        name="title"
                    >
                        {data.title}
                    </Form.Item>
                    {
                        data.metadata?.description ?
                        <Form.Item
                            label="描述(不可编辑)"
                            name="description"
                        >
                            {data.description}
                        </Form.Item>
                        :
                        <Form.Item
                            label="描述"
                            name="description"
                        >
                            <TextArea
                                autoSize={{
                                    minRows: 3,
                                    maxRows: 5,
                                }}
                            />
                        </Form.Item>
                    }
                    <Form.Item
                        label="难度"
                        name="difficulty"
                    >
                        <Select
                            options={[
                                {label: "简单", value: 0},
                                {label: "中等", value: 1},
                                {label: "困难", value: 2},
                            ]}
                        />
                    </Form.Item>
                    {/* <Form.Item
                        label="分类"
                        name="category"
                    > */}
                    <div style={{display: "flex", alignItems: "center", gap: "8px", marginBottom: "24px"}}>
                        <div style={{width: "190px", textAlign: "right"}}>
                            <lable>分类:</lable>
                        </div>
                        <Select
                            options={categoryOption}
                            mode="tags"
                            onChange={changeCategory}
                            value={category}
                            style={{width: "600px"}}
                        />
                    </div>
                    {/* </Form.Item> */}
                    <Form.Item
                        label="权重"
                        name="sort"
                    >
                        <InputNumber controls={false} />
                    </Form.Item>
                    <Form.Item
                        label="时长"
                        name="estimateTime"
                    >
                        <InputNumber controls={false} addonAfter="min" />
                    </Form.Item>
                    <Form.Item
                        label="合辑名称"
                        name="collection_id"
                    >
                        <Select
                            allowClear
                            options={collection}
                        />
                    </Form.Item>

                    {/* 新增：题目列表编辑 */}
                    <div style={{ marginTop: 24, marginBottom: 24 }}>
                        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                            <h3>{canEditQuestData ? '题目列表' : '题目列表（不可编辑）'}</h3>
                            {canEditQuestData && (
                                <Button
                                    type="dashed"
                                    icon={<PlusOutlined />}
                                    onClick={addQuestion}
                                >
                                    添加题目
                                </Button>
                            )}
                        </div>

                        {questions.map((question, qIndex) => (
                            <Card
                                key={qIndex}
                                title={`题目 ${qIndex + 1}`}
                                extra={canEditQuestData && (
                                    <Button
                                        type="text"
                                        danger
                                        icon={<DeleteOutlined />}
                                        onClick={() => deleteQuestion(qIndex)}
                                    >
                                        删除
                                    </Button>
                                )}
                                style={{ marginBottom: 16 }}
                            >
                                <Space direction="vertical" style={{ width: '100%' }} size="middle">
                                    {/* 题目类型 */}
                                    <div>
                                        <label>题目类型：</label>
                                        <Radio.Group
                                            value={question.type}
                                            onChange={(e) => updateQuestion(qIndex, 'type', e.target.value)}
                                            disabled={!canEditQuestData}
                                        >
                                            <Radio value="multiple_choice">单选题</Radio>
                                            <Radio value="multiple_response">多选题</Radio>
                                            <Radio value="open_quest">开放题</Radio>
                                        </Radio.Group>
                                    </div>

                                    {/* 题目标题 */}
                                    <div>
                                        <label>题目标题：</label>
                                        <ReactQuill
                                            value={question.title}
                                            onChange={(value) => updateQuestion(qIndex, 'title', value)}
                                            readOnly={!canEditQuestData}
                                            theme="snow"
                                            placeholder="请输入题目标题"
                                            style={{
                                                background: canEditQuestData ? 'white' : '#f5f5f5',
                                                minHeight: '100px'
                                            }}
                                        />
                                    </div>

                                    {/* 题目分数 */}
                                    <div>
                                        <label>题目分数：</label>
                                        <InputNumber
                                            value={question.score}
                                            onChange={(value) => updateQuestion(qIndex, 'score', value)}
                                            min={0}
                                            disabled={!canEditQuestData}
                                        />
                                    </div>

                                    {/* 选项（仅非开放题显示） */}
                                    {question.type !== 'open_quest' && (
                                        <div>
                                            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                                                <label>选项：</label>
                                                {canEditQuestData && (
                                                    <Button
                                                        size="small"
                                                        type="dashed"
                                                        icon={<PlusOutlined />}
                                                        onClick={() => addOption(qIndex)}
                                                    >
                                                        添加选项
                                                    </Button>
                                                )}
                                            </div>
                                            {question.options?.map((option, oIndex) => (
                                                <div key={oIndex} style={{ display: 'flex', gap: 8, marginBottom: 8 }}>
                                                    <span style={{ lineHeight: '32px' }}>选项 {oIndex + 1}:</span>
                                                    <Input
                                                        value={option}
                                                        onChange={(e) => updateOption(qIndex, oIndex, e.target.value)}
                                                        placeholder={`请输入选项 ${oIndex + 1}`}
                                                        disabled={!canEditQuestData}
                                                        style={{ flex: 1 }}
                                                    />
                                                    {canEditQuestData && question.options.length > 1 && (
                                                        <Button
                                                            danger
                                                            icon={<DeleteOutlined />}
                                                            onClick={() => deleteOption(qIndex, oIndex)}
                                                        />
                                                    )}
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </Space>
                            </Card>
                        ))}

                        {questions.length === 0 && (
                            <div style={{ textAlign: 'center', padding: '40px 0', color: '#999' }}>
                                暂无题目
                            </div>
                        )}
                    </div>

                    <Form.Item>
                        <Button type="primary" htmlType="submit" loading={loading}>
                            保存
                        </Button>
                    </Form.Item>
                </Form>
            }
        </div>
    )
}