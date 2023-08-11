import {
    MenuOutlined,
    BookOutlined,
    FolderOutlined,
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
                label: "打包管理",
                key: "tutorials/build",
                icon: <FolderOutlined />,
            },
        ]
    },
]