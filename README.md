![Getting Started](godzilla.png)

<div style="border: 1px solid white; padding: 10px;">

IDL_management_page contains the user interface

Gateway-Generator contains the Code generator

pages-backend-test contains the backend server for interacting with the database

</div>

> Requirements: 
> [Docker Desktop](https://www.docker.com/products/docker-desktop/)
> [Command Tool of Hertz](https://www.cloudwego.io/docs/hertz/getting-started/#install-the-command-tool-of-hz) installed
> [Node.js](https://nodejs.org/en)
> [Nginx](https://www.nginx.com/)
> MacOS or Windows OS

For first time set up:
- Run build.sh
- Type "yes" to reset schema if prompted
- Indicate nginx path when promoted (For windows)

To run:
- Run start.sh

Possible Issues:
- If start.sh is unable to start up the IDL management page, try navigate to IDL_management_page/page/my-app and run "npm i"