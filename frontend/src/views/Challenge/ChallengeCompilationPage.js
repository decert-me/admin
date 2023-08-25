import { Button, Popconfirm, Space, Switch, Table } from "antd";
import { useNavigate } from "react-router-dom";
import "./index.scss";

const isTest = window.location.host.indexOf("localhost") === -1;
const host = isTest ? "https://decert.me" : "http://192.168.1.10:8087";

export default function ChallengeCompilationPage(params) {
    
    const navigateTo = useNavigate();

    const columns = [
        {
          title: 'ID',
          dataIndex: 'id',
          render: (tokenId) => (
            <a className="underline" href={`${host}/quests/${tokenId}`} target="_blank">{tokenId}</a>
          )
        },
        {
          title: '合辑图片',
          dataIndex: 'img',
          render: ({image}) => (
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
            title: '作者',
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
                    // onChange={(checked) => handleChangeStatus({status: checked ? 1 : 2, id: quest.id}, quest.key)}
                />
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
            //   <p>{formatTimestamp(addTs * 1000)}</p>
                {addTs}
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
                <Button 
                  type="link" 
                  className="p0"
                  onClick={() => navigateTo(`/dashboard/challenge/modify/${quest.id}/${quest.tokenId}`)}
                >排序</Button>
                <Popconfirm
                  title="移除合辑"
                  description="确定要移除该合辑吗?"
                //   onConfirm={() => deleteT(Number(quest.id))}
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
                // dataSource={data} 
                // rowClassName={(record) => record.top && "toTop"}
                // pagination={{
                //     current: pageConfig.page, 
                //     total: pageConfig.total, 
                //     pageSize: pageConfig.pageSize, 
                //     onChange: (page) => getList(page)
                // }} 
            />
        </div>
    )
}