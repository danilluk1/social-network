import {
  Routes,
  Route,
  Navigate,
  Outlet,
  BrowserRouter,
} from "react-router-dom";
import Main from "../../pages/Main";
import Login from "../../pages/Login";
import Confirm from "../../pages/Confirm";
import Reset from "../../pages/Reset";
import Register from "../../pages/Register";

interface Props {
  isAuth: boolean;
  redirectPath: string;
}

const ProtectedRoute = ({ isAuth, redirectPath }: Props) => {
  if (!isAuth) {
    return <Navigate to={redirectPath} replace />;
  }

  return <Outlet />;
};

const Router = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<ProtectedRoute isAuth={false} redirectPath="login" />}>
          <Route index element={<Main />} />
        </Route>
        <Route path="login" element={<Login />} />
        <Route path="register" element={<Register />} />
        <Route path="confirm" element={<Confirm />} />
        <Route path="reset" element={<Reset />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;
