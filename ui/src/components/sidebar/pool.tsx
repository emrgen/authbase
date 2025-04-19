import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

export const SelectPool = () => {
  return (
    <Select value={'default'}>
      <SelectTrigger className="w-full">
        <SelectValue placeholder="Select a fruit" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Pool</SelectLabel>
          <SelectItem value="default">Default</SelectItem>
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}