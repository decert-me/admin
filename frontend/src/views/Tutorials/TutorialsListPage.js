import { Button, Popconfirm, Space, Switch, Table, Tag, message } from "antd";
import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
    VideoCameraOutlined,
    ReadOutlined
  } from '@ant-design/icons';
import { deleteTutorial, getTutorialList } from "../../request/api/tutorial";
import { getLabelList } from "../../request/api/tags";

export default function TutorialsListPage(params) {
    
    const navigateTo = useNavigate();
    let [tags, setTags] = useState([]);
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
          title: '上架状态',
          key: 'status',
          dataIndex: 'status',
          render: (status) => (
              <Switch checkedChildren="已上架" unCheckedChildren="待上架" defaultChecked={status == 0 ? true : false} />
          )
        },
        {
          title: '分类',
          dataIndex: 'category',
          key: 'category',
          render: (category) => (
              category.map(tag => 
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
              <Tag color={language === "zh" ? "#2db7f5" : "#87d068"}>
                  {language === "zh" ? "中文" : "英文"}
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
          title: 'Action',
          key: 'action',
          render: (_, tutorial) => (
            <Space size="middle">
              <Link to={`/dashboard/tutorials/modify/${tutorial.ID}`}>修改</Link>
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

    function getList() {
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
      getList()
    }

    useEffect(() => {
        init();
    },[])

    return (
        <div className="tutorials tutorials-list">
          <div className="tabel-title">
            <h2>教程列表</h2>
            <Button 
              type="primary"
              onClick={() => navigateTo("/dashboard/tutorials/add")}
            >创建教程</Button>
          </div>
            <Table columns={columns} dataSource={data} pagination={{current: pageConfig.page, total: pageConfig.total, pageSize: pageConfig.pageSize}} />
        </div>
    )
}