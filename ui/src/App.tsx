import "./App.css";
import {AppSidebar} from "@/components/sidebar/app-sidebar.tsx";
import {SidebarInset, SidebarProvider, SidebarTrigger} from "@/components/ui/sidebar.tsx";
import {useUserStore} from "@/store/user.ts";
import {useEffect} from "react";
import {Outlet, useNavigate} from "react-router";
import {authbase} from "./api/client.ts";

import './main.styl';
import {rotateAccessToken} from "./service/authbase.ts";
import {useProjectStore} from "./store/project.ts";

// This is the main entry point of the application
function App() {
  const navigate = useNavigate()
  const isAuthenticated = useUserStore((state) => state.isAuthenticated);
  const setProjects = useProjectStore((state) => state.setProjects);
  const setListProjectState = useProjectStore((state) => state.setListProjectState);

  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, navigate]);

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
    <SidebarProvider>
      <AppSidebar/>
      <SidebarInset>
        <header
          className="fixed flex h-16 shrink-0 items-center gap-2 transition-[width,height] ease-linear group-has-[[data-collapsible=icon]]/sidebar-wrapper:h-12 bg-white z-10 w-full
             shadow-sm">
          <div className="flex items-center gap-2 px-4">
            <SidebarTrigger className="-ml-1"/>
          </div>
        </header>
        <main className="flex-1 px-4 mt-20">
         <Outlet/>
        </main>
      </SidebarInset>
    </SidebarProvider>
    // <Layout>
    //   <Sidebar>
    //     {/*show projects only for the admins*/}
    //     <Heading>
    //       <ProjectSelect/>
    //     </Heading>
    //     {/*show pools only for non admin users*/}
    //     <Heading>
    //       <PoolSelect/>
    //     </Heading>
    //     <Separator/>
    //     <Stack h={'full'}>
    //       <Stack flex={1}>
    //         <SidebarItem path={'/'}>Dashboard</SidebarItem>
    //         <SidebarItem path={'/pool'}>Pools</SidebarItem>
    //         <SidebarItem path={'/account'}>Accounts</SidebarItem>
    //         <SidebarItem path={'/session'}>Sessions</SidebarItem>
    //         <SidebarItem path={'/provider'}>Providers</SidebarItem>
    //         <SidebarItem path={'/clients'}>Clients</SidebarItem>
    //         <SidebarItem path={'/access-key'}>Access Tokens</SidebarItem>
    //       </Stack>
    //
    //       <SidebarItem>
    //         <HStack gap={4}>
    //           <IconButton borderRadius={'50%'} size={'xs'} bg={'gray.200'}/>
    //           <Text>
    //             Username
    //           </Text>
    //         </HStack>
    //       </SidebarItem>
    //     </Stack>
    //   </Sidebar>
    //   <Content>
    //     <Outlet/>
    //   </Content>
    // </Layout>
  );
}

export default App;
