import {Heading, HStack, IconButton, Separator, Stack, Text} from "@chakra-ui/react";
import "./App.css";
import Content from "./components/authbase/Content";
import {Layout} from "./components/authbase/Layout";
import {ProjectSelect} from "./components/authbase/ProjectSelect.tsx";
import {Sidebar} from "./components/authbase/Sidebar";
import {SidebarItem} from "./components/authbase/SidebarItem.tsx";
import {Users} from "./pages/Users.tsx";
import './main.styl';

// This is the main entry point of the application
function App() {

  return (
    <Layout>
      <Sidebar>
        {/*show projects only for the admins*/}
        <Heading>
          <ProjectSelect/>
        </Heading>
        <Separator/>
        {/*show pools only for non admin users*/}
        {/*<Heading>*/}
        {/*  <PooltSelect/>*/}
        {/*</Heading>*/}
        <Separator/>
        <Stack h={'full'}>
          <Stack flex={1}>
            <SidebarItem isActive={true}>Dashboard</SidebarItem>
            <SidebarItem>Accounts</SidebarItem>
            <SidebarItem>Providers</SidebarItem>
            <SidebarItem>Access Token</SidebarItem>
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
        <Users/>
      </Content>
    </Layout>
  );
}

export default App;
