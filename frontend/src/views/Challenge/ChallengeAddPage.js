import { Link, useNavigate } from "react-router-dom";
import {
    ArrowLeftOutlined,
    PlusOutlined
  } from '@ant-design/icons';
import { Button, Form, Input, Select, Upload, message } from "antd";
import { UploadProps } from "../../utils/props";
import { useState } from "react";
import { createCollection } from "../../request/api/quest";
const { TextArea } = Input;


export default function ChallengeAddPage(params) {

    const [loading, setLoading] = useState(false);
    const navigateTo = useNavigate();


    function onFinish(values) {
        try {
            const cover = values.cover.file.response.data.hash;
            createCollection({...values, cover})
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
        } catch (error) {
            setLoading(false);
            message.error(error)
        }

    }


    
    return (
        <div className="challenge-add challenge">
            <Link to={`/dashboard/challenge/compilation`}>
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
                    label="合辑标题"
                    name="title"
                    rules={[{
                        required: true,
                        message: '请输入标题!',
                    }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="合辑简介"
                    name="description"
                >
                    <TextArea autoSize={{ minRows: 5 }} />
                </Form.Item>

                <Form.Item 
                    label="封面图" 
                    name="cover"
                    valuePropName="cover" 
                    rules={[{
                        required: true,
                        message: '请上传图片!',
                    }]}
                >
                    <Upload
                        listType="picture-card"
                        {...UploadProps}
                    >
                        <div>
                        <PlusOutlined />
                        <div style={{ marginTop: 8 }}>
                            Upload
                        </div>
                        </div>
                    </Upload>
                </Form.Item>

                <Form.Item
                    label="合辑作者"
                    name="author"
                    rules={[{
                        required: true,
                        message: '请输入合辑作者!',
                    }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="难度"
                    name="difficulty"
                >
                    <Select
                        placeholder="请选择难度"
                        options={[
                            {label: "困难", value: 2},
                            {label: "中等", value: 1},
                            {label: "简单", value: 0}
                        ]}
                    />
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading}>
                        添加挑战合辑
                    </Button>
                </Form.Item>
            </Form>
        </div>
    )
}