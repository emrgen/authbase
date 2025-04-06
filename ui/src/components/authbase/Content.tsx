import {Stack} from '@chakra-ui/react'
import React from 'react'

interface ContentProps {
  children: React.ReactNode
}

// This is the main content area of the application
export default function Content(props: ContentProps) {
  return (
    <Stack flex={1} h="100vh" p={0} pos={'relative'}>
      {props.children}
    </Stack>
  )
}
