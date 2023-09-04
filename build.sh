# check if docker is running
if command docker info &> /dev/null; then
    echo "Docker is running/installed."
else
    echo "Docker is not running/installed.."
    exit 1
fi

if node --version &>/dev/null; then
    echo "Node.js is installed."
    npm install -g prisma
    echo "prisma installed"
else
    echo "Node.js is not installed."
    exit 1
fi


#check nginx installation
if [[ $OSTYPE == 'darwin'* ]]; then
    if nginx -v 2>&1 >/dev/null; then
        echo "nginx is installed."
    else
        echo "nginx is not installed."
        exit 1
    fi
fi

# Check nginx installation on Windows
if [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "msys" ]] || [[ "$(uname -s)" == "MINGW"* ]]; then
    echo "Please enter your nginx path:"
    read arg2
    if [[ ! "$arg2" =~ ^[A-Za-z]:/ ]]; then
        echo "Please provide an absolute path."
        exit 1
    fi
    cd "$arg2" || { echo "Failed to change directory to $arg2"; exit 1; }
    if ./nginx.exe -v 2>&1 >/dev/null; then
        echo "nginx is installed."
        # Here, you probably meant:
        cmd.exe /C "setx NGINX_PATH \"$arg2\" /M"
        cd -
    else
        echo "nginx is not installed/found"
        exit 1
    fi
fi



# Pull the etcd image from dockerhub
docker pull quay.io/coreos/etcd:v3.5.0

# Create an etcd-net network
docker network create --driver bridge etcd-net

# Start the docker container for etcd


# Pull the latest postgres image
docker pull postgres:latest
container_name="godzilla"


if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
    # Check if the container is not running
    if ! docker ps --format '{{.Names}}' | grep -q "^${container_name}$"; then
        # Start the existing container
        docker start ${container_name}
    else
        # Container is already running
        echo "Container ${container_name} is already running."
    fi
else
    # Run the container if it doesn't exist
    echo "Creating and starting the ${container_name} container..."
    docker run --name ${container_name} -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=mydatabase -p 5432:5432 postgres
fi




echo "postgres containers started."

echo "Move into the frontend folder"

# Navigate to the directory
cd IDL_management_page/page/my-app/ || { echo "Failed to change directory to frontend"; exit 1; }

# Ensure dependencies are up to date
npm install
npm install react react-dom next 
if [ $? -ne 0 ]; then
    echo "Error during npm install. Exiting."
    exit 1
fi

# Apply the migration using the schema to the database
npx prisma migrate dev --name init --schema=prisma/schema.prisma
if [ $? -ne 0 ]; then
    echo "Error during prisma migration. Exiting."
    exit 1
fi

# Generate a prisma client using the schema
npx prisma generate --schema=prisma/schema.prisma
if [ $? -ne 0 ]; then
    echo "Error during prisma generation. Exiting."
    exit 1
fi

# Run the user interface on port 3000

#move into backend folder``
cd -
cd pages-backend-test/pages-backend || { echo "Failed to change directory to the backend"; exit 1; }

npx prisma migrate dev --name init --schema=prisma/schema.prisma

npx prisma generate --schema=prisma/schema.prisma



