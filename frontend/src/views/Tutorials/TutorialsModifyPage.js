import {
    ArrowLeftOutlined,
    PlusOutlined,
    MinusCircleOutlined
  } from '@ant-design/icons';
import { useEffect, useState } from 'react';
import { Link, useNavigate, useParams } from 'react-router-dom';
import { Button, Form, Input, InputNumber, Select, Upload, message } from 'antd';
import { getTutorial, updateTutorial } from '../../request/api/tutorial';
import { getLabelList } from '../../request/api/tags';
import { UploadProps } from '../../utils/props';
const { TextArea } = Input;



export default function TutorialsModifyPage(params) {
    
    const { id } = useParams();
    const [form] = Form.useForm();
    const navigateTo = useNavigate();
    const videoCategory = Form.useWatch("videoCategory", form);
    let [fields, setFields] = useState([]);
    let [tutorial, setTutorial] = useState();
    let [category, setCategory] = useState();     //  类别 选择器option
    let [lang, setLang] = useState();     //  语种 选择器option
    let [doctype, setDoctype] = useState("doc");
    const [loading, setLoading] = useState(false);
    
    const onFinish = (values) => {
        // console.log('Success:', values);
        // console.log();
        // updateTutorial(values)
        setLoading(true);
        const {
            repoUrl, label, catalogueName, docType, desc, 
            challenge, branch, docPath, commitHash, 
            category, language, difficulty,
        } = values;
        const img = values.img?.file ? values.img.file.response.data.hash : tutorial.img;
        if (doctype === "doc") {
            const obj = {
                repoUrl, label, catalogueName, docType, img, desc, 
                challenge, branch, docPath, commitHash,
                category, language, difficulty
            }
            addArticle(obj)
        }else{
            // addVideo(values)
        }

    };

    function addArticle(obj) {
        updateTutorial({...obj, id: Number(id)})
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

    async function init() {
        await getTutorial({id: Number(id)})
        .then(res => {
            if (res.code === 0) {
                tutorial = res.data
                setTutorial({...tutorial});
            }else{
                navigateTo(-1);
            }
        })
        .catch(err => {
            navigateTo(-1);
        })
        fields = [
            {
                name: ['label'],
                value: tutorial.label
            },
            {
                name: ['desc'],
                value: tutorial.desc
            },
            {
                name: ['img'],
                value: "https://ipfs.decert.me/"+tutorial.img
            },
            {
                name: ['challenge'],
                value: tutorial?.challenge
            },
            {
                name: ['category'],
                value: tutorial?.category
            },
            {
                name: ['language'],
                value: tutorial?.language
            },
            {
                name: ['estimateTime'],
                value: tutorial?.estimateTime
            },
            {
                name: ['difficulty'],
                value: tutorial?.difficulty
            },

            // 文档
            {
                name: ['repoUrl'],
                value: tutorial?.repoUrl
            },
            {
                name: ['branch'],
                value: tutorial?.branch
            },
            {
                name: ['docPath'],
                value: tutorial?.docPath
            },
            {
                name: ['commitHash'],
                value: tutorial?.commitHash
            },

            // 视频
            {
                name: ['url'],
                value: tutorial?.url
            },
            {
                name: ['videoCategory'],
                value: tutorial?.videoCategory
            },
        ]
        setFields([...fields]);
        optionsInit()
    }

    async function getOption(type) {
        return await getLabelList(type)
        .then(res => {
            if (res.code === 0) {
                return res.data ? res.data : [];
            }
        })
    }

    // 选择器option初始化
    async function optionsInit(params) {
        // 类别初始化
        let categoryOption = await getOption({type: "category"});
        categoryOption.forEach(ele => {
            ele.label = ele.Chinese;
            ele.value = ele.ID
        })
        category = categoryOption;
        setCategory([...category])
    
        // 语种
        let langOption = await getOption({type: "language"});
        langOption.forEach(ele => {
            ele.label = ele.Chinese;
            ele.value = ele.ID
        })
        lang = langOption;
        setLang([...lang]);
    }

    useEffect(() => {
        init()
    },[])

    return (
        <div className="tutorials-modify">
            <Link to={`/dashboard/tutorials/list`}>
                <ArrowLeftOutlined />
            </Link>
            {
                tutorial &&
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
                            defaultFileList={[{
                                uid: '-1',
                                name: 'image.png',
                                status: 'done',
                                url: "https://ipfs.decert.me/"+tutorial.img,
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
                        name="estimateTime"
                    >
                        <InputNumber addonAfter="min" controls={false} />
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
                            修改教程
                        </Button>
                    </Form.Item>

                </Form>
            }
        </div>
    )
}