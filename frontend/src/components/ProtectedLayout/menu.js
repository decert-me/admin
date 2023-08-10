import {
    MenuOutlined,
    AuditOutlined,
    BookOutlined,
    FileAddOutlined,
    HomeOutlined
  } from '@ant-design/icons';

export const menu = [
    {
        label: "首页",
        key: "home",
        icon: <HomeOutlined />,
    },
    {
        label: "教程管理",
        key: "tutorials",
        icon: <BookOutlined />,
        children: [
            {
                label: "教程列表",
                key: "tutorials/list",
                icon: <MenuOutlined />,
            },
            {
                label: "添加教程",
                key: "tutorials/add",
                icon: <FileAddOutlined />,
            }
        ]
    },
]