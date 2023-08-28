import { Link, useNavigate, useParams } from "react-router-dom";
import {
    ArrowLeftOutlined,
  } from '@ant-design/icons';
import { Button, Form, InputNumber, Radio, Select, message } from "antd";
import { useEffect, useState } from "react";
import { getCollectionList, getQuest, updateQuest } from "../../request/api/quest";


export default function ChallengeModifyPage(params) {

    const [form] = Form.useForm();
    const { id, tokenId } = useParams();
    const navigateTo = useNavigate();
    const type = Form.useWatch("type", form);


    let [data, setData] = useState();
    let [fields, setFields] = useState([]);
    let [collection, setCollection] = useState([]);
    const [loading, setLoading] = useState(false);
    
    function onFinish({difficulty, estimateTime, collection_id, type}) {
        updateQuest({
            id: Number(id), 
            difficulty, 
            estimate_time: estimateTime * 60,
            collection_id: type === "default" ? 0 : collection_id
        })
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

    function init(params) {
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
        getQuest({id: Number(tokenId)})
        .then(res => {
            if (res.status === 0) {
                data = res.data;
                setData({...data});
                fields = [
                    {name: ["difficulty"], value: data.metadata.attributes.difficulty},
                    {name: ["estimateTime"], value: data.quest_data.estimateTime / 60},
                    {name: ["type"], value: data.collection_id === 0 ? "default" : "compilation"},
                    {name: ["collection_id"], value: data.collection_id}
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
                        label="教程(不可编辑)"
                        name="title"
                    >
                        {data.title}
                    </Form.Item>
                    <Form.Item
                        label="难度"
                        name="difficulty"
                    >
                        <Select
                            options={[
                                {label: "简单", value: 0},
                                {label: "一般", value: 1},
                                {label: "困难", value: 2},
                            ]}
                        />
                    </Form.Item>
                    <Form.Item
                        label="时长"
                        name="estimateTime"
                    >
                        <InputNumber controls={false} addonAfter="min" />
                    </Form.Item>
                    <Form.Item
                        label="挑战类型"
                        name="type"
                    >
                        <Radio.Group 
                            options={[
                                {label: "独立挑战", value: "default"},
                                {label: "选择合辑", value: "compilation"}
                            ]} 
                        />
                    </Form.Item>
                    {
                        type === "compilation" && 
                        <Form.Item
                            label="合辑名称"
                            name="collection_id"
                        >
                            <Select
                                options={collection}
                            />
                        </Form.Item>
                    }
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