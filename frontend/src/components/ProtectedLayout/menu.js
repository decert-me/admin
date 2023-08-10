import {
    MenuOutlined,
    AuditOutlined,
    BookOutlined
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
                label: "教程审核",
                key: "tutorials/audit",
                icon: <AuditOutlined />,
            }
        ]
    }
]