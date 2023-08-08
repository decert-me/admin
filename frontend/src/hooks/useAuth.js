import { createContext, useContext, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import { useLocalStorage } from "./useLocalStorage";

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useLocalStorage("user", null);
  const [token, setToken] = useLocalStorage("token", null);
  const [auth, setAuth] = useLocalStorage("auth", null);

  const navigate = useNavigate();

  const login = async (token, data, auth) => {
    setUser(data);
    setToken(token);
    // setAuth(auth)
    navigate("/dashboard/profile", { replace: true });
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    // setAuth(null);
    navigate("/", { replace: true });
  };

  const value = useMemo(
    () => ({
      user,
      auth,
      token,
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
