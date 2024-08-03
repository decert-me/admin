import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom"
import { getTagList, getUsersInfo, getUsersList, updateUsersInfo } from "../../request/api/userTags";
import { Button, Input, Select, message } from "antd";


export default function UserTagModify() {
    
    const {address} = useParams();
    const navigateTo = useNavigate();
    const [tags, setTags] = useState([]);
    const [userTags, setUserTags] = useState([]);
    const [userInfo, setUserInfo] = useState();
    const [nickname, setNickName] = useState("");
    const [userid, setUserId] = useState("");
    const [loading, setLoading] = useState(false);

    function onFinish() {
        setLoading(true);
        updateUsersInfo({
            "user_id": userid,
            "name": nickname, 
            "tag_ids": userTags
        })
        .then(res => {
            message.success(res.msg);
            setTimeout(() => {
                navigateTo(-1);
            }, 500);
        })
        .catch(err => {
            message.error(err?.msg);
            setLoading(false);
        })
    }

    function changeName(e) {
        setNickName(e.target.value);
    }

    function handleChange(params) {
        setUserTags([...params]);
    }

    async function init() {
        // 获取user_id
        const res = await getUsersList({
            "page": 1,
            "pageSize": 50, 
            "search_tag": "", 
            "search_address": address
        })
        if (res.code !== 0 || res.data.list.length === 0) {
            return
        }
        setUserId(res.data.list[0].user_id);
        Promise.all([
            // 获取user tags
            getUsersInfo({user_id: res.data.list[0].user_id})
            .then(res => {
                if (res.code === 0) {
                    const list = res.data.tag || [];
                    const select = list.map(res => res.id)
                    setUserTags([...select])
                }
            }), 
            // 获取全部tags
            await getTagList({})
            .then(res => {
                if (res.code === 0) {
                    const list = res.data.list || [];
                    const options = list.map(res => {
                        return {
                            label: res.name,
                            value: res.id
                        }
                    })
                    setTags([...options])
                }
            })
        ]).then(() => {
            setUserInfo({...res.data.list[0]});
            setNickName(res.data.list[0]?.name || "")
        })
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div>
            <div className="tabel-title">
                <h2>用户管理/编辑用户信息</h2>
            </div>
            {
                userInfo &&
                <div className="form-inner">
                    <div className="inner-item">
                        <p className="label">挑战者地址</p>
                        <Input disabled value={address} />
                    </div>
                    <div className="inner-item">
                        <p className="label">昵称</p>
                        <Input defaultValue={nickname} onChange={changeName} />
                    </div>
                    <div className="inner-item">
                        <p className="label">标签</p>
                        <Select
                            mode="multiple"
                            allowClear
                            style={{
                                width: '100%',
                            }}
                            value={userTags}
                            onChange={handleChange}
                            options={tags}
                        />
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