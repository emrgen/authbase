import {Box, Heading, Stack, Table} from "@chakra-ui/react";

export const Users = () => {
    return (
        <Stack p={4} w={'full'} h='full' gap={4}>
            <Heading>Users</Heading>
            <Box flex={1}>
              <Box overflow={'auto'} h={'full'}>
                  <Table.Root>
                      <Table.Header>
                          <Table.Row>
                              <Table.ColumnHeader>
                                  Username
                              </Table.ColumnHeader>
                              <Table.ColumnHeader>
                                  Email
                              </Table.ColumnHeader>
                              <Table.ColumnHeader>
                                  Last Login
                              </Table.ColumnHeader>
                              <Table.ColumnHeader>
                                  Status
                              </Table.ColumnHeader>
                              <Table.ColumnHeader w={20}>
                                  Action
                              </Table.ColumnHeader>
                          </Table.Row>
                      </Table.Header>

                      <Table.Body>
                          {Array.from({length: 100}).map(() => (
                              <Table.Row>
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
              </Box>
            </Box>
        </Stack>
    )
}