import { Link, useLocation, useNavigate, useParams } from "react-router-dom";
import { ArrowLeftOutlined } from "@ant-design/icons";
import { Button, Input, message } from "antd";
import { useEffect, useState } from "react";
import { addUserTag, getTagInfo, tagUpdate } from "../../request/api/userTags";


export default function UserTagInfoPage(params) {

    const location = useLocation();
    const navigateTo = useNavigate();
    const [tagName, setTagName] = useState("");
    const [label, setLabel] = useState("");
    const [pageMode, setPageMode] = useState("");
    const [loading, setLoading] = useState(false);
    const [info, setInfo] = useState();
    const { id } = useParams();

    function addFunc() {
        // 添加
        addUserTag({name: tagName})
        .then(res => {
            message.success(res.msg);
            setTimeout(() => {
                navigateTo(`/dashboard/user/tag`);
            }, 500);
        })
        .catch(err => {
            message.error(err?.msg);
            setLoading(false);
        })
    }

    function modifyFunc() {
        // 修改
        tagUpdate({
            "id": info?.id, 
            "name": tagName
        })
        .then(res => {
            message.success(res.msg);
            setTimeout(() => {
                navigateTo(`/dashboard/user/tag`);
            }, 500);
        })
        .catch(err => {
            message.error(err?.msg);
            setLoading(false);
        })
    }
    
    function onFinish() {
        setLoading(true);
        switch (pageMode) {
            case "add":
                addFunc();
                break;
            default:
                modifyFunc();
                break;
        }
    }

    function changeTagName(value) {
        setTagName(value);
    }

    function init() {
        getTagInfo({tag_id: Number(id)})
        .then(res => {
            res.code === 0 ?
            setInfo({...res?.data})
            :
            message.error(res.msg)
        })
        .catch(err => {
            message.error(err?.msg);
        })
    }

    useEffect(() => {
        const arr = location.pathname.split('/');
        const mode = arr[arr.length-1];
        let str = "";
        
        switch (mode) {
            case "add":
                str = "标签管理/添加标签";
                break;
            default:
                str = "标签管理/修改标签"
                break;
        }
        setPageMode(mode);
        setLabel(str);
        id && init();
    },[])

    return (
        <div className="challenge-add challenge">
            <div className="tabel-title left-side">
                <Link to={`/dashboard/user/tag`}>
                    <ArrowLeftOutlined />
                </Link>
                <h2>{label}</h2>
            </div>
            {
                (!id || info) &&
                <div className="form-inner">
                    <div className="inner-item">
                        <p className="label">标签名</p>
                        <Input onChange={(e) => changeTagName(e.target.value)} defaultValue={info?.name || tagName} />
                    </div>
                    <div className="inner-btns">
                        <Button type="primary" htmlType="submit" loading={loading} onClick={() => onFinish()} >
                            保存
                        </Button>
                    </div>
                </div>
            }
        </div>
    )
}