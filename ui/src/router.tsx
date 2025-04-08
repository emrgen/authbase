import {createBrowserRouter} from "react-router";
import App from "./App.tsx";
import {AccessKey} from "./pages/AccessKey.tsx";
import {Accounts} from "./pages/Accounts.tsx";
import {Clients} from "./pages/Clients.tsx";
import {Dashboard} from "./pages/Dashboard.tsx";
import {LoginPage} from "./pages/Login.tsx";
import {Pool} from "./pages/Pool.tsx";
import {Provider} from "./pages/Provider.tsx";
import {Session} from "./pages/Session.tsx";


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
        path: '/pool',
        element: <Pool/>,
      },
      {
        path: '/session',
        element: <Session/>,
      },
      {
        path: '/account',
        element: <Accounts/>,
      },
      {
        path: '/clients',
        element: <Clients/>,
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
