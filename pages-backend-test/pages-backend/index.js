import { spawn } from "child_process";
import path from "path";
import { PrismaClient } from "@prisma/client";
import express from "express";
import pkg from "formidable";
const { formidable } = pkg;
import fs from "fs";
import cors from "cors";
import { fileURLToPath } from "url";
import os from "os";

const prisma = new PrismaClient();

const app = express();

const corsOptions = {
  origin: "http://localhost:3000", // Allow requests from a specific origin
  methods: ["GET", "POST"], // Allow specific HTTP methods
  allowedHeaders: ["Content-Type", "Authorization"], // Allow specific headers
};

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
app.use(cors(corsOptions));
app.use(express.json());

// Register and set up the middleware
app.use(express.urlencoded({ extended: true }));

//takes svcname, vname, upstreamurl and a file, creates a svc and a version
app.post("/createSvc", async (req, res) => {
  const form = new formidable.IncomingForm();
  // Parse `req` and upload all associated files

  form.parse(req, async function (err, fields, files) {
    if (err) {
      return res.status(400).json({ error: err.message });
    }
    const filepath = files.filetoupload.filepath; //filetoupload is the name of the file in the form in the req
    const { svcname, vname } = fields;
    const filename = files.filetoupload.originalFilename;
    try {
      const microsvc = await prisma.microservice.create({
        data: {
          svcname: svcname,
          versions: {
            create: [
              {
                vname: vname,
                idlfile: await convertFileToBytes(filepath),
                idlname: filename,
              },
            ],
          },
        },
        include: {
          versions: true,
        },
      });

      res.json({ yoursvc: microsvc });
      console.log("idlfile uploaded");
      res.status(200);
    } catch (err) {
      console.log(err);
      res.status(400).json({ error: err });
    }
  });
});

//takes a svcname and returns the svc object
app.post("/findSvc", async (req, res) => {
  const form = new formidable.IncomingForm();
  // Parse `req` and upload all associated files

  form.parse(req, async function (err, fields, files) {
    if (err) {
      return res.status(400).json({ error: err.message });
    }

    const { svcname } = fields;

    const microsvc = await prisma.microservice.findUnique({
      where: {
        svcname: svcname,
      },
    });

    res.json({ yoursvc: microsvc });

    res.status(200);
  });
});

app.post("/deleteSvc", async (req, res) => {
  const form = new formidable.IncomingForm();

  form.parse(req, async function (err, fields, files) {
    try {
      const { microserviceId } = fields;

      const microservice = await prisma.microservice.findUnique({
        where: {
          id: microserviceId,
        },
      });

      if (!microservice) {
        return res.status(404).json({ message: "Microservice not found." });
      }

      await prisma.microservice.delete({
        where: {
          id: microserviceId,
        },
      });

      return res
        .status(200)
        .json({ message: "Microservice deleted successfully." });
    } catch (err) {
      console.error(err);
      return res
        .status(500)
        .json({ error: "An error occurred while deleting the microservice." });
    }
  });
});

//takes microserviceId, vname, upstreamurl and a file, creates a version attached to that svc
app.post("/createVer", async (req, res) => {
  const form = new formidable.IncomingForm();
  // Parse `req` and upload all associated files

  form.parse(req, async function (err, fields, files) {
    if (err) {
      return res.status(400).json({ error: err.message });
    }
    const filepath = files.filetoupload.filepath; //filetoupload is the name of the file in the form in the req
    const { svcname, vname } = fields;
    const filename = files.filetoupload.originalFilename;
    const microsvc = await prisma.microservice.findUnique({
      where: {
        svcname: svcname,
      },
    });

    if (!microsvc) {
      console.error("Microservice not found");
      res.status(400).json({ error: "Microservice not found" });
      return;
    } else {
      // Create a new version and associate it with the microservice
      const version = await prisma.version.create({
        data: {
          vname: vname, // Replace with the actual version name
          idlfile: await convertFileToBytes(filepath),
          idlname: filename,
          microservice: { connect: { id: microsvc.id } },
        },
      });

      res
        .status(200)
        .json({ message: "Successfully created version for microsvc! " });
    }
  });
});

//returns all the svcs
app.get("/findAllSvc", async (req, res) => {
  try {
    const microsvcs = await prisma.microservice.findMany();

    res.status(200).json({ microsvcs: microsvcs });
  } catch (err) {
    res.status(400).json({ error: err });
  }
});

//returns all the versions of a svc
app.post("/findAllSvcVer", async (req, res) => {
  const form = new formidable.IncomingForm();
  // Parse `req` and upload all associated files

  form.parse(req, async function (err, fields, files) {
    if (err) {
      return res.status(400).json({ error: err.message });
    }

    const { svcname } = fields;

    const microsvc = await prisma.microservice.findUnique({
      where: {
        svcname: svcname,
      },
      include: {
        versions: true,
      },
    });

    if (!microsvc) {
      console.error("Microservice not found");
      res.status(400).json({ error: "Microservice not found" });
      return;
    } else {
      res.status(200).json({ microsvc: microsvc });
    }
  });
});

//returns all the svcs and their versions
app.get("/findAllInfo", async (req, res) => {
  try {
    const microsvcs = await prisma.microservice.findMany({
      include: {
        versions: true,
      },
    });

    res.status(200).json({ microsvcs: microsvcs });
  } catch (err) {
    res.status(400).json({ error: err });
  }
});

//returns all the svcs and their versions
app.post("/findSvcVerIDL", async (req, res) => {
  //req should provide vname and microserviceId
  try {
    const form = new formidable.IncomingForm();

    form.parse(req, async function (err, fields, files) {
      if (err) {
        return res.status(400).json({ error: err.message });
      }

      const { microserviceId, vname } = fields;

      const microservice = await prisma.microservice.findUnique({
        where: {
          id: microserviceId,
        },
      });

      const version = await prisma.version.findFirst({
        where: {
          microservice,
          vname: vname,
        },
        select: {
          idlfile: true,
        },
      });

      if (version) {
        const idlfile = Buffer.from(version.idlfile);
        res.status(200).json({ idlfile: idlfile });
        // await fs.promises.writeFile('output.thrift', idlfile);
        // Process the idlfile as needed
      } else {
        console.log("No version found for the microservice.");
      }
    });
  } catch (err) {
    res.status(400).json({ error: err });
  }
});

//Delete Specifc Version
app.post("/delSvcVer", async (req, res) => {
  try {
    const form = new formidable.IncomingForm();

    form.parse(req, async function (err, fields, files) {
      if (err) {
        return res.status(400).json({ error: err.message });
      }

      const { microserviceId, vname } = fields;

      const microservice = await prisma.microservice.findUnique({
        where: {
          id: microserviceId,
        },
      });

      const versions = await prisma.version.findMany({
        where: {
          microservice,
        },
      });

      const version = versions.find((version) => version.vname === vname);

      if (version) {
        await prisma.version.delete({
          where: {
            id: version.id,
          },
        });

        // if this was the only version for the microservice, send a request to the delete service route

        res
          .status(200)
          .json({
            message: "Version deleted successfully.",
            name: microservice.svcname,
          });
      } else {
        console.log("No version found for the microservice.");
        res
          .status(404)
          .json({ message: "No version found for the microservice." });
      }
    });
  } catch (err) {
    res.status(400).json({ error: err });
  }
});

//start gateway instances on ports 8888 8889 8890 with serverstart
//start nginx server with nstart
app.post("/start", async (req, res) => {
  try {
    let parentDirServerStart = path.join(
      __dirname,
      "..",
      "..",
      "/Gateway-Generator/serverstart"
    );

    let child = spawn("./serverstart", [8888, 8889, 8890], { cwd: parentDirServerStart });

    child.on("error", (error) => {
      console.error("Error starting child process:", error);
      // res.status(500).json({outcome: "Error Starting Child Process",error: error.message,});
      throw new Error(error.message);
    });

    child.stdout.on("data", (data) => {
      console.log(data.toString());
    });

    child.stderr.on("data", (data) => {
      console.error(data.toString());
    });

    child.on("close", (code) => {
      if (code == 0) {
        console.log("Servers Started successfully");
        // res.status(200).json({outcome: 'Gateway Started' }) cannot set headers yet
      } else {
        console.log("Error Starting Servers");
        // res.status(500).json({ outcome: "Error Starting Servers" });
        throw new Error("Error Starting Servers");
      }
    });

    //run nstart
    let parentDirNStart = path.join(
      __dirname,
      "..",
      "..",
      "/Gateway-Generator/nstart"
    );

    const osType = os.platform();

    console.log("os is " + osType);
    if (osType === "win32") {
      child = spawn("./nstart.exe", { cwd: parentDirNStart });
    } else {
      child = spawn("./nstart", { cwd: parentDirNStart });
    }

    child.on("error", (error) => {
      console.error("Error starting child process:", error);
      // res.status(500).json({outcome: "Error Starting Child Process",error: error.message,});
        throw new Error(error.message);
    });

    child.stdout.on("data", (data) => {
      console.log(data.toString());
    });

    child.stderr.on("data", (data) => {
      console.error(data.toString());
    });

    child.on("close", (code) => {
      if (code == 0) {
        console.log("Gateway Started successfully");
        res.status(200).json({ outcome: "Gateway Started" });
      } else {
        console.log("Error generating Gateway");
        // res.status(500).json({ outcome: "Error Starting Gateway" });
        throw new Error("Error Starting Gateway");
      }
    });
  } catch (err) {
    console.error("Error in route:", err);
    res.status(500).json({ outcome: "Error in route", error: err.message });
    return;
  }
});

//shutdown gateway instances on ports 8888 8889 8890 with shutdown
//shutdown nginx server with nstop
app.post("/stop", async (req, res) => {
  try {
    const parentDirShutDown = path.join(
      __dirname,
      "..",
      "..",
      "/Gateway-Generator/shutdown"
    );

    const osType = os.platform();

    let child;

    console.log("os is " + osType);
    if (osType === "win32") {
      child = spawn("./shutdown.exe", [8888, 8889, 8890], { cwd: parentDirShutDown });
    } else {
      child = spawn("./shutdown", [8888, 8889, 8890], { cwd: parentDirShutDown });
    }

    child.on("error", (error) => {
      console.error("Error stopping child process:", error);
      // res.status(500).json({outcome: "Error Stopping Child Process",error: error.message,});
      throw new Error(error.message)
    });

    child.stdout.on("data", (data) => {
      console.log(data.toString());
    });

    child.stderr.on("data", (data) => {
      console.error(data.toString());
    });

    child.on("close", (code) => {
      if (code == 0) {
        console.log("Servers Stopped");
        // res.status(200).json({outcome: 'Servers Stopped'}) cannot set headers yet
      } else {
        console.log("Error Stopping Servers");
        // res.status(500).json({ outcome: "Error Stopping Servers" });
        throw new Error("Error Stopping Servers")
      }
    });

    const parentDirNStop = path.join(
      __dirname,
      "..",
      "..",
      "/Gateway-Generator/nstop"
    );

    console.log("os is " + osType);
    if (osType === "win32") {
      child = spawn("./nstop.exe", { cwd: parentDirNStop });
    } else {
      child = spawn("./nstop", { cwd: parentDirNStop });
    }

    child.on("error", (error) => {
      console.error("Error stopping child process:", error);
      // res.status(500).json({outcome: "Error Stopping Child Process",error: error.message,});
      throw new Error(error.message)
    });

    child.stdout.on("data", (data) => {
      console.log(data.toString());
    });

    child.stderr.on("data", (data) => {
      console.error(data.toString());
    });

    child.on("close", (code) => {
      if (code == 0) {
        console.log("Gateway Stopped");
        res.status(200).json({ outcome: "Gateway Stopped" });
      } else {
        console.log("Error Stopping Gateway");
        // res.status(500).json({ outcome: "Error Stopping Gateway" });
        throw new Error("Error Stopping Gateway")
      }
    });
  } catch (err) {
    res.status(500).json({ outcome: "Error Stopping Gateway", error: err.message });
    return;
  }
});

//
app.post("/update", async (req, res) => {
  try {
    const form = new formidable.IncomingForm();
    form.parse(req, async function (err, fields) {
      if (err) {
        // res.status(400).json({ error: err.message });
        throw new Error(err.message)
      }
      const { url, lb } = fields;
      console.log("The url for update is: " + url)
      const parentDir = path.join(
        __dirname,
        "..",
        "..",
        "/Gateway-Generator/update"
      );


      const platform = os.platform()

      let child 
      console.log("os is " + platform)
      if (platform === "win32") {
        child = spawn("./update.exe", [url, lb], { cwd: parentDir });
      } else {
        child = spawn("./update", [url, lb], { cwd: parentDir });
      }
        

      child.on("error", (error) => {
        console.error("Error stopping child process:", error);
        // res.status(500).json({outcome: "Error Stopping Child Process",error: error.message,});
        throw new Error(error.message)
      });

      child.stdout.on("data", (data) => {
        console.log(data.toString());
      });

      child.stderr.on("data", (data) => {
        console.error(data.toString());
      });

      child.on("close", (code) => {
        if (code == 0) {
          console.log("Gateway Updated");
          res.status(200).json({ outcome: "Gateway Updated" });
        } else {
          console.log("Error Updating Gateway");
          // res.status(400).json({ outcome: "Error Updating Gateway" });
          throw new Error("Error Updating Gateway")
        }
      });
    });
  } catch (err) {
    console.log(err);
    res.status(500).json({ outcome: "Error Updating Gateway", error: err.message });
  }
});

//generate the gateway in Gateway-Generator/gateway using ./gen
app.post("/gen", async (req, res) => {
  try {
    const form = new formidable.IncomingForm();
    console.log(req.body.url);
    console.log(req.body.lb);

    form.parse(req, async function (err, fields) {

      if (err) {
        // res.status(400).json({ error: err.message });
        throw new Error(err.message)
      }
      const { url, lb } = fields;
      const platform = os.platform();

      const parentDir = path.join(__dirname, "..", "..", "/Gateway-Generator");
      let child;

      if (platform === "win32") {
        child = spawn("./gen.exe", [url, "start", lb], { cwd: parentDir });
      } else {
        child = spawn("./gen", [url, "start", lb], { cwd: parentDir });
      }
     

      child.on("error", (error) => {
        console.error("Error stopping child process:", error);
        // res.status(500).json({outcome: "Error Stopping Child Process",error: error.message,});
        throw new Error(error.message);
      });

      child.stdout.on("data", (data) => {
        console.log(data.toString());
      });

      child.stderr.on("data", (data) => {
        console.error(data.toString());
      });

      child.on("close", (code) => {
        if (code == 0) {
          console.log("Gateway generated successfully");
          res.status(200).json({ outcome: "Gateway Generated" });
        } else {
          console.log("Error generating Gateway");
          // res.status(400).json({ outcome: "Error Generating Gateway" });
          throw new Error("Error Generating Gateway");
        }
      });
    });
  } catch (error) {
    console.log(error)
    res.status(400).json({ outcome: "Error Generating Gateway", error: error.message });
  }
});

//delete can only be run when the gateway is stopped (i.e. nginx and gateway instances stopped)
//deletes the files in the gateway
app.post("/del", async (req, res) => {
  try {
    const parentDir = path.join(
      __dirname,
      "..",
      "..",
      "/Gateway-Generator/del"
    );
    const platform = os.platform();

    let child;

    if (platform === "win32") {
      child = spawn("./del.exe", { cwd: parentDir });
    } else {
      child = spawn("./del", { cwd: parentDir });
    }

    child.on("error", (error) => {
      console.error("Error stopping child process:", error);
      // res.status(500).json({outcome: "Error Stopping Child Process",error: error.message,});
      throw new Error(error.message)
    });

    child.stdout.on("data", (data) => {
      console.log(data.toString());
    });

    child.stderr.on("data", (data) => {
      console.error(data.toString());
    });

    child.on("close", (code) => {
      if (code == 0) {
        console.log("Gateway Deleted");
        return res.status(200).json({ outcome: "Gateway Deleted" });
      } else {
        console.log("Error Deleting Gateway");
        // res.status(400).json({ outcome: "Error Deleting Gateway" });
        throw new Error("Error Deleting Gateway");
      }
    });
  } catch (err) {
    console.log(err)
    return res.status(500).json({ outcome: "Error Deleting Gateway" });
  }
});

app.get("/", async (req, res) => {
  // const user = await User.findById(req.session!.userId).select('-password -__v -createdAt -updatedAt')

  // res.json(user)
  res.status(200).json({ message: "get req works!" });
});

const server = app.listen(3333, () => {
  console.log("Server is running on port 3333");
});
server.timeout = 300000;

async function convertFileToBytes(filePath) {
  return new Promise((resolve, reject) => {
    fs.readFile(filePath, (err, data) => {
      if (err) {
        reject(err);
      } else {
        resolve(data);
      }
    });
  });
}
