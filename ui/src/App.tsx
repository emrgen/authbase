import {Heading, Separator} from "@chakra-ui/react";
import "./App.css";
import Content from "./components/authbase/Content";
import {Layout} from "./components/authbase/Layout";
import {OrgSelect} from "./components/authbase/OrgSelect";
import {Sidebar} from "./components/authbase/Sidebar";
import {SidebarItem} from "./components/authbase/SidebarItem.tsx";
import {Users} from "./pages/Users.tsx";

function App() {
    return (
        <>
            <Layout>
                <Sidebar>
                    <Heading>
                        <OrgSelect/>
                    </Heading>
                    <Separator/>
                    <SidebarItem>Organizations</SidebarItem>
                    <SidebarItem>Sessions</SidebarItem>
                    <SidebarItem isActive={true}>Users</SidebarItem>
                    <SidebarItem>Members</SidebarItem>
                    <SidebarItem>Providers</SidebarItem>
                    <SidebarItem>Tokens</SidebarItem>
                </Sidebar>
                <Content>
                    <Users/>
                </Content>
            </Layout>
        </>
    );
}

export default App;
