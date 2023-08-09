import React from "react";
import { Navigate, useLocation } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";

const AuthGuard = ({ permissions, children }) => {

  const location = useLocation();
  const { auth } = useAuth();

  // function checkPermissions(permissions) {
  //   // 获取当前路由
  //   const path = location.pathname.split("/").pop();
  //   const isTrue = auth.some(e => e === path);
  //   return isTrue
  // }

  // // 根据权限判断是否有访问权限
  // const hasAccess = checkPermissions(permissions);

  // if (!hasAccess) {
  //   // 没有访问权限，可以重定向到其他页面或显示未授权信息
  //   return <Navigate to="/dashboard/profile" />;
  // }

  // 有访问权限，渲染内容
  return children;
};

export default AuthGuard;
