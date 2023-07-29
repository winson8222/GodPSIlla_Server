//main page
import { AppProps } from "next/app";
import SelectServiceComponent from "../components/SelectService";
import { ChakraProvider } from "@chakra-ui/react";
import 'react-toastify/dist/ReactToastify.css';
import { ToastContainer } from 'react-toastify';

function MyApp({ Component, pageProps }: AppProps) {
    return (
        <ChakraProvider>
            <Component {...pageProps} />
            < ToastContainer />
        </ChakraProvider>
    )
}

export default MyApp;