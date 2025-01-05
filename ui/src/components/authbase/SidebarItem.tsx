import {Box} from "@chakra-ui/react";
import {ReactNode} from "react";

interface SidebarItemProps {
    children: ReactNode;
    isActive?: boolean;
}

export const SidebarItem = (props: SidebarItemProps) => {
    const {children, isActive = false} = props;

    return (
        <Box
            _hover={{ bg: "gray.700" }}
            p={2}
            bg={isActive ? "gray.700" : "transparent"}
            borderRadius="md"
            cursor="pointer"
        >
            {children}
        </Box>
    );
}