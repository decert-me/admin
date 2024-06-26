import { Link, useNavigate, useParams } from "react-router-dom";
import {
    ArrowLeftOutlined,
  } from '@ant-design/icons';
import { Button, Form, Input, InputNumber, Select, message } from "antd";
import { useEffect, useState } from "react";
import { getCollectionList, getQuest, updateQuest } from "../../request/api/quest";
import { getLabelList } from "../../request/api/tags";
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
    
    function onFinish({difficulty, estimateTime, collection_id, type, sort, description}) {
        const obj = {
            id: Number(id), 
            difficulty, 
            estimate_time: estimateTime && estimateTime !== 0 ? estimateTime * 60 : null,
            sort,
            category,
            collection_id: collection_id ? [collection_id] : [],
            description
        }
        if (data.metadata.description) {
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
        getQuest({id: tokenId})
        .then(res => {
            if (res.code === 0) {
                data = res.data;
                setData({...data});
                setCategory([...data.category||[]]);
                fields = [
                    {name: ["difficulty"], value: data.metadata.attributes.difficulty},
                    {name: ["estimateTime"], value: data.quest_data.estimateTime / 60},
                    {name: ["sort"], value: Number(data.sort)},
                    {name: ["type"], value: data.collection_id.length === 0 ? "default" : "compilation"},
                    {name: ["collection_id"], value: data.collection_id[0]},
                    {name: ["description"], value: data.description}
                ];
                setFields([...fields]);
            }
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
                        <img src={data.metadata.image.replace("ipfs://", "https://ipfs.decert.me/")} alt="" style={{height: "100px"}} />
                    </Form.Item>
                    <Form.Item
                        label="标题(不可编辑)"
                        name="title"
                    >
                        {data.title}
                    </Form.Item>
                    {
                        data.metadata.description ? 
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