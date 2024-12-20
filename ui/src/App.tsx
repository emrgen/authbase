import { Heading } from "@chakra-ui/react";
import "./App.css";
import Content from "./components/authbase/Content";
import { Layout } from "./components/authbase/Layout";
import { OrgSelect } from "./components/authbase/OrgSelect";
import { Sidebar } from "./components/authbase/Sidebar";

function App() {
  return (
    <>
      <Layout>
        <Sidebar>
          <Heading>
            <OrgSelect />
          </Heading>
        </Sidebar>
        <Content>1231</Content>
      </Layout>
    </>
  );
}

export default App;
