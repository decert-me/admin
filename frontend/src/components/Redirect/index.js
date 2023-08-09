import { Navigate } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";

export const Redirect = () => {
  const { user } = useAuth();

  if (user) {
    return <Navigate to="/dashboard/tutorials/list" replace />;
  } else {
    return <Navigate to="/login" replace />;
  }
};
