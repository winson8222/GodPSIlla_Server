import { Container, SubContainer, SubContainer3, SubContainer2, VerticalContainer, SubContainerStyle, TitleText} from "./styles"
import { Button, Spinner, FormControl, FormHelperText, FormLabel, Input, Select, HStack} from "@chakra-ui/react"
import BackgroundComponent from "../Background"
import MicrosvcDropdownComponent from "../MicrosvcDropdown"
import axios from "axios";
import { Entry } from "../../types/entry";
import React, { useEffect, useState } from "react";
import { useRouter } from "next/router";
import {Routes} from "../../constants/Routes.enum"
import { flattenDiagnosticMessageText } from "typescript";
import {toast} from "react-toastify"
import { DragControls } from "framer-motion";


export default function SelectServiceComponent(props: any) {
    
    const [lb, setLb] = useState("ROUND_ROBIN")
    const [url, setUrl] = useState("");

    const router = useRouter();
    const [loading, setLoading] = useState(false);
    const [isProcessingDone, setisProcessingDone] = useState(false)
    const [success, setsuccess] = useState('')
    const [generated, setgenerated] = useState(false)
    const [running, setRunning] = useState(false)
    const [stopping , setStopping] = useState(false)
    const [starting, setStarting] = useState(false)
    const [updating, setUpdating] = useState(false)
    const [deleting, setDeleting] = useState(false)
    const [cangen, setCangen] = useState(true)
    console.log(starting)
    const invalidPorts = [8888, 8889, 8890, 20000, 5432, 3333 ,3000];


    useEffect(() => {
      const generated = localStorage.getItem('generated');
      const running = localStorage.getItem('running')
      const port = localStorage.getItem("port")
      if (props.entries.length == 0 ) {
        setCangen(false)
      } else {
        setCangen(true)
      }

      if (generated === 'y') {
        setgenerated(true);
      } else {
        setgenerated(false)
      }

      if (running === 'y') {
        setRunning(true)
      } else {
        setRunning(false)
      }

      if (port === null) {
        setUrl("")
      } else {
        setUrl(port)
      }
    }, []);

  

    const handleStop = (e: any) => {

      e.preventDefault();
      setStopping(true);

      try {
            
        // Perform any necessary operations or API calls
        axios
        .post("http://localhost:3333/stop",)
        .then((response : any) => {
          if (response.data.outcome === 'Gateway Stopped') {
            localStorage.setItem("running", "n")
            setRunning(false)
            const toastId = toast.success(`Gateway Stopped Successfully`);

            const redirectInterval = setInterval(() => {
              
            if (!toast.isActive(toastId)) {
              
                clearInterval(redirectInterval);
                  
                }
              }, 500);
          }
        }).then(() => setStopping(false))
        .catch((error) => {
          console.log(error); 

          if (error.response) {
            toast.error('Stopping Gateway failed');
          }
          setStopping(false)
          
          
          // Handle the upload error
        })
        // After the operation is complete, navigate to the new version page
      } catch (error) {
        // Handle any errors that occur during the operation
        console.log(error)
      } finally {
        
      //   router.push(Routes.HOME);
      }

      console.log("stopped")
    }

    const handleStart = (e: any) => {
      e.preventDefault();
      setStarting(true);

      try {
            
        // Perform any necessary operations or API calls
        axios
        .post("http://localhost:3333/start")
        .then((response : any) => {
          if (response.data.outcome === 'Gateway Started') {
            localStorage.setItem("running", "y")
            setRunning(true)
            const toastId = toast.success(`Gateway Started Successfully`);

            const redirectInterval = setInterval(() => {
              
            if (!toast.isActive(toastId)) {
              
                clearInterval(redirectInterval);
                  
                }
              }, 500);
          
          }

          // Handle the upload success
        }).then(() => setStarting(false))
        .catch((error) => {
          console.log(error); 
          if (error.response) {
            toast.error('Gateway fails to Start, Try regenerate Gateway');
          }
          // Handle the upload error
        });

  
      } catch (error) {
        console.error(error);
        setStarting(false)
        // Handle any errors that occur during the operation
      }

      console.log("started")

    }

    const handleDel = (e: any) => {
      e.preventDefault();
      setDeleting(true);
      try {
        axios
        .post("http://localhost:3333/del")
        .then((response : any) => {
          if (response.data.outcome == 'Gateway Deleted') {
            const toastId = toast.success(`Gateway Deleted`);

            const redirectInterval = setInterval(() => {
              
            if (!toast.isActive(toastId)) {
              
                clearInterval(redirectInterval);
                  
                }
              }, 500);
          }
        }).then(() => {
          localStorage.setItem("generated", "n")
          localStorage.removeItem("port")
          setisProcessingDone(false)
          setgenerated(false);
          setDeleting(false)

        })
        .catch((error) => {
          console.log(error); 
          if (error.response) {
            toast.error('Deleting Gateway failed');
          }
          setDeleting(false)
          

        })
  
      } catch (error) {
        console.error(error);
        // Handle any errors that occur during the operation
      } finally {
        
      //   router.push(Routes.HOME);
      }
    }

    const handleUpdate = (e: any) => {
      if (!cangen) {
        toast.error("Please add a service")
      } else {
        console.log("updating")
      e.preventDefault();
      setUpdating(true);

      const formdata = {url, lb}

      try {
            
        // Perform any necessary operations or API calls
        axios
        .post("http://localhost:3333/update", formdata, {
          headers: {
            "Content-Type": "multipart/form-data",
          }})
        .then((response : any) => {
          if (response.data.outcome == 'Gateway Updated') {
            const toastId = toast.success(`Gateway Updated Successfully`);

            const redirectInterval = setInterval(() => {
              
            if (!toast.isActive(toastId)) {
              
                clearInterval(redirectInterval);
                  
                }
              }, 500);
          }
        }).then(() => setUpdating(false))
        .catch((error) => {
          console.log(error); 
          if (error.response) {
            toast.error('Updating Gateway failed');
          }
          setUpdating(false)
          

        })
  
      } catch (error) {
        console.error(error);
        // Handle any errors that occur during the operation
      } finally {
        
      //   router.push(Routes.HOME);
      }
      }
      
    }


    //Handler generation of gateway file
    const handleGen = (e: any) => {

      if(!cangen) {
        toast.error(`Please add a service`);
        return
      } else {
        if (url && lb) {
          const num = Number(url)
          if (invalidPorts.includes(num)) {
            toast.error(`Port is forbidden`);
  
          } else if (num < 0 || num > 65536) {
            toast.error('Please enter a port number between 0 and 65536.');
          } else {
            const id = toast.loading("Generating Gateway... Please do not leave this page")
            setisProcessingDone(true)
            e.preventDefault();
            setLoading(true);
            localStorage.setItem("port", url)
            
  
            const formData = { url , lb };
            
            
            try {
              // Perform any necessary operations or API calls
              axios
              .post("http://localhost:3333/gen", formData, {
                headers: {
                  "Content-Type": "multipart/form-data",
                }})
              .then((response : any) => {
                setLoading(false);
                console.log(response.data);
                if (response.data.outcome == 'Gateway Generated') {
                  localStorage.setItem("generated", "y")
                  localStorage.setItem("port", url)
                  setsuccess("Gateway Generated")
                  setgenerated(true);
                  toast.update(id, {render: "Gateway Generated", type: "success", isLoading: false, autoClose: 5000})
                }
                // Handle the upload success
              })
              .catch((error) => {
                console.log(error); //POSSIBLE ERRORS: NOT UNIQUE SVC NAME
                setLoading(false);
                setsuccess("Error")
                toast.update(id, {render: "Gateway Generation Error", type: "error", isLoading: false, autoClose: 5000})
                // Handle the upload error
              });
        
              // After the operation is complete, navigate to the new version page
            } catch (error) {
              console.error(error);
              toast.update(id, {render: "Gateway Generation Error", type: "error", isLoading: false, autoClose: 5000})
              // Handle any errors that occur during the operation
            }
          }
        } else {
          setsuccess("Fields not filled")
        }
      }
      
 
        
      };


  
      return (
        <>
          <BackgroundComponent />
          <Container
            style={{ display: "flex", flexDirection: "column", gap: "1rem" }}
          >
            <SubContainer style={{ height: "60%"}}>
              <MicrosvcDropdownComponent entries={props.entries} />
            </SubContainer>
            <SubContainer2 style={{ width: "30%", height: "60%", paddingBottom : "1%"}}>
              <FormControl>
                <div style={{ margin: "1rem" }}>
                  <FormLabel>Gateway Port</FormLabel>
                  <Input
                    type="number"
                    name="url"
                    value={url}
                    onChange={(e) => setUrl(e.target.value)}
                    style={{ width: "100%" }}
                    disabled={generated}
                    pattern="\d*"
                    required
                  />
    
                  <FormHelperText>Mandatory</FormHelperText>
                </div>
                <div style={{ margin: "1rem" }}>
                  <FormLabel>Load Balancing Type</FormLabel>
                  <Select
                    name="lb"
                    value={lb}
                    onChange={(e) => setLb(e.target.value)}
                    style={{ width: "100%" }}
                  >
                    <option value="ROUND_ROBIN">Round Robin</option>
                    <option value="">None</option>
                  </Select>
                  <FormHelperText>Optional</FormHelperText>
                </div>
              </FormControl>
              <HStack>
                <Button
                  onClick={handleGen}
                  size="md"
                  colorScheme="blue"
                  isLoading={loading}
                  loadingText="Creating..."
                  spinner={<Spinner size="md" />}
                  disabled={loading || !cangen}
                  style={{
                    opacity:  loading || generated ? 0.5 : 1,
                    cursor:
                    loading || generated
                        ? "not-allowed"
                        : "pointer",
                  }}
                  marginLeft="2rem"
                >
                  {isProcessingDone
                    ? success
                    : generated
                    ? "Gateway Generated"
                    : "Generate"}
                </Button>
                <Button
                  onClick={handleDel}
                  size="md"
                  colorScheme="red"
                  margin="0rem"
                  isLoading={deleting}
                  loadingText="Deleting..."
                  spinner={<Spinner size="md" />}
                  disabled={running || !generated}
                  style={{
                    opacity:
                      stopping || starting || updating || running || !generated
                        ? 0.5
                        : 1,
                    cursor:
                      stopping || starting || updating || running || !generated
                        ? "not-allowed"
                        : "pointer",
                  }}
                >
                  Delete
                </Button>
              </HStack>
            </SubContainer2>
    
            {generated && (
              <SubContainer3 style={{ height: "60%"}}>
                <div>
                  <TitleText margin="1rem">Control Panel</TitleText>
                </div>
                <Button
                  onClick={handleStop}
                  size="sm"
                  colorScheme="blue"
                  margin="0.5rem auto"
                  width="50%"
                  isLoading={stopping}
                  loadingText="Stopping"
                  spinner={<Spinner size="md" />}
                  disabled={starting || updating || stopping || !running}
                  style={{
                    opacity: stopping || starting || updating || !running ? 0.5 : 1,
                    cursor:
                      stopping || starting || updating || !running
                        ? "not-allowed"
                        : "pointer",
                  }}
                >
                  {" "}
                  Stop{" "}
                </Button>
                <Button
                  onClick={handleStart}
                  size="sm"
                  colorScheme="blue"
                  margin="0.5rem auto"
                  width="50%"
                  isLoading={starting}
                  loadingText="Starting"
                  spinner={<Spinner size="md" />}
                  disabled={starting || updating || stopping || running}
                  style={{
                    opacity: starting || stopping || updating || running ? 0.5 : 1,
                    cursor:
                      starting || stopping || updating || running
                        ? "not-allowed"
                        : "pointer",
                  }}
                >
                  {" "}
                  Start{" "}
                </Button>
                <Button
                  onClick={handleUpdate}
                  size="sm"
                  colorScheme="blue"
                  margin="0.5rem auto"
                  width="50%"
                  isLoading={updating}
                  loadingText="Updating"
                  spinner={<Spinner size="md" />}
                  disabled={starting || updating || stopping || !running || !cangen}
                  style={{
                    opacity: updating || stopping || starting || !running || !cangen ? 0.5 : 1,
                    cursor:
                      updating || stopping || starting || !running || !cangen
                        ? "not-allowed"
                        : "pointer",
                  }}
                >
                  {" "}
                  Update{" "}
                </Button>
              </SubContainer3>
            )}
          </Container>
        </>
    );
  }
 