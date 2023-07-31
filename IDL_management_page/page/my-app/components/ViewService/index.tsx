import {
    Container,
    VStackSvcInfo,
    SubContainer,
    TextCategory,
    TextInfo,
    TitleText,
    FlexInputs,
    MyIcon,
  } from "./styles";

  import BackgroundComponent from "../Background";

import { Entry } from "../../types/entry";
import IDLDropdownComponent from "./IDLDropdown";
import { Button, HStack } from "@chakra-ui/react";
import { useRouter } from "next/router";
import { ChevronLeftIcon } from "@chakra-ui/icons";
import axios from "axios";
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
  
  export default function ViewServiceComponent({ id, svcname, versions }: Entry) {
    const router = useRouter();

    const handleNewVerClick = (e: any) => {
      e.preventDefault();
      
      const newURL = `${router.asPath}/newVersion`;
      router.push(newURL, undefined, { shallow: true });
    };

    const handleDeleteSvc = (e: any) => {
      e.preventDefault();

      const formData = new FormData();
      formData.append('microserviceId', id)
      axios
      .post("http://localhost:3333/deleteSvc", formData, {
          headers: {
          "Content-Type": "multipart/form-data",
          },
      })
      .then((response: any) => {
        if (response.status === 200) {

          const toastId = toast.success('Service deleted successfully');

          router.push('/');
        }
      })
      .catch((error) => {
          console.log(error); //POSSIBLE ERRORS: NOT UNIQUE SVC NAME
          // Handle the upload error
          if (error.response) {
            toast.error('Deletion failed');
          }
      });
    }

    return (
      <>
        <BackgroundComponent></BackgroundComponent>
        <Container>
          <SubContainer>
            <FlexInputs>
              
              <TitleText>
              <MyIcon as={ChevronLeftIcon} onClick={() => router.back()} />Microservice Details
              </TitleText>
  
              <VStackSvcInfo>
                <div>
                  <TextCategory>Microservice Name</TextCategory>
                  <TextInfo>{svcname}</TextInfo>
                </div>  

                <HStack>
                <IDLDropdownComponent id={id} svcname={svcname} versions={versions}/>
                <Button onClick={handleNewVerClick} size="md" colorScheme="blue">New Version</Button>
                
                </HStack>
                <Button onClick={handleDeleteSvc} size="md" colorScheme="red">Delete Service</Button>
              </VStackSvcInfo>
            </FlexInputs>
          </SubContainer>
        </Container>
      </>
    );
  }
  