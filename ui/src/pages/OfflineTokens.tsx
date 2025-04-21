import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table"
import {useAccountStore} from "@/store/account.ts";
import {useListAccounts} from "@/store/hooks/account.tsx";
import {useOfflineTokenStore} from "@/store/offline-tokens.ts";
import {usePoolStore} from "@/store/pool.ts";

export const OfflineTokens = () => {
  return (
    <div>
      <OfflineTokensTable/>
    </div>
  )
}

const OfflineTokensTable = () => {
  const tokens = useOfflineTokenStore(state => state.tokens);

  return (
    <Table>
      {/*<TableCaption>List of your accounts</TableCaption>*/}
      <TableHeader className={'bg-gray-50'}>
        <TableRow>
          <TableHead className="w-[300px]">Name</TableHead>
          <TableHead>Created By</TableHead>
          <TableHead>Created At</TableHead>
          <TableHead>Expires At</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {tokens.map((token) => (
          <TableRow key={token.id}>
            <TableCell className="font-medium">{token.name}</TableCell>
            <TableCell className="font-medium">{token.createdBy}</TableCell>
            <TableCell className={"font-medium"}>
              {token.createdAt}
            </TableCell>
            <TableCell className="font-medium">{token.expiresAt}</TableCell>
            <TableCell className="text-right">Edit</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}