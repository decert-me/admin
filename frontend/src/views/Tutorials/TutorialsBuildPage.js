import { Space, Table } from "antd";
import { useEffect, useState } from "react";
import { mock } from "../../mock";



export default function TutorialsBuildPage(params) {
    
    const { tutorials } = mock();
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
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status) => (
                <div className="build-status">
                    <div className={`point ${status === 0 ? "success" : "error"}`}></div>
                    <p>{status === 0 ? "打包成功" : "打包失败"}</p>
                </div>
            )
        },
        {
            title: '分支',
            dataIndex: 'branch',
            key: 'branch',
            render: (branch) => (
              <p>{branch ? branch : "main"}</p>
            )
        },
        {
            title: '文档目录',
            dataIndex: 'docPath',
            key: 'docPath',
            render: (docPath) => (
              <p>{docPath ? docPath : "/"}</p>
            )
        },
        {
            title: '教程文档CommitHash',
            dataIndex: 'commitHash',
            key: 'commitHash',
            render: (commitHash) => (
                <p>{commitHash}</p>
            )
        },
        {
            title: '创建时间',
            dataIndex: 'createDate',
            key: 'createDate',
            render: () => (
                <p>2023-01-01 11:59:00</p>
            )
        },
        {
            title: '最新打包时间',
            dataIndex: 'updateDate',
            key: 'updateDate',
            render: () => (
                <p>2023-01-01 11:59:00</p>
            )
        },
        {
          title: 'Action',
          key: 'action',
          render: () => (
            <Space size="middle">
              <a>打包</a>
              <a>打包记录</a>
            </Space>
          ),
        },
    ];

    function init() {
        data = tutorials.build;
        setData([...data]);
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="tutorials tutorials-build">
            <div className="tabel-title">
                <h2>打包管理</h2>
            </div>
            <Table columns={columns} dataSource={data} />
        </div>
    )
}