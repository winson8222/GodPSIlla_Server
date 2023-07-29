import { Box, Flex, VStack, chakra, Text } from "@chakra-ui/react";

export const ContainerVStack = chakra(VStack, {
    baseStyle: {
    }
})

export const SelectText = chakra(Text, {
    baseStyle: {
        fontWeight: 'bold',
        fontSize: 'lg', //["2xs", "2xs", "lg"] for responsive
    }
})

export const AddNewText = chakra(Text, {
    baseStyle: {
        fontWeight: 'bold',
        fontSize: 'lg', //["2xs", "2xs", "lg"] for responsive
        textDecor: 'underline',
        _hover: { 
            cursor: 'pointer',
            textDecor: 'None'
        }
    }
})