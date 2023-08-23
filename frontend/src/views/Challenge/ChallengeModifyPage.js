import { Link, useNavigate, useParams } from "react-router-dom";
import {
    ArrowLeftOutlined,
  } from '@ant-design/icons';
import { Button, Form, InputNumber, Select, message } from "antd";
import { useEffect, useState } from "react";
import { getQuest, updateQuest } from "../../request/api/quest";


export default function ChallengeModifyPage(params) {

    const { id, tokenId } = useParams();
    const navigateTo = useNavigate();
    let [data, setData] = useState();
    let [fields, setFields] = useState([]);
    const [loading, setLoading] = useState(false);
    
    function onFinish({difficulty, estimateTime}) {
        updateQuest({id: Number(id), difficulty, estimate_time: estimateTime * 60})
        .then(res => {
            if (res.code === 0) {
                message.success(res.msg);
                setTimeout(() => {
                    navigateTo("/dashboard/challenge");
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
        getQuest({id: Number(tokenId)})
        .then(res => {
            if (res.status === 0) {
                data = res.data;
                setData({...data});
                fields = [
                    {name: ["difficulty"], value: data.metadata.attributes.difficulty},
                    {name: ["estimateTime"], value: data.quest_data.estimateTime}
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
            <Link to={`/dashboard/challenge`}>
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