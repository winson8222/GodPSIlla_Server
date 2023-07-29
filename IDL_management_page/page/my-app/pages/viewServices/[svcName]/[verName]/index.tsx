import axios from "axios";
import BackgroundComponent from "../../../../components/Background";
import { Version } from "../../../../types/version";
import ViewVersionComponent from "../../../../components/ViewVersion";
import { Entry } from "../../../../types/entry";
// import { entryInfoArray } from "../../../data/db";


function VerPage(props: any) {
  return (
    <>
    <BackgroundComponent></BackgroundComponent>
    <ViewVersionComponent 
    id={props.selectedVersion.id} 
    vname={props.selectedVersion.vname} 
    upstreamurl={props.selectedVersion.upstreamurl} 
    idlfile={props.selectedVersion.idlfile}
    microserviceId={props.selectedVersion.microserviceId}/>
    </>
  );
}

//Get all versions of [svcName]
export async function getStaticPaths() {
  try {
    const response = await axios.get("http://localhost:3333/findAllInfo");

    const entries: Entry[] = response.data["microsvcs"]; // Assuming the response data is an array
    console.log(entries)

    const paths :any = []
    entries.forEach((entry) => {
        entry.versions.forEach((version) => {
          paths.push({
            params: {
              svcName: entry.svcname,
              verName: version.vname,
            },
          });
        });
      });

    return {
      fallback: false,
      paths: paths,
    };
  } catch (error) {
    console.log(error);
    // Handle the error
    return {
      fallback: false,
      paths: [],
    };
  }
}

//get the data for each version
export async function getStaticProps(context: any) {
  try {
    const actualSvcName = context.params.svcName
    const actualVerName = context.params.verName
    console.log(actualVerName)
    const formData = new FormData()
    formData.append('svcname', actualSvcName)
    const response = await axios.post("http://localhost:3333/findAllSvcVer", formData, {
      headers: {
        "Content-Type": "multipart/form-data",
      },
    })

    console.log(response)

    const entries: Version[] = response.data["microsvc"]["versions"]; // Assuming the response data is an array
  
    let selectedVersion;
  
    for (let i = 0; i < entries.length; i++) {
      if (entries[i].vname == actualVerName) {
        selectedVersion = entries[i];
        console.log(selectedVersion)
      }
    }
    
  
    //fetch data for a single ver
    return {
      props: {
          selectedVersion,
      },
      revalidate: 5,
    };
  } catch(err) {
    console.log(err)
  }
}

export default VerPage;