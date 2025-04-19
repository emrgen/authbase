import {Box, Button, Heading, HStack, Stack, Table} from "@chakra-ui/react";
import {useAccountStore} from "../store/account.ts";
import {useListAccounts} from "../store/hooks/account.tsx";

export const Accounts = () => {
  useListAccounts()

  return (
    <Box h={'full'} w={'full'}>
      <AccountTable/>
    </Box>
  )
}

const AccountTable = () => {
  const accounts = useAccountStore(state => state.accounts);
  return (
    <Stack p={4} pos={'relative'} w={'full'} h='full' gap={4}>
      <HStack px={2} justifyContent={'space-between'}>
        <Heading>
          Accounts
        </Heading>
        <Button size={'sm'} colorScheme={'blue'}>Create Account</Button>
      </HStack>
      <Box flex={1} overflow={'hidden'}>
        <Box h={'full'}>
          <Table.ScrollArea height={'100%'}>
            <Table.Root stickyHeader>
              <Table.Header>
                <Table.Row>
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
                {accounts.map((account) => (
                  <Table.Row key={account.id}>
                    <Table.Cell>
                      {account.username}
                    </Table.Cell>
                    <Table.Cell>
                      {account.email}
                    </Table.Cell>
                    <Table.Cell>
                      {/*{account.last_login}*/}
                    </Table.Cell>
                    <Table.Cell>
                      Active
                    </Table.Cell>
                    <Table.Cell>
                      <Button size={'2xs'}>Edit</Button>
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