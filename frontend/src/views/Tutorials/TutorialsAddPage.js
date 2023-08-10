import { Button, Form, Input, InputNumber, Select, Upload } from 'antd';
import {
    PlusOutlined,
    MinusCircleOutlined
  } from '@ant-design/icons';
import { useEffect, useState } from 'react';
const { TextArea } = Input;


export default function TutorialsAddPage(params) {
    
    const table = require("./category_tabel.json");
    const [form] = Form.useForm();
    const videoCategory = Form.useWatch("videoCategory", form);
    const docType = Form.useWatch("docType", form);

    let [category, setCategory] = useState();     //  类别 选择器option
    let [theme, setTheme] = useState();     //  主题 选择器option
    let [lang, setLang] = useState();     //  语种 选择器option

    const onFinish = (values) => {
        console.log('Success:', values);
        console.log();

    };

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
                        key: key,
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
                        key: key,
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
                        key: key,
                        label: table.language[key]
                    })
                }
            }
            lang = langOption;
            setLang([...lang])
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
                    action="/upload.do" 
                    listType="picture-card"
                    maxCount={1}
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
                name="docType"
                rules={[{
                    required: true,
                    message: '请选择教程类型!',
                }]}
            >
                <Select
                    options={[
                        {label: "docusaurus", value: "docusaurus"},
                        {label: "gitbook", value: "gitbook"},
                        {label: "mdBook", value: "mdBook"},
                        {label: "video", value: "video"}
                    ]}
                />
            </Form.Item>                
            {
                docType === "video" ?
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
                        <Input />
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
                <InputNumber addonAfter="ms" controls={false} />
            </Form.Item>

            <Form.Item
                label="难度"
                name="difficulty"
            >
                <Select
                    placeholder="请至少选择一项主题"
                    options={[
                        {label: "困难", value: 2},
                        {label: "一般", value: 1},
                        {label: "简单", value: 0}
                    ]}
                />
            </Form.Item>

            <Form.Item>
                <Button type="primary" htmlType="submit">
                    Submit
                </Button>
            </Form.Item>

        </Form>
    )
}