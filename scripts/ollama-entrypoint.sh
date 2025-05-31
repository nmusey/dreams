#!/bin/sh
set -e

echo "Starting Ollama service..."

export OLLAMA_HOST=0.0.0.0
export OLLAMA_ORIGINS=*

# Start Ollama in the background
/usr/bin/ollama serve &
OLLAMA_PID=$!

# Function to check if model exists
check_model() {
    ollama ls | grep -q "llava"
    return $?
}

# Pull the llava model if it doesn't exist
echo "Checking for llava model..."
if ! check_model; then
    echo "Pulling llava model..."
    /usr/bin/ollama pull llava
    
    # Verify the model was pulled successfully
    echo "Verifying llava model..."
    for i in $(seq 1 5); do
        if check_model; then
            echo "llava model successfully pulled and verified"
            break
        fi
        if [ $i -eq 5 ]; then
            echo "Failed to verify llava model after pulling"
            exit 1
        fi
        echo "Waiting for model to be available..."
        sleep 2
    done
else
    echo "llava model already exists"
fi

# Keep the container running and monitor the Ollama process
wait $OLLAMA_PID 
