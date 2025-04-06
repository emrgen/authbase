import {createBrowserRouter} from "react-router";
import App from "./App.tsx";
import {LoginPage} from "./pages/Login.tsx";


export const router = createBrowserRouter([
  {
    path: "/",
    element: <App/>,
  },
  {
    path: '/login',
    element: <LoginPage/>,
  },
]);
