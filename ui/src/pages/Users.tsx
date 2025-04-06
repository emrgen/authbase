import {Box, Heading, HStack, Stack, Table} from "@chakra-ui/react";

export const Users = () => {

  return (
    <Stack p={4} pos={'relative'} w={'full'} h='full' gap={4}>
      <HStack px={2}>
        <Heading>
          Users
        </Heading>
      </HStack>
      <Box flex={1} overflow={'hidden'}>
        <Box h={'full'}>
          <Table.ScrollArea height={'100%'}>
            <Table.Root stickyHeader>
              <Table.Header>
                <Table.Row bg={'#333'}>
                  <Table.ColumnHeader fontWeight={'bold'}>
                    Username
                  </Table.ColumnHeader>
                  <Table.ColumnHeader fontWeight={'bold'}>
                    Email
                  </Table.ColumnHeader>
                  <Table.ColumnHeader fontWeight={'bold'}>
                    Last Login
                  </Table.ColumnHeader>
                  <Table.ColumnHeader fontWeight={'bold'}>
                    Status
                  </Table.ColumnHeader>
                  <Table.ColumnHeader w={20} fontWeight={'bold'}>
                    Action
                  </Table.ColumnHeader>
                </Table.Row>
              </Table.Header>

              <Table.Body>
                {Array.from({length: 100}).map(() => (
                  <Table.Row key={Math.random()}>
                    <Table.Cell>
                      admin
                    </Table.Cell>
                    <Table.Cell>
                      admin@mail.com
                    </Table.Cell>
                    <Table.Cell>
                      2021-10-01 12:00:00
                    </Table.Cell>
                    <Table.Cell>
                      Active
                    </Table.Cell>
                    <Table.Cell>
                      Edit
                    </Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table.Root>
          </Table.ScrollArea>
        </Box>
      </Box>
    </Stack>
  )
}