import { Box, Flex, VStack, chakra, Text, extendTheme } from "@chakra-ui/react";

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
        flexDirection: 'column',
        justifyContent: 'center',
        bg: 'white',
        borderRadius: '2xl',
        boxShadow: '0 3px 6px rgba(0,0,0,0.16)',
      },
})

export const SubContainer2 = chakra(Flex, {
  baseStyle: {
      width: '40%',
      height: '60%',
      flexDirection: 'column',
      justifyContent: 'center',
      bg: 'white',
      borderRadius: '2xl',
      boxShadow: '0 3px 6px rgba(0,0,0,0.16)',
    },
})

export const VerticalContainer = {
  display: 'flex',
  gap: '1rem',
};

export const SubContainerStyle = {
  width: '100%', // or any fixed width you prefer
};


export const SubContainer3 = chakra(Flex, {
  baseStyle: {
      width: '10%',
      height: '50%',
      flexDirection: 'column',
      justifyContent: 'center',
      bg: 'white',
      borderRadius: '2xl',
      boxShadow: '0 3px 6px rgba(0,0,0,0.16)',
    },
})


export const TitleText = chakra(Text, {
  baseStyle: {
      fontWeight: 'bold',
      fontSize: 'lg', //["2xs", "2xs", "lg"] for responsive
  }
})

export const theme = extendTheme({
  components: {
    Button: {
      variants: {
        disabled: {
          bg: "gray",
          color: "white",
          _hover: {
            bg: "gray",
          },
          _focus: {
            bg: "gray",
            boxShadow: "none",
          },
          _active: {
            bg: "gray",
          },
        },
      },
    },
  },
});
