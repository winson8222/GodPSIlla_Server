import { Box, chakra, Flex } from "@chakra-ui/react";

export const Container = chakra(Flex, {
    baseStyle: {
        width: '100vw',
        height: '100vh',
        flexDirection: 'row',
        justifyContent: 'space-between',
        flexWrap: 'nowrap',
        bg: '#044454',
        position: 'absolute',
        overflow: 'hidden',
        zIndex: '-1'
      },
})

export const BoxOne = chakra(Box, {
    baseStyle: {
        width: '50vw',
        height: '100vh',
        bg: 'orange',
        right: '65vw',
        clipPath: 'polygon(25% 0%, 100% 0, 75% 100%, 0% 100%)',
        position: 'absolute',
      },
})

export const BoxThree = chakra(Box, {
    baseStyle: {
        width: '50vw',
        height: '100vh',
        bg: '#047f84',
        left: '65vw',
        clipPath: 'polygon(25% 0%, 100% 0, 75% 100%, 0% 100%)',
        position: 'absolute',
      },
})