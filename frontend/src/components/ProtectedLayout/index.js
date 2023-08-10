import { Navigate, useLocation, useNavigate, useOutlet } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import { useEffect, useState } from "react";
import {
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    GithubOutlined
} from '@ant-design/icons';
import { Layout, Menu, Button, theme } from 'antd';
import "./index.scss";
import { menu } from "./menu";
const { Header, Sider, Content } = Layout;

export const ProtectedLayout = () => {

  const { user, auth, logout } = useAuth();
  const outlet = useOutlet();
  const location = useLocation();
  const navigateTo = useNavigate();

  const { token: { colorBgContainer }, } = theme.useToken();
  const [collapsed, setCollapsed] = useState(false);

  if (!user) {
    return <Navigate to="/" />;
  }
//   自动生成侧边栏
  // function menu() {
    // const arr = []
    // TODO: menu生成
    // 鉴权生成
    // auth.forEach((element, index) => {
    //     arr.push({
    //         key: element,
    //         icon: "",
    //         label: element
    //     })
    // });
    
  //   return arr
  // }

  function defaultSelectedKeys() {
    // 获取当前路由 menu
    const path = location.pathname;
    let key = "";
    let open = "";
    menu.forEach(e => {
      if (e.children) {
        const res = e.children.filter(ele => path.includes(ele.key))
        if (res.length !== 0) {
          key = res[0].key
          open = e.key;
        }
      }else{
        if (path.includes(e.key)) {
          key = e.key
        }
      }
    })
    return {
      key,
      open
    }
  }
  
  return (
    <Layout style={{height: "100vh"}} className="main">
        {/* 侧边栏 */}
      <Sider trigger={null} collapsible collapsed={collapsed} className="main-sidebar">
        <div className="demo-logo-vertical" >
            <GithubOutlined />
        </div>
        <Menu
          theme="dark"
          mode="inline"
          defaultSelectedKeys={[defaultSelectedKeys().key]}  //   初始值
          defaultOpenKeys={[defaultSelectedKeys().open]}  //  默认展开
          items={menu}
          onClick={(item) => {navigateTo(item.key)}}  //   跳转
        />
      </Sider>
      {/* 正文 */}
      <Layout className="main-content">
        {/* 正文导航栏 */}
        <Header
          style={{
            padding: 0,
            background: colorBgContainer,
          }}
          className="content-navbar"
        >
            {/* menu闭合按钮 */}
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{
              fontSize: '16px',
              width: 64,
              height: 64,
            }}
          />

          {/* 控制台 */}
          <div className="operate">
            {!!user && (
                <button key={"logout"} onClick={logout}>
                  <p >Logout</p>
                </button>
            )}
          </div>
        </Header>
        {/* 正文内容 */}
        <Content
          style={{
            margin: '24px 16px',
            padding: 24,
            minHeight: 280,
            background: colorBgContainer,
            overflow: "auto"
          }}
          className="content-info"
        >
          {outlet}
        </Content>
      </Layout>
    </Layout>
  );
};
