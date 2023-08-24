import { Button, Popconfirm, Space, Switch, Table, message } from "antd";
import { useEffect, useState } from "react";
import "./index.scss";
import { deleteQuest, getQuestList, topQuest, updateQuestStatus } from "../../request/api/quest";
import { format } from "../../utils/format";
import { useNavigate } from "react-router-dom";

const isTest = window.location.host.indexOf("localhost") === -1;
const host = isTest ? "https://decert.me" : "http://192.168.1.10:8087";
const opensea = isTest ? "https://opensea.io/assets/matic/0xc8E9cd4921E54c4163870092Ca8d9660e967B53d" : "https://testnets.opensea.io/assets/mumbai/0x66C54CB10Ef3d038aaBA2Ac06d2c25B326be8142"

export default function ChallengeListPage(params) {

    const { formatTimestamp } = format();
    const navigateTo = useNavigate();

    const [selectedRowKeys, setSelectedRowKeys] = useState([]);     //  多选框: 选中的挑战
    const [topLoad, setTopLoad] = useState(false);    //  置顶等待
    let [data, setData] = useState([]);
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 10, total: 0
    });

    const columns = [
        {
          title: '挑战编号',
          dataIndex: 'tokenId',
          render: (tokenId) => (
            <a className="underline" href={`${host}/quests/${tokenId}`} target="_blank">{tokenId}</a>
          )
        },
        {
          title: 'NFT',
          dataIndex: 'metadata',
          render: ({image}, quest) => (
            quest.claim_num !== 0 ?
            <a href={`${opensea}/${quest.tokenId}`} target="_blank">
                <img src={image.replace("ipfs://", "https://ipfs.decert.me/")} alt="" style={{height: "53px"}} />
            </a>
            :
            <img src={image.replace("ipfs://", "https://ipfs.decert.me/")} alt="" style={{height: "53px"}} />
          )
        },
        {
          title: '标题',
          dataIndex: 'title',
          render: (title, quest) => (
            <a className="underline" href={`${host}/quests/${quest.tokenId}`} target="_blank">{title}</a>
          )
        },
        {
            title: '发布者',
            dataIndex: 'creator',
            render: (creator) => (
                <a className="underline" href={`${host}/user/${creator}`} target="_blank">{creator.substring(0,5) + "..." + creator.substring(38,42)}</a>
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
            dataIndex: 'metadata',
            render: ({attributes}) => (
              <p>{attributes.difficulty === 0 ? "简单" : attributes.difficulty === 1 ? "一般" : attributes.difficulty === 2 ? "困难" : "/"}</p>
            )
        },
        {
            title: '时长',
            dataIndex: 'quest_data',
            render: ({estimateTime}) => (
              <p>{estimateTime ? (estimateTime / 60) + "min" : "/"}</p>
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
            title: '挑战人次',
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
            title: 'Action',
            key: 'action',
            render: (_, quest) => (
              <Space size="middle">
                <Button 
                  type="link" 
                  className="p0"
                  onClick={() => navigateTo(`/dashboard/challenge/modify/${quest.id}/${quest.tokenId}`)}
                >编辑</Button>
                <Popconfirm
                  title="移除挑战"
                  description="确定要移除该挑战吗?"
                  onConfirm={() => deleteT(Number(quest.id))}
                  okText="确定"
                  cancelText="取消"
                >
                  <Button 
                  type="link" 
                  className="p0"
                >删除</Button>
                </Popconfirm>
              </Space>
            ),
        }
    ];

    // 上下架
    function handleChangeStatus({status, id}, key) {
      const index = data.findIndex((item) => item.key === key);
      updateQuestStatus({id, status})
      .then(res => {
        if (res.code === 0) {
          message.success(res.msg);
          data[index].status = status;
          setData([...data]);
        }
      })
    }

    // 删除教程
    async function deleteT(id) {
        await deleteQuest({id})
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
    
    const onSelectChange = (newSelectedRowKeys) => {
        setSelectedRowKeys(newSelectedRowKeys);
    };
    const rowSelection = {
        selectedRowKeys,
        onChange: onSelectChange,
    };
    const hasSelected = selectedRowKeys.length > 0;


    // 挑战置顶
    function toTop(status) {
        setTopLoad(true);
        const statusArr = Array(selectedRowKeys.length).fill(status);
        topQuest({id: selectedRowKeys, top: statusArr})
        .then(res => {
          setTopLoad(false);
          if (res.code === 0) {
            message.success(res.msg);
            setSelectedRowKeys([...[]]);
            getList()
          }
        })
    }

    function getList(page) {
        if (page) {
          pageConfig.page = page;
          setPageConfig({...pageConfig});
        }
        // 获取教程列表
        getQuestList(pageConfig)
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
        .catch(err => {
            message.error(err)
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
                <h2>挑战列表</h2>
                <Space size="large">
                    <Button 
                        onClick={() => toTop(true)} 
                        disabled={!hasSelected}
                        loading={topLoad}
                    >
                        置顶
                    </Button>
                    <Button 
                        onClick={() => toTop(false)} 
                        disabled={!hasSelected}
                        loading={topLoad}
                    >
                        取消置顶
                    </Button>
                </Space>
            </div>
            <Table 
                rowSelection={rowSelection} 
                columns={columns} 
                dataSource={data} 
                rowClassName={(record) => record.top && "toTop"}
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