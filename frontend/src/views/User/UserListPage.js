import { Button, Form, Input, Space, Table, message } from "antd";
import { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom"
import { getUsersList } from "../../request/api/userTags";
import { useUpdateEffect } from "ahooks";



export default function UserListPage(params) {
    
    const navigateTo = useNavigate();
    const location = useLocation();
    const [formProps] = Form.useForm();
    const [data, setData] = useState([]);
    const [search_key, setSearch_key] = useState(""); //  搜索
    const [form, setForm] = useState({tag: ""}); //  搜索
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 10, total: 0
    });

    const columns = [
        {
          title: "ID",
          dataIndex: "user_id",
        },
        {
          title: "用户地址",
          dataIndex: "address",
          render: (address) => (
            <a target="_blank" href={`${process.env.REACT_APP_LINK_URL || "https://decert.me"}/user/${address}`}>{address}</a>
          )
        },
        {
          title: "昵称",
          dataIndex: "name"
        },
        {
          title: "标签",
          dataIndex: "tags",
          ellipsis: true
        },
        {
            title: "创建时间",
            dataIndex: "created_at",
            render: (time) => (
              time.indexOf("0001-01-01T") === -1 ?
              time.replace("T", " ").split(".")[0].split("+")[0]
              :"-"
          )
        },
        {
          title: "操作",
          key: "action",
          render: (_, tag) => (
            <Space size="middle">
              <Button
                type="link"
                className="p0"
                  onClick={() => navigateTo(`/dashboard/user/tag/modifyuser/${tag.address}`)}
              >
                编辑
              </Button>
            </Space>
          ),
        },
    ];

    function onFinish(params) {
        setForm({...params});
    }

    async function getList(page) {
        if (page) {
          pageConfig.page = page;
          setPageConfig({ ...pageConfig });
        }
        // 获取教程列表
        let res = await getUsersList({ ...pageConfig, search_address: form?.challenger, search_tag: form?.tag });
    
        if (res.code === 0) {
          const list = res.data.list || [];
          // 添加key
          list.forEach((ele) => {
            ele.key = ele.id;
          });
          setData([...list]);
          pageConfig.total = res.data.total;
          setPageConfig({ ...pageConfig });
        } else {
          message.success(res.msg);
        }
    }

    function init(params) {
        pageConfig.page += 1;
        setPageConfig({ ...pageConfig });
        getList();
    }

    // getTagUserList
    useEffect(() => {
        pageConfig = {
          page: 0,
          pageSize: 10,
          total: 0,
        };
        setPageConfig({ ...pageConfig });
        init();
    }, [form]);

    return (
        <div>
            <div className="tabel-title">
                <h2>用户管理</h2>
            </div>
            <div>
                <div className="operat">
                    <Form
                        name="horizontal_login"
                        layout="inline"
                        form={formProps}
                        onFinish={onFinish}
                    >
                        <Form.Item label="标签" name="tag">
                            <Input />
                        </Form.Item>
                        <Form.Item label="挑战者" name="challenger">
                            <Input />
                        </Form.Item>
                        <Form.Item shouldUpdate>
                        {() => (
                            <Button type="primary" htmlType="submit">
                            搜索
                            </Button>
                        )}
                        </Form.Item>
                    </Form>
                </div>
            </div>
            <Table 
                columns={columns} 
                dataSource={data}         
                pagination={{
                    current: pageConfig.page, 
                    total: pageConfig.total, 
                    pageSize: pageConfig.pageSize, 
                    onChange: (page) => getList(page)
                }} />
        </div>
    )
}