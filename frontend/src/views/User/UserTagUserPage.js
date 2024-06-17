import { Button, Form, Input, Popconfirm, Space, Table, message } from "antd";
import { PlusOutlined, DeleteOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { useLocation, useNavigate, useParams } from "react-router-dom"
import { getTagInfo, getTagUserList, getUsersList, tagUserDeleteBatch } from "../../request/api/userTags";
import { useUpdateEffect } from "ahooks";



export default function UserTagUserPage(params) {
    
    const {tagid} = useParams();
    const navigateTo = useNavigate();
    const [formProps] = Form.useForm();
    const [label, setLabel] = useState("");
    const [data, setData] = useState([]);
    const [selectedRowKeys, setSelectedRowKeys] = useState([]);
    const [form, setForm] = useState({}); //  搜索
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 10, total: 0
    });

    const rowSelection = {
        selectedRowKeys,
        onChange: onSelectChange,
    };

    const columns = [
        {
          title: "ID",
          dataIndex: "id",
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
          dataIndex: "nickname",
          ellipsis: true
        },
        {
            title: "创建时间",
            dataIndex: "createdAt",
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

    function onSelectChange(newSelectedRowKeys) {
        setSelectedRowKeys(newSelectedRowKeys);
    }

    function deleteTags() {
        tagUserDeleteBatch({user_ids: selectedRowKeys, tag_id: Number(tagid)})
        .then(res => {
            message.success(res.msg);
            getList();
        })
        .catch(err => {
            message.error(err?.msg);
        })
    }

    async function getList(page) {
        if (page) {
          pageConfig.page = page;
          setPageConfig({ ...pageConfig });
        }
        // 获取教程列表
        let res = await getTagUserList({ ...pageConfig, search_val: form?.challenger, tag_id: Number(tagid) });
    
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

    async function init(params) {
        pageConfig.page += 1;
        setPageConfig({ ...pageConfig });
        getList();

        const res = await getTagInfo({tag_id: Number(tagid)})
        if (res.code === 0) {
            setLabel(res.data.name);
        }
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

    useUpdateEffect(() => {
        navigateTo(0);
    },[tagid])

    return (
        <div>
            <div className="tabel-title">
                <h2>用户管理/{label}</h2>
            </div>
            <div>
                <div className="operat">
                    <div className="btns">
                        <Button
                        icon={<PlusOutlined />}
                        onClick={() => navigateTo(`/dashboard/user/tag/adduser/${tagid}`)}
                        />
                        <Popconfirm
                            title="确认删除?"
                            onConfirm={deleteTags}
                            okText="确认"
                            cancelText="取消"
                        >
                            <Button icon={<DeleteOutlined />} />
                        </Popconfirm>
                    </div>
                    <Form
                        name="horizontal_login"
                        layout="inline"
                        form={formProps}
                        onFinish={onFinish}
                    >
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
                rowSelection={rowSelection} 
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