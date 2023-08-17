import { Space, Table, message } from "antd";
import { useEffect, useState } from "react";
import { mock } from "../../mock";
import { getPackList } from "../../request/api/tutorial";
import { Link } from "react-router-dom";



export default function TutorialsBuildPage(params) {
    
    const { tutorials } = mock();
    let [data, setData] = useState([]);
    let [pageConfig, setPageConfig] = useState({
      page: 0, pageSize: 10, total: 0
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
                <p>{commitHash}</p>
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
              <a>打包</a>
              <Link to={`/dashboard/tutorials/buildlog/${tutorial.ID}`}>
                打包记录
              </Link>
            </Space>
          ),
        },
    ];

    function init() {
      pageConfig.page += 1;
      setPageConfig({...pageConfig});
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
            <Table columns={columns} dataSource={data} />
        </div>
    )
}