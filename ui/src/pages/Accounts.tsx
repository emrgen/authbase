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
import {usePoolStore} from "@/store/pool.ts";

export const Accounts = () => {
  useListAccounts()
  return (
    <div>
      <AccountsTable/>
    </div>
  )
}

const AccountsTable = () => {
  const accounts = useAccountStore(state => state.accounts);
  const activePool = usePoolStore(state => state.activePool);

  return (
    <Table>
      {/*<TableCaption>List of your accounts in the <b>{activePool?.name ?? '-'}</b> pool</TableCaption>*/}
      <TableHeader className={'bg-gray-50'}>
        <TableRow>
          <TableHead className="w-[200px]">Username</TableHead>
          <TableHead>Email</TableHead>
          <TableHead>Status</TableHead>
          <TableHead>Pool</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {accounts.map((account) => (
          <TableRow key={account.id}>
            <TableCell className="font-medium">{account.name}</TableCell>
            <TableCell className="font-medium">{account.email}</TableCell>
            <TableCell className="font-medium">{account.status ?? '-'}</TableCell>
            <TableCell className="font-medium">{activePool?.name ?? '-'}</TableCell>
            <TableCell className="text-right">Edit</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}