import { Button, Space, Table, message } from "antd";
import { useEffect, useState } from "react";
import { buildTutorial, getPackList } from "../../request/api/tutorial";
import { Link } from "react-router-dom";



export default function TutorialsBuildPage(params) {
    
    let [loading, setLoading] = useState(false);
    let [data, setData] = useState([]);
    let [pageConfig, setPageConfig] = useState({
      page: 1, total: 0, pageSize: 10
    });

    const columns = [
        {
          title: '封面图',
          dataIndex: 'img',
          key: 'img',
          render: (img) => (
            <img src={`https://ipfs.decert.me/${img}`} alt="" style={{height: "40px"}} />
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
            dataIndex: 'pack_status',
            key: 'pack_status',
            render: (pack_status) => (
                <div className="build-status">
                    <div className={`point ${pack_status === 2 ? "success" : "error"}`}></div>
                    <p>{pack_status === 2 ? "打包成功" : "打包失败"}</p>
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
                <p>{commitHash ? commitHash.substring(0,7) : ""}</p>
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
            title: '最新打包时间',
            dataIndex: 'UpdatedAt',
            key: 'UpdatedAt',
            render: (UpdatedAt) => (
                <p>{UpdatedAt.replace("T", " ").split(".")[0]}</p>
            )
        },
        {
          title: 'Action',
          key: 'action',
          render: (_, tutorial) => (
            <Space size="middle">
              <Button type="primary" onClick={() => build(tutorial.id)} loading={loading}>
                打包
              </Button>
              <Link to={`/dashboard/tutorials/buildlog/${tutorial.id}`}>
                <Button type="primary">
                  打包记录
                </Button>
              </Link>
            </Space>
          ),
        },
    ];

    async function build(id) {
      setLoading(true);
      await buildTutorial({id})
      .then(res => {
        if (res.code === 0) {
          
        }
      })
      setLoading(false);
    }

    function init(page) {
      if (page) {        
        pageConfig.page = page;
        setPageConfig({...pageConfig});
      }
      getPackList(pageConfig)
      .then(res => {
        if (res.code === 0) {
          const list = res.data.list;
          data = list ? list : [];
          // 添加key
          data.forEach((ele, index) => {
            ele.key = index
          })
          setData([...data]);
          pageConfig.total = res.data.total;
          setPageConfig({...pageConfig});
        }else{
            message.success(res.msg);
        }
      })
      .catch(err => {
          message.error(err)
      })
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="tutorials tutorials-build">
            <div className="tabel-title">
                <h2>打包管理</h2>
            </div>
            <Table 
              columns={columns}
              dataSource={data} 
              pagination={{
                current: pageConfig.page, 
                total: pageConfig.total, 
                pageSize: pageConfig.pageSize, 
                onChange: (page) => init(page)
              }} 
            />
        </div>
    )
}