import { Button, Form, Input, Select, message } from "antd";
import { Link, useNavigate } from "react-router-dom";
import {
    ArrowLeftOutlined,
  } from '@ant-design/icons';
import { createLabel } from "../../request/api/tags";
import { useState } from "react";


export default function TagsAddPage(params) {

    const navigateTo = useNavigate();
    const [loading, setLoading] = useState(false);

    const onFinish = (values) => {
        setLoading(true);
        createLabel(values)
        .then(res => {
            if (res.code === 0) {
                message.success(res.msg);
                setTimeout(() => {
                    navigateTo("/dashboard/tags")
                }, 500);
            }else{
                setLoading(false);
                message.error(res.msg)
            }
        }).catch(err => {
            setLoading(false);
            message.error(err)
        })
    };
    
    return (
        <div className="tags-add">
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
            >
                <Form.Item
                    label="请选择父级标签"
                    name="type"
                    rules={[{
                        required: true,
                        message: '请输入标题!',
                    }]}
                >
                    <Select
                        style={{
                            width: 120,
                        }}
                        options={[
                            { value: 'category', label: '分类' },
                            { value: 'language', label: '语言' }
                        ]}
                    />
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
                        添加标签
                    </Button>
                </Form.Item>

            </Form>
        </div>
    )
}