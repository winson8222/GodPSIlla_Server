#!/bin/bash

# Check if Docker containers are running
if docker-compose ps | grep -q "Up"; then
    echo "etcd services are up and running!"
else
    docker-compose up > docker-compose.log 2>&1 &
    # Give it a few seconds to ensure services start
    sleep 5
fi

docker start godzilla

# Check if the operating system is macOS or Linux (assuming you're using this on macOS, but it will also work on Linux)
if [[ "$OSTYPE" == "darwin"* ]]; then
    echo "Running on macOS..."

    terminate_processes() {
        echo "Terminating processes..."
        for pid in "${pids[@]}"; do
            if ps -p $pid > /dev/null; then
                kill -TERM $pid
            fi
        done
    }

    trap terminate_processes SIGINT SIGTERM

    cd pages-backend-test/pages-backend
    node index.js &
    echo backend started

    cd -
    cd IDL_management_page/page/my-app
    npm run dev
    echo frontend started

# Rest of your script
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    stop_background_processes_windows() {
        echo "Stopping background processes..."
        
        taskkill /F /PID $backend_pid
        taskkill /F /PID $frontend_pid
        
        wait "$backend_pid" "$frontend_pid" 2>/dev/null
    }

    trap 'stop_background_processes_windows; exit' SIGINT
    echo "Running on Windows..."

    cd pages-backend-test/pages-backend && node index.js &
    backend_pid=$!
    echo "Backend started (PID: $backend_pid)"

    # Return to the original directory and then navigate to IDL_management_page/page/my-app
    cd -
    cd IDL_management_page/page/my-app

    # Run npm run dev in the foreground
    npm run dev &
    frontend_pid=$!
    echo "Frontend started (PID: $frontend_pid)"

    # Keep the script running until terminated
    wait
fi
