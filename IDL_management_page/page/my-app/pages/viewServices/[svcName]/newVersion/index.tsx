import axios from "axios";
import BackgroundComponent from "../../../../components/Background";
// import { Version } from "../../../../types/version";
// import ViewVersionComponent from "../../../../components/ViewVersion";
// import { Entry } from "../../../../types/entry";
import CreateVersionComponent from "../../../../components/NewVersion";
import { Version } from "../../../../types/version";
import { Entry } from "../../../../types/entry";
// import { entryInfoArray } from "../../../data/db";

function NewVerPage(props: any) {
  return (
    <>
      <BackgroundComponent></BackgroundComponent>
      <CreateVersionComponent 
    id={props.selectedEntry.id} svcname={props.selectedEntry.svcname} versions={props.selectedEntry.versions}/>
    </>
  );
}

//Get all versions of [svcName]
export async function getStaticPaths() {
  try {
    const response = await axios.get("http://localhost:3333/findAllInfo");

    const entries: Entry[] = response.data["microsvcs"]; // Assuming the response data is an array
    console.log(entries);

    const paths: any = [];
    entries.forEach((entry) => {
      paths.push({
        params: {
          svcName: entry.svcname,
        },
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

// //get the data for each version
export async function getStaticProps(context: any) {
  const actualSvcName = context.params.svcName;
  try {
    const response = await axios.get("http://localhost:3333/findAllInfo");

    const entries: Entry[] = response.data["microsvcs"]; // Assuming the response data is an array
    console.log(entries)
  
    let selectedEntry;
  
    for (let i = 0; i < entries.length; i++) {
      if (entries[i].svcname == actualSvcName) {
        selectedEntry = entries[i];
        console.log(selectedEntry)
      }
    }
  
    //fetch data for a single svc
    return {
      props: {
          selectedEntry,
      },
      revalidate: 10,
    };
  } catch(err) {
    console.log(err)
  }
}

export default NewVerPage;
