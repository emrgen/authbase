import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table"
import {useListProviders} from "@/store/hooks/provider.tsx";
import {useProviderStore} from "@/store/provider.ts";

export const  Providers = () => {
  useListProviders();
  return (
    <div>
      <ProvidersTable/>
    </div>
  )
}

const ProvidersTable = () => {
  const providers = useProviderStore(state => state.providers);

  return (
    <Table>
      {/*<TableCaption>List of your accounts</TableCaption>*/}
      <TableHeader className={'bg-gray-50'}>
        <TableRow>
          <TableHead className="w-[200px]">Name</TableHead>
          <TableHead className="text-right">Actions</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {providers.map((provider) => (
          <TableRow key={provider.name}>
            <TableCell className="font-medium">{provider.name}</TableCell>
            <TableCell className="text-right">Edit</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}