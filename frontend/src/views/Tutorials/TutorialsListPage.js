import { Button, Space, Switch, Table, Tag } from "antd";
import { useEffect, useState } from "react";
import { mock } from "../../mock";
import { Link, useNavigate } from "react-router-dom";
import {
    VideoCameraOutlined,
    ReadOutlined
  } from '@ant-design/icons';

export default function TutorialsListPage(params) {
    
    const { tutorials } = mock();
    const navigateTo = useNavigate();
    const table = require("./category_tabel.json");
    let [data, setData] = useState([]);

    const columns = [
        {
          title: '封面图',
          dataIndex: 'img',
          key: 'img',
          render: (img) => (
            <img src={img} alt="" style={{height: "40px"}} />
          )
        },
        {
          title: '标题',
          dataIndex: 'label',
          key: 'label',
          render: (text) => (
            <p className="tabel-item-title newline-omitted">{text}</p>
          )
        },
        {
          title: '上架状态',
          key: 'status',
          dataIndex: 'status',
          render: (status) => (
              <Switch checkedChildren="已上架" unCheckedChildren="待上架" defaultChecked={status == 0 ? true : false} />
          )
        },
        {
          title: '分类-主题',
          dataIndex: 'category',
          key: 'category',
          render: (category, tutorial) => (
            <>
            {
              category.map(tag => 
                  <Tag color="geekblue" key={tag}>
                      {table.category[tag]}
                  </Tag>    
              )
            }
            {
              tutorial.theme.map(tag => 
                  <Tag color="green" key={tag}>
                      {table.theme[tag]}
                  </Tag>    
              )
            }
            </>
          )
        },
        {
          title: '语言',
          key: 'language',
          dataIndex: 'language',
          render: (language) => (
              <Tag color={language === "zh" ? "#2db7f5" : "#87d068"}>
                  {language === "zh" ? "中文" : "英文"}
              </Tag>    
          )
        },
        {
            title: '媒体类型',
            key: 'docType',
            dataIndex: 'docType',
            render: (docType) => (
              <div style={{lineHeight: "20px"}}>
                <Tag icon={docType === "video" ? <VideoCameraOutlined /> : <ReadOutlined />} color="default">
                    {docType === "video" ? "视频" : "文章"}
                </Tag>
              </div>
            )
        },
        {
          title: 'Action',
          key: 'action',
          render: (_, tutorial) => (
            <Space size="middle">
              <Link to={`/dashboard/tutorials/modify/${tutorial.id}`}>修改</Link>
              <a>删除</a>
            </Space>
          ),
        },
    ];

    function init() {
        data = tutorials.list;
        setData([...data]);
    }


    useEffect(() => {
        init();
    },[])

    return (
        <div className="tutorials tutorials-list">
          <div className="tabel-title">
            <h2>教程列表</h2>
            <Button 
              type="primary"
              onClick={() => navigateTo("/dashboard/tutorials/add")}
            >创建教程</Button>
          </div>
            <Table columns={columns} dataSource={data} />
        </div>
    )
}