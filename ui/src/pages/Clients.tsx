import {
  Table,
  TableBody,
  TableCaption,
  TableCell,
  TableHead,
  TableHeader,
  TableRow
} from "@/components/ui/table"
import {useClientStore} from "@/store/clients.ts";
import {useListClients} from "@/store/hooks/client.tsx";
import { toast } from "sonner";
import {usePoolStore} from "@/store/pool.ts";
import {Button} from "@/components/ui/button";
import dayjs from "dayjs";
import {ClipboardCopy,} from "lucide-react";

export const Clients = () => {
  useListClients();
  return (
    <div>
      <ClientsTable/>
    </div>
  )
}

const ClientsTable = () => {
  const clients = useClientStore(state => state.clients);
  const activePool = usePoolStore(state => state.activePool);

  return (
    <Table>
      <TableCaption>A list of your accounts in pool: {activePool?.name ?? '-'}</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead className="w-[200px]">Name</TableHead>
          <TableHead>Client ID</TableHead>
          <TableHead>Client Secret</TableHead>
          <TableHead>Created At</TableHead>
          <TableHead>Created By</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {clients.map((client) => (
          <TableRow key={client.id}>
            <TableCell className="font-medium">{client.name}</TableCell>
            <TableCell className="font-medium w-[300px]">
              <div className={'flex gap-2 items-center justify-normal'}>
                <div
                  className={'text-xs text-gray-500 overflow-hidden whitespace-nowrap overflow-ellipsis'}>
                  {client.id}
                </div>
                {/* tiny button with clipboard icons*/}
                <Button
                  variant={'ghost'} size={'xs'}
                  className={'text-gray-500 hover:text-gray-700 cursor-pointer'}
                  onClick={async () => {
                    await navigator.clipboard.writeText(client.id ?? '');
                  }}
                >
                  <ClipboardCopy className={'w-4 h-4'}/>
                </Button>
              </div>
            </TableCell>
            <TableCell className="font-medium">
              <div className={'flex gap-2 items-center justify-normal'}>
                <div
                  className={'text-xs text-gray-500 overflow-hidden whitespace-nowrap overflow-ellipsis w-[300px]'}>
                  {/*mask client secret*/}
                  {client.clientSecret?.replace(/.(?=.{4})/g, '*')}
                </div>
                {/* tiny button with clipboard icons*/}
                <Button
                  variant={'ghost'} size={'xs'}
                  className={'text-gray-500 hover:text-gray-700 cursor-pointer'}
                  onClick={async () => {
                    await navigator.clipboard.writeText(client.clientSecret ?? '');
                  }}
                >
                  <ClipboardCopy className={'w-4 h-4'}/>
                </Button>
              </div>
            </TableCell>
            <TableCell className="font-medium">
              {client.createdAt ? dayjs(client.createdAt).format('YYYY-MM-DD HH:mm:ss') : '-'}
            </TableCell>
            <TableCell className="font-medium">
              {client.createdBy ? client.createdBy : '-'}
            </TableCell>
            <TableCell className="text-right">Rename</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}