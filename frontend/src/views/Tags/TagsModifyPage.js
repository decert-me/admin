import { Link, useLocation, useNavigate, useParams } from "react-router-dom"
import {
    ArrowLeftOutlined,
  } from '@ant-design/icons';
import { Button, Form, Input, InputNumber, Select, message } from "antd";
import { useEffect, useState } from "react";
import { updateLabel } from "../../request/api/tags";

export default function TagsModifyPage(params) {

    const { type, id } = useParams();
    const navigateTo = useNavigate();
    const location = useLocation();
    const [loading, setLoading] = useState(false);
    let [fields, setFields] = useState([]);
    
    const onFinish = (values) => {
        setLoading(true);
        updateLabel({...values, id: Number(id), weight: Number(values.weight)})
        .then(res => {
            if (res.code === 0) {
                message.success(res.msg);
                setTimeout(() => {
                    navigateTo(`/dashboard/tags`);
                }, 1000);
            }else{
                setLoading(false);
            }
        })
        .catch(err => {
            setLoading(false);
            message.error(err);
        })
    }

    function init() {
        const searchParams = new URLSearchParams(location.search);
        const params = {};
        for (const [key, value] of searchParams) {
            params[key] = value;
        }
        fields = [
            {
                name: ['weight'],
                value: params.weight
            },
            {
                name: ['english'],
                value: params.english
            },
            {
                name: ['chinese'],
                value: params.chinese
            }
        ]
        setFields([...fields]);
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="tags-modify">
            <Link to={`/dashboard/tags`}>
                <ArrowLeftOutlined />
            </Link>
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
                    label="权重"
                    name="weight"
                    rules={[{
                        required: true,
                        message: '请输入权重!',
                    }]}
                >
                    <InputNumber controls={false} />
                </Form.Item>

                <Form.Item
                    label="中文标题"
                    name="chinese"
                    rules={[{
                        required: true,
                        message: '请输入中文标题!',
                    }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="英文标题"
                    name="english"
                    rules={[{
                        required: true,
                        message: '请输入英文标题!',
                    }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading}>
                        修改标签
                    </Button>
                </Form.Item>

            </Form>
        </div>
    )
}