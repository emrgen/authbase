import {Box} from "@chakra-ui/react";

export const SidebarItem = ({ children, isActive }) => {
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