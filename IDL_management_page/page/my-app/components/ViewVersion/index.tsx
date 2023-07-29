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
import { toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

import { Entry } from "../../types/entry";
import { Version } from "../../types/version";
import { Button } from "@chakra-ui/react";
import axios from "axios";
import { ChevronLeftIcon } from "@chakra-ui/icons";
import router from "next/router";
import { confirmAlert } from 'react-confirm-alert'; // Import
import 'react-confirm-alert/src/react-confirm-alert.css'; // Import css

export default function ViewVersionComponent({ id, vname, idlfile, microserviceId}: Version) {
    const handleDownload = async (e: any) => {
        e.preventDefault();

        const formData = new FormData();
        formData.append('vname', vname)
        formData.append('microserviceId', microserviceId)

        axios
        .post("http://localhost:3333/findSvcVerIDL", formData, {
            headers: {
            "Content-Type": "multipart/form-data",
            },
        })
        .then((response: any) => {
            const byteCharacters = Uint8Array.from(response.data.idlfile.data);
            const blob = new Blob([byteCharacters]);
            const downloadLink = document.createElement('a');
            downloadLink.href = URL.createObjectURL(blob);
            downloadLink.download = 'output.thrift';
            downloadLink.click();
            // Handle the upload success
        })
        .catch((error) => {
            console.log(error); //POSSIBLE ERRORS: NOT UNIQUE SVC NAME
            // Handle the upload error
        });
        };

        // const handleDeleteVer = async (e: any) => {
        //   e.preventDefault();

        //   confirmAlert({
        //     title: 'Confirm to delete',
        //     message: 'Are you sure you want to delete this version?',
        //     buttons: [
        //       {
        //         label: 'Yes',
        //         onClick: async () => {
        //           // ... your delete function
        //           const formData = new FormData();
        //         formData.append('vname', vname)
        //         formData.append('microserviceId', microserviceId)
                
        //       axios
        //       .post("http://localhost:3333/delSvcVer", formData, {
        //           headers: {
        //           "Content-Type": "multipart/form-data",
        //           },
        //       })
        //       .then((response: any) => {
        //           // Handle the deletion success
        //             if (response.data.message === 'Version deleted successfully.') {
        //                     // redirect to the homepage
        //                     const deletedName = response.data.name;

        //                     const toastId = toast.success(`${deletedName} deleted successfully`);

        //                     const redirectInterval = setInterval(() => {
                              
        //                     if (!toast.isActive(toastId)) {
        //                         // If the toast message is no longer displayed, clear the interval and redirect to '/'
        //                         clearInterval(redirectInterval);
        //                         router.push(`/viewServices/${deletedName}`);
                                  
        //                         }
        //                       }, 500);
        //                     }        
                        

        //             if (response.data.message === 'Version and microservice deleted successfully.') {
        //               const toastId = toast.success(`Service deleted successfully`);

        //               const redirectInterval = setInterval(() => {
                        
        //               if (!toast.isActive(toastId)) {
        //                   // If the toast message is no longer displayed, clear the interval and redirect to '/'
        //                   clearInterval(redirectInterval);
        //                   router.push(`/`);
                            
        //                   }
        //                 }, 500);
        //               router.push(`/`);
        //             }
        //           }).catch((error) => {
        //                   console.log(error);
        //                   // Handle the deletion error
        //                   if (error.response) {
        //                     toast.error('Deletion failed');
        //                   }
        //               });
        //             }
        //           },
        //       {
        //         label: 'No',
        //         onClick: () => {}
        //       }
        //     ]
        //   });
        // };
      
          
      ;
      

  return (
    <>
      <BackgroundComponent></BackgroundComponent>
      <Container>
        <SubContainer>
          <FlexInputs>
            <TitleText>
            <MyIcon as={ChevronLeftIcon} onClick={() => router.back()} />
            Version Details</TitleText>

            <VStackSvcInfo>
              <div>
                <TextCategory>Microservice Version Name</TextCategory>
                <TextInfo>{vname}</TextInfo>
              </div>
              <div>
                <Button onClick={handleDownload} colorScheme="blue">Download IDL File</Button>
              </div>
              {/* <div>
                <Button onClick={handleDeleteVer} colorScheme="blue">Delete Version</Button>
              </div> */}
            </VStackSvcInfo>
          </FlexInputs>
        </SubContainer>
        </Container>
    </>
  );
}
