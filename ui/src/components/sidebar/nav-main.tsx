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
import {BookOpen, Bot, ChevronRight, Users} from "lucide-react"
import {TbBrandOauth} from "react-icons/tb";

const navMain: any[] = [
  {
    title: "Accounts",
    url: "#",
    icon: Users,
    isActive: true,
  },
  {
    title: "Clients",
    url: "#",
    icon: Bot,
  },
  {
    title: "Offline tokens",
    url: "#",
    icon: BookOpen,
  },
  {
    title: "Providers",
    url: "#",
    icon: TbBrandOauth,
  },
];


export function NavMain() {
  return (
    <SidebarGroup>

      <SelectPool/>

      <div className={'h-4'}/>

      {/*<SidebarGroupLabel>Pool</SidebarGroupLabel>*/}
      <SidebarMenu>
        {navMain.map((item) => {
          if (!item.items?.length) {
            return (
              <SidebarMenuItem key={item.title} className={'cursor-pointer'}>
                <SidebarMenuButton tooltip={item.title}>
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
