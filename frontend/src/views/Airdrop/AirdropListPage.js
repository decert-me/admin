import { Button, Table, message } from "antd";
import { useEffect, useState } from "react";
import { getAirdropList, runAirdrop } from "../../request/api/airdrop";
import "./index.scss";
import { CHAINS } from "../../config/CHAINS";

export default function AirdropList(params) {
    
    const isDev = process.env.REACT_APP_IS_DEV;
    let [data, setData] = useState([]);
    let [status, setStatus] = useState(3);
    let [pageConfig, setPageConfig] = useState({
        page: 0, pageSize: 10, total: 0
    });

    const handleChange = (pagination, filters, sorter) => {
        const { pageSize } = pagination
        const newStatus = Array.isArray(filters.status) ? filters.status[0] : null;
        if (status !== newStatus) {
            status = newStatus;
            setStatus(newStatus);
            getList(1);
        }
        if (pageSize !== pageConfig.pageSize) {
            pageConfig.pageSize = pageSize;
            setPageConfig({...pageConfig});
            getList();
        }
    };

    const addressHref = (addr, {app}) => {
        if (app === "decert") {
            const prefix = isDev ? "https://mumbai.polygonscan.com" : "https://polygonscan.com"
            return prefix + "/address/" + addr
        }else if (app === "decert_solana") {
            const suffix = isDev ? "?cluster=devnet" : "";
            return `https://solscan.io/account/${addr}${suffix}`
        }
    }

    const txhashHref = (hash, {app, params}) => {
        if (app === "decert") {
            const prefix = isDev ? "https://mumbai.polygonscan.com" : "https://polygonscan.com"
            return prefix + "/tx/" + hash
        }
        if (app === "decert_solana") {
            const suffix = isDev ? "?cluster=devnet" : "";
            return `https://solscan.io/tx/${hash}${suffix}`
        }
        if (app === "decert_v2") {
            const chain = CHAINS.filter(e => e.chainID == params.params.chain_id);
            return `${chain[0].url}${hash}`
        }
    }

    const chain = (app, info) => {
        if (app === "decert_solana") {
            return "Solana"
        }
        if (app === "decert") {
            return "Polygon"
        }
        if (app === "decert_v2") {
            const chain = CHAINS.filter(e => e.chainID == info.params.params.chain_id);
            return chain[0]?.name
        }
    }

    const columns = [
        {
            title: 'ID',
            dataIndex: 'ID',
            key: 'ID',
            render: (ID) => (
              <p>{ID}</p>
            )
        },
        {
            title: '链',
            dataIndex: 'app',
            key: 'app',
            render: (app, info) => (
              <p>{chain(app, info)}</p>
            )
        },
        {
            title: '地址',
            dataIndex: 'params',
            key: 'params',
            render: ({params}, record) => (
                <a 
                    href={addressHref(params.receiver, record)} 
                    target="_blank" >
                    {params.receiver.substring(0,5) + "..." + params.receiver.substring(38,42)}
                </a>
            )
        },
        {
            title: '交易哈希',
            dataIndex: 'airdrop_hash',
            key: 'airdrop_hash',
            render: (airdrop_hash, record) => (
                airdrop_hash &&
                <a 
                    href={
                        txhashHref(airdrop_hash, record)
                        // process.env.REACT_APP_IS_DEV ? 
                        // `https://mumbai.polygonscan.com/tx/${airdrop_hash}`
                        // :
                        // `https://polygonscan.com/tx/${airdrop_hash}`
                    } 
                    target="_blank" >
                    {airdrop_hash.substring(0,5) + "..." + airdrop_hash.substring(61,66)}
                </a>
            )
        },
        {
            title: '信息',
            dataIndex: 'params',
            key: 'params',
            ellipsis: true,
            render: (params) => (
                JSON.stringify(params)
            )
        },
        {
            title: '创建时间',
            dataIndex: 'CreatedAt',
            key: 'CreatedAt',
            render: (CreatedAt) => (
                CreatedAt.replace("T", " ").split(".")[0]
            )
        },
        {
            title: '空投时间',
            dataIndex: 'airdrop_time',
            key: 'airdrop_time',
            render: (airdrop_time) => (
                airdrop_time.indexOf("0001-01-01T") === -1 ?
                <p>{airdrop_time.replace("T", " ").split(".")[0]}</p>
                :
                "-"
            )
        },
        {
            title: '错误信息',
            dataIndex: 'msg',
            key: 'msg'
        },
        {
            title: `状态:${status === 1 ? "待空投" : status === 2 ? "已空投" : status === 3 ? "空投失败" : "全部"}`,
            dataIndex: 'status',
            key: 'status',
            filters: [
                { text: '待空投', value: 1 },
                { text: '已空投', value: 2 },
                { text: '空投失败', value: 3 },
            ],
            filterMultiple: false,
            filteredValue: [status],
            render: (status) => (
                <p style={{
                    color: status === 2 ? "#09CD92" : status === 3 ? "#FF0000" : "000"
                }}>{status === 1 ? "待空投" : status === 2 ? "已空投" : "空投失败"}</p>
            )
        }
    ];

    function goAirdrop() {
        runAirdrop({
            app: "decert_all"
        })
        .then(res => {
            if (res.code === 0) {
                message.success("操作成功!")
            }
        })
    }

    function getList(page) {
        if (page) {
            pageConfig.page = page;
            setPageConfig({...pageConfig});
        }
        getAirdropList({
            status,
            app: "decert_all",
            ...pageConfig
        })
        .then(res => {
            const list = res.data.list;
            data = list ? list : [];
            // 添加key
            data.forEach(ele => {
                ele.key = ele.ID
            })
            setData([...data]);
            pageConfig.total = res.data.total;
            setPageConfig({...pageConfig});
        })
    }

    function init() {
        pageConfig.page += 1;
        setPageConfig({...pageConfig});
        getList();
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="airdrop">
            <div className="airdrop-btn">
                <Button
                    type="primary"
                    onClick={goAirdrop}
                >立即空投</Button>  
            </div>
            
            <Table 
                columns={columns} 
                dataSource={data} 
                onChange={handleChange}
                pagination={{
                    current: pageConfig.page, 
                    total: pageConfig.total, 
                    pageSize: pageConfig.pageSize, 
                    onChange: (page) => {
                        page !== pageConfig.page && getList(page)
                    }
                }} 
            />
        </div>
    )
}