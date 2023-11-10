import { useEffect, useState } from "react"
import { Link, useNavigate } from "react-router-dom";
import { Button, Popconfirm, Space, Table, message } from "antd";
import { PlusOutlined } from '@ant-design/icons';
import { deleteUser, getUserList } from "../../request/api/user";
import "./index.scss";
import { useAuth } from "../../hooks/useAuth";


export default function PersonelListPage(params) {

    const { user } = useAuth();
    const navigator = useNavigate();
    let [data, setData] = useState([]);
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 10, total: 0
    });

    const columns = [
        {
            title: '成员名称',
            dataIndex: 'username',
            key: 'username',
            render: (username, record) => (
                <Space size="middle">
                    {
                        record.headerImg ?
                        <div className="img">
                            <img src={process.env.REACT_APP_BASE_URL+"/"+record.headerImg} alt="" />
                        </div>
                        :
                        <div className="img">{username[0].toUpperCase()}</div>
                    }
                    <p>{username}</p>
                </Space>
            )
        },
        {
            title: '钱包地址',
            dataIndex: 'address',
            key: 'address',
            // render: (address) => (
            //     <p>{address.substring(0,5) + "..." + address.substring(38,42)}</p>
            // )
        },
        {
            title: '角色',
            dataIndex: 'authority',
            key: 'authority',
            render: ({authorityName}) => (
                <p>{authorityName}</p>
            )
        },
        {
            title: '操作',
            key: 'action',
            render: (_, record) => (
              <Space size="middle">
                {
                    (user.authority.authorityId === "888" || 
                    user.address === record.address) &&
                    <Link 
                        to={`/dashboard/personnel/edit?id=${record.id}`}
                    >
                        编辑
                    </Link>
                }
                {
                    user.authority.authorityId === "888" &&
                    <Popconfirm
                        title="删除成员"
                        description="确定要删除该成员吗?"
                        onConfirm={() => goDelete({id: record.id})}
                        okText="确定"
                        cancelText="取消"
                    >
                        <a>删除</a>
                    </Popconfirm>
                }
              </Space>
            ),
        }
    ];

    // 删除
    async function goDelete({id}) {
        await deleteUser({id})
        .then(res => {
            message.success(res.msg);
        })
        getList(1);
    }

    // 获取用户列表
    function getList(page) {
        if (page) {
            pageConfig.page = page;
            setPageConfig({...pageConfig});
        }
        getUserList({...pageConfig})
        .then(res => {
            const list = res.data.list;
            data = list ? list : [];
            // 添加key
            data.forEach(ele => {
                ele.key = ele.id
            })
            setData([...data]);
            pageConfig.total = res.data.total;
            setPageConfig({...pageConfig});
        })
    }

    function init(params) {
        pageConfig.page += 1;
        setPageConfig({...pageConfig});
        getList();
    }

    useEffect(() => {
        init();
    },[]);
    
    return (
        <div className="personel">
            <div className="personel-btn">
                <Button
                    type="primary"
                    onClick={() => navigator("/dashboard/personnel/add")}
                    icon={<PlusOutlined />}
                >添加</Button>  
            </div>
            
            <Table
                columns={columns} 
                dataSource={data} 
                pagination={{
                    current: pageConfig.page, 
                    total: pageConfig.total, 
                    pageSize: pageConfig.pageSize, 
                    onChange: (page) => getList(page)
                }} 
            />
        </div>
    )
}