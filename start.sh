
if docker-compose ps | grep -q "Up"; then
    echo "etcd services are up and running!"
else
    docker-compose up > docker-compose.log 2>&1 &
    # Give it a few seconds to ensure services start
    sleep 5
fi



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

# Check if the operating system is Windows (under WSL)
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo "Running on Windows..."

    cd pages-backend-test/pages-backend && node index.js &
    echo backend stated

    # Return to the original directory and then navigate to IDL_management_page/page/my-app
    cd -
    cd IDL_management_page/page/my-app

    # Run npm run dev in the foreground
    npm run dev
    echo frontend stated
fi
exit