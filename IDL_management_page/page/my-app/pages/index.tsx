//main page
import { AppProps } from "next/app";
import SelectServiceComponent from "../components/SelectService";
import { ChakraProvider } from "@chakra-ui/react";
// import { entryInfoArray } from "../data/db";
import { Entry } from "../types/entry";
import { PrismaClient } from '@prisma/client'


function Home(props: any) {
    return (
        <SelectServiceComponent entries={props.response}/>
    )
}

export async function getServerSideProps() {
    //upon router.push() from
    const prisma = new PrismaClient();
    const response = await prisma.microservice.findMany();
    console.log(response);
    // const userInfo = await response.json();
  
    return {
      props: {
        response
      },
    };
  }

export default Home;