import { Button, Popconfirm, Space, Table, message } from "antd";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { deleteLabel, getLabelList } from "../../request/api/tags";



export default function TagsPage(params) {
    
    const navigateTo = useNavigate();
    let [categoryData, setCategoryData] = useState([]);
    let [langData, setlangData] = useState([]);

    const columns = [
        {
            title: '中文',
            dataIndex: 'Chinese',
            key: 'Chinese',
            render: (Chinese) => (
              <p>{Chinese}</p>
            )
        },
        {
            title: '英文',
            dataIndex: 'English',
            key: 'English',
            render: (English) => (
              <p>{English}</p>
            )
        },
        {
            title: '权重',
            dataIndex: 'Weight',
            key: 'Weight',
            render: (Weight) => (
                <p>{Weight}</p>
              )
        },
        {
            title: '创建时间',
            dataIndex: 'CreatedAt',
            key: 'CreatedAt',
            render: (CreatedAt) => (
                <p>{CreatedAt.replace("T", " ").split(".")[0]}</p>
            )
        },
        {
            title: '最新更新时间',
            dataIndex: 'UpdatedAt',
            key: 'UpdatedAt',
            render: (UpdatedAt) => (
                <p>{UpdatedAt.replace("T", " ").split(".")[0]}</p>
            )
        },
        {
          title: 'Action',
          key: 'action',
          render: (_, tags) => (
            <Space size="middle">
                <Link to={`/dashboard/tags/modify/${tags.type}/${tags.ID}?weight=${tags.Weight}&english=${tags.English}&chinese=${tags.Chinese}`}>
                    修改
                </Link>
                <Popconfirm
                  title="删除标签"
                  description="确定要删除该标签吗?"
                  onConfirm={() => deleteTags({type: tags.type, id: tags.ID})}
                  okText="确定"
                  cancelText="取消"
                >
                  <a>删除</a>
                </Popconfirm>
            </Space>
          ),
        },
    ];

    async function deleteTags(obj) {
        await deleteLabel(obj)
        .then(res => {
            if (res.code === 0) {
                message.success(res.msg);
            }
        })
        .catch(err => {
            message.error(err);
        })

        if (obj.type === "language") {            
            langData = await getData({type: obj.type});
            setlangData([...langData]);
        }else{
            categoryData = await getData({type: obj.type});
            setCategoryData([...categoryData]);
        }
    }

    async function getData(obj) {
        return await getLabelList(obj)
        .then(res => {
            if (res.code === 0) {
            const list = res.data;
            const data = list ? list : [];
            // 添加key
            data.forEach(ele => {
                ele.key = ele.ID
                ele.type = obj.type
            })
            return data
            }else{
                message.success(res.msg);
            }
        })
        .catch(err => {
            message.error(err)
        })
    }

    async function init() {
        langData = await getData({type: "language"});
        setlangData([...langData]);
        categoryData = await getData({type: "category"});
        setCategoryData([...categoryData]);
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="tags">
            <Button
                type="primary"
                onClick={() => navigateTo("/dashboard/tags/add")}
            >创建标签</Button>  
            
            <h2>分类</h2>
            <Table columns={columns} dataSource={categoryData} pagination={false} />
            <h2>语种</h2>
            <Table columns={columns} dataSource={langData} pagination={false} />
        </div>
    )
}