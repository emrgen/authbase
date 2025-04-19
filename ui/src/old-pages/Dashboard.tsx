import {Box, Heading, HStack, Stack} from "@chakra-ui/react";

export const Dashboard = () => {
  return (
    <Box h={'full'} w={'full'}>
      <Stack p={4} pos={'relative'} w={'full'} h='full' gap={4}>
        <HStack px={2}>
          <Heading>
            Dashboard
          </Heading>
        </HStack>
        <Box flex={1} overflow={'hidden'}>
          <Box h={'full'}>
            {/* Add your dashboard content here */}
          </Box>
        </Box>
      </Stack>
    </Box>
  )
}