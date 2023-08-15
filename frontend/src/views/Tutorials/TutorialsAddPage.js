import { Button, Form, Input, InputNumber, Select, Upload, message } from 'antd';
import {
    PlusOutlined,
    MinusCircleOutlined
  } from '@ant-design/icons';
import { useEffect, useState } from 'react';
import "./index.scss"
import { UploadProps } from '../../utils/props';
import { createTutorial } from '../../request/api/tutorial';
import { useNavigate } from 'react-router-dom';
const { TextArea } = Input;



export default function TutorialsAddPage(params) {
    
    const table = require("./category_tabel.json");
    const [form] = Form.useForm();
    const videoCategory = Form.useWatch("videoCategory", form);
    const navigateTo = useNavigate();

    let [category, setCategory] = useState();     //  类别 选择器option
    let [theme, setTheme] = useState();     //  主题 选择器option
    let [lang, setLang] = useState();     //  语种 选择器option
    let [doctype, setDoctype] = useState("doc");
    const [loading, setLoading] = useState(false);

    const onFinish = (values) => {
        setLoading(true);
        const {
            repoUrl, label, catalogueName, docType, desc, 
            challenge, branch, docPath, commitHash, 
            category, theme, language, difficulty,
        } = values;
        const img = values.img.file.response.data.hash;

        if (doctype === "doc") {
            const obj = {
                repoUrl, label, catalogueName, docType, img, desc, 
                challenge, branch, docPath, commitHash,
                category, theme, language, difficulty
            }
            addArticle(obj)
        }else{
            // addVideo(values)
        }
    };

    function addArticle(obj) {
        createTutorial(obj)
        .then(res => {
            if (res.code === 0) {
                message.success(res.msg);
                setTimeout(() => {
                    navigateTo("/dashboard/tutorials/list");
                }, 500);
            }else{
                setLoading(false);
                message.success(res.msg);
            }
        })
        .catch(err => {
            setLoading(false);
            message.error(err)
        })
    }

    function addVideo({}) {
        
    }

    function init() {
        optionsInit()
    }

    // 选择器option初始化
    function optionsInit(params) {
        // 类别初始化
        let categoryOption = [];
        for (const key in table.category) {
                if (Object.hasOwnProperty.call(table.category, key)) {
                    categoryOption.push({
                        value: key,
                        label: table.category[key]
                    })
                }
        }
        category = categoryOption;
        setCategory([...category])
        // 主题初始化
        let themeOption = [];
        for (const key in table.theme) {
                if (Object.hasOwnProperty.call(table.theme, key)) {
                    themeOption.push({
                        value: key,
                        label: table.theme[key]
                    })
                }
        }
        theme = themeOption;
        setTheme([...theme])
    
        // 语种
        let langOption = [];
        for (const key in table.language) {
                if (Object.hasOwnProperty.call(table.language, key)) {
                    langOption.push({
                        value: key,
                        label: table.language[key]
                    })
                }
        }
        lang = langOption;
        setLang([...lang]);
    }
 
    useEffect(() => {
        init()
    },[])

    return (
        <Form
            name="basic"
            labelCol={{ span: 6 }}
            wrapperCol={{ span: 18 }}
            style={{ maxWidth: 800 }}
            onFinish={onFinish}
            autoComplete="off"
            form={form}
        >

            <Form.Item
                label="标题"
                name="label"
                rules={[{
                    required: true,
                    message: '请输入标题!',
                }]}
            >
                <Input />
            </Form.Item>

            <Form.Item
                label="描述"
                name="desc"
                rules={[{
                    required: true,
                    message: '请输入描述!',
                }]}
            >
                <TextArea autoSize={{ minRows: 3 }} />
            </Form.Item>

            <Form.Item 
                label="图片" 
                name="img"
                valuePropName="img" 
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
                label="挑战编号"
                name="challenge"
            >
                <InputNumber controls={false} />
            </Form.Item>

            <Form.Item
                label="教程类型"
            >
                <div className="doctype">
                    <div className={`box ${doctype === "video" ? "active-box" : ""}`} onClick={() => setDoctype("video")}>
                        视频
                    </div>
                    <div className={`box ${doctype !== "video" ? "active-box" : ""}`} onClick={() => setDoctype("doc")}>
                        文档
                    </div>
                </div>
            </Form.Item>
            {
                doctype === "video" ?
                <>
                    <Form.Item
                        label="视频地址"
                        name="url"
                        rules={[{
                            required: true,
                            message: '请输入视频地址!',
                        }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="视频类型"
                        name="videoCategory"
                        rules={[{
                            required: true,
                            message: '请输入视频类型!',
                        }]}
                    >
                        <Select
                            placeholder="请选择视频类型!"
                            options={[
                                {label: "youtube", value: "youtube"},
                                {label: "bilibili", value: "bilibili"}
                            ]}
                        />
                    </Form.Item>
                    {
                        videoCategory === "bilibili" &&
                        <Form.Item
                            label="视频列表"
                        >
                            <Form.List
                                name="videoItems"
                            >
                                {(fields, { add, remove }, { errors }) => (
                                <>
                                    {fields.map((field, index) => (
                                    <Form.Item
                                        required={false}
                                        key={field.key}
                                    >
                                        <Form.Item
                                        {...field}
                                        validateTrigger={['onChange', 'onBlur']}
                                        rules={[
                                            {
                                            required: true,
                                            whitespace: true,
                                            message: "请输入视频链接或删除该输入框！",
                                            },
                                        ]}
                                        noStyle
                                        >
                                        <Input
                                            placeholder="bilibili的嵌入代码"
                                            style={{
                                            width: '90%',
                                            }}
                                        />
                                        </Form.Item>
                                        {fields.length > 1 ? (
                                        <MinusCircleOutlined
                                            className="dynamic-delete-button"
                                            onClick={() => remove(field.name)}
                                        />
                                        ) : null}
                                    </Form.Item>
                                    ))}
                                    <Form.Item>
                                    <Button
                                        type="dashed"
                                        onClick={() => add()}
                                        style={{
                                        width: '60%',
                                        }}
                                        icon={<PlusOutlined />}
                                    >
                                        Add field
                                    </Button>
                                    <Form.ErrorList errors={errors} />
                                    </Form.Item>
                                </>
                                )}
                            </Form.List>
                        </Form.Item>
                    }
                </>
                :
                <>
                    <Form.Item
                        label="教程地址"
                        name="repoUrl"
                        rules={[{
                            required: true,
                            message: '请输入教程地址!',
                        }]}
                    >
                        <Input addonBefore={
                            <Form.Item
                                // label="教程类型"
                                name="docType"
                                rules={[{
                                    required: true,
                                    message: '请选择教程类型!',
                                }]}
                                noStyle
                            >
                                <Select
                                    placeholder="教程类型"
                                    options={[
                                        {label: "docusaurus", value: "docusaurus"},
                                        {label: "gitbook", value: "gitbook"},
                                        {label: "mdBook", value: "mdBook"},
                                    ]}
                                    style={{
                                        width: 120
                                    }}
                                />
                            </Form.Item>                
                        } />
                    </Form.Item>
                    <Form.Item
                        label="分支"
                        name="branch"
                    >
                        <Input placeholder="默认为main分支" />
                    </Form.Item>
                    <Form.Item
                        label="教程文档目录"
                        name="docPath"
                    >
                        <Input placeholder="默认为根目录" />
                    </Form.Item>
                    <Form.Item
                        label="教程文档commitHash"
                        name="commitHash"
                    >
                        <Input placeholder="默认为最新" />
                    </Form.Item>
                </>
            }

            <Form.Item
                label="类别"
                name="category"
            >
                <Select
                    mode="multiple"
                    placeholder="请至少选择一项类别"
                    options={category}
                />
            </Form.Item>
            
            <Form.Item
                label="主题"
                name="theme"
            >
                <Select
                    mode="multiple"
                    placeholder="请至少选择一项主题"
                    options={theme}
                />
            </Form.Item>

            <Form.Item
                label="语种"
                name="language"
            >
                <Select
                    placeholder="请至少选择一项主题"
                    options={lang}
                />
            </Form.Item>

            <Form.Item
                label="预估时间"
                name="time"
            >
                <InputNumber addonAfter="min" controls={false} />
            </Form.Item>

            <Form.Item
                label="难度"
                name="difficulty"
            >
                <Select
                    placeholder="请选择难度"
                    options={[
                        {label: "困难", value: 2},
                        {label: "一般", value: 1},
                        {label: "简单", value: 0}
                    ]}
                />
            </Form.Item>

            <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                    添加教程
                </Button>
            </Form.Item>

        </Form>
    )
}