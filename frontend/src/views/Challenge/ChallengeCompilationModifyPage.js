import { Link, useNavigate, useParams } from "react-router-dom";
import {
    ArrowLeftOutlined,
    PlusOutlined
  } from '@ant-design/icons';
import { Button, Form, Input, InputNumber, Select, Upload, message } from "antd";
import { UploadProps } from "../../utils/props";
import { useEffect, useState } from "react";
import { getCollectionDetail, updateCollection } from "../../request/api/quest";
const { TextArea } = Input;


export default function ChallengeCompilationModifyPage(params) {
    
    const navigateTo = useNavigate();
    const {id} = useParams();
    const [loading, setLoading] = useState(false);
    let [fields, setFields] = useState([]);
    let [data, setData] = useState();

    function onFinish(values) {
        try {
            const cover = values.cover?.file?.response?.data?.hash || data.cover;
            updateCollection({...values, cover, id: Number(id)})
            .then(res => {
                if (res.code === 0) {
                    message.success(res.msg);
                    setTimeout(() => {
                        navigateTo("/dashboard/challenge/compilation");
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

    function init(params) {
        getCollectionDetail({id: Number(id)})
        .then(res => {
            if (res.code === 0) {
                data = res.data;
                setData({...data});
                fields = [
                    {name: ["title"], value: res.data.title},
                    {name: ["description"], value: res.data.description},
                    {name: ["author"], value: res.data.author},
                    {name: ["cover"], value: "https://ipfs.decert.me/"+res.data.cover},
                    {name: ["difficulty"], value: res.data.difficulty},
                    {name: ["sort"], value: Number(data.sort)},
                ]
                setFields([...fields]);
            }
        })
    }

    useEffect(() => {
        init();
    },[])

    return (
        data &&
        <div className="challenge">
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
                fields={fields}
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
                    rules={[{
                        required: true,
                        message: '请输入简介!',
                    }]}
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
                        defaultFileList={[{
                            uid: '-1',
                            name: 'image.png',
                            status: 'done',
                            url: "https://ipfs.decert.me/"+data.cover,
                        }]}
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

                <Form.Item
                    label="权重"
                    name="sort"
                >
                    <InputNumber controls={false} />
                </Form.Item>

                <Form.Item>
                    <Button type="primary" htmlType="submit" loading={loading}>
                        保存
                    </Button>
                </Form.Item>
            </Form>
        </div>
    )
}