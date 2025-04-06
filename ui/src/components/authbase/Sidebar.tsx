import { Stack } from '@chakra-ui/react'
import React from 'react'

interface SidebarProps {
  children: React.ReactNode
}

// Sidebar component that wraps the sidebar items
export function Sidebar(props: SidebarProps) {
  return (
    <Stack maxW="240px" w="240px" h="100vh" color="white" p={3} py={4}>
      {props.children}
    </Stack>
  )
}
