import {ChakraProvider, createSystem, defaultConfig, defineConfig} from "@chakra-ui/react";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import {RouterProvider} from "react-router";
import {router} from "./router.tsx";

const config = defineConfig({
  theme: {
    tokens: {
      colors: {},
    },
  },
})

const system = createSystem(defaultConfig, config)

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ChakraProvider value={system}>
      <RouterProvider router={router} />
    </ChakraProvider>
  </StrictMode>
);
