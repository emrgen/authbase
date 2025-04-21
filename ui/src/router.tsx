import {LoginForm} from "@/old-pages/Login.tsx";
import {Accounts} from "@/pages/Accounts.tsx";
import {Clients} from "@/pages/Clients.tsx";
import {Dashboard} from "@/pages/Dashboard.tsx";
import {Providers} from "@/pages/Providers.tsx";
import {createBrowserRouter} from "react-router";
import App from "./App.tsx";


export const router = createBrowserRouter([
  {
    path: "/",
    element: <App/>,
    children: [
      {
        index: true,
        element: <Dashboard/>,
      },
      // {
      //   path: '/session',
      //   element: <Session/>,
      // },
      {
        path: '/account',
        element: <Accounts/>,
      },
      {
        path: '/client',
        element: <Clients/>,
      },
      {
        path: '/provider',
        element: <Providers/>,
      },
      // {
      //   path: '/access-key',
      //   element: <AccessKey/>,
      // },
    ]
  },
  {
    path: '/login',
    element: <LoginForm/>,
  },
]);
