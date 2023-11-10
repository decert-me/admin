import { useEffect, useState } from "react";
import { Link, useLocation, useNavigate, useParams } from "react-router-dom"
import { Button, Form, Input, Select, Upload, message } from "antd";
import { UserOutlined, PlusOutlined, ArrowUpOutlined, DeleteOutlined, ArrowLeftOutlined } from '@ant-design/icons';
import { getUserInfo, registerUser, updateOtherInfo, updatePersonalInfo, uploadAvatar } from "../../request/api/user"; 
import { useAuth } from "../../hooks/useAuth";

const formItemLayout = {
    labelCol: {
      xs: { span: 24 },
      sm: { span: 6 },
    },
    wrapperCol: {
      xs: { span: 24 },
      sm: { span: 17 },
    },
};

const beforeUpload = (file) => {
    const isJpgOrPng = file.type === "image/jpeg" || file.type === "image/png";
    if (!isJpgOrPng) {
      message.error("You can only upload JPG/PNG file!");
    }
    const isLt2M = file.size / 1024 / 1024 < 5;
    if (!isLt2M) {
      message.error("Image must smaller than 5MB!");
    }
    return isJpgOrPng && isLt2M;
};


export default function PersonelEditPage(params) {
    
    const { type } = useParams();
    const { user } = useAuth();
    const [form] = Form.useForm();
    const location = useLocation();
    const navigator = useNavigate();
    const queryParams = new URLSearchParams(location.search);
    let [userAvatar, setUserAvatar] = useState();
    let [fields, setFields] = useState([]);

    // 清除头像
    function goUpload() {
        // 获取要点击的div
        const div = document.getElementById('Upload');
        // 创建一个新的点击事件
        const clickEvent = new MouseEvent('click', {
            // bubbles: true,
            // cancelable: true,
            view: window
        });

        // 分派点击事件
        div.dispatchEvent(clickEvent);
    }

    // 清除头像
    function clearAvatar() {
        userAvatar = null;
        setUserAvatar(userAvatar);
    }

    // 更新成员信息
    function updateUser(params) {
        const id = queryParams.get("id");
        const { userName } = form.getFieldsValue();
        // 判断当前身份 => 是否是超级管理员?
        if (user.authority.authorityId === "888") {
            updateOtherInfo({
                id, userName, headerimg: userAvatar
            })
            .then(res => {
                message.success(res.msg)
                navigator("/dashboard/personnel/list");
            })
        }else{
            updatePersonalInfo({
                userName, headerimg: userAvatar
            })
            .then(res => {
                message.success(res.msg)
                navigator("/dashboard/personnel/list");
            })
        }
    }

    // 添加成员信息
    function addUser(params) {
        const { userName, address, AuthorityId } = form.getFieldsValue();
        // 判断当前身份 => 是否是超级管理员?
        if (user.authority.authorityId === "888") {
            registerUser({
                username: userName, address, authorityId: AuthorityId, headerimg: userAvatar
            })
            .then(res => {
                message.success(res.msg);
                navigator("/dashboard/personnel/list");
            })
        }
    }

    // 更新 || 添加
    async function onFinish(params) {
        type === "add" ? addUser() : updateUser()
    }

    // 上传logo
    const handleChange = async(info) => {
        if (info.file.response?.code === 0) {
            userAvatar = info.file.response.data.url;
            setUserAvatar(userAvatar);
            message.success("上传成功!")
        }
    };

    // 编辑初始化
    function editInit(params) {
        const id = queryParams.get("id");
        getUserInfo({id})
        .then(res => {
            const { username, address, headerImg, authority } = res.data?.user;
            userAvatar = headerImg;
            setUserAvatar(userAvatar);
            fields = [
                { "name": ["userName"], "value": username },
                { "name": ["address"], "value": address },
                { "name": ["AuthorityId"], "value": authority.authorityId }
            ];
            setFields([...fields]);
        })
        
    }

    // 添加初始化
    function addInit(params) {
        form.setFieldValue("AuthorityId", "111");
    }

    useEffect(() => {
        type === "edit" && editInit();
        type === "add" && addInit();
    },[])

    return (
        <div className="personelEdit">
            {/* <h1>PersonelEditPage {type}</h1> */}
            <Link to={`/dashboard/personnel/list`}>
                <ArrowLeftOutlined />
            </Link>
            <h1>
                {
                    type === "edit" ? "编辑信息" : "添加成员"
                }
            </h1>
            <div className="container">
                <Form
                    className="form"
                    layout="vertical"
                    {...formItemLayout}
                    form={form}
                    fields={fields}
                >
                    <div style={{
                        position: "relative"
                    }}>
                        <Form.Item
                            label="姓名"
                            name="userName"
                            className="pl56"
                        >
                            <Input placeholder="用户名" />
                        </Form.Item>
                        <UserOutlined className="custom-icon" />
                    </div>

                    <Form.Item
                        label="钱包地址"
                        name="address"
                    >
                        <Input placeholder="钱包地址" disabled={type==="edit"} />
                    </Form.Item>

                    <Form.Item
                        label="角色"
                        name="AuthorityId"
                    >
                        <Select
                            className="custom-select"
                            disabled
                            options={[
                                { value: '888', label: '超级管理员' },
                                { value: '111', label: '管理员' }
                            ]}
                        />
                    </Form.Item>

                </Form>
                <div className="avatar">
                    {
                        userAvatar &&
                        <div className="img">
                            <img src={process.env.REACT_APP_BASE_URL+"/"+userAvatar} alt="" />
                            <div 
                                className="operate" 
                                style={{right: "51px"}}
                                onClick={goUpload}
                            ><ArrowUpOutlined /></div>
                            <div 
                                className="operate" 
                                style={{right: "14px"}}
                                onClick={clearAvatar}
                            ><DeleteOutlined /></div>
                        </div>
                    }
                        <div
                            style={{
                                display: userAvatar ? "none" : "block"
                            }}
                        >
                            <Upload
                                id="Upload"
                                name="file"
                                multiple={false}
                                maxCount={1}
                                className="custom-upload"
                                beforeUpload={beforeUpload}
                                onChange={handleChange}
                                customRequest={({
                                    data,
                                    file,
                                    onSuccess
                                }) => {
                                    const formData = new FormData();
                                    if (data) {
                                    Object.keys(data).forEach(key => {
                                        formData.append(key, data[key]);
                                    });
                                    }
                                    formData.append('file', file);
                                    uploadAvatar(formData)
                                    .then(res => {
                                    onSuccess(res, file);
                                    })
                                    .catch(err => {
                                    console.log(err);
                                    });
                                    return {
                                    abort() {
                                        console.log("upload progress is aborted.");
                                    }
                                    };
                                }}
                            >
                                <div className="img-info">
                                    <PlusOutlined />
                                    <p>上传图片</p>
                                </div>
                            </Upload>
                        </div>
                </div>
            </div>

            <Button type="primary" className="submit-btn" onClick={onFinish}>
                {
                    type === "edit" ? "保存" : "添加"
                }
            </Button>
        </div>
    )
}