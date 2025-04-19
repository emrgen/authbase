import {createBrowserRouter} from "react-router";
import App from "./App.tsx";


export const router = createBrowserRouter([
  {
    path: "/",
    element: <App/>,
    children: [
      // {
      //   index: true,
      //   element: <Dashboard/>,
      // },
      // {
      //   path: '/pool',
      //   element: <Pool/>,
      // },
      // {
      //   path: '/session',
      //   element: <Session/>,
      // },
      // {
      //   path: '/account',
      //   element: <Accounts/>,
      // },
      // {
      //   path: '/clients',
      //   element: <Clients/>,
      // },
      // {
      //   path: '/provider',
      //   element: <Provider/>,
      // },
      // {
      //   path: '/access-key',
      //   element: <AccessKey/>,
      // },
    ]
  },
  {
    path: '/login',
    element: <div/>,
  },
]);
