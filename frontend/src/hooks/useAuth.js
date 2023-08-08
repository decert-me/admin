import { createContext, useContext, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { useLocalStorage } from "./useLocalStorage";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useLocalStorage("user", null);
  const [auth, setAuth] = useLocalStorage("auth", null);

  const navigate = useNavigate();

  const login = async (data, auth) => {
    setUser(data);
    setAuth(auth)
    navigate("/dashboard/profile", { replace: true });
  };

  const logout = () => {
    setUser(null);
    setAuth(null);
    navigate("/", { replace: true });
  };

  const value = useMemo(
    () => ({
      user,
      auth,
      login,
      logout
    }),
    [user]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuth = () => {
  return useContext(AuthContext);
};
