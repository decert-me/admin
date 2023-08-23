import { Button, Modal, Popconfirm, Space, Switch, Table, Tag, message } from "antd";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
    VideoCameraOutlined,
    ReadOutlined
  } from '@ant-design/icons';
import { buildTutorial, deleteTutorial, getTutorialList, topTutorial, updateTutorialStatus } from "../../request/api/tutorial";
import { getLabelList } from "../../request/api/tags";

export default function TutorialsListPage(params) {
    
    const navigateTo = useNavigate();
    const [selectedRowKeys, setSelectedRowKeys] = useState([]);   //  列表所选item
    const [topLoad, setTopLoad] = useState(false);    //  置顶等待
    let [loading, setLoading] = useState(false);    //  打包loading
    let [tags, setTags] = useState([]);
    let [lang, setLang] = useState([]);
    let [data, setData] = useState([]);
    let [pageConfig, setPageConfig] = useState({
      page: 0, pageSize: 10, total: 0
    });
    let [log, setLog] = useState();
    const [selectKey, setSelectKey] = useState('');
    const isSelect = (record) => record.key === selectKey;
    const [isModalOpen, setIsModalOpen] = useState(false);


    // 教程上下架
    const handleChangeStatus = ({id, checked}, key) => {
      const index = data.findIndex((item) => item.key === key);
      updateTutorialStatus({id, status: checked ? 2 : 1})
      .then(res => {
        if (res.code === 0) {
          message.success(res.msg);
          data[index].status = checked ? 2 : 1;
          setData([...data]);
        }
      })
    };

    // 打包
    async function build(id, key) {
      setSelectKey(key);
      setLoading(true);
      await buildTutorial({id})
      .then(res => {
        if (res.code === 0) {
          message.success(res.msg);
        }
      })
      setLoading(false);
      getList()
    }

    function goBuild(tutorial) {
      const selet = isSelect(tutorial);
      return (
        <Button 
          type="link" 
          className="p0" 
          loading={selet && loading} 
          disabled={selectKey !== tutorial.key && loading}
          onClick={() => build(tutorial.ID, tutorial.key)}
        >
          打包
        </Button>
      )
    }

    // 日志
    const showModal = (logs) => {
      log = logs;
      setLog(log);
      setIsModalOpen(true);
    };

    const handleCancel = () => {
      setIsModalOpen(false);
    };

    // 教程置顶
    function toTop(status) {
      setTopLoad(true);
      const statusArr = Array(selectedRowKeys.length).fill(status);
      topTutorial({id: selectedRowKeys, top: statusArr})
      .then(res => {
        setTopLoad(false);
        if (res.code === 0) {
          message.success(res.msg);
          setSelectedRowKeys([...[]]);
          getList()
        }
      })
    }

    const hasSelected = selectedRowKeys.length > 0;

    const onSelectChange = (newSelectedRowKeys) => {
      setSelectedRowKeys(newSelectedRowKeys);
    };

    const rowSelection = {
      selectedRowKeys,
      onChange: onSelectChange,
    };

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
          render: (text, tutorial) => (
            tutorial.status == 2 ?
            <a className="tabel-item-title newline-omitted underline" href={`${window.location.host.indexOf("localhost") === -1 ? "https://decert.me" : "http://192.168.1.10:8087"}/tutorial/${tutorial.startPage}/`} target="_blank">{tutorial.top ? "【置顶】" : ""}{text}</a>
            :
            <p className="tabel-item-title newline-omitted">{tutorial.top ? "【置顶】" : ""}{text}</p>
          )
        },
        {
          title: '上架状态',
          key: 'status',
          dataIndex: 'status',
          render: (status, tutorial) => (
              <Switch 
                checkedChildren="已上架" 
                unCheckedChildren="待上架" 
                checked={status == 2 ? true : false}
                onChange={(checked) => handleChangeStatus({checked: checked, id: tutorial.ID}, tutorial.key)}
              />
          )
        },
        {
          title: '分类',
          dataIndex: 'category',
          key: 'category',
          render: (category) => (
              category && category.map(tag => 
                  <Tag color="geekblue" key={tag}>
                    {
                      tags.filter(e => e.ID === tag).length !== 0 &&
                        tags.filter(e => e.ID === tag)[0].Chinese
                    }
                  </Tag>    
              )
          )
        },
        {
          title: '语言',
          key: 'language',
          dataIndex: 'language',
          render: (language) => (
              <Tag>
                  {lang.filter(e => e.ID === language)[0].Chinese}
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
            title: '创建时间',
            key: 'CreatedAt',
            dataIndex: 'CreatedAt',
            render: (CreatedAt) => (
              <p>{CreatedAt.replace("T", " ").split(".")[0]}</p>
            )
        },
        {
          title: 'Action',
          key: 'action',
          render: (_, tutorial) => (
            <Space size="middle">
              <Link to={`/dashboard/tutorials/modify/${tutorial.ID}`}>编辑</Link>
              {goBuild(tutorial)}
              <Button 
                type="link" 
                className="p0"
                onClick={() => showModal(tutorial.pack_log)}
              >日志</Button>
              <Popconfirm
                title="删除教程"
                description="确定要删除这篇教程吗?"
                onConfirm={() => deleteT(tutorial.ID)}
                okText="确定"
                cancelText="取消"
              >
                <a>删除</a>
              </Popconfirm>
            </Space>
          ),
        },
    ];

    async function deleteT(id) {
      await deleteTutorial({id})
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
      // 获取教程列表
      getTutorialList(pageConfig)
      .then(res => {
        if (res.code === 0) {
          const list = res.data.list;
          data = list ? list : [];
          // 添加key
          data.forEach(ele => {
            ele.key = ele.ID
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

    async function init() {
      pageConfig.page += 1;
      setPageConfig({...pageConfig});
      // 获取标签列表
      getLabelList({type: "category"})
      .then(res => {
        if (res.code === 0) {
          tags = res.data ? res.data : [];
          setTags([...tags]);
        }
      })
      getLabelList({type: "language"})
      .then(res => {
        if (res.code === 0) {
          lang = res.data ? res.data : [];
          setLang([...lang]);
        }
      })
      getList()
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="tutorials tutorials-list">
          <div className="tabel-title">
            <h2>教程列表</h2>
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
              <Button 
                type="primary"
                onClick={() => navigateTo("/dashboard/tutorials/add")}
              >创建教程</Button>
            </Space>
          </div>

            <Table 
              rowSelection={rowSelection}
              columns={columns} 
              dataSource={data} 
              pagination={{
                current: pageConfig.page, 
                total: pageConfig.total, 
                pageSize: pageConfig.pageSize, 
                onChange: (page) => getList(page)
              }} 
            />
            <Modal width={800} open={isModalOpen} onCancel={handleCancel} footer={null}>
              <p dangerouslySetInnerHTML={{__html: log}}></p>
            </Modal>
        </div>
    )
}