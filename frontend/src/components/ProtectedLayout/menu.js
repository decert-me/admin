import {
    MenuOutlined,
    AuditOutlined,
    BookOutlined,
    FileAddOutlined
  } from '@ant-design/icons';

export const menu = [
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
    }
]