import { Button, Form, Input, Popconfirm, Space, Table, message } from "antd";
import { PlusOutlined, DeleteOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getTagList, tagDeleteBatch } from "../../request/api/userTags";
import "./index.scss";

export default function UserTagsPage(params) {

  const navigateTo = useNavigate();
  const [selectedRowKeys, setSelectedRowKeys] = useState([]);
  const [data, setData] = useState([]);
  const [form, setForm] = useState({}); //  搜索
  let [pageConfig, setPageConfig] = useState({
    page: 0, pageSize: 50, total: 0
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
      title: "标签名",
      dataIndex: "name",
      ellipsis: {showTitle: false}
    },
    {
      title: "用户数量",
      dataIndex: "userNum",
    },
    {
      title: "创建时间",
      dataIndex: "createdAt",
      render: (time) => (
        time.indexOf("0001-01-01T") === -1 ?
        time.replace("T", " ").split(".")[0]
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
              onClick={() => navigateTo(`/dashboard/user/list/${tag.id}`)}
          >
            查看
          </Button>
          <Button
            type="link"
            className="p0"
              onClick={() => navigateTo(`/dashboard/user/tag/adduser/${tag.id}`)}
          >
            添加
          </Button>
          <Button
            type="link"
            className="p0"
              onClick={() => navigateTo(`/dashboard/user/tag/modify/${tag.id}`)}
          >
            编辑
          </Button>
        </Space>
      ),
    },
  ];

  function onSelectChange(newSelectedRowKeys) {
    setSelectedRowKeys(newSelectedRowKeys);
  }

  function deleteTags() {
    tagDeleteBatch({tag_ids: selectedRowKeys})
    .then(res => {
        message.success(res.msg);
        getList();
    })
    .catch(err => {
        message.error(err?.msg);
    })
  }

  function onFinish(params) {
    setForm({...params});
  }

  async function getList(page) {
    if (page) {
      pageConfig.page = page;
      setPageConfig({ ...pageConfig });
    }
    // 获取教程列表
    let res = await getTagList({ ...pageConfig, search_val: form.tag });

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

  useEffect(() => {
    pageConfig = {
      page: 0,
      pageSize: 50,
      total: 0,
    };
    setPageConfig({ ...pageConfig });
    init();
  }, [form]);

  return (
    <div>
      <div className="tabel-title">
        <h2>标签管理</h2>
      </div>
      <div>
        <div className="operat">
          <div className="btns">
            <Button
              icon={<PlusOutlined />}
              onClick={() => navigateTo(`/dashboard/user/tag/add`)}
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
              onFinish={onFinish}
          >
            <Form.Item label="标签" name="tag">
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
        }} 
      />
    </div>
  );
}
