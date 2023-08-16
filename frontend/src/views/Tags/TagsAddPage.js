import { Button, Form, Input, Select } from "antd";
import { Link } from "react-router-dom";
import {
    ArrowLeftOutlined,
  } from '@ant-design/icons';


export default function TagsAddPage(params) {

    const onFinish = (values) => {
        console.log('Success:', values);
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
                    name="tag"
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
                    name="content_zh"
                    rules={[{
                        required: true,
                        message: '请输入中文标题!',
                    }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="英文标题"
                    name="content_en"
                    rules={[{
                        required: true,
                        message: '请输入英文标题!',
                    }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit">
                        Submit
                    </Button>
                </Form.Item>

            </Form>
        </div>
    )
}