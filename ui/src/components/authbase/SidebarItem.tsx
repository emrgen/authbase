import {Box} from "@chakra-ui/react";
import {ReactNode, useEffect, useState} from "react";
import {useLocation, useNavigate} from "react-router";

interface SidebarItemProps {
    children: ReactNode;
    isActive?: boolean;
    path?: string;
}

export const SidebarItem = (props: SidebarItemProps) => {
    const {children, path} = props;
    const navigate = useNavigate();
    const [isActive, setIsActive] = useState(false);
    const location = useLocation();

    useEffect(() => {
        setIsActive(location.pathname === path);
    }, [location.pathname, path]);

    return (
        <Box
            _hover={{ bg: "gray.100" }}
            p={2}
            bg={isActive ? "gray.100" : "transparent"}
            borderRadius="md"
            cursor="pointer"
            onClick={() => {
                if (path) {
                    navigate(path);
                }
            }}
        >
            {children}
        </Box>
    );
}