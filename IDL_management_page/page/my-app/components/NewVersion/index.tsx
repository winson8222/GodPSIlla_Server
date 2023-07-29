import {
    Container,
    MyFormControl,
    SubContainer,
    SvcNameInput,
    SvcUpstreamInput,
    VersionInput,
    FileInput,
    TitleText,
    FlexInputs,
    MyIcon,
  } from "./styles";
  import BackgroundComponent from "../Background";
  import {
    Input,
    FormControl,
    FormLabel,
    FormErrorMessage,
    FormHelperText,
    Button,
    ButtonGroup,
    Icon,
  } from "@chakra-ui/react";
  import { ChevronLeftIcon } from "@chakra-ui/icons";
  import { Routes } from "../../constants/Routes.enum";
  import { useRouter } from "next/router";
  import { useRef, useState } from "react";
  import axios from "axios";
import { Version } from "../../types/version";
import { Entry } from "../../types/entry";
import {toast} from "react-toastify"

  
  export default function CreateVersionComponent({ id, svcname, versions }: Entry) {
    const vNameRef = useRef<HTMLInputElement>();
    const [selectedFile, setSelectedFile] = useState(null);
    const fileRef = useRef<HTMLInputElement>(null);
  
    const handleFileChange = (event: any) => {
      setSelectedFile(event.target.files[0]);
    };
  
    const handleSubmit = async (e: any) => {
      e.preventDefault();
  
      const formElement = document.querySelector('#myForm') as HTMLFormElement | undefined;
      const formData = new FormData(formElement);
      formData.append('svcname', svcname)
  
      axios
      .post("http://localhost:3333/createVer", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      })
      .then((response) => {
        console.log(response);
        // Handle the upload success
        const toastId = toast.success(`Version added Successfully`);

        const redirectInterval = setInterval(() => {
          
        if (!toast.isActive(toastId)) {
            // If the toast message is no longer displayed, clear the interval and redirect to '/'
            clearInterval(redirectInterval);
              
            }
          }, 500);
        router.push(`/`);
      })
      .catch((error) => {
        console.log(error); //POSSIBLE ERRORS: NOT UNIQUE SVC NAME
        // Handle the upload error
      });
    };
  
    const router = useRouter();
    const handleBackClick = (e: any) => {
      e.preventDefault();
      router.push(Routes.HOME);
    };
  
    return (
      <>
        <BackgroundComponent></BackgroundComponent>
        <Container>
          <SubContainer>
            <FlexInputs>
              <TitleText>
                <MyIcon as={ChevronLeftIcon} onClick={() => router.back()} />
                New Version for {svcname}
              </TitleText>
              <form id="myForm">
                <MyFormControl>  
                  <div>
                    <FormLabel>Version Name</FormLabel>
                    <VersionInput
                      type="text"
                      name="vname"
                      ref={vNameRef}
                    ></VersionInput>
                    <FormHelperText>Your version name</FormHelperText>
                  </div>

                  <div>
                    <FormLabel>IDL File Upload</FormLabel>
                    <FileInput
                      type="file"
                      name="filetoupload"
                      onChange={handleFileChange}
                    ></FileInput>
                  </div>
                </MyFormControl>
                <Button type="submit" onClick={handleSubmit}>Submit</Button>
              </form>
            </FlexInputs>
          </SubContainer>
        </Container>
      </>
    );
  }
  