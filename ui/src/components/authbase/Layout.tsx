import { HStack } from "@chakra-ui/react"

interface LayoutProps{
  children: React.ReactNode
}

export const Layout = (props: LayoutProps) => {
  return (
    <HStack h='full' align={'start'} gap={0}>
      {props.children}
    </HStack>
  )
}
