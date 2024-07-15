import "./index.scss"
import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom"
import { getPackLog, getTutorial } from "../../request/api/tutorial";
import {
    ArrowLeftOutlined,
  } from '@ant-design/icons';
import { Table } from "antd";


export default function TutorialsBuildLogPage(params) {

    const { id } = useParams();
    let [tutorial, setTutorial] = useState();
    let [log, setLog] = useState();
    let [pageConfig, setPageConfig] = useState({
        page: 1, total: 0, pageSize: 50
    });

    const columns = [
        {
          title: '标题',
          dataIndex: 'label',
          key: 'label',
          render: () => (
            <p>{tutorial.label}</p>
          )
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status) => (
                <div className="flex">
                    <div className={`point ${status === 2 ? "point-success" : "point-error"}`}></div>
                    <p>{ status === 2 ? "打包成功" : "打包失败" }</p>
                </div>
            )
        },
        {
            title: '分支',
            key: 'branch',
            dataIndex: 'branch',
            render: (branch) => (
              branch ? branch : "main"
            )
        },
        {
            title: '文档目录',
            key: 'startPage',
            dataIndex: 'startPage',
            render: () => (
              <p>{tutorial.startPage}</p>
            )
        },
        {
            title: '教程文档CommitHash',
            key: 'commitHash',
            dataIndex: 'commitHash',
            render: (commitHash) => (
                <p>{commitHash}</p>
            )
        },
        {
            title: '打包时间',
            key: 'UpdatedAt',
            dataIndex: 'UpdatedAt',
            render: (UpdatedAt) => (
                <p>{UpdatedAt.replace("T", " ").split(".")[0]}</p>
            )
        }
    ];

    function getLog(num) {
        if (num) {                
            pageConfig.page = num;
            setPageConfig({...pageConfig});
        }
        const { page, pageSize } = pageConfig;
        //  获取日志
        getPackLog({id: Number(id), page, pageSize})
        .then(res => {
            if (res.code === 0) {
                const list = res.data.list;
                log = list ? list : [];
                log.forEach(e => {
                    e.key = e.ID
                });
                setLog([...log]);
                pageConfig.total = res.data.total;
                setPageConfig({...pageConfig});
            }
        })
    }

    function init() {
        //  获取详情
        getTutorial({id: Number(id)})
        .then(res => {
            if (res.code === 0) {
                tutorial = res.data;
                setTutorial({...tutorial});
            }
        })
        getLog();
    }

    useEffect(() => {
        init();
    },[])
    
    return (
        tutorial && 
        <div className="tutorials">
            <Link to={`/dashboard/tutorials/build`}>
                <ArrowLeftOutlined />
            </Link>
            <div className="tabel-title">
                <h2>打包记录</h2>
            </div>
            <Table 
                columns={columns} 
                dataSource={log}
                pagination={{
                    current: pageConfig.page, 
                    total: pageConfig.total, 
                    pageSize: pageConfig.pageSize, 
                    onChange: (page) => getLog(page)
                }} 
            />
        </div>
    )
}