import { Box, Flex, VStack, chakra, Input, Text, FormControl, Icon } from "@chakra-ui/react";

export const Container = chakra(Flex, {
    baseStyle: {
        width: '100vw',
        height: '100vh',
        pb: '10%',
        flexDirection: 'column',
        justifyContent: 'center',
        flexWrap: 'wrap',
        alignItems: 'center',
        alignContent: 'center',
      },
}) 

export const SubContainer = chakra(Flex, {
    baseStyle: {
        width: '30%',
        height: '50%',
        pl: '5%',
        pr: '5%',
        flexDirection: 'column',
        justifyContent: 'center',
        bg: 'white',
        borderRadius: '2xl',
        boxShadow: '0 3px 6px rgba(0,0,0,0.16)',
      },
})

export const VStackInputs = chakra(VStack, {
  baseStyle: {

  }
})

export const TitleText = chakra(Text, {
  baseStyle: {
      fontWeight: 'bold',
      fontSize: 'lg', //["2xs", "2xs", "lg"] for responsive
  }
})

export const FlexInputs = chakra(Flex, {
  baseStyle: {
    flexDir: 'column',

    "& > *": {
      pt: '10%',
    }
  }
})

export const VStackSvcInfo = chakra(FormControl, {
  baseStyle: {
    "& > *": {
      // Styles for child elements
      mb: '5%',
    },
  }
})

export const TextCategory = chakra(Text, {
    baseStyle: {
        fontWeight: 'medium',
        fontSize: 'lg', //["2xs", "2xs", "lg"] for responsive
    }
  })

  export const TextInfo = chakra(Text, {
    baseStyle: {
        fontWeight: 'small',
        fontStyle: 'italic',
        fontSize: 'lg', //["2xs", "2xs", "lg"] for responsive
    }
  })

export const VStackForm = chakra(VStack, {
  baseStyle: {
  }
})

export const SvcNameInput = chakra(Input, {
    baseStyle: {
      },
})

export const SvcUpstreamInput = chakra(Input, {
  baseStyle: {
    },
})

export const MyIcon = chakra(Icon, {
  baseStyle: {
    boxSize: '6',
    _hover: { 
      cursor: 'pointer',
      color: 'gray'
  }
  }
})