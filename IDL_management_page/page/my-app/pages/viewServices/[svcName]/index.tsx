import axios from "axios";
import BackgroundComponent from "../../../components/Background";
import ViewServiceComponent from "../../../components/ViewService";
import { Entry } from "../../../types/entry";
// import { entryInfoArray } from "../../../data/db";


function SvcPage(props: any) {
  return (
    <>
    <BackgroundComponent></BackgroundComponent>
    <ViewServiceComponent id={props.selectedEntry.id} svcname={props.selectedEntry.svcname} versions={props.selectedEntry.versions}/>
    </>
  );
}

//Must be used with getStaticProps to define for which values of [meetupId] are pre-rendered pages
export async function getStaticPaths() {
  //Must be used with getStaticProps to define for which values of [meetupId] are pre-rendered pages
  try {
    const response = await axios.get("http://localhost:3333/findAllSvc");
    const entries = response.data; // Assuming the response data is an array
    // console.log(entries)

    return {
      fallback: false,
      paths: entries["microsvcs"].map((entry: any) => ({
        params: { svcName: entry.svcname },
      })),
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

//runs at page build to generate pre-rendered pages on the server
//Will pre-render all possible pages of /[meetupId] in advance
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
export default SvcPage;