#!/bin/sh

# Start Ollama in the background
ollama serve &

# Wait for Ollama to be ready (port 11434)
while ! nc -z localhost 11434 2>/dev/null; do
    echo "Waiting for Ollama to start..."
    sleep 1
done

echo "Ollama is ready. Pulling llava model..."
ollama pull llava

# Keep the container running
wait 