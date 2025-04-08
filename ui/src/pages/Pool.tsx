import {Box, Button, Heading, HStack, Stack, Table} from "@chakra-ui/react"
import {usePoolStore} from "../store/pool.ts";

export const Pool = () => {
  return (
    <Box h={'full'} w={'full'}>
      <PoolTable/>
    </Box>
  )
}

export const PoolTable = () => {
  const pools = usePoolStore(state => state.pools);
  return (
    <Stack p={4} pos={'relative'} w={'full'} h='full' gap={4}>
      <HStack px={2} justifyContent={'space-between'}>
        <Heading>
          Pools
        </Heading>
        <Button size={'sm'} colorScheme={'blue'}>Create Pool</Button>
      </HStack>
      <Box flex={1} overflow={'hidden'}>
        <Box h={'full'}>
          <Table.ScrollArea height={'100%'}>
            <Table.Root stickyHeader>
              <Table.Header>
                <Table.Row>
                  <Table.ColumnHeader fontWeight={'bold'}>
                    Pool Name
                  </Table.ColumnHeader>
                  <Table.ColumnHeader fontWeight={'bold'}>
                    Pool ID
                  </Table.ColumnHeader>
                  <Table.ColumnHeader fontWeight={'bold'}>
                    Created At
                  </Table.ColumnHeader>
                  <Table.ColumnHeader fontWeight={'bold'} w={20}>
                    Actions
                  </Table.ColumnHeader>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {/* Add your table rows here */}
                {pools.map((pool) => (
                  <Table.Row key={pool.id}>
                    <Table.Cell>
                      {pool.name}
                    </Table.Cell>
                    <Table.Cell>
                      {pool.id}
                    </Table.Cell>
                    <Table.Cell>
                      {/*{pool.created_at}*/}
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