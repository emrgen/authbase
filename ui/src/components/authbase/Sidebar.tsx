import { Stack } from '@chakra-ui/react'
import React from 'react'

interface SidebarProps {
  children: React.ReactNode
}

export function Sidebar(props: SidebarProps) {
  return (
    <Stack maxW="240px" w="240px" h="100vh" bg="gray.800" color="white" p={4}>
      {props.children}
    </Stack>
  )
}
