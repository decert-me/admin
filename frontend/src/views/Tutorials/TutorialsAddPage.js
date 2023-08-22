import "./index.scss"
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import { Button, Form, Input, InputNumber, Select, Upload, message, Space } from 'antd';
import {
    PlusOutlined,
    MinusCircleOutlined,
    ArrowLeftOutlined,
    MenuOutlined
  } from '@ant-design/icons';
import { useEffect, useState } from 'react';
import { UploadProps } from '../../utils/props';
import { createTutorial } from '../../request/api/tutorial';
import { Link, useNavigate } from 'react-router-dom';
import { getLabelList } from '../../request/api/tags';
import { getYouTubePlayList } from '../../request/api/public';
import { useUpdateEffect } from 'ahooks';
import { getQuest } from "../../request/api/quest";
const { TextArea } = Input;



export default function TutorialsAddPage(params) {
    
    const [form] = Form.useForm();
    const videoCategory = Form.useWatch("videoCategory", form);
    const videoUrl = Form.useWatch("url", form);
    const navigateTo = useNavigate();

    const [loading, setLoading] = useState(false);
    const [parseLoading, setParseLoading] = useState(false);
    let [category, setCategory] = useState();     //  类别 选择器option
    let [lang, setLang] = useState();     //  语种 选择器option
    let [doctype, setDoctype] = useState("doc");
    let [videoList, updateVideoList] = useState([]);

    function parseVideoList() {
        const link = form.getFieldValue("url");
        if (!link) {
            message.error("请输入正确的视频地址!")
            return
        }
        setParseLoading(true);
        getYouTubePlayList({link})
        .then(res => {
            setParseLoading(false)
            if (res.code === 0) {
                message.success(res.msg);
                videoList = res.data;
                updateVideoList([...videoList]);
            }
        })
        .catch(err => {
            setParseLoading(false)
            message.error(err);
        })
    }

    function handleOnDragEnd(result) {
        if (!result.destination) return;
    
        const items = Array.from(videoList);
        const [reorderedItem] = items.splice(result.source.index, 1);
        items.splice(result.destination.index, 0, reorderedItem);
    
        updateVideoList(items);
    }

    const onFinish = async(values) => {
        setLoading(true);
        const {
            repoUrl, label, catalogueName, docType, desc, 
            challenge, branch, docPath, commitHash, url, videoCategory, videoItems,
            category, language, difficulty, estimateTime
        } = values;
        const img = values.img.file.response.data.hash;

        let flag = true;
        if (challenge) {
            await getQuest({id: Number(challenge)})
            .then(res => {
                if (res.status !== 0) {
                    flag = false;
                    return;
                }
            })
        }
        if (!flag) {
            // 终止
            message.error("请输入正确的挑战编号!")
            setLoading(false);
            return
        }

        if (doctype === "doc") {
            const obj = {
                repoUrl, label, catalogueName, docType, img, desc, 
                challenge, branch, docPath, commitHash,
                category, language, difficulty, estimateTime
            }
            create(obj)
        }else{
            const obj = {
                url, label, catalogueName, img, desc, 
                challenge, videoCategory,
                category, language, difficulty, estimateTime, docType: "video"
            }
            if (videoCategory === "bilibili") {
                create({...obj, video: videoItems})
            }else{
                create({...obj, video: videoList})
            }
        }
    };

    function create(obj) {
        createTutorial(obj)
        .then(res => {
            if (res.code === 0) {
                message.success(res.msg);
                setTimeout(() => {
                    navigateTo("/dashboard/tutorials/list");
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

    function init() {
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

    useUpdateEffect(() => {
        if (videoList.length !== 0) {
            videoList = [];
            updateVideoList([...videoList]);
        }
    },[videoUrl])

    useUpdateEffect(() => {
        if (videoCategory === "bilibili") {
            videoList = [];
            updateVideoList([...videoList]);
        }
    },[videoCategory])

    return (
        <div className="tutorials-add tutorials">
            <Link to={`/dashboard/tutorials/list`}>
                <ArrowLeftOutlined />
            </Link>
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
                <Form.Item
                    label="视频地址"
                    name="url"
                    rules={[{
                        required: true,
                        message: '请输入视频地址!',
                    }]}
                    hidden={doctype !== "video"}
                >
                    <Space.Compact
                        style={{
                            width: '100%',
                        }}
                    >
                    <Input />
                    {
                        videoCategory === "youtube" &&
                        <Button type="primary" onClick={() => parseVideoList()} loading={parseLoading} >解析</Button>
                    }
                    </Space.Compact>
                </Form.Item>
                {
                    doctype === "video" ?
                    <>
                        {
                            videoList.length !== 0 &&
                            <Form.Item label="视频排序">
                                <DragDropContext onDragEnd={handleOnDragEnd}>
                                    <Droppable droppableId="characters">
                                        {(provided) => (
                                        <ul className="characters video-list" {...provided.droppableProps} ref={provided.innerRef}>
                                            {videoList.map(({id, img, label, url}, index) => {
                                            return (
                                                <Draggable key={id} draggableId={id} index={index}>
                                                {(provided) => (
                                                    <li ref={provided.innerRef} {...provided.draggableProps} {...provided.dragHandleProps}>
                                                    <p className="number">{index + 1}</p>
                                                    <img src={img} />
                                                    <div>
                                                        <p className="newline-omitted">{label}</p>
                                                        <p className="newline-omitted">{url}</p>
                                                    </div>
                                                    </li>
                                                )}
                                                </Draggable>
                                            );
                                            })}
                                            {provided.placeholder}
                                        </ul>
                                        )}
                                    </Droppable>
                                </DragDropContext>
                            </Form.Item>
                        }
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
                                    {label: "YouTube", value: "youtube"},
                                    {label: "Bilibili", value: "bilibili"}
                                ]}
                            />
                        </Form.Item>
                        {
                            videoCategory === "bilibili" &&
                            <Form.Item
                                label="视频列表"
                            >
                                <Form.List name="videoItems" >
                                    {(fields, { add, remove }, { errors }) => (
                                    <>
                                        {fields.map((field, index) => (
                                        <Space
                                            key={index}
                                            align="baseline"
                                            className="bte"
                                        >
                                            
                                            <Form.Item
                                                {...field}
                                                validateTrigger={['onChange', 'onBlur']}
                                                name={[field.name, 'label']}
                                                rules={[
                                                    {
                                                    required: true,
                                                    whitespace: true,
                                                    message: "请输入视频标题！",
                                                    },
                                                ]}
                                                noStyle
                                            >
                                                <Input placeholder="视频标题"/>
                                            </Form.Item>
                                            <Form.Item
                                                {...field}
                                                name={[field.name, 'code']}
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
                                                <Input placeholder="bilibili的嵌入代码" />
                                            </Form.Item>
                                            <MinusCircleOutlined
                                                className="dynamic-delete-button"
                                                onClick={() => remove(field.name)}
                                            />
                                        </Space>

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
                                            {label: "Docusaurus", value: "docusaurus"},
                                            {label: "GitBook", value: "gitbook"},
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
                    label="目录名"
                    name="catalogueName"
                    rules={[{
                        required: true,
                        message: '请输入目录名!',
                    }]}
                >
                    <Input />
                </Form.Item>

                <Form.Item
                    label="分类"
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
                    rules={[{
                        required: true,
                        message: '请输入语种!',
                    }]}
                >
                    <Select
                        placeholder="请至少选择一项语种"
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
                        添加教程
                    </Button>
                </Form.Item>

            </Form>
        </div>
    )
}