import { HStack } from '@chakra-ui/react'
import React from 'react'

interface ContentProps {
  children: React.ReactNode
}

export default function Content(props: ContentProps) {
  return (
    <HStack flex={1} h={'full'} bg='red'>{props.children}</HStack>
  )
}
