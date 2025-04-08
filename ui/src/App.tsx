import {Heading, HStack, IconButton, Separator, Stack, Text} from "@chakra-ui/react";
import "./App.css";
import {useEffect} from "react";
import {Outlet} from "react-router";
import {authbase} from "./api/client.ts";
import Content from "./components/authbase/Content";
import {Layout} from "./components/authbase/Layout";
import {PoolSelect} from "./components/authbase/PoolSelect.tsx";
import {ProjectSelect} from "./components/authbase/ProjectSelect.tsx";
import {Sidebar} from "./components/authbase/Sidebar";
import {SidebarItem} from "./components/authbase/SidebarItem.tsx";
import './main.styl';
import {rotateAccessToken} from "./service/authbase.ts";
import {useProjectStore} from "./store/project.ts";

// This is the main entry point of the application
function App() {
  const setProjects = useProjectStore((state) => state.setProjects);
  const setListProjectState = useProjectStore((state) => state.setListProjectState);
  useEffect(() => {
    console.log("App mounted");
    rotateAccessToken();
  }, []);

  useEffect(() => {
    setListProjectState('loading');
    authbase.project.listProjects({}).then((res) => {
      const {data} = res;
      const projects = data.projects?.map((project) => ({
        id: project.id!,
        name: project.name!,
      })) || [];

      console.log("Projects", data);
      setProjects(projects);
    }).finally(() => {
      setListProjectState('success');
    })
  }, [setListProjectState, setProjects]);

  return (
    <Layout>
      <Sidebar>
        {/*show projects only for the admins*/}
        <Heading>
          <ProjectSelect/>
        </Heading>
        <Separator/>
        {/*show pools only for non admin users*/}
        <Heading>
          <PoolSelect/>
        </Heading>
        <Separator/>
        <Stack h={'full'}>
          <Stack flex={1}>
            <SidebarItem path={'/'}>Dashboard</SidebarItem>
            <SidebarItem path={'/account'}>Accounts</SidebarItem>
            <SidebarItem path={'/provider'}>Providers</SidebarItem>
            <SidebarItem path={'/access-key'}>Access Token</SidebarItem>
          </Stack>

          <SidebarItem>
            <HStack gap={4}>
              <IconButton borderRadius={'50%'} size={'xs'}/>
              <Text>
                Username
              </Text>
            </HStack>
          </SidebarItem>
        </Stack>
      </Sidebar>
      <Content>
        <Outlet/>
      </Content>
    </Layout>
  );
}

export default App;
