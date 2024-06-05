import { Link, useParams } from "react-router-dom";
import { ArrowLeftOutlined } from "@ant-design/icons";
import { Button, Input, message } from "antd";
import { useEffect, useState } from "react";
import { getTagInfo, getUsersList, tagUserUpdate } from "../../request/api/userTags";
import { useRequest, useUpdateEffect } from "ahooks";

export default function UserTagAddPage(params) {

    const { id } = useParams();
    const [info, setInfo] = useState();
    const [loading, setLoading] = useState(false);
    const [searchItem, setSearchItem] = useState({});
    const [key, setKey] = useState(0);
    let [formItem, setFormItem] = useState({addr: "", name: ""});

    const { run } = useRequest(changeAddr, {
        debounceWait: 500,
        manual: true,
    });
    
    async function onFinish() {
        setLoading(true)
        if (!searchItem?.user_id) {
            message.error("该地址不存在，请确认后再添加!")
            setLoading(false);
            return
        }
        await tagUserUpdate({
            tag_id: Number(id),
            user_id: Number(searchItem.user_id),
            name: formItem.name
        })
        .then(res => {
            message.success(res.msg);
            setFormItem({addr: "", name: ""});
            setSearchItem({});
            setKey(key+1);
        })
        .catch(err => {
            message.error(err?.msg);
        })
        setLoading(false);
    }

    function changeForm(key, value) {
        formItem[key] = value;
        setFormItem({...formItem});
    }

    function changeAddr(addr) {
        let searchItem = {};
        getUsersList({
            "page": 1,
            "pageSize": 10, //每页数量
            "search_tag": "", 
            "search_address": addr
        })
        .then(res => {
            if (res.code === 0 && res.data.list.length > 0 && res.data.list[0].address === addr) {
                searchItem = res.data.list[0];
            }
            formItem.addr = addr;
            setFormItem({...formItem})
            setSearchItem({...searchItem});
        })
        .catch(err => {
            message.error(err.msg);
        })
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

    useUpdateEffect(() => {
        const { name: nickname} = searchItem;
        const {addr, name} = formItem;
        const obj = {addr, name: nickname || name}
        setFormItem({...obj})
    },[searchItem])

    useEffect(() => {
        init();
    },[])

    return(
        <div className="challenge-add challenge" key={key}>
            <div className="tabel-title left-side">
                <Link to={`/dashboard/user/tag`}>
                    <ArrowLeftOutlined />
                </Link>
                <h2>标签管理/用户列表/添加用户</h2>
            </div>
            <div className="form-inner">
                <div className="inner-item">
                    <p className="label">标签</p>
                    <Input disabled value={info?.name} />
                </div>
                <div className="inner-item">
                    <p className="label">挑战者地址</p>
                    <Input onChange={(e) => run(e.target.value)} />
                </div>
                <div className="inner-item">
                    <p className="label">昵称</p>
                    <Input value={formItem.name} onChange={(e) => changeForm("name", e.target.value)} />
                </div>
                <div className="inner-btns">
                    <Button type="primary" htmlType="submit" loading={loading} onClick={() => onFinish()} >
                        保存
                    </Button>
                </div>
            </div>
        </div>
    )
}