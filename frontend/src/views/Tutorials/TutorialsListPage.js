import { Space, Table, Tag } from "antd";
import { useEffect, useState } from "react";
import { mock } from "../../mock";
import {
    VideoCameraOutlined,
    ReadOutlined
  } from '@ant-design/icons';

export default function TutorialsListPage(params) {
    
    const { tutorials } = mock();
    const table = require("./category_tabel.json");
    let [data, setData] = useState([]);

    const columns = [
        {
          title: '标题',
          dataIndex: 'label',
          key: 'label',
          render: (text, tutorial) => (
            <a href={tutorial.repoUrl} target="_blank">{text}</a>
          )
        },
        {
          title: '图片',
          dataIndex: 'img',
          key: 'img',
          render: (img) => (
            <img src={img} alt="" style={{height: "100px"}} />
          )
        },
        {
          title: '分类',
          dataIndex: 'category',
          key: 'category',
          render: (category) => (
            category.map(tag => 
                <Tag color="geekblue" key={tag}>
                    {table[tag]}
                </Tag>    
            )
          )
        },
        {
          title: '主题',
          key: 'theme',
          dataIndex: 'theme',
          render: (theme) => (
            theme.map(tag => 
                <Tag color="green" key={tag}>
                    {table[tag]}
                </Tag>    
            )
          )
        },
        {
            title: '媒体类型',
            key: 'docType',
            dataIndex: 'docType',
            render: (docType) => (
                <Tag icon={docType === "video" ? <VideoCameraOutlined /> : <ReadOutlined />} color="default">
                    {docType}
                </Tag>    
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
          title: 'Action',
          key: 'action',
          render: (_, record) => (
            <Space size="middle">
              <a>修改 {record.name}</a>
              <a>隐藏</a>
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
        // <h1>TutorialsListPage</h1>
        <div className="tutorials-list">
            <Table columns={columns} dataSource={data} />
        </div>
    )
}