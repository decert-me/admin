import { Button, Popconfirm, Space, Switch, Table, message } from "antd";
import { useNavigate } from "react-router-dom";
import "./index.scss";
import { useEffect, useState } from "react";
import { deleteCollection, getCollectionList, updateStatusCollection } from "../../request/api/quest";
import { format } from "../../utils/format";

const location = window.location.host;
const host = ((location.indexOf("localhost") !== -1) || (location.indexOf("192.168.1.10") !== -1)) ? "http://192.168.1.10:8087" : "https://decert.me";

export default function ChallengeCompilationPage(params) {
    
    const { formatTimestamp } = format();
    const navigateTo = useNavigate();
    let [data, setData] = useState();
    let [pageConfig, setPageConfig] = useState({
      page: 0, pageSize: 50, total: 0
    });

    const handleChange = (pagination, filters, sorter) => {
      const { pageSize } = pagination
      if (pageSize !== pageConfig.pageSize) {
          pageConfig.pageSize = pageSize;
          setPageConfig({...pageConfig});
          getList();
        }
    };

    const columns = [
        {
          title: '权重',
          dataIndex: 'sort'
        },
        {
          title: 'ID',
          dataIndex: 'id',
          render: (tokenId) => (
            <a className="underline" href={`${host}/collection/${tokenId}`} target="_blank">{tokenId}</a>
          )
        },
        {
          title: '合辑图片',
          dataIndex: 'cover',
          render: (image) => (
            <img src={"https://ipfs.decert.me/"+image} alt="" style={{height: "53px"}} />
          )
        },
        {
          title: '标题',
          dataIndex: 'title',
          render: (title, quest) => (
            <a className="underline" href={`${host}/collection/${quest.id}`} target="_blank">{title}</a>
          )
        },
        {
            title: '作者',
            dataIndex: 'author',
            render: (author) => (
              <p>{author}</p>
            )
        },
        {
            title: '上架状态',
            dataIndex: 'status',
            render: (status, quest) => (
                <Switch
                    checkedChildren="已上架" 
                    unCheckedChildren="待上架" 
                    checked={status == 1 ? true : false}
                    onChange={(checked) => handleChangeStatus({status: checked ? 1 : 2, id: quest.id}, quest.key)}
                />
            )
        },
        {
          title: '难度',
          dataIndex: 'difficulty',
          render: (difficulty) => (
            <p>{difficulty === 0 ? "简单" : difficulty === 1 ? "中等" : difficulty === 2 ? "困难" : "/"}</p>
          )
        },
        {
          title: '时长',
          dataIndex: 'estimate_time',
          render: (estimate_time) => (
            <p>{estimate_time / 60}min</p>
          )
        },
        {
          title: '挑战数量',
          dataIndex: 'quest_num',
          render: (quest_num) => (
            <p>{quest_num}</p>
          )
        },
        {
          title: '铸造数量',
          dataIndex: 'claim_num',
          render: (claim_num) => (
            <p>{claim_num}</p>
          )
        },
        {
            title: '挑战人数',
            dataIndex: 'challenge_num',
            render: (challenge_num) => (
              <p>{challenge_num}次</p>
            )
        },
        {
          title: '创建时间',
          dataIndex: 'addTs',
          render: (addTs) => (
            <p>{formatTimestamp(addTs * 1000)}</p>
          )
      },
        {
            title: '操作',
            key: 'action',
            render: (_, quest) => (
              <Space size="middle">
                <Button 
                  type="link" 
                  className="p0"
                  onClick={() => navigateTo(`/dashboard/challenge/compilation/modify/${quest.id}`)}
                >编辑</Button>
                <Button 
                  type="link" 
                  className="p0"
                  onClick={() => navigateTo(`/dashboard/challenge/list/${quest.id}?label=${quest.title}`)}
                >目录</Button>
                <Popconfirm
                  title="移除合辑"
                  description="确定要移除该合辑吗?"
                  onConfirm={() => deleteT(Number(quest.id))}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button 
                  type="link" 
                  className="p0"
                >移除</Button>
                </Popconfirm>
              </Space>
            ),
        }
    ];

    // 上下架
    function handleChangeStatus({id, status}, key){
      const index = data.findIndex((item) => item.key === key);
      updateStatusCollection({id, status})
      .then(res => {
        if (res.code === 0) {
          message.success(res.msg);
          data[index].status = status;
          setData([...data]);
        }
      })
    };

    async function deleteT(id) {
      await deleteCollection({id})
      .then(res => {
        if (res.code === 0) {
            message.success(res.msg);
        }
      })
      .catch(err => {
          message.error(err);
      })
      getList()
    }

    function getList(page) {
      if (page) {
        pageConfig.page = page;
        setPageConfig({...pageConfig});
      }
      // 获取挑战合辑列表
      getCollectionList(pageConfig)
      .then(res => {
        if (res.code === 0) {
          const list = res.data.list;
          data = list ? list : [];
          // 添加key
          data.forEach(ele => {
            ele.key = ele.id
          })
          setData([...data]);
          pageConfig.total = res.data.total;
          setPageConfig({...pageConfig});
        }else{
          message.success(res.msg);
        }
      })
    }

    function init(params) {
      pageConfig.page += 1;
      setPageConfig({...pageConfig});
      getList()
    }

    useEffect(() => {
      init();
    },[])

    return (
        <div className="challenge">
            <div className="tabel-title">
                <h2>挑战合辑</h2>
                <Button
                    type="primary"
                    onClick={() => navigateTo("/dashboard/challenge/add")} 
                >
                    创建合辑
                </Button>
            </div>
            <Table
                columns={columns} 
                dataSource={data} 
                onChange={handleChange}
                // rowClassName={(record) => record.top && "toTop"}
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