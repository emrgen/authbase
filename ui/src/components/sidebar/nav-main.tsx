"use client"

import {SelectPool} from "@/components/sidebar/pool.tsx";

import {Collapsible, CollapsibleContent, CollapsibleTrigger,} from "@/components/ui/collapsible"
import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
} from "@/components/ui/sidebar"
import {useAppStore} from "@/store/app.ts";
import {
  BookOpen,
  Bot,
  ChevronRight, IdCard,
  LayoutDashboardIcon,
  Unplug,
  Users,
  WifiOff
} from "lucide-react"
import {useEffect} from "react";
import {TbBrandOauth} from "react-icons/tb";
import {useLocation, useNavigate} from "react-router";

export const navMain: any[] = [
  {
    title: "Dashboard",
    url: "/",
    icon: LayoutDashboardIcon,
  },
  {
    title: "Accounts",
    url: "/account",
    icon: Users,
    isActive: true,
  },
  {
    title: "Clients",
    url: "/client",
    icon: Unplug,
  },
  {
    title: "Offline tokens",
    url: "/offline-token",
    icon: WifiOff,
  },
  {
    title: "Identity Providers",
    url: "/provider",
    icon: IdCard,
  },
];


export function NavMain() {
  const navigate = useNavigate();
  const setActiveSidebarItem = useAppStore((state) => state.setActiveSidebarItem);
  const activeSidebarItem = useAppStore((state) => state.activeSidebarItem);

  const location = useLocation();

  useEffect(() => {
    const path = location.pathname;
    const activeNavItem = navMain.find((item) => item.url.indexOf(path) !== -1);
    if (activeNavItem) {
      setActiveSidebarItem({
        title: activeNavItem.title,
        url: activeNavItem.url,
        icon: activeNavItem.icon,
      });
    } else {
      setActiveSidebarItem(null);
    }
  }, [location, setActiveSidebarItem]);


  return (
    <SidebarGroup>

      <SelectPool/>

      <div className={'h-4'}/>

      {/*<SidebarGroupLabel>Pool</SidebarGroupLabel>*/}
      <SidebarMenu>
        {navMain.map((item) => {
          if (!item.items?.length) {
            return (
              <SidebarMenuItem key={item.title}>
                <SidebarMenuButton
                  className={'cursor-pointer'}
                  isActive={item.url === activeSidebarItem?.url}
                  tooltip={item.title}
                  onClick={() => {
                    if (item.url) {
                      navigate(item.url);
                    }
                  }}>
                  {item.icon && <item.icon/>}
                  <span>{item.title}</span>
                </SidebarMenuButton>
              </SidebarMenuItem>
            )
          }

          return (
            <Collapsible
              key={item.title}
              asChild
              defaultOpen={item.isActive}
              className="group/collapsible"
            >
              <SidebarMenuItem>
                <CollapsibleTrigger asChild>
                  <SidebarMenuButton tooltip={item.title}>
                    {item.icon && <item.icon/>}
                    <span>{item.title}</span>
                    {item.items?.length && <ChevronRight
                      className="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"/>}
                  </SidebarMenuButton>
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <SidebarMenuSub>
                    {item.items?.map((subItem) => (
                      <SidebarMenuSubItem key={subItem.title}>
                        <SidebarMenuSubButton asChild>
                          <a href={subItem.url}>
                            <span>{subItem.title}</span>
                          </a>
                        </SidebarMenuSubButton>
                      </SidebarMenuSubItem>
                    ))}
                  </SidebarMenuSub>
                </CollapsibleContent>
              </SidebarMenuItem>
            </Collapsible>
          );
        })}
      </SidebarMenu>
    </SidebarGroup>
  )
}
