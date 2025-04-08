import {createBrowserRouter} from "react-router";
import App from "./App.tsx";
import {AccessKey} from "./pages/AccessKey.tsx";
import {Accounts} from "./pages/Accounts.tsx";
import {Dashboard} from "./pages/Dashboard.tsx";
import {LoginPage} from "./pages/Login.tsx";
import {Provider} from "./pages/Provider.tsx";


export const router = createBrowserRouter([
  {
    path: "/",
    element: <App/>,
    children: [
      {
        index: true,
        element: <Dashboard/>,
      },
      {
        path: '/account',
        element: <Accounts/>,
      },
      {
        path: '/provider',
        element: <Provider/>,
      },
      {
        path: '/access-key',
        element: <AccessKey/>,
      },
    ]
  },
  {
    path: '/login',
    element: <LoginPage/>,
  },
]);
