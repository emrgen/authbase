
import {Box, Heading, HStack, Stack} from "@chakra-ui/react"

export const AccessKey = () => {
  return (
    <Box h={'full'} w={'full'}>
      <Stack p={4} pos={'relative'} w={'full'} h='full' gap={4}>
        <HStack px={2}>
          <Heading>
            AccessKey
          </Heading>
        </HStack>
        <Box flex={1} overflow={'hidden'}>
          <Box h={'full'}>
            {/* Add your provider content here */}
          </Box>
        </Box>
      </Stack>
    </Box>
  )
}