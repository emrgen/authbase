import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import {RouterProvider} from "react-router";
import { Provider } from "./components/ui/provider.tsx";
import {router} from "./router.tsx";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Provider themes={['dark', 'light']}>
      <RouterProvider router={router} />
    </Provider>
  </StrictMode>
);
