import { HStack } from "@chakra-ui/react"

interface LayoutProps{
  children: React.ReactNode
}

// This is the main layout component for the application
// It is used to wrap the sidebar and content components and provide a consistent layout
export const Layout = (props: LayoutProps) => {
  return (
    <HStack h='full' align={'start'} gap={0} pos={'relative'} overflow={'hidden'}>
      {props.children}
    </HStack>
  )
}
