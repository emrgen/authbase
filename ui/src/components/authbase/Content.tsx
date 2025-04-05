import { HStack } from '@chakra-ui/react'
import React from 'react'

interface ContentProps {
  children: React.ReactNode
}

// This is the main content area of the application
export default function Content(props: ContentProps) {
  return (
    <HStack flex={1} h={'full'}>{props.children}</HStack>
  )
}
