mkdir -p models/ldm/stable-diffusion-v1/
cd models/ldm/stable-diffusion-v1/
wget https://huggingface.co/runwayml/stable-diffusion-v1-5/resolve/main/v1-5-pruned-emaonly.ckpt -O model.ckpt
wget https://huggingface.co/runwayml/stable-diffusion-v1-5/resolve/main/v1-inference.yaml -O model.yaml
